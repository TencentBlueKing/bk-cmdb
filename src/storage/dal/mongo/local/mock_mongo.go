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
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"configcenter/src/storage/dal"
	"configcenter/src/storage/types"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mock implement client.DALRDB interface
type Mock struct {
	retval *MockResult
	cache  map[string]*MockResult
}

var _ dal.DB = new(Mock)

// NewMock returns new DB
func NewMock() *Mock {
	return &Mock{
		cache: map[string]*MockResult{},
	}
}

// Close replica client
func (c *Mock) Close() error {
	return nil
}

// Ping replica client
func (c *Mock) Ping() error {
	// TODO
	return nil
}

// Clone return the new client
func (c *Mock) Clone() dal.DB {
	nc := Mock{
		cache: c.cache,
	}
	return &nc
}

// MockResult the result mock
type MockResult struct {
	Err        error
	OK         bool
	RawResult  []byte
	Count      int64
	SequenceID uint64
	Info       types.Transaction
	Indexs     []dal.Index
}

// Mock mock method
func (c *Mock) Mock(retval MockResult) *Mock {
	if len(c.cache) <= 0 {
		c.cache = map[string]*MockResult{}
	}
	c.retval = &retval
	return c
}

// IsDuplicatedError returns whether error is Duplicated Error
func (c *Mock) IsDuplicatedError(err error) bool {
	return err == dal.ErrDuplicated || mgo.IsDup(err)
}

// IsNotFoundError returns whether error is Not Found Error
func (c *Mock) IsNotFoundError(err error) bool {
	return err == dal.ErrDocumentNotFound || err == mgo.ErrNotFound
}

// Table collection operation
func (c *Mock) Table(collName string) dal.Table {
	col := MockCollection{}
	col.collName = collName
	col.Mock = c
	return &col
}

// MockCollection implement client.Collection interface
type MockCollection struct {
	collName string // 集合名
	*Mock
}

// Find 查询多个并反序列化到 Result
func (c *MockCollection) Find(filter dal.Filter) dal.Find {
	return &MockFind{MockCollection: c, filter: filter, projection: map[string]interface{}{"_id": false}}
}

// MockFind define a find operation
type MockFind struct {
	*MockCollection `json:"-"`
	projection      map[string]interface{}
	filter          dal.Filter
	start           int64
	limit           int64
	sort            []string
}

// Fields 查询字段
func (f *MockFind) Fields(fields ...string) dal.Find {

	for _, field := range fields {
		if len(field) <= 0 {
			continue
		}
		f.projection[field] = true
	}
	return f
}

// Sort 查询排序
func (f *MockFind) Sort(sort string) dal.Find {
	if sort != "" {
		f.sort = strings.Split(sort, ",")
	}
	return f
}

// Start 查询上标
func (f *MockFind) Start(start uint64) dal.Find {
	f.start = int64(start)
	return f
}

// Limit 查询限制
func (f *MockFind) Limit(limit uint64) dal.Find {
	f.limit = int64(limit)
	return f
}

// All 查询多个
func (f *MockFind) All(ctx context.Context, result interface{}) error {
	out, err := json.Marshal(f)
	if err != nil {
		return err
	}
	key := "FINDALL:" + f.collName + ":" + string(out)

	if retval, ok := f.Mock.cache[string(key)]; ok {
		raw := bson.Raw{Kind: 4, Data: retval.RawResult}
		err = raw.Unmarshal(result)
		if err != nil {
			return err
		}
		return retval.Err
	}

	bsonout, err := bson.Marshal(result)
	if err != nil {
		return err
	}
	f.Mock.retval.RawResult = bsonout
	f.Mock.cache[string(key)] = f.Mock.retval
	f.Mock.retval = nil
	return nil

}

// One 查询一个
func (f *MockFind) One(ctx context.Context, result interface{}) error {
	out, err := json.Marshal(f)
	if err != nil {
		return err
	}
	key := "FINDONE:" + f.collName + ":" + string(out)
	if retval, ok := f.Mock.cache[string(key)]; ok {
		err = bson.Unmarshal(retval.RawResult, result)
		if err != nil {
			return err
		}
		return retval.Err
	}

	bsonout, err := bson.Marshal(result)
	if err != nil {
		return err
	}
	f.Mock.retval.RawResult = bsonout
	f.Mock.cache[string(key)] = f.Mock.retval
	f.Mock.retval = nil
	return err
}

// Count 统计数量(非事务)
func (f *MockFind) Count(ctx context.Context) (uint64, error) {
	out, err := json.Marshal(f)
	if err != nil {
		return 0, err
	}

	key := "FINDCOUNT:" + f.collName + ":" + string(out)

	if retval, ok := f.Mock.cache[string(key)]; ok {
		return uint64(retval.Count), retval.Err
	}

	return uint64(f.Mock.retval.Count), err
}

// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
func (c *MockCollection) Insert(ctx context.Context, docs interface{}) error {
	bsonout, err := bson.Marshal(docs)
	if err != nil {
		return err
	}

	key := "INSERT:" + c.collName + ":" + string(bsonout)
	if retval, ok := c.Mock.cache[key]; ok {
		return retval.Err
	}

	c.Mock.cache[key] = c.Mock.retval
	c.Mock.retval = nil

	return nil
}

// Update 更新数据
func (c *MockCollection) Update(ctx context.Context, filter dal.Filter, doc interface{}) error {
	bsonout, err := bson.Marshal([]interface{}{filter, doc})
	if err != nil {
		return err
	}

	key := "UPDATE:" + c.collName + ":" + string(bsonout)
	if retval, ok := c.Mock.cache[key]; ok {
		return retval.Err
	}

	c.Mock.cache[key] = c.Mock.retval
	c.Mock.retval = nil

	return nil
}

// UpdateModifyCount 更新数据,返回更新的条数
func (c *MockCollection) UpdateModifyCount(ctx context.Context, filter dal.Filter, doc interface{}) (int64, error) {
	bsonout, err := bson.Marshal([]interface{}{filter, doc})
	if err != nil {
		return 0, err
	}

	key := "UPDATE:" + c.collName + ":" + string(bsonout)
	retval, ok := c.Mock.cache[key]
	if ok {
		return 0, retval.Err
	}

	c.Mock.cache[key] = c.Mock.retval
	c.Mock.retval = nil

	return retval.Count, nil
}

// Update or insert data
func (c *MockCollection) Upsert(ctx context.Context, filter dal.Filter, doc interface{}) error {
	panic("unimplemented operation")
}

// UpdateMultiModel Update data based on operators.
func (c *MockCollection) UpdateMultiModel(ctx context.Context, filter dal.Filter, updateModel ...dal.ModeUpdate) error {

	return errors.New("not support UpdateOp")
}

// Delete 删除数据
func (c *MockCollection) Delete(ctx context.Context, filter dal.Filter) error {
	bsonout, err := bson.Marshal(filter)
	if err != nil {
		return err
	}

	key := "DELETE:" + c.collName + ":" + string(bsonout)
	if retval, ok := c.Mock.cache[key]; ok {
		return retval.Err
	}

	c.Mock.cache[key] = c.Mock.retval
	c.Mock.retval = nil

	return nil
}

// NextSequence 获取新序列号(非事务)
func (c *Mock) NextSequence(ctx context.Context, sequenceName string) (uint64, error) {

	key := "NEXT_SEQUENCE:" + sequenceName
	if retval, ok := c.cache[key]; ok {
		seq := retval.SequenceID
		retval.SequenceID++
		return seq, retval.Err
	}

	c.cache[key] = c.retval
	c.retval = nil

	return 0, nil

}

// StartSession 开启新会话
func (c *Mock) StartSession() (dal.DB, error) {
	return c, nil
}

// EndSession 结束会话
func (c *Mock) EndSession(ctx context.Context) error {
	return nil
}

