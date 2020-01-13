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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
)

func (m *associationKind) isExists(kit *rest.Kit, associationKindID string) (origin *metadata.AssociationKind, exists bool, err error) {
	cond := mongo.NewCondition()
	origin = &metadata.AssociationKind{}
	cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: associationKindID})
	err = m.dbProxy.Table(common.BKTableNameAsstDes).Find(cond.ToMapStr()).One(kit.Ctx, origin)
	if m.dbProxy.IsNotFoundError(err) {
		return origin, !m.dbProxy.IsNotFoundError(err), nil
	}
	return origin, !m.dbProxy.IsNotFoundError(err), err
}

func (m *associationKind) hasModel(kit *rest.Kit, cond mapstr.MapStr) (cnt uint64, exists bool, err error) {
	cnt, err = m.dbProxy.Table(common.BKTableNameObjDes).Find(cond).Count(kit.Ctx)
	exists = 0 != cnt
	return cnt, exists, err
}

func (m *associationKind) update(kit *rest.Kit, data mapstr.MapStr, cond mapstr.MapStr) error {
	return m.dbProxy.Table(common.BKTableNameAsstDes).Update(kit.Ctx, cond, data)
}

func (m *associationKind) countInstanceAssociation(kit *rest.Kit, cond mapstr.MapStr) (count uint64, err error) {
	count, err = m.dbProxy.Table(common.BKTableNameAsstDes).Find(cond).Count(kit.Ctx)

	return count, err
}

func (m *associationKind) isPreAssociationKind(kit *rest.Kit, cond metadata.DeleteOption) (exists bool, err error) {
	condition := mapstr.MapStr{}
	for key, val := range cond.Condition {
		condition[key] = val
	}
	condition[common.BKIsPre] = true
	innerCnt, err := m.dbProxy.Table(common.BKTableNameAsstDes).Find(condition).Count(kit.Ctx)
	exists = 0 != innerCnt
	return exists, err
}

func (m *associationKind) isApplyToObject(kit *rest.Kit, cond metadata.DeleteOption) (cnt uint64, exists bool, err error) {

	innerCnt, err := m.dbProxy.Table(common.BKTableNameAsstDes).Find(cond).Count(kit.Ctx)
	exists = 0 != innerCnt
	return innerCnt, exists, err
}

func (m *associationKind) save(kit *rest.Kit, associationKind metadata.AssociationKind) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(kit.Ctx, common.BKTableNameAsstDes)
	if err != nil {
		return id, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	associationKind.ID = int64(id)
	associationKind.OwnerID = kit.SupplierAccount

	err = m.dbProxy.Table(common.BKTableNameAsstDes).Insert(kit.Ctx, associationKind)
	return id, err
}

func (m *associationKind) searchAssociationKind(kit *rest.Kit, inputParam metadata.QueryCondition) (results []metadata.AssociationKind, err error) {
	results = []metadata.AssociationKind{}
	instHandler := m.dbProxy.Table(common.BKTableNameAsstDes).Find(inputParam.Condition).Fields(inputParam.Fields...)
	err = instHandler.Start(uint64(inputParam.Page.Start)).Limit(uint64(inputParam.Page.Limit)).Sort(inputParam.Page.Sort).All(kit.Ctx, &results)

	return results, err
}
