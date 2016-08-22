// auth
package locker

import (
	"gilgamesh/protos"

	"github.com/liuhanlcj/mylog"
)

type _Locker struct {
	Key     string
	Session uint64
}

type Service struct {
	logger mylog.Logger

	globalLock map[string]*_Locker
}

func NewService(
	logger mylog.Logger) *Service {
	return &Service{
		logger:     logger,
		globalLock: make(map[string]*_Locker, 10000),
	}
}

func (c *Service) On_AcquireRequest(caller string, session uint64, in *protos.Internal_Global_Locker_AcquireRequest, responser func(out *protos.Internal_Global_Response, e error)) {
	_, ok := c.globalLock[in.Key]
	if !ok {
		c.logger.Infof("[%s %d] locked\n", in.Key, session)

		c.globalLock[in.Key] = &_Locker{
			Key:     in.Key,
			Session: session,
		}
		responser(&protos.Internal_Global_Response{
			Success: true,
		}, nil)
		return
	}

	responser(&protos.Internal_Global_Response{
		Success: false,
	}, nil)
}

func (c *Service) On_ReleaseRequest(caller string, session uint64, in *protos.Internal_Global_Locker_ReleaseRequest, responser func(out *protos.Internal_Global_Response, e error)) {
	locker, ok := c.globalLock[in.Key]
	if !ok {
		responser(&protos.Internal_Global_Response{
			Success: false,
		}, nil)
		return
	}

	if locker.Session != session {
		responser(&protos.Internal_Global_Response{
			Success: false,
		}, nil)
		return
	}

	delete(c.globalLock, in.Key)

	responser(&protos.Internal_Global_Response{
		Success: true,
	}, nil)

	c.logger.Infof("[%s %d] unlocked\n", in.Key, session)
}
