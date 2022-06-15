package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/patriciabonaldy/sports-news/internal/config"
	"github.com/patriciabonaldy/sports-news/internal/platform/logger"
)

const (
	lsbnTZ = "Europe/Lisbon"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.New()
	loc, err := time.LoadLocation(lsbnTZ)
	if err != nil {
		loc = time.Local
	}

	c := cron.New(cron.WithLocation(loc))
	if err = sync(cfg, c, logger); err != nil {
		log.Fatal(err)
	}

	c.Start()
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
