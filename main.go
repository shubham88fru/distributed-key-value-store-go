package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/shubham88fru/distributed-key-value-store-go/db"
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

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer close()

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.FormValue("key")
		value, err := db.GetKey(key)

		fmt.Fprintf(w, "Value is = %q, error = %v", value, err)
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.FormValue("key")
		value := []byte(r.FormValue("value"))
		err := db.SetKey(key, value)

		fmt.Fprintf(w, "error = %v", err)
	})

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
