// main
package main

import (
	"errors"
	"gilgamesh/utility/config"
	"gilgamesh/utility/models"
	"log"

	"github.com/liuhanlcj/mylog"
)

var (
	logger mylog.Logger = mylog.NewLogger(`Tools`, 4, log.LstdFlags)

	ErrLoadConfigFailed error = errors.New("load config failed")
)

func main() {
	databaseOption, err := loadConfig()
	if err != nil {
		return
	}

	models.Init(databaseOption)
	models.Install()
}

func loadConfig() (*config.DatabaseOption, error) {
	var failed bool

	databaseOption, err := config.LoadDatabaseOption("database.json")
	if err != nil {
		config.GenerateDefaultDatabaseOption("database.json")
		logger.Error(`load database option failed, generate default option to "database.json"`)
		failed = true
	}

	if failed {
		return nil, ErrLoadConfigFailed
	}

	return databaseOption, nil
}
