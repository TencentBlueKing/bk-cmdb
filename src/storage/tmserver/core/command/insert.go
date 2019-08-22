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
	"configcenter/src/common/util"
	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/tmserver/core"
	"configcenter/src/storage/types"
)

func init() {
	core.GCommands.SetCommand(types.OPInsertCode, &insert{})
}

var _ core.SetDBProxy = (*insert)(nil)

type insert struct {
	dbProxy mongodb.Client
}

func (d *insert) SetDBProxy(db mongodb.Client) {
	d.dbProxy = db
}

func (d *insert) Execute(ctx core.ContextParams, decoder rpc.Request) (*types.OPReply, error) {

	msg := types.OPInsertOperation{}
	reply := &types.OPReply{}
	reply.RequestID = ctx.Header.RequestID
	if err := decoder.Decode(&msg); nil != err {
		reply.Message = err.Error()
		return reply, err
	}
	blog.V(4).Infof("[MONGO OPERATION] db insert, info:%#v, rid:%s", msg, msg.RequestID)

	var targetCol mongodb.CollectionInterface
	if nil != ctx.Session {
		targetCol = ctx.Session.Collection(msg.Collection)
	} else {
		targetCol = d.dbProxy.Collection(msg.Collection)
	}

	slice := util.ConverToInterfaceSlice(msg.DOCS)
	err := targetCol.InsertMany(ctx, slice, nil)
	if nil == err {
		reply.Success = true
	} else {
		blog.ErrorJSON("insert execute error.  errr: %s, raw data: %s, rid:%s", err.Error(), msg, msg.RequestID)
		reply.Message = err.Error()
	}
	return reply, err
}
