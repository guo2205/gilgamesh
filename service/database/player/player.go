// player
package player

import (
	"fractal/fractal"
	"gilgamesh/protos"
	"gilgamesh/utility/utils"

	"github.com/golang/protobuf/proto"
	"github.com/liuhanlcj/mylog"
)

type Service struct {
	fractal.DefaultServiceProvider
	logger mylog.Logger
	f      *fractal.Fractal
}

func NewService(
	logger mylog.Logger,
	f *fractal.Fractal) *Service {
	return &Service{
		logger: logger,
		f:      f,
	}
}

func (c *Service) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	obj, ptype, err := utils.Unmarshal(data)
	if err != nil {
		return []byte{}, err
	}

	switch ptype {
	case proto.MessageName((*protos.Public_Cts_Player_Query)(nil)):
		return c.do_Public_Cts_Player_Query(session, obj.(*protos.Public_Cts_Player_Query))
	case proto.MessageName((*protos.Internal_Database_Player_Create)(nil)):
		return c.do_Internal_Database_Player_Create(session, obj.(*protos.Internal_Database_Player_Create))
	case proto.MessageName((*protos.Internal_Database_Player_Modify)(nil)):
		return c.do_Internal_Database_Player_Modify(session, obj.(*protos.Internal_Database_Player_Modify))
	case proto.MessageName((*protos.Internal_Database_Player_Query)(nil)):
		return c.do_Internal_Database_Player_Query(session, obj.(*protos.Internal_Database_Player_Query))
	}
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Player_Query(session uint64, obj *protos.Public_Cts_Player_Query) ([]byte, error) {
	players, err := GetAccountsPlayer(obj.AccountList)
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Player_QueryResponse{
			PlayerList: []*protos.Public_Player{},
		}), nil
	}

	return utils.Marshal(&protos.Public_Stc_Player_QueryResponse{
		PlayerList: players,
	}), nil
}

func (c *Service) do_Internal_Database_Player_Create(session uint64, obj *protos.Internal_Database_Player_Create) ([]byte, error) {
	ok, err := CreateAccountPlayer(obj.Account, obj.Player)
	if err != nil || !ok {
		return utils.Marshal(&protos.Public_Stc_Player_CreateResponse{
			State: false,
		}), nil
	}

	return utils.Marshal(&protos.Public_Stc_Player_CreateResponse{
		State: true,
	}), nil
}

func (c *Service) do_Internal_Database_Player_Modify(session uint64, obj *protos.Internal_Database_Player_Modify) ([]byte, error) {
	ok, err := ModifyAccountPlayer(obj.Account, obj.Player)
	if err != nil || !ok {
		return utils.Marshal(&protos.Public_Stc_Player_ModifyResponse{
			State: false,
		}), nil
	}

	return utils.Marshal(&protos.Public_Stc_Player_ModifyResponse{
		State: true,
	}), nil
}

func (c *Service) do_Internal_Database_Player_Query(session uint64, obj *protos.Internal_Database_Player_Query) ([]byte, error) {
	py, ok, err := GetAccountPlayer(obj.Account)
	if err != nil || !ok {
		return []byte{}, nil
	}
	return utils.Marshal(&protos.Public_Stc_Player_QueryResponse{
		PlayerList: []*protos.Public_Player{py},
	}), nil
}
