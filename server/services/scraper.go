package services

import (
	"apertoire.net/mediabase/server/bus"
	"apertoire.net/mediabase/server/helper"
	"apertoire.net/mediabase/server/message"
	"github.com/apertoire/go-tmdb"
	"github.com/apertoire/mlog"
	"github.com/goinggo/workpool"
	"time"
)

type Scraper struct {
	Bus      *bus.Bus
	Config   *helper.Config
	tmdb     *tmdb.Tmdb
	workpool *workpool.WorkPool
}

func (self *Scraper) Start() {
	mlog.Info("starting scraper service ...")

	var err error
	self.tmdb, err = tmdb.NewClient("e610ded10c3f47d05fe797961d90fea6", false)
	if err != nil {
		mlog.Fatalf("unable to create tmdb client: %s", err)
	}

	self.workpool = workpool.New(12, 4000)

	go self.react()

	// go self.workpool.Balance()

	mlog.Info("scraper service started")
}

func (self *Scraper) Stop() {
	self.workpool.Shutdown("scraper")
	mlog.Info("scraper service stopped")
}

func (self *Scraper) react() {
	for {
		select {
		case msg := <-self.Bus.ScrapeMovie:
			// self.doScrapeMovie(msg)
			go self.requestWork(msg)

		case msg := <-self.Bus.RescrapeMovies:
			go self.fixMoviesWork(msg)
		}
	}
}

func (self *Scraper) fixMoviesWork(movies []*message.Movie) {
	mlog.Info("FIX MOVIES WORK REQUESTED FOR [%d] movies", len(movies))

	c := make(chan *message.Media)

	for i := range movies {
		gig := &FixMovieGig{
			self.Bus,
			self.tmdb,
			&message.Media{BaseUrl: "", SecureBaseUrl: "", BasePath: "", Movie: movies[i], Forced: true},
			c,
		}

		self.workpool.PostWork("fixMovieGig", gig)

		// mlog.Info("[%s] RUNNING SCRAPING [%s]", movie.Title)
		media := <-c

		// mlog.Info("[%s] FINISHED SCRAPING", media.Movie.Title)
		self.Bus.MovieRescraped <- media
	}

	mlog.Info("FIX MOVIES WORK COMPLETED FOR [%d]", len(movies))
}

func (self *Scraper) requestWork(movie *message.Movie) {
	mlog.Info("WORK REQUESTED [%s]", movie.Title)

	c := make(chan *message.Media)

	gig := &Gig{
		self.Bus,
		self.tmdb,
		&message.Media{BaseUrl: "", SecureBaseUrl: "", BasePath: "", Movie: movie, Forced: false},
		c,
	}

	self.workpool.PostWork("gig", gig)

	// mlog.Info("[%s] RUNNING SCRAPING [%s]", movie.Title)
	media := <-c

	// mlog.Info("[%s] FINISHED SCRAPING", media.Movie.Title)
	self.Bus.MovieScraped <- media

	mlog.Info("WORK COMPLETED [%s]", movie.Title)
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

	mlog.Info("STARTED SCRAPING [%s]", self.media.Movie.Title)
	movies, err := self.tmdb.SearchMovie(self.media.Movie.Title)
	if err != nil {
		mlog.Error(err)
		return
	}

	if movies.Total_Results == 0 {
		mlog.Info("TMDB: NO MATCH FOUND [%s]", self.media.Movie.Title)
		return
	} else if movies.Total_Results > 1 {
		mlog.Info("TMDB: MORE THAN ONE [%s]", self.media.Movie.Title)
	}

	id := movies.Results[0].Id

	// log.Printf("before getmovie [%d] %s", id, media.Movie.Title)
	// mlog.Info("[%s] before getmovie [%s]", self.media.Movie.Title)
	gmr, err := self.tmdb.GetMovie(id)
	if err != nil {
		mlog.Info("FAILED GETTING MOVIE [%s]", self.media.Movie.Title)
		return
	}

	self.media.Movie.Title = gmr.Title
	self.media.Movie.Original_Title = gmr.Original_Title
	self.media.Movie.Runtime = gmr.Runtime
	self.media.Movie.Tmdb_Id = gmr.Id
	self.media.Movie.Imdb_Id = gmr.Imdb_Id
	self.media.Movie.Overview = gmr.Overview
	self.media.Movie.Tagline = gmr.Tagline
	self.media.Movie.Cover = gmr.Poster_Path
	self.media.Movie.Backdrop = gmr.Backdrop_Path

	for i := 0; i < len(gmr.Genres); i++ {
		attr := &gmr.Genres[i]
		if self.media.Movie.Genres == "" {
			self.media.Movie.Genres = attr.Name
		} else {
			self.media.Movie.Genres += " " + attr.Name
		}
	}

	self.media.Movie.Vote_Average = gmr.Vote_Average
	self.media.Movie.Vote_Count = gmr.Vote_Count

	for i := 0; i < len(gmr.Production_Countries); i++ {
		attr := &gmr.Production_Countries[i]
		if self.media.Movie.Production_Countries == "" {
			self.media.Movie.Production_Countries = attr.Name
		} else {
			self.media.Movie.Production_Countries += "|" + attr.Name
		}
	}

	now := time.Now().Format(time.RFC3339)
	self.media.Movie.Added = now
	self.media.Movie.Modified = now

	self.media.BaseUrl = self.tmdb.BaseUrl
	self.media.SecureBaseUrl = self.tmdb.SecureBaseUrl

	mlog.Info("FINISHED SCRAPING [%s]", self.media.Movie.Title)
	// return media
	// self.Bus.MovieScraped <- &message.Media{self.tmdb.BaseUrl, self.tmdb.SecureBaseUrl, "", movie}
}

