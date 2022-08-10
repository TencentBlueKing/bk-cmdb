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
	"configcenter/src/storage/dal/table"
)

// PodSpecFieldsDescriptor pod spec's fields descriptors.
var PodSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: table.String, IsRequired: true, IsEditable: false},
	{Field: NamespaceField, Type: table.String, IsRequired: true, IsEditable: false},
	{Field: PriorityField, Type: table.Numeric, IsRequired: false, IsEditable: true},
	{Field: LabelsField, Type: table.MapString, IsRequired: false, IsEditable: true},
	{Field: IPField, Type: table.String, IsRequired: false, IsEditable: true},
	{Field: IPsField, Type: table.String, IsRequired: false, IsEditable: true},
	{Field: ControlledBy, Type: table.String, IsRequired: false, IsEditable: true},
	{Field: ContainerUIDField, Type: table.Array, IsRequired: false, IsEditable: true},
	{Field: VolumesField, Type: table.Object, IsRequired: false, IsEditable: true},
	{Field: QOSClassField, Type: table.Enum, IsRequired: false, IsEditable: true},
	{Field: NodeSelectorsField, Type: table.MapString, IsRequired: false, IsEditable: true},
	{Field: TolerationsField, Type: table.Object, IsRequired: false, IsEditable: true},
}

// ContainerSpecFieldsDescriptor container spec's fields descriptors.
var ContainerSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: table.String, IsRequired: true, IsEditable: false},
	{Field: ContainerUIDField, Type: table.String, IsRequired: true, IsEditable: false},
	{Field: ImageField, Type: table.String, IsRequired: true, IsEditable: false},
	{Field: PortsField, Type: table.String, IsRequired: false, IsEditable: true},
	{Field: HostPortsField, Type: table.String, IsRequired: false, IsEditable: true},
	{Field: ArgsField, Type: table.String, IsRequired: false, IsEditable: true},
	{Field: StartedField, Type: table.Numeric, IsRequired: false, IsEditable: true},
	{Field: RequestsField, Type: table.Object, IsRequired: false, IsEditable: true},
	{Field: LimitsField, Type: table.Object, IsRequired: false, IsEditable: true},
	{Field: LivenessField, Type: table.Object, IsRequired: false, IsEditable: true},
	{Field: EnvironmentField, Type: table.Object, IsRequired: false, IsEditable: true},
	{Field: MountsField, Type: table.Object, IsRequired: false, IsEditable: true},
}

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
