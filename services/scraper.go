package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"github.com/apertoire/go-tmdb"
	"log"
)

type Scraper struct {
	Bus    *bus.Bus
	Config *helper.Config
	tmdb   *tmdb.Tmdb
}

func (self *Scraper) Start() {
	log.Println("starting scraper service ...")

	var err error
	self.tmdb, err = tmdb.NewClient("e610ded10c3f47d05fe797961d90fea6")
	if err != nil {
		log.Fatal(err)
	}

	go self.react()

	log.Println("scraper service started")
}

func (self *Scraper) Stop() {

}

func (self *Scraper) react() {
	for {
		select {
		case msg := <-self.Bus.ScrapeMovie:
			self.doScrapeMovie(msg)
		}
	}
}

func (self *Scraper) doScrapeMovie(movie *message.Movie) {
	res, err := self.tmdb.SearchMovie(movie.Title)
	if err != nil {
		log.Println(err)
		return
	}

	if res.Total_Results != 1 {
		log.Println("more than one")
	}

	id := res.Results[0].Id

	gmr, err := self.tmdb.GetMovie(id)
	if err != nil {
		log.Println(err)
		return
	}

	movie.Original_Title = gmr.Original_Title
	movie.Runtime = gmr.Runtime
	movie.Tmdb_Id = gmr.Id
	movie.Imdb_Id = gmr.Imdb_Id
	movie.Overview = gmr.Overview
	movie.Tagline = gmr.Tagline
	movie.Cover = gmr.Poster_Path
	movie.Backdrop = gmr.Backdrop_Path

	self.Bus.CacheMedia <- &message.Media{self.tmdb.BaseUrl, self.tmdb.SecureBaseUrl, movie}
}
