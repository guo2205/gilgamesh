// auth
package auth

import (
	"fractal/fractal"
	"gilgamesh/utility/mylog"
)

type AuthService struct {
	fractal.DefaultServiceProvider
	logger *mylog.Logger
	f      *fractal.Fractal
}

func NewService(
	logger *mylog.Logger,
	f *fractal.Fractal) *AuthService {
	return &AuthService{
		logger: logger,
		f:      f,
	}
}

func (c *AuthService) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	c.logger.Debug(caller, _type, session, data)
	return []byte{}, nil
}
