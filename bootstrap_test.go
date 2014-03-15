package main

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	// "apertoire.net/mediabase/message"
	"apertoire.net/mediabase/services"
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"
)

type myGenre struct {
	id   int
	name string
}

func TestDb(t *testing.T) {
	log.Printf("numproc %d", runtime.NumCPU())

	tiempo := time.Now()
	// tiempo := time.Date(2013, time.December, 15,34,0,0,0, time.

	log.Printf("tiempo-bare: %s\n", tiempo)

	log.Printf("tiempo-fmt: %s\n", tiempo.Format(time.RFC3339))

	log.Printf("tiempo-fmt2: %s\n", tiempo.Format(time.RFC1123Z))

	test, err := time.Parse(time.RFC3339, "2013-12-14T16:18:59-05:00")
	if err != nil {
		panic(err)
	}

	log.Println(test.Format(time.RFC3339))

	// n := []myGenre{
	// 	{id: 1, name: "alfa"},
	// 	{id: 2, name: "beta"},
	// }

	// var genres string
	// for i := 0; i < len(n); i++ {
	// 	attr := &n[i]
	// 	if genres == "" {
	// 		genres = attr.name
	// 	} else {
	// 		genres += "|" + attr.name
	// 	}
	// }

	// log.Println(genres)

	log.Printf("starting up ...")

	config := helper.Config{}
	config.Init()

	bus := bus.Bus{}
	logger := services.Logger{Bus: &bus, Config: &config}
	dal := services.Dal{Bus: &bus, Config: &config}

	bus.Start()
	logger.Start()
	dal.Start()

	// bus.StoreMovie <- &message.Movie{Title: "september morning"}
	// bus.StoreMovie <- &message.Movie{Title: "remember how we danced"}
	// bus.StoreMovie <- &message.Movie{Title: "something happened"}
	// bus.StoreMovie <- &message.Movie{Title: "what can you do"}
	// bus.StoreMovie <- &message.Movie{Title: "stella"}
	// bus.StoreMovie <- &message.Movie{Title: "or else"}
	// bus.StoreMovie <- &message.Movie{Title: "find out about"}

	log.Printf("press enter to stop ...")
	var input string
	fmt.Scanln(&input)

	dal.Stop()
	logger.Stop()
	// bus.Stop()
}
