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

package dal

import (
	"context"
	"errors"
	"net/http"

	ccErr "configcenter/src/common/errors"
	"configcenter/src/storage/types"
)

// Errors defind
var (
	ErrTransactionStated   = errors.New("transaction already started")
	ErrTransactionNotFound = errors.New("not in transaction environment")
	ErrDocumentNotFound    = errors.New("document not found")
	ErrNotImplemented      = errors.New("not implemented")
	ErrDuplicated          = errors.New("duplicated")
	ErrSessionNotStarted   = errors.New("session is not started")

	UpdateOpAddToSet = "addToSet"
	UpdateOpPull     = "pull"
)

// RDB rename the RDB into DB
// Compatible stock code
type RDB DB

// DB db operation interface
type DB interface {
	Clone() DB
	// Table collection 操作
	Table(collection string) Table

	// NextSequence 获取新序列号(非事务)
	NextSequence(ctx context.Context, sequenceName string) (uint64, error)
	// Ping 健康检查
	Ping() error // 健康检查

	// HasTable 判断是否存在集合
	HasTable(ctx context.Context, tablename string) (bool, error)
	// DropTable 移除集合
	DropTable(ctx context.Context, tablename string) error
	// CreateTable 创建集合
	CreateTable(ctx context.Context, tablename string) error

	IsDuplicatedError(error) bool
	IsNotFoundError(error) bool

	Close() error

	//StartSession 开启会话
	StartSession() (DB, error)
	// EndSession 结束会话
	EndSession(ctx context.Context) error

	Transaction
}

// Transcation db transcation interface
type Transaction interface {
	// StartTransaction 开启新事务
	StartTransaction(context.Context) error
	// CommitTransaction 提交事务
	CommitTransaction(context.Context) error
	// AbortTransaction 取消事务
	AbortTransaction(context.Context) error

	// Start 开启新事务
	Start(ctx context.Context) (Transaction, error)
	// Commit 提交事务
	Commit(context.Context) error
	// Abort 取消事务
	Abort(context.Context) error
	// TxnInfo 当前事务信息，用于事务发起者往下传递
	TxnInfo() (*types.Transaction, error)

	// AutoRun Interface for automatic processing of encapsulated transactions
	// f func return error, abort commit, other commit transcation. transcation commit can be error.
	// f func parameter http.header, the handler must be accepted and processed. Subsequent passthrough to call subfunctions and APIs
	AutoRun(ctx context.Context, opt TxnWrapperOption, f func(header http.Header) error) error
}

// Table collection operation interface
type Table interface {
	// Find 查询多个并反序列化到 Result
	Find(filter Filter) Find
	// Aggregate 聚合查询
	AggregateOne(ctx context.Context, pipeline interface{}, result interface{}) error
	AggregateAll(ctx context.Context, pipeline interface{}, result interface{}) error
	// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
	Insert(ctx context.Context, docs interface{}) error
	// Update 更新数据
	Update(ctx context.Context, filter Filter, doc interface{}) error
	// update or insert data
	Upsert(ctx context.Context, filter Filter, doc interface{}) error
	// UpdateMultiModel  data based on operators.
	UpdateMultiModel(ctx context.Context, filter Filter, updateModel ...ModeUpdate) error
	// UpdateModifyCount 更新数据,返回更新的条数
	UpdateModifyCount(ctx context.Context, filter Filter, doc interface{}) (int64, error)

	// Delete 删除数据
	Delete(ctx context.Context, filter Filter) error

	// CreateIndex 创建索引
	CreateIndex(ctx context.Context, index Index) error
	// DropIndex 移除索引
	DropIndex(ctx context.Context, indexName string) error
	// Indexes 查询索引
	Indexes(ctx context.Context) ([]Index, error)

	// AddColumn 添加字段
	AddColumn(ctx context.Context, column string, value interface{}) error
	// RenameColumn 重命名字段
	RenameColumn(ctx context.Context, oldName, newColumn string) error
	// DropColumn 移除字段
	DropColumn(ctx context.Context, field string) error
}

// JoinOption defind join transaction options
type JoinOption struct {
	TxnID     string // 事务ID,uuid
	RequestID string // 请求ID,可选项
	Processor string // 处理进程号，结构为"IP:PORT-PID"用于识别事务session被存于那个TM多活实例

	TMAddr string // TMServer IP. 存放事务对应的db session 存在TMServer地址的IP

	SessionID    string // 会话ID
	SessionState string // 会话状态
	TxnNumber    string // 事务Number

}

type TxnWrapperOption struct {
	Header http.Header
	CCErr  ccErr.DefaultCCErrorIf
}

// Find find operation interface
type Find interface {
	// Fields 设置查询字段
	Fields(fields ...string) Find
	// Sort 设置查询排序
	Sort(sort string) Find
	// Start 设置限制查询上标
	Start(start uint64) Find
	// Limit 设置查询数量
	Limit(limit uint64) Find
	// All 查询多个
	All(ctx context.Context, result interface{}) error
	// One 查询单个
	One(ctx context.Context, result interface{}) error
	// Count 统计数量(非事务)
	Count(ctx context.Context) (uint64, error)
}

// Index define the DB index struct
type Index struct {
	Keys       map[string]int32 `json:"keys" bson:"key"`
	Name       string           `json:"name" bson:"name"`
	Unique     bool             `json:"unique" bson:"unique"`
	Background bool             `json:"background" bson:"background"`
}

// ModeUpdate  根据不同的操作符去更新数据
type ModeUpdate struct {
	Op  string
	Doc interface{}
}
