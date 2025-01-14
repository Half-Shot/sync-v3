package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	syncv3 "github.com/matrix-org/sync-v3"
	"github.com/matrix-org/sync-v3/sync2"
	"github.com/matrix-org/sync-v3/sync3"
)

var (
	flagDestinationServer = flag.String("server", "", "The destination v2 matrix server")
	flagBindAddr          = flag.String("port", ":8008", "Bind address")
	flagPostgres          = flag.String("db", "user=postgres dbname=syncv3 sslmode=disable", "Postgres DB connection string (see lib/pq docs)")
)

func main() {
	flag.Parse()
	if *flagDestinationServer == "" {
		flag.Usage()
		os.Exit(1)
	}
	// pprof
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			panic(err)
		}
	}()
	h, err := sync3.NewSync3Handler(&sync2.HTTPClient{
		Client: &http.Client{
			Timeout: 5 * time.Minute,
		},
		DestinationServer: *flagDestinationServer,
	}, *flagPostgres)
	if err != nil {
		panic(err)
	}
	syncv3.RunSyncV3Server(h, *flagBindAddr)
}
