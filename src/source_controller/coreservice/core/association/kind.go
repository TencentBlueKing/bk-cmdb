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

package association

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type associationKind struct {
	dbProxy dal.RDB
}

func (m *associationKind) isExists(ctx core.ContextParams, associationKindID string) (bool, error) {
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: associationKindID})
	cnt, err := m.dbProxy.Table(common.BKTableNameAsstDes).Find(cond.ToMapStr()).Count(ctx)
	return 0 != cnt, err
}

func (m *associationKind) isPrPreAssociationKind(ctx core.ContextParams, associationKindID string) (bool, error) {
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: associationKindID}, &mongo.Eq{Key: common.BKIsPre, Val: true})
	cnt, err := m.dbProxy.Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).Count(ctx)
	return 0 != cnt, err
}

func (m *associationKind) isApply2Object(ctx core.ContextParams, associationKindID string) (bool, error) {
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: associationKindID})
	cnt, err := m.dbProxy.Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).Count(ctx)
	return 0 != cnt, err
}

func (m *associationKind) isApply2Instance(ctx core.ContextParams, associationKindID string) (bool, error) {
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: associationKindID})
	cnt, err := m.dbProxy.Table(common.BKTableNameInstAsst).Find(cond.ToMapStr()).Count(ctx)
	return 0 != cnt, err
}

func (m *associationKind) save(ctx core.ContextParams, associationKind metadata.AssociationKind) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(ctx, common.BKTableNameAsstDes)
	if err != nil {
		return id, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	associationKind.ID = int64(id)

	err = m.dbProxy.Table(common.BKTableNameAsstDes).Insert(ctx, associationKind)
	return id, err
}

func (m *associationKind) CreateAssociationKind(ctx core.ContextParams, inputParam metadata.CreateAssociationKind) (*metadata.CreateOneDataResult, error) {
	exists, err := m.isExists(ctx, inputParam.Data.AssociationKindID)
	if nil != err {
		return nil, err
	}
	if exists {
		blog.Errorf("association kind (%v)is duplicated", inputParam.Data)
		return nil, ctx.Error.Error(common.CCErrCommDuplicateItem)
	}

	id, err := m.save(ctx, inputParam.Data)
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
}

func (m *associationKind) CreateManyAssociationKind(ctx core.ContextParams, inputParam metadata.CreateManyAssociationKind) (*metadata.CreateManyAssociationKind, error) {
	return nil, nil
}
func (m *associationKind) SetAssociationKind(ctx core.ContextParams, inputParam metadata.SetAssociationKind) (*metadata.SetDataResult, error) {
	return nil, nil
}
func (m *associationKind) SetManyAssociationKind(ctx core.ContextParams, inputParam metadata.SetManyAssociationKind) (*metadata.SetDataResult, error) {
	return nil, nil
}
func (m *associationKind) UpdateAssociationKind(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {
	return nil, nil
}
func (m *associationKind) DeleteAssociationKind(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	return nil, nil
}
func (m *associationKind) CascadeDeleteAssociationKind(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {
	return nil, nil
}
func (m *associationKind) SearchAssociationKind(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	return nil, nil
}
