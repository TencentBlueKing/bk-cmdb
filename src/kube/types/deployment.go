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

// DeploymentFields merge the fields of the Deployment and the details corresponding to the fields together.
var DeploymentFields = table.MergeFields(CommonSpecFieldsDescriptor, WorkLoadBaseFieldsDescriptor,
	DeploymentSpecFieldsDescriptor)

// DeploymentSpecFieldsDescriptor Deployment spec's fields descriptors.
var DeploymentSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: LabelsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: SelectorField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: ReplicasField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: StrategyTypeField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: MinReadySecondsField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: RollingUpdateStrategyField, Type: enumor.Object, IsRequired: false, IsEditable: true},
}

// DeploymentStrategyType deployment strategy type
type DeploymentStrategyType string

const (
	// RecreateDeploymentStrategyType kill all existing pods before creating new ones.
	RecreateDeploymentStrategyType DeploymentStrategyType = "Recreate"

	// RollingUpdateDeploymentStrategyType replace the old ReplicaSets by new one using rolling update
	// i.e gradually scale down the old ReplicaSets and scale up the new one.
	RollingUpdateDeploymentStrategyType DeploymentStrategyType = "RollingUpdate"
)

// RollingUpdateDeployment spec to control the desired behavior of rolling update.
type RollingUpdateDeployment struct {
	// The maximum number of pods that can be unavailable during the update.
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	// This can not be 0 if MaxSurge is 0.
	MaxUnavailable *IntOrString `json:"max_unavailable" bson:"max_unavailable"`

	// The maximum number of pods that can be scheduled above the desired number of pods.
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	// This can not be 0 if MaxUnavailable is 0.
	MaxSurge *IntOrString `json:"max_surge" bson:"max_surge"`
}

// Deployment define the deployment struct.
type Deployment struct {
	WorkloadBase          `json:",inline" bson:",inline"`
	Labels                *map[string]string       `json:"labels,omitempty" bson:"labels"`
	Selector              *LabelSelector           `json:"selector,omitempty" bson:"selector"`
	Replicas              *int64                   `json:"replicas,omitempty" bson:"replicas"`
	MinReadySeconds       *int64                   `json:"min_ready_seconds,omitempty" bson:"min_ready_seconds"`
	StrategyType          *DeploymentStrategyType  `json:"strategy_type,omitempty" bson:"strategy_type"`
	RollingUpdateStrategy *RollingUpdateDeployment `json:"rolling_update_strategy,omitempty" bson:"rolling_update_strategy"`
}

// GetWorkloadBase get workload base
func (d *Deployment) GetWorkloadBase() WorkloadBase {
	return d.WorkloadBase
}

// SetWorkloadBase set workload base
func (d *Deployment) SetWorkloadBase(wl WorkloadBase) {
	d.WorkloadBase = wl
}

// ValidateCreate validate create workload
func (w *Deployment) ValidateCreate() ccErr.RawErrorInfo {
	if w == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommHTTPInputInvalid,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateCreate(*w, DeploymentFields); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// ValidateUpdate validate update workload
func (w *Deployment) ValidateUpdate() ccErr.RawErrorInfo {
	if w == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommHTTPInputInvalid,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateUpdate(*w, DeploymentFields); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// BuildUpdateData build deployment update data
func (w *Deployment) BuildUpdateData(user string) (map[string]interface{}, error) {
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
