// helper
package deck

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gilgamesh/protos"
	"gilgamesh/utility/models"
)

var (
	ErrDeckDataIsNil      error = errors.New("deck data is nil")
	ErrDeckDataOutOfRange error = errors.New("deck data out of range")
)

func CreateDeck(account, name, desc, avatarHash, dataHash string) error {
	return models.CreateDeck(account, name, desc, avatarHash, dataHash)
}

func RemoveDeck(id uint64, account string) error {
	return models.RemoveDeck(id, account)
}

func QueryDeck(account string) ([]*protos.Public_Deck, error) {
	dkList, err := models.QueryDeck(account)
	if err != nil {
		return nil, err
	}

	pdkList := []*protos.Public_Deck{}
	for _, v := range dkList {
		pdkList = append(pdkList, &protos.Public_Deck{
			Id:         uint64(v.Id),
			DataHashId: v.DataHashId,
			Details: &protos.Public_DeckDetails{
				Name:         v.Name,
				Desc:         v.Desc,
				AvatarHashId: v.AvatarHashId,
			},
		})
	}

	return pdkList, nil
}

func GetDeck(id uint64, account string) (*protos.Public_Deck, bool, error) {
	dk, being, err := models.GetDeck(id, account)
	if err != nil {
		return nil, false, err
	}
	if !being {
		return nil, false, nil
	}

	return &protos.Public_Deck{
		Id:         uint64(dk.Id),
		DataHashId: dk.DataHashId,
		Details: &protos.Public_DeckDetails{
			Name:         dk.Name,
			Desc:         dk.Desc,
			AvatarHashId: dk.AvatarHashId,
		},
	}, true, nil
}

func ConvertDeckToData(deck *protos.Public_DeckData) ([]byte, error) {
	if deck.Main == nil || deck.Extra == nil || deck.Side == nil {
		return nil, ErrDeckDataIsNil
	}
	if len(deck.Main) > 70 || len(deck.Extra) > 15 || len(deck.Side) > 15 {
		return nil, ErrDeckDataIsNil
	}

	wb := bytes.NewBuffer(make([]byte, 0, 412))

	binary.Write(wb, binary.LittleEndian, uint16(len(deck.Main)))
	for _, v := range deck.Main {
		binary.Write(wb, binary.LittleEndian, v)
	}

	binary.Write(wb, binary.LittleEndian, uint16(len(deck.Extra)))
	for _, v := range deck.Extra {
		binary.Write(wb, binary.LittleEndian, v)
	}

	binary.Write(wb, binary.LittleEndian, uint16(len(deck.Side)))
	for _, v := range deck.Side {
		binary.Write(wb, binary.LittleEndian, v)
	}

	return wb.Bytes(), nil
}

func ConvertDataToDeck(d []byte) (*protos.Public_DeckData, error) {
	dk := protos.Public_DeckData{}

	r := bytes.NewReader(d)

	var (
		MainLen uint32
		ExLen   uint32
		SideLen uint32
	)

	err := binary.Read(r, binary.LittleEndian, &MainLen)
	if err != nil {
		return nil, err
	}
	dk.Main = make([]uint32, int(MainLen))
	for i := 0; i < int(MainLen); i++ {
		err := binary.Read(r, binary.LittleEndian, &dk.Main[i])
		if err != nil {
			return nil, err
		}
	}

	err = binary.Read(r, binary.LittleEndian, &ExLen)
	if err != nil {
		return nil, err
	}
	dk.Extra = make([]uint32, int(ExLen))
	for i := 0; i < int(ExLen); i++ {
		err := binary.Read(r, binary.LittleEndian, &dk.Extra[i])
		if err != nil {
			return nil, err
		}
	}

	err = binary.Read(r, binary.LittleEndian, &SideLen)
	if err != nil {
		return nil, err
	}
	dk.Side = make([]uint32, int(SideLen))
	for i := 0; i < int(SideLen); i++ {
		err := binary.Read(r, binary.LittleEndian, &dk.Side[i])
		if err != nil {
			return nil, err
		}
	}

	return &dk, nil
}
