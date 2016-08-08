// node
package config

import (
	"encoding/json"
	"io/ioutil"
)

type NodeOption struct {
	RemoteAddr string
	LocalAddr  string
	Cookie     uint64
	Timeout    uint32
}

func LoadNodeOption(filename string) (*NodeOption, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	op := NodeOption{}
	err = json.Unmarshal(d, &op)
	if err != nil {
		return nil, err
	}

	return &op, nil
}

func GenerateDefaultNodeOption(filename string) error {
	d, err := json.Marshal(generateDefaultNodeOption())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, d, 777)
}

func generateDefaultNodeOption() *NodeOption {
	return &NodeOption{
		RemoteAddr: "127.0.0.1:6679",
		LocalAddr:  "127.0.0.1:0",
		Cookie:     12345678987654321,
		Timeout:    10,
	}
}
