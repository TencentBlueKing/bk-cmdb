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

package core

import (
	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/tmserver/core/transaction"
	"configcenter/src/storage/types"
)

// Core core operation methods
type Core interface {
	ExecuteCommand(ctx ContextParams, input rpc.Request) (*types.OPReply, error)
}

type core struct {
	txn *transaction.Manager
}

// SetTransaction set txc method interface
type SetTransaction interface {
	SetTxn(txn *transaction.Manager)
}

// SetDBProxy set db proxy
type SetDBProxy interface {
	SetDBProxy(db mongodb.Client)
}

// New create a core instance
func New(txnMgr *transaction.Manager, db mongodb.Client) Core {

	for _, cmd := range GCommands.cmds {
		switch tmp := cmd.(type) {
		case SetTransaction:
			tmp.SetTxn(txnMgr)
		case SetDBProxy:
			tmp.SetDBProxy(db)
		}
	}

	return &core{txn: txnMgr}
}

func (c *core) ExecuteCommand(ctx ContextParams, input rpc.Request) (*types.OPReply, error) {

	cmd, ok := GCommands.cmds[ctx.Header.OPCode]
	if !ok {
		reply := types.OPReply{}
		reply.Message = "unknow operation"
		return &reply, nil
	}

	if ctx.Header.TxnID != "" {
		session := c.txn.GetSession(ctx.Header.TxnID)
		if nil == session {
			reply := &types.OPReply{}
			reply.Message = "session not found"
			return reply, nil
		}
		ctx.Session = session.Session
	}

	reply, err := cmd.Execute(ctx, input)
	reply.RequestID = ctx.Header.RequestID
	reply.TxnID = ctx.Header.TxnID
	reply.Processor = ctx.ListenIP
	return reply, err

	/*

		switch header.OPCode {
		case types.OPStartTransaction:
			session, err := t.man.CreateTransaction(header.RequestID)
			if nil != err {
				reply.Message = err.Error()
				return &reply, nil
			}
			reply.Success = true
			reply.TxnID = session.Txninst.TxnID
			reply.Processor = session.Txninst.Processor
			return &reply, nil
		case types.OPCommit:
			err := t.man.Commit(header.TxnID)
			if nil != err {
				reply.Message = err.Error()
				return &reply, nil
			}
			reply.Success = true
			return &reply, nil
		case types.OPAbort:
			err := t.man.Abort(header.TxnID)
			if nil != err {
				reply.Message = err.Error()
				return &reply, nil
			}
			reply.Success = true
			return &reply, nil
		case types.OPInsert, types.OPUpdate, types.OPDelete, types.OPFind, types.OPFindAndModify, types.OPCount:
			var collectionFunc = t.db.Collection
			if transaction != nil {
				collectionFunc = transaction.Collection
			}
			return ExecuteCollection(t.ctx, collectionFunc, header.OPCode, input, &reply)
		default:
			reply.Message = "unknow operation"
			return &reply, nil
		}
	*/
}
