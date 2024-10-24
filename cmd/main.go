package main

import (
	"log"

	"github.com/Toolnado/alligator/cache"
	"github.com/Toolnado/alligator/server"
)

func main() {
	cche := cache.New()
	opts := server.Options{
		Addr:  ":3000",
		Cache: cche,
	}
	svr := server.New(opts, true)
	log.Fatal(svr.ListenAndServe())
}
