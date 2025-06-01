package config

import (
	"os"
	"reflect"
	"testing"
)

func createConfig(t *testing.T, configStr string) *Config {
	t.Helper()

	f, err := os.CreateTemp(os.TempDir(), "config.toml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer f.Close()

	name := f.Name()
	defer os.Remove(name)

	_, err = f.WriteString(configStr)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	c, err := ParseConfigFile(name)
	if err != nil {
		t.Fatalf("Failed to parse config file: %v", err)
	}

	return c

}

func TestConfigParse(t *testing.T) {
	config := `[[shards]]
	name = "shard-1"
	idx = 0
	address = "localhost:8081"

[[shards]]
	name = "shard-2"
	idx = 1
	address = "localhost:8082"
	`
	c := createConfig(t, config)

	want := Config{
		Shards: []Shard{
			{Name: "shard-1", Idx: 0, Address: "localhost:8081"},
			{Name: "shard-2", Idx: 1, Address: "localhost:8082"},
		},
	}

	if !reflect.DeepEqual(c, &want) {
		t.Errorf("Parsed config does not match expected value.\nGot: %+v\nWant: %+v", c, &want)
	}

}

func TestParseShards(t *testing.T) {
	config := `[[shards]]
	name = "shard-1"
	idx = 0
	address = "localhost:8081"

[[shards]]
	name = "shard-2"
	idx = 1
	address = "localhost:8082"
	`
	c := createConfig(t, config)

	shards, err := ParseShards(c.Shards, "shard-1")
	if err != nil {
		t.Fatalf("Failed to parse shards: %v", err)
	}

	want := &Shards{
		Addrs:   map[int]string{0: "localhost:8081", 1: "localhost:8082"},
		Count:   2,
		CurrIdx: 0,
	}

	if !reflect.DeepEqual(shards, want) {
		t.Errorf("Parsed shards do not match expected value.\nGot: %+v\nWant: %+v", shards, want)
	}
}
