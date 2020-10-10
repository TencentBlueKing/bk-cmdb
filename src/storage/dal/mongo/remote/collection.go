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
	"reflect"

	"configcenter/src/common"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/mongodb"
	//"configcenter/src/storage/rpc"
	"configcenter/src/storage/types"
)

// Collection implement client.Collection interface
type Collection struct {
	RequestID  string // 请求ID,可选项
	Processor  string // 处理进程号，结构为"IP:PORT-PID"用于识别事务session被存于那个TM多活实例
	TxnID      string // 事务ID,uuid
	collection string // 集合名
	rpc        *pool  //rpc.Client
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
func (c *Collection) Upsert(ctx context.Context, filter dal.Filter, doc interface{}) error {
	panic("unimplemented Upsert operation")
}

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
	err := c.rpc.Option(&opt).Call(types.CommandRDBOperation, &msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return nil
}

// UpdateMultiModel 根据操作符更新数据。
func (c *Collection) UpdateMultiModel(ctx context.Context, filter dal.Filter, updateModel ...dal.ModeUpdate) error {

	doc := make(map[string]interface{}, 0)
	for _, item := range updateModel {
		if _, ok := doc[item.Op]; ok {
			return errors.New(item.Op + " appear multiple times")
		}
		doc["$"+item.Op] = item.Doc
	}

	// build msg
	msg := types.OPUpdateOperation{}
	msg.OPCode = types.OPUpdateByOperatorCode
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
	err := c.rpc.Option(&opt).Call(types.CommandRDBOperation, &msg, &reply)
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
	err := c.rpc.Option(&opt).Call(types.CommandRDBOperation, &msg, &reply)
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
	err := c.rpc.Option(&opt).Call(types.CommandRDBOperation, &msg, &reply)
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
	msg := types.OPDDLOperation{
		Command:    types.OPDDLCreateIndexCommand,
		Collection: c.collection,
		Index:      mongodb.Index{Name: index.Name, Keys: index.Keys, Unique: index.Unique, Background: index.Background},
		MsgHeader:  types.MsgHeader{OPCode: types.OPDDLCode},
	}

	// call
	reply := types.OPReply{}
	err := c.rpc.Option(nil).Call(types.CommandRDBOperation, msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}

	return err
}

// DropIndex 移除索引
func (c *Collection) DropIndex(ctx context.Context, indexName string) error {
	msg := types.OPDDLOperation{
		Command:    types.OPDDLDropIndexCommand,
		Collection: c.collection,
		Index:      mongodb.Index{Name: indexName},
		MsgHeader:  types.MsgHeader{OPCode: types.OPDDLCode},
	}

	// call
	reply := types.OPReply{}
	err := c.rpc.Option(nil).Call(types.CommandRDBOperation, msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}

	return err
}

// Indexes 查询索引
func (c *Collection) Indexes(ctx context.Context) ([]dal.Index, error) {
	msg := types.OPDDLOperation{
		Command:    types.OPDDLIndexCommand,
		Collection: c.collection,
		MsgHeader:  types.MsgHeader{OPCode: types.OPDDLCode},
	}

	// call
	reply := types.OPReply{}
	err := c.rpc.Option(nil).Call(types.CommandRDBOperation, msg, &reply)
	if err != nil {
		return nil, err
	}
	if !reply.Success {
		return nil, errors.New(reply.Message)
	}
	indexs := make([]dal.Index, 0)
	err = reply.Docs.Decode(&indexs)
	return indexs, err
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
	// build msg
	msg := types.OPUpdateOperation{}
	msg.OPCode = types.OPUpdateUnsetCode
	msg.Collection = c.collection
	unsetField := map[string]interface{}{field: ""}
	if err := msg.DOC.Encode(unsetField); err != nil {
		return err
	}
	if err := msg.Selector.Encode(nil); err != nil {
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
	err := c.rpc.Option(&opt).Call(types.CommandRDBOperation, &msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return nil
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
	err := c.rpc.Option(&opt).Call(types.CommandRDBOperation, msg, &reply)
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
	err := c.rpc.Option(&opt).Call(types.CommandRDBOperation, msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}

	if len(reply.Docs) <= 0 {
		return dal.ErrDocumentNotFound
	}
	return reply.Docs.Decode(result)
}

// Distinct 查询不重复的值
func (c *Collection) Distinct(ctx context.Context, field string, filter dal.Filter, result interface{}) error {
	// build msg
	msg := types.OPDistinctOperation{}
	msg.OPCode = types.OPDistinctCode
	msg.Field = field
	msg.Collection = c.collection
	msg.Filter.Encode(filter)

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
	err := c.rpc.Option(&opt).Call(types.CommandRDBOperation, msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}

	return decodeDistinctIntoSlice(ctx, reply.DistinctRes, result)
}

func decodeDistinctIntoSlice(ctx context.Context, dbResults []interface{}, results interface{}) error {
	resultv := reflect.ValueOf(results)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		return errors.New("result argument must be a slice address")
	}

	elemt := resultv.Elem().Type().Elem()

	isInt := false
	isUint := false
	isStr := false
	switch elemt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		isInt = true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		isUint = true
	case reflect.String:
		isStr = false
	default:
		return errors.New("not support decode distinct result to " + elemt.Kind().String())
	}

	slice := reflect.MakeSlice(resultv.Elem().Type(), 0, len(dbResults))

	for _, item := range dbResults {

		elemp := reflect.New(elemt)
		switch val := item.(type) {
		case int:
			if isInt {
				elemp.Elem().SetInt(int64(val))
			} else if isUint {
				elemp.Elem().SetUint(uint64(val))
			} else {
				return errors.New("not can int to " + elemt.Kind().String())
			}
		case int8:
			if isInt {
				elemp.Elem().SetInt(int64(val))
			} else if isUint {
				elemp.Elem().SetUint(uint64(val))
			} else {
				return errors.New("not can " + " to int8" + elemt.Kind().String())
			}
		case int16:
			if isInt {
				elemp.Elem().SetInt(int64(val))
			} else if isUint {
				elemp.Elem().SetUint(uint64(val))
			} else {
				return errors.New("not can " + " to int16" + elemt.Kind().String())
			}
		case int32:
			if isInt {
				elemp.Elem().SetInt(int64(val))
			} else if isUint {
				elemp.Elem().SetUint(uint64(val))
			} else {
				return errors.New("not can " + " to int32" + elemt.Kind().String())
			}
		case int64:
			if isInt {
				elemp.Elem().SetInt(val)
			} else if isUint {
				elemp.Elem().SetUint(uint64(val))
			} else {
				return errors.New("not can " + " to int64" + elemt.Kind().String())
			}
		case uint:
			if isInt {
				elemp.Elem().SetInt(int64(val))
			} else if isUint {
				elemp.Elem().SetUint(uint64(val))
			} else {
				return errors.New("not can " + " to uint" + elemt.Kind().String())
			}
		case uint8:
			if isInt {
				elemp.Elem().SetInt(int64(val))
			} else if isUint {
				elemp.Elem().SetUint(uint64(val))
			} else {
				return errors.New("not can " + " to uint8" + elemt.Kind().String())
			}
		case uint16:
			if isInt {
				elemp.Elem().SetInt(int64(val))
			} else if isUint {
				elemp.Elem().SetUint(uint64(val))
			} else {
				return errors.New("not can " + " to uint16" + elemt.Kind().String())
			}
		case uint32:
			if isInt {
				elemp.Elem().SetInt(int64(val))
			} else if isUint {
				elemp.Elem().SetUint(uint64(val))
			} else {
				return errors.New("not can uint32 to" + elemt.Kind().String())
			}
		case uint64:
			if isInt {
				elemp.Elem().SetInt(int64(val))
			} else if isUint {
				elemp.Elem().SetUint(uint64(val))
			} else {
				return errors.New("not can uint64 to" + elemt.Kind().String())
			}
		case string:
			if !isStr {
				return errors.New("not can string to" + elemt.Kind().String())
			}
			elemp.Elem().SetString(val)
		default:
			return errors.New("not can " + elemt.Kind().String())
		}
		slice = reflect.Append(slice, elemp.Elem())

	}
	resultv.Elem().Set(slice)

	return nil
}
