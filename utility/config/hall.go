// hall
package config

import (
	"encoding/json"
	"io/ioutil"
)

type HallOption struct {
	RoomExe string
}

func LoadHallOption(filename string) (*HallOption, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	op := HallOption{}
	err = json.Unmarshal(d, &op)
	if err != nil {
		return nil, err
	}

	return &op, nil
}

func GenerateDefaultHallOption(filename string) error {
	d, err := json.Marshal(generateDefaultHallOption())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, d, 777)
}

func generateDefaultHallOption() *HallOption {
	return &HallOption{
		RoomExe: "./room.exe",
	}
}
