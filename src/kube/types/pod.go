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

	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	ccErr "configcenter/src/common/errors"
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

const (
	// PodQueryLimit limit on the number of pod query.
	PodQueryLimit = 500
	// createPodsLimit the maximum number of pods to be created at one time.
	createPodsLimit = 200
)

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
func (p *PodQueryReq) Validate() ccErr.RawErrorInfo {
	if (p.ClusterID != nil || p.NamespaceID != nil || (p.Ref != nil && p.Ref.ID != nil) || p.NodeID != 0) &&
		(p.ClusterUID != nil || p.Namespace != nil || (p.Ref != nil && p.Ref.Name != nil) || p.NodeName != "") {

		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrorTopoIdentificationIllegal,
		}
	}

	if p.Ref != nil && ((p.Ref.Name == nil && p.Ref.ID == nil) || p.Ref.Kind == nil ||
		!IsInnerWorkload(WorkloadType(*p.Ref.Kind))) {

		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{RefField},
		}
	}

	if err := p.Page.ValidateWithEnableCount(false, PodQueryLimit); err.ErrCode != 0 {
		return err
	}

	// todo validate Filter
	return ccErr.RawErrorInfo{}
}

// BuildCond build query pod condition
func (p *PodQueryReq) BuildCond(bizID int64, supplierAccount string) (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKAppIDField:      bizID,
		common.BkSupplierAccount: supplierAccount,
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

// Pod pod details
type Pod struct {
	// cc的自增主键
	ID          int64 `json:"id"`
	SysSpec     `json:",inline" bson:",inline"`
	PodCoreInfo `json:",inline" bson:",inline"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// PodCoreInfo pod core details
type PodCoreInfo struct {
	Name          *string           `json:"name"`
	Priority      *int32            `json:"priority,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	IP            *string           `json:"ip,omitempty"`
	IPs           []PodIP           `json:"ips,omitempty"`
	Volumes       []Volume          `json:"volumes,omitempty"`
	QOSClass      PodQOSClass       `json:"qos_class,omitempty"`
	NodeSelectors map[string]string `json:"node_selectors,omitempty"`
	Tolerations   []Toleration      `json:"tolerations,omitempty"`
}

// Container container details
type Container struct {
	// cc的自增主键
	ID                int64 `json:"id"`
	PodID             int64 `json:"bk_pod_id"`
	SysSpec           `json:",inline"`
	ContainerCoreInfo `json:",inline"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// ContainerCoreInfo container core details
type ContainerCoreInfo struct {
	Name            string          `json:"name"`
	ContainerID     string          `json:"container_uid"`
	Image           string          `json:"image,omitempty"`
	Ports           []ContainerPort `json:"ports,omitempty"`
	HostPorts       []ContainerPort `json:"host_ports,omitempty"`
	Args            []string        `json:"args,omitempty"`
	Started         int64           `json:"started,omitempty"`
	Limits          ResourceList    `json:"limits,omitempty"`
	ReqSysSpecuests ResourceList    `json:"requests,omitempty"`
	Liveness        *Probe          `json:"liveness,omitempty"`
	Environment     []EnvVar        `json:"environment,omitempty"`
	Mounts          []VolumeMount   `json:"mounts,omitempty"`
}

// SysSpec 存放cc的容器相关的关系信息，所有类型共用这个结构体
type SysSpec struct {
	BizID           *int64  `json:"bk_biz_id"`
	SupplierAccount *string `json:"bk_supplier_account"`
	ClusterID       *int64  `json:"bk_cluster_id,omitempty"`
	// 冗余的cluster id
	ClusterUID  *string `json:"cluster_uid,omitempty"`
	NameSpaceID *int64  `json:"bk_namespace_id,omitempty"`
	// 冗余的namespace名称
	NameSpace *string `json:"namespace,omitempty"`
	Workload  *Ref    `json:"workload,omitempty"`
	HostID    *int64  `json:"bk_host_id,omitempty"`
	NodeID    *int64  `json:"bk_node_id,omitempty"`
	// 冗余的node名称
	Node *string `json:"node,omitempty"`
	// 所有容器相关数据用相同的relation结构体，pod不需要这两个字段，仅container需要这两个字段
	PodID *int64  `json:"bk_pod_id,omitempty"`
	Pod   *string `json:"pod_name,omitempty"`
}

// Ref pod-related workload association information.
type Ref struct {
	Kind string `json:"kind"`
	// 冗余的workload名称
	Name string `json:"name,omitempty"`
	// ID workload在cc中的ID
	ID int64 `json:"id,omitempty"`
}

// PodsInfo details of creating pods.
type PodsInfo struct {
	KubeSpecInfo *KubeSpec `json:"kube_spec"`
	CmdbSpecInfo *CmdbSpec `json:"cmdb_spec"`
	HostID       *int64    `json:"bk_host_id"`
	PodCoreInfo  `json:",inline"`
	Containers   []ContainerCoreInfo `json:"containers"`
}

// Validate validate the PodCoreInfo
func (option *PodCoreInfo) Validate() error {

	return nil
}

// Validate validate the ContainerCoreInfo
func (option *ContainerCoreInfo) Validate() error {

	return nil
}

// CreatePodsOption create Pods request
type CreatePodsOption struct {
	Pods []PodsInfo `json:"pods"`
}

// Validate validate the KubeSpec
func (option *KubeSpec) Validate() error {

	if option.ClusterUID == nil {
		return errors.New("cluster uid must be set")
	}
	if option.Namespace == nil {
		return errors.New("namespace must be set")
	}
	if option.Node == nil {
		return errors.New("node must be set")
	}
	if option.WorkloadKind == nil {
		return errors.New("workload kind must be set")
	}
	if option.WorkloadName == nil {
		return errors.New("workload name must be set")
	}
	return nil
}

// Validate validate the CmdbSpec
func (option *CmdbSpec) Validate() error {

	if option.ClusterID == nil {
		return errors.New("cluster id must be set")
	}
	if option.NamespaceID == nil {
		return errors.New("namespace id must be set")
	}
	if option.NodeID == nil {
		return errors.New("node id must be set")
	}
	if option.WorkloadKind == nil {
		return errors.New("workload kind must be set")
	}
	if option.WorkloadID == nil {
		return errors.New("workload id must be set")
	}

	return nil
}

// Validate validate the CreatePodsReq
func (option *CreatePodsOption) Validate() error {

	if len(option.Pods) == 0 {
		return errors.New("param cannot be empty")
	}

	if len(option.Pods) > createPodsLimit {
		return errors.New("param cannot be empty")
	}

	for _, pod := range option.Pods {
		if pod.KubeSpecInfo == nil && pod.CmdbSpecInfo == nil {
			return errors.New("kube spec and cmdb spec cannot be empty at the same time")
		}
		if pod.KubeSpecInfo != nil && pod.CmdbSpecInfo != nil {
			return errors.New("kube spec and cmdb spec cannot be set at the same time")
		}

		if pod.CmdbSpecInfo != nil {
			if err := pod.CmdbSpecInfo.Validate(); err != nil {
				return err
			}
		}

		if pod.KubeSpecInfo != nil {
			if err := pod.KubeSpecInfo.Validate(); err != nil {
				return err
			}
		}

		if err := pod.Validate(); err != nil {
			return err
		}

		for _, container := range pod.Containers {
			if err := container.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}
