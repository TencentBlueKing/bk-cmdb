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
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"

	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// Core Provides management interfaces for models and instances
type Core interface {
	CreateClassification(params types.LogicParams, data frtypes.MapStr) (model.Classification, error)
	CreateObject(params types.LogicParams, data frtypes.MapStr) (model.Object, error)
	CreateObjectAttribute(params types.LogicParams, data frtypes.MapStr) (model.Attribute, error)
	CreateObjectGroup(params types.LogicParams, data frtypes.MapStr) (model.Group, error)
	CreateInst(params types.LogicParams, obj model.Object, data frtypes.MapStr) (inst.Inst, error)
	CreateAssociation(params types.LogicParams, data frtypes.MapStr) (model.Association, error)

	DeleteClassification(params types.LogicParams, cond condition.Condition) error
	DeleteObject(params types.LogicParams, cond condition.Condition) error
	DeleteObjectAttribute(params types.LogicParams, cond condition.Condition) error
	DeleteObjectGroup(params types.LogicParams, cond condition.Condition) error
	DeleteInst(params types.LogicParams, cond condition.Condition) error
	DeleteAssociation(params types.LogicParams, cond condition.Condition) error

	FindClassification(params types.LogicParams, cond condition.Condition) ([]model.Classification, error)
	FindObject(params types.LogicParams, cond condition.Condition) ([]model.Object, error)
	FindObjectAttribute(params types.LogicParams, cond condition.Condition) ([]model.Attribute, error)
	FindObjectGroup(params types.LogicParams, cond condition.Condition) ([]model.Group, error)
	FindInst(params types.LogicParams, cond condition.Condition) ([]inst.Inst, error)

	UpdateClassification(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error
	UpdateObject(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error
	UpdateObjectAttribute(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error
	UpdateObjectGroup(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error
	UpdateInst(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error
	UpdateAssociation(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error
}
