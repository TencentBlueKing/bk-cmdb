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

// PodsWorkloadFields merge the fields of the PodsWorkload and the details corresponding to the fields together.
var PodsWorkloadFields = table.MergeFields(CommonSpecFieldsDescriptor, NamespaceBaseRefDescriptor,
	ClusterBaseRefDescriptor, PodsWorkloadSpecFieldsDescriptor)

// PodsWorkloadSpecFieldsDescriptor PodsWorkload spec's fields descriptors.
var PodsWorkloadSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: LabelsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: SelectorField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: ReplicasField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: StrategyTypeField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: MinReadySecondsField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: RollingUpdateStrategyField, Type: enumor.Object, IsRequired: false, IsEditable: true},
}

// PodsWorkload define the pods workload struct.
type PodsWorkload struct {
	WorkloadBase    `json:",inline" bson:",inline"`
	Labels          *map[string]string `json:"labels,omitempty" bson:"labels"`
	Selector        *LabelSelector     `json:"selector,omitempty" bson:"selector"`
	Replicas        *int64             `json:"replicas,omitempty" bson:"replicas"`
	MinReadySeconds *int64             `json:"min_ready_seconds,omitempty" bson:"min_ready_seconds"`
}

// GetWorkloadBase get workload base
func (p *PodsWorkload) GetWorkloadBase() WorkloadBase {
	return p.WorkloadBase
}

// SetWorkloadBase set workload base
func (p *PodsWorkload) SetWorkloadBase(wl WorkloadBase) {
	p.WorkloadBase = wl
}

// ValidateCreate validate create workload
func (w *PodsWorkload) ValidateCreate() ccErr.RawErrorInfo {
	if w == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateCreate(*w, PodsWorkloadFields); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// ValidateUpdate validate update workload
func (w *PodsWorkload) ValidateUpdate() ccErr.RawErrorInfo {
	if w == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateUpdate(*w, PodsWorkloadFields); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// BuildUpdateData build workload pods update data
func (w *PodsWorkload) BuildUpdateData(user string) (map[string]interface{}, error) {
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
