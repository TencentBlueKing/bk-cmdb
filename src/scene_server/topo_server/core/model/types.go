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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	metadata "configcenter/src/common/metadata"
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
	Update(data mapstr.MapStr) error
	Save(data mapstr.MapStr) error
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
	Parse(data mapstr.MapStr) (*metadata.Association, error)

	GetType() AssociationType
	SetTopo(parent, child Object) error
	GetTopo(obj Object) (Topo, error)
	ToMapStr() (mapstr.MapStr, error)
}

// Factory used to create object  classification attribute etd.
type Factory interface {
	CreateObject(kit *rest.Kit) Object
	CreateClassification(kit *rest.Kit) Classification
	CreateAttribute(kit *rest.Kit) AttributeInterface
	CreateGroup(kit *rest.Kit, bizID int64) GroupInterface
	CreateCommonAssociation(kit *rest.Kit, obj Object, asstKey string, asstObj Object) Association
	CreateMainLineAssociation(kit *rest.Kit, obj Object, asstKey string, asstObj Object) Association
}

type ObjectAssoPair struct {
	Object      Object
	Association metadata.Association
}
