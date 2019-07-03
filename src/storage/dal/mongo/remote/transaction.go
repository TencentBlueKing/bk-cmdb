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

// Start create a new transaction
func (c *Mongo) Start(ctx context.Context) (dal.Transcation, error) {

	// build msg
	msg := types.OPStartTransactionOperation{}
	msg.OPCode = types.OPStartTransactionCode

	// set txn
	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		msg.RequestID = opt.RequestID
	}

	// call
	reply := types.OPReply{}
	addr, err := c.rpc.CallInfo(types.CommandRDBOperation, &msg, &reply)
	if err != nil {
		return nil, err
	}
	if !reply.Success {
		return nil, errors.New(reply.Message)
	}

	clone := c.Clone().(*Mongo)
	clone.TxnID = reply.TxnID
	clone.RequestID = reply.RequestID
	clone.tmAddr = addr
	return clone, nil
}

// Commit 提交事务
func (c *Mongo) Commit(ctx context.Context) error {
	msg := types.OPCommitOperation{}
	msg.OPCode = types.OPCommitCode
	msg.RequestID = c.RequestID

	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		msg.TxnID = opt.TxnID
		msg.RequestID = opt.RequestID
	}

	reply := types.OPReply{}
	err := c.rpc.Option(&opt).Call(types.CommandRDBOperation, &msg, &reply)
	c.TxnID = "" // clear TxnID
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return nil
}

// Abort 取消事务
func (c *Mongo) Abort(ctx context.Context) error {
	msg := types.OPAbortOperation{}
	msg.OPCode = types.OPAbortCode
	msg.RequestID = c.RequestID

	opt, ok := ctx.Value(common.CCContextKeyJoinOption).(dal.JoinOption)
	if ok {
		msg.TxnID = opt.TxnID
		msg.RequestID = opt.RequestID
	}

	reply := types.OPReply{}
	err := c.rpc.Option(&opt).Call(types.CommandRDBOperation, &msg, &reply)
	c.TxnID = "" // clear TxnID
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New(reply.Message)
	}
	return nil
}

// TxnInfo 当前事务信息，用于事务发起者往下传递
func (c *Mongo) TxnInfo() *types.Transaction {
	return &types.Transaction{
		RequestID: c.RequestID,
		TxnID:     c.TxnID,
		TMAddr:    c.tmAddr,
	}
}

func (c *Mongo) DB() dal.RDB {
	return c
}
