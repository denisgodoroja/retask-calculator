package main

import (
	"fmt"

	"denisgodoroja/retask/config"
	"denisgodoroja/retask/webservice"
)

func main() {
	cfg, err := config.LoadConfigs()
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	webservice.Start(cfg)
}
