// player
package models

import (
	"fmt"
	"gilgamesh/protos"
	"strings"
)

type Player struct {
	Id       int64
	Account  string `xorm:"notnull"`
	Level    uint8  `xorm:"notnull"`
	Vip      bool   `xorm:"notnull"`
	NickName string `xorm:"varchar(32) notnull"`
	Avatar   string `xorm:"varchar(32) notnull"`
	UseDeck  uint64 `xorm:"notnull"`
}

func ModifyAccountPlayer(account string, data protos.Public_PlayerData) (bool, error) {
	py := Player{
		Account: account,
	}
	ok, err := engine.Get(&py)
	if err != nil {
		return false, err
	}
	if ok {
		return false, nil
	}

	sql := "update player set "
	if data.NickName != "" {
		sql += fmt.Sprintf(`nick_name = '%s' `, data.NickName)
	}
	if data.AvatarHashId != "" {
		sql += fmt.Sprintf(`avatar = '%s' `, data.AvatarHashId)
	}
	if data.DeckId != 0 {
		sql += fmt.Sprintf(`use_deck = '%d' `, data.DeckId)
	}
	sql += "where id = ?"

	_, err = engine.Exec(sql, py.Id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateAccountPlayer(account string, data protos.Public_PlayerData) (bool, error) {
	py := Player{
		Account: account,
	}
	ok, err := engine.Get(&py)
	if err != nil {
		return false, err
	}
	if ok {
		return false, nil
	}

	py = Player{
		Account:  account,
		Level:    0,
		Vip:      false,
		NickName: data.NickName,
		Avatar:   data.AvatarHashId,
		UseDeck:  data.DeckId,
	}

	_, err = engine.InsertOne(&py)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetAccountPlayer(account string) (*Player, bool, error) {
	player := Player{
		Account: account,
	}
	ok, err := engine.Get(&player)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	return &player, true, nil
}

func GetAccountsPlayer(accounts []string) ([]*Player, error) {
	players := make([]*Player, 0, len(accounts))
	for i := range accounts {
		accounts[i] = "'" + accounts[i] + "'"
	}
	err := engine.Sql(fmt.Sprintf("select * from player where account in (%s)", strings.Join(accounts, ","))).Find(&players)
	if err != nil {
		return nil, err
	}
	return players, nil
}
