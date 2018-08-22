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

package mongoc

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/mongobyc"
	"configcenter/src/storage/mongobyc/findopt"
	"configcenter/src/storage/types"
)

var ErrSessionMissing = errors.New("session missing")

// Client implement client.DALRDB interface
type Client struct {
	txc     mongobyc.Client
	session mongobyc.Session
	pool    mongobyc.ClientPool
}

var _ dal.RDB = new(Client)

var initMongoc sync.Once

// NewClient returns new RDB
func NewClient(uri string) (*Client, error) {
	initMongoc.Do(mongobyc.InitMongoc)

	pool := mongobyc.NewClientPool(uri)
	err := pool.Open()
	if err != nil {
		return nil, err
	}
	return &Client{
		pool: pool,
	}, nil
}

// Close replica client
func (c *Client) Close() error {
	return c.pool.Close()
}

// Ping replica client
func (c *Client) Ping() error {
	dbc := c.pool.Pop()
	err := dbc.Ping()
	c.pool.Push(dbc)
	return err
}

// Clone return the new client
func (c *Client) Clone() dal.RDB {
	nc := Client{
		pool: c.pool,
	}
	return &nc
}

// Table collection operation
func (c *Client) Table(collName string) dal.Table {
	col := Collection{}
	col.collName = collName
	col.Client = c
	return &col
}

// Collection implement client.Collection interface
type Collection struct {
	collName string // 集合名
	*Client
}

// Find 查询多个并反序列化到 Result
func (c *Collection) Find(filter dal.Filter) dal.Find {
	return &Find{Collection: c, filter: filter}
}

// Find define a find operation
type Find struct {
	*Collection
	projection types.Document
	filter     dal.Filter
	start      uint64
	limit      uint64
	sort       string
}

// Fields 查询字段
func (f *Find) Fields(fields ...string) dal.Find {
	projection := types.Document{}
	for _, field := range fields {
		projection[field] = true
	}
	f.projection = projection
	return f
}

// Sort 查询排序
func (f *Find) Sort(sort string) dal.Find {
	f.sort = sort
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
	opt := findopt.Many{}
	opt.Skip = int64(f.start)
	opt.Limit = int64(f.limit)
	opt.Fields = mapstr.MapStr(f.projection)

	p, table := f.getCollection(f.collName)
	err := table.Find(ctx, f.filter, &opt, result)
	p.push()
	return err
}

// One 查询一个
func (f *Find) One(ctx context.Context, result interface{}) error {
	opt := findopt.One{}
	opt.Skip = int64(f.start)
	opt.Limit = int64(f.limit)
	opt.Fields = mapstr.MapStr(f.projection)

	p, table := f.getCollection(f.collName)
	err := table.FindOne(ctx, f.filter, &opt, result)
	p.push()
	return err
}

// Count 统计数量(非事务)
func (f *Find) Count(ctx context.Context) (uint64, error) {
	dbc := f.pool.Pop()
	count, err := dbc.Collection(f.collName).Count(ctx, f.filter)
	f.pool.Push(dbc)
	return count, err
}

// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
func (c *Collection) Insert(ctx context.Context, docs interface{}) error {
	p, table := c.getCollection(c.collName)
	err := table.InsertMany(ctx, util.ConverToInterfaceSlice(docs), nil)
	p.push()
	return err
}

// Update 更新数据
func (c *Collection) Update(ctx context.Context, filter dal.Filter, doc interface{}) error {
	p, table := c.getCollection(c.collName)
	_, err := table.UpdateMany(ctx, filter, doc, nil)
	p.push()
	return err
}

// Delete 删除数据
func (c *Collection) Delete(ctx context.Context, filter dal.Filter) error {
	p, table := c.getCollection(c.collName)
	_, err := table.DeleteMany(ctx, filter, nil)
	p.push()
	return err
}

// NextSequence 获取新序列号(非事务)
func (c *Client) NextSequence(ctx context.Context, sequenceName string) (uint64, error) {
	data := types.Document{
		"$inc": types.Document{"SequenceID": 1},
	}
	filter := types.Document{
		"_id": sequenceName,
	}

	opt := findopt.FindAndModify{}
	opt.Upsert = true
	opt.New = true

	results := types.Documents{}
	dbc := c.pool.Pop()
	err := dbc.Collection(common.BKTableNameIDgenerator).FindAndModify(ctx, filter, data, &opt, &results)
	c.pool.Push(dbc)
	if nil != err {
		return 0, err
	}

	if len(results) <= 0 {
		return 0, dal.ErrDocumentNotFound
	}

	return strconv.ParseUint(fmt.Sprint(results[0]["SequenceID"]), 10, 64)
}

