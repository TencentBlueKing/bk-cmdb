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

package rpc

import (
	"configcenter/src/txnframe/client"
	"configcenter/src/txnframe/types"
	"context"
	"net/rpc"
)

type RPCClient struct {
	*rpc.Client
}

var _ client.DALClient = new(RPCClient)
var _ client.TxDALClient = new(RPCTxClient)

// Find 查询多个并反序列化到 Result
func (c *RPCClient) Find(ctx context.Context, result interface{}, filter types.Filter) error {
	return nil
}

// FindOne 查询单个并反序列化到 Result
func (c *RPCClient) FindOne(ctx context.Context, result interface{}, filter types.Filter) error {
	return nil
}

// Insert 插入单个，如果tag有id, 则回设
func (c *RPCClient) Insert(ctx context.Context, doc types.Document) error {
	return nil
}

// InsertMulti 插入多个, 如果tag有id, 则回设
func (c *RPCClient) InsertMulti(ctx context.Context, docs []types.Document) error {
	return nil
}

// Update 更新数据
func (c *RPCClient) Update(ctx context.Context, doc types.Document, filter types.Filter) error {
	return nil
}

// Delete 删除数据
func (c *RPCClient) Delete(ctx context.Context, filter types.Filter) error {
	return nil
}

// Count 统计数量(非事务)
func (c *RPCClient) Count(ctx context.Context, filter types.Filter) (uint64, error) {
	return 0, nil
}

// NextSequence 获取新序列号(非事务)
func (c *RPCClient) NextSequence(ctx context.Context, sequenceName string) (int64, error) {
	return 0, nil
}

// StartTransaction 开启新事务
func (c *RPCClient) StartTransaction(ctx context.Context) (client.TxDALClient, error) {
	return new(RPCTxClient), nil
}

// JoinTransaction 加入事务, controller 加入某个事务
func (c *RPCClient) JoinTransaction(client.JoinOption) client.TxDALClient {
	return new(RPCTxClient)
}

// Ping 健康检查
func (c *RPCClient) Ping() error {
	return nil
}

type RPCTxClient struct {
	RPCClient
}

// Commit 提交事务
func (c *RPCTxClient) Commit() error {
	return nil

}

// Abort 取消事务
func (c *RPCTxClient) Abort() error {
	return nil

}

// TxnInfo 当前事务信息，用于事务发起者往下传递
func (c *RPCTxClient) TxnInfo() *types.Tansaction {
	return nil
}
