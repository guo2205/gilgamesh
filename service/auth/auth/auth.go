// auth
package auth

import (
	"encoding/hex"
	"gilgamesh/protos"
	"gilgamesh/utility/models"

	"github.com/liuhanlcj/mylog"
)

type Service struct {
	logger mylog.Logger
}

func NewService(logger mylog.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

func (c *Service) On_Auth(caller string, session uint64, in *protos.Auth_AuthRequest, responser func(out *protos.Auth_AuthResponse, e error)) {
	reason, ok, err := models.AccountVerifyPassword(in.Account, hex.EncodeToString(in.Password))
	if err != nil {
		c.logger.Warningf("[%s %d] verify database failed :%v\n", in.Account, session, err)
		responser(nil, err)
		return
	}

	if !ok {
		responser(&protos.Auth_AuthResponse{
			Success: false,
			Reason:  reason,
		}, nil)
		return
	}

	responser(&protos.Auth_AuthResponse{
		Success: true,
	}, nil)
}