// StartTransaction 开启新事务
func (c *Client) StartTransaction(ctx context.Context) error {
	txc := c.pool.Pop()
	c.txc = txc
	session := txc.Session().Create()
	if err := session.Open(); err != nil {
		session.Close()
		return err
	}
	c.session = session
	err := session.StartTransaction()
	if err != nil {
		session.Close()
	}
	return err
}

// Commit 提交事务
func (c *Client) Commit() error {
	if c.session == nil {
		return ErrSessionMissing
	}
	commitErr := c.session.CommitTransaction()
	if commitErr == nil {
		closeErr := c.session.Close()
		c.pool.Push(c.txc)
		if closeErr != nil {
			blog.Warnf("[mongoc dal] session close faile: %v", closeErr)
		}
		c.session = nil
		return nil
	}
	return commitErr
}

// Abort 取消事务
func (c *Client) Abort() error {
	if c.session == nil {
		return ErrSessionMissing
	}
	abortErr := c.session.AbortTransaction()
	closeErr := c.session.Close()
	c.pool.Push(c.txc)
	if closeErr != nil {
		blog.Warnf("[mongoc dal] session close faile: %v", closeErr)
	}
	c.session = nil
	return abortErr
}

// TxnInfo 当前事务信息，用于事务发起者往下传递
func (c *Client) TxnInfo() *types.Transaction {
	return &types.Transaction{}
}

// HasTable 判断是否存在集合
func (c *Client) HasTable(collName string) (bool, error) {
	dbc := c.pool.Pop()
	exists, err := dbc.Database().HasCollection(collName)
	c.pool.Push(dbc)
	return exists, err
}

// DropTable 移除集合
func (c *Client) DropTable(collName string) error {
	dbc := c.pool.Pop()
	err := dbc.Database().DropCollection(collName)
	c.pool.Push(dbc)
	return err
}

// CreateTable 创建集合
func (c *Client) CreateTable(collName string) error {
	dbc := c.pool.Pop()
	err := dbc.Database().CreateEmptyCollection(collName)
	c.pool.Push(dbc)
	return err
}

// CreateIndex 创建索引
func (c *Collection) CreateIndex(ctx context.Context, index dal.Index) error {
	i := mongobyc.Index{
		Keys:       mapstr.MapStr(index.Keys),
		Name:       index.Name,
		Unique:     index.Unique,
		Backgroupd: index.Backgroupd,
	}

	dbc := c.pool.Pop()
	err := dbc.Collection(c.collName).CreateIndex(i)
	c.pool.Push(dbc)
	return err
}

// DropIndex 移除索引
func (c *Collection) DropIndex(ctx context.Context, indexName string) error {
	dbc := c.pool.Pop()
	err := dbc.Collection(c.collName).DropIndex(indexName)
	c.pool.Push(dbc)
	return err
}

// AddColumn 添加字段
func (c *Collection) AddColumn(ctx context.Context, column string, value interface{}) error {
	selector := types.Document{column: types.Document{"$exists": false}}
	datac := types.Document{"$set": types.Document{column: value}}

	dbc := c.pool.Pop()
	_, err := dbc.Collection(c.collName).UpdateMany(ctx, selector, datac, nil)
	c.pool.Push(dbc)
	return err
}

// RenameColumn 重命名字段
func (c *Collection) RenameColumn(ctx context.Context, oldName, newColumn string) error {
	datac := types.Document{"$rename": types.Document{oldName: newColumn}}

	dbc := c.pool.Pop()
	_, err := dbc.Collection(c.collName).UpdateMany(ctx, nil, datac, nil)
	c.pool.Push(dbc)
	return err
}

// DropColumn 移除字段
func (c *Collection) DropColumn(ctx context.Context, field string) error {
	datac := types.Document{"$unset": types.Document{field: "1"}}
	dbc := c.pool.Pop()
	_, err := dbc.Collection(c.collName).UpdateMany(ctx, nil, datac, nil)
	c.pool.Push(dbc)
	return err
}

type pusher struct {
	pool mongobyc.ClientPool
	dbc  mongobyc.Client
}

func (c *Client) getCollection(collName string) (*pusher, mongobyc.CollectionInterface) {
	var table mongobyc.CollectionInterface
	var p = new(pusher)
	if c.session == nil {
		p.dbc = c.pool.Pop()
		p.pool = c.pool
		table = p.dbc.Collection(collName)
	} else {
		table = c.session.Collection(collName)
	}
	return p, table
}

func (p *pusher) push() {
	if p.dbc != nil {
		p.pool.Push(p.dbc)
	}
}
