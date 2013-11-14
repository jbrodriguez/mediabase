package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/model"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var re = regexp.MustCompile(`/Volumes/films/(?P<Resolution>.*?)/(?P<Name>.*?)/(?:.*/)*.*\.(?P<FileType>bdmv|iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms)$`)

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

func visit(path string, f os.FileInfo, err error) error {
	log.Printf("p: %s", path)

	match := re.FindStringSubmatch(path)
	if len(match) != 3 {
		return nil
	}

	log.Printf("p: %s", path)
	return nil
}

func (self *Scanner) doMovieScan(user *model.MovieScanReq, reply chan *model.MovieScanRep) {
	log.Printf("i got here")

	err := filepath.Walk("/Volumes/films", visit)
	log.Println("err: %s", err)

}
