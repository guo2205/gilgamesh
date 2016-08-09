// gate
package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type GateOption struct {
	Auth               bool
	LocalAddr          string
	Cookie             uint64
	Timeout            time.Duration
	PerSecondMaxPacket uint32
}

func LoadGateOption(filename string) (*GateOption, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	op := GateOption{}
	err = json.Unmarshal(d, &op)
	if err != nil {
		return nil, err
	}

	return &op, nil
}

func GenerateDefaultGateOption(filename string) error {
	d, err := json.Marshal(generateDefaultGateOption())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, d, 777)
}

func generateDefaultGateOption() *GateOption {
	return &GateOption{
		Auth:               true,
		LocalAddr:          "127.0.0.1:11001",
		Cookie:             12345678987654321,
		Timeout:            time.Second * 12,
		PerSecondMaxPacket: 20,
	}
}
