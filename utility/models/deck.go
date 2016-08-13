// deck
package models

type Deck struct {
	Id           int64
	Account      string `xorm:"varchar(32) notnull"`
	Name         string `xorm:"varchar(32) notnull"`
	Desc         string `xorm:"varchar(128) notnull"`
	AvatarHashId string `xorm:"varchar(32) notnull"`
	DataHashId   string `xorm:"varchar(32) notnull"`
	Created      int64  `xorm:"created notnull"`
}

func CreateDeck(account, name, desc, avatarHashId, dataHashId string) error {
	dk := Deck{
		Account:      account,
		Name:         name,
		Desc:         desc,
		AvatarHashId: avatarHashId,
		DataHashId:   dataHashId,
	}
	_, err := engine.InsertOne(&dk)
	if err != nil {
		return err
	}

	return nil
}

func RemoveDeck(id uint64, account string) error {
	dk := Deck{
		Id:      int64(id),
		Account: account,
	}
	_, err := engine.Delete(&dk)
	if err != nil {
		return err
	}

	return nil
}

func QueryDeck(account string) ([]*Deck, error) {
	dkList := []*Deck{}
	dk := Deck{
		Account: account,
	}

	err := engine.Find(&dkList, &dk)
	if err != nil {
		return nil, err
	}

	return dkList, nil
}

func GetDeck(id uint64, account string) (*Deck, bool, error) {
	dk := Deck{
		Id:      int64(id),
		Account: account,
	}
	being, err := engine.Get(&dk)
	if err != nil {
		return nil, false, err
	}

	if !being {
		return nil, false, nil
	}

	return &dk, true, nil
}
