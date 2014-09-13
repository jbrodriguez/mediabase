package bus

import (
	"apertoire.net/mediabase/server/message"
	"github.com/apertoire/mlog"
)

type Bus struct {
	ImportMovies chan *message.Status
	ScanMovies   chan *message.ScanMovies
	ScrapeMovie  chan *message.Movie

	MovieFound     chan *message.Movie
	MovieScraped   chan *message.Media
	MovieRescraped chan *message.Media

	ImportMoviesFinished chan int

	PruneMovies chan *message.PruneMovies

	GetMovies      chan *message.GetMovies
	ListMovies     chan *message.ListMovies
	ShowDuplicates chan *message.Movies
	ListByRuntime  chan *message.Movies
	SearchMovies   chan *message.SearchMovies
	CheckMovie     chan *message.CheckMovie
	FixMovies      chan int
	GetMoviesToFix chan *message.Movies
	RescrapeMovies chan []*message.Movie

	WatchedMovie chan *message.SingleMovie
	FixMovie     chan *message.SingleMovie

	StoreMovie  chan *message.Movie
	DeleteMovie chan *message.Movie
	CacheMedia  chan *message.Media
	UpdateMovie chan *message.Movie
}

func (self *Bus) Start() {
	mlog.Info("bus starting up ...")

	self.ImportMovies = make(chan *message.Status)
	self.ScanMovies = make(chan *message.ScanMovies)
	self.ScrapeMovie = make(chan *message.Movie)

	self.MovieFound = make(chan *message.Movie)
	self.MovieScraped = make(chan *message.Media)
	self.MovieRescraped = make(chan *message.Media)

	self.ImportMoviesFinished = make(chan int)

	self.PruneMovies = make(chan *message.PruneMovies)

	self.GetMovies = make(chan *message.GetMovies)
	self.ListMovies = make(chan *message.ListMovies)
	self.ShowDuplicates = make(chan *message.Movies)
	self.ListByRuntime = make(chan *message.Movies)
	self.SearchMovies = make(chan *message.SearchMovies)
	self.CheckMovie = make(chan *message.CheckMovie)
	self.FixMovies = make(chan int)
	self.GetMoviesToFix = make(chan *message.Movies)
	self.RescrapeMovies = make(chan []*message.Movie)

	self.WatchedMovie = make(chan *message.SingleMovie)
	self.FixMovie = make(chan *message.SingleMovie)

	self.StoreMovie = make(chan *message.Movie)
	self.DeleteMovie = make(chan *message.Movie)
	self.CacheMedia = make(chan *message.Media)
	self.UpdateMovie = make(chan *message.Movie)

	mlog.Info("bus started ...")
}

// type Msg struct {
// 	id int
// }

// channel = make(chan *[]Msg)

// cannot use make(chan *[]Msg, 0) (type chan *[]Msg) as type chan *[]Msg in assignment
