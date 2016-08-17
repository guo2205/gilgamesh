// gate main.go
package main

import (
	"errors"
	"fractal/fractal"
	"gilgamesh/service/gate/entry"
	"gilgamesh/service/gate/service"
	"gilgamesh/utility/config"
	"log"
	"os"
	"time"

	"github.com/liuhanlcj/mylog"
)

var (
	logger     mylog.Logger = mylog.NewLogger(`Gate Node`, 4, log.LstdFlags)
	gateLogger mylog.Logger = mylog.NewLogger(`Gate Service`, 4, log.LstdFlags)

	ErrLoadConfigFailed error = errors.New("load config failed")
)

func main() {
	nodeOption, gateOption, err := loadConfig()
	if err != nil {
		return
	}

	f := fractal.NewFractal()

	f.SetLogger(log.New(os.Stdout, "[fractal]", log.Ltime))

	err = f.StartTransport(false, nodeOption.LocalAddr, nodeOption.RemoteAddr, "public.gate", nodeOption.Cookie, nodeOption.Timeout)
	if err != nil {
		logger.Error("start Fractal Transport failed :", err)
		return
	}
	defer f.StopTransport()

	gateEntry := entry.NewGateEntry(f, gateOption)

	gateService := service.NewGateService(gateLogger, f, gateOption,
		gateEntry.WritePacket, gateEntry.CloseConn)

	err = f.NewService("gate", gateService)
	if err != nil {
		logger.Error("gate service new failed :", err)
		return
	}
	defer f.StopService("gate")

	err = gateEntry.Run()
	if err != nil {
		logger.Error("gate entry run failed :", err)
		return
	}
	logger.Info("gate entry start at", gateOption.LocalAddr)

	for {
		time.Sleep(time.Hour)
	}
}

func loadConfig() (*config.NodeOption, *config.GateOption, error) {
	var failed bool

	nodeOption, err := config.LoadNodeOption("node.json")
	if err != nil {
		config.GenerateDefaultNodeOption("node.json")
		logger.Error(`load node option failed, generate default option to "node.json"`)
		failed = true
	}

	gateOption, err := config.LoadGateOption("gate.json")
	if err != nil {
		config.GenerateDefaultGateOption("gate.json")
		logger.Error(`load gate option failed, generate default option to "gate.json"`)
		failed = true
	}

	if failed {
		return nil, nil, ErrLoadConfigFailed
	}

	return nodeOption, gateOption, nil
}
