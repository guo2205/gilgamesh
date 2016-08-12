// auth
package locker

import (
	"fractal/fractal"
	"gilgamesh/protos"
	"gilgamesh/utility/mylog"
	"gilgamesh/utility/utils"

	"github.com/golang/protobuf/proto"
)

type LockerService struct {
	fractal.DefaultServiceProvider
	logger *mylog.Logger
	f      *fractal.Fractal

	globalLock map[string]bool
}

func NewService(
	logger *mylog.Logger,
	f *fractal.Fractal) *LockerService {
	return &LockerService{
		logger:     logger,
		f:          f,
		globalLock: make(map[string]bool, 10000),
	}
}

func (c *LockerService) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	obj, ptype, err := utils.Unmarshal(data)
	if err != nil {
		return []byte{}, err
	}

	switch ptype {
	case proto.MessageName((*protos.Internal_AcquireLocker)(nil)):
		return c.do_Internal_AcquireLocker(session, obj.(*protos.Internal_AcquireLocker))
	}
	return []byte{}, nil
}

func (c *LockerService) do_Internal_AcquireLocker(session uint64, obj *protos.Internal_AcquireLocker) ([]byte, error) {
	if obj.Lock {
		state, ok := c.globalLock[obj.Account]
		if !ok {
			c.globalLock[obj.Account] = true
			return []byte("success"), nil
		}

		if state {
			return []byte("failed"), nil
		} else {
			c.globalLock[obj.Account] = true
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
			return []byte("success"), nil
		}
	}

}
