package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"github.com/goinggo/tracelog"
	"log"
)

type Logger struct {
	Bus    *bus.Bus
	Config *helper.Config
}

func (self *Logger) Start() {
	log.Println("starting logger service ...")

	tracelog.StartFile(tracelog.LEVEL_TRACE, "./mediabase.log", 1)

	go self.react()

	log.Println("logger service started")
}

func (self *Logger) Stop() {
	tracelog.Stop()
}

func (self *Logger) react() {
	for {
		select {
		case msg := <-self.Bus.Log:
			go self.doLog(msg)
		}
	}
}

func (self *Logger) doLog(msg string) {
	tracelog.INFO("mb", "lg", msg)
}
