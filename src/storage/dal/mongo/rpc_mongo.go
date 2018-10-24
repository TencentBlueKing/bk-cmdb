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
	"errors"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/types"
)

// RPC implement client.DALRDB interface
type RPC struct {
	RequestID string // 请求ID,可选项
	TxnID     string // 事务ID,uuid
	rpc       *rpc.Client
	getServer types.GetServerFunc

	parent *RPC
}

var _ dal.RDB = new(RPC)

// NewRPCWithDiscover returns new RDB
func NewRPCWithDiscover(getServer types.GetServerFunc) (*RPC, error) {
	servers, err := getServer()
	if err != nil {
		return nil, err
	}

	rpccli, err := rpc.DialHTTPPath("tcp", servers[0], "/txn/v3/rpc")
	if err != nil {
		return nil, err
	}
	return &RPC{
		rpc:       rpccli,
		getServer: getServer,
	}, nil
}

// NewRPC returns new RDB
func NewRPC(uri string) (*RPC, error) {
	rpccli, err := rpc.DialHTTPPath("tcp", uri, "/txn/v3/rpc")
	if err != nil {
		return nil, err
	}
	return &RPC{
		rpc: rpccli,
	}, nil
}

// Close replica client
func (c *RPC) Close() error {
	return c.rpc.Close()
}

// Ping replica client
func (c *RPC) Ping() error {
	return c.rpc.Ping()
}

func (c *RPC) Clone() dal.RDB {
	nc := RPC{
		TxnID:     c.TxnID,
		RequestID: c.RequestID,
		rpc:       c.rpc,
		parent:    c,
	}
	return &nc
}

func (c *RPC) IsDuplicatedError(error) bool {
	return false
}
func (c *RPC) IsNotFoundError(error) bool {
	return false
}

// Table collection operation
func (c *RPC) Table(collection string) dal.Table {
	col := RPCCollection{
		RequestID: c.RequestID,
		TxnID:     c.TxnID,
	}
	col.collection = collection
	col.rpc = c.rpc

	return &col
}

// RPCCollection implement client.Collection interface
type RPCCollection struct {
	RequestID  string // 请求ID,可选项
	Processor  string // 处理进程号，结构为"IP:PORT-PID"用于识别事务session被存于那个TM多活实例
	TxnID      string // 事务ID,uuid
	collection string // 集合名
	rpc        *rpc.Client
}

// Find 查询多个并反序列化到 Result
func (c *RPCCollection) Find(filter dal.Filter) dal.Find {
	// build msg
	msg := types.OPFIND{}
	msg.OPCode = types.OPFind
	msg.Collection = c.collection
	msg.Selector.Encode(filter)

	find := RPCFind{RPCCollection: c, msg: &msg}
	find.RequestID = c.RequestID
	find.Processor = c.Processor
	find.TxnID = c.TxnID
	return &find
}

// RPCFind define a find operation
type RPCFind struct {
	*RPCCollection
	msg *types.OPFIND
}

// Fields 查询字段
func (f *RPCFind) Fields(fields ...string) dal.Find {
	projection := types.Document{}
	for _, field := range fields {
		projection[field] = true
	}
	f.msg.Projection = projection
	return f
}

// Sort 查询排序
func (f *RPCFind) Sort(sort string) dal.Find {
	f.msg.Sort = sort
	return f
}

// Start 查询上标
func (f *RPCFind) Start(start uint64) dal.Find {
	f.msg.Start = start
	return f
}

// Limit 查询限制
func (f *RPCFind) Limit(limit uint64) dal.Find {
	f.msg.Limit = limit
	return f
}

// All 查询多个
func (f *RPCFind) All(ctx context.Context, result interface{}) error {
	// set txn
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		f.msg.RequestID = opt.RequestID
		f.msg.TxnID = opt.TxnID
	}
	if f.TxnID != "" {
		f.msg.TxnID = f.TxnID
	}

	// call
	reply := types.OPREPLY{}
	err := f.rpc.Call(types.CommandRDBOperation, f.msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return reply.Docs.Decode(result)
}

