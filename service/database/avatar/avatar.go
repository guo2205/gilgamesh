// avatar
package avatar

import (
	"fractal/fractal"
	"gilgamesh/protos"
	"gilgamesh/utility/utils"
	wdk "gilgamesh/utility/weed-sdk"

	"github.com/golang/protobuf/proto"
	"github.com/liuhanlcj/mylog"
)

type Service struct {
	fractal.DefaultServiceProvider
	logger mylog.Logger
	f      *fractal.Fractal

	sdk *wdk.WeedSdk
}

func NewService(
	logger mylog.Logger,
	f *fractal.Fractal,
	sdk *wdk.WeedSdk) *Service {
	return &Service{
		logger: logger,
		f:      f,
		sdk:    sdk,
	}
}

func (c *Service) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	obj, ptype, err := utils.Unmarshal(data)
	if err != nil {
		return []byte{}, err
	}

	switch ptype {
	case proto.MessageName((*protos.Public_Cts_Resource_GetAvatar)(nil)):
		return c.do_Public_Cts_Resource_GetAvatar(session, obj.(*protos.Public_Cts_Resource_GetAvatar))
	case proto.MessageName((*protos.Public_Cts_Resource_UploadAvatar)(nil)):
		return c.do_Public_Cts_Resource_UploadAvatar(session, obj.(*protos.Public_Cts_Resource_UploadAvatar))
	}
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Resource_GetAvatar(session uint64, obj *protos.Public_Cts_Resource_GetAvatar) ([]byte, error) {
	d, err := c.sdk.GetFile(obj.HashId)
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Resource_GetAvatarResponse{
			AvatarHashId: obj.HashId,
		}), nil
	}
	return utils.Marshal(&protos.Public_Stc_Resource_GetAvatarResponse{
		AvatarHashId: obj.HashId,
		Data:         d,
	}), nil
}

func (c *Service) do_Public_Cts_Resource_UploadAvatar(session uint64, obj *protos.Public_Cts_Resource_UploadAvatar) ([]byte, error) {
	fid, _, err := c.sdk.SaveFile(obj.Data)
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Resource_UploadAvatarResponse{}), nil
	}
	return utils.Marshal(&protos.Public_Stc_Resource_UploadAvatarResponse{
		AvatarHashId: fid,
	}), nil
}
