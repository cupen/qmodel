package mongo

import (
	"context"
	"errors"
	urllib "net/url"
	"time"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	M = qmgo.M
	E = qmgo.E
)

type Connection struct {
	dbName  string
	session *qmgo.Session
	client  *qmgo.Client
	db      *qmgo.Database
}

// @See https://docs.mongodb.com/manual/reference/connection-string/
// Usage:
//
//	Connect("mongodb://user:passwd@127.0.0.1:27017/db")
func New(url string, timeout time.Duration) (*Connection, error) {
	ctx := context.Background()
	timeoutMS := int64(timeout / time.Millisecond)
	client, err := qmgo.NewClient(ctx, &qmgo.Config{
		Uri:              url,
		ConnectTimeoutMS: &timeoutMS,
		SocketTimeoutMS:  &timeoutMS,
	})
	if err != nil {
		return nil, err
	}
	u, _ := urllib.Parse(url)
	conn := Connection{}
	conn.dbName = u.EscapedPath()[1:]
	session, err := client.Session()
	if err != nil {
		return nil, err
	}
	conn.client = client
	conn.session = session
	conn.db = client.Database(conn.dbName)
	return &conn, nil
}

func (conn *Connection) Close() {
	conn.client.Close(context.TODO())
}

func (conn *Connection) GetDatabase() *qmgo.Database {
	return conn.db
}

func (conn *Connection) GetDatabaseName() string {
	return conn.dbName
}

func (conn *Connection) UseCollection(collection string) *qmgo.Collection {
	return conn.db.Collection(collection)
}

// func (conn *Connection) UseCollectionWithTimeout(collection string, timeout time.Duration) *qmgo.Collection {
// 	return conn.db.Collection(collection)
// }

// Ping ...
func (conn *Connection) Ping() error {
	c := conn.UseCollection("_ping_")
	ctx := context.Background()
	_, err := c.UpsertId(ctx, "ping", bson.M{"ping": "pong", "updateAt": time.Now()})
	if err != nil {
		return err
	}

	rs := map[string]interface{}{}
	err = c.Find(ctx, bson.M{"_id": "ping"}).One(&rs)
	if err != nil {
		return err
	}
	if rs["ping"].(string) != "pong" {
		return errors.New("ping failed")
	}
	return nil
}
