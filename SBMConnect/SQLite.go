package SBMConnect

import (
	"database/sql"
	"log"

	// SQLite
	_ "github.com/mattn/go-sqlite3"

	"github.com/BestianRU/SABModules/SBMSystem"
)

type SQLite struct {
	D *sql.DB
}

func (_s *SQLite) Init(conf SBMSystem.ReadJSONConfig, initDB string) int {
	var err error

	_s.D, err = sql.Open("sqlite3", conf.Conf.SQLite_DB)
	if err != nil {
		log.Printf("SQLite::Open() error: %v\n", err)
		return -1
	}

	err = _s.D.Ping()
	if err != nil {
		log.Printf("SQLite::Ping() error: %v\n", err)
		return -1
	}

	if len(initDB) > 10 {
		_, err = _s.D.Exec(initDB)
		if err != nil {
			log.Printf("SQLite::Query() InitDB error: %v\n", err)
			return -1
		}
	}

	return 0
}

func (_s *SQLite) Close() {
	_s.D.Close()
}
