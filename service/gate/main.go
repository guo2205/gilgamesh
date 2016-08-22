// gate main.go
package main

import (
	"errors"
	"gilgamesh/protos"
	"gilgamesh/service/gate/entry"
	"gilgamesh/service/gate/service"
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
	logger     mylog.Logger = mylog.NewLogger(`Gate Node`, _DEBUG_LEVEL, log.LstdFlags)
	gateLogger mylog.Logger = mylog.NewLogger(`Gate Service`, _DEBUG_LEVEL, log.LstdFlags)

	ErrLoadConfigFailed error = errors.New("load config failed")
)

func main() {
	nodeOption, gateOption, err := loadConfig()
	if err != nil {
		return
	}

	f := fsdk.NewFractal(logger)

	err = f.StartHarbour(nodeOption.LocalAddr, nodeOption.RemoteAddr, "/public/gate/?", nodeOption.Cookie, nodeOption.Timeout)
	if err != nil {
		logger.Error("start Fractal Harbour failed :", err)
		return
	}
	defer f.StopHarbour()

	gateEntry := entry.NewEntry(gateOption.LocalAddr,
		f.GenerateSession,
		func(session uint64, data []byte) error {
			_, err := protos.New_GameGateService_ServiceClient(f, "gate").Call_EntryData("entry", session, &protos.Internal_GameGate_EntryDataRequest{
				Data: data,
			})
			return err
		},
		&entry.Option{
			PerSecondMaxPacket: gateOption.PerSecondMaxPacket,
		})

	gateService := service.NewService(gateLogger, f, gateOption,
		gateEntry.WriteConn, gateEntry.CloseConn)

	err = f.NewService("gate", protos.New_GameGateService_ServiceServer(f, gateService))
	if err != nil {
		logger.Error("gate service new failed :", err)
		return
	}
	defer f.StopService("gate")

	err = gateEntry.Start()
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
