package main

import (
	"apertoire.net/mediabase/server/bus"
	"apertoire.net/mediabase/server/model"
	"apertoire.net/mediabase/server/services"
	"fmt"
	"github.com/apertoire/mlog"
)

var Version string

func main() {
	config := model.Config{}
	config.Init(Version)

	bus := bus.Bus{}
	dal := services.Dal{Bus: &bus, Config: &config}
	server := services.Server{Bus: &bus, Config: &config}
	scanner := services.Scanner{Bus: &bus, Config: &config}
	scraper := services.Scraper{Bus: &bus, Config: &config}
	pruner := services.Pruner{Bus: &bus, Config: &config}
	cache := services.Cache{Bus: &bus, Config: &config}
	core := services.Core{Bus: &bus, Config: &config}

	bus.Start()
	dal.Start()
	server.Start()
	scanner.Start()
	scraper.Start()
	pruner.Start()
	cache.Start()
	core.Start()

	// dal.ImportOmdb()

	mlog.Info("press enter to stop ...")
	var input string
	fmt.Scanln(&input)

	core.Stop()
	cache.Stop()
	pruner.Stop()
	scraper.Stop()
	scanner.Stop()
	server.Stop()
	dal.Stop()
	// bus.Stop()
}
