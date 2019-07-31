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

package command

import (
	"configcenter/src/common/blog"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/tmserver/core"
	"configcenter/src/storage/tmserver/core/session"
	"configcenter/src/storage/types"
)

func init() {
	core.GCommands.SetCommand(types.OPStartTransactionCode, &startTransaction{})
}

type startTransaction struct {
	txn *session.Manager
}

var _ core.SetTransaction = (*startTransaction)(nil)

func (d *startTransaction) SetTxn(txn *session.Manager) {
	d.txn = txn
}

func (d *startTransaction) Execute(ctx core.ContextParams, decoder rpc.Request) (*types.OPReply, error) {
	blog.V(4).Infof("[MONGO OPERATION] %+v", &ctx.Header)
	reply := &types.OPReply{}
	reply.RequestID = ctx.Header.RequestID
	sess, err := d.txn.CreateTransaction(ctx.Header.RequestID)
	if nil != err {
		reply.Message = err.Error()
		return reply, err
	}
	reply.Success = true
	reply.TxnID = sess.Txninst.TxnID
	reply.Processor = sess.Txninst.Processor

	return reply, nil
}
