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
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/txnframe/mongobyc"
	"configcenter/src/txnframe/rpc"
	"configcenter/src/txnframe/server/manager"
	"configcenter/src/txnframe/types"
	"context"
	"fmt"
	"github.com/rs/xid"
)

type TXRPC struct {
	*backbone.Engine
	ctx    context.Context
	rpcsrv *rpc.Server
	man    *manager.TxnManager
	db     mongobyc.Client
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

	rpcsrv.Handle("StartTransaction", txrpc.StartTransaction)
	rpcsrv.Handle("DBOperation", txrpc.DBOperation)
	return txrpc
}

func (t *TXRPC) StartTransaction(input *rpc.Message) (interface{}, error) {
	blog.V(3).Infof("StartTransaction %#v", input)
	fmt.Printf("StartTransaction %#v", input)
	txn := &types.Tansaction{
		TxnID: xid.New().String(),
	}

	blog.V(3).Infof("t: %#v", t)
	blog.V(3).Infof("t.CCErr: %#v", t.CCErr)
	return txn, nil
	// return txn, t.CCErr.CreateDefaultCCErrorIf("zh").Error(common.CCErrCommJSONUnmarshalFailed)
}

func (t *TXRPC) DBOperation(input *rpc.Message) (interface{}, error) {
	blog.V(3).Infof("DBOperation %#v", input)

	reply := new(types.OPREPLY)

	header := types.MsgHeader{}
	err := input.Decode(&header)
	if nil != err {
		reply.Message = err.Error()
		return reply, nil
	}

	var transaction mongobyc.Session
	if header.TxnID != "" {
		session := t.man.GetSession(header.TxnID)
		if nil == session {
			reply.Message = "session not found"
			return reply, nil
		}
		transaction = session.Session
	}

	switch header.OPCode {
	case types.OPStartTransaction:
		session, err := t.man.CreateTransaction(header.RequestID, buildProcessor(t.ServerInfo.IP, t.ServerInfo.Port, t.ServerInfo.Pid))
		if nil != err {
			reply.Message = err.Error()
			return reply, nil
		}
		reply.OK = true
		reply.TxnID = session.Txninst.TxnID
		reply.Processor = session.Txninst.Processor
		return reply, nil
	case types.OPCommit:
		err := t.man.Commit(header.TxnID)
		if nil != err {
			reply.Message = err.Error()
			return reply, nil
		}
		reply.OK = true
		return reply, nil
	case types.OPAbort:
		err := t.man.Abort(header.TxnID)
		if nil != err {
			reply.Message = err.Error()
			return reply, nil
		}
		reply.OK = true
		return reply, nil
	case types.OPInsert, types.OPUpdate, types.OPDelete, types.OPFind, types.OPFindAndModify, types.OPCount:
		var collectionFunc = t.db.Collection
		if header.TxnID != "" {
			collectionFunc = transaction.Collection
		}
		return ExecuteCollection(t.ctx, collectionFunc, header.OPCode, input)
	default:
		reply.Message = "unknow operation"
		return reply, nil
	}
}

func (*TXRPC) Watch(input interface{}, output string) error {
	blog.V(3).Infof("Watch %#v", input)
	return nil
}
func (*TXRPC) Search(input interface{}, output string) error {
	blog.V(3).Infof("Search %#v", input)
	return nil
}
func (*TXRPC) Healthz(input interface{}, output string) error {
	blog.V(3).Infof("Healthz %#v", input)
	return nil
}
func (*TXRPC) Metrics(input interface{}, output string) error {
	blog.V(3).Infof("Metrics %#v", input)
	return nil
}

func buildProcessor(ip string, port uint, pid int) string {
	return fmt.Sprintf("%s:%d-%d", ip, port, pid)
}
