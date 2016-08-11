// global main.go
package main

import (
	"fractal/fractal"
	"gilgamesh/service/global/auth"
	"gilgamesh/service/global/locker"
	"gilgamesh/service/global/online"
	"gilgamesh/utility/config"
	"gilgamesh/utility/mylog"
	"log"
	"os"
	"time"
)

var (
	logger       *mylog.Logger = mylog.NewLogger(`Global Node`, 4)
	authLogger   *mylog.Logger = mylog.NewLogger(`Auth Service`, 4)
	lockerLogger *mylog.Logger = mylog.NewLogger(`Locker Service`, 4)
	onlineLogger *mylog.Logger = mylog.NewLogger(`Online Service`, 4)
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

	err = f.NewService("auth", auth.NewService(authLogger, f))
	if err != nil {
		logger.Error("auth service new failed :", err)
		return
	}
	defer f.StopService("auth")

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
	nodeOption, err := config.LoadNodeOption("node.json")
	if err != nil {
		config.GenerateDefaultNodeOption("node.json")
		logger.Error(`load node option failed, generate default option to "node.json"`)
		return nil, err
	}

	return nodeOption, nil
}
