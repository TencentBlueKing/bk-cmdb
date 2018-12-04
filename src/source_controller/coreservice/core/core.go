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

package core

import (
	"configcenter/src/common/metadata"
)

// ModelClassification model classification methods definitions
type ModelClassification interface {
	CreateOneModelClassification(ctx ContextParams, inputParam metadata.CreateOneModelClassification) (*metadata.CreateOneDataResult, error)
	CreateManyModelClassification(ctx ContextParams, inputParam metadata.CreateManyModelClassifiaction) (*metadata.CreateManyDataResult, error)
	SetManyModelClassification(ctx ContextParams, inputParam metadata.SetManyModelClassification) (*metadata.SetDataResult, error)
	SetOneModelClassification(ctx ContextParams, inputParam metadata.SetOneModelClassification) (*metadata.SetDataResult, error)
	UpdateModelClassification(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error)
	DeleteModelClassificaiton(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error)
	CascadeDeleteModeClassification(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error)
	SearchModelClassification(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
}

// ModelAttribute model attribute methods definitions
type ModelAttribute interface {
	CreateModelAttributes(ctx ContextParams, objID string, inputParam metadata.CreateModelAttributes) (*metadata.CreateManyDataResult, error)
	SetModelAttributes(ctx ContextParams, objID string, inputParam metadata.SetModelAttributes) (*metadata.SetDataResult, error)
	UpdateModelAttributes(ctx ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error)
	DeleteModelAttributes(ctx ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error)
	SearchModelAttributes(ctx ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
}

// ModelOperation model methods
type ModelOperation interface {
	ModelClassification
	ModelAttribute

	CreateModel(ctx ContextParams, inputParam metadata.CreateModel) (*metadata.CreateOneDataResult, error)
	SetModel(ctx ContextParams, inputParam metadata.SetModel) (*metadata.SetDataResult, error)
	UpdateModel(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error)
	DeleteModel(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error)
	CascadeDeleteModel(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error)
	SearchModel(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
}

// InstanceOperation instance methods
type InstanceOperation interface {
	CreateModelInstance(ctx ContextParams, objID string, inputParam metadata.CreateModelInstance) (*metadata.CreateOneDataResult, error)
	CreateManyModelInstance(ctx ContextParams, objID string, inputParam metadata.CreateManyModelInstance) (*metadata.CreateManyDataResult, error)
	SetModelInstance(ctx ContextParams, objID string, inputParam metadata.SetModelInstance) (*metadata.SetDataResult, error)
	SetManyModelInstance(ctx ContextParams, objID string, inputParam metadata.SetManyModelInstance) (*metadata.SetDataResult, error)
	UpdateModelInstance(ctx ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error)
	SearchModelInstance(ctx ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
	DeleteModelInstance(ctx ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error)
	CascadeDeleteModelInstance(ctx ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error)
}

// AssociationKind association kind methods
type AssociationKind interface {
	CreateAssociationKind(ctx ContextParams, inputParam metadata.CreateAssociationKind) (*metadata.CreateOneDataResult, error)
	CreateManyAssociationKind(ctx ContextParams, inputParam metadata.CreateManyAssociationKind) (*metadata.CreateManyAssociationKind, error)
	SetAssociationKind(ctx ContextParams, inputParam metadata.SetAssociationKind) (*metadata.SetDataResult, error)
	SetManyAssociationKind(ctx ContextParams, inputParam metadata.SetManyAssociationKind) (*metadata.SetDataResult, error)
	UpdateAssociationKind(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error)
	DeleteAssociationKind(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error)
	CascadeDeleteAssociationKind(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteOption, error)
	SearchAssociationKind(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
}

// ModelAssociation manager model association
type ModelAssociation interface {
	CreateModelAssociation(ctx ContextParams, inputParam metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error)
	SetModelAssociation(ctx ContextParams, inputParam metadata.SetModelAssociation) (*metadata.SetDataResult, error)
	UpdateModelAssociation(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error)
	SearchModelAssociation(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
	DeleteModelAssociation(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error)
	CascadeDeleteModelAssociation(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error)
}

// InstanceAssociation manager instance association
type InstanceAssociation interface {
	CreateOneInstanceAssociation(ctx ContextParams, inputParam metadata.CreateOneInstanceAssociation) (*metadata.CreateOneDataResult, error)
	SetOneInstanceAssociation(ctx ContextParams, inputParam metadata.SetOneInstanceAssociation) (*metadata.SetDataResult, error)
	CreateManyInstanceAssociation(ctx ContextParams, inputParam metadata.CreateManyInstanceAssociation) (*metadata.CreateManyDataResult, error)
	SetManyInstanceAssociation(ctx ContextParams, inputParam metadata.SetManyInstanceAssociation) (*metadata.SetDataResult, error)
	UpdateInstanceAssociation(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdateDataResult, error)
	SearchInstanceAssociation(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
	DeleteInstanceAssociation(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeleteDataResult, error)
}

// AssociationOperation association methods
type AssociationOperation interface {
	AssociationKind
	ModelAssociation
	InstanceAssociation
}

// Core core itnerfaces methods
type Core interface {
	ModelOperation() ModelOperation
	InstanceOperation() InstanceOperation
	AssociationOperation() AssociationOperation
}

type core struct {
	model        ModelOperation
	instance     InstanceOperation
	associaction AssociationOperation
}

// New create core
func New(model ModelOperation, instance InstanceOperation, association AssociationOperation) Core {
	return &core{
		model:        model,
		instance:     instance,
		associaction: association,
	}
}

func (m *core) ModelOperation() ModelOperation {
	return m.model
}

func (m *core) InstanceOperation() InstanceOperation {
	return m.instance
}

func (m *core) AssociationOperation() AssociationOperation {
	return m.associaction
}
