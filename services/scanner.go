package services

import (
	"apertoire.net/moviebase/bus"
	"apertoire.net/moviebase/helper"
	"apertoire.net/moviebase/model"
	"log"
)

type Scanner struct {
	Bus    *bus.Bus
	Config helper.Config
}

func (self *Scanner) Start() {
	log.Printf("starting scanner service ...")

	go self.react()

	log.Printf("scanner service started")
}

func (self *Scanner) Stop() {
	// nothing right now
}

func (self *Scanner) react() {
	for {
		select {
		case msg := <-self.Bus.MovieScan:
			go self.doMovieScan(msg.Payload, msg.Reply)
		}
	}
}

func (self *Scanner) doMovieScan(user *model.MovieScanReq, reply chan *model.MovieScanRep) {
	log.Printf("i got here")
}
