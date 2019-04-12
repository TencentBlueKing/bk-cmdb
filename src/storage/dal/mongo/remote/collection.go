/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package remote

import (
	"context"
	"errors"

	"configcenter/src/common"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/types"
)

// Collection implement client.Collection interface
type Collection struct {
	RequestID  string // 请求ID,可选项
	Processor  string // 处理进程号，结构为"IP:PORT-PID"用于识别事务session被存于那个TM多活实例
	TxnID      string // 事务ID,uuid
	collection string // 集合名
	rpc        rpc.Client
}

// Find 查询多个并反序列化到 Result
func (c *Collection) Find(filter dal.Filter) dal.Find {
	// build msg
	msg := types.OPFindOperation{}
	msg.OPCode = types.OPFindCode
	msg.Collection = c.collection
	msg.Selector.Encode(filter)

	find := Find{Collection: c, msg: &msg}
	find.RequestID = c.RequestID
	find.Processor = c.Processor
	find.TxnID = c.TxnID
	return &find
}

// Update 更新数据
func (c *Collection) Update(ctx context.Context, filter dal.Filter, doc interface{}) error {
	// build msg
	msg := types.OPUpdateOperation{}
	msg.OPCode = types.OPUpdateCode
	msg.Collection = c.collection
	if err := msg.DOC.Encode(doc); err != nil {
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
	reply := types.OPReply{}
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
func (c *Collection) Delete(ctx context.Context, filter dal.Filter) error {
	// build msg
	msg := types.OPDeleteOperation{}
	msg.OPCode = types.OPDeleteCode
	msg.Collection = c.collection
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
	reply := types.OPReply{}
	err := c.rpc.Call(types.CommandRDBOperation, &msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return nil
}

// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
func (c *Collection) Insert(ctx context.Context, docs interface{}) error {

	// build msg
	msg := types.OPInsertOperation{}
	msg.OPCode = types.OPInsertCode
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
	reply := types.OPReply{}
	err := c.rpc.Call(types.CommandRDBOperation, &msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return nil
}

// CreateIndex 创建索引
func (c *Collection) CreateIndex(ctx context.Context, index dal.Index) error {
	return dal.ErrNotImplemented
}

// DropIndex 移除索引
func (c *Collection) DropIndex(ctx context.Context, indexName string) error {
	return dal.ErrNotImplemented
}

// Indexes 查询索引
func (c *Collection) Indexes(ctx context.Context) ([]dal.Index, error) {
	return nil, dal.ErrNotImplemented
}

// AddColumn 添加字段
func (c *Collection) AddColumn(ctx context.Context, column string, value interface{}) error {
	return dal.ErrNotImplemented
}

// RenameColumn 重命名字段
func (c *Collection) RenameColumn(ctx context.Context, oldName, newColumn string) error {
	return dal.ErrNotImplemented
}

// DropColumn 移除字段
func (c *Collection) DropColumn(ctx context.Context, field string) error {
	return dal.ErrNotImplemented
}

// AggregateOne 聚合查询
func (c *Collection) AggregateOne(ctx context.Context, pipeline interface{}, result interface{}) error {
	// build msg
	msg := types.OPAggregateOperation{}
	msg.OPCode = types.OPAggregateCode
	msg.Collection = c.collection

	if err := msg.Pipiline.Encode(pipeline); err != nil {
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
	reply := types.OPReply{}
	err := c.rpc.Call(types.CommandRDBOperation, msg, &reply)
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

// AggregateAll 聚合查询
func (c *Collection) AggregateAll(ctx context.Context, pipeline interface{}, result interface{}) error {
	return dal.ErrNotImplemented
}
