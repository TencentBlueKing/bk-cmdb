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
	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/tmserver/core"
	"configcenter/src/storage/types"
)

func init() {
	core.GCommands.SetCommand(types.OPUpdateCode, &update{})
	core.GCommands.SetCommand(types.OPUpdateUnsetCode, &updateUnset{})
	core.GCommands.SetCommand(types.OPUpdateByOperatorCode, &updateByOperator{})
}

var _ core.SetDBProxy = (*update)(nil)

type update struct {
	dbProxy mongodb.Client
}

func (d *update) SetDBProxy(db mongodb.Client) {
	d.dbProxy = db
}

func (d *update) Execute(ctx core.ContextParams, decoder rpc.Request) (*types.OPReply, error) {

	msg := types.OPUpdateOperation{}
	reply := &types.OPReply{}
	reply.RequestID = ctx.Header.RequestID
	if err := decoder.Decode(&msg); nil != err {
		reply.Message = err.Error()
		return reply, err
	}
	blog.V(4).Infof("[MONGO OPERATION] %+v, rid:%s", &msg, msg.RequestID)

	var targetCol mongodb.CollectionInterface
	if nil != ctx.Session {
		targetCol = ctx.Session.Collection(msg.Collection)
	} else {
		targetCol = d.dbProxy.Collection(msg.Collection)
	}

	_, err := targetCol.UpdateMany(ctx, msg.Selector, msg.DOC, nil)
	if nil == err {
		reply.Success = true
	} else {
		blog.ErrorJSON("update execute error.  errr: %s, raw data: %s, rid:%s", err.Error(), msg, msg.RequestID)
		reply.Message = err.Error()
	}
	return reply, err
}

var _ core.SetDBProxy = (*update)(nil)

type updateUnset struct {
	dbProxy mongodb.Client
}

func (d *updateUnset) SetDBProxy(db mongodb.Client) {
	d.dbProxy = db
}

func (d *updateUnset) Execute(ctx core.ContextParams, decoder rpc.Request) (*types.OPReply, error) {

	msg := types.OPUpdateOperation{}
	reply := &types.OPReply{}
	reply.RequestID = ctx.Header.RequestID
	if err := decoder.Decode(&msg); nil != err {
		reply.Message = err.Error()
		return reply, err
	}
	blog.V(4).Infof("[MONGO OPERATION] %+v, rid:%s", &msg, msg.RequestID)

	var targetCol mongodb.CollectionInterface
	if nil != ctx.Session {
		targetCol = ctx.Session.Collection(msg.Collection)
	} else {
		targetCol = d.dbProxy.Collection(msg.Collection)
	}

	updateData := map[string]interface{}{"$unset": msg.DOC}
	_, err := targetCol.Update(ctx, msg.Selector, updateData, nil)
	if nil == err {
		reply.Success = true
	} else {
		blog.ErrorJSON("update execute error.  errr: %s, raw data: %s, rid:%s", err.Error(), msg, msg.RequestID)
		reply.Message = err.Error()
	}
	return reply, err
}

var _ core.SetDBProxy = (*update)(nil)

type updateByOperator struct {
	dbProxy mongodb.Client
}

func (d *updateByOperator) SetDBProxy(db mongodb.Client) {
	d.dbProxy = db
}

func (d *updateByOperator) Execute(ctx core.ContextParams, decoder rpc.Request) (*types.OPReply, error) {

	msg := types.OPUpdateOperation{}
	reply := &types.OPReply{}
	reply.RequestID = ctx.Header.RequestID
	if err := decoder.Decode(&msg); nil != err {
		reply.Message = err.Error()
		return reply, err
	}
	blog.V(4).Infof("[MONGO OPERATION] %+v, rid:%s", &msg, msg.RequestID)

	var targetCol mongodb.CollectionInterface
	if nil != ctx.Session {
		targetCol = ctx.Session.Collection(msg.Collection)
	} else {
		targetCol = d.dbProxy.Collection(msg.Collection)
	}

	_, err := targetCol.Update(ctx, msg.Selector, msg.DOC, nil)
	if nil == err {
		reply.Success = true
	} else {
		blog.ErrorJSON("update execute error.  errr: %s, raw data: %s, rid:%s", err.Error(), msg, msg.RequestID)
		reply.Message = err.Error()
	}
	return reply, err
}
