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

package model

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/source_controller/coreservice/core"
)

func (g *modelAttributeGroup) count(ctx core.ContextParams, cond universalsql.Condition) (count int64, err error) {

	iCount, err := g.dbProxy.Table(common.BKTableNamePropertyGroup).Find(cond.ToMapStr()).Count(ctx)
	return int64(iCount), err
}

func (g *modelAttributeGroup) save(ctx core.ContextParams, group metadata.Group) (uint64, error) {

	id, err := g.dbProxy.NextSequence(ctx, common.BKTableNamePropertyGroup)
	if err != nil {
		return id, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	group.ID = int64(id)
	group.OwnerID = ctx.SupplierAccount

	err = g.dbProxy.Table(common.BKTableNamePropertyGroup).Insert(ctx, group)
	return id, err
}

func (g *modelAttributeGroup) delete(ctx core.ContextParams, cond universalsql.Condition) (uint64, error) {

	cnt, err := g.dbProxy.Table(common.BKTableNamePropertyGroup).Find(cond.ToMapStr()).Count(ctx)
	if nil != err {
		return cnt, err
	}

	err = g.dbProxy.Table(common.BKTableNamePropertyGroup).Delete(ctx, cond.ToMapStr())
	return cnt, err
}

func (g *modelAttributeGroup) search(ctx core.ContextParams, cond universalsql.Condition) ([]metadata.Group, error) {

	dataResult := make([]metadata.Group, 0)
	err := g.dbProxy.Table(common.BKTableNamePropertyGroup).Find(cond.ToMapStr()).All(ctx, &dataResult)
	return dataResult, err
}

func (g *modelAttributeGroup) update(ctx core.ContextParams, data mapstr.MapStr, cond universalsql.Condition) (uint64, error) {

	cnt, err := g.dbProxy.Table(common.BKTableNamePropertyGroup).Find(cond.ToMapStr()).Count(ctx)
	if nil != err {
		return cnt, err
	}
	err = g.dbProxy.Table(common.BKTableNamePropertyGroup).Update(ctx, cond.ToMapStr(), data)
	return cnt, err
}
