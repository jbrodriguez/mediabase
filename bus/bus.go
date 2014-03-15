package bus

import (
	"apertoire.net/mediabase/message"
	"log"
)

type Bus struct {
	ScanMovies  chan *message.ScanMovies
	ScrapeMovie chan *message.Movie
	PruneMovies chan *message.PruneMovies

	MovieFound   chan *message.Movie
	MovieScraped chan *message.Media

	GetMovies    chan *message.GetMovies
	ListMovies   chan *message.ListMovies
	SearchMovies chan *message.SearchMovies
	CheckMovie   chan *message.CheckMovie

	StoreMovie  chan *message.Movie
	DeleteMovie chan *message.Movie
	CacheMedia  chan *message.Media
	// UpdateMovie  chan *message.Picture
	Log chan string
}

func (self *Bus) Start() {
	log.Println("bus starting up ...")

	self.ScanMovies = make(chan *message.ScanMovies)
	self.ScrapeMovie = make(chan *message.Movie)
	self.PruneMovies = make(chan *message.PruneMovies)

	self.MovieFound = make(chan *message.Movie)
	self.MovieScraped = make(chan *message.Media)

	self.GetMovies = make(chan *message.GetMovies)
	self.ListMovies = make(chan *message.ListMovies)
	self.SearchMovies = make(chan *message.SearchMovies)
	self.CheckMovie = make(chan *message.CheckMovie)

	self.StoreMovie = make(chan *message.Movie)
	self.DeleteMovie = make(chan *message.Movie)
	self.CacheMedia = make(chan *message.Media)

	// self.UpdateMovie = make(chan *message.Picture)

	self.Log = make(chan string)
}

// type Msg struct {
// 	id int
// }

// channel = make(chan *[]Msg)

// cannot use make(chan *[]Msg, 0) (type chan *[]Msg) as type chan *[]Msg in assignment
