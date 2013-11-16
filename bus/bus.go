package bus

import (
	"apertoire.net/mediabase/message"
	"log"
)

type Bus struct {
	ScanMovie  chan *message.ScanMovie
	MovieFound chan *message.MovieFound
}

func (self *Bus) Start() {
	log.Println("bus starting up ...")

	self.ScanMovie = make(chan *message.ScanMovie)
	self.MovieFound = make(chan *message.MovieFound)
}

// type Bus struct {
// 	UserAuth chan *message.UserAuth
// 	UserData chan *message.UserData
// }

// func (self *Bus) Start() {
// 	log.Println("bus starting up ...")
// 	self.UserAuth = make(chan *message.UserAuth)
// 	self.UserData = make(chan *message.UserData)
// }
