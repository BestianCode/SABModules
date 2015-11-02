package SBMSystem

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
)

type PidFile struct {
	f *os.File
}

func (_s *PidFile) ON(conf ReadJSONConfig) {
	var err error

	_s._check(conf)

	_s.f, err = os.OpenFile(conf.Conf.PID_File, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Error create PID file: %s (%v)\n", conf.Conf.PID_File, err)
	}
	_s.f.WriteString(fmt.Sprintf("%d", os.Getpid()))
	_s.f.Close()
}

func (_s *PidFile) OFF(conf ReadJSONConfig) {
	var err error
	err = os.Remove(conf.Conf.PID_File)
	if err != nil {
		log.Fatalf("Error remove PID file: %s (%v)\n", conf.Conf.PID_File, err)
	}
}

func (_s PidFile) _check(conf ReadJSONConfig) {
	var err error

	_s.f, err = os.OpenFile(conf.Conf.PID_File, os.O_RDONLY, 0666)
	if err != nil {
		return
	}

	defer _s.f.Close()

	pid_read := make([]byte, 10)

	pid_bytes, err := _s.f.Read(pid_read)
	if err != nil {
		log.Printf("WR1/ > Remove old pid file")
		err = os.Remove(conf.Conf.PID_File)
		if err != nil {
			log.Fatalf("ER1/Error remove PID file: %s (%v)\n", conf.Conf.PID_File, err)
		}
		return
	}

	if pid_bytes > 0 {
		pid_read_int, err := strconv.Atoi(fmt.Sprintf("%s", pid_read[0:pid_bytes]))
		if err != nil {
			log.Printf("WR2/ > Remove old pid file")
			err = os.Remove(conf.Conf.PID_File)
			if err != nil {
				log.Fatalf("ER2/Error remove PID file: %s (%v)\n", conf.Conf.PID_File, err)
			}
			return
		}

		pid_proc, err := os.FindProcess(pid_read_int)
		if err != nil {
			log.Printf("WR3/ > Remove old pid file")
			err = os.Remove(conf.Conf.PID_File)
			if err != nil {
				log.Fatalf("ER3/Error remove PID file: %s (%v)\n", conf.Conf.PID_File, err)
			}
			return
		}

		err = pid_proc.Signal(syscall.Signal(0))
		if err != nil {
			log.Printf("WR4/ > Remove old pid file")
			err = os.Remove(conf.Conf.PID_File)
			if err != nil {
				log.Fatalf("ER4/Error remove PID file: %s (%v)\n", conf.Conf.PID_File, err)
			}
			return
		}

		log.Printf("<< ! Another copy of the program with PID %d is running! Exiting... ! >>", pid_read_int)
		os.Exit(1)
	} else {
		log.Printf("WR5/ > Remove old pid file")
		err = os.Remove(conf.Conf.PID_File)
		if err != nil {
			log.Fatalf("ER5/Error remove PID file: %s (%v)\n", conf.Conf.PID_File, err)
		}
		return
	}
}
