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
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type modelClassification struct {
	dbProxy dal.RDB
}

func (m *modelClassification) IsExists(ctx core.ContextParams, classificationID string) (bool, error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.ClassFieldClassificationID, Val: ctx.SupplierAccount}, &mongo.Eq{Key: metadata.ClassFieldClassificationID, Val: classificationID})
	cnt, err := m.dbProxy.Table(common.BKTableNameObjClassifiction).Find(cond.ToMapStr()).Count(ctx)
	return 0 != cnt, err
}

func (m *modelClassification) Save(ctx core.ContextParams, classification metadata.Classification) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(ctx, common.BKTableNameObjClassifiction)
	if err != nil {
		return id, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	classification.ID = int64(id)
	classification.OwnerID = ctx.SupplierAccount

	err = m.dbProxy.Table(common.BKTableNameObjClassifiction).Insert(ctx, classification)
	return id, err
}

func (m *modelClassification) CreateOneModelClassification(ctx core.ContextParams, inputParam metadata.CreateOneModelClassification) (*metadata.CreateOneDataResult, error) {

	exists, err := m.IsExists(ctx, inputParam.Data.ClassificationID)
	if nil != err {
		return nil, err
	}
	if exists {
		blog.Errorf("classification (%v)is duplicated", inputParam.Data)
		return nil, ctx.Error.Error(common.CCErrCommDuplicateItem)
	}

	id, err := m.Save(ctx, inputParam.Data)
	return &metadata.CreateOneDataResult{Created: metadata.CreatedDataResult{ID: id}}, err
}

func (m *modelClassification) CreateManyModelClassification(ctx core.ContextParams, inputParam metadata.CreateManyModelClassifiaction) (*metadata.CreateManyDataResult, error) {
	return nil, nil
}
func (m *modelClassification) SetManyModelClassification(ctx core.ContextParams, inputParam metadata.SetManyModelClassification) (*metadata.SetDataResult, error) {
	return nil, nil
}
func (m *modelClassification) SetOneModelClassification(ctx core.ContextParams, inputParam metadata.SetOneModelClassification) (*metadata.SetDataResult, error) {
	return nil, nil
}
func (m *modelClassification) UpdateModelClassification(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error) {
	return nil, nil
}
func (m *modelClassification) DeleteModelClassificaiton(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error) {
	return nil, nil
}
func (m *modelClassification) CascadeDeleteModeClassification(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error) {
	return nil, nil
}
func (m *modelClassification) SearchModelClassification(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {
	return nil, nil
}
