// hall main.go
package main

import (
	"errors"
	"fractal/fractal"
	"gilgamesh/service/hall/hall"
	"gilgamesh/utility/config"
	"log"
	"os"
	"time"

	"github.com/liuhanlcj/mylog"
)

var (
	logger     mylog.Logger = mylog.NewLogger(`Hall Node`, 4, log.LstdFlags)
	hallLogger mylog.Logger = mylog.NewLogger(`Hall Service`, 4, log.LstdFlags)

	ErrLoadConfigFailed error = errors.New("load config failed")
)

func main() {
	nodeOption, err := loadConfig()
	if err != nil {
		return
	}

	f := fractal.NewFractal()

	f.SetLogger(log.New(os.Stdout, "[fractal]", log.Ltime))

	err = f.StartTransport(true, nodeOption.LocalAddr, nodeOption.RemoteAddr, "ygo.hall", nodeOption.Cookie, nodeOption.Timeout)
	if err != nil {
		logger.Error("start Fractal Transport failed :", err)
		return
	}
	defer f.StopTransport()

	err = f.NewService("hall", hall.NewService(hallLogger, f))
	if err != nil {
		logger.Error("hall service new failed :", err)
		return
	}
	defer f.StopService("hall")

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
