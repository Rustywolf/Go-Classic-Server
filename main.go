package main

import (
	"classicserver/classic"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	server, err := classic.NewClassicServer()
	if err != nil {
		log.Fatalln(err)
	} else {
		server.Run()
	}
}
