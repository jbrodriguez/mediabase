package bus

import (
	"apertoire.net/moviebase/message"
	"log"
)

type Bus struct {
	UserAuth chan *message.UserAuth
	UserData chan *message.UserData
}

func (self *Bus) Start() {
	log.Println("bus starting up ...")
	self.UserAuth = make(chan *message.UserAuth)
	self.UserData = make(chan *message.UserData)
}
