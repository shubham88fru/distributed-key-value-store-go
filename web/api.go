package web

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/http"

	"github.com/shubham88fru/distributed-key-value-store-go/db"
)

type server struct {
	db       *db.Database
	shardIdx int
	shards   int
	addrs    map[int]string
}

func NewServer(db *db.Database, shardIdx, shards int, addrs map[int]string) *server {
	return &server{
		db,
		shardIdx,
		shards,
		addrs,
	}
}

func (s *server) redirect(w http.ResponseWriter, r *http.Request, destShard int) {
	log.Printf("Redirecting to shard %d", destShard)
	http.Redirect(w, r, "http://"+s.addrs[destShard]+r.RequestURI, http.StatusSeeOther)
}

func (s *server) GetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetHandler called on shard", s.shardIdx)

	r.ParseForm()
	key := r.FormValue("key")

	destShard := s.getShard(key)
	if s.shardIdx != destShard { //redirect to the appropriate shard
		s.redirect(w, r, destShard)
		return
	}

	value, err := s.db.GetKey(key)
	fmt.Fprintf(w, "Current shard = %d, Key shard = %d, Value is = %q, Error = %v", s.shardIdx, destShard, value, err)
}

func (s *server) SetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SetHandler called on shard", s.shardIdx)

	r.ParseForm()
	key := r.FormValue("key")
	value := []byte(r.FormValue("value"))
	err := s.db.SetKey(key, value)

	destShard := s.getShard(key)
	if s.shardIdx != destShard { //redirect to the appropriate shard
		s.redirect(w, r, destShard)
		return
	}

	fmt.Fprintf(w, "Current shard = %d, Destination shard = %d, Error = %v", s.shardIdx, destShard, err)
}

func (s *server) getShard(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.shards))
}
