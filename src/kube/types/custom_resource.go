/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

// CustomResourceFields merge the fields of the CustomResource and the details corresponding to the fields together.
var CustomResourceFields = table.MergeFields(CommonSpecFieldsDescriptor, NamespaceBaseRefDescriptor,
	ClusterBaseRefDescriptor, CustomResourceSpecFieldsDescriptor)

// CustomResourceSpecFieldsDescriptor CustomResource spec's fields descriptors.
var CustomResourceSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: CRKindField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: CRApiVersionField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: LabelsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: SelectorField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: ReplicasField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: StrategyTypeField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: MinReadySecondsField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: RollingUpdateStrategyField, Type: enumor.Object, IsRequired: false, IsEditable: true},
}

// CustomResource define the custom resource workload struct.
type CustomResource struct {
	WorkloadBase    `json:",inline" bson:",inline"`
	CRKind          *string            `json:"cr_kind,omitempty" bson:"cr_kind"`
	CRApiVersion    *string            `json:"cr_api_version,omitempty" bson:"cr_api_version"`
	Labels          *map[string]string `json:"labels,omitempty" bson:"labels"`
	Selector        *LabelSelector     `json:"selector,omitempty" bson:"selector"`
	Replicas        *int64             `json:"replicas,omitempty" bson:"replicas"`
	MinReadySeconds *int64             `json:"min_ready_seconds,omitempty" bson:"min_ready_seconds"`
}

// GetWorkloadBase get workload base
func (cr *CustomResource) GetWorkloadBase() WorkloadBase {
	return cr.WorkloadBase
}

// SetWorkloadBase set workload base
func (cr *CustomResource) SetWorkloadBase(wl WorkloadBase) {
	cr.WorkloadBase = wl
}

// ValidateCreate validate create workload
func (cr *CustomResource) ValidateCreate() ccErr.RawErrorInfo {
	if cr == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"data"},
		}
	}

	if cr.BizID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKAppIDField},
		}
	}

	if err := ValidateCreate(*cr, CustomResourceFields); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// ValidateUpdate validate update workload
func (cr *CustomResource) ValidateUpdate() ccErr.RawErrorInfo {
	if cr == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateUpdate(*cr, CustomResourceFields); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// BuildUpdateData build custom resource update data
func (cr *CustomResource) BuildUpdateData(user string) (map[string]interface{}, error) {
	if cr == nil {
		return nil, errors.New("update param is invalid")
	}

	now := time.Now().Unix()
	opts := orm.NewFieldOptions().AddIgnoredFields(append(wlIgnoreField, CRKindField, CRApiVersionField)...)
	updateData, err := orm.GetUpdateFieldsWithOption(cr, opts)
	if err != nil {
		return nil, err
	}
	updateData[common.LastTimeField] = now
	updateData[common.ModifierField] = user
	return updateData, err
}
