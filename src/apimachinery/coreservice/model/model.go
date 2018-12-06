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

package model

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type ModelClientInterface interface {
	CreateManyModelClassification(ctx context.Context, h http.Header, input *metadata.CreateManyModelClassifiaction) (resp *metadata.CreateManyOptionResult, err error)
	CreateModelClassification(ctx context.Context, h http.Header, input *metadata.CreateOneModelClassification) (resp *metadata.CreatedOneOptionResult, err error)
	SetManyModelClassification(ctx context.Context, h http.Header, input *metadata.SetManyModelClassification) (resp *metadata.SetOptionResult, err error)
	SetModelClassification(ctx context.Context, h http.Header, input *metadata.SetOneModelClassification) (resp *metadata.SetOptionResult, err error)
	UpdateModelClassification(ctx context.Context, h http.Header, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error)
	DeleteModelClassification(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error)
	DeleteModelClassificationCascade(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error)
	ReadModelClassification(ctx context.Context, h http.Header, input *metadata.QueryCondition) (resp *metadata.ReadModelClassifitionResult, err error)
	CreateModel(ctx context.Context, h http.Header, input *metadata.CreateModel) (resp *metadata.CreatedOneOptionResult, err error)
	SetModel(ctx context.Context, h http.Header, input *metadata.SetModel) (resp *metadata.SetOptionResult, err error)
	UpdateModel(ctx context.Context, h http.Header, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error)
	DeleteModel(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error)
	DeleteModelCascade(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error)
	ReadModel(ctx context.Context, h http.Header, input *metadata.QueryCondition) (resp *metadata.ReadModelResult, err error)
	CreateModelAttrs(ctx context.Context, h http.Header, objID string, input *metadata.CreateModelAttributes) (resp *metadata.CreateManyOptionResult, err error)
	SetModelAttrs(ctx context.Context, h http.Header, objID string, input *metadata.SetModelAttributes) (resp *metadata.CreateManyOptionResult, err error)
	UpdateModelAttrs(ctx context.Context, h http.Header, objID string, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error)
	DeleteModelAttr(ctx context.Context, h http.Header, objID string, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error)
	ReadModelAttr(ctx context.Context, h http.Header, objID string, input *metadata.QueryCondition) (resp *metadata.QueryConditionResult, err error)
}

func NewModelClientInterface(client rest.ClientInterface) ModelClientInterface {
	return &model{client: client}
}

type model struct {
	client rest.ClientInterface
}
