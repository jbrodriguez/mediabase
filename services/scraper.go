package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"github.com/apertoire/go-tmdb"
	"log"
)

type Scraper struct {
	Bus      *bus.Bus
	Config   *helper.Config
	tmdb     *tmdb.Tmdb
	workpool *helper.Workpool
}

func (self *Scraper) Start() {
	log.Println("starting scraper service ...")

	var err error
	self.tmdb, err = tmdb.NewClient("e610ded10c3f47d05fe797961d90fea6", false)
	if err != nil {
		log.Fatal(err)
	}

	go self.react()

	self.workpool = helper.NewWorkpool(10, 1)
	go self.workpool.Balance()

	log.Println("scraper service started")
}

func (self *Scraper) Stop() {
	log.Printf("scraper service stopped")
}

func (self *Scraper) react() {
	for {
		select {
		case msg := <-self.Bus.ScrapeMovie:
			// self.doScrapeMovie(msg)
			go self.requestWork(msg)
		}
	}
}

func (self *Scraper) requestWork(movie *message.Movie) {
	c := make(chan *message.Media)

	self.workpool.Work <- helper.Request{self.scrapeMovie, &message.Media{"", "", "", movie}, c}
	media := <-c

	self.Bus.MovieScraped <- media
}

func (self *Scraper) scrapeMovie(media *message.Media) *message.Media {
	log.Printf("before searchmovie %s", media.Movie.Title)
	res, err := self.tmdb.SearchMovie(media.Movie.Title)
	if err != nil {
		log.Println(err)
		return nil
	}

	if res.Total_Results != 1 {
		log.Println("more than one")
	}

	id := res.Results[0].Id

	log.Printf("before getmovie [%d] %s", id, media.Movie.Title)
	gmr, err := self.tmdb.GetMovie(id)
	if err != nil {
		log.Println(err)
		return nil
	}

	media.Movie.Original_Title = gmr.Original_Title
	media.Movie.Runtime = gmr.Runtime
	media.Movie.Tmdb_Id = gmr.Id
	media.Movie.Imdb_Id = gmr.Imdb_Id
	media.Movie.Overview = gmr.Overview
	media.Movie.Tagline = gmr.Tagline
	media.Movie.Cover = gmr.Poster_Path
	media.Movie.Backdrop = gmr.Backdrop_Path

	media.BaseUrl = self.tmdb.BaseUrl
	media.SecureBaseUrl = self.tmdb.SecureBaseUrl

	return media
	// self.Bus.MovieScraped <- &message.Media{self.tmdb.BaseUrl, self.tmdb.SecureBaseUrl, "", movie}
}

// func (self *Scraper) doScrapeMovie(movie *message.Movie) {
// 	res, err := self.tmdb.SearchMovie(movie.Title)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	if res.Total_Results != 1 {
// 		log.Println("more than one")
// 	}

// 	id := res.Results[0].Id

// 	gmr, err := self.tmdb.GetMovie(id)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	movie.Original_Title = gmr.Original_Title
// 	movie.Runtime = gmr.Runtime
// 	movie.Tmdb_Id = gmr.Id
// 	movie.Imdb_Id = gmr.Imdb_Id
// 	movie.Overview = gmr.Overview
// 	movie.Tagline = gmr.Tagline
// 	movie.Cover = gmr.Poster_Path
// 	movie.Backdrop = gmr.Backdrop_Path

// 	self.Bus.MovieScraped <- &message.Media{self.tmdb.BaseUrl, self.tmdb.SecureBaseUrl, "", movie}
// }
