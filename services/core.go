package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
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
		case msg := <-self.Bus.CachePicture:
			go self.doCachePicture(msg)
		}
	}
}

func (self *Core) doMovieFound(movie *message.Movie) {
	log.Printf("found: %s (%s) [%s, %s, %s]", movie.Name, movie.Year, movie.Resolution, movie.Type, movie.Path)
	self.Bus.StoreMovie <- movie

	self.Bus.CachePicture <- &message.Picture{Path: movie.Path}
}

func (self *Core) doCachePicture(picture *message.Picture) {
	log.Println("estoy dentro de cachepicture")

	ext := filepath.Ext(picture.Path)
	name := picture.Path[0 : len(picture.Path)-len(ext)]

	h := sha1.New()
	h.Write([]byte(picture.Path))
	picture.Id = hex.EncodeToString(h.Sum(nil)) + ".jpg"

	picPath := filepath.Join(self.Config.AppDir, "/web/img/", picture.Id)
	log.Printf("picpath: %s", picPath)

	if _, err := os.Stat(picPath); err == nil {
		log.Printf("err: %s", picPath)
		return
	}

	helper.Copy(filepath.Join(name, ".jpg"), picPath)

	self.Bus.UpdateMovie <- picture
}
