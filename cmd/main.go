package main

import (
	"log"

	"github.com/Toolnado/alligator/cache"
	"github.com/Toolnado/alligator/server"
)

func main() {
	instance := cache.New()
	opts := server.Options{
		Addr:  ":3000",
		Cache: instance,
	}
	svr := server.New(opts, true)
	log.Fatal(svr.ListenAndServe())
}
