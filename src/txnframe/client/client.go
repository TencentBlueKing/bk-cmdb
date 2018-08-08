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

var ErrDocumentNotFount = errors.New("document not found")

type DALClient interface {
	Collection(collection string) Collection
	StartTransaction(ctx context.Context, opt JoinOption) (TxDALClient, error) // 开启新事务
	JoinTransaction(JoinOption) TxDALClient                                    // 加入事务, controller 加入某个事务
	NextSequence(ctx context.Context, sequenceName string) (uint64, error)     // 获取新序列号(非事务)
	Ping() error                                                               // 健康检查
}

type Collection interface {
	Find(ctx context.Context, filter types.Filter, result interface{}) error    // 查询多个并反序列化到 Result
	FindOne(ctx context.Context, filter types.Filter, result interface{}) error // 查询单个并反序列化到 Result
	Insert(ctx context.Context, doc interface{}) error                          // 插入单个，如果tag有id, 则回设
	InsertMulti(ctx context.Context, docs []interface{}) error                  // 插入多个, 如果tag有id, 则回设
	Update(ctx context.Context, filter types.Filter, doc interface{}) error     // 更新数据
	Delete(ctx context.Context, filter types.Filter) error                      // 删除数据
	Count(ctx context.Context, filter types.Filter) (uint64, error)             // 统计数量(非事务)
}

type TxDALClient interface {
	Commit() error              // 提交事务
	Abort() error               // 取消事务
	TxnInfo() *types.Tansaction // 当前事务信息，用于事务发起者往下传递
	DALClient
}

type JoinOption struct {
	TxnID     string // 事务ID,uuid
	RequestID string // 请求ID,可选项
	Processor string // 处理进程号，结构为"IP:PORT-PID"用于识别事务session被存于那个TM多活实例
}
