// gate
package service

import (
	"errors"
	"fractal/fractal"
	"gilgamesh/utility/config"

	"github.com/liuhanlcj/mylog"
)

var (
	ErrUnknownMailType error = errors.New("unknown mail type")
)

type _Client struct {
	Session uint64
	Account string
}

type Service struct {
	fractal.DefaultServiceProvider
	logger mylog.Logger
	f      *fractal.Fractal
	option *config.GateOption
	writer func(session uint64, d []byte) error
	closer func(session uint64)

	clients map[uint64]*_Client
}

func NewGateService(logger mylog.Logger,
	f *fractal.Fractal,
	option *config.GateOption,
	writer func(session uint64, d []byte) error,
	closer func(session uint64)) *Service {
	return &Service{
		logger:  logger,
		f:       f,
		option:  option,
		writer:  writer,
		closer:  closer,
		clients: make(map[uint64]*_Client, 100),
	}
}

func (c *Service) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	switch caller {
	case "entry":
		return c.doEntry(_type, session, data)
	default:
		return c.doOther(caller, _type, session, data)
	}
}
