/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mongo

import (
	"context"
	"strings"

	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo implement client.DALRDB interface
type Mongo struct {
	dbc    *mgo.Session
	dbname string
}

var _ dal.RDB = new(Mongo)

// NewMgo returns new RDB
func NewMgo(uri string) (*Mongo, error) {
	cs, err := mgo.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	client, err := mgo.DialWithInfo(cs)
	if err != nil {
		return nil, err
	}
	return &Mongo{
		dbc:    client,
		dbname: cs.Database,
	}, nil
}

// Close replica client
func (c *Mongo) Close() error {
	c.dbc.Close()
	return nil
}

// Ping replica client
func (c *Mongo) Ping() error {
	// TODO
	return nil
}

// Clone return the new client
func (c *Mongo) Clone() dal.RDB {
	nc := Mongo{
		dbc:    c.dbc,
		dbname: c.dbname,
	}
	return &nc
}

func (c *Mongo) IsDuplicatedError(err error) bool {
	return err == dal.ErrDuplicated || mgo.IsDup(err)
}
func (c *Mongo) IsNotFoundError(err error) bool {
	return err == dal.ErrDocumentNotFound || err == mgo.ErrNotFound
}

// Table collection operation
func (c *Mongo) Table(collName string) dal.Table {
	col := Collection{}
	col.collName = collName
	col.Mongo = c
	return &col
}

// Collection implement client.Collection interface
type Collection struct {
	collName string // 集合名
	*Mongo
}

// Find 查询多个并反序列化到 Result
func (c *Collection) Find(filter dal.Filter) dal.Find {
	return &Find{Collection: c, filter: filter, projection: types.Document{"_id": false}}
}

// Find define a find operation
type Find struct {
	*Collection
	projection types.Document
	filter     dal.Filter
	start      uint64
	limit      uint64
	sort       []string
}

// Fields 查询字段
func (f *Find) Fields(fields ...string) dal.Find {
	for _, field := range fields {
		f.projection[field] = true
	}
	return f
}

// Sort 查询排序
func (f *Find) Sort(sort string) dal.Find {
	f.sort = strings.Split(sort, ",")
	return f
}

// Start 查询上标
func (f *Find) Start(start uint64) dal.Find {
	f.start = start
	return f
}

// Limit 查询限制
func (f *Find) Limit(limit uint64) dal.Find {
	f.limit = limit
	return f
}

// All 查询多个
func (f *Find) All(ctx context.Context, result interface{}) error {
	f.dbc.Refresh()
	query := f.dbc.DB(f.dbname).C(f.collName).Find(f.filter)
	query = query.Skip(int(f.start))
	query = query.Limit(int(f.limit))
	query = query.Sort(f.sort...)
	return query.All(result)
}

// One 查询一个
func (f *Find) One(ctx context.Context, result interface{}) error {
	f.dbc.Refresh()

	err := f.dbc.DB(f.dbname).C(f.collName).Find(f.filter).One(result)
	if err == mgo.ErrNotFound {
		err = dal.ErrDocumentNotFound
	}
	return err
}

// Count 统计数量(非事务)
func (f *Find) Count(ctx context.Context) (uint64, error) {
	count, err := f.dbc.DB(f.dbname).C(f.collName).Find(f.filter).Count()
	return uint64(count), err
}

// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
func (c *Collection) Insert(ctx context.Context, docs interface{}) error {
	c.dbc.Refresh()
	return c.dbc.DB(c.dbname).C(c.collName).Insert(util.ConverToInterfaceSlice(docs)...)
}

// Update 更新数据
func (c *Collection) Update(ctx context.Context, filter dal.Filter, doc interface{}) error {
	c.dbc.Refresh()
	data := bson.M{"$set": doc}
	_, err := c.dbc.DB(c.dbname).C(c.collName).UpdateAll(filter, data)
	return err
}

// Delete 删除数据
func (c *Collection) Delete(ctx context.Context, filter dal.Filter) error {
	c.dbc.Refresh()
	_, err := c.dbc.DB(c.dbname).C(c.collName).RemoveAll(filter)
	return err
}

// NextSequence 获取新序列号(非事务)
func (c *Mongo) NextSequence(ctx context.Context, sequenceName string) (uint64, error) {
	c.dbc.Refresh()
	coll := c.dbc.DB(c.dbname).C("cc_idgenerator")
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"SequenceID": int64(1)}},
		ReturnNew: true,
		Upsert:    true,
	}
	doc := Idgen{}

	_, err := coll.Find(bson.M{"_id": sequenceName}).Apply(change, &doc)
	if err != nil {
		return 0, err
	}
	return doc.SequenceID, err
}

