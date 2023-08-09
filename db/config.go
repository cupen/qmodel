package db

import "fmt"

type Config struct {
	KeySpace string `toml:"keyspace"`
	MongoURL string `toml:"mongo_url"`
	RedisURL string `toml:"redis_url"`
	CacheURL string `toml:"cache_url"`
	QueueURL string `toml:"queue_url"`
}

func (c *Config) Clone() *Config {
	cloned := *(c)
	clonedPtr := &cloned
	if clonedPtr == c {
		panic(fmt.Errorf("cloned failed: dbv2.Config"))
	}
	return clonedPtr
}
