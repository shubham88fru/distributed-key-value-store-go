package config

type Shard struct {
	Name    string `json:"name"`
	Idx     int    `json:"idx"`
	Address string `json:"address"`
}

type Config struct {
	Shards []Shard `json:"shards"`
}
