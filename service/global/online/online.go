// auth
package online

import (
	"fractal/fractal"
	"gilgamesh/protos"
	"gilgamesh/utility/mylog"
	"gilgamesh/utility/utils"

	"github.com/golang/protobuf/proto"
)

type _AccountPosition struct {
	State   bool
	Where   string
	Session uint64
}

type Service struct {
	fractal.DefaultServiceProvider
	logger *mylog.Logger
	f      *fractal.Fractal

	accountStateMap map[string]*_AccountPosition
}

func NewService(
	logger *mylog.Logger,
	f *fractal.Fractal) *Service {
	return &Service{
		logger:          logger,
		f:               f,
		accountStateMap: make(map[string]*_AccountPosition, 1000),
	}
}

func (c *Service) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	obj, ptype, err := utils.Unmarshal(data)
	if err != nil {
		return []byte{}, err
	}

	switch ptype {
	case proto.MessageName((*protos.Internal_Global_Online_Query)(nil)):
		return c.do_Internal_QueryOnline(session, obj.(*protos.Internal_Global_Online_Query))
	case proto.MessageName((*protos.Internal_Global_Online_Set)(nil)):
		return c.do_Internal_SetOnline(caller, session, obj.(*protos.Internal_Global_Online_Set))
	}
	return []byte{}, nil
}

func (c *Service) do_Internal_QueryOnline(session uint64, obj *protos.Internal_Global_Online_Query) ([]byte, error) {
	response := protos.Internal_Global_Online_QueryResponse{}

	state, ok := c.accountStateMap[obj.Account]
	if !ok {
		response.State = false
	} else {
		response.State = true
		response.Where = state.Where
		response.Session = state.Session
	}
	return utils.Marshal(&response), nil
}

func (c *Service) do_Internal_SetOnline(caller string, session uint64, obj *protos.Internal_Global_Online_Set) ([]byte, error) {
	if obj.State {
		c.accountStateMap[obj.Account] = &_AccountPosition{
			State:   true,
			Where:   caller,
			Session: session,
		}
		c.logger.Info("account ", obj.Account, "online")
	} else {
		delete(c.accountStateMap, obj.Account)
		c.logger.Info("account ", obj.Account, "offline")
	}
	return []byte{}, nil
}
