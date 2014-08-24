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
	logger := services.Logger{Bus: &bus, Config: &config}
	dal := services.Dal{Bus: &bus, Config: &config}
	server := services.Server{Bus: &bus, Config: &config}
	scanner := services.Scanner{Bus: &bus}
	scraper := services.Scraper{Bus: &bus, Config: &config}
	pruner := services.Pruner{Bus: &bus, Config: &config}
	cache := services.Cache{Bus: &bus, Config: &config}
	core := services.Core{Bus: &bus, Config: &config}

	bus.Start()
	logger.Start()
	dal.Start()
	server.Start()
	scanner.Start()
	scraper.Start()
	pruner.Start()
	cache.Start()
	core.Start()

	log.Printf("press enter to stop ...")
	var input string
	fmt.Scanln(&input)

	core.Stop()
	cache.Stop()
	pruner.Stop()
	scraper.Stop()
	scanner.Stop()
	server.Stop()
	dal.Stop()
	logger.Stop()
	// bus.Stop()
}
