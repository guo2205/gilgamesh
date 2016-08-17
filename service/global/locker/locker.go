// auth
package locker

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

	globalLock map[string]bool
}

func NewService(
	logger mylog.Logger,
	f *fractal.Fractal) *Service {
	return &Service{
		logger:     logger,
		f:          f,
		globalLock: make(map[string]bool, 10000),
	}
}

func (c *Service) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	obj, ptype, err := utils.Unmarshal(data)
	if err != nil {
		return []byte{}, err
	}

	switch ptype {
	case proto.MessageName((*protos.Internal_Global_Locker_Acquire)(nil)):
		return c.do_Internal_AcquireLocker(session, obj.(*protos.Internal_Global_Locker_Acquire))
	}
	return []byte{}, nil
}

func (c *Service) do_Internal_AcquireLocker(session uint64, obj *protos.Internal_Global_Locker_Acquire) ([]byte, error) {
	c.logger.Debug("account ", obj.Account, "lock req :", obj.Lock)

	if obj.Lock {
		state, ok := c.globalLock[obj.Account]
		if !ok {
			c.globalLock[obj.Account] = true

			c.logger.Info("account ", obj.Account, "lock")

			return []byte("success"), nil
		}

		if state {
			return []byte("failed"), nil
		} else {
			c.globalLock[obj.Account] = true

			c.logger.Info("account ", obj.Account, "lock")

			return []byte("success"), nil
		}
	} else {
		state, ok := c.globalLock[obj.Account]
		if !ok {
			return []byte("failed"), nil
		}
		if !state {
			return []byte("failed"), nil
		} else {
			c.globalLock[obj.Account] = false

			c.logger.Info("account ", obj.Account, "unlock")

			return []byte("success"), nil
		}
	}

}
