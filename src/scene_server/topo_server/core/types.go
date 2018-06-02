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

package core

import (
	frcommon "configcenter/src/framework/common"
	frtypes "configcenter/src/framework/core/types"

	"configcenter/src/common/errors"
	"configcenter/src/common/language"

	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

// LogicParams the logic function params
type LogicParams struct {
	Err  errors.CCErrorIf
	Lang language.CCLanguageIf
}

// LogicFunc the core logic function definition
type LogicFunc func(params LogicParams, data frtypes.MapStr) (frtypes.MapStr, error)

// Core Provides management interfaces for models and instances
type Core interface {
	CreateClassification(data frtypes.MapStr) (model.Classification, error)
	CreateObject(data frtypes.MapStr) (model.Object, error)
	CreateObjectAttribute(data frtypes.MapStr) (model.Attribute, error)
	CreateObjectGroup(data frtypes.MapStr) (model.Group, error)
	CreateInst(obj model.Object, data frtypes.MapStr) (inst.Inst, error)
	CreateAssociation(data frtypes.MapStr) (model.Association, error)

	DeleteClassification(cond frcommon.Condition) error
	DeleteObject(cond frcommon.Condition) error
	DeleteObjectAttribute(cond frcommon.Condition) error
	DeleteObjectGroup(cond frcommon.Condition) error
	DeleteInst(cond frcommon.Condition) error
	DeleteAssociation(cond frcommon.Condition) error

	FindClassification(cond frcommon.Condition) ([]model.Classification, error)
	FindObject(cond frcommon.Condition) ([]model.Object, error)
	FindObjectAttribute(cond frcommon.Condition) ([]model.Attribute, error)
	FindObjectGroup(cond frcommon.Condition) ([]model.Group, error)
	FindInst(cond frcommon.Condition) ([]inst.Inst, error)

	UpdateClassification(data frtypes.MapStr, cond frcommon.Condition) error
	UpdateObject(data frtypes.MapStr, cond frcommon.Condition) error
	UpdateObjectAttribute(data frtypes.MapStr, cond frcommon.Condition) error
	UpdateObjectGroup(data frtypes.MapStr, cond frcommon.Condition) error
	UpdateInst(data frtypes.MapStr, cond frcommon.Condition) error
	UpdateAssociation(data frtypes.MapStr, cond frcommon.Condition) error
}
