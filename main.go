package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	bolt "go.etcd.io/bbolt"
)

var (
	dbLocation = flag.String("db-location", "", "Path to the bolt db")
	httpAddr   = flag.String("http-addr", "localhost:8080", "Host and port for the HTTP server")
)

func parseFlags() {
	flag.Parse()
	if *dbLocation == "" {
		log.Fatal("db-location is required.")
	}
}

func main() {
	parseFlags()

	db, err := bolt.Open(*dbLocation, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET request received")
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("SET request received")
	})

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
