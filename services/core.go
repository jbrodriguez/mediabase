package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/goinggo/tracelog"
	"log"
)

type Core struct {
	Bus    *bus.Bus
	Config *helper.Config
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
		case msg := <-self.Bus.MovieFound:
			go self.doMovieFound(msg)
		}
	}
}

func (self *Core) doMovieFound(movie *message.Movie) {
	log.Printf("found: %s (%s) [%s, %s, %s]", movie.Name, movie.Year, movie.Resolution, movie.Type, movie.Path)
	tracelog.INFO("mb", "core", fmt.Sprintf("found: %s (%s) [%s, %s, %s]", movie.Name, movie.Year, movie.Resolution, movie.Type, movie.Path))
	// calculate hex sha1 for the full movie path
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s|%s", movie.Name, movie.Year)))
	movie.Picture = hex.EncodeToString(h.Sum(nil)) + ".jpg"

	// go func() {
	// 	self.Bus.StoreMovie <- movie
	// }()

	// go func() {
	// 	self.Bus.CachePicture <- &message.Picture{Path: movie.Path, Id: movie.Picture}
	// }()

	self.Bus.StoreMovie <- movie

	self.Bus.CachePicture <- &message.Picture{Path: movie.Path, Id: movie.Picture, Name: movie.Name}
}