// One 查询一个
func (f *RPCFind) One(ctx context.Context, result interface{}) error {
	// set txn
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		f.msg.RequestID = opt.RequestID
		f.msg.TxnID = opt.TxnID
	}
	if f.TxnID != "" {
		f.msg.TxnID = f.TxnID
	}

	// call
	reply := types.OPREPLY{}
	err := f.rpc.Call(types.CommandRDBOperation, f.msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}

	if len(reply.Docs) <= 0 {
		return dal.ErrDocumentNotFound
	}
	return reply.Docs[0].Decode(result)
}

// Count 统计数量(非事务)
func (f *RPCFind) Count(ctx context.Context) (uint64, error) {
	// build msg
	f.msg.OPCode = types.OPCount

	// set txn
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		f.msg.RequestID = opt.RequestID
	}

	// call
	reply := types.OPREPLY{}
	err := f.rpc.Call(types.CommandRDBOperation, f.msg, &reply)
	if err != nil {
		return 0, err
	}
	if !reply.Success {
		return 0, errors.New(reply.Message)
	}
	return reply.Count, nil
}

// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
func (c *RPCCollection) Insert(ctx context.Context, docs interface{}) error {
	// build msg
	msg := types.OPINSERT{}
	msg.OPCode = types.OPInsert
	msg.Collection = c.collection
	if err := msg.DOCS.Encode(docs); err != nil {
		return err
	}

	// set txn
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		msg.RequestID = opt.RequestID
		msg.TxnID = opt.TxnID
	}
	if c.TxnID != "" {
		msg.TxnID = c.TxnID
	}

	// call
	reply := types.OPREPLY{}
	err := c.rpc.Call(types.CommandRDBOperation, &msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return nil
}

// Update 更新数据
func (c *RPCCollection) Update(ctx context.Context, filter dal.Filter, doc interface{}) error {
	// build msg
	msg := types.OPUPDATE{}
	msg.OPCode = types.OPUpdate
	msg.Collection = c.collection
	if err := msg.DOC.Encode(types.Document{
		"$set": doc,
	}); err != nil {
		return err
	}
	if err := msg.Selector.Encode(filter); err != nil {
		return err
	}

	// set txn
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		msg.RequestID = opt.RequestID
		msg.TxnID = opt.TxnID
	}
	if c.TxnID != "" {
		msg.TxnID = c.TxnID
	}

	// call
	reply := types.OPREPLY{}
	err := c.rpc.Call(types.CommandRDBOperation, &msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return nil
}

// Delete 删除数据
func (c *RPCCollection) Delete(ctx context.Context, filter dal.Filter) error {
	// build msg
	msg := types.OPDELETE{}
	msg.OPCode = types.OPDelete
	msg.Collection = c.collection
	if err := msg.Selector.Encode(filter); err != nil {
		return err
	}
	if c.TxnID != "" {
		msg.TxnID = c.TxnID
	}

	// set txn
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		msg.RequestID = opt.RequestID
		msg.TxnID = opt.TxnID
	}
	if c.TxnID != "" {
		msg.TxnID = c.TxnID
	}

	// call
	reply := types.OPREPLY{}
	err := c.rpc.Call(types.CommandRDBOperation, &msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return nil
}

// NextSequence 获取新序列号(非事务)
func (c *RPC) NextSequence(ctx context.Context, sequenceName string) (uint64, error) {
	// build msg
	msg := types.OPFINDANDMODIFY{}
	msg.OPCode = types.OPFindAndModify
	msg.Collection = common.BKTableNameIDgenerator
	if err := msg.DOC.Encode(types.Document{
		"$inc": types.Document{"SequenceID": 1},
	}); err != nil {
		return 0, err
	}
	if err := msg.Selector.Encode(types.Document{
		"_id": sequenceName,
	}); err != nil {
		return 0, err
	}
	msg.Upsert = true
	msg.ReturnNew = true

	// set txn
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		msg.RequestID = opt.RequestID
		// msg.TxnID = opt.TxnID // because NextSequence was not supported for transaction in mongo
	}

	// call
	reply := types.OPREPLY{}
	err := c.rpc.Call(types.CommandRDBOperation, &msg, &reply)
	if err != nil {
		return 0, err
	}
	if !reply.Success {
		return 0, errors.New(reply.Message)
	}

	if len(reply.Docs) <= 0 {
		return 0, dal.ErrDocumentNotFound
	}

	return strconv.ParseUint(fmt.Sprint(reply.Docs[0]["SequenceID"]), 10, 64)
}