// StartTransaction 开启新事务
func (c *Mock) StartTransaction(ctx context.Context) error {
	key := "StartTransaction"
	if retval, ok := c.cache[key]; ok {
		return retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil
	return nil
}

// CommitTransaction 提交事务
func (c *Mock) CommitTransaction(ctx context.Context) error {
	key := "COMMIT"
	if retval, ok := c.cache[key]; ok {
		return retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil
	return nil
}

// AbortTransaction 取消事务
func (c *Mock) AbortTransaction(ctx context.Context) error {
	key := "ABORT"
	if retval, ok := c.cache[key]; ok {
		return retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil
	return nil
}

// StartTransaction 开启新事务
//func (c *Mock) StartTransaction(ctx context.Context) (dal.DB, error) {
//	key := "StartTransaction"
//	if retval, ok := c.cache[key]; ok {
//		return c, retval.Err
//	}
//	c.cache[key] = c.retval
//	c.retval = nil
//	return c, nil
//}

func (c *Mock) Start(ctx context.Context) (dal.Transaction, error) {
	key := "StartTransaction"
	if retval, ok := c.cache[key]; ok {
		return c, retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil
	return c, nil
}

// Commit 提交事务
func (c *Mock) Commit(ctx context.Context) error {
	key := "COMMIT"
	if retval, ok := c.cache[key]; ok {
		return retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil
	return nil
}

// Abort 取消事务
func (c *Mock) Abort(ctx context.Context) error {
	key := "ABORT"
	if retval, ok := c.cache[key]; ok {
		return retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil
	return nil
}

// TxnInfo 当前事务信息，用于事务发起者往下传递
func (c *Mock) TxnInfo() (*types.Transaction, error) {
	key := "TxnInfo"
	if retval, ok := c.cache[key]; ok {
		return &retval.Info, nil
	}
	c.cache[key] = c.retval
	c.retval = nil
	return &types.Transaction{}, nil
}

// TxnInfo 当前事务信息，用于事务发起者往下传递
//func (c *Mock) TxnInfo() *types.Transaction {
//	key := "TxnInfo"
//	if retval, ok := c.cache[key]; ok {
//		return &retval.Info
//	}
//	c.cache[key] = c.retval
//	c.retval = nil
//	return &types.Transaction{}
//}

// AutoRun Interface for automatic processing of encapsulated transactions
// f func return error, abort commit, other commit transcation. transcation commit can be error.
// f func parameter http.header, the handler must be accepted and processed. Subsequent passthrough to call subfunctions and APIs
func (c *Mock) AutoRun(ctx context.Context, opt dal.TxnWrapperOption, f func(header http.Header) error) error {
	panic("transcation wrapper not implemented")
}

// HasTable 判断是否存在集合
func (c *Mock) HasTable(ctx context.Context, collName string) (bool, error) {
	key := "HAS_TABLE" + collName
	if retval, ok := c.cache[key]; ok {
		return retval.OK, retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil
	return false, nil
}

// DropTable 移除集合
func (c *Mock) DropTable(ctx context.Context, collName string) error {
	key := "HAS_TABLE:" + collName
	if retval, ok := c.cache[key]; ok {
		return retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil
	return nil

}

// CreateTable 创建集合
func (c *Mock) CreateTable(ctx context.Context, collName string) error {
	key := "CREATE_TABLE:" + collName
	if retval, ok := c.cache[key]; ok {
		return retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil
	return nil
}

// CreateIndex 创建索引
func (c *MockCollection) CreateIndex(ctx context.Context, index dal.Index) error {
	bsonout, err := bson.Marshal(index)
	if err != nil {
		return err
	}
	key := "CREATE_INDEX:" + c.collName + ":" + string(bsonout)
	if retval, ok := c.cache[key]; ok {
		return retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil
	return nil
}

// DropIndex 移除索引
func (c *MockCollection) DropIndex(ctx context.Context, indexName string) error {
	key := "DROP_INDEX:" + c.collName + ":" + indexName
	if retval, ok := c.cache[key]; ok {
		return retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil
	return nil
}

func (c *MockCollection) Indexes(ctx context.Context) ([]dal.Index, error) {
	key := "INDEXES:" + c.collName
	if retval, ok := c.cache[key]; ok {
		return retval.Indexs, retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil
	return nil, nil
}

// AddColumn 添加字段
func (c *MockCollection) AddColumn(ctx context.Context, column string, value interface{}) error {
	bsonout, err := bson.Marshal(value)
	if err != nil {
		return err
	}

	key := "ADD_COLUMN:" + c.collName + ":" + column + ":" + string(bsonout)
	if retval, ok := c.cache[key]; ok {
		return retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil

	return err
}

// RenameColumn 重命名字段
func (c *MockCollection) RenameColumn(ctx context.Context, oldName, newColumn string) error {
	key := "RENAME_COLUMN:" + c.collName + ":" + oldName + ":" + newColumn
	if retval, ok := c.cache[key]; ok {
		return retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil

	return nil
}

// DropColumn 移除字段
func (c *MockCollection) DropColumn(ctx context.Context, field string) error {
	key := "DROP_COLUMN:" + c.collName + ":" + field
	if retval, ok := c.cache[key]; ok {
		return retval.Err
	}
	c.cache[key] = c.retval
	c.retval = nil

	return nil
}

func (c *MockCollection) AggregateAll(ctx context.Context, pipeline interface{}, result interface{}) error {
	out, err := json.Marshal(pipeline)
	if err != nil {
		return err
	}
	key := "AGGREGATE:" + c.collName + ":" + string(out)

	if retval, ok := c.Mock.cache[string(key)]; ok {
		raw := bson.Raw{Kind: 4, Data: retval.RawResult}
		err = raw.Unmarshal(result)
		if err != nil {
			return err
		}
		return retval.Err
	}

	bsonout, err := bson.Marshal(result)
	if err != nil {
		return err
	}
	c.Mock.retval.RawResult = bsonout
	c.Mock.cache[string(key)] = c.Mock.retval
	c.Mock.retval = nil
	return nil
}

func (c *MockCollection) AggregateOne(ctx context.Context, pipeline interface{}, result interface{}) error {
	out, err := json.Marshal(pipeline)
	if err != nil {
		return err
	}
	key := "AGGREGATE:" + c.collName + ":" + string(out)

	if retval, ok := c.Mock.cache[string(key)]; ok {
		raw := bson.Raw{Kind: 4, Data: retval.RawResult}
		err = raw.Unmarshal(result)
		if err != nil {
			return err
		}
		return retval.Err
	}

	bsonout, err := bson.Marshal(result)
	if err != nil {
		return err
	}
	c.Mock.retval.RawResult = bsonout
	c.Mock.cache[string(key)] = c.Mock.retval
	c.Mock.retval = nil
	return nil
}
