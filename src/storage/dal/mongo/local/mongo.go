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

package local

import (
	"context"
	"errors"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
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

var _ dal.DB = new(Mongo)

// NewMgo returns new RDB
func NewMgo(uri string, timeout time.Duration) (*Mongo, error) {
	cs, err := mgo.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	client, err := mgo.DialWithTimeout(uri, time.Second*10)
	if err != nil {
		return nil, err
	}
	client.SetSyncTimeout(timeout)
	client.SetSocketTimeout(timeout)
	client.SetPoolLimit(1000)
	if err := client.Ping(); err != nil {
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
	return c.dbc.Ping()
}

// Clone return the new client
func (c *Mongo) Clone() dal.DB {
	nc := Mongo{
		dbc:    c.dbc,
		dbname: c.dbname,
	}
	return &nc
}

// IsDuplicatedError check duplicated error
func (c *Mongo) IsDuplicatedError(err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "The existing index") {
			return true
		}
		if strings.Contains(err.Error(), "There's already an index with name") {
			return true
		}
	}
	return err == dal.ErrDuplicated || mgo.IsDup(err)
}

// IsNotFoundError check the not found error
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
		if len(field) <= 0 {
			continue
		}
		f.projection[field] = true
	}
	return f
}

// Sort 查询排序
func (f *Find) Sort(sort string) dal.Find {
	if sort != "" {
		f.sort = strings.Split(sort, ",")
	}
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
	sess := f.dbc.Clone()

	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	query := sess.DB(f.dbname).C(f.collName).Find(f.filter)
	query = query.Select(f.projection)
	query = query.Skip(int(f.start))
	query = query.Limit(int(f.limit))
	query = query.Sort(f.sort...)
	err := query.All(result)
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo find-all cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	return err
}

// One 查询一个
func (f *Find) One(ctx context.Context, result interface{}) error {
	sess := f.dbc.Clone()
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	err := sess.DB(f.dbname).C(f.collName).Find(f.filter).One(result)
	if err == mgo.ErrNotFound {
		err = dal.ErrDocumentNotFound
	}
	sess.Close()

	blog.V(4).InfoDepthf(1, "mongo find-one cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	return err
}

// Count 统计数量(非事务)
func (f *Find) Count(ctx context.Context) (uint64, error) {
	sess := f.dbc.Clone()
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	count, err := sess.DB(f.dbname).C(f.collName).Find(f.filter).Count()
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo count cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)

	return uint64(count), err
}

// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
func (c *Collection) Insert(ctx context.Context, docs interface{}) error {
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	sess := c.dbc.Clone()
	err := sess.DB(c.dbname).C(c.collName).Insert(util.ConverToInterfaceSlice(docs)...)
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo insert cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	return err
}

// Update 更新数据
func (c *Collection) Update(ctx context.Context, filter dal.Filter, doc interface{}) error {
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	sess := c.dbc.Clone()
	data := bson.M{"$set": doc}
	_, err := sess.DB(c.dbname).C(c.collName).UpdateAll(filter, data)
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo update cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)

	return err
}

// upsert 更新数据
func (c *Collection) Upsert(ctx context.Context, filter dal.Filter, doc interface{}) error {
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	sess := c.dbc.Clone()
	data := bson.M{"$set": doc}
	_, err := sess.DB(c.dbname).C(c.collName).Upsert(filter, data)
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo upsert cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	return err
}

// UpdateMultiModel 根据不同的操作符去更新数据
func (c *Collection) UpdateMultiModel(ctx context.Context, filter dal.Filter, updateModel ...dal.ModeUpdate) error {

	data := bson.M{}
	for _, item := range updateModel {
		if _, ok := data[item.Op]; ok {
			return errors.New(item.Op + " appear multiple times")
		}
		data["$"+item.Op] = item.Doc
	}

	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	sess := c.dbc.Clone()
	_, err := sess.DB(c.dbname).C(c.collName).UpdateAll(filter, data)
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo update-multi-model cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)

	return err
}

// Delete 删除数据
func (c *Collection) Delete(ctx context.Context, filter dal.Filter) error {
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	sess := c.dbc.Clone()
	_, err := sess.DB(c.dbname).C(c.collName).RemoveAll(filter)
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo delete cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	return err
}

// NextSequence 获取新序列号(非事务)
func (c *Mongo) NextSequence(ctx context.Context, sequenceName string) (uint64, error) {
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	sess := c.dbc.Clone()
	coll := sess.DB(c.dbname).C("cc_idgenerator")
	change := mgo.Change{
		Update: bson.M{
			"$inc":         bson.M{"SequenceID": int64(1)},
			"$setOnInsert": bson.M{"create_time": time.Now()},
			"$set":         bson.M{"last_time": time.Now()},
		},
		ReturnNew: true,
		Upsert:    true,
	}
	doc := Idgen{}

	_, err := coll.Find(bson.M{"_id": sequenceName}).Apply(change, &doc)
	if err != nil {
		sess.Close()
		return 0, err
	}
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo next-sequence cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	return doc.SequenceID, err
}

type Idgen struct {
	ID         string `bson:"_id"`
	SequenceID uint64 `bson:"SequenceID"`
}

