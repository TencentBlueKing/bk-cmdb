/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package model TODO
package model

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

// ModelClientInterface TODO
type ModelClientInterface interface {
	CreateManyModelClassification(ctx context.Context, h http.Header, input *metadata.CreateManyModelClassifiaction) (
		resp *metadata.CreatedManyOptionResult, err error)
	CreateModelClassification(ctx context.Context, h http.Header, input *metadata.CreateOneModelClassification) (
		*metadata.CreateOneDataResult, error)
	SetManyModelClassification(ctx context.Context, h http.Header, input *metadata.SetManyModelClassification) (
		resp *metadata.SetOptionResult, err error)
	SetModelClassification(ctx context.Context, h http.Header, input *metadata.SetOneModelClassification) (
		resp *metadata.SetOptionResult, err error)
	UpdateModelClassification(ctx context.Context, h http.Header, input *metadata.UpdateOption) (*metadata.UpdatedCount,
		error)
	DeleteModelClassification(ctx context.Context, h http.Header, input *metadata.DeleteOption) (*metadata.DeletedCount,
		error)
	ReadModelClassification(ctx context.Context, h http.Header, input *metadata.QueryCondition) (
		*metadata.QueryModelClassificationDataResult, error)
	CreateModel(ctx context.Context, h http.Header, input *metadata.CreateModel) (*metadata.CreateOneDataResult, error)
	SetModel(ctx context.Context, h http.Header, input *metadata.SetModel) (resp *metadata.SetOptionResult, err error)
	UpdateModel(ctx context.Context, h http.Header, input *metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteModel(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult,
		err error)
	DeleteModelCascade(ctx context.Context, h http.Header, modelID int64) (*metadata.DeletedCount, error)
	ReadModelWithAttribute(ctx context.Context, h http.Header, input *metadata.QueryCondition) (
		resp *metadata.ReadModelWithAttributeResult, err error)
	// ReadModel TODO
	// TODO replace the calling of ReadModelWithAttribute that do not need the object's attributes with this
	ReadModel(ctx context.Context, h http.Header, input *metadata.QueryCondition) (*metadata.QueryModelDataResult,
		error)
	CreateModelAttrs(ctx context.Context, h http.Header, objID string, input *metadata.CreateModelAttributes) (
		*metadata.CreateManyDataResult, error)
	SetModelAttrs(ctx context.Context, h http.Header, objID string, input *metadata.SetModelAttributes) (
		resp *metadata.SetOptionResult, err error)
	UpdateModelAttrs(ctx context.Context, h http.Header, objID string, input *metadata.UpdateOption) (
		*metadata.UpdatedCount, error)
	UpdateModelAttrsIndex(ctx context.Context, h http.Header, objID string, input *metadata.UpdateOption) (
		*metadata.UpdateAttrIndexData, error)
	// UpdateModelAttrsByCondition TODO
	// deprecated, only for old api
	UpdateModelAttrsByCondition(ctx context.Context, h http.Header, input *metadata.UpdateOption) (
		*metadata.UpdatedCount, error)
	DeleteModelAttr(ctx context.Context, h http.Header, objID string, input *metadata.DeleteOption) (
		resp *metadata.DeletedOptionResult, err error)
	ReadModelAttr(ctx context.Context, h http.Header, objID string, input *metadata.QueryCondition) (
		*metadata.QueryModelAttributeDataResult, error)
	// ReadModelAttrByCondition TODO
	// deprecated, only for old api
	ReadModelAttrByCondition(ctx context.Context, h http.Header, input *metadata.QueryCondition) (
		*metadata.QueryModelAttributeDataResult, error)
	GetModelStatistics(ctx context.Context, h http.Header) (resp *metadata.Response, err error)

	ReadAttributeGroup(ctx context.Context, h http.Header, objID string, input metadata.QueryCondition) (
		*metadata.QueryModelAttributeGroupDataResult, error)
	ReadAttributeGroupByCondition(ctx context.Context, h http.Header, input metadata.QueryCondition) (
		*metadata.QueryModelAttributeGroupDataResult, error)
	CreateAttributeGroup(ctx context.Context, h http.Header, objID string, input metadata.CreateModelAttributeGroup) (
		*metadata.CreateOneDataResult, error)
	UpdateAttributeGroup(ctx context.Context, h http.Header, objID string, input metadata.UpdateOption) (
		resp metadata.UpdatedOptionResult, err error)
	UpdateAttributeGroupByCondition(ctx context.Context, h http.Header, input metadata.UpdateOption) (
		*metadata.UpdatedCount, error)
	SetAttributeGroup(ctx context.Context, h http.Header, objID string, input metadata.SetModelAttributes) (
		resp metadata.SetOptionResult, err error)
	DeleteAttributeGroup(ctx context.Context, h http.Header, objID string, input metadata.DeleteOption) (
		resp metadata.DeletedOptionResult, err error)
	DeleteAttributeGroupByCondition(ctx context.Context, h http.Header, input metadata.DeleteOption) (
		*metadata.DeletedCount, error)

	CreateModelAttrUnique(ctx context.Context, h http.Header, objID string, data metadata.CreateModelAttrUnique) (
		*metadata.CreateOneDataResult, error)
	UpdateModelAttrUnique(ctx context.Context, h http.Header, objID string, id uint64,
		data metadata.UpdateModelAttrUnique) (*metadata.UpdatedCount, error)
	DeleteModelAttrUnique(ctx context.Context, h http.Header, objID string, id uint64) (*metadata.DeletedCount, error)
	ReadModelAttrUnique(ctx context.Context, h http.Header, inputParam metadata.QueryCondition) (
		*metadata.QueryUniqueResult, error)

	CreateModelTables(ctx context.Context, h http.Header, input *metadata.CreateModelTable) (err error)
}

// NewModelClientInterface TODO
func NewModelClientInterface(client rest.ClientInterface) ModelClientInterface {
	return &model{client: client}
}

type model struct {
	client rest.ClientInterface
}
