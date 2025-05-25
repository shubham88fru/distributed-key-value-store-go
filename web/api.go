package web

import (
	"fmt"
	"net/http"

	"github.com/shubham88fru/distributed-key-value-store-go/db"
)

type server struct {
	db *db.Database
}

func NewServer(db *db.Database) *server {
	return &server{
		db,
	}
}

func (s *server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.FormValue("key")
	value, err := s.db.GetKey(key)

	fmt.Fprintf(w, "Value is = %q, error = %v", value, err)
}

func (s *server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.FormValue("key")
	value := []byte(r.FormValue("value"))
	err := s.db.SetKey(key, value)

	fmt.Fprintf(w, "error = %v", err)
}
