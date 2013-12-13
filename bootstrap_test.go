package main

import (
	// "apertoire.net/mediabase/bus"
	// "apertoire.net/mediabase/helper"
	// "apertoire.net/mediabase/message"
	// "apertoire.net/mediabase/services"
	// "fmt"
	"log"
	"runtime"
	"testing"
)

func TestDb(t *testing.T) {
	log.Printf("numproc %d", runtime.NumCPU())
	// log.Printf("starting up ...")

	// config := helper.Config{}
	// config.Init()

	// bus := bus.Bus{}
	// dal := services.Dal{Bus: &bus, Config: &config}

	// bus.Start()
	// dal.Start()

	// bus.StoreMovie <- &message.Movie{Title: "september morning"}
	// bus.StoreMovie <- &message.Movie{Title: "remember how we danced"}
	// // bus.StoreMovie <- &message.Movie{Title: "something happened"}
	// // bus.StoreMovie <- &message.Movie{Title: "what can you do"}
	// // bus.StoreMovie <- &message.Movie{Title: "stella"}
	// // bus.StoreMovie <- &message.Movie{Title: "or else"}
	// // bus.StoreMovie <- &message.Movie{Title: "find out about"}

	// log.Printf("press enter to stop ...")
	// var input string
	// fmt.Scanln(&input)

	// dal.Stop()
	// // bus.Stop()
}
