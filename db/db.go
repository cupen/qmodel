package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cupen/qmodel/db/mongo"

	"github.com/go-redis/redis/v8"
	"github.com/qiniu/qmgo"
)

type DB struct {
	cfg   *Config
	Mongo *mongo.Connection
	Redis *redis.Client
	Cache *redis.Client
	Queue *redis.Client
}

func New(cfg *Config) *DB {
	mongo, err := mongo.New(cfg.MongoURL, 8*time.Second)
	if err != nil {
		panic(err)
	}
	return &DB{
		cfg:   cfg,
		Mongo: mongo,
		Redis: newRedis(cfg.RedisURL),
		Cache: newRedis(cfg.CacheURL),
		Queue: newRedis(cfg.QueueURL),
	}
}

func (d *DB) BuildKey(names ...string) string {
	baseKey := d.cfg.KeySpace
	return strings.Join(append([]string{baseKey}, names...), ":")
}

func (d *DB) GetKeySpace() string {
	return d.cfg.KeySpace
}

func (d *DB) GetConfig() *Config {
	return d.cfg.Clone()
}

func (d *DB) Ping() error {
	err1 := d.Mongo.Ping()
	err2 := d.Redis.Ping(context.TODO()).Err()
	var show string
	if err1 != nil {
		show = fmt.Sprintf("mongo: %v", err1)
	}
	if err2 != nil {
		show += fmt.Sprintf(" redis: %v", err2)
	}
	if show != "" {
		return errors.New(show)
	}
	return nil
}

func (d *DB) Clone() *DB {
	cloned := (*d)
	if &cloned == d {
		panic(fmt.Errorf("clone db object failed"))
	}
	return &cloned
}

func newRedis(url string) *redis.Client {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opt)
}

func IsErrNotFound(err error) bool {
	return err == qmgo.ErrNoSuchDocuments
}
