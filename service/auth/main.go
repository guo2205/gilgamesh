// auth main.go
package main

import (
	"errors"
	"gilgamesh/protos"
	"gilgamesh/service/auth/auth"
	"gilgamesh/utility/config"
	"gilgamesh/utility/models"
	"log"
	"time"

	"github.com/liuhanlcj/fractal/fractal/sdk"
	"github.com/liuhanlcj/mylog"
)

const (
	_DEBUG_LEVEL = 4
)

var (
	logger     mylog.Logger = mylog.NewLogger(`Auth Node`, _DEBUG_LEVEL, log.LstdFlags)
	authLogger mylog.Logger = mylog.NewLogger(`Auth Service`, _DEBUG_LEVEL, log.LstdFlags)

	ErrLoadConfigFailed error = errors.New("load config failed")
)

func main() {
	nodeOption, databaseOption, err := loadConfig()
	if err != nil {
		return
	}

	models.Init(databaseOption)

	f := fsdk.NewFractal(logger)

	err = f.StartHarbour(nodeOption.LocalAddr, nodeOption.RemoteAddr, "/public/auth/?", nodeOption.Cookie, nodeOption.Timeout)
	if err != nil {
		logger.Error("start Fractal Harbour failed :", err)
		return
	}
	defer f.StopHarbour()

	err = f.NewService("auth", protos.New_AuthService_ServiceServer(f, auth.NewService(authLogger)))
	if err != nil {
		logger.Error("auth service new failed :", err)
		return
	}
	defer f.StopService("auth")

	for {
		time.Sleep(time.Hour)
	}
}

func loadConfig() (*config.NodeOption, *config.DatabaseOption, error) {
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

	if failed {
		return nil, nil, ErrLoadConfigFailed
	}

	return nodeOption, databaseOption, nil
}
