package config

import (
	"fmt"
	"hash/fnv"

	"github.com/BurntSushi/toml"
)

type Shard struct {
	Name    string `json:"name"`
	Idx     int    `json:"idx"`
	Address string `json:"address"`
}

type Config struct {
	Shards []Shard `json:"shards"`
}

type Shards struct {
	Count   int            `json:"count"`
	CurrIdx int            `json:"curr_idx"`
	Addrs   map[int]string `json:"addrs"`
}

func ParseConfigFile(configFile string) (*Config, error) {
	var c Config
	if _, err := toml.DecodeFile(configFile, &c); err != nil {
		return &Config{}, err
	}

	return &c, nil
}

func ParseShards(shards []Shard, currShardName string) (*Shards, error) {
	shardCount := len(shards)
	shardIdx := -1
	addrs := make(map[int]string)

	for _, s := range shards {
		if _, ok := addrs[s.Idx]; ok {
			return nil, fmt.Errorf("duplicate shard index found: %d", s.Idx)
		}

		addrs[s.Idx] = s.Address
		if s.Name == currShardName {
			shardIdx = s.Idx
		}

	}

	for i := 0; i < shardCount; i++ {
		if _, ok := addrs[i]; !ok {
			return nil, fmt.Errorf("shard index %d not found in config", i)
		}
	}

	if shardIdx < 0 {
		return nil, fmt.Errorf("shard %s not found in config", currShardName)
	}

	return &Shards{
		Addrs:   addrs,
		Count:   shardCount,
		CurrIdx: shardIdx,
	}, nil

}

func (s *Shards) GetShard(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.Count))
}
