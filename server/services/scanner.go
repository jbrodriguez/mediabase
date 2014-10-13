package services

import (
	"apertoire.net/mediabase/server/bus"
	"apertoire.net/mediabase/server/helper"
	"apertoire.net/mediabase/server/message"
	"apertoire.net/mediabase/server/model"
	"github.com/apertoire/mlog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Scanner struct {
	Bus    *bus.Bus
	Config *model.Config

	// re           [3]*helper.Rexp
	re           []*helper.Rexp
	includedMask string
}

func (self *Scanner) Start() {
	mlog.Info("starting scanner service ...")

	mlog.Info("compiling regular expressions ...")

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

	self.re = make([]*helper.Rexp, 0)

	for _, regex := range self.Config.MediaRegexs {
		self.re = append(self.re, &helper.Rexp{Exp: regexp.MustCompile(regex)})
	}

	// self.re[0] = &helper.Rexp{Exp: regexp.MustCompile(`(?i)/volumes/.*?/(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*bdmv/index.(?P<FileType>bdmv)$`)}
	// self.re[1] = &helper.Rexp{Exp: regexp.MustCompile(`(?i)/volumes/.*?/(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*.*\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv)$`)}
	// self.re[2] = &helper.Rexp{Exp: regexp.MustCompile(`(?i)/volumes/.*?/(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*(?:video_ts|hv000i01)\.(?P<FileType>ifo)$`)}

	self.includedMask = ".bdmv|.iso|.img|.nrg|.mkv|.avi|.xvid|.ts|.mpg|.dvr-ms|.mdf|.wmv|.ifo"

	go self.react()

	mlog.Info("scanner service started")
}

func (self *Scanner) Stop() {
	// nothing right now
	mlog.Info("scanner service stopped")
}

func (self *Scanner) react() {
	for {
		select {
		case msg := <-self.Bus.ScanMovies:
			go self.doScanMovies(msg.Reply)
		}
	}
}

func (self *Scanner) doScanMovies(reply chan string) {
	mlog.Info("inside ScanMovies")

	reply <- "Movie scannning process started ..."

	for _, folder := range self.Config.MediaFolders {
		err := self.walker(folder)
		if err != nil {
			mlog.Info("Unable to scan movies: %s", err)
		}

		mlog.Info("Completed scan of folder: %s", folder)
	}

	self.Bus.ImportMoviesFinished <- 1
}

func (self *Scanner) walker(folder string) error {

	if folder[len(folder)-1] != '/' {
		folder = folder + "/"
	}

	err := filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			mlog.Info("from-start err: %s", err)
		}

		// mlog.Info("maldito: %s", path)

		if !strings.Contains(self.includedMask, strings.ToLower(filepath.Ext(path))) {
			// mlog.Info("[%s] excluding %s", filepath.Ext(path), path)
			return nil
		}

		comparePath := strings.TrimPrefix(path, folder)
		// mlog.Info("folder: %s, comparePath: %s", folder, comparePath)
		// TODO: remove folder from path to match against rexp

		for i := 0; i < 3; i++ {
			// match := self.re[i].FindStringSubmatch(strings.ToLower(path))
			// if match == nil {
			// 	continue
			// }
			var rmap = self.re[i].Match(comparePath)
			if rmap == nil {
				continue
			}

			var resolution string
			var ok bool
			if resolution, ok = rmap["Resolution"]; !ok {
				resolution = ""
			}

			movie := &message.Movie{Title: rmap["Name"], File_Title: rmap["Name"], Year: rmap["Year"], Resolution: resolution, FileType: rmap["FileType"], Location: path}
			mlog.Info("FOUND [%s] (%s)", movie.Title, movie.Location)

			self.Bus.MovieFound <- movie

			return nil
		}

		return nil
	})

	return err
}

// err := filepath.Walk("/Volumes/hal-films", self.visit)
// if err != nil {
// 	mlog.Info("err: %s", err)
// }

// mlog.Info("completed scannning hal for movies")

// err = filepath.Walk("/Volumes/wopr-films", self.visit)
// if err != nil {
// 	mlog.Info("err: %s", err)
// }

// mlog.Info("completed scannning wopr for movies")

// self.Bus.ImportMoviesFinished <- 1
// }

// func (self *Scanner) visit(path string, f os.FileInfo, err error) error {
// 	if err != nil {
// 		mlog.Info("from-start err: %s", err)
// 	}
//
// 	// mlog.Info("maldito: %s", path)
//
// 	if !strings.Contains(self.includedMask, strings.ToLower(filepath.Ext(path))) {
// 		// mlog.Info("[%s] excluding %s", filepath.Ext(path), path)
// 		return nil
// 	}
//
// 	for i := 0; i < 3; i++ {
// 		// match := self.re[i].FindStringSubmatch(strings.ToLower(path))
// 		// if match == nil {
// 		// 	continue
// 		// }
// 		var rmap = self.re[i].Match(path)
// 		if rmap == nil {
// 			continue
// 		}
//
// 		movie := &message.Movie{Title: rmap["Name"], File_Title: rmap["Name"], Year: rmap["Year"], Resolution: rmap["Resolution"], FileType: rmap["FileType"], Location: path}
// 		mlog.Info("FOUND [%s] (%s)", movie.Title, movie.Location)
//
// 		self.Bus.MovieFound <- movie
//
// 		return nil
// 	}
//
// 	return nil
// }
