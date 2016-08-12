// global main.go
package main

import (
	"fractal/fractal"
	"gilgamesh/service/database/auth"
	"gilgamesh/utility/config"
	"gilgamesh/utility/mylog"
	"log"
	"os"
	"time"
)

var (
	logger     *mylog.Logger = mylog.NewLogger(`Database Node`, 4)
	authLogger *mylog.Logger = mylog.NewLogger(`Auth Service`, 4)
)

func main() {
	nodeOption, err := loadConfig()
	if err != nil {
		return
	}

	f := fractal.NewFractal()

	f.SetLogger(log.New(os.Stdout, "[fractal]", log.Ltime))

	err = f.StartTransport(true, nodeOption.LocalAddr, nodeOption.RemoteAddr, "ygo.database", nodeOption.Cookie, nodeOption.Timeout)
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

func loadConfig() (*config.NodeOption, error) {
	nodeOption, err := config.LoadNodeOption("node.json")
	if err != nil {
		config.GenerateDefaultNodeOption("node.json")
		logger.Error(`load node option failed, generate default option to "node.json"`)
		return nil, err
	}

	return nodeOption, nil
}
