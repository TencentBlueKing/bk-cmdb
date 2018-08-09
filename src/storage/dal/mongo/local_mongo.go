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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/mongobyc"
	"configcenter/src/storage/mongobyc/findopt"
	"configcenter/src/storage/types"
	"context"
	"fmt"
	"strconv"
	"sync"
)

// Client implement client.DALRDB interface
type Client struct {
	dbc     mongobyc.CommonClient
	session mongobyc.Session
}

var _ dal.RDB = new(Client)
var _ dal.RDBTxn = new(ClientTxn)

var initMongoc sync.Once

// NewClient returns new RDB
func NewClient(uri string) (*Client, error) {
	initMongoc.Do(mongobyc.InitMongoc)

	client := mongobyc.NewClient(uri)
	return &Client{
		dbc: client,
	}, nil
}

// Close replica client
func (c *Client) Close() error {
	return c.dbc.Close()
}

// Ping replica client
func (c *Client) Ping() error {
	c.dbc.Database()
	return c.dbc.Ping()
}

func (c *Client) clone() *Client {
	nc := Client{
		dbc: c.dbc,
	}
	return &nc
}

// Collection collection operation
func (c *Client) Collection(collection string) dal.Collection {
	col := Collection{}
	col.collection = collection

	if c.session == nil {
		col.table = c.dbc.Collection
	} else {
		col.table = c.session.Collection
	}
	return &col
}

// Collection implement client.Collection interface
type Collection struct {
	RequestID  string // 请求ID,可选项
	Processor  string // 处理进程号，结构为"IP:PORT-PID"用于识别事务session被存于那个TM多活实例
	TxnID      string // 事务ID,uuid
	collection string // 集合名
	table      func(collName string) mongobyc.CollectionInterface
}

// Find 查询多个并反序列化到 Result
func (c *Collection) Find(ctx context.Context, filter types.Filter) dal.Find {
	return &Find{Collection: c, filter: filter, ctx: ctx}
}

// Find define a find operation
type Find struct {
	*Collection
	ctx        context.Context
	projection types.Document
	filter     types.Filter
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
func (f *Find) All(result interface{}) error {
	opt := findopt.Many{}
	opt.Skip = int64(f.start)
	opt.Limit = int64(f.limit)
	opt.Fields = mapstr.MapStr(f.projection)
	return f.table(f.collection).Find(f.ctx, f.filter, &opt, result)
}

// One 查询一个
func (f *Find) One(result interface{}) error {
	opt := findopt.One{}
	opt.Skip = int64(f.start)
	opt.Limit = int64(f.limit)
	opt.Fields = mapstr.MapStr(f.projection)
	return f.table(f.collection).FindOne(f.ctx, f.filter, &opt, result)
}

// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
func (c *Collection) Insert(ctx context.Context, docs interface{}) error {
	return c.table(c.collection).InsertMany(ctx, util.ConverToInterfaceSlice(docs), nil)
}

// Update 更新数据
func (c *Collection) Update(ctx context.Context, filter types.Filter, doc interface{}) error {
	_, err := c.table(c.collection).UpdateMany(ctx, filter, doc, nil)
	return err
}

// Delete 删除数据
func (c *Collection) Delete(ctx context.Context, filter types.Filter) error {
	c.table(c.collection).DeleteMany(ctx, filter, nil)
	return nil
}

// Count 统计数量(非事务)
func (c *Collection) Count(ctx context.Context, filter types.Filter) (uint64, error) {
	return c.table(c.collection).Count(ctx, filter)
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
	err := c.dbc.Collection(common.BKTableNameIDgenerator).FindAndModify(ctx, filter, data, &opt, &results)
	if nil != err {
		return 0, err
	}

	if len(results) <= 0 {
		return 0, dal.ErrDocumentNotFound
	}

	return strconv.ParseUint(fmt.Sprint(results[0]["SequenceID"]), 10, 64)
}

// StartTransaction 开启新事务
func (c *Client) StartTransaction(ctx context.Context, opt dal.JoinOption) (dal.RDBTxn, error) {
	session := c.dbc.Session().Create()
	if err := session.Open(); err != nil {
		return nil, err
	}
	txn := &ClientTxn{Client: c.clone()}
	txn.session = session
	return txn, session.StartTransaction()
}

// JoinTransaction 加入事务, controller 加入某个事务
func (c *Client) JoinTransaction(opt dal.JoinOption) dal.RDBTxn {
	blog.Fatalf("not support JoinTransaction")
	return nil
}

// ClientTxn implement dal.ClientTxn
type ClientTxn struct {
	*Client
}

// Commit 提交事务
func (c *ClientTxn) Commit() error {
	err := c.session.CommitTransaction()
	c.session.Close()
	c.session = nil
	return err
}

// Abort 取消事务
func (c *ClientTxn) Abort() error {
	err := c.session.AbortTransaction()
	c.session.Close()
	c.session = nil
	return err
}

// TxnInfo 当前事务信息，用于事务发起者往下传递
func (c *ClientTxn) TxnInfo() *types.Tansaction {
	return &types.Tansaction{}
}

// HasCollection 判断是否存在集合
func (c *Client) HasCollection(collName string) (bool, error) {
	return c.dbc.Database().HasCollection(collName)
}

// DropCollection 移除集合
func (c *Client) DropCollection(collName string) error {
	return c.dbc.Database().DropCollection(collName)
}

// CreateCollection 创建集合
func (c *Client) CreateCollection(collName string) error {
	return c.dbc.Database().CreateEmptyCollection(collName)
}

// CreateIndex 创建索引
func (c *Collection) CreateIndex(ctx context.Context, index dal.Index) error {
	i := mongobyc.Index{
		Keys:       mapstr.MapStr(index.Keys),
		Name:       index.Name,
		Unique:     index.Unique,
		Backgroupd: index.Backgroupd,
	}

	return c.table(c.collection).CreateIndex(i)
}

// DropIndex 移除索引
func (c *Collection) DropIndex(ctx context.Context, indexName string) error {
	return c.table(c.collection).DropIndex(indexName)
}

// AddColumn 添加字段
func (c *Collection) AddColumn(ctx context.Context, column string, value interface{}) error {
	selector := types.Document{column: types.Document{"$exists": false}}
	datac := types.Document{"$set": types.Document{column: value}}
	_, err := c.table(c.collection).UpdateMany(ctx, selector, datac, nil)
	return err
}

// RenameColumn 重命名字段
func (c *Collection) RenameColumn(ctx context.Context, oldName, newColumn string) error {
	datac := types.Document{"$rename": types.Document{oldName: newColumn}}
	_, err := c.table(c.collection).UpdateMany(ctx, nil, datac, nil)
	return err
}

// DropColumn 移除字段
func (c *Collection) DropColumn(ctx context.Context, field string) error {
	datac := types.Document{"$unset": types.Document{field: "1"}}
	_, err := c.table(c.collection).UpdateMany(ctx, nil, datac, nil)
	return err
}
