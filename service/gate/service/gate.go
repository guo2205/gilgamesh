// gate
package service

import (
	"fractal/fractal"
	"gilgamesh/utility/config"
	"gilgamesh/utility/mylog"
)

type GateService struct {
	logger *mylog.Logger
	f      *fractal.Fractal
	option *config.GateOption
	writer func(session uint64, d []byte) error
}

func NewGateService(logger *mylog.Logger,
	f *fractal.Fractal,
	option *config.GateOption,
	writer func(session uint64, d []byte) error) *GateService {
	return &GateService{
		logger: logger,
		f:      f,
		option: option,
		writer: writer,
	}
}

func (c *GateService) OnStart() error {
	return nil
}

func (c *GateService) OnMail(caller string, session uint64, data []byte) ([]byte, error) {
	c.logger.Debug(caller, session, data)
	return []byte{}, nil
}

func (c *GateService) OnClose() error {
	return nil
}
