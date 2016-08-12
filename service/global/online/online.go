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

type OnlineService struct {
	fractal.DefaultServiceProvider
	logger *mylog.Logger
	f      *fractal.Fractal

	accountStateMap map[string]*_AccountPosition
}

func NewService(
	logger *mylog.Logger,
	f *fractal.Fractal) *OnlineService {
	return &OnlineService{
		logger:          logger,
		f:               f,
		accountStateMap: make(map[string]*_AccountPosition, 1000),
	}
}

func (c *OnlineService) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	obj, ptype, err := utils.Unmarshal(data)
	if err != nil {
		return []byte{}, err
	}

	switch ptype {
	case proto.MessageName((*protos.Internal_QueryOnline)(nil)):
		return c.do_Internal_QueryOnline(session, obj.(*protos.Internal_QueryOnline))
	case proto.MessageName((*protos.Internal_SetOnline)(nil)):
		return c.do_Internal_SetOnline(caller, session, obj.(*protos.Internal_SetOnline))
	}
	return []byte{}, nil
}

func (c *OnlineService) do_Internal_QueryOnline(session uint64, obj *protos.Internal_QueryOnline) ([]byte, error) {
	response := protos.Internal_QueryOnlineResponse{}

	state, ok := c.accountStateMap[obj.Account]
	if !ok {
		response.State = false
	} else {
		response.State = true
		response.Where = state.Where
		response.Session = state.Session
	}
	d, _ := proto.Marshal(&response)
	return d, nil
}

func (c *OnlineService) do_Internal_SetOnline(caller string, session uint64, obj *protos.Internal_SetOnline) ([]byte, error) {
	if obj.State {
		c.accountStateMap[obj.Account] = &_AccountPosition{
			State:   true,
			Where:   caller,
			Session: session,
		}
	} else {
		delete(c.accountStateMap, obj.Account)
	}
	return []byte{}, nil
}
