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
	server := services.Server{Bus: &bus, Config: &config}
	dal := services.Dal{Bus: &bus, Config: &config}
	cache := services.Cache{Bus: &bus, Config: &config}
	scanner := services.Scanner{Bus: &bus}
	core := services.Core{Bus: &bus, Config: &config}

	bus.Start()
	dal.Start()
	server.Start()
	cache.Start()
	scanner.Start()
	core.Start()

	log.Printf("press enter to stop ...")
	var input string
	fmt.Scanln(&input)

	core.Stop()
	scanner.Stop()
	cache.Stop()
	server.Stop()
	dal.Stop()
	// bus.Stop()
}
