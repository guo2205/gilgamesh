// global main.go
package main

import (
	"errors"
	"fractal/fractal"
	"gilgamesh/service/database/avatar"
	"gilgamesh/service/database/deck"
	"gilgamesh/service/database/lflist"
	"gilgamesh/service/database/player"
	"gilgamesh/service/database/videotape"
	"gilgamesh/utility/config"
	clflist "gilgamesh/utility/lflist"
	"gilgamesh/utility/models"
	"gilgamesh/utility/mylog"
	wdk "gilgamesh/utility/weed-sdk"
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

	ErrLoadConfigFailed error = errors.New("load config failed")
)

func main() {
	nodeOption, databaseOption, resourceOption, err := loadConfig()
	if err != nil {
		return
	}

	models.Init(databaseOption)

	f := fractal.NewFractal()

	f.SetLogger(log.New(os.Stdout, "[fractal]", log.Ltime))

	err = f.StartTransport(false, nodeOption.LocalAddr, nodeOption.RemoteAddr, "ygo.database", nodeOption.Cookie, nodeOption.Timeout)
	if err != nil {
		logger.Error("start Fractal Transport failed :", err)
		return
	}
	defer f.StopTransport()

	err = f.NewService("avatar", avatar.NewService(avatarLogger, f, wdk.NewWeedSdk(resourceOption)))
	if err != nil {
		logger.Error("avatar service new failed :", err)
		return
	}
	defer f.StopService("avatar")

	err = f.NewService("deck", deck.NewService(deckLogger, f, wdk.NewWeedSdk(resourceOption)))
	if err != nil {
		logger.Error("deck service new failed :", err)
		return
	}
	defer f.StopService("deck")

	lf, err := clflist.NewLFListContainer(resourceOption)
	if err != nil {
		logger.Error("lflist load failed :", err)
		return
	}

	err = f.NewService("lflist", lflist.NewService(lflistLogger, f, lf))
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

func loadConfig() (*config.NodeOption, *config.DatabaseOption, *config.ResourceOption, error) {
	var failed bool

	nodeOption, err := config.LoadNodeOption("node.json")
	if err != nil {
		config.GenerateDefaultNodeOption("node.json")
		logger.Error(`load node option failed, generate default option to "node.json"`)
		failed = true
	}

	databaseOption, err := config.LoadDatabaseOption("database.json")
	if err != nil {
		config.GenerateDefaultDatabaseOption("database.json")
		logger.Error(`load database option failed, generate default option to "database.json"`)
		failed = true
	}

	resourceOption, err := config.LoadResourceOption("resource.json")
	if err != nil {
		config.GenerateDefaultResourceOption("resource.json")
		logger.Error(`load resource option failed, generate default option to "resource.json"`)
		failed = true
	}

	if failed {
		return nil, nil, nil, ErrLoadConfigFailed
	}

	return nodeOption, databaseOption, resourceOption, nil
}
