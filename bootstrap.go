package main

import (
	// "dal"
	"apertoire.net/moviebase/server"
	"fmt"
	"log"
)

func main() {
	log.Printf("starting up ...")

	go server.Start()
	// go dal.Start()

	log.Printf("press enter to stop ...")
	var input string
	fmt.Scanln(&input)
}
