// auth main.go
package main

import (
	"errors"
	"gilgamesh/protos"
	"gilgamesh/service/chat/chat"
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
	logger     mylog.Logger = mylog.NewLogger(`Chat Node`, _DEBUG_LEVEL, log.LstdFlags)
	chatLogger mylog.Logger = mylog.NewLogger(`Chat Service`, _DEBUG_LEVEL, log.LstdFlags)

	ErrLoadConfigFailed error = errors.New("load config failed")
)

func main() {
	nodeOption, err := loadConfig()
	if err != nil {
		return
	}

	f := fsdk.NewFractal(logger)

	err = f.StartHarbour(nodeOption.LocalAddr, nodeOption.RemoteAddr, "/ygo/chat", nodeOption.Cookie, nodeOption.Timeout)
	if err != nil {
		logger.Error("start Fractal Harbour failed :", err)
		return
	}
	defer f.StopHarbour()

	err = f.NewService("chat", protos.New_ChatService_ServiceServer(f, chat.NewService(chatLogger, f)))
	if err != nil {
		logger.Error("chat service new failed :", err)
		return
	}
	defer f.StopService("chat")

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
