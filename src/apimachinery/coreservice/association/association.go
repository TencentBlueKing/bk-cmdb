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

package association

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type AssociationClientInterface interface {
	UpdateAssociationType(ctx context.Context, h http.Header, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error)
	DeleteAssociationType(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error)
	ReadAssociationType(ctx context.Context, h http.Header, input *metadata.QueryCondition) (resp *metadata.SearchAssociationTypeResult, err error)
	CreateAssociationType(ctx context.Context, h http.Header, input *metadata.CreateAssociationKind) (resp *metadata.CreatedOneOptionResult, err error)

	CreateManyAssociation(ctx context.Context, h http.Header, input *metadata.CreateManyAssociationKind) (resp *metadata.CreatedManyOptionResult, err error)
	SetAssociation(ctx context.Context, h http.Header, input *metadata.SetAssociationKind) (resp *metadata.SetOptionResult, err error)
	SetManyAssociation(ctx context.Context, h http.Header, input *metadata.SetManyAssociationKind) (resp *metadata.SetOptionResult, err error)
	DeleteAssociationCascade(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error)
	CreateModelAssociation(ctx context.Context, h http.Header, input *metadata.CreateModelAssociation) (resp *metadata.CreatedOneOptionResult, err error)
	CreateMainlineModelAssociation(ctx context.Context, h http.Header, input *metadata.CreateModelAssociation) (resp *metadata.CreatedOneOptionResult, err error)
	SetModelAssociation(ctx context.Context, h http.Header, input *metadata.SetModelAssociation) (resp *metadata.SetOptionResult, err error)
	UpdateModelAssociation(ctx context.Context, h http.Header, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error)
	ReadModelAssociation(ctx context.Context, h http.Header, input *metadata.QueryCondition) (resp *metadata.ReadModelAssociationResult, err error)
	DeleteModelAssociation(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error)
	DeleteModelAssociationCascade(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error)
	CreateInstAssociation(ctx context.Context, h http.Header, input *metadata.CreateOneInstanceAssociation) (resp *metadata.CreatedOneOptionResult, err error)
	SetInstAssociation(ctx context.Context, h http.Header, input *metadata.SetOneInstanceAssociation) (resp *metadata.SetOptionResult, err error)
	UpdateInstAssociation(ctx context.Context, h http.Header, input *metadata.UpdateOption) (resp *metadata.UpdatedOptionResult, err error)
	ReadInstAssociation(ctx context.Context, h http.Header, input *metadata.QueryCondition) (resp *metadata.ReadInstAssociationResult, err error)
	DeleteInstAssociation(ctx context.Context, h http.Header, input *metadata.DeleteOption) (resp *metadata.DeletedOptionResult, err error)
}

func NewAssociationClientInterface(client rest.ClientInterface) AssociationClientInterface {
	return &association{client: client}
}

type association struct {
	client rest.ClientInterface
}
