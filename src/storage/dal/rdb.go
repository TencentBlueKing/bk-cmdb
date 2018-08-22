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

	"configcenter/src/storage/types"
)

// Errors defind
var (
	ErrTransactionStated   = errors.New("transaction already started")
	ErrTransactionNotFound = errors.New("not in transaction environment")
	ErrDocumentNotFound    = errors.New("document not found")
	ErrNotImplemented      = errors.New("not implemented")
	ErrDuplicated          = errors.New("duplicated")
)

// RDB db operation interface
type RDB interface {
	Clone() RDB
	// Table collection 操作
	Table(collection string) Table
	// StartTransaction 开启新事务
	StartTransaction(ctx context.Context) (RDB, error)
	// Commit 提交事务
	Commit(context.Context) error
	// Abort 取消事务
	Abort(context.Context) error
	// TxnInfo 当前事务信息，用于事务发起者往下传递
	TxnInfo() *types.Transaction
	// NextSequence 获取新序列号(非事务)
	NextSequence(ctx context.Context, sequenceName string) (uint64, error)
	// Ping 健康检查
	Ping() error // 健康检查

	// HasTable 判断是否存在集合
	HasTable(tablename string) (bool, error)
	// DropTable 移除集合
	DropTable(tablename string) error
	// CreateTable 创建集合
	CreateTable(tablename string) error

	IsDuplicatedError(error) bool
	IsNotFoundError(error) bool
}

// Table collection operation interface
type Table interface {
	// Find 查询多个并反序列化到 Result
	Find(filter Filter) Find
	// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
	Insert(ctx context.Context, docs interface{}) error
	// Update 更新数据
	Update(ctx context.Context, filter Filter, doc interface{}) error
	// Delete 删除数据
	Delete(ctx context.Context, filter Filter) error

	// CreateIndex 创建索引
	CreateIndex(ctx context.Context, index Index) error
	// DropIndex 移除索引
	DropIndex(ctx context.Context, indexName string) error
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

type Index struct {
	Keys       map[string]interface{}
	Name       string
	Unique     bool
	Background bool
}
