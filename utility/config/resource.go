// resource
package config

import (
	"encoding/json"
	"io/ioutil"
)

type ResourceOption struct {
	LFList      string
	StoreServer string
}

func LoadResourceOption(filename string) (*ResourceOption, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	op := ResourceOption{}
	err = json.Unmarshal(d, &op)
	if err != nil {
		return nil, err
	}

	return &op, nil
}

func GenerateDefaultResourceOption(filename string) error {
	d, err := json.Marshal(generateDefaultResourceOption())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, d, 777)
}

func generateDefaultResourceOption() *ResourceOption {
	return &ResourceOption{
		LFList:      "./card/lflist.conf",
		StoreServer: "localhost:10003",
	}
}
