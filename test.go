package main

import (
	"github.com/BestianRU/SABModules/SBMSystem"
	"log"
)

func main() {
	var (
		jsonConfig  SBMSystem.ReadJSONConfig
		logRedirect SBMSystem.LogFile
		pid         SBMSystem.PidFile
	)

	jsonConfig.Init()

	pid.ON(jsonConfig)
	defer pid.OFF(jsonConfig)

	logRedirect.ON(jsonConfig)
	defer logRedirect.OFF()

	log.Printf("%v\n", jsonConfig.Conf)
}
