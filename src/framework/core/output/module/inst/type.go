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

package inst

import (
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

// FieldName the field name
type FieldName string

const (
	// BusinessID the business id
	BusinessID = "bk_biz_id"
	// InstID the common inst id
	InstID = "bk_inst_id"

	// InstName the common inst name
	InstName = "bk_inst_name"

	// SetID the set id
	SetID = "bk_set_id"

	// SetName the set name
	SetName = "bk_set_name"

	// ModuleID the module id
	ModuleID = "bk_module_id"

	// ModuleName the module name
	ModuleName = "bk_module_name"
)

// Topo the inst topo structure
type Topo interface {
	Pre() Inst
	Next() Inst
}

// Iterator the iterator interface for the Inst
type Iterator interface {
	Next() (Inst, error)
	ForEach(callbackItem func(item Inst)) error
}

// Inst the inst interface
type Inst interface {
	types.Saver

	GetModel() model.Model

	IsMainLine() bool

	GetAssociationModels() ([]model.Model, error)

	GetInstID() int
	GetInstName() string

	SetValue(key string, value interface{}) error
	GetValues() (types.MapStr, error)

	GetAssociationsByModleID(modleID string) ([]Inst, error)
	GetAllAssociations() (map[model.Model][]Inst, error)

	SetParent(parentInstID int) error
	GetParent() ([]Topo, error)
	GetChildren() ([]Topo, error)
}
