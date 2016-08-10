// gate main.go
package main

import (
	"fractal/fractal"
	"gilgamesh/service/gate/entry"
	"gilgamesh/service/gate/service"
	"gilgamesh/utility/config"
	"gilgamesh/utility/mylog"
	"log"
	"os"
	"time"
)

var (
	logger *mylog.Logger = mylog.NewLogger(`Gate Service`, 4)
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

	gateService := service.NewGateService(logger, f, gateOption,
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
	nodeOption, err := config.LoadNodeOption("node.json")
	if err != nil {
		config.GenerateDefaultNodeOption("node.json")
		logger.Error(`load node option failed, generate default option to "node.json"`)
		return nil, nil, err
	}

	gateOption, err := config.LoadGateOption("gate.json")
	if err != nil {
		config.GenerateDefaultGateOption("gate.json")
		logger.Error(`load gate option failed, generate default option to "gate.json"`)
		return nil, nil, err
	}

	return nodeOption, gateOption, nil
}
