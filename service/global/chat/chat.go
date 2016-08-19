// chat
package chat

import (
	"fractal/fractal"
	"gilgamesh/protos"
	"gilgamesh/utility/utils"

	"github.com/golang/protobuf/proto"
	"github.com/liuhanlcj/mylog"
)

type _AccountPosition struct {
	Where   string
	Session uint64
}

type Service struct {
	fractal.DefaultServiceProvider
	logger mylog.Logger
	f      *fractal.Fractal

	accountMap map[string]*_AccountPosition
}

func NewService(
	logger mylog.Logger,
	f *fractal.Fractal) *Service {
	return &Service{
		logger:     logger,
		f:          f,
		accountMap: make(map[string]*_AccountPosition, 1000),
	}
}

func (c *Service) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	obj, ptype, err := utils.Unmarshal(data)
	if err != nil {
		return []byte{}, err
	}

	switch ptype {
	case proto.MessageName((*protos.Public_Chat)(nil)):
		return c.do_Public_Chat(session, data, obj.(*protos.Public_Chat))
	case proto.MessageName((*protos.Internal_Global_Online_Set)(nil)):
		return c.do_Internal_SetOnline(caller, session, obj.(*protos.Internal_Global_Online_Set))
	}
	return []byte{}, nil
}

func (c *Service) do_Public_Chat(session uint64, data []byte, obj *protos.Public_Chat) ([]byte, error) {
	for _, acc := range c.accountMap {
		c.f.PostMail(acc.Where, 0, "chat", acc.Session, data)
	}
	return []byte{}, nil
}

func (c *Service) do_Internal_SetOnline(caller string, session uint64, obj *protos.Internal_Global_Online_Set) ([]byte, error) {
	if obj.State {
		c.accountMap[obj.Account] = &_AccountPosition{
			Where:   caller,
			Session: session,
		}
		c.logger.Info("account ", obj.Account, "enter chat")
	} else {
		delete(c.accountMap, obj.Account)
		c.logger.Info("account ", obj.Account, "exit chat")
	}
	return []byte{}, nil
}
