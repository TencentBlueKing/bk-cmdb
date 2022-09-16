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
	"fmt"
	"reflect"

	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/filter"
	"configcenter/src/storage/dal/table"
)

// PodFields merge the fields of the cluster and the details corresponding to the fields together.
var PodFields = table.MergeFields(PodSpecFieldsDescriptor)

// PodSpecFieldsDescriptor pod spec's fields descriptors.
var PodSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	//	{Field: NamespaceField, Type: enumor.String, IsRequired: true, IsEditable: false},
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

// ContainerFields merge the fields of the cluster and the details corresponding to the fields together.
var ContainerFields = table.MergeFields(ContainerSpecFieldsDescriptor)

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

	if p.Ref != nil {
		err := p.Ref.Kind.Validate()
		if (p.Ref.Name == nil && p.Ref.ID == nil) || p.Ref.Kind == nil || err != nil {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{RefField},
			}
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
	ID            int64 `json:"id,omitempty" bson:"id"`
	SysSpec       `json:",inline" bson:",inline"`
	PodBaseFields `json:",inline" bson:",inline"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// PodBaseFields pod core details
type PodBaseFields struct {
	Name          *string            `json:"name,omitempty" bson:"name"`
	Priority      *int32             `json:"priority,omitempty" bson:"priority"`
	Labels        *map[string]string `json:"labels,omitempty"  bson:"labels"`
	IP            *string            `json:"ip,omitempty"  bson:"ip"`
	IPs           *[]PodIP           `json:"ips,omitempty"  bson:"ips"`
	Volumes       *[]Volume          `json:"volumes,omitempty"  bson:"volumes"`
	QOSClass      *PodQOSClass       `json:"qos_class,omitempty"  bson:"qos_class"`
	NodeSelectors *map[string]string `json:"node_selectors,omitempty"  bson:"node_selectors"`
	Tolerations   *[]Toleration      `json:"tolerations,omitempty" bson:"tolerations"`
}

// Container container details
type Container struct {
	// cc的自增主键
	ID                  int64 `json:"id,omitempty" bson:"id"`
	PodID               int64 `json:"bk_pod_id,omitempty" bson:"bk_pod_id"`
	ContainerBaseFields `json:",inline" bson:",inline"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// ContainerBaseFields container core details
type ContainerBaseFields struct {
	Name            *string          `json:"name,omitempty" bson:"name"`
	ContainerID     *string          `json:"container_uid,omitempty" bson:"container_uid"`
	Image           *string          `json:"image,omitempty" bson:"image"`
	Ports           *[]ContainerPort `json:"ports,omitempty" bson:"ports"`
	HostPorts       *[]ContainerPort `json:"host_ports,omitempty" bson:"host_ports"`
	Args            *[]string        `json:"args,omitempty" bson:"args"`
	Started         *int64           `json:"started,omitempty" bson:"started"`
	Limits          *ResourceList    `json:"limits,omitempty" bson:"limits"`
	ReqSysSpecuests *ResourceList    `json:"requests,omitempty" bson:"requests"`
	Liveness        *Probe           `json:"liveness,omitempty" bson:"liveness"`
	Environment     *[]EnvVar        `json:"environment,omitempty" bson:"environment"`
	Mounts          *[]VolumeMount   `json:"mounts,omitempty" bson:"mounts"`
}

// SysSpec the relationship information related to the container
// that stores the cc, all types share this structure.
type SysSpec struct {
	BizID           *int64  `json:"bk_biz_id,omitempty" bson:"bk_biz_id"`
	SupplierAccount *string `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
	ClusterID       *int64  `json:"bk_cluster_id,omitempty" bson:"bk_cluster_id"`
	// redundant cluster id
	ClusterUID  *string `json:"cluster_uid,omitempty" bson:"cluster_uid"`
	NameSpaceID *int64  `json:"bk_namespace_id,omitempty" bson:"bk_namespace_id"`
	// redundant namespace names
	NameSpace *string `json:"namespace,omitempty" bson:"namespace"`
	Workload  *Ref    `json:"ref,omitempty" bson:"ref"`
	HostID    *int64  `json:"bk_host_id,omitempty" bson:"bk_host_id"`
	NodeID    *int64  `json:"bk_node_id,omitempty" bson:"bk_node_id"`
	// redundant node names
	Node *string `json:"node_name,omitempty" bson:"node_name"`
	// all container related data use the same relation structure, pod
	// does not need these two fields, only container needs these two fields.
	PodID *int64  `json:"bk_pod_id,omitempty" bson:"bk_pod_id ,omitempty"`
	Pod   *string `json:"pod_name,omitempty" bson:"pod_name ,omitempty"`
}

// Ref pod-related workload association information.
type Ref struct {
	Kind string `json:"kind"`
	// redundant workload names
	Name string `json:"name,omitempty"`
	// ID workload ID in cc
	ID int64 `json:"id,omitempty"`
}

// PodsInfo details of creating pods.
type PodsInfo struct {
	KubeSpecInfo  *KubeSpec `json:"kube_spec"`
	CmdbSpecInfo  *CmdbSpec `json:"cmdb_spec"`
	HostID        int64     `json:"bk_host_id"`
	PodBaseFields `json:",inline"`
	Containers    []ContainerBaseFields `json:"containers"`
}

// CreateValidate validate the PodBaseFields
func (option *PodBaseFields) CreateValidate() error {

	if option == nil {
		return errors.New("pod information must be given")
	}

	// first get a list of required fields.
	requireMap := make(map[string]struct{}, 0)
	requires := PodFields.RequiredFields()
	for field, required := range requires {
		if required {
			requireMap[field] = struct{}{}
		}
	}

	typeOfOption := reflect.TypeOf(*option)
	valueOfOption := reflect.ValueOf(*option)
	for i := 0; i < typeOfOption.NumField(); i++ {
		tag := typeOfOption.Field(i).Tag.Get("json")
		if PodFields.IsFieldRequiredByField(tag) {
			fieldValue := valueOfOption.Field(i)
			if fieldValue.IsNil() {
				return fmt.Errorf("required fields cannot be empty, %s", tag)
			}
			delete(requireMap, tag)
		}
	}

	if len(requireMap) > 0 {
		return fmt.Errorf("required fields cannot be empty")
	}
	return nil
}

// CreateValidate validate the ContainerBaseFields
func (option *ContainerBaseFields) CreateValidate() error {

	// first get a list of required fields.
	requireMap := make(map[string]struct{}, 0)
	requires := ContainerFields.RequiredFields()
	for field, required := range requires {
		if required {
			requireMap[field] = struct{}{}
		}
	}

	if option == nil {
		return errors.New("node information must be given")
	}
	typeOfOption := reflect.TypeOf(*option)
	valueOfOption := reflect.ValueOf(*option)
	for i := 0; i < typeOfOption.NumField(); i++ {
		tag := typeOfOption.Field(i).Tag.Get("json")

		if ContainerFields.IsFieldRequiredByField(tag) {
			fieldValue := valueOfOption.Field(i)
			if fieldValue.IsNil() {
				return fmt.Errorf("required fields cannot be empty, %s", tag)
			}
			delete(requireMap, tag)
		}
	}

	if len(requireMap) > 0 {
		return fmt.Errorf("required fields cannot be empty")
	}
	return nil

}

// CreatePodsOption create pods option
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

// Validate validate the CreatePodsOption
func (option *CreatePodsOption) Validate() error {

	if len(option.Pods) == 0 {
		return errors.New("params cannot be empty")
	}

	if len(option.Pods) > createPodsLimit {
		return fmt.Errorf("the maximum number of pods created at one time cannot exceed %d", createPodsLimit)
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

		if err := pod.CreateValidate(); err != nil {
			return err
		}

		for _, container := range pod.Containers {
			if err := container.CreateValidate(); err != nil {
				return err
			}
		}
	}
	return nil
}

// PodInstResp pod instance response
type PodInstResp struct {
	metadata.BaseResp `json:",inline"`
	Data              PodDataResp `json:"data"`
}

// PodDataResp pod data response
type PodDataResp struct {
	Info []Pod `json:"info"`
}

// ContainerInstResp container instance response
type ContainerInstResp struct {
	metadata.BaseResp `json:",inline"`
	Data              ContainerDataResp `json:"data"`
}

// ContainerDataResp container data response
type ContainerDataResp struct {
	Info []Container `json:"info"`
}
