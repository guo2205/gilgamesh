// main
package main

import (
	"errors"
	"gilgamesh/service/register/entry2"
	"gilgamesh/utility/config"
	"gilgamesh/utility/models"
	"log"

	"github.com/liuhanlcj/mylog"
)

var (
	logger mylog.Logger = mylog.NewLogger(`Register Node`, 4, log.LstdFlags)

	ErrLoadConfigFailed error = errors.New("load config failed")
)

func main() {
	registerOption, databaseOption, err := loadConfig()
	if err != nil {
		return
	}

	models.Init(databaseOption)

	gateEntry := entry2.NewGateEntry2(registerOption)

	logger.Info("register entry start at", registerOption.LocalAddr)

	gateEntry.Run()
}

func loadConfig() (*config.RegisterOption, *config.DatabaseOption, error) {
	var failed bool

	registerOption, err := config.LoadRegisterOption("register.json")
	if err != nil {
		config.GenerateDefaultRegisterOption("register.json")
		logger.Error(`load register option failed, generate default option to "register.json"`)
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

	return registerOption, databaseOption, nil
}
