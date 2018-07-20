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
	"configcenter/src/framework/common"
)

const (
	// PropertyID the property identifier for a object
	PropertyID = "bk_property_id"
	// PropertyName the property name for a object
	PropertyName = "bk_property_name"
	// PropertyGroup the property group for a object
	PropertyGroup = "bk_property_group"
	// PropertyIndex the property index for a object
	PropertyIndex = "bk_property_index"
	// Unit the unit for a object
	Unit = "unit"
	// PlaceHolder the placeholder for the property
	PlaceHolder = "placeholder"
	// IsEditable the editable for the property
	IsEditable = "editable"
	// IsRequired  mark the property status which must be set
	IsRequired = "isrequired"
	// IsReadOnly mark the property status which can not be editable
	IsReadOnly = "isreadonly"
	// IsOnly mark the property is a key
	IsOnly = "isonly"
	// IsSystem mark the property is the system inner used
	IsSystem = "bk_issystem"
	// IsApi mark the property is the api param
	IsApi = "bk_isapi"
	// PropertyType the property type definition
	PropertyType = "bk_property_type"
	// Option the field configuration information
	Option = "option"

	// GroupID the group identifier
	GroupID = "bk_group_id"
	// GroupName the group name
	GroupName = "bk_group_name"
	// GroupIndex the group index
	GroupIndex = "bk_group_index"
	// IsDefault true is default group
	IsDefault = "bk_isdefault"

	// ObjectIcon the icon name for the object
	ObjectIcon = "bk_obj_icon"
	// ObjectID the id for the object
	ObjectID = "bk_obj_id"
	// ObjectName the name for the object
	ObjectName = "bk_obj_name"
	// IsPre mark the inner object
	IsPre = "ispre"
	// IsPaused mark the object status
	IsPaused = "bk_ispaused"
	// Position the position to draw the object in the page
	Position = "position"
	// SupplierAccount the business id
	SupplierAccount = "bk_supplier_account"
	// Description to introduced object
	Description = "description"
	// Creator the creator for the object
	Creator = "creator"
	// Modifier the last modifier
	Modifier = "modifier"

	// ClassificationID the const definition
	ClassificationID = "bk_classification_id"
	// ClassificationName the const definition
	ClassificationName = "bk_classification_name"
	// ClassificationType the const definition
	ClassificationType = "bk_classification_type"
	// ClassificationIcon the const definition
	ClassificationIcon = "bk_classification_icon"
)

// FieldDataType type definition
type FieldDataType string

const (
	// FieldTypeSingleChar the single char filed type
	FieldTypeSingleChar FieldDataType = "singlechar"

	// FieldTypeLongChar the long char field type
	FieldTypeLongChar FieldDataType = "longchar"

	// FieldTypeInt the int field type
	FieldTypeInt FieldDataType = "int"

	// FieldTypeEnum the enum field type
	FieldTypeEnum FieldDataType = "enum"

	// FieldTypeDate the date field type
	FieldTypeDate FieldDataType = "date"

	// FieldTypeTime the time field type
	FieldTypeTime FieldDataType = "time"

	// FieldTypeUser the user field type
	FieldTypeUser FieldDataType = "objuser"

	// FieldTypeSingleAsst the single association
	FieldTypeSingleAsst FieldDataType = "singleasst"

	// FieldTypeMultiAsst the multi association
	FieldTypeMultiAsst FieldDataType = "multiasst"

	// FieldTypeTimeZone the timezone field type
	FieldTypeTimeZone FieldDataType = "timezone"

	// FieldTypeBool the bool type
	FieldTypeBool FieldDataType = "bool"
)

// default
const (
	classificationIconDefault = "icon-cc-default-class"
	objectIconDefault         = "icon-cc-business"
)

// GroupIterator the group iterator
type GroupIterator interface {
	Next() (Group, error)
	ForEach(itemCallback func(item Group) error) error
}

// Group the interface declaration for model maintence
type Group interface {
	IsExists() (bool, error)
	Create() error
	Update() error
	Save() error

	SetID(id string)
	GetID() string
	SetName(name string)
	GetName() string
	SetIndex(idx int)
	GetIndex() int
	SetSupplierAccount(ownerID string)
	GetSupplierAccount() string
	SetDefault()
	SetNonDefault()
	GetDefault() bool

	GetRecordID() int

	CreateAttribute() Attribute
	FindAttributesLikeName(supplierAccount string, attributeName string) (AttributeIterator, error)
	FindAttributesByCondition(supplierAccount string, cond common.Condition) (AttributeIterator, error)
}

// ClassificationIterator the classification iterator
type ClassificationIterator interface {
	Next() (Classification, error)
	ForEach(itemCallback func(item Classification) error) error
}

// Classification the interface declaration for model classification
type Classification interface {
	IsExists() (bool, error)
	Create() error
	Update() error
	Save() error

	GetRecordID() int
	SetID(id string)
	GetID() string
	SetName(name string)
	GetName() string
	SetIcon(iconName string)
	GetIcon() string

	CreateModel() Model
	FindModelsLikeName(supplierAccount string, modelName string) (Iterator, error)
	FindModelsByCondition(supplierAccount string, cond common.Condition) (Iterator, error)
}

// Iterator the model iterator
type Iterator interface {
	Next() (Model, error)
	ForEach(itemCallback func(item Model) error) error
}

// Model the interface declaration for model maintence
type Model interface {
	IsExists() (bool, error)
	Create() error
	Update() error
	Save() error

	SetIcon(iconName string)
	GetIcon() string
	SetID(id string)
	GetID() string
	SetName(name string)
	GetName() string

	SetPaused()
	SetNonPaused()
	Paused() bool

	SetPosition(position string)
	GetPosition() string
	SetSupplierAccount(ownerID string)
	GetSupplierAccount() string
	SetDescription(desc string)
	GetDescription() string
	SetCreator(creator string)
	GetCreator() string
	SetModifier(modifier string)
	GetModifier() string

	CreateAttribute() Attribute
	CreateGroup() Group

	Attributes() ([]Attribute, error)

	FindAttributesLikeName(attributeName string) (AttributeIterator, error)
	FindAttributesByCondition(cond common.Condition) (AttributeIterator, error)

	FindGroupsLikeName(groupName string) (GroupIterator, error)
	FindGroupsByCondition(cond common.Condition) (GroupIterator, error)
}

// AttributeIterator the attribute iterator
type AttributeIterator interface {
	Next() (Attribute, error)
	ForEach(itemCallback func(item Attribute) error) error
}

// Attribute the interface declaration for model attribute maintence
type Attribute interface {
	IsExists() (bool, error)
	Create() error
	Update() error
	Save() error

	GetRecordID() int
	SetID(id string)
	GetID() string
	SetName(name string)
	GetName() string
	SetUnit(unit string)
	GetUnit() string
	SetPlaceholder(placeHolder string)
	GetPlaceholder() string
	SetEditable()
	SetNonEditable()
	GetEditable() bool
	SetRequired()
	SetNonRequired()
	GetRequired() bool
	SetKey(isKey bool)
	GetKey() bool
	SetOption(option interface{})
	GetOption() interface{}
	SetDescrition(des string)
	GetDescription() string

	SetType(dataType FieldDataType)
	GetType() FieldDataType
}
