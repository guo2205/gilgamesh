// videotape
package videotape

import (
	"fractal/fractal"
	"gilgamesh/protos"
	"gilgamesh/utility/mylog"
	"gilgamesh/utility/utils"
	wdk "gilgamesh/utility/weed-sdk"

	"github.com/golang/protobuf/proto"
)

type Service struct {
	fractal.DefaultServiceProvider
	logger *mylog.Logger
	f      *fractal.Fractal

	sdk *wdk.WeedSdk
}

func NewService(
	logger *mylog.Logger,
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
	case proto.MessageName((*protos.Internal_Database_Videotape_Get)(nil)):
		return c.do_Internal_Database_Videotape_Get(session, obj.(*protos.Internal_Database_Videotape_Get))
	case proto.MessageName((*protos.Internal_Database_Videotape_QueryList)(nil)):
		return c.do_Internal_Database_Videotape_QueryList(session, obj.(*protos.Internal_Database_Videotape_QueryList))
	}
	return []byte{}, nil
}

func (c *Service) do_Internal_Database_Videotape_Get(session uint64, obj *protos.Internal_Database_Videotape_Get) ([]byte, error) {
	d, err := c.sdk.GetFile(obj.HashId)
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Videotape_VideoTapeData{
			HashId: obj.HashId,
		}), nil
	}

	return utils.Marshal(&protos.Public_Stc_Videotape_VideoTapeData{
		HashId: obj.HashId,
		Data:   d,
	}), nil
}

func (c *Service) do_Internal_Database_Videotape_QueryList(session uint64, obj *protos.Internal_Database_Videotape_QueryList) ([]byte, error) {
	vtList, err := QueryVideotapeList(obj.Account)
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Videotape_VideoTapeList{}), nil
	}

	return utils.Marshal(&protos.Public_Stc_Videotape_VideoTapeList{
		VideoTapeList: vtList,
	}), nil
}
