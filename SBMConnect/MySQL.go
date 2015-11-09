package SBMConnect

import (
	"database/sql"
	"fmt"
	"log"

	// MySQL
	_ "github.com/ziutek/mymysql/godrv"

	"github.com/BestianRU/SABModules/SBMSystem"
)

type MySQL struct {
	D *sql.DB
}

func (_s *MySQL) Init(conf SBMSystem.ReadJSONConfig, initDB string) int {
	var err error

	_s.D, err = sql.Open("mymysql", conf.Conf.MY_DSN)
	if err != nil {
		log.Printf("MySQL::Open() error: %v\n", err)
		return -1
	}

	if len(initDB) > 10 {
		_, err = _s.D.Query(initDB)
		if err != nil {
			log.Printf("MySQL::Query() InitDB error: %v\n", err)
			return -1
		}
	}

	return 0
}

func (_s *MySQL) Close() {
	_s.D.Close()
}

func (_s *MySQL) QSimple(query ...interface{}) int {
	var (
		i       int
		queryGo = string("")
	)
	for _, x := range query {
		queryGo = fmt.Sprintf("%s%v", queryGo, x)
	}
	res, err := _s.D.Query(queryGo)
	if err != nil {
		log.Printf("SQL::Query() QSimple error: %v\n", err)
		log.Printf("Query: %s\n", queryGo)
		return -1
	}

	res.Next()
	res.Scan(&i)

	return i
}
