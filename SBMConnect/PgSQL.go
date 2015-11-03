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

func (_s *PgSQL) Init(conf SBMSystem.ReadJSONConfig) int {
	var err error
	_s.D, err = sql.Open("postgres", conf.Conf.PG_DSN)
	if err != nil {
		log.Printf("PG::Open() error: %v\n", err)
		return -1
	}

	return 0
}

func (_s *PgSQL) Close() {
	_s.D.Close()
}
