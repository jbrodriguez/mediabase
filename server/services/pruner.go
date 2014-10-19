package services

import (
	"apertoire.net/mediabase/server/bus"
	"apertoire.net/mediabase/server/message"
	"apertoire.net/mediabase/server/model"
	"github.com/apertoire/mlog"
	"os"
)

type Pruner struct {
	Bus    *bus.Bus
	Config *model.Config
}

func (self *Pruner) Start() {
	mlog.Info("starting Pruner service ...")

	go self.react()

	mlog.Info("Pruner service started")
}

func (self *Pruner) Stop() {
	// nothing right now
	mlog.Info("Pruner service stopped")
}

func (self *Pruner) react() {
	for {
		select {
		case msg := <-self.Bus.PruneMovies:
			go self.doPruneMovies(msg.Reply)
		}
	}
}

func (self *Pruner) ConfigChanged(conf *model.Config) {
	self.Config = conf
}

func (self *Pruner) doPruneMovies(reply chan string) {
	mlog.Info("Looking for something to prune")

	options := message.Options{Current: 0, Limit: 99999999999999, SortBy: "title", SortOrder: "asc"}
	msg := message.Movies{Options: options, Reply: make(chan *message.MoviesDTO)}
	self.Bus.GetMovies <- &msg
	dto := <-msg.Reply

	for _, item := range dto.Movies {
		if _, err := os.Stat(item.Location); err != nil {
			if os.IsNotExist(err) {
				mlog.Info("UP FOR DELETION: [%d] %s (%s))", item.Id, item.Title, item.Location)
				self.Bus.DeleteMovie <- item
			}
		}

	}

	reply <- "Tudo Bom"
}
