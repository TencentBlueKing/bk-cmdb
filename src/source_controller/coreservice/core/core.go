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

package core

import (
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// ModelAttributeGroup model attribute group methods definitions
type ModelAttributeGroup interface {
	CreateModelAttributeGroup(ctx ContextParams, objID string, inputParam metadata.CreateModelAttributeGroup) (*metadata.CreateOneDataResult, error)
	SetModelAttributeGroup(ctx ContextParams, objID string, inputParam metadata.SetModelAttributeGroup) (*metadata.SetDataResult, error)
	UpdateModelAttributeGroup(ctx ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	UpdateModelAttributeGroupByCondition(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	SearchModelAttributeGroup(ctx ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeGroupDataResult, error)
	SearchModelAttributeGroupByCondition(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeGroupDataResult, error)
	DeleteModelAttributeGroup(ctx ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	DeleteModelAttributeGroupByCondition(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
}

// ModelClassification model classification methods definitions
type ModelClassification interface {
	CreateOneModelClassification(ctx ContextParams, inputParam metadata.CreateOneModelClassification) (*metadata.CreateOneDataResult, error)
	CreateManyModelClassification(ctx ContextParams, inputParam metadata.CreateManyModelClassifiaction) (*metadata.CreateManyDataResult, error)
	SetManyModelClassification(ctx ContextParams, inputParam metadata.SetManyModelClassification) (*metadata.SetDataResult, error)
	SetOneModelClassification(ctx ContextParams, inputParam metadata.SetOneModelClassification) (*metadata.SetDataResult, error)
	UpdateModelClassification(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteModelClassificaiton(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteModeClassification(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	SearchModelClassification(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelClassificationDataResult, error)
}

// ModelAttribute model attribute methods definitions
type ModelAttribute interface {
	CreateModelAttributes(ctx ContextParams, objID string, inputParam metadata.CreateModelAttributes) (*metadata.CreateManyDataResult, error)
	SetModelAttributes(ctx ContextParams, objID string, inputParam metadata.SetModelAttributes) (*metadata.SetDataResult, error)
	UpdateModelAttributes(ctx ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	UpdateModelAttributesByCondition(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteModelAttributes(ctx ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	SearchModelAttributes(ctx ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeDataResult, error)
	SearchModelAttributesByCondition(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeDataResult, error)
}

// ModelAttrUnique model attribute  unique methods definitions
type ModelAttrUnique interface {
	CreateModelAttrUnique(ctx ContextParams, objID string, data metadata.CreateModelAttrUnique) (*metadata.CreateOneDataResult, error)
	UpdateModelAttrUnique(ctx ContextParams, objID string, id uint64, data metadata.UpdateModelAttrUnique) (*metadata.UpdatedCount, error)
	DeleteModelAttrUnique(ctx ContextParams, objID string, id uint64) (*metadata.DeletedCount, error)
	SearchModelAttrUnique(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryUniqueResult, error)
}

// ModelOperation model methods
type ModelOperation interface {
	ModelClassification
	ModelAttributeGroup
	ModelAttribute
	ModelAttrUnique

	CreateModel(ctx ContextParams, inputParam metadata.CreateModel) (*metadata.CreateOneDataResult, error)
	SetModel(ctx ContextParams, inputParam metadata.SetModel) (*metadata.SetDataResult, error)
	UpdateModel(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteModel(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteModel(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	SearchModel(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelDataResult, error)
	SearchModelWithAttribute(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelWithAttributeDataResult, error)
}

// InstanceOperation instance methods
type InstanceOperation interface {
	CreateModelInstance(ctx ContextParams, objID string, inputParam metadata.CreateModelInstance) (*metadata.CreateOneDataResult, error)
	CreateManyModelInstance(ctx ContextParams, objID string, inputParam metadata.CreateManyModelInstance) (*metadata.CreateManyDataResult, error)
	UpdateModelInstance(ctx ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	SearchModelInstance(ctx ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
	DeleteModelInstance(ctx ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteModelInstance(ctx ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
}

// AssociationKind association kind methods
type AssociationKind interface {
	CreateAssociationKind(ctx ContextParams, inputParam metadata.CreateAssociationKind) (*metadata.CreateOneDataResult, error)
	CreateManyAssociationKind(ctx ContextParams, inputParam metadata.CreateManyAssociationKind) (*metadata.CreateManyDataResult, error)
	SetAssociationKind(ctx ContextParams, inputParam metadata.SetAssociationKind) (*metadata.SetDataResult, error)
	SetManyAssociationKind(ctx ContextParams, inputParam metadata.SetManyAssociationKind) (*metadata.SetDataResult, error)
	UpdateAssociationKind(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	DeleteAssociationKind(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteAssociationKind(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	SearchAssociationKind(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
}

// ModelAssociation manager model association
type ModelAssociation interface {
	CreateModelAssociation(ctx ContextParams, inputParam metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error)
	CreateMainlineModelAssociation(ctx ContextParams, inputParam metadata.CreateModelAssociation) (*metadata.CreateOneDataResult, error)
	SetModelAssociation(ctx ContextParams, inputParam metadata.SetModelAssociation) (*metadata.SetDataResult, error)
	UpdateModelAssociation(ctx ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error)
	SearchModelAssociation(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
	DeleteModelAssociation(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
	CascadeDeleteModelAssociation(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
}

// InstanceAssociation manager instance association
type InstanceAssociation interface {
	CreateOneInstanceAssociation(ctx ContextParams, inputParam metadata.CreateOneInstanceAssociation) (*metadata.CreateOneDataResult, error)
	CreateManyInstanceAssociation(ctx ContextParams, inputParam metadata.CreateManyInstanceAssociation) (*metadata.CreateManyDataResult, error)
	SearchInstanceAssociation(ctx ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryResult, error)
	DeleteInstanceAssociation(ctx ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error)
}

// DataSynchronize manager data synchronize interface
type DataSynchronizeOperation interface {
	SynchronizeInstanceAdapter(ctx ContextParams, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error)
	SynchronizeModelAdapter(ctx ContextParams, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error)
	SynchronizeAssociationAdapter(ctx ContextParams, syncData *metadata.SynchronizeParameter) ([]metadata.ExceptionResult, error)
	Find(ctx ContextParams, find *metadata.SynchronizeFindInfoParameter) ([]mapstr.MapStr, uint64, error)
	ClearData(ctx ContextParams, input *metadata.SynchronizeClearDataParameter) error
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
	DataSynchronizeOperation() DataSynchronizeOperation
}

type core struct {
	model           ModelOperation
	instance        InstanceOperation
	associaction    AssociationOperation
	dataSynchronize DataSynchronizeOperation
}

// New create core
func New(model ModelOperation, instance InstanceOperation, association AssociationOperation, dataSynchronize DataSynchronizeOperation) Core {
	return &core{
		model:           model,
		instance:        instance,
		associaction:    association,
		dataSynchronize: dataSynchronize,
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

func (m *core) DataSynchronizeOperation() DataSynchronizeOperation {
	return m.dataSynchronize
}