// Start 开启新事务
func (c *Mongo) Start(ctx context.Context) (dal.Transcation, error) {
	return c, nil
}

// Commit 提交事务
func (c *Mongo) Commit(ctx context.Context) error {
	return nil
}

// Abort 取消事务
func (c *Mongo) Abort(ctx context.Context) error {
	return nil
}

// TxnInfo 当前事务信息，用于事务发起者往下传递
func (c *Mongo) TxnInfo() *types.Transaction {
	return &types.Transaction{}
}

// HasTable 判断是否存在集合
func (c *Mongo) HasTable(collName string) (bool, error) {
	sess := c.dbc.Clone()
	defer sess.Close()
	colls, err := sess.DB(c.dbname).CollectionNames()
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
	sess := c.dbc.Clone()
	defer sess.Close()
	return sess.DB(c.dbname).C(collName).DropCollection()
}

// CreateTable 创建集合
func (c *Mongo) CreateTable(collName string) error {
	sess := c.dbc.Clone()
	defer sess.Close()
	return sess.DB(c.dbname).C(collName).Create(&mgo.CollectionInfo{})
}

// DB get dal interface
func (c *Mongo) DB(collName string) dal.RDB {
	return c
}

// CreateIndex 创建索引
func (c *Collection) CreateIndex(ctx context.Context, index dal.Index) error {
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
	sess := c.dbc.Clone()
	defer sess.Close()
	return sess.DB(c.dbname).C(c.collName).EnsureIndex(i)
}

// DropIndex remove index by name
func (c *Collection) DropIndex(ctx context.Context, indexName string) error {
	sess := c.dbc.Clone()
	defer sess.Close()
	return sess.DB(c.dbname).C(c.collName).DropIndexName(indexName)
}

// Indexes get all indexes for the collection
func (c *Collection) Indexes(ctx context.Context) ([]dal.Index, error) {
	sess := c.dbc.Clone()
	defer sess.Close()
	dbindexs, err := sess.DB(c.dbname).C(c.collName).Indexes()
	if err != nil {
		return nil, err
	}

	indexs := []dal.Index{}
	for _, dbindex := range dbindexs {
		keys := map[string]int32{}
		for _, key := range dbindex.Key {
			if strings.HasPrefix(key, "-") {
				key = strings.TrimLeft(key, "-")
				keys[key] = -1
			} else {
				keys[key] = 1
			}
		}

		index := dal.Index{}
		index.Name = dbindex.Name
		index.Unique = dbindex.Unique
		index.Background = dbindex.Background
		index.Keys = keys
		indexs = append(indexs, index)
	}
	return indexs, nil
}

// AddColumn add a new column for the collection
func (c *Collection) AddColumn(ctx context.Context, column string, value interface{}) error {
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	sess := c.dbc.Clone()

	selector := types.Document{column: types.Document{"$exists": false}}
	datac := types.Document{"$set": types.Document{column: value}}
	_, err := sess.DB(c.dbname).C(c.collName).UpdateAll(selector, datac)
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo add-column cost: %sms, rid: %s", time.Since(start)/time.Millisecond, rid)
	return err
}

// RenameColumn rename a column for the collection
func (c *Collection) RenameColumn(ctx context.Context, oldName, newColumn string) error {
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	sess := c.dbc.Clone()
	datac := types.Document{"$rename": types.Document{oldName: newColumn}}
	_, err := sess.DB(c.dbname).C(c.collName).UpdateAll(types.Document{}, datac)
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo rename-column cost: %sms, rid: %s", time.Since(start)/time.Millisecond, rid)
	return err
}

// DropColumn remove a column by the name
func (c *Collection) DropColumn(ctx context.Context, field string) error {
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	sess := c.dbc.Clone()
	datac := types.Document{"$unset": types.Document{field: ""}}
	_, err := sess.DB(c.dbname).C(c.collName).UpdateAll(types.Document{}, datac)
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo drop-column cost: %sms, rid: %s", time.Since(start)/time.Millisecond, rid)
	return err
}

// DropColumn remove a column by the name
func (c *Collection) DropColumns(ctx context.Context, filter dal.Filter, fields []string) error {
	c.dbc.Refresh()
	unsetFields := make(map[string]interface{})
	for _, field := range fields {
		unsetFields[field] = ""
	}
	datac := types.Document{"$unset": unsetFields}
	_, err := c.dbc.DB(c.dbname).C(c.collName).UpdateAll(filter, datac)
	return err
}

// AggregateAll aggregate all operation
func (c *Collection) AggregateAll(ctx context.Context, pipeline interface{}, result interface{}) error {
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	sess := c.dbc.Clone()
	err := sess.DB(c.dbname).C(c.collName).Pipe(pipeline).All(result)
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo aggregate-all cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	return err
}

// AggregateOne aggregate one operation
func (c *Collection) AggregateOne(ctx context.Context, pipeline interface{}, result interface{}) error {
	rid := ctx.Value(common.ContextRequestIDField)
	start := time.Now()
	sess := c.dbc.Clone()
	err := sess.DB(c.dbname).C(c.collName).Pipe(pipeline).One(result)
	sess.Close()
	blog.V(4).InfoDepthf(1, "mongo aggregate-one cost %dms, rid: %v", time.Since(start)/time.Millisecond, rid)
	return err
}
