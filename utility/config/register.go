// register
package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type RegisterOption struct {
	LocalAddr string
	Cookie    uint64
	Timeout   time.Duration
}

func LoadRegisterOption(filename string) (*RegisterOption, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	op := RegisterOption{}
	err = json.Unmarshal(d, &op)
	if err != nil {
		return nil, err
	}

	return &op, nil
}

func GenerateDefaultRegisterOption(filename string) error {
	d, err := json.Marshal(generateDefaultRegisterOption())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, d, 777)
}

func generateDefaultRegisterOption() *RegisterOption {
	return &RegisterOption{
		LocalAddr: "127.0.0.1:11000",
		Cookie:    12345678987654321,
		Timeout:   time.Second * 12,
	}
}
