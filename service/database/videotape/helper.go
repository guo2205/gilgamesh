// helper
package videotape

import (
	"gilgamesh/protos"
	"gilgamesh/utility/models"
)

func CreateVideotape(account, name, hashName string) error {
	return models.CreateVideotape(account, name, hashName)
}

func QueryVideotapeList(account string) ([]*protos.Public_VideoTape, error) {
	vtList, err := models.QueryVideotapeList(account)
	if err != nil {
		return nil, err
	}

	pvtList := []*protos.Public_VideoTape{}
	for _, v := range vtList {
		pvtList = append(pvtList, &protos.Public_VideoTape{
			HashId:    v.Hash,
			Name:      v.Name,
			Timestamp: uint64(v.Created),
		})
	}

	return pvtList, nil
}
