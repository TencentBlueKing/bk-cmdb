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
	"configcenter/src/storage/types"
)

// Find define a find operation
type Find struct {
	*Collection
	msg *types.OPFindOperation
}

// Fields 查询字段
func (f *Find) Fields(fields ...string) dal.Find {
	projection := types.Document{}
	for _, field := range fields {
		projection[field] = true
	}
	f.msg.Projection = projection
	return f
}

// Sort 查询排序
func (f *Find) Sort(sort string) dal.Find {
	f.msg.Sort = sort
	return f
}

// Start 查询上标
func (f *Find) Start(start uint64) dal.Find {
	f.msg.Start = start
	return f
}

// Limit 查询限制
func (f *Find) Limit(limit uint64) dal.Find {
	f.msg.Limit = limit
	return f
}

// All 查询多个
func (f *Find) All(ctx context.Context, result interface{}) error {
	// set txn
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		f.msg.RequestID = opt.RequestID
		f.msg.TxnID = opt.TxnID
	}
	if f.TxnID != "" {
		f.msg.TxnID = f.TxnID
	}

	// call
	reply := types.OPReply{}
	err := f.rpc.Call(types.CommandRDBOperation, f.msg, &reply)
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return reply.Docs.Decode(result)
}

// One 查询一个
func (f *Find) One(ctx context.Context, result interface{}) error {
	// set txn
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		f.msg.RequestID = opt.RequestID
		f.msg.TxnID = opt.TxnID
	}
	if f.TxnID != "" {
		f.msg.TxnID = f.TxnID
	}

	// call
	reply := types.OPReply{}
	err := f.rpc.Call(types.CommandRDBOperation, f.msg, &reply)
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

// Count 统计数量(非事务)
func (f *Find) Count(ctx context.Context) (uint64, error) {
	// build msg
	f.msg.OPCode = types.OPCountCode

	// set txn
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		f.msg.RequestID = opt.RequestID
	}

	// call
	reply := types.OPReply{}
	err := f.rpc.Call(types.CommandRDBOperation, f.msg, &reply)
	if err != nil {
		return 0, err
	}
	if !reply.Success {
		return 0, errors.New(reply.Message)
	}
	return reply.Count, nil
}
