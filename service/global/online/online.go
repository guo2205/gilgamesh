// auth
package online

import (
	"fractal/fractal"
	"gilgamesh/utility/mylog"
)

type OnlineService struct {
	fractal.DefaultServiceProvider
	logger *mylog.Logger
	f      *fractal.Fractal
}

func NewService(
	logger *mylog.Logger,
	f *fractal.Fractal) *OnlineService {
	return &OnlineService{
		logger: logger,
		f:      f,
	}
}

func (c *OnlineService) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	c.logger.Debug(caller, _type, session, data)
	return []byte{}, nil
}
