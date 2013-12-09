package bus

import (
	"apertoire.net/mediabase/message"
	"log"
)

type Bus struct {
	ScanMovies  chan *message.ScanMovies
	ScrapeMovie chan *message.Movie

	MovieFound   chan *message.Movie
	MovieScraped chan *message.Media

	GetMovies    chan *message.GetMovies
	SearchMovies chan *message.SearchMovies

	StoreMovie chan *message.Movie
	CacheMedia chan *message.Media
	// UpdateMovie  chan *message.Picture
	Log chan string
}

func (self *Bus) Start() {
	log.Println("bus starting up ...")

	self.ScanMovies = make(chan *message.ScanMovies)
	self.ScrapeMovie = make(chan *message.Movie)

	self.MovieFound = make(chan *message.Movie)
	self.MovieScraped = make(chan *message.Media)

	self.GetMovies = make(chan *message.GetMovies)
	self.SearchMovies = make(chan *message.SearchMovies)

	self.StoreMovie = make(chan *message.Movie)
	self.CacheMedia = make(chan *message.Media)

	// self.UpdateMovie = make(chan *message.Picture)

	self.Log = make(chan string)
}

// type Msg struct {
// 	id int
// }

// channel = make(chan *[]Msg)

// cannot use make(chan *[]Msg, 0) (type chan *[]Msg) as type chan *[]Msg in assignment
