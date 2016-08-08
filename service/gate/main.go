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
	nodeOption, err := config.LoadNodeOption("node.json")
	if err != nil {
		config.GenerateDefaultNodeOption("node.json")
		logger.Error(`load node option failed, generate default option to "node.json"`)
		return
	}

	f := fractal.NewFractal()

	f.SetLogger(log.New(os.Stdout, "[fractal]", log.Ltime))

	err = f.StartTransport(nodeOption.LocalAddr, nodeOption.RemoteAddr, "public.gate", nodeOption.Cookie, time.Second*time.Duration(nodeOption.Timeout))
	if err != nil {
		logger.Error("start Fractal Transport failed :", err)
		return
	}
	defer f.StopTransport()

	gateEntry := &entry.GateEntry{
		Fractal: f,
	}

	gateService := &service.GateService{
		Fractal:     f,
		WritePacket: gateEntry.WritePacket,
	}

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
	defer gateEntry.Close()

	for {
		time.Sleep(time.Hour)
	}
}
