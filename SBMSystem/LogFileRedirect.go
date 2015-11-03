package SBMSystem

import (
	"log"
	"os"
)

const (
	LLError = iota
	LLWarning
	LLInfo
	LLTrace
)

type LogFile struct {
	LL      int
	flog    *os.File
	lineLog *log.Logger
}

func (_s *LogFile) ON(conf ReadJSONConfig) {
	var err error

	_s.flog, err = os.OpenFile(conf.Conf.LOG_File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error open log file: %s (%v)\n", conf.Conf.LOG_File, err)
	}

	if conf.Conf.LogLevel > 0 && conf.Conf.LogLevel <= 3 {
		_s.LL = conf.Conf.LogLevel
	} else {
		_s.LL = 0
	}

	_s.lineLog = log.New(_s.flog, "", log.Ldate|log.Ltime)

	if _s.LL > 1 {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
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

func (_s *LogFile) Log(msg string) {
	_s.lineLog.Println(msg)
}

func (_s *LogFile) Hello(pName, pVer string) {
	_s.lineLog.Printf(".")
	_s.lineLog.Printf(">")
	_s.lineLog.Printf("-> %s V%s", pName, pVer)
	_s.lineLog.Printf("--> Go!")
}

func (_s *LogFile) Bye() {
	_s.lineLog.Printf("--> To Sleep...")
	_s.lineLog.Printf("->")
	_s.lineLog.Printf(">")
	_s.lineLog.Printf(".")
}
