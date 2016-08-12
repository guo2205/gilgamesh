// global main.go
package main

import (
	"fractal/fractal"
	"gilgamesh/service/database/avatar"
	"gilgamesh/service/database/deck"
	"gilgamesh/service/database/lflist"
	"gilgamesh/service/database/player"
	"gilgamesh/service/database/videotape"
	"gilgamesh/utility/config"
	"gilgamesh/utility/mylog"
	"log"
	"os"
	"time"
)

var (
	logger          *mylog.Logger = mylog.NewLogger(`Database Node`, 4)
	avatarLogger    *mylog.Logger = mylog.NewLogger(`Avatar Service`, 4)
	deckLogger      *mylog.Logger = mylog.NewLogger(`Deck Service`, 4)
	lflistLogger    *mylog.Logger = mylog.NewLogger(`Lflist Service`, 4)
	playerLogger    *mylog.Logger = mylog.NewLogger(`Player Service`, 4)
	videotapeLogger *mylog.Logger = mylog.NewLogger(`Videotape Service`, 4)
)

func main() {
	nodeOption, err := loadConfig()
	if err != nil {
		return
	}

	f := fractal.NewFractal()

	f.SetLogger(log.New(os.Stdout, "[fractal]", log.Ltime))

	err = f.StartTransport(false, nodeOption.LocalAddr, nodeOption.RemoteAddr, "ygo.database", nodeOption.Cookie, nodeOption.Timeout)
	if err != nil {
		logger.Error("start Fractal Transport failed :", err)
		return
	}
	defer f.StopTransport()

	err = f.NewService("avatar", avatar.NewService(avatarLogger, f))
	if err != nil {
		logger.Error("avatar service new failed :", err)
		return
	}
	defer f.StopService("avatar")

	err = f.NewService("deck", deck.NewService(deckLogger, f))
	if err != nil {
		logger.Error("deck service new failed :", err)
		return
	}
	defer f.StopService("deck")

	err = f.NewService("lflist", lflist.NewService(lflistLogger, f))
	if err != nil {
		logger.Error("lflist service new failed :", err)
		return
	}
	defer f.StopService("lflist")

	err = f.NewService("player", player.NewService(playerLogger, f))
	if err != nil {
		logger.Error("player service new failed :", err)
		return
	}
	defer f.StopService("player")

	err = f.NewService("videotape", videotape.NewService(videotapeLogger, f))
	if err != nil {
		logger.Error("videotape service new failed :", err)
		return
	}
	defer f.StopService("videotape")

	for {
		time.Sleep(time.Hour)
	}
}

func loadConfig() (*config.NodeOption, error) {
	nodeOption, err := config.LoadNodeOption("node.json")
	if err != nil {
		config.GenerateDefaultNodeOption("node.json")
		logger.Error(`load node option failed, generate default option to "node.json"`)
		return nil, err
	}

	return nodeOption, nil
}
