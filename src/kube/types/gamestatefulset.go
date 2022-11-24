/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package types

import (
	"errors"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/kube/orm"
	"configcenter/src/storage/dal/table"
)

// GameStatefulSetFields merge the fields of the GameStatefulSet and the details corresponding to the fields together.
var GameStatefulSetFields = table.MergeFields(CommonSpecFieldsDescriptor, NamespaceBaseRefDescriptor,
	ClusterBaseRefDescriptor, GameStatefulSetSpecFieldsDescriptor)

// GameStatefulSetSpecFieldsDescriptor GameStatefulSet spec's fields descriptors.
var GameStatefulSetSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: LabelsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: SelectorField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: ReplicasField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: StrategyTypeField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: MinReadySecondsField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: RollingUpdateStrategyField, Type: enumor.Object, IsRequired: false, IsEditable: true},
}

// GameStatefulSetUpdateStrategyType is a string enumeration type that enumerates
// all possible update strategies for the StatefulSet controller.
type GameStatefulSetUpdateStrategyType string

const (
	// RollingUpdateGameStatefulSetStrategyType indicates that update will be
	// applied to all Pods in the StatefulSet with respect to the StatefulSet
	// ordering constraints. When a scale operation is performed with this
	// strategy, new Pods will be created from the specification version indicated
	// by the StatefulSet's updateRevision.
	RollingUpdateGameStatefulSetStrategyType = "RollingUpdate"
	// OnDeleteGameStatefulSetStrategyType triggers the legacy behavior. Version
	// tracking and ordered rolling restarts are disabled. Pods are recreated
	// from the StatefulSetSpec when they are manually deleted. When a scale
	// operation is performed with this strategy,specification version indicated
	// by the StatefulSet's currentRevision.
	OnDeleteGameStatefulSetStrategyType = "OnDelete"
	// InplaceUpdateGameStatefulSetStrategyType indicates that update will be
	// applied to all Pods in the StatefulSet with respect to the StatefulSet
	// ordering constraints. When a scale operation is performed with this
	// strategy, new Pods will be created from the specification version indicated
	// by the StatefulSet's updateRevision.
	InplaceUpdateGameStatefulSetStrategyType = "InplaceUpdate"
	// HotPatchGameStatefulSetStrategyType indicates that pods in the GameStatefulSet will be update hot-patch
	HotPatchGameStatefulSetStrategyType = "HotPatchUpdate"
)

// RollingUpdateGameStatefulSetStrategy spec to control the desired behavior of rolling update.
type RollingUpdateGameStatefulSetStrategy struct {
	// Partition indicates the ordinal at which the StatefulSet should be partitioned for updates.
	Partition *int32 `json:"partition" bson:"partition"`

	// The maximum number of pods that can be unavailable during the update.
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	// This can not be 0 if MaxSurge is 0.
	MaxUnavailable *IntOrString `json:"max_unavailable" bson:"max_unavailable"`

	// The maximum number of pods that can be scheduled above the desired number of pods.
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	// This can not be 0 if MaxUnavailable is 0.
	MaxSurge *IntOrString `json:"max_surge" bson:"max_surge"`
}

// GameStatefulSet define the gameStatefulSet struct.
type GameStatefulSet struct {
	WorkloadBase          `json:",inline" bson:",inline"`
	Labels                *map[string]string                    `json:"labels,omitempty" bson:"labels"`
	Selector              *LabelSelector                        `json:"selector,omitempty" bson:"selector"`
	Replicas              *int64                                `json:"replicas,omitempty" bson:"replicas"`
	MinReadySeconds       *int64                                `json:"min_ready_seconds,omitempty" bson:"min_ready_seconds"`
	StrategyType          *GameStatefulSetUpdateStrategyType    `json:"strategy_type,omitempty" bson:"strategy_type"`
	RollingUpdateStrategy *RollingUpdateGameStatefulSetStrategy `json:"rolling_update_strategy,omitempty" bson:"rolling_update_strategy"`
}

// GetWorkloadBase get workload base
func (g *GameStatefulSet) GetWorkloadBase() WorkloadBase {
	return g.WorkloadBase
}

// SetWorkloadBase set workload base
func (g *GameStatefulSet) SetWorkloadBase(wl WorkloadBase) {
	g.WorkloadBase = wl
}

// ValidateCreate validate create workload
func (w *GameStatefulSet) ValidateCreate() ccErr.RawErrorInfo {
	if w == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateCreate(*w, GameStatefulSetFields); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// ValidateUpdate validate update workload
func (w *GameStatefulSet) ValidateUpdate() ccErr.RawErrorInfo {
	if w == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateUpdate(*w, GameStatefulSetFields); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// BuildUpdateData build gameStatefulSet update data
func (w *GameStatefulSet) BuildUpdateData(user string) (map[string]interface{}, error) {
	if w == nil {
		return nil, errors.New("update param is invalid")
	}

	now := time.Now().Unix()
	opts := orm.NewFieldOptions().AddIgnoredFields(wlIgnoreField...)
	updateData, err := orm.GetUpdateFieldsWithOption(w, opts)
	if err != nil {
		return nil, err
	}
	updateData[common.LastTimeField] = now
	updateData[common.ModifierField] = user
	return updateData, err
}
