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

// Package association TODO
package association

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

// AssociationClientInterface TODO
type AssociationClientInterface interface {
	UpdateAssociationType(ctx context.Context, h http.Header, input *metadata.UpdateOption) (*metadata.UpdatedCount,
		error)
	DeleteAssociationType(ctx context.Context, h http.Header, input *metadata.DeleteOption) (
		*metadata.DeletedCount, error)
	ReadAssociationType(ctx context.Context, h http.Header, input *metadata.QueryCondition) (
		*metadata.SearchAssociationType, error)
	CreateAssociationType(ctx context.Context, h http.Header, input *metadata.CreateAssociationKind) (
		*metadata.CreateOneDataResult, error)

	CreateManyAssociation(ctx context.Context, h http.Header, input *metadata.CreateManyAssociationKind) (
		*metadata.CreateManyDataResult, error)
	SetAssociation(ctx context.Context, h http.Header, input *metadata.SetAssociationKind) (
		resp *metadata.SetOptionResult, err error)
	SetManyAssociation(ctx context.Context, h http.Header, input *metadata.SetManyAssociationKind) (
		resp *metadata.SetOptionResult, err error)
	DeleteAssociationCascade(ctx context.Context, h http.Header, input *metadata.DeleteOption) (
		resp *metadata.DeletedOptionResult, err error)
	CreateModelAssociation(ctx context.Context, h http.Header, input *metadata.CreateModelAssociation) (
		*metadata.CreateOneDataResult, error)
	CreateMainlineModelAssociation(ctx context.Context, h http.Header, input *metadata.CreateModelAssociation) (
		*metadata.CreateOneDataResult, error)
	SetModelAssociation(ctx context.Context, h http.Header, input *metadata.SetModelAssociation) (
		resp *metadata.SetOptionResult, err error)
	UpdateModelAssociation(ctx context.Context, h http.Header, input *metadata.UpdateOption) (*metadata.UpdatedCount,
		error)
	ReadModelAssociation(ctx context.Context, h http.Header, input *metadata.QueryCondition) (
		*metadata.QueryModelAssociationResult, error)
	DeleteModelAssociation(ctx context.Context, h http.Header, input *metadata.DeleteOption) (*metadata.DeletedCount,
		error)
	DeleteModelAssociationCascade(ctx context.Context, h http.Header, input *metadata.DeleteOption) (
		resp *metadata.DeletedOptionResult, err error)
	CreateInstAssociation(ctx context.Context, h http.Header, input *metadata.CreateOneInstanceAssociation) (
		*metadata.CreateOneDataResult, error)
	CreateManyInstAssociation(ctx context.Context, header http.Header, input *metadata.CreateManyInstanceAssociation) (
		*metadata.CreateManyDataResult, error)
	SetInstAssociation(ctx context.Context, h http.Header, input *metadata.SetOneInstanceAssociation) (
		resp *metadata.SetOptionResult, err error)
	UpdateInstAssociation(ctx context.Context, h http.Header, input *metadata.UpdateOption) (
		resp *metadata.UpdatedOptionResult, err error)
	ReadInstAssociation(ctx context.Context, h http.Header, input *metadata.InstAsstQueryCondition) (
		resp *metadata.QueryInstAssociationResult, err error)
	DeleteInstAssociation(ctx context.Context, h http.Header, input *metadata.InstAsstDeleteOption) (
		*metadata.DeletedCount, error)

	// CountInstanceAssociations counts model instance associations num.
	CountInstanceAssociations(ctx context.Context, header http.Header, objID string, input *metadata.Condition) (
		*metadata.CountResponseContent, error)
}

// NewAssociationClientInterface TODO
func NewAssociationClientInterface(client rest.ClientInterface) AssociationClientInterface {
	return &association{client: client}
}

type association struct {
	client rest.ClientInterface
}
