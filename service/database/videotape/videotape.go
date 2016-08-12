// videotape
package videotape

import (
	"fractal/fractal"
	"gilgamesh/utility/mylog"
)

type Service struct {
	fractal.DefaultServiceProvider
	logger *mylog.Logger
	f      *fractal.Fractal
}

func NewService(
	logger *mylog.Logger,
	f *fractal.Fractal) *Service {
	return &Service{
		logger: logger,
		f:      f,
	}
}

func (c *Service) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	c.logger.Debug(caller, _type, session, data)
	return []byte{}, nil
}
