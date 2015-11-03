package SBMConnect

import (
	"database/sql"
	"log"

	// PostgreSQL
	_ "github.com/lib/pq"

	"github.com/BestianRU/SABModules/SBMSystem"
)

type PgSQL struct {
	D *sql.DB
}

func (_s *PgSQL) Init(conf SBMSystem.ReadJSONConfig, initDB string) int {
	var err error

	_s.D, err = sql.Open("postgres", conf.Conf.PG_DSN)
	if err != nil {
		log.Printf("PgSQL::Open() error: %v\n", err)
		return -1
	}

	if len(initDB) > 10 {
		_, err = _s.D.Query(initDB)
		if err != nil {
			log.Printf("PgSQL::Query() InitDB error: %v\n", err)
			return -1
		}
	}

	return 0
}

func (_s *PgSQL) Close() {
	_s.D.Close()
}
