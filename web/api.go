package web

import (
	"fmt"
	"hash/fnv"
	"net/http"

	"github.com/shubham88fru/distributed-key-value-store-go/db"
)

type server struct {
	db       *db.Database
	shardIdx int
	shards   int
}

func NewServer(db *db.Database, shardIdx, shards int) *server {
	return &server{
		db,
		shardIdx,
		shards,
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

	h := fnv.New64()
	h.Write([]byte(key))
	destShard := h.Sum64() % uint64(s.shards)
	fmt.Fprintf(w, "error = %v, hash = %d, destination shard = %d", err, h.Sum64(), destShard)
}
