// gate
package service

import (
	"errors"
	"gilgamesh/protos"
	"gilgamesh/utility/config"
	"gilgamesh/utility/utils"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/liuhanlcj/fractal/fractal/sdk"
	"github.com/liuhanlcj/mylog"
)

type _Client struct {
	Session uint64
	Account string
}

type Service struct {
	logger mylog.Logger
	f      *fsdk.Fractal
	option *config.GateOption
	writer func(session uint64, d []byte) error
	closer func(session uint64) error

	clients map[uint64]*_Client
}

var (
	ErrUnknownProtoType error = errors.New("unknown proto type")
)

func NewService(logger mylog.Logger,
	f *fsdk.Fractal,
	option *config.GateOption,
	writer func(session uint64, d []byte) error,
	closer func(session uint64) error) *Service {
	return &Service{
		logger:  logger,
		f:       f,
		option:  option,
		writer:  writer,
		closer:  closer,
		clients: make(map[uint64]*_Client, 100),
	}
}

func (c *Service) On_EntryDataRequest(caller string, session uint64, in *protos.Internal_GameGate_EntryDataRequest, responser func(e error)) {
	obj, family, err := utils.Unmarshal(in.Data)
	if err != nil {
		c.logger.Debug("unmarshal failed:", err)
		responser(err)
		return
	}

	switch family {
	case proto.MessageName((*protos.Auth_AuthRequest)(nil)):
		c.do_Auth_AuthRequest(session, obj.(*protos.Auth_AuthRequest), responser)

	case proto.MessageName((*protos.Chat_HallRequest)(nil)):
		c.do_Chat_HallRequest(session, obj.(*protos.Chat_HallRequest), responser)

	case proto.MessageName((*protos.Hall_EnterRoomRequest)(nil)):
		c.do_Hall_EnterRoomRequest(session, obj.(*protos.Hall_EnterRoomRequest), responser)
	case proto.MessageName((*protos.Hall_CreateRoomRequest)(nil)):
		c.do_Hall_CreateRoomRequest(session, obj.(*protos.Hall_CreateRoomRequest), responser)

	case proto.MessageName((*protos.Duel_DataTransfer)(nil)):
		c.do_Duel_DataTransfer(session, obj.(*protos.Duel_DataTransfer), responser)

	default:
		responser(ErrUnknownProtoType)
	}

	return
}

func (c *Service) On_PassThroughRequest(caller string, session uint64, in *protos.Internal_GameGate_PassThroughRequest, responser func(e error)) {
	responser(c.writer(session, in.Data))
}

func (c *Service) On_KickRequest(caller string, session uint64, in *protos.Internal_GameGate_KickRequest, responser func(e error)) {
	responser(c.closer(in.Session))
}

func (c *Service) On_MyOfflineRequest(caller string, session uint64, responser func(e error)) {
	client, ok := c.clients[session]
	if ok {
		// 通知大厅用户离开
		protos.New_HallService_ServiceClient(c.f, "/ygo/hall.hall").Call_LeaveRequest_Cps("gate", session, time.Second*3, func(out *protos.Internal_Hall_LeaveResponse, to string, e error) {
			if e != nil {
				c.logger.Warningf("[%s %d] leave hall failed : %v\n", client.Account, client.Session, e)
			} else if !out.Success {
				c.logger.Warningf("[%s %d] leave hall failed : unknown reason\n", client.Account, client.Session)
			}

			// 设置用户离线
			protos.New_GlobalOnlineStateService_ServiceClient(c.f, "/public/global.online").Call_SetRequest_Cps("gate", session, time.Second*3, &protos.Internal_Global_OnlineState_SetRequest{
				Account: client.Account,
				State:   false,
			}, func(out *protos.Internal_Global_Response, to string, e error) {
				if e != nil {
					c.logger.Warningf("[%s %d] set offline failed : %v\n", client.Account, client.Session, e)
				} else if !out.Success {
					c.logger.Warningf("[%s %d] set offline failed : unknown reason\n", client.Account, client.Session)
				}

				// 清理
				c.f.PostEvent("gate", func() {
					delete(c.clients, session)
					responser(nil)
				})
			})
		})
	}
}
