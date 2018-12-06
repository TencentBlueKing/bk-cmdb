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
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *modelClassification) IsExists(ctx core.ContextParams, classificationID string) (origin *metadata.Classification, exists bool, err error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.ClassFieldClassificationID, Val: ctx.SupplierAccount}, &mongo.Eq{Key: metadata.ClassFieldClassificationID, Val: classificationID})
	err = m.dbProxy.Table(common.BKTableNameObjClassifiction).Find(cond.ToMapStr()).One(ctx, origin)
	return origin, !m.dbProxy.IsNotFoundError(err), err
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

func (m *modelClassification) Update(ctx core.ContextParams, data mapstr.MapStr, cond mapstr.MapStr) error {

	return m.dbProxy.Table(common.BKTableNameObjClassifiction).Update(ctx, cond, data)
}
