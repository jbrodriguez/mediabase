package main

import (
	"apertoire.net/moviebase/bus"
	"apertoire.net/moviebase/helper"
	"apertoire.net/moviebase/services"
	"fmt"
	"log"
)

func main() {
	log.Printf("starting up ...")

	config := helper.Config{}
	config.Init()

	bus := bus.Bus{}
	server := services.Server{Bus: &bus, Config: config}
	dal := services.Dal{Bus: &bus}

	bus.Start()
	dal.Start()
	server.Start()

	log.Printf("press enter to stop ...")
	var input string
	fmt.Scanln(&input)

	server.Stop()
	dal.Stop()
}
