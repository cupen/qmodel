package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{
		KeySpace: "db-test",
		MongoURL: "mongodb://root:00674q0TlsQfbxlEawUNSgIkmdzwuX2y@127.0.0.1:27017/test_qmgo?authSource=admin",
		RedisURL: "redis://127.0.0.1:6379/0",
		CacheURL: "redis://127.0.0.1:6379/1",
		QueueURL: "redis://127.0.0.1:6379/2",
	}
	db := New(&cfg)
	err := db.Ping()
	assert.NoError(err)
}
