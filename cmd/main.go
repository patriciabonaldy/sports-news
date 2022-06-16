package main

import (
	"log"

	"github.com/patriciabonaldy/sports-news/cmd/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
