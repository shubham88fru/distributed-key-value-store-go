package config

type Shard struct {
	Name string `json:"name"`
	Idx  int    `json:"idx"`
}

type Config struct {
	Shards []Shard `json:"shards"`
}
