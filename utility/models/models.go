// init
package models

import (
	"fmt"
	"gilgamesh/utility/config"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	engine *xorm.Engine
)

func Init(option *config.DatabaseOption) {
	var (
		connString string
	)

	switch option.Type {
	case "sqlite3":
		connString = option.Sqlite3.Filename
	case "mssql":
		connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
			option.Mssql.Addr, option.Mssql.User, option.Mssql.Password, option.Mssql.Port, option.Mssql.Name)
	default:
		log.Fatalln("unknown database type", option.Type)
	}

	e, err := xorm.NewEngine(option.Type, connString)
	if err != nil {
		log.Fatalln(err)
	}
	engine = e
}

func Install() {
	tables := []interface{}{
		new(Account),
		new(Player),
		new(Deck),
		new(Videotape),
	}

	err := engine.DropTables(tables...)
	if err != nil {
		log.Fatalln(err)
	}

	err = engine.Sync2(tables...)
	if err != nil {
		log.Fatalln(err)
	}
}
