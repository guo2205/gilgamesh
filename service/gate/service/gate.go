// gate
package service

import (
	"errors"
	"fractal/fractal"
	"gilgamesh/utility/config"
	"gilgamesh/utility/mylog"
)

var (
	ErrUnknownMailType error = errors.New("unknown mail type")
)

type _Client struct {
	Session    uint64
	LoginCount int
}

type GateService struct {
	fractal.DefaultServiceProvider
	logger *mylog.Logger
	f      *fractal.Fractal
	option *config.GateOption
	writer func(session uint64, d []byte) error
	closer func(session uint64)

	clients map[uint64]*_Client
}

func NewGateService(logger *mylog.Logger,
	f *fractal.Fractal,
	option *config.GateOption,
	writer func(session uint64, d []byte) error,
	closer func(session uint64)) *GateService {
	return &GateService{
		logger:  logger,
		f:       f,
		option:  option,
		writer:  writer,
		closer:  closer,
		clients: make(map[uint64]*_Client, 100),
	}
}

func (c *GateService) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	c.logger.Debug(caller, session, data)
	switch caller {
	case "entry":
		return c.doEntry(_type, session, data)
	default:
		return c.doOther(caller, _type, session, data)
	}
}
