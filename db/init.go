package db

import "sync/atomic"

var (
	instance atomic.Value
)

func Init(cfg *Config) {
	db := New(cfg)
	instance.Store(db)
}

func Global() *DB {
	return instance.Load().(*DB)
}
