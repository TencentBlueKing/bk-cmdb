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
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type modelClassification struct {
	model   *modelManager
	dbProxy dal.RDB
}

func (m *modelClassification) CreateOneModelClassification(ctx core.ContextParams, inputParam metadata.CreateOneModelClassification) (*metadata.CreateOneDataResult, error) {

	_, exists, err := m.IsExists(ctx, inputParam.Data.ClassificationID)
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

	dataResult := &metadata.CreateManyDataResult{}
	for itemIdx, item := range inputParam.Data {

		_, exists, err := m.IsExists(ctx, item.ClassificationID)
		if nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}

		if exists {
			dataResult.Repeated = append(dataResult.Repeated, metadata.RepeatedDataResult{OriginIndex: int64(itemIdx), Data: mapstr.NewFromStruct(item, "field")})
			continue
		}

		id, err := m.Save(ctx, item)
		if nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}

		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{
			ID: id,
		})

	}

	return dataResult, nil
}
func (m *modelClassification) SetManyModelClassification(ctx core.ContextParams, inputParam metadata.SetManyModelClassification) (*metadata.SetDataResult, error) {

	dataResult := &metadata.SetDataResult{}
	for itemIdx, item := range inputParam.Data {

		origin, exists, err := m.IsExists(ctx, item.ClassificationID)
		if nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}

		if exists {

			cond := mongo.NewCondition()
			if err := m.Update(ctx, mapstr.NewFromStruct(item, "field"), cond.Element(&mongo.Eq{Key: metadata.ClassificationFieldID, Val: origin.ID}).ToMapStr()); nil != err {
				dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
					Message:     err.Error(),
					Code:        int64(err.(errors.CCErrorCoder).GetCode()),
					Data:        item,
					OriginIndex: int64(itemIdx),
				})
				continue
			}

			dataResult.UpdatedCount.Count++
			continue
		}

		id, err := m.Save(ctx, item)
		if nil != err {

			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        item,
				OriginIndex: int64(itemIdx),
			})
			continue
		}

		dataResult.CreatedCount.Count++
		dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{
			ID: id,
		})

	}

	return dataResult, nil
}

func (m *modelClassification) SetOneModelClassification(ctx core.ContextParams, inputParam metadata.SetOneModelClassification) (*metadata.SetDataResult, error) {

	origin, exists, err := m.IsExists(ctx, inputParam.Data.ClassificationID)
	if nil != err {
		return nil, err
	}

	dataResult := &metadata.SetDataResult{}

	if exists {

		cond := mongo.NewCondition()
		if err := m.Update(ctx, mapstr.NewFromStruct(inputParam.Data, "field"), cond.Element(&mongo.Eq{Key: metadata.ClassificationFieldID, Val: origin.ID}).ToMapStr()); nil != err {
			dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
				Message:     err.Error(),
				Code:        int64(err.(errors.CCErrorCoder).GetCode()),
				Data:        inputParam.Data,
				OriginIndex: 0,
			})
			return dataResult, nil
		}
		dataResult.Updated = append(dataResult.Updated, metadata.UpdatedDataResult{ID: uint64(origin.ID)})
		return dataResult, err
	}

	id, err := m.Save(ctx, inputParam.Data)
	if nil != err {
		dataResult.Exceptions = append(dataResult.Exceptions, metadata.ExceptionResult{
			Message:     err.Error(),
			Code:        int64(err.(errors.CCErrorCoder).GetCode()),
			Data:        origin,
			OriginIndex: 0,
		})
	}
	dataResult.Created = append(dataResult.Created, metadata.CreatedDataResult{ID: id})
	return dataResult, err
}

func (m *modelClassification) UpdateModelClassification(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	cnt, err := m.dbProxy.Table(common.BKTableNameObjClassifiction).Find(inputParam.Condition).Count(ctx)
	if nil != err {
		return &metadata.UpdatedCount{}, err
	}
	if err := m.Update(ctx, inputParam.Data, inputParam.Condition); nil != err {
		return &metadata.UpdatedCount{}, err
	}
	return &metadata.UpdatedCount{Count: cnt}, nil
}

func (m *modelClassification) DeleteModelClassificaiton(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	cnt, exists, err := m.hasModel(ctx, inputParam.Condition)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}
	if exists {
		return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrTopoObjectClassificationHasObject)
	}

	m.dbProxy.Table(common.BKTableNameObjClassifiction).Delete(ctx, inputParam.Condition)
	return &metadata.DeletedCount{Count: cnt}, nil
}

func (m *modelClassification) CascadeDeleteModeClassification(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	classificationItems, err := m.searchClassification(ctx, inputParam.Condition)
	if nil != err {
		return &metadata.DeletedCount{}, err
	}

	for _, item := range classificationItems {
		cond := mongo.NewCondition()
		cond.Element(&mongo.Eq{Key: metadata.ModelFieldObjCls, Val: item.ClassificationID})
		if _, err := m.model.cascadeDeleteModel(ctx, cond.ToMapStr()); nil != err {
			return &metadata.DeletedCount{}, err
		}
	}

	return &metadata.DeletedCount{Count: uint64(len(classificationItems))}, nil
}

func (m *modelClassification) SearchModelClassification(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error) {

	classificationItems, err := m.searchClassification(ctx, inputParam.Condition)
	if nil != err {
		return &metadata.QueryResult{}, err
	}

	dataResult := &metadata.QueryResult{}
	dataResult.Count = uint64(len(classificationItems))
	for item := range classificationItems {
		dataResult.Info = append(dataResult.Info, mapstr.NewFromStruct(item, "field"))
	}

	return dataResult, nil
}
