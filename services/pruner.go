package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"fmt"
	"github.com/goinggo/tracelog"
	"log"
	"os"
	// "path/filepath"
	// "regexp"
	// "strings"
)

type Pruner struct {
	Bus    *bus.Bus
	Config *helper.Config
}

func (self *Pruner) Start() {
	log.Printf("starting Pruner service ...")

	go self.react()

	log.Printf("Pruner service started")
}

func (self *Pruner) Stop() {
	// nothing right now
	log.Printf("Pruner service stopped")
}

func (self *Pruner) react() {
	for {
		select {
		case msg := <-self.Bus.PruneMovies:
			go self.doPruneMovies(msg.Reply)
		}
	}
}

func (self *Pruner) doPruneMovies(reply chan string) {
	tracelog.TRACE("mb", "pruner", "Looking for something to prune")

	msg := message.ListMovies{make(chan []*message.Movie)}
	self.Bus.ListMovies <- &msg
	items := <-msg.Reply

	for _, item := range items {

		if _, err := os.Stat(item.Location); err != nil {
			if os.IsNotExist(err) {
				tracelog.TRACE("mb", "pruner", fmt.Sprintf("UP FOR DELETION: [%d] %s (%s))", item.Id, item.Title, item.Location))
				self.Bus.DeleteMovie <- item
			}
		}

	}
}
