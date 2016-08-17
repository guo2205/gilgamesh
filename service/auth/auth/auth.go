// auth
package auth

import (
	"encoding/hex"
	"fractal/fractal"
	"gilgamesh/protos"
	"gilgamesh/utility/models"
	"gilgamesh/utility/utils"

	"github.com/golang/protobuf/proto"
	"github.com/liuhanlcj/mylog"
)

type Service struct {
	fractal.DefaultServiceProvider
	logger mylog.Logger
	f      *fractal.Fractal
}

func NewService(
	logger mylog.Logger,
	f *fractal.Fractal) *Service {
	return &Service{
		logger: logger,
		f:      f,
	}
}

func (c *Service) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	obj, ptype, err := utils.Unmarshal(data)
	if err != nil {
		return []byte{}, err
	}

	switch ptype {
	case proto.MessageName((*protos.Public_Cts_Login)(nil)):
		return c.do_Public_Cts_Login(session, obj.(*protos.Public_Cts_Login))
	}
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Login(session uint64, obj *protos.Public_Cts_Login) ([]byte, error) {
	reason, ok, err := models.AccountVerifyPassword(obj.Account, hex.EncodeToString(obj.Password))
	if err != nil || !ok {
		return utils.Marshal(&protos.Public_Stc_LoginResponse{
			State:  false,
			Reason: reason,
		}), nil
	}

	return utils.Marshal(&protos.Public_Stc_LoginResponse{
		State: true,
	}), nil
}
