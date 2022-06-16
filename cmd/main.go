package main

import (
	"github.com/patriciabonaldy/sports-news/cmd/bootstrap"
	"log"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
