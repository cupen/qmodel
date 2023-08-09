package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func (t *TestModel) DBIndexes() []Index {
	return []Index{
		{
			Key:         []string{"id", "createdAt"},
			Background:  true,
			ExpireAfter: 1 * time.Hour,
		},
	}
}

func TestIndex(t *testing.T) {
	obj := &TestModel{ID: 2}
	conn := testConn()

	cleanup := func() { conn.DeleteM(obj) }
	t.Cleanup(cleanup)
	conn.DeleteM(obj)

	t.Run("DropCollection", func(t *testing.T) {
		assert := assert.New(t)
		err := conn.GetDatabase().DropDatabase(context.TODO())
		assert.NoError(err)
	})

	t.Run("CreateIndex/AutoCreateCollection", func(t *testing.T) {
		assert := assert.New(t)
		err := conn.CreateIndex(&TestModel{})
		assert.NoError(err)
	})
}
