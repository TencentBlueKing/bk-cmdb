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
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/txnframe/rpc"
	"configcenter/src/txnframe/types"
	"fmt"
	"github.com/rs/xid"
)

type TXRPC struct {
	*backbone.Engine

	rpcsrv *rpc.Server
}

func (t *TXRPC) SetEngine(engine *backbone.Engine) {
	t.Engine = engine
}

func NewTXRPC(rpcsrv *rpc.Server) *TXRPC {
	txrpc := new(TXRPC)
	txrpc.rpcsrv = rpcsrv

	rpcsrv.Handle("StartTransaction", txrpc.StartTransaction)
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
	return txn, t.CCErr.CreateDefaultCCErrorIf("zh").Error(common.CCErrCommJSONUnmarshalFailed)
}

func (*TXRPC) HandleDB(input interface{}, output string) error {
	blog.V(3).Infof("HandleDB %#v", input)
	return nil
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
