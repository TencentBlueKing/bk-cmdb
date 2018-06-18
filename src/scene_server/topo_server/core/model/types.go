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
	frtypes "configcenter/src/common/mapstr"
	metadata "configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// AssociationType the association type
type AssociationType string

const (
	// MainLineAssociation the main line association
	MainLineAssociation AssociationType = "mainline"

	// CommonAssociation the common association
	CommonAssociation AssociationType = "commonasso"
)

// Operation the saver interface method
type Operation interface {
	IsExists() (bool, error)
	Create() error
	Delete() error
	Update() error
	Save() error
}

// Topo the object topo interface
type Topo interface {
	Current() Object
	Prev() Object
	Next() Object
}

// Association association operation interface declaration
type Association interface {
	Operation
	Parse(data frtypes.MapStr) (*metadata.Association, error)

	GetType() AssociationType
	SetTopo(parent, child Object) error
	GetTopo(obj Object) (Topo, error)
	ToMapStr() (frtypes.MapStr, error)
}

// Group group opeartion interface declaration
type Group interface {
	Operation

	Parse(data frtypes.MapStr) (*metadata.Group, error)
	CreateAttribute() Attribute

	GetAttributes() ([]Attribute, error)

	SetID(groupID string)
	GetID() string

	SetName(groupName string)
	GetName() string

	SetIndex(groupIndex int64)
	GetIndex() int64

	SetSupplierAccount(supplierAccount string)
	GetSupplierAccount() string

	SetDefault(isDefault bool)
	GetDefault() bool

	SetIsPre(isPre bool)
	GetIsPre() bool

	ToMapStr() (frtypes.MapStr, error)
}

// Attribute attribute opeartion interface declaration
type Attribute interface {
	Operation
	Parse(data frtypes.MapStr) (*metadata.Attribute, error)

	SetSupplierAccount(supplierAccount string)
	GetSupplierAccount() string

	GetParentObject() (Object, error)
	GetChildObject() (Object, error)

	SetParentObject(objID string) error
	SetChildObject(objID string) error

	SetObjectID(objectID string)
	GetObjectID() string

	SetID(attributeID string)
	GetID() string

	SetName(attributeName string)
	GetName() string

	SetGroup(grp Group)
	GetGroup() (Group, error)

	SetGroupIndex(attGroupIndex int64)
	GetGroupIndex() int64

	SetUnint(unit string)
	GetUnint() string

	SetPlaceholder(placeHolder string)
	GetPlaceholder() string

	SetIsEditable(isEditable bool)
	GetIsEditable() bool

	SetIsPre(isPre bool)
	GetIsPre() bool

	SetIsReadOnly(isReadOnly bool)
	GetIsReadOnly() bool

	SetIsOnly(isOnly bool)
	GetIsOnly() bool

	SetIsSystem(isSystem bool)
	GetIsSystem() bool

	SetIsAPI(isAPI bool)
	GetIsAPI() bool

	SetType(attributeType string)
	GetType() string

	SetOption(attributeOption interface{})
	GetOption() interface{}

	SetDescription(attributeDescription string)
	GetDescription() string

	SetCreator(attributeCreator string)
	GetCreator() string

	ToMapStr() (frtypes.MapStr, error)
}

// Classification classification operation interface declaration
type Classification interface {
	Operation
	Parse(data frtypes.MapStr) (*metadata.Classification, error)

	GetObjects() ([]Object, error)

	SetID(classificationID string)
	GetID() string

	SetName(classificationName string)
	GetName() string

	SetType(classificationType string)
	GetType() string

	SetSupplierAccount(supplierAccount string)
	GetSupplierAccount() string

	SetIcon(classificationIcon string)
	GetIcon() string

	ToMapStr() (frtypes.MapStr, error)
}

// Object model operation interface declaration
type Object interface {
	Operation

	Parse(data frtypes.MapStr) (*metadata.Object, error)

	IsCommon() bool

	GetMainlineParentObject() (Object, error)
	GetMainlineChildObject() (Object, error)

	GetParentObject() ([]Object, error)
	GetChildObject() ([]Object, error)

	SetMainlineParentObject(objID string) error
	SetMainlineChildObject(objID string) error

	CreateGroup() Group
	CreateAttribute() Attribute

	GetGroups() ([]Group, error)
	GetAttributes() ([]Attribute, error)

	SetClassification(class Classification)
	GetClassification() (Classification, error)

	SetIcon(objectIcon string)
	GetIcon() string

	SetID(objectID string)
	GetID() string

	SetName(objectName string)
	GetName() string

	SetIsPre(isPre bool)
	GetIsPre() bool

	SetIsPaused(isPaused bool)
	GetIsPaused() bool

	SetPosition(position string)
	GetPosition() string

	SetSupplierAccount(supplierAccount string)
	GetSupplierAccount() string

	SetDescription(description string)
	GetDescription() string

	SetCreator(creator string)
	GetCreator() string

	SetModifier(modifier string)
	GetModifier() string

	ToMapStr() (frtypes.MapStr, error)

	GetInstIDFieldName() string
	GetInstNameFieldName() string
	GetObjectType() string
}

// Factory used to create object  classification attribute etd.
type Factory interface {
	CreaetObject(params types.LogicParams) Object
	CreaetClassification(params types.LogicParams) Classification
	CreateAttribute(params types.LogicParams) Attribute
	CreateGroup(params types.LogicParams) Group
	CreateCommonAssociation(params types.LogicParams, obj Object, asstKey string, asstObj Object) Association
	CreateMainLineAssociatin(params types.LogicParams, obj Object, asstKey string, asstObj Object) Association
}
