// hall main.go
package main

import (
	"errors"
	"log"
	"time"

	"github.com/liuhanlcj/fractal/fractal/sdk"
	"github.com/liuhanlcj/gilgamesh/protos"
	"github.com/liuhanlcj/gilgamesh/service/hall/hall"
	"github.com/liuhanlcj/gilgamesh/utility/config"
	"github.com/liuhanlcj/mylog"
)

const (
	_DEBUG_LEVEL = 4
)

var (
	logger     mylog.Logger = mylog.NewLogger(`Hall Node`, _DEBUG_LEVEL, log.LstdFlags)
	hallLogger mylog.Logger = mylog.NewLogger(`Hall Service`, _DEBUG_LEVEL, log.LstdFlags)

	ErrLoadConfigFailed error = errors.New("load config failed")
)

func main() {
	nodeOption, hallOption, err := loadConfig()
	if err != nil {
		return
	}

	f := fsdk.NewFractal(logger)

	err = f.StartHarbour(nodeOption.LocalAddr, nodeOption.RemoteAddr, "/ygo/hall", nodeOption.Cookie, nodeOption.Timeout)
	if err != nil {
		logger.Error("start Fractal Harbour failed :", err)
		return
	}
	defer f.StopHarbour()

	err = f.NewService("hall", protos.New_HallService_ServiceServer(f, hall.NewService(hallLogger, f, hallOption)))
	if err != nil {
		logger.Error("hall service new failed :", err)
		return
	}
	defer f.StopService("hall")

	for {
		time.Sleep(time.Hour)
	}
}

func loadConfig() (*config.NodeOption, *config.HallOption, error) {
	var failed bool

	nodeOption, err := config.LoadNodeOption("node.json")
	if err != nil {
		config.GenerateDefaultNodeOption("node.json")
		logger.Error(`load node option failed, generate default option to "node.json"`)
		failed = true
	}

	hallOption, err := config.LoadHallOption("hall.json")
	if err != nil {
		config.GenerateDefaultHallOption("hall.json")
		logger.Error(`load hall option failed, generate default option to "hall.json"`)
		failed = true
	}

	if failed {
		return nil, nil, ErrLoadConfigFailed
	}

	return nodeOption, hallOption, nil
}
