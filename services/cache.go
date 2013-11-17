package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"log"
	"os"
	"path/filepath"
)

type Cache struct {
	Bus    *bus.Bus
	Config *helper.Config
}

func (self *Cache) Start() {
	log.Println("starting cache service ...")

	go self.react()

	log.Println("cache service started")
}

func (self *Cache) Stop() {

}

func (self *Cache) react() {
	for {
		select {
		case msg := <-self.Bus.CachePicture:
			go self.doCachePicture(msg)
		}
	}
}

func (self *Cache) doCachePicture(picture *message.Picture) {
	picPath := filepath.Join(self.Config.AppDir, "/web/img/", picture.Id)
	if _, err := os.Stat(picPath); err == nil {
		log.Printf("SKIP: picture in cache: %s", picPath)
		return
	}

	ext := filepath.Ext(picture.Path)
	name := picture.Path[0:len(picture.Path)-len(ext)] + ".jpg"

	err := helper.Copy(name, picPath)
	if err != nil {
		log.Printf("ERR: couldn't copy %s", name)
		return
	}

	log.Printf("INFO: copied %s to %s", name, picPath)
}
