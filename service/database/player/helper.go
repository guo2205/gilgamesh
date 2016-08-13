// helper
package player

import (
	"gilgamesh/protos"
	"gilgamesh/utility/models"
)

func ModifyAccountPlayer(account string, data *protos.Public_PlayerData) (bool, error) {
	return models.ModifyAccountPlayer(account, data)
}

func CreateAccountPlayer(account string, data *protos.Public_PlayerData) (bool, error) {
	return models.CreateAccountPlayer(account, data)
}

func GetAccountPlayer(account string) (*protos.Public_Player, bool, error) {
	player, being, err := models.GetAccountPlayer(account)
	if err != nil {
		return nil, false, err
	}
	if !being {
		return nil, false, nil
	}
	return &protos.Public_Player{
		Account: player.Account,
		Data: &protos.Public_PlayerData{
			NickName:     player.NickName,
			AvatarHashId: player.Avatar,
		},
		GameData: &protos.Public_PlayerGameData{
			Level: uint32(player.Level),
			VIP:   player.Vip,
		},
	}, true, nil
}

func GetAccountsPlayer(accounts []string) ([]*protos.Public_Player, error) {
	players, err := models.GetAccountsPlayer(accounts)
	if err != nil {
		return nil, err
	}
	nplayers := []*protos.Public_Player{}
	for _, player := range players {
		nplayers = append(nplayers, &protos.Public_Player{
			Account: player.Account,
			Data: &protos.Public_PlayerData{
				NickName:     player.NickName,
				AvatarHashId: player.Avatar,
			},
			GameData: &protos.Public_PlayerGameData{
				Level: uint32(player.Level),
				VIP:   player.Vip,
			},
		})
	}
	return nplayers, nil
}
