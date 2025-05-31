package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shubham88fru/distributed-key-value-store-go/config"
	"github.com/shubham88fru/distributed-key-value-store-go/db"
)

type server struct {
	db     *db.Database
	shards *config.Shards
}

func NewServer(db *db.Database, s *config.Shards) *server {
	return &server{
		db:     db,
		shards: s,
	}
}

func (s *server) redirect(w http.ResponseWriter, r *http.Request, destShard int) {
	log.Printf("Redirecting to shard %d", destShard)
	http.Redirect(w, r, "http://"+s.shards.Addrs[destShard]+r.RequestURI, http.StatusSeeOther)
}

func (s *server) GetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetHandler called on shard", s.shards.CurrIdx)

	r.ParseForm()
	key := r.FormValue("key")

	destShard := s.shards.GetShard(key)
	if s.shards.CurrIdx != destShard { //redirect to the appropriate shard
		s.redirect(w, r, destShard)
		return
	}

	value, err := s.db.GetKey(key)
	fmt.Fprintf(w, "Current shard = %d, Key shard = %d, Value is = %q, Error = %v", s.shards.CurrIdx, destShard, value, err)
}

func (s *server) SetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SetHandler called on shard", s.shards.CurrIdx)

	r.ParseForm()
	key := r.FormValue("key")
	value := []byte(r.FormValue("value"))
	err := s.db.SetKey(key, value)

	destShard := s.shards.GetShard(key)
	if s.shards.CurrIdx != destShard { //redirect to the appropriate shard
		s.redirect(w, r, destShard)
		return
	}

	fmt.Fprintf(w, "Current shard = %d, Destination shard = %d, Error = %v", s.shards.CurrIdx, destShard, err)
}
