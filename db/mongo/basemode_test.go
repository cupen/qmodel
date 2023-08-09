package mongo

import (
	"testing"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/stretchr/testify/assert"
)

func TestBaseModel_GetM(t *testing.T) {
	assert := assert.New(t)

	obj := &TestModel{ID: 1}
	conn := testConn()

	cleanup := func() { conn.DeleteM(obj) }
	t.Cleanup(cleanup)
	var err error
	conn.DeleteM(obj)
	t.Run("GetM", func(t *testing.T) {
		err = conn.GetM(obj)
		assert.Equal(qmgo.ErrNoSuchDocuments, err)
	})

	t.Run("CreateM", func(t *testing.T) {
		t.Cleanup(cleanup)
		obj.Field_int64 = 101
		err := conn.CreateM(obj)
		assert.NoError(err)
		{
			obj2 := &TestModel{ID: obj.ID}
			err = conn.GetM(obj2)
			assert.NoError(err)
			assert.Equal(obj, obj2)
		}
	})

	t.Run("UpdateM", func(t *testing.T) {
		t.Cleanup(cleanup)
		err = conn.GetM(obj)
		assert.Equal(qmgo.ErrNoSuchDocuments, err)
		err = conn.UpdateM(obj)
		assert.Equal(qmgo.ErrNoSuchDocuments, err)
		conn.CreateM(obj)

		obj.Field_int64 = time.Now().Unix()
		err = conn.UpdateM(obj)
		{
			assert.NoError(err)
			obj2 := &TestModel{ID: obj.ID}
			conn.GetM(obj2)
			assert.Equal(obj, obj2)
		}
	})

	t.Run("UpdateOrCreateM", func(t *testing.T) {
		t.Cleanup(cleanup)
		err = conn.GetM(obj)
		assert.Equal(qmgo.ErrNoSuchDocuments, err)
		isCreated, err := conn.UpdateOrCreateM(obj)
		assert.NoError(err)
		assert.True(isCreated, "未更新到，则创建之")

		isCreated, err = conn.UpdateOrCreateM(obj)
		assert.NoError(err)
		assert.False(isCreated, "已更新到，无需创建")

		isCreated, err = conn.UpdateOrCreateM(obj)
		assert.NoError(err)
		assert.False(isCreated, "已更新到，无需创建")
	})

}
