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
	"errors"

	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/storage/mongobyc"
	"configcenter/src/storage/mongobyc/findopt"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/types"
)

var ErrNoSupported = errors.New("not supported")

type CollectionFunc func(collName string) mongobyc.CollectionInterface

// ExecuteCollection execute collection operation to db
func ExecuteCollection(ctx context.Context, collFunc CollectionFunc, opcode types.OPCode, decoder rpc.Request, reply *types.OPREPLY) (*types.OPREPLY, error) {
	executor := &collectionExecutor{ctx: ctx, collection: collFunc, opcode: opcode, msg: decoder, reply: reply}
	executor.execute()
	return executor.reply, executor.execerr
}

type collectionExecutor struct {
	ctx context.Context

	opcode types.OPCode
	msg    rpc.Request

	collection CollectionFunc

	reply   *types.OPREPLY
	execerr error
}

func (e *collectionExecutor) execute() {
	switch e.opcode {
	case types.OPInsert:
		e.insert()
	case types.OPUpdate:
		e.update()
	case types.OPDelete:
		e.delete()
	case types.OPFind:
		e.find()
	case types.OPFindAndModify:
		e.findAndModify()
	case types.OPCount:
		e.count()
	default:
		e.execerr = ErrNoSupported
	}
	if e.execerr != nil {
		e.reply.Success = false
		e.reply.Message = e.execerr.Error()
	} else {
		e.reply.Success = true
	}
}

func (e *collectionExecutor) insert() {
	msg := types.OPINSERT{}
	e.msg.Decode(&msg)
	slice := util.ConverToInterfaceSlice(msg.DOCS)
	e.execerr = e.collection(msg.Collection).InsertMany(e.ctx, slice, nil)
}
func (e *collectionExecutor) update() {
	msg := types.OPUPDATE{}
	e.msg.Decode(&msg)
	_, e.execerr = e.collection(msg.Collection).UpdateMany(e.ctx, msg.Selector, msg.DOC, nil)
}
func (e *collectionExecutor) delete() {
	msg := types.OPDELETE{}
	e.msg.Decode(&msg)
	_, e.execerr = e.collection(msg.Collection).DeleteMany(e.ctx, msg.Selector, nil)
}

func (e *collectionExecutor) find() {
	msg := types.OPFIND{}
	e.msg.Decode(&msg)

	opt := findopt.Many{}
	opt.Skip = int64(msg.Start)
	opt.Limit = int64(msg.Limit)
	opt.Sort = msg.Sort

	blog.Infof("[collectionExecutor] execute %+v", msg)
	e.execerr = e.collection(msg.Collection).Find(e.ctx, msg.Selector, &opt, &e.reply.Docs)
	blog.Infof("[collectionExecutor] find result: %+v, err: [%v]", e.reply.Docs, e.execerr)
}

func (e *collectionExecutor) findAndModify() {
	msg := types.OPFINDANDMODIFY{}
	e.msg.Decode(&msg)
	opt := findopt.FindAndModify{}
	opt.Upsert = msg.Upsert
	opt.Remove = msg.Remove
	opt.New = msg.ReturnNew
	e.execerr = e.collection(msg.Collection).FindAndModify(e.ctx, msg.Selector, msg.DOC, nil, &e.reply.Docs)
}
func (e *collectionExecutor) count() {
	msg := types.OPDELETE{}
	e.msg.Decode(&msg)
	e.reply.Count, e.execerr = e.collection(msg.Collection).Count(e.ctx, msg.Selector)
}
