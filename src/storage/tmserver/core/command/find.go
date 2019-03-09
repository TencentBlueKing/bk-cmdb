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
	"configcenter/src/storage/mongodb/options/findopt"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/tmserver/core"
	"configcenter/src/storage/types"
)

func init() {
	core.GCommands.SetCommand(types.OPFindCode, &find{})
	core.GCommands.SetCommand(types.OPFindAndModifyCode, &findAndModify{})
}

var _ core.SetDBProxy = (*find)(nil)

type find struct {
	dbProxy mongodb.Client
}

func (d *find) SetDBProxy(db mongodb.Client) {
	d.dbProxy = db
}
func (d *find) Execute(ctx core.ContextParams, decoder rpc.Request) (*types.OPReply, error) {

	msg := types.OPFindOperation{}
	reply := &types.OPReply{}
	if err := decoder.Decode(&msg); nil != err {
		reply.Message = err.Error()
		return reply, err
	}
	blog.V(4).Infof("[MONGO OPERATION] %+v", &msg)

	opt := findopt.Many{}
	opt.Skip = int64(msg.Start)
	opt.Limit = int64(msg.Limit)
	//opt.Sort = msg.Sort

	err := d.dbProxy.Collection(msg.Collection).Find(ctx, msg.Selector, &opt, &reply.Docs)
	if nil == err {
		reply.Success = true
	} else {
		reply.Message = err.Error()
	}

	return reply, err
}

var _ core.SetDBProxy = (*findAndModify)(nil)

type findAndModify struct {
	dbProxy mongodb.Client
}

func (d *findAndModify) SetDBProxy(db mongodb.Client) {
	d.dbProxy = db
}

func (d *findAndModify) Execute(ctx core.ContextParams, decoder rpc.Request) (*types.OPReply, error) {

	msg := types.OPFindAndModifyOperation{}
	reply := &types.OPReply{}
	if err := decoder.Decode(&msg); nil != err {
		reply.Message = err.Error()
		return reply, err
	}
	blog.V(4).Infof("[MONGO OPERATION] %+v", &msg)

	opt := findopt.FindAndModify{}
	opt.Upsert = msg.Upsert
	opt.Remove = msg.Remove
	opt.New = msg.ReturnNew

	var targetCol mongodb.CollectionInterface
	if nil != ctx.Session {
		targetCol = ctx.Session.Collection(msg.Collection)
	} else {
		targetCol = d.dbProxy.Collection(msg.Collection)
	}

	reply.Docs = types.Documents{types.Document{}}
	err := targetCol.FindOneAndModify(ctx, msg.Selector, msg.DOC, &opt, &reply.Docs[0])
	if nil == err {
		reply.Success = true
	} else {
		reply.Message = err.Error()
	}
	return reply, err
}
