package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/shubham88fru/distributed-key-value-store-go/config"
	"github.com/shubham88fru/distributed-key-value-store-go/db"
	"github.com/shubham88fru/distributed-key-value-store-go/web"
)

var (
	dbLocation  = flag.String("db-location", "", "Path to the bolt db")
	httpAddr    = flag.String("http-addr", "localhost:8080", "Host and port for the HTTP server")
	shardConfig = flag.String("shard-config", "shard-config.toml", "Path to the (static) shard config file")
	shard       = flag.String("shard", "", "shard name")
)

func parseFlags() {
	flag.Parse()
	if *dbLocation == "" {
		log.Fatal("db-location is required.")
	}

	if *shard == "" {
		log.Fatal("shard is required.")
	}
}

func main() {
	parseFlags()

	var c *config.Config
	c, err := config.ParseConfigFile(*shardConfig)

	if err != nil {
		log.Fatalf("Failed to parse shard config file: %v", err)
	}

	shards, err := config.ParseShards(c.Shards, *shard)

	if err != nil {
		log.Fatalf("Failed to parse shard config: %v", err)
	}

	log.Println("Total shards: ", shards.Count, " Shard index: ", shards.CurrIdx)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer close()

	server := web.NewServer(db, shards)

	http.HandleFunc("/get", server.GetHandler)
	http.HandleFunc("/set", server.SetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