type Idgen struct {
	ID         string `bson:"_id"`
	SequenceID uint64 `bson:"SequenceID"`
}

// StartTransaction 开启新事务
func (c *Mongo) StartTransaction(ctx context.Context) error {
	return nil
}

// Commit 提交事务
func (c *Mongo) Commit() error {
	return nil
}

// Abort 取消事务
func (c *Mongo) Abort() error {
	return nil
}

// TxnInfo 当前事务信息，用于事务发起者往下传递
func (c *Mongo) TxnInfo() *types.Tansaction {
	return &types.Tansaction{}
}

// HasTable 判断是否存在集合
func (c *Mongo) HasTable(collName string) (bool, error) {
	c.dbc.Refresh()
	colls, err := c.dbc.DB(c.dbname).CollectionNames()
	if err != nil {
		return false, err
	}

	for _, coll := range colls {
		if coll == collName {
			return true, nil
		}
	}
	return false, err
}

// DropTable 移除集合
func (c *Mongo) DropTable(collName string) error {
	c.dbc.Refresh()
	return c.dbc.DB(c.dbname).C(collName).DropCollection()
}

// CreateTable 创建集合
func (c *Mongo) CreateTable(collName string) error {
	c.dbc.Refresh()
	return c.dbc.DB(c.dbname).C(collName).Create(&mgo.CollectionInfo{})
}

// CreateIndex 创建索引
func (c *Collection) CreateIndex(ctx context.Context, index dal.Index) error {
	c.dbc.Refresh()
	keys := []string{}
	for key := range index.Keys {
		keys = append(keys, key)
	}

	i := mgo.Index{
		Key:        keys,
		Name:       index.Name,
		Unique:     index.Unique,
		Background: index.Background,
	}
	return c.dbc.DB(c.dbname).C(c.collName).EnsureIndex(i)
}

// DropIndex 移除索引
func (c *Collection) DropIndex(ctx context.Context, indexName string) error {
	c.dbc.Refresh()
	return c.dbc.DB(c.dbname).C(c.collName).DropIndexName(indexName)
}

// AddColumn 添加字段
func (c *Collection) AddColumn(ctx context.Context, column string, value interface{}) error {
	c.dbc.Refresh()
	selector := types.Document{column: types.Document{"$exists": false}}
	datac := types.Document{"$set": types.Document{column: value}}

	_, err := c.dbc.DB(c.dbname).C(c.collName).UpdateAll(selector, datac)
	return err
}

// RenameColumn 重命名字段
func (c *Collection) RenameColumn(ctx context.Context, oldName, newColumn string) error {
	c.dbc.Refresh()
	datac := types.Document{"$rename": types.Document{oldName: newColumn}}
	_, err := c.dbc.DB(c.dbname).C(c.collName).UpdateAll(nil, datac)
	return err
}

// DropColumn 移除字段
func (c *Collection) DropColumn(ctx context.Context, field string) error {
	c.dbc.Refresh()
	datac := types.Document{"$unset": types.Document{field: "1"}}
	_, err := c.dbc.DB(c.dbname).C(c.collName).UpdateAll(nil, datac)
	return err
}
