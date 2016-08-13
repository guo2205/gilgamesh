// database
package config

import (
	"encoding/json"
	"io/ioutil"
)

type Sqlite3DatabaseOption struct {
	Filename string
}

type MssqlDatabaseOption struct {
	Addr     string
	Port     int
	User     string
	Password string
	Name     string
}

type DatabaseOption struct {
	Type    string
	Sqlite3 Sqlite3DatabaseOption
	Mssql   MssqlDatabaseOption
}

func LoadDatabaseOption(filename string) (*DatabaseOption, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	op := DatabaseOption{}
	err = json.Unmarshal(d, &op)
	if err != nil {
		return nil, err
	}

	return &op, nil
}

func GenerateDefaultDatabaseOption(filename string) error {
	d, err := json.Marshal(generateDefaultGateOption())
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, d, 777)
}

func generateDefaultDatabaseOption() *DatabaseOption {
	return &DatabaseOption{
		Type: "sqlite3",
		Sqlite3: Sqlite3DatabaseOption{
			Filename: "./db/db3.sqlite3",
		},
		Mssql: MssqlDatabaseOption{
			Addr:     "192.168.2.2",
			Port:     1433,
			User:     "sa",
			Password: "lwzw611116",
			Name:     "ygo",
		},
	}
}
