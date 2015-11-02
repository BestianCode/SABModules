package SBMSystem

import (
	"log"
	"os"
)

type LogFile struct {
	flog *os.File
}

func (_s *LogFile) ON(conf ReadJSONConfig) {
	var err error

	_s.flog, err = os.OpenFile(conf.Conf.LOG_File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error open log file: %s (%v)\n", conf.Conf.LOG_File, err)
	}

	log.SetOutput(_s.flog)
}

func (_s *LogFile) OFF() {
	var err error
	err = _s.flog.Close()
	if err != nil {
		log.Fatalf("Error close log file: (%v)\n", err)
	}
}
