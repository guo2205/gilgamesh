// lflist
package lflist

import (
	"fractal/fractal"
	"gilgamesh/protos"
	"gilgamesh/utility/lflist"
	"gilgamesh/utility/utils"

	"github.com/golang/protobuf/proto"
	"github.com/liuhanlcj/mylog"
)

type Service struct {
	fractal.DefaultServiceProvider
	logger mylog.Logger
	f      *fractal.Fractal

	lflist *lflist.LFListContainer
}

func NewService(
	logger mylog.Logger,
	f *fractal.Fractal,
	lflist *lflist.LFListContainer) *Service {
	return &Service{
		logger: logger,
		f:      f,
		lflist: lflist,
	}
}

func (c *Service) OnMail(caller string, _type uint32, session uint64, data []byte) ([]byte, error) {
	obj, ptype, err := utils.Unmarshal(data)
	if err != nil {
		return []byte{}, err
	}

	switch ptype {
	case proto.MessageName((*protos.Public_Cts_Resource_GetLFList)(nil)):
		return c.do_Public_Cts_Resource_GetLFList(session, obj.(*protos.Public_Cts_Resource_GetLFList))
	case proto.MessageName((*protos.Public_Cts_Resource_GetLFListData)(nil)):
		return c.do_Public_Cts_Resource_GetLFListData(session, obj.(*protos.Public_Cts_Resource_GetLFListData))
	}
	return []byte{}, nil
}

func (c *Service) do_Public_Cts_Resource_GetLFList(session uint64, obj *protos.Public_Cts_Resource_GetLFList) ([]byte, error) {
	return utils.Marshal(&protos.Public_Stc_Resource_GetLFListResponse{
		LFListList: c.lflist.GetList(),
	}), nil
}

func (c *Service) do_Public_Cts_Resource_GetLFListData(session uint64, obj *protos.Public_Cts_Resource_GetLFListData) ([]byte, error) {
	lf, err := c.lflist.GetData(int(obj.Id))
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Resource_GetLFListDataResponse{}), nil
	}

	return utils.Marshal(&protos.Public_Stc_Resource_GetLFListDataResponse{
		LFList: lf,
	}), nil
}
