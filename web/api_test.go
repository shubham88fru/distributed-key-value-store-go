package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/shubham88fru/distributed-key-value-store-go/config"
	"github.com/shubham88fru/distributed-key-value-store-go/db"
)

func createShardDb(t *testing.T, idx int) *db.Database {
	t.Helper()

	f, err := os.CreateTemp(os.TempDir(), "shard-"+strconv.Itoa(idx)+".db")
	if err != nil {
		t.Fatalf("Could not create temp file for shard %d: %v", idx, err)
	}
	fileName := f.Name()
	defer f.Close()
	t.Cleanup(func() { os.Remove(fileName) })

	d, close, err := db.NewDatabase(fileName)
	if err != nil {
		t.Fatalf("Could not create database for shard %d: %v", idx, err)
	}
	t.Cleanup(func() { close() })

	return d
}

func createShardServer(t *testing.T, idx int, addrs map[int]string) (*db.Database, *server) {
	t.Helper()

	d := createShardDb(t, idx)

	cfg := &config.Shards{
		Addrs:   addrs,
		Count:   len(addrs),
		CurrIdx: idx,
	}

	s := NewServer(d, cfg)
	return d, s
}

func TestWebServer(t *testing.T) {
	var shard1Get, shard1Set func(w http.ResponseWriter, r *http.Request)
	var shard2Get, shard2Set func(w http.ResponseWriter, r *http.Request)

	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/get") {
			shard1Get(w, r)
		} else if strings.HasPrefix(r.RequestURI, "/set") {
			shard1Set(w, r)
		}
	}))
	defer ts1.Close()

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/get") {
			shard2Get(w, r)
		} else if strings.HasPrefix(r.RequestURI, "/set") {
			shard2Set(w, r)
		}
	}))
	defer ts2.Close()

	addrs := map[int]string{
		0: strings.TrimPrefix(ts1.URL, "http://"),
		1: strings.TrimPrefix(ts2.URL, "http://"),
	}

	db1, web1 := createShardServer(t, 0, addrs)
	db2, web2 := createShardServer(t, 1, addrs)

	keys := map[string]int{
		"a": 0, //dest shard
		"b": 1, //dest shard
	}

	shard1Get = web1.GetHandler
	shard1Set = web1.SetHandler
	shard2Get = web2.GetHandler
	shard2Set = web2.SetHandler

	for key := range keys {
		_, err := http.Get(fmt.Sprintf("%s/set?key=%s&value=value-%s", ts1.URL, key, key))
		if err != nil {
			t.Fatalf("Failed to set key %s on shard 1: %v", key, err)
		}
	}

	for key := range keys {
		resp, err := http.Get(fmt.Sprintf("%s/get?key=%s", ts1.URL, key))
		if err != nil {
			t.Fatalf("Failed to get key %s from shard 1: %v", key, err)
		}

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body for key %s: %v", key, err)
		}

		if !bytes.Contains(contents, []byte(fmt.Sprintf("value-%s", key))) {
			t.Errorf("Expected value for key %s to contain 'value-%s', got '%s'", key, key, contents)
		}
	}

	val1, err := db1.GetKey("a")
	if err != nil {
		t.Fatalf("Failed to get key 'a' from shard 1: %v", err)
	}
	if !bytes.Equal(val1, []byte("value-a")) {
		t.Errorf("Expected value for key 'a' to be 'value-a', got '%s'", val1)
	}

	val2, err := db2.GetKey("b")
	if err != nil {
		t.Fatalf("Failed to get key 'b' from shard 2: %v", err)
	}
	if !bytes.Equal(val2, []byte("value-b")) {
		t.Errorf("Expected value for key 'b' to be 'value-b', got '%s'", val2)
	}
}
