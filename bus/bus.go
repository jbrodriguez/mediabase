package bus

import (
	"apertoire.net/mediabase/message"
	"log"
)

type Bus struct {
	MovieScan chan *message.MovieScan
}

func (self *Bus) Start() {
	log.Println("bus starting up ...")
	self.MovieScan = make(chan *message.MovieScan)
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
