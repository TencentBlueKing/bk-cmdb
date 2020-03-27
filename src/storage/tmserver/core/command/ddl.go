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
	"fmt"
	"strings"

	"configcenter/src/common/blog"
	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/tmserver/core"
	"configcenter/src/storage/types"
)

// DDL(Data Definition Language)，是用于描述数据库中要存储的现实世界实体的语言。

func init() {
	core.GCommands.SetCommand(types.OPDDLCode, &ddl{})
}

var _ core.SetDBProxy = (*ddl)(nil)

type ddl struct {
	dbProxy mongodb.Client
}

func (d *ddl) SetDBProxy(db mongodb.Client) {
	d.dbProxy = db
}

func (m *ddl) Execute(ctx core.ContextParams, decoder rpc.Request) (*types.OPReply, error) {

	msg := types.OPDDLOperation{}
	reply := &types.OPReply{}
	reply.RequestID = ctx.Header.RequestID
	if err := decoder.Decode(&msg); nil != err {
		reply.Message = err.Error()
		return reply, err
	}
	blog.V(4).Infof("[MONGO OPERATION] execute ddl operater. msg:%#v, rid:%s", msg, msg.RequestID)

	var db mongodb.Database
	var targetCol mongodb.CollectionInterface
	if nil != ctx.Session {
		db = ctx.Session
		targetCol = ctx.Session.Collection(msg.Collection)
	} else {
		db = m.dbProxy.Database()
		targetCol = m.dbProxy.Collection(msg.Collection)
	}
	var execErr error
	switch msg.Command {
	case types.OPDDLHasCollectCommand:
		var exist bool
		exist, execErr = db.HasCollection(msg.Collection)

		if exist {
			reply.Count = 1
		}
	case types.OPDDLDropCollectCommand:
		execErr = db.DropCollection(msg.Collection)
	case types.OPDDLCreateCollectCommand:
		execErr = db.CreateEmptyCollection(msg.Collection)
	case types.OPDDLCreateIndexCommand:
		// new version mongodb driver, name not support name
		if msg.Index.Name == "" {
			var name string
			for key, val := range msg.Index.Keys {
				name = name + fmt.Sprintf("_%s_%d", key, val)
			}
			msg.Index.Name = strings.Trim(name, "_")
		}
		execErr = targetCol.CreateIndex(msg.Index)
	case types.OPDDLIndexCommand:
		var dbIndexs *mongodb.QueryIndexResult
		dbIndexs, execErr = targetCol.GetIndexes()
		if execErr == nil {
			for _, dbIndex := range dbIndexs.Indexes {
				keys := map[string]int32{}
				for _, key := range dbIndex.Key {
					if strings.HasPrefix(key, "-") {
						key = strings.TrimLeft(key, "-")
						keys[key] = -1
					} else {
						keys[key] = 1
					}
				}
				reply.Docs = append(reply.Docs, types.Document{
					"name": dbIndex.Name,
					"keys": keys,
				})
			}
		}
	case types.OPDDLDropIndexCommand:
		execErr = targetCol.DropIndex(msg.Index.Name)
	default:
		execErr = fmt.Errorf("db data definition language execute operate type %s not support", msg.OPCode)
	}

	//err := targetCol.Find(ctx, msg.Selector, &opt, &reply.Docs)
	if nil == execErr {
		reply.Success = true
	} else {
		reply.Message = execErr.Error()
	}

	return reply, execErr
}
