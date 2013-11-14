package main

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/services"
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
	scanner := services.Scanner{Bus: &bus}

	bus.Start()
	dal.Start()
	server.Start()
	scanner.Start()

	log.Printf("press enter to stop ...")
	var input string
	fmt.Scanln(&input)

	scanner.Stop()
	server.Stop()
	dal.Stop()
	// bus.Stop()
}
