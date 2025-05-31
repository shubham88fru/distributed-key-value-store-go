package config

import "testing"

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
	_ = config

}
