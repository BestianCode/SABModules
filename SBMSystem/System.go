package SBMSystem

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func Exit(conf ReadJSONConfig, signalType os.Signal, pid PidFile) {
	var logRedirect LogFile

	logRedirect.ON(conf)
	defer logRedirect.OFF()

	log.Println(".")
	log.Println("..")
	log.Println("...")
	log.Println("Exit command received. Exiting...")
	log.Println("Signal type: ", signalType)
	log.Println("Bye...")
	log.Println("...")
	log.Println("..")
	log.Println(".")

	pid.OFF(conf)

	os.Exit(0)
}

func Fork(conf ReadJSONConfig) {
	var err error

	if conf.Daemon_mode != "YES" {
		return
	}

	err = exec.Command(os.Args[0], "-daemon=GO", fmt.Sprintf("-config=%s", conf.Config_file), " &").Start()
	if err != nil {
		fmt.Println("\tFork daemon error: %v\n\n\n", err)
		os.Exit(1)
	} else {
		fmt.Println("\tForked!\n\n\n")
		os.Exit(0)
	}
}

func Signal(conf ReadJSONConfig, pid PidFile) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		signalType := <-ch
		signal.Stop(ch)
		Exit(conf, signalType, pid)
	}()
}
