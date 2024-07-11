package main

import (
	"log"

	"github.com/praveenmahasena/aiserver/internal"
)

func main() {
	if err := internal.Start(); err != nil {
		log.Fatalln(err)
	}
}
