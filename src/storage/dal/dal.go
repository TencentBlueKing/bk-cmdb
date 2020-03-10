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
	"net/http"

	"configcenter/src/storage/dal/types"
)

// RDB rename the RDB into DB
// Compatible stock code
// Deprecated: do not use anymore.
type RDB DB

// DB db operation interface
type DB interface {
	// Table collection 操作
	Table(collection string) types.Table

	// NextSequence 获取新序列号(非事务)
	NextSequence(ctx context.Context, sequenceName string) (uint64, error)
	// Ping 健康检查
	Ping() error // 健康检查

	// HasTable 判断是否存在集合
	HasTable(ctx context.Context, name string) (bool, error)
	// DropTable 移除集合
	DropTable(ctx context.Context, name string) error
	// CreateTable 创建集合
	CreateTable(ctx context.Context, name string) error

	IsDuplicatedError(error) bool
	IsNotFoundError(error) bool

	Close() error
	
	Transaction
}

// Transaction interface
type Transaction interface {
	// StartTransaction 开启新事务
	StartTransaction(ctx *context.Context, h http.Header, opts ...types.TxnOption) (Transaction, error)
	// CommitTransaction 提交事务
	CommitTransaction(context.Context) error
	// AbortTransaction 取消事务
	AbortTransaction(context.Context) error

	// AutoRun Interface for automatic processing of encapsulated transactions
	// f func return error, abort commit, other commit transaction. transaction commit can be error.
	// f func parameter http.header, the handler must be accepted and processed. Subsequent passthrough to call subfunctions and APIs
	// AutoRun(ctx context.Context, opt TxnWrapperOption, f func(header http.Header) error) error
}
