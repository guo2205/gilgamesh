// auth
package locker

import (
	"fractal/fractal"
	"gilgamesh/utility/mylog"
)

type LockerService struct {
	fractal.DefaultServiceProvider
	logger *mylog.Logger
	f      *fractal.Fractal
}

func NewService(
	logger *mylog.Logger,
	f *fractal.Fractal) *LockerService {
	return &LockerService{
		logger: logger,
		f:      f,
	}
}

func (c *LockerService) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	c.logger.Debug(caller, _type, session, data)
	return []byte{}, nil
}
