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

// GameDeploymentFields merge the fields of the GameDeployment and the details corresponding to the fields together.
var GameDeploymentFields = table.MergeFields(CommonSpecFieldsDescriptor, NamespaceBaseRefDescriptor,
	ClusterBaseRefDescriptor, GameDeploymentSpecFieldsDescriptor)

// GameDeploymentSpecFieldsDescriptor GameDeployment spec's fields descriptors.
var GameDeploymentSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: LabelsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: SelectorField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: ReplicasField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: StrategyTypeField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: MinReadySecondsField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: RollingUpdateStrategyField, Type: enumor.Object, IsRequired: false, IsEditable: true},
}

// GameDeploymentUpdateStrategyType defines strategies for pods in-place update.
type GameDeploymentUpdateStrategyType string

const (
	// RollingGameDeploymentUpdateStrategyType indicates that we always delete Pod and create new Pod
	// during Pod update, which is the default behavior.
	RollingGameDeploymentUpdateStrategyType GameDeploymentUpdateStrategyType = "RollingUpdate"

	// InPlaceGameDeploymentUpdateStrategyType indicates that we will in-place update Pod instead of
	// recreating pod. Currently we only allow image update for pod spec. Any other changes to the pod spec will be
	// rejected by kube-apiserver
	InPlaceGameDeploymentUpdateStrategyType GameDeploymentUpdateStrategyType = "InplaceUpdate"

	// HotPatchGameDeploymentUpdateStrategyType indicates that we will hot patch container image with pod being active.
	// Currently we only allow image update for pod spec. Any other changes to the pod spec will be
	// rejected by kube-apiserver
	HotPatchGameDeploymentUpdateStrategyType GameDeploymentUpdateStrategyType = "HotPatchUpdate"
)

// RollingUpdateGameDeployment gameDeployment update strategy
type RollingUpdateGameDeployment struct {
	// Partition is the desired number of pods in old revisions. It means when partition
	// is set during pods updating, (replicas - partition) number of pods will be updated.
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

// GameDeployment define the gameDeployment struct.
type GameDeployment struct {
	WorkloadBase          `json:",inline" bson:",inline"`
	Labels                *map[string]string                `json:"labels,omitempty" bson:"labels"`
	Selector              *LabelSelector                    `json:"selector,omitempty" bson:"selector"`
	Replicas              *int64                            `json:"replicas,omitempty" bson:"replicas"`
	MinReadySeconds       *int64                            `json:"min_ready_seconds,omitempty" bson:"min_ready_seconds"`
	StrategyType          *GameDeploymentUpdateStrategyType `json:"strategy_type,omitempty" bson:"strategy_type"`
	RollingUpdateStrategy *RollingUpdateGameDeployment      `json:"rolling_update_strategy,omitempty" bson:"rolling_update_strategy"`
}

// GetWorkloadBase get workload base
func (d *GameDeployment) GetWorkloadBase() WorkloadBase {
	return d.WorkloadBase
}

// SetWorkloadBase set workload base
func (g *GameDeployment) SetWorkloadBase(wl WorkloadBase) {
	g.WorkloadBase = wl
}

// ValidateCreate validate create workload
func (w *GameDeployment) ValidateCreate() ccErr.RawErrorInfo {
	if w == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateCreate(*w, GameDeploymentFields); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// ValidateUpdate validate update workload
func (w *GameDeployment) ValidateUpdate() ccErr.RawErrorInfo {
	if w == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateUpdate(*w, GameDeploymentFields); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// BuildUpdateData build gameDeployment update data
func (w *GameDeployment) BuildUpdateData(user string) (map[string]interface{}, error) {
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
