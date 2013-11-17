package bus

import (
	"apertoire.net/mediabase/message"
	"log"
)

type Bus struct {
	ScanMovies   chan *message.ScanMovies
	MovieFound   chan *message.Movie
	StoreMovie   chan *message.Movie
	CachePicture chan *message.Picture
	UpdateMovie  chan *message.Picture
}

func (self *Bus) Start() {
	log.Println("bus starting up ...")

	self.ScanMovies = make(chan *message.ScanMovies)
	self.MovieFound = make(chan *message.Movie)

	self.StoreMovie = make(chan *message.Movie)
	self.CachePicture = make(chan *message.Picture)

	self.UpdateMovie = make(chan *message.Picture)
}
