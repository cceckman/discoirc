package main

import (
	"log"
	"os"

	backend "github.com/cceckman/discoirc/backend/girc"
	"github.com/lrstanley/girc"
)

func main() {
	b := backend.NewNetwork("localhost", &girc.Config{
		Server: "localhost",
		Port:   6667,
		Nick:   "cceckman",
		Name:   "Charles Eckman",
		User:   "cceckman",
		Debug:  os.Stdout,
	})

	if err := b.Connect(); err != nil {
		log.Fatal(err)
	}
}
