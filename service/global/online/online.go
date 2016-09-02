// auth
package online

import (
	"github.com/liuhanlcj/gilgamesh/protos"
	"github.com/liuhanlcj/mylog"
)

type _AccountPosition struct {
	State   bool
	Where   string
	Session uint64
}

type Service struct {
	logger mylog.Logger

	accountStateMap map[string]*_AccountPosition
}

func NewService(
	logger mylog.Logger) *Service {
	return &Service{
		logger:          logger,
		accountStateMap: make(map[string]*_AccountPosition, 1000),
	}
}

func (c *Service) On_Query(caller string, session uint64, in *protos.Internal_Global_OnlineState_QueryRequest, responser func(out *protos.Internal_Global_OnlineState_QueryResponse, e error)) {
	response := protos.Internal_Global_OnlineState_QueryResponse{}

	state, ok := c.accountStateMap[in.Account]
	if !ok {
		response.State = false
	} else {
		response.State = true
		response.Where = state.Where
		response.Session = state.Session
	}
	responser(&response, nil)
}

func (c *Service) On_Set(caller string, session uint64, in *protos.Internal_Global_OnlineState_SetRequest, responser func(out *protos.Internal_Global_Response, e error)) {
	if in.State {
		_, ok := c.accountStateMap[in.Account]
		if ok {
			responser(&protos.Internal_Global_Response{
				Success: false,
			}, nil)
			return
		}

		c.accountStateMap[in.Account] = &_AccountPosition{
			State:   true,
			Where:   caller,
			Session: session,
		}
		responser(&protos.Internal_Global_Response{
			Success: true,
		}, nil)

		c.logger.Infof("[%s %d] online\n", in.Account, session)
	} else {
		delete(c.accountStateMap, in.Account)
		responser(&protos.Internal_Global_Response{
			Success: true,
		}, nil)

		c.logger.Infof("[%s %d] offline\n", in.Account, session)
	}
}
