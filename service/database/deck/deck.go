// deck
package deck

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
	case proto.MessageName((*protos.Internal_Database_Deck_Download)(nil)):
		return c.do_Internal_Database_Deck_Download(session, obj.(*protos.Internal_Database_Deck_Download))
	case proto.MessageName((*protos.Internal_Database_Deck_Query)(nil)):
		return c.do_Internal_Database_Deck_Query(session, obj.(*protos.Internal_Database_Deck_Query))
	case proto.MessageName((*protos.Internal_Database_Deck_Remove)(nil)):
		return c.do_Internal_Database_Deck_Remove(session, obj.(*protos.Internal_Database_Deck_Remove))
	case proto.MessageName((*protos.Internal_Database_Deck_Upload)(nil)):
		return c.do_Internal_Database_Deck_Upload(session, obj.(*protos.Internal_Database_Deck_Upload))
	}
	return []byte{}, nil
}

func (c *Service) do_Internal_Database_Deck_Download(session uint64, obj *protos.Internal_Database_Deck_Download) ([]byte, error) {
	dk, being, err := GetDeck(obj.Id, obj.Account)
	if err != nil || !being {
		return utils.Marshal(&protos.Public_Stc_Deck_DownloadResponse{
			State: false,
		}), nil
	}

	d, err := c.sdk.GetFile(dk.DataHashId)
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Deck_DownloadResponse{
			State: false,
		}), nil
	}

	dkd, err := ConvertDataToDeck(d)
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Deck_DownloadResponse{
			State: false,
		}), nil
	}

	dk.Data = dkd

	return utils.Marshal(&protos.Public_Stc_Deck_DownloadResponse{
		State: true,
		Deck:  dk,
	}), nil
}

func (c *Service) do_Internal_Database_Deck_Query(session uint64, obj *protos.Internal_Database_Deck_Query) ([]byte, error) {
	dkList, err := QueryDeck(obj.Account)
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Deck_QueryListResponse{}), nil
	}

	return utils.Marshal(&protos.Public_Stc_Deck_QueryListResponse{
		DeckList: dkList,
	}), nil
}

func (c *Service) do_Internal_Database_Deck_Remove(session uint64, obj *protos.Internal_Database_Deck_Remove) ([]byte, error) {
	err := RemoveDeck(obj.Id, obj.Account)
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Deck_RemoveResponse{
			State: false,
		}), nil
	}

	return utils.Marshal(&protos.Public_Stc_Deck_RemoveResponse{
		State: true,
	}), nil
}

func (c *Service) do_Internal_Database_Deck_Upload(session uint64, obj *protos.Internal_Database_Deck_Upload) ([]byte, error) {
	if obj.Deck.Data == nil {
		return utils.Marshal(&protos.Public_Stc_Deck_UploadResponse{
			State: false,
		}), nil
	}

	d, err := ConvertDeckToData(obj.Deck.Data)
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Deck_UploadResponse{
			State: false,
		}), nil
	}

	fid, _, err := c.sdk.SaveFile(d)
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Deck_UploadResponse{
			State: false,
		}), nil
	}

	err = CreateDeck(obj.Account, obj.Deck.Details.Name, obj.Deck.Details.Desc, obj.Deck.Details.AvatarHashId, fid)
	if err != nil {
		return utils.Marshal(&protos.Public_Stc_Deck_UploadResponse{
			State: false,
		}), nil
	}

	return utils.Marshal(&protos.Public_Stc_Deck_UploadResponse{
		State: true,
	}), nil
}
