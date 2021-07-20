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

	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/redis"
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

	//NextSequences 批量获取新序列号(非事务)
	NextSequences(ctx context.Context, sequenceName string, num int) ([]uint64, error)

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

	// CommitTransaction 提交事务
	CommitTransaction(context.Context, *metadata.TxnCapable) error
	// AbortTransaction 取消事务
	AbortTransaction(context.Context, *metadata.TxnCapable) error

	// InitTxnManager TxnID management of initial transaction
	InitTxnManager(r redis.Client) error
}
