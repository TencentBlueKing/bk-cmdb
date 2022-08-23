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
	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/filter"
	"configcenter/src/storage/dal/table"
)

// PodSpecFieldsDescriptor pod spec's fields descriptors.
var PodSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: NamespaceField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: PriorityField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: LabelsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: IPField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: IPsField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: ControlledBy, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: ContainerUIDField, Type: enumor.Array, IsRequired: false, IsEditable: true},
	{Field: VolumesField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: QOSClassField, Type: enumor.Enum, IsRequired: false, IsEditable: true},
	{Field: NodeSelectorsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: TolerationsField, Type: enumor.Object, IsRequired: false, IsEditable: true},
}

// ContainerSpecFieldsDescriptor container spec's fields descriptors.
var ContainerSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: ContainerUIDField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: ImageField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: PortsField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: HostPortsField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: ArgsField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: StartedField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: RequestsField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: LimitsField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: LivenessField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: EnvironmentField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: MountsField, Type: enumor.Object, IsRequired: false, IsEditable: true},
}

// Container 容器表定义
type Container struct {
	// cc的自增主键
	ID int64 `json:"bk_container_id"`
	// PodID pod id in cc
	PodID *int64 `json:"bk_pod_id" bson:"bk_pod_id"`

	// Pod pod name in third party platform
	Pod     *string `json:"pod" bson:"pod"`
	OwnerID string  `json:"bk_supplier_account"`
	Name    string  `json:"name"`
	// 容器ID
	ContainerID string `json:"container_id"`
	Image       string `json:"image,omitempty"`
	// 确认下这两个端口有什么区别
	//Ports     []v1.ContainerPort `json:"ports,omitempty"`
	//HostPorts []v1.ContainerPort `json:"host_ports,omitempty"`
	Args []string `json:"args,omitempty"`
	// 启动时间，unix时间戳
	Started int64 `json:"started,omitempty"`
	//Limits  v1.ResourceList `json:"limits,omitempty"`
	//Requests    v1.ResourceList  `json:"requests,omitempty"`
	//Liveness    *v1.Probe        `json:"liveness,omitempty"`
	//Environment []v1.EnvVar      `json:"environment,omitempty"`
	//Mounts      []v1.VolumeMount `json:"mounts,omitempty"`
	// cc时间，unix时间戳（或者按之前的用时间类型？）
	LastTime   int64 `json:"last_time"`
	CreateTime int64 `json:"create_time"`
}

// PodQueryReq pod query request
type PodQueryReq struct {
	WorkloadSpec `json:",inline" bson:",inline"`
	HostID       int64              `json:"bk_host_id"`
	NodeID       int64              `json:"bk_node_id"`
	NodeName     string             `json:"node_name"`
	Filter       *filter.Expression `json:"filter"`
	Fields       []string           `json:"fields,omitempty"`
	Page         metadata.BasePage  `json:"page,omitempty"`
}

// Validate validate PodQueryReq
func (p *PodQueryReq) Validate() errors.RawErrorInfo {
	if (p.ClusterID != nil || p.NamespaceID != nil || (p.Ref != nil && p.Ref.ID != nil) || p.NodeID != 0) &&
		(p.ClusterUID != nil || p.Namespace != nil || (p.Ref != nil && p.Ref.Name != nil) || p.NodeName != "") {

		return errors.RawErrorInfo{
			ErrCode: common.CCErrorTopoIdentificationIllegal,
		}
	}

	if p.Ref != nil && ((p.Ref.Name == nil && p.Ref.ID == nil) || p.Ref.Kind == nil ||
		!IsInnerWorkload(WorkloadType(*p.Ref.Kind))) {

		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{RefField},
		}
	}

	if errInfo, err := p.Page.Validate(false); err != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{errInfo},
		}
	}

	// todo validate Filter
	return errors.RawErrorInfo{}
}

// BuildCond build query pod condition
func (p *PodQueryReq) BuildCond(bizID int64, supplierAccount string) (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}
	if supplierAccount != "" {
		cond[common.BkSupplierAccount] = supplierAccount
	}

	if p.ClusterID != nil {
		cond[BKClusterIDFiled] = p.ClusterID
	}

	if p.ClusterUID != nil {
		cond[ClusterUIDField] = p.ClusterUID
	}

	if p.NamespaceID != nil {
		cond[BKNamespaceIDField] = p.NamespaceID
	}

	if p.Namespace != nil {
		cond[NamespaceField] = p.Namespace
	}

	if p.Ref != nil {
		if p.Ref.Kind != nil {
			cond[RefKindField] = p.Ref.Kind
		}

		if p.Ref.Name != nil {
			cond[RefNameField] = p.Ref.Name
		}

		if p.Ref.ID != nil {
			cond[RefIDField] = p.Ref.ID
		}
	}

	if p.HostID != 0 {
		cond[common.BKHostIDField] = p.HostID
	}

	if p.NodeID != 0 {
		cond[BKNodeIDField] = p.NodeID
	}

	if p.NodeName != "" {
		cond[NodeField] = p.NodeName
	}

	if p.Filter != nil {
		filterCond, err := p.Filter.ToMgo()
		if err != nil {
			return nil, err
		}
		cond = mapstr.MapStr{common.BKDBAND: []mapstr.MapStr{cond, filterCond}}
	}
	return cond, nil
}

// ContainerAttrsRsp 容器资源属性回应
type ContainerAttrsRsp struct {
	Field    string `json:"field"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}
