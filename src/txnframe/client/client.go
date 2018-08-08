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

package client

import (
	"configcenter/src/txnframe/types"
	"context"
	"errors"
)

// Errors defind
var (
	ErrDocumentNotFount = errors.New("document not found")
)

// DALClient db operation interface
type DALClient interface {
	// Collection collection 操作
	Collection(collection string) Collection
	// StartTransaction 开启新事务
	StartTransaction(ctx context.Context, opt JoinOption) (TxDALClient, error)
	// JoinTransaction 加入事务, controller 加入某个事务
	JoinTransaction(JoinOption) TxDALClient
	// NextSequence 获取新序列号(非事务)
	NextSequence(ctx context.Context, sequenceName string) (uint64, error)
	// Ping 健康检查
	Ping() error // 健康检查
}

// Collection collection operation interface
type Collection interface {
	// Find 查询多个并反序列化到 Result
	Find(ctx context.Context, filter types.Filter) Find
	// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
	Insert(ctx context.Context, docs interface{}) error
	// Update 更新数据
	Update(ctx context.Context, filter types.Filter, doc interface{}) error
	// Delete 删除数据
	Delete(ctx context.Context, filter types.Filter) error
	// Count 统计数量(非事务)
	Count(ctx context.Context, filter types.Filter) (uint64, error)
}

// TxDALClient transaction operation interface
type TxDALClient interface {
	DALClient
	// Commit 提交事务
	Commit() error
	// Abort 取消事务
	Abort() error
	// TxnInfo 当前事务信息，用于事务发起者往下传递
	TxnInfo() *types.Tansaction
}

// JoinOption defind join transaction options
type JoinOption struct {
	TxnID     string // 事务ID,uuid
	RequestID string // 请求ID,可选项
	Processor string // 处理进程号，结构为"IP:PORT-PID"用于识别事务session被存于那个TM多活实例
}

// Find find operation interface
type Find interface {
	Fields(fields ...string) Find
	Sort(sort string) Find
	Start(start uint64) Find
	Limit(limit uint64) Find
	All(result interface{}) error
	One(result interface{}) error
}