type FixMovieGig struct {
	bus   *bus.Bus
	tmdb  *tmdb.Tmdb
	media *message.Media
	ret   chan *message.Media
}

func (self *FixMovieGig) DoWork(workRoutine int) {
	var result *message.Media

	defer func() {
		self.ret <- result
	}()

	mlog.Info("FIXMOVIE: STARTED SCRAPING [%s]", self.media.Movie.Title)
	result = self.media

	id := self.media.Movie.Tmdb_Id

	// log.Printf("before getmovie [%d] %s", id, media.Movie.Title)
	// mlog.Info("[%s] before getmovie [%s]", self.media.Movie.Title)
	gmr, err := self.tmdb.GetMovie(id)
	if err != nil {
		mlog.Info("FIXMOVIE: FAILED GETTING MOVIE [%s]", self.media.Movie.Title)
		return
	}

	self.media.Movie.Title = gmr.Title
	self.media.Movie.Original_Title = gmr.Original_Title
	self.media.Movie.Runtime = gmr.Runtime
	self.media.Movie.Tmdb_Id = gmr.Id
	self.media.Movie.Imdb_Id = gmr.Imdb_Id
	self.media.Movie.Overview = gmr.Overview
	self.media.Movie.Tagline = gmr.Tagline
	self.media.Movie.Cover = gmr.Poster_Path
	self.media.Movie.Backdrop = gmr.Backdrop_Path

	for i := 0; i < len(gmr.Genres); i++ {
		attr := &gmr.Genres[i]
		if self.media.Movie.Genres == "" {
			self.media.Movie.Genres = attr.Name
		} else {
			self.media.Movie.Genres += " " + attr.Name
		}
	}

	self.media.Movie.Vote_Average = gmr.Vote_Average
	self.media.Movie.Vote_Count = gmr.Vote_Count

	for i := 0; i < len(gmr.Production_Countries); i++ {
		attr := &gmr.Production_Countries[i]
		if self.media.Movie.Production_Countries == "" {
			self.media.Movie.Production_Countries = attr.Name
		} else {
			self.media.Movie.Production_Countries += "|" + attr.Name
		}
	}

	now := time.Now().Format(time.RFC3339)
	self.media.Movie.Modified = now

	self.media.BaseUrl = self.tmdb.BaseUrl
	self.media.SecureBaseUrl = self.tmdb.SecureBaseUrl

	mlog.Info("FIXMOVIE: FINISHED SCRAPING [%s]", self.media.Movie.Title)
	// return media
	// self.Bus.MovieScraped <- &message.Media{self.tmdb.BaseUrl, self.tmdb.SecureBaseUrl, "", movie}
}
