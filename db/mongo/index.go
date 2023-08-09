package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/qiniu/qmgo/options"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

// 暂时只提供必要的参数，其他高级选项一般用不到，如确实需要可酌情添加。
type Index struct {
	Key         []string
	Background  bool
	Unique      bool
	Sparse      bool
	ExpireAfter time.Duration
}

// mgo 的索引结构体
// func (i *Index) AsMgo()  {
// }

// qmgo 的索引结构体
func (i *Index) AsQmgo() options.IndexModel {
	_eaSecs := int64(i.ExpireAfter / time.Second)
	if _eaSecs > (0xafffffff) {
		panic(fmt.Errorf("invalid ExpireAfterSeconds: %d", i.ExpireAfter))
	}
	eaSecs := int32(_eaSecs)
	if len(i.Key) <= 0 {
		panic(fmt.Errorf("invalid index.Key: %v", i.Key))
	}
	return options.IndexModel{
		Key: i.Key,
		IndexOptions: &opts.IndexOptions{
			Background:         &i.Background,
			Unique:             &i.Unique,
			Sparse:             &i.Sparse,
			ExpireAfterSeconds: &eaSecs,
		},
	}
}

func (this *Connection) CreateIndexWith(coll string, indexes ...Index) error {
	if coll == "" {
		return fmt.Errorf("empty collection name")
	}
	c := this.UseCollection(coll)
	for _, index := range indexes {
		err := c.CreateOneIndex(context.TODO(), index.AsQmgo())
		if err != nil {
			return err
		}
	}
	return nil
}

type IndexsGetter interface {
	GetCollection() string
	DBIndexes() []Index
}

func (this *Connection) CreateIndex(getter IndexsGetter) error {
	coll := getter.GetCollection()
	indexes := getter.DBIndexes()
	return this.CreateIndexWith(coll, indexes...)
}

// func (this *Connection) HasIndex(getter IndexsGetter) (bool, error) {
// 	coll := getter.GetCollection()
// 	indexes := getter.DBIndexes()
// }
