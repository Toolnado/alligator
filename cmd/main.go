package main

import (
	"flag"
	"log"

	"github.com/Toolnado/alligator/cache"
	"github.com/Toolnado/alligator/server"
)

func main() {
	var (
		addr       = flag.String("addr", ":3000", "listen address of the server")
		leaderAddr = flag.String("laddr", "", "listen address of the leader server")
		leader     = flag.Bool("l", true, "is the current server a leader")
	)
	opts := server.Options{
		Addr:       *addr,
		LeaderAddr: *leaderAddr,
		Leader:     *leader,
	}
	svr := server.New(opts, cache.New())
	log.Fatal(svr.ListenAndServe())
}
