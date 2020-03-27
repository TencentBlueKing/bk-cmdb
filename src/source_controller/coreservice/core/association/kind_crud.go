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

package association

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *associationKind) isExists(ctx core.ContextParams, associationKindID string) (origin *metadata.AssociationKind, exists bool, err error) {
	cond := mongo.NewCondition()
	origin = &metadata.AssociationKind{}
	cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: associationKindID})
	err = m.dbProxy.Table(common.BKTableNameAsstDes).Find(cond.ToMapStr()).One(ctx, origin)
	if m.dbProxy.IsNotFoundError(err) {
		return origin, !m.dbProxy.IsNotFoundError(err), nil
	}
	return origin, !m.dbProxy.IsNotFoundError(err), err
}

func (m *associationKind) hasModel(ctx core.ContextParams, cond mapstr.MapStr) (cnt uint64, exists bool, err error) {
	cnt, err = m.dbProxy.Table(common.BKTableNameObjDes).Find(cond).Count(ctx)
	exists = 0 != cnt
	return cnt, exists, err
}

func (m *associationKind) update(ctx core.ContextParams, data mapstr.MapStr, cond mapstr.MapStr) error {
	return m.dbProxy.Table(common.BKTableNameAsstDes).Update(ctx, cond, data)
}

func (m *associationKind) countInstanceAssociation(ctx core.ContextParams, cond mapstr.MapStr) (count uint64, err error) {
	count, err = m.dbProxy.Table(common.BKTableNameAsstDes).Find(cond).Count(ctx)

	return count, err
}

func (m *associationKind) isPreAssociationKind(ctx core.ContextParams, cond metadata.DeleteOption) (exists bool, err error) {
	condition := mapstr.MapStr{}
	for key, val := range cond.Condition {
		condition[key] = val
	}
	condition[common.BKIsPre] = true
	innerCnt, err := m.dbProxy.Table(common.BKTableNameAsstDes).Find(condition).Count(ctx)
	exists = 0 != innerCnt
	return exists, err
}

func (m *associationKind) isApplyToObject(ctx core.ContextParams, cond metadata.DeleteOption) (cnt uint64, exists bool, err error) {

	innerCnt, err := m.dbProxy.Table(common.BKTableNameAsstDes).Find(cond).Count(ctx)
	exists = 0 != innerCnt
	return innerCnt, exists, err
}

func (m *associationKind) save(ctx core.ContextParams, associationKind metadata.AssociationKind) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(ctx, common.BKTableNameAsstDes)
	if err != nil {
		return id, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	associationKind.ID = int64(id)
	associationKind.OwnerID = ctx.SupplierAccount

	err = m.dbProxy.Table(common.BKTableNameAsstDes).Insert(ctx, associationKind)
	return id, err
}

func (m *associationKind) searchAssociationKind(ctx core.ContextParams, inputParam metadata.QueryCondition) (results []metadata.AssociationKind, err error) {
	results = []metadata.AssociationKind{}
	instHandler := m.dbProxy.Table(common.BKTableNameAsstDes).Find(inputParam.Condition).Fields(inputParam.Fields...)
	for _, sort := range inputParam.SortArr {
		fileld := sort.Field
		if sort.IsDsc {
			fileld = "-" + fileld
		}
		instHandler = instHandler.Sort(fileld)
	}
	err = instHandler.Start(uint64(inputParam.Limit.Offset)).Limit(uint64(inputParam.Limit.Limit)).All(ctx, &results)

	return results, err
}
