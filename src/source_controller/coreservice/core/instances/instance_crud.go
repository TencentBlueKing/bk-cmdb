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

package instances

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *instanceManager) save(ctx core.ContextParams, objID string, inputParam mapstr.MapStr) (id uint64, err error) {

	tableName := common.GetInstTableName(objID)
	id, err = m.dbProxy.NextSequence(ctx, tableName)
	if nil != err {
		return id, err
	}
	instIDFieldName := common.GetInstIDField(objID)
	inputParam[instIDFieldName] = id
	if !util.IsInnerObject(objID) {
		inputParam[common.BKObjIDField] = objID
	}
	ts := time.Now()
	inputParam.Set(common.BKOwnerIDField, ctx.SupplierAccount)
	inputParam.Set(common.CreateTimeField, ts)
	inputParam.Set(common.LastTimeField, ts)
	err = m.dbProxy.Table(tableName).Insert(ctx, inputParam)
	return id, err
}

func (m *instanceManager) update(ctx core.ContextParams, objID string, data mapstr.MapStr, cond mapstr.MapStr) (cnt uint64, err error) {
	tableName := common.GetInstTableName(objID)
	if !util.IsInnerObject(objID) {
		cond.Set(common.BKObjIDField, objID)
	}
	cnt, err = m.dbProxy.Table(tableName).Find(cond).Count(ctx)
	if nil != err {
		return cnt, err
	}
	ts := time.Now()
	data.Set(common.LastTimeField, ts)
	data.Remove(common.BKObjIDField)
	err = m.dbProxy.Table(tableName).Update(ctx, cond, data)
	return cnt, err
}

func (m *instanceManager) getInsts(ctx core.ContextParams, objID string, cond mapstr.MapStr) (origins []mapstr.MapStr, exists bool, err error) {
	origins = make([]mapstr.MapStr, 0)
	tableName := common.GetInstTableName(objID)
	if !util.IsInnerObject(objID) {
		cond.Set(common.BKObjIDField, objID)
	}
	if nil != err {
		return origins, false, err
	}
	err = m.dbProxy.Table(tableName).Find(cond).All(ctx, &origins)
	return origins, !m.dbProxy.IsNotFoundError(err), err
}

func (m *instanceManager) getInstDataByID(ctx core.ContextParams, objID string, instID uint64, instanceManager *instanceManager) (origin mapstr.MapStr, err error) {
	tableName := common.GetInstTableName(objID)
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.GetInstIDField(objID), Val: instID})
	if common.GetInstTableName(objID) == common.BKTableNameBaseInst {
		cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: objID})
	}
	err = m.dbProxy.Table(tableName).Find(cond.ToMapStr()).One(ctx, &origin)
	if nil != err {
		return nil, err
	}
	return origin, nil
}

func (m *instanceManager) searchInstance(ctx core.ContextParams, objID string, inputParam metadata.QueryCondition) (results []mapstr.MapStr, err error) {
	results = []mapstr.MapStr{}
	tableName := common.GetInstTableName(objID)
	condition, err := mongo.NewConditionFromMapStr(inputParam.Condition)
	results = make([]mapstr.MapStr, 0)
	if nil != err {
		return results, err
	}
	if tableName == common.BKTableNameBaseInst {
		condition.And(&mongo.Eq{Key: common.BKObjIDField, Val: objID})
	}
	condsMap := util.SetQueryOwner(condition.ToMapStr(), ctx.SupplierAccount)
	blog.V(9).Infof("searchInstance with table: %s and parameters: %#v, rid:%s", tableName, condition.ToMapStr(), ctx.ReqID)
	instHandler := m.dbProxy.Table(tableName).Find(condsMap)
	for _, sort := range inputParam.SortArr {
		fileld := sort.Field
		if sort.IsDsc {
			fileld = "-" + fileld
		}
		instHandler = instHandler.Sort(fileld)
	}
	err = instHandler.Start(uint64(inputParam.Limit.Offset)).Limit(uint64(inputParam.Limit.Limit)).Fields(inputParam.Fields...).All(ctx, &results)
	blog.V(9).Infof("searchInstance with table: %s and parameters: %s, results: %+v, rid: %s", tableName, condition.ToMapStr(), results, ctx.ReqID)

	return results, err
}

func (m *instanceManager) countInstance(ctx core.ContextParams, objID string, cond mapstr.MapStr) (count uint64, err error) {
	tableName := common.GetInstTableName(objID)
	condition, err := mongo.NewConditionFromMapStr(cond)
	if nil != err {
		return 0, err
	}
	if tableName == common.BKTableNameBaseInst {
		condition.And(&mongo.Eq{Key: common.BKObjIDField, Val: objID})
	}

	condsMap := util.SetQueryOwner(condition.ToMapStr(), ctx.SupplierAccount)
	count, err = m.dbProxy.Table(tableName).Find(condsMap).Count(ctx)

	return count, err
}
