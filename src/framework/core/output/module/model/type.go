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
	"configcenter/src/framework/core/types"
)

// GroupIterator the group iterator
type GroupIterator interface {
	Next() (Group, error)
}

// Group the interface declaration for model maintence
type Group interface {
	types.Saver

	SetID(id string)
	GetID() string
	SetName(name string)
	SetIndex(idx int)
	GetIndex() int
	SetSupplierAccount(ownerID string)
	GetSupplierAccount() string
	SetDefault()
	SetNonDefault()
	Default() bool

	CreateAttribute() Attribute
	FindAttributesLikeName(attributeName string) (AttributeIterator, error)
	FindAttributesByCondition(condition *common.Condition) (AttributeIterator, error)
}

// ClassificationIterator the classification iterator
type ClassificationIterator interface {
	Next() (Classification, error)
}

// Classification the interface declaration for model classification
type Classification interface {
	types.Saver

	SetID(id string)
	SetName(name string)
	SetIcon(iconName string)
	GetID() string

	CreateModel() Model
	FindModelsLikeName(modelName string) (Iterator, error)
	FindModelsByCondition(condition *common.Condition) (Iterator, error)
}

// Iterator the model iterator
type Iterator interface {
	Next() (Model, error)
}

// Model the interface declaration for model maintence
type Model interface {
	types.Saver
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

	FindAttributesLikeName(attributeName string) (AttributeIterator, error)
	FindAttributesByCondition(condition *common.Condition) (AttributeIterator, error)

	FindGroupsLikeName(groupName string) (GroupIterator, error)
	FindGroupsByCondition(condition *common.Condition) (GroupIterator, error)
}

// AttributeIterator the attribute iterator
type AttributeIterator interface {
	Next() (Attribute, error)
}

// Attribute the interface declaration for model attribute maintence
type Attribute interface {
	types.Saver

	SetID(id string)
	SetName(name string)
	SetUnit(unit string)
	SetPlaceholer(placeHoler string)
	SetEditable()
	SetNonEditable()
	Editable() bool
	SetRequired()
	SetNonRequired()
	Required() bool
	SetKey(isKey bool)
	SetOption(option string)
	SetDescrition(des string)
}
