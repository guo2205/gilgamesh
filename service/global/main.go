// global main.go
package main

import (
	"errors"
	"gilgamesh/protos"
	"gilgamesh/service/global/locker"
	"gilgamesh/service/global/online"
	"gilgamesh/utility/config"
	"log"
	"time"

	"github.com/liuhanlcj/fractal/fractal/sdk"
	"github.com/liuhanlcj/mylog"
)

const (
	_DEBUG_LEVEL = 4
)

var (
	logger       mylog.Logger = mylog.NewLogger(`Global Node`, _DEBUG_LEVEL, log.LstdFlags)
	lockerLogger mylog.Logger = mylog.NewLogger(`Locker Service`, _DEBUG_LEVEL, log.LstdFlags)
	onlineLogger mylog.Logger = mylog.NewLogger(`Online Service`, _DEBUG_LEVEL, log.LstdFlags)

	ErrLoadConfigFailed error = errors.New("load config failed")
)

func main() {
	nodeOption, err := loadConfig()
	if err != nil {
		return
	}

	f := fsdk.NewFractal(logger)

	err = f.StartHarbour(nodeOption.LocalAddr, nodeOption.RemoteAddr, "/public/global", nodeOption.Cookie, nodeOption.Timeout)
	if err != nil {
		logger.Error("start Fractal Harbour failed :", err)
		return
	}
	defer f.StopHarbour()

	err = f.NewService("locker", protos.New_GlobalLockerService_ServiceServer(f, locker.NewService(lockerLogger)))
	if err != nil {
		logger.Error("locker service new failed :", err)
		return
	}
	defer f.StopService("locker")

	err = f.NewService("online", protos.New_GlobalOnlineStateService_ServiceServer(f, online.NewService(onlineLogger)))
	if err != nil {
		logger.Error("online service new failed :", err)
		return
	}
	defer f.StopService("online")

	for {
		time.Sleep(time.Hour)
	}
}

func loadConfig() (*config.NodeOption, error) {
	var failed bool

	nodeOption, err := config.LoadNodeOption("node.json")
	if err != nil {
		config.GenerateDefaultNodeOption("node.json")
		logger.Error(`load node option failed, generate default option to "node.json"`)
		failed = true
	}

	if failed {
		return nil, ErrLoadConfigFailed
	}

	return nodeOption, nil
}
