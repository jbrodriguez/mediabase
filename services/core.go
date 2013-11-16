package services

import (
	"apertoire.net/mediabase/bus"
	// "apertoire.net/mediabase/helper"
	"log"
)

type Core struct {
	Bus *bus.Bus
}

func (self *Core) Start() {
	log.Printf("starting core service ...")

	// some initialization

	go self.react()

	log.Printf("core service started")
}

func (self *Core) Stop() {
	// some deinitialization
}

func (self *Core) react() {
	for {
		select {
		case msg := <-self.Bus.MediaScanned:
			// go self.doAuthenticate(msg.Payload, msg.Reply)
		}
	}
}
