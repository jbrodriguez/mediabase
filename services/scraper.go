package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"fmt"
	"github.com/apertoire/go-tmdb"
	"github.com/goinggo/tracelog"
	"github.com/goinggo/workpool"
	"log"
)

type Scraper struct {
	Bus      *bus.Bus
	Config   *helper.Config
	tmdb     *tmdb.Tmdb
	workpool *workpool.WorkPool
}

func (self *Scraper) Start() {
	log.Println("starting scraper service ...")

	var err error
	self.tmdb, err = tmdb.NewClient("e610ded10c3f47d05fe797961d90fea6", false)
	if err != nil {
		log.Fatal(err)
	}

	self.workpool = workpool.New(12, 2000)

	go self.react()

	// go self.workpool.Balance()

	log.Println("scraper service started")
}

func (self *Scraper) Stop() {
	self.workpool.Shutdown("scraper")
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
	tracelog.TRACE("mb", "scraper", "WORK REQUESTED [%s]", movie.Title)

	c := make(chan *message.Media)

	gig := &Gig{
		self.Bus,
		self.tmdb,
		&message.Media{"", "", "", movie},
		c,
	}

	self.workpool.PostWork("gig", gig)

	// tracelog.TRACE("mb", "scraper", "[%s] RUNNING SCRAPING [%s]", movie.Title)
	media := <-c

	// tracelog.TRACE("mb", "scraper", "[%s] FINISHED SCRAPING", media.Movie.Title)
	self.Bus.MovieScraped <- media

	tracelog.TRACE("mb", "scraper", "WORK COMPLETED [%s]", movie.Title)
}

type Gig struct {
	bus   *bus.Bus
	tmdb  *tmdb.Tmdb
	media *message.Media
	ret   chan *message.Media
}

func (self *Gig) DoWork(workRoutine int) {
	var result *message.Media

	defer func() {
		self.ret <- result
	}()

	result = self.media

	tracelog.TRACE("mb", "scraper", "STARTED SCRAPING [%s]", self.media.Movie.Title)
	res, err := self.tmdb.SearchMovie(self.media.Movie.Title)
	if err != nil {
		log.Println(err)
		return
	}

	if res.Total_Results == 0 {
		tracelog.TRACE("mb", "scraper", fmt.Sprintf("TMDB: NO MATCH FOUND [%s]", self.media.Movie.Title))
		return
	} else if res.Total_Results > 1 {
		tracelog.TRACE("mb", "scraper", fmt.Sprintf("TMDB: MORE THAN ONE [%s]", self.media.Movie.Title))
	}

	id := res.Results[0].Id

	// log.Printf("before getmovie [%d] %s", id, media.Movie.Title)
	// tracelog.TRACE("mb", "scraper", "[%s] before getmovie [%s]", self.media.Movie.Title)
	gmr, err := self.tmdb.GetMovie(id)
	if err != nil {
		log.Println(err)
		return
	}

	self.media.Movie.Original_Title = gmr.Original_Title
	self.media.Movie.Runtime = gmr.Runtime
	self.media.Movie.Tmdb_Id = gmr.Id
	self.media.Movie.Imdb_Id = gmr.Imdb_Id
	self.media.Movie.Overview = gmr.Overview
	self.media.Movie.Tagline = gmr.Tagline
	self.media.Movie.Cover = gmr.Poster_Path
	self.media.Movie.Backdrop = gmr.Backdrop_Path

	self.media.BaseUrl = self.tmdb.BaseUrl
	self.media.SecureBaseUrl = self.tmdb.SecureBaseUrl

	tracelog.TRACE("mb", "scraper", "FINISHED SCRAPING [%s]", self.media.Movie.Title)
	// return media
	// self.Bus.MovieScraped <- &message.Media{self.tmdb.BaseUrl, self.tmdb.SecureBaseUrl, "", movie}
}

// func (self *Scraper) scrapeMovie(media *message.Media) *message.Media {
// 	tracelog.INFO("mb", "scraper", "before searchmovie %s", media.Movie.Title)
// 	res, err := self.tmdb.SearchMovie(media.Movie.Title)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	if res.Total_Results != 1 {
// 		log.Println("more than one")
// 	}

// 	id := res.Results[0].Id

// 	// log.Printf("before getmovie [%d] %s", id, media.Movie.Title)
// 	tracelog.INFO("mb", "scraper", "before gethmovie %s", media.Movie.Title)
// 	gmr, err := self.tmdb.GetMovie(id)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}

// 	media.Movie.Original_Title = gmr.Original_Title
// 	media.Movie.Runtime = gmr.Runtime
// 	media.Movie.Tmdb_Id = gmr.Id
// 	media.Movie.Imdb_Id = gmr.Imdb_Id
// 	media.Movie.Overview = gmr.Overview
// 	media.Movie.Tagline = gmr.Tagline
// 	media.Movie.Cover = gmr.Poster_Path
// 	media.Movie.Backdrop = gmr.Backdrop_Path

// 	media.BaseUrl = self.tmdb.BaseUrl
// 	media.SecureBaseUrl = self.tmdb.SecureBaseUrl

// 	tracelog.INFO("mb", "scraper", "before finalizing %s", media.Movie.Title)
// 	return media
// 	// self.Bus.MovieScraped <- &message.Media{self.tmdb.BaseUrl, self.tmdb.SecureBaseUrl, "", movie}
// }

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
