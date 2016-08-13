// auth main.go
package main

import (
	"errors"
	"fractal/fractal"
	"gilgamesh/service/auth/auth"
	"gilgamesh/utility/config"
	"gilgamesh/utility/models"
	"gilgamesh/utility/mylog"
	"log"
	"os"
	"time"
)

var (
	logger     *mylog.Logger = mylog.NewLogger(`Auth Node`, 4)
	authLogger *mylog.Logger = mylog.NewLogger(`Auth Service`, 4)

	ErrLoadConfigFailed error = errors.New("load config failed")
)

func main() {
	nodeOption, databaseOption, err := loadConfig()
	if err != nil {
		return
	}

	models.Init(databaseOption)

	f := fractal.NewFractal()

	f.SetLogger(log.New(os.Stdout, "[fractal]", log.Ltime))

	err = f.StartTransport(false, nodeOption.LocalAddr, nodeOption.RemoteAddr, "public.auth", nodeOption.Cookie, nodeOption.Timeout)
	if err != nil {
		logger.Error("start Fractal Transport failed :", err)
		return
	}
	defer f.StopTransport()

	err = f.NewService("auth", auth.NewService(authLogger, f))
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