// StartTransaction 开启新事务
func (c *RPC) StartTransaction(ctx context.Context) (dal.RDB, error) {
	if c.TxnID != "" {
		return nil, dal.ErrTransactionStated
	}
	// build msg
	msg := types.OPSTARTTTRANSATION{}
	msg.OPCode = types.OPStartTransaction

	// set txn
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		msg.RequestID = opt.RequestID
	}

	// call
	reply := types.OPREPLY{}
	err := c.rpc.Call(types.CommandRDBOperation, &msg, &reply)
	if err != nil {
		return nil, err
	}
	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	clone := c.Clone().(*RPC)
	clone.TxnID = reply.TxnID
	clone.RequestID = reply.RequestID
	return clone, nil
}

// Commit 提交事务
func (c *RPC) Commit(ctx context.Context) error {
	if c.TxnID == "" {
		return dal.ErrTransactionNotFound
	}
	msg := types.OPCOMMIT{}
	msg.OPCode = types.OPCommit
	msg.RequestID = c.RequestID
	msg.TxnID = c.TxnID

	reply := types.OPREPLY{}
	err := c.rpc.Call(types.CommandRDBOperation, &msg, &reply)
	c.TxnID = "" // clear TxnID
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return nil
}

// Abort 取消事务
func (c *RPC) Abort(ctx context.Context) error {
	if c.TxnID == "" {
		return dal.ErrTransactionNotFound
	}
	msg := types.OPABORT{}
	msg.OPCode = types.OPAbort
	msg.RequestID = c.RequestID
	msg.TxnID = c.TxnID

	reply := types.OPREPLY{}
	err := c.rpc.Call(types.CommandRDBOperation, &msg, &reply)
	c.TxnID = "" // clear TxnID
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return nil
}

// TxnInfo 当前事务信息，用于事务发起者往下传递
func (c *RPC) TxnInfo() *types.Transaction {
	return &types.Transaction{
		RequestID: c.RequestID,
		TxnID:     c.TxnID,
	}
}

// HasTable 判断是否存在集合
func (c *RPC) HasTable(tablename string) (bool, error) {
	return false, dal.ErrNotImplemented
}

// DropTable 移除集合
func (c *RPC) DropTable(tablename string) error {
	return dal.ErrNotImplemented
}

// CreateTable 创建集合
func (c *RPC) CreateTable(tablename string) error {
	return dal.ErrNotImplemented
}

// CreateIndex 创建索引
func (c *RPCCollection) CreateIndex(ctx context.Context, index dal.Index) error {
	return dal.ErrNotImplemented
}

// DropIndex 移除索引
func (c *RPCCollection) DropIndex(ctx context.Context, indexName string) error {
	return dal.ErrNotImplemented
}

// Indexes 查询索引
func (c *RPCCollection) Indexes(ctx context.Context) ([]dal.Index, error) {
	return nil, dal.ErrNotImplemented
}

// AddColumn 添加字段
func (c *RPCCollection) AddColumn(ctx context.Context, column string, value interface{}) error {
	return dal.ErrNotImplemented
}

// RenameColumn 重命名字段
func (c *RPCCollection) RenameColumn(ctx context.Context, oldName, newColumn string) error {
	return dal.ErrNotImplemented
}

// DropColumn 移除字段
func (c *RPCCollection) DropColumn(ctx context.Context, field string) error {
	return dal.ErrNotImplemented
}
