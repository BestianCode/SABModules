package main

import (
	"log"

	// SQLite
	_ "github.com/mattn/go-sqlite3"

	"github.com/BestianRU/SABModules/SBMConnect"
	"github.com/BestianRU/SABModules/SBMSystem"
)

func main() {
	var (
		jsonConfig  SBMSystem.ReadJSONConfig
		sqlite      SBMConnect.SQLite
		sqlite_init = string(`
PRAGMA journal_mode=WAL;
CREATE TABLE IF NOT EXISTS test (
	test1 varchar(255) PRIMARY KEY,
	test2 varchar(255)
);
			`)
	)

	jsonConfig.Init()

	if sqlite.Init(jsonConfig, sqlite_init) != 0 {
		return
	}
	defer sqlite.Close()

	_, err := sqlite.D.Exec("insert into test (test1, test2) values ('X1','X2');")
	if err != nil {
		log.Printf("SQLite::Query() Insert error: %v\n", err)
		return
	}
}
