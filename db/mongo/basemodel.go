package mongo

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/qiniu/qmgo"
)

var (
	ErrEmptyObjectID = errors.New("empty objectId form m.GetID")

	// emptyStr = reflect.Zero(reflect.TypeOf("")).Interface()
)

type BaseModel interface {
	// Get ID of document
	GetID() interface{}

	// Get Name of Collection
	GetCollection() string
}

func (this *Connection) CreateM(m BaseModel) error {
	c := this.UseCollection(m.GetCollection())
	objId := m.GetID()
	if objId == nil {
		return ErrEmptyObjectID
	}
	selector := qmgo.M{"_id": objId}
	count, err := c.Find(context.TODO(), selector).Limit(1).Count()
	if err != nil {
		if err != qmgo.ErrNoSuchDocuments {
			return err
		}
	}
	if count >= 1 {
		return nil
	}
	rs, err := c.UpsertId(context.TODO(), objId, m)
	_ = rs
	return err
}

func (this *Connection) GetM(m BaseModel) error {
	c := this.UseCollection(m.GetCollection())
	objId := m.GetID()
	if objId == "" {
		return ErrEmptyObjectID
	}
	selector := qmgo.M{"_id": objId}
	return c.Find(context.TODO(), selector).One(m)
}

func (this *Connection) GetMby(m BaseModel, selector interface{}) error {
	c := this.UseCollection(m.GetCollection())
	return c.Find(context.TODO(), selector).One(m)
}

func (this *Connection) GetList(coll string, selector interface{}, skip, limit int, result interface{}) error {
	c := this.UseCollection(coll)
	query := c.Find(context.TODO(), selector)
	if skip > 0 {
		query = query.Skip(int64(skip))
	}
	if limit > 0 {
		query = query.Limit(int64(limit))
	}
	err := query.All(result)
	if err == qmgo.ErrNoSuchDocuments {
		return nil
	}
	return err
}

func (this *Connection) GetOrCreateM(m BaseModel) (bool, error) {
	if err := this.GetM(m); err == nil {
		return false, nil
	} else if err == qmgo.ErrNoSuchDocuments {
		return true, this.CreateM(m)
	} else {
		return false, err
	}
}

func (this *Connection) UpdateM(m BaseModel) error {
	c := this.UseCollection(m.GetCollection())
	objId := m.GetID()
	if objId == "" {
		return ErrEmptyObjectID
	}
	selector := qmgo.M{"_id": objId}
	updates := qmgo.M{
		"$set": m,
	}
	return c.UpdateOne(context.TODO(), selector, updates)
}

func (this *Connection) UpdateOrCreateM(m BaseModel) (bool, error) {
	c := this.UseCollection(m.GetCollection())
	objId := m.GetID()
	if objId == "" {
		return false, ErrEmptyObjectID
	}
	res, err := c.UpsertId(context.TODO(), objId, m)
	return res.UpsertedCount > 0, err
}

func (this *Connection) DeleteM(m BaseModel) error {
	c := this.UseCollection(m.GetCollection())
	objId := m.GetID()
	if objId == "" {
		return ErrEmptyObjectID
	}
	return c.RemoveId(context.TODO(), objId)
}

// TODO: 实现 emptyomit
func (this *Connection) UpdateMField(m BaseModel, fieldName string) error {
	c := this.UseCollection(m.GetCollection())
	// FIXME: GetField 的逻辑转移到一个独立模块
	// 获取字段
	field := reflect.ValueOf(m)
	val := reflect.Indirect(field).FieldByName(fieldName)
	if !val.IsValid() {
		panic(fmt.Errorf("Unexist field(%s)", fieldName))
	}
	// 获取字段 bson 键名
	fieldType := reflect.TypeOf(m).Elem()
	keyType, _ := fieldType.FieldByName(fieldName)
	bsonKey := keyType.Tag.Get("bson")

	objId := m.GetID()
	err := c.UpdateId(context.TODO(), objId, qmgo.M{
		"$set": qmgo.M{
			bsonKey: val.Interface(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}
