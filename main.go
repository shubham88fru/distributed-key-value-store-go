package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/shubham88fru/distributed-key-value-store-go/db"
	"github.com/shubham88fru/distributed-key-value-store-go/web"
)

var (
	dbLocation  = flag.String("db-location", "", "Path to the bolt db")
	httpAddr    = flag.String("http-addr", "localhost:8080", "Host and port for the HTTP server")
	shardConfig = flag.String("shard-config", "shard-config.toml", "Path to the (static) shard config file")
)

func parseFlags() {
	flag.Parse()
	if *dbLocation == "" {
		log.Fatal("db-location is required.")
	}
}

func main() {
	parseFlags()

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer close()

	server := web.NewServer(db)

	http.HandleFunc("/get", server.GetHandler)
	http.HandleFunc("/set", server.SetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
