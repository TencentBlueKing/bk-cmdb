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

package service

import (
	"context"

	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/storage/mongobyc"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/server/manager"
	"configcenter/src/storage/types"
)

type TXRPC struct {
	*backbone.Engine
	ctx    context.Context
	rpcsrv *rpc.Server
	man    *manager.TxnManager
	db     mongobyc.Client
	listen string //listening address, use for processor
}

func (t *TXRPC) SetEngine(engine *backbone.Engine) {
	t.Engine = engine
}

func (t *TXRPC) SetDB(db mongobyc.Client) {
	t.db = db
}

func (t *TXRPC) SetMan(man *manager.TxnManager) {
	t.man = man
}

func NewTXRPC(rpcsrv *rpc.Server) *TXRPC {
	txrpc := new(TXRPC)
	txrpc.rpcsrv = rpcsrv

	rpcsrv.Handle(types.CommandRDBOperation, txrpc.RDBOperation)
	rpcsrv.HandleStream(types.CommandWatchTransactionOperation, txrpc.WatchTransaction)
	return txrpc
}

func (t *TXRPC) RDBOperation(input rpc.Request) (interface{}, error) {

	reply := types.OPREPLY{}

	header := types.MsgHeader{}
	err := input.Decode(&header)
	if nil != err {
		reply.Message = err.Error()
		return &reply, nil
	}
	reply.RequestID = header.RequestID
	reply.TxnID = header.TxnID
	reply.Processor = t.listen

	blog.V(3).Infof("RDBOperation %+v", header)

	var transaction mongobyc.Session
	if header.TxnID != "" {
		session := t.man.GetSession(header.TxnID)
		if nil == session {
			reply.Message = "session not found"
			return &reply, nil
		}
		transaction = session.Session
	}

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
}

func (t *TXRPC) WatchTransaction(input rpc.Request, stream rpc.ServerStream) (err error) {
	ch := make(chan *types.Transaction, 100)
	t.man.Subscribe(ch)
	defer t.man.UnSubscribe(ch)
	for txn := range ch {
		if err = stream.Send(txn); err != nil {
			return err
		}
	}
	return nil
}
