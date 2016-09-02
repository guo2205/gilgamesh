// auth
package chat

import (
	"github.com/liuhanlcj/fractal/fractal/sdk"
	"github.com/liuhanlcj/gilgamesh/protos"
	"github.com/liuhanlcj/gilgamesh/utility/utils"
	"github.com/liuhanlcj/mylog"
)

type Service struct {
	logger mylog.Logger
	f      *fsdk.Fractal

	visitors map[uint64]string
}

func NewService(logger mylog.Logger, f *fsdk.Fractal) *Service {
	return &Service{
		logger:   logger,
		f:        f,
		visitors: make(map[uint64]string, 1000),
	}
}

func (c *Service) On_EnterHall(caller string, session uint64, responser func(e error)) {
	c.visitors[session] = caller
	responser(nil)

	c.logger.Infof("[%d %s] enter chat\n", session, caller)
}

func (c *Service) On_LeaveHall(caller string, session uint64, responser func(e error)) {
	delete(c.visitors, session)
	responser(nil)

	c.logger.Infof("[%d] leave chat\n", session)
}

func (c *Service) On_HallChat(caller string, session uint64, in *protos.Internal_Chat_HallRequest, responser func(e error)) {
	d := utils.Marshal(&protos.Chat_Hall{
		Account: in.Account,
		Content: in.Content,
	})
	for session, where := range c.visitors {
		protos.New_GameGateService_ServiceClient(c.f, where).Call_PassThrough("chat", session, &protos.Internal_GameGate_PassThroughRequest{
			Data: d,
		})
	}
	responser(nil)
}
