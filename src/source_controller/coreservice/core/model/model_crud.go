/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
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
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *modelManager) isExists(ctx core.ContextParams, cond universalsql.Condition) (oneModel *metadata.ObjectDes, exists bool, err error) {

	oneModel = &metadata.ObjectDes{}
	err = m.dbProxy.Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).One(ctx, oneModel)
	if nil != err && m.dbProxy.IsNotFoundError(err) {
		return oneModel, exists, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}
	exists = !m.dbProxy.IsNotFoundError(err)
	return oneModel, exists, err
}

func (m *modelManager) save(ctx core.ContextParams, model *metadata.ObjectDes) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(ctx, common.BKTableNameObjDes)
	if err != nil {
		return id, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	err = m.dbProxy.Table(common.BKTableNameObjDes).Insert(ctx, model)
	return id, err
}

func (m *modelManager) isValid(ctx core.ContextParams, objID string) (isValid bool, err error) {

	checkCond := mongo.NewCondition()
	checkCond.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})
	checkCond.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: objID})

	cnt, err := m.dbProxy.Table(common.BKTableNameObjDes).Find(checkCond.ToMapStr()).Count(ctx)
	isValid = (0 != cnt)

	return isValid, err
}

func (m *modelManager) update(ctx core.ContextParams, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.dbProxy.Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).Count(ctx)
	if nil != err {
		return cnt, err
	}

	err = m.dbProxy.Table(common.BKTableNameObjDes).Update(ctx, cond.ToMapStr(), data)
	return cnt, err
}
