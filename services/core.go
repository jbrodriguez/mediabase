package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	// // "crypto/sha1"
	// // "encoding/hex"
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
	log.Printf("core service stopped")
}

func (self *Core) react() {
	for {
		select {
		case msg := <-self.Bus.MovieFound:
			go self.doMovieFound(msg)
		case msg := <-self.Bus.MovieScraped:
			go self.doMovieScraped(msg)
		case msg := <-self.Bus.MovieRescraped:
			go self.doMovieRescraped(msg)
		case msg := <-self.Bus.FixMovies:
			go self.doFixMovies(msg)
		}
	}
}

func (self *Core) doMovieFound(movie *message.Movie) {
	// log.Printf("found: %s (%s) [%s, %s, %s]", movie.Title, movie.Year, movie.Resolution, movie.FileType, movie.Location)
	// tracelog.INFO("mb", "core", fmt.Sprintf("found: %s (%s) [%s, %s, %s]", movie.Title, movie.Year, movie.Resolution, movie.FileType, movie.Location))
	// calculate hex sha1 for the full movie path
	// h := sha1.New()
	// h.Write([]byte(fmt.Sprintf("%s|%s", movie.Title, movie.Year)))
	// movie.Picture = hex.EncodeToString(h.Sum(nil)) + ".jpg"

	// go func() {
	// 	self.Bus.StoreMovie <- movie
	// }()

	// go func() {
	// 	self.Bus.CachePicture <- &message.Picture{Path: movie.Path, Id: movie.Picture}
	// }()

	c := make(chan bool)

	self.Bus.CheckMovie <- &message.CheckMovie{movie, c}
	exists := <-c

	if exists {
		tracelog.TRACE("mb", "core", fmt.Sprintf("SKIPPED: present in db [%s] (%s)", movie.Title, movie.Location))
		return
	}

	self.Bus.ScrapeMovie <- movie

	// self.Bus.StoreMovie <- movie

	// self.Bus.CachePicture <- &message.Picture{Path: movie.Path, Id: movie.Picture, Title: movie.Title}
}

func (self *Core) doMovieScraped(media *message.Media) {
	go func() {
		tracelog.TRACE("mb", "core", fmt.Sprintf("STORING MOVIE [%s]", media.Movie.Title))
		self.Bus.StoreMovie <- media.Movie
	}()

	go func() {
		tracelog.TRACE("mb", "core", fmt.Sprintf("CACHING MEDIA [%s]", media.Movie.Title))
		media.BasePath = self.Config.AppDir
		self.Bus.CacheMedia <- media
	}()
}

func (self *Core) doMovieRescraped(media *message.Media) {
	go func() {
		tracelog.TRACE("mb", "core", fmt.Sprintf("UPDATING MOVIE [%s]", media.Movie.Title))
		self.Bus.UpdateMovie <- media.Movie
	}()

	go func() {
		tracelog.TRACE("mb", "core", fmt.Sprintf("CACHING MEDIA [%s]", media.Movie.Title))
		media.BasePath = self.Config.AppDir
		self.Bus.CacheMedia <- media
	}()
}

func (self *Core) doFixMovies(flag int) {
	msg := message.Movies{make(chan []*message.Movie)}
	self.Bus.GetMoviesToFix <- &msg

	tracelog.TRACE("mb", "core", fmt.Sprintf("AFTER GET MOVIES TO FIX [%v]", msg.Reply))

	reply := <-msg.Reply

	tracelog.TRACE("mb", "core", fmt.Sprintf("WAITING FOR REPLY [%v]", reply))

	self.Bus.RescrapeMovies <- reply
}
