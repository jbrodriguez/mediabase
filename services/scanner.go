package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/model"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`/Volumes/films/(?P<Resolution>.*?)/(?P<Name>.*?)/(?:.*/)*.*\.(?P<FileType>bdmv|iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms)$`)

type Scanner struct {
	Bus    *bus.Bus
	Config helper.Config

	re [3]*regexp.Regexp
}

func (self *Scanner) Start() {
	log.Printf("starting scanner service ...")

	log.Printf("compiling regular expressions ...")

	// test:="I am leaving from home in a while"
	// prepositionsRegex := make([]*regexp.Regexp, len(preps))
	// for i := 0; i < len(preps); i++ {
	// prepositionsRegex[i]=regexp.MustCompile(`\b`+preps[i]+`\b`)
	// }

	// for i := 0; i < len(prepositionsRegex); i++ {
	// fmt.Println(prepositionsRegex[i].String())
	// if loc := prepositionsRegex[i].FindStringIndex(test); loc != nil{
	// fmt.Println(test[loc[0]:loc[1]], "found at: ", loc[0])
	// break

	self.re[0] = regexp.MustCompile(`/volumes/films/(?P<Resolution>.*?)/(?P<Name>.*?)/(?:.*/)*bdmv/index.bdmv$`)
	self.re[1] = regexp.MustCompile(`/volumes/films/(?P<Resolution>.*?)/(?P<Name>.*?)/(?:.*/)*.*\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms)$`)
	self.re[2] = regexp.MustCompile(`/volumes/films/(?P<Resolution>.*?)/(?P<Name>.*?)/(?:.*/)*(?:video_ts|hv000i01)\.(?P<FileType>ifo)$`)

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

func (self *Scanner) visit(path string, f os.FileInfo, err error) error {
	for i := 0; i < 3; i++ {
		match := self.re[i].FindStringSubmatch(strings.ToLower(path))
		if match == nil {
			continue
		}

		log.Printf("p: %s", path)
		return nil
	}

	return nil
}

func (self *Scanner) doMovieScan(user *model.MovieScanReq, reply chan *model.MovieScanRep) {
	log.Printf("i got here")

	err := filepath.Walk("/Volumes/films", self.visit)
	log.Println("err: %s", err)

}
