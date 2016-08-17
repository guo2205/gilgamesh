// global main.go
package main

import (
	"errors"
	"fractal/fractal"
	"gilgamesh/service/global/locker"
	"gilgamesh/service/global/online"
	"gilgamesh/utility/config"
	"log"
	"os"
	"time"

	"github.com/liuhanlcj/mylog"
)

var (
	logger       mylog.Logger = mylog.NewLogger(`Global Node`, 4, log.LstdFlags)
	lockerLogger mylog.Logger = mylog.NewLogger(`Locker Service`, 4, log.LstdFlags)
	onlineLogger mylog.Logger = mylog.NewLogger(`Online Service`, 4, log.LstdFlags)

	ErrLoadConfigFailed error = errors.New("load config failed")
)

func main() {
	nodeOption, err := loadConfig()
	if err != nil {
		return
	}

	f := fractal.NewFractal()

	f.SetLogger(log.New(os.Stdout, "[fractal]", log.Ltime))

	err = f.StartTransport(true, nodeOption.LocalAddr, nodeOption.RemoteAddr, "public.global", nodeOption.Cookie, nodeOption.Timeout)
	if err != nil {
		logger.Error("start Fractal Transport failed :", err)
		return
	}
	defer f.StopTransport()

	err = f.NewService("locker", locker.NewService(lockerLogger, f))
	if err != nil {
		logger.Error("locker service new failed :", err)
		return
	}
	defer f.StopService("locker")

	err = f.NewService("online", online.NewService(onlineLogger, f))
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
