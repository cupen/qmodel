package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	ID            uint64    `bson:"id"`
	CreatedAt     time.Time `bson:"createdAt"`
	Field_uint8   uint8
	Field_uint16  uint16
	Field_uint32  uint32
	Field_int64   int64
	Field_string  string
	Field_float32 float32
	Field_float64 float64
}

func (t *TestModel) GetID() interface{} {
	return t.ID
}

func (t *TestModel) GetCollection() string {
	return "mongodb-test-basemodel"
}

func testConn() *Connection {
	url := "mongodb://root:0a8q09l0Q9bXlrUl@127.0.0.1:27017/test_qmgo?authSource=admin"
	conn, err := New(url, 8*time.Second)
	if err != nil {
		panic(err)
	}
	return conn
}

func TestConnection(t *testing.T) {
	assert := assert.New(t)
	conn := testConn()
	c := conn.UseCollection("TestConnection")
	c.DropCollection(context.TODO())

	uid := uint64(time.Now().Unix())
	m := TestModel{}
	m.ID = uid
	m.Field_uint8 = 0xff
	m.Field_uint16 = 0xffff
	m.Field_uint32 = 0xffffffff
	m.Field_int64 = 0xffffffff << 31
	{
		resp, err := c.InsertOne(context.TODO(), &m)
		assert.NoError(err)
		assert.NotEmpty(resp.InsertedID)
	}

	m2 := TestModel{}
	err := c.Find(context.TODO(), qmgo.M{"id": uid}).One(&m2)
	assert.NoError(err)
	assert.Equal(uint8(0xff), m.Field_uint8)
	assert.Equal(uint16(0xffff), m.Field_uint16)
	assert.Equal(uint32(0xffffffff), m.Field_uint32)
	assert.Equal(int64(0xffffffff<<31), m.Field_int64)
}
