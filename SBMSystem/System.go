package SBMSystem

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func Exit(conf ReadJSONConfig, signalType os.Signal, pid PidFile) {
	var rLog LogFile

	rLog.ON(conf)
	defer rLog.OFF()

	rLog.Log(".")
	rLog.Log("..")
	rLog.Log("...")
	rLog.Log("Exit command received. Exiting...")
	rLog.Log("Signal type: ", signalType)
	rLog.Log("Bye...")
	rLog.Log("...")
	rLog.Log("..")
	rLog.Log(".")

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
