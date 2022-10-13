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
	{Field: PriorityField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: LabelsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: IPField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: IPsField, Type: enumor.Array, IsRequired: false, IsEditable: true},
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
	// podQueryLimit limit on the number of pod query.
	podQueryLimit = 500
	// createPodsLimit the maximum number of pods to be created at one time.
	createPodsLimit = 200
	// containerQueryLimit limit on the number of container query
	containerQueryLimit = 500
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
	if (p.ClusterID != 0 || p.NamespaceID != 0 || p.Ref.ID != 0 || p.NodeID != 0) &&
		(p.ClusterUID != "" || p.Namespace != "" || p.Ref.Name != "" || p.NodeName != "") {

		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrorTopoIdentificationIllegal,
		}
	}

	err := p.Ref.Kind.Validate()
	if (p.Ref.Name != "" || p.Ref.ID != 0) && err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{RefField},
		}
	}

	if err := p.Page.ValidateWithEnableCount(false, podQueryLimit); err.ErrCode != 0 {
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

	if p.ClusterID != 0 {
		cond[BKClusterIDFiled] = p.ClusterID
	}

	if p.ClusterUID != "" {
		cond[ClusterUIDField] = p.ClusterUID
	}

	if p.NamespaceID != 0 {
		cond[BKNamespaceIDField] = p.NamespaceID
	}

	if p.Namespace != "" {
		cond[NamespaceField] = p.Namespace
	}

	if p.Ref.Kind != "" {
		cond[RefKindField] = p.Ref.Kind
	}

	if p.Ref.Name != "" {
		cond[RefNameField] = p.Ref.Name
	}

	if p.Ref.ID != 0 {
		cond[RefIDField] = p.Ref.ID
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
	Name          *string            `json:"name,omitempty" bson:"name"`
	Priority      *int32             `json:"priority,omitempty" bson:"priority"`
	Labels        *map[string]string `json:"labels,omitempty"  bson:"labels"`
	IP            *string            `json:"ip,omitempty"  bson:"ip"`
	IPs           *[]PodIP           `json:"ips,omitempty"  bson:"ips"`
	Volumes       *[]Volume          `json:"volumes,omitempty"  bson:"volumes"`
	QOSClass      *PodQOSClass       `json:"qos_class,omitempty"  bson:"qos_class"`
	NodeSelectors *map[string]string `json:"node_selectors,omitempty"  bson:"node_selectors"`
	Tolerations   *[]Toleration      `json:"tolerations,omitempty" bson:"tolerations"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// createValidate validate the PodBaseFields
func (option *Pod) createValidate() error {

	if option == nil {
		return errors.New("pod information must be set")
	}

	if option.Name == nil || *option.Name == "" {
		return errors.New("pod name must be set")
	}

	// first get a list of required fields.
	requires := PodFields.RequiredFields()
	for _, required := range requires {
		if !required {
			continue
		}
		typeOfOption := reflect.TypeOf(*option)
		valueOfOption := reflect.ValueOf(*option)
		for i := 0; i < typeOfOption.NumField(); i++ {
			tag, flag := getFieldTag(typeOfOption, i)
			if flag {
				continue
			}

			if !PodFields.IsFieldRequiredByField(tag) {
				continue
			}

			if err := isRequiredField(tag, valueOfOption, i); err != nil {
				return err
			}
		}
	}
	return nil
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

// createValidate validate the ContainerBaseFields
func (option *ContainerBaseFields) createValidate() error {

	if option == nil {
		return errors.New("container information must be set")
	}

	if option.Name == nil || *option.Name == "" {
		return errors.New("container name must be set")
	}

	if option.ContainerID == nil || *option.ContainerID == "" {
		return errors.New("container name must be set")
	}

	typeOfOption := reflect.TypeOf(*option)
	valueOfOption := reflect.ValueOf(*option)
	for i := 0; i < typeOfOption.NumField(); i++ {

		tag, flag := getFieldTag(typeOfOption, i)
		if flag {
			continue
		}

		if !ContainerFields.IsFieldRequiredByField(tag) {
			continue
		}

		if err := isRequiredField(tag, valueOfOption, i); err != nil {
			return err
		}
	}

	return nil
}

// SysSpec the relationship information related to the container
// that stores the cc, all types share this structure.
type SysSpec struct {
	BizID           int64  `json:"bk_biz_id" bson:"bk_biz_id"`
	SupplierAccount string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	ClusterID       int64  `json:"bk_cluster_id,omitempty" bson:"bk_cluster_id"`
	// redundant cluster id
	ClusterUID  string `json:"cluster_uid,omitempty" bson:"cluster_uid"`
	NameSpaceID int64  `json:"bk_namespace_id,omitempty" bson:"bk_namespace_id"`
	// redundant namespace names
	NameSpace string `json:"namespace,omitempty" bson:"namespace"`
	Workload  Ref    `json:"ref,omitempty" bson:"ref"`
	HostID    int64  `json:"bk_host_id,omitempty" bson:"bk_host_id"`
	NodeID    int64  `json:"bk_node_id,omitempty" bson:"bk_node_id"`
	// redundant node names
	Node string `json:"node_name,omitempty" bson:"node_name"`
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
	Spec       SpecInfo `json:"spec"`
	HostID     int64    `json:"bk_host_id"`
	Pod        `json:",inline"`
	Containers []ContainerBaseFields `json:"containers"`
}

// CreatePodsOption create pods option
type CreatePodsOption struct {
	Data []PodsInfoArray `json:"data"`
}

// PodsInfoArray create pods option
type PodsInfoArray struct {
	BizID int64      `json:"bk_biz_id"`
	Pods  []PodsInfo `json:"pods"`
}

// Validate validate the CreatePodsOption
func (option *CreatePodsOption) Validate() error {

	if len(option.Data) == 0 {
		return errors.New("params cannot be empty")
	}
	var podsLen int
	for _, data := range option.Data {
		podsLen += len(data.Pods)
	}
	if podsLen > createPodsLimit {
		return fmt.Errorf("the maximum number of pods created at one time cannot exceed %d", createPodsLimit)
	}

	for _, data := range option.Data {
		for _, pod := range data.Pods {
			if err := pod.Spec.validate(); err != nil {
				return err
			}
			if pod.HostID == 0 {
				return errors.New("host id must be set")
			}

			if err := pod.createValidate(); err != nil {
				return err
			}

			for _, container := range pod.Containers {
				if err := container.createValidate(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ContainerQueryReq container query request
type ContainerQueryReq struct {
	PodID  int64              `json:"bk_pod_id"`
	Filter *filter.Expression `json:"filter"`
	Fields []string           `json:"fields,omitempty"`
	Page   metadata.BasePage  `json:"page,omitempty"`
}

// Validate validate ContainerQueryReq
func (p *ContainerQueryReq) Validate() ccErr.RawErrorInfo {
	if p.PodID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{BKPodIDField},
		}
	}

	if err := p.Page.ValidateWithEnableCount(false, containerQueryLimit); err.ErrCode != 0 {
		return err
	}

	// todo validate Filter
	return ccErr.RawErrorInfo{}
}

// BuildCond build query container condition
func (p *ContainerQueryReq) BuildCond() (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		BKPodIDField: p.PodID,
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

// CreatePodsResult create pods results in batches.
type CreatePodsResult struct {
	metadata.BaseResp
	Info []Pod `json:"data" bson:"data"`
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

// DeletePodsOption delete pods option, pods are aggregated by biz id
type DeletePodsOption struct {
	// Data array of delete pod data that defines pods to be deleted in one biz
	Data []DeletePodData `json:"data"`
}

// DeletePodData delete pods data, including biz id and pods in it
type DeletePodData struct {
	// BizID biz id
	BizID int64 `json:"bk_biz_id"`
	// PodIDs pod cc id array
	PodIDs []int64 `json:"ids"`
}

// Validate delete pods option
func (d *DeletePodsOption) Validate() ccErr.RawErrorInfo {
	if len(d.Data) == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"data"}}
	}

	if len(d.Data) > common.BKMaxWriteOpLimit {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommXXExceedLimit, Args: []interface{}{
			"data", common.BKMaxWriteOpLimit}}
	}

	// validate that all delete pods count must not exceed 200
	podsCnt := 0
	for _, data := range d.Data {
		if data.BizID == 0 {
			return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{common.BKAppIDField}}
		}

		if len(data.PodIDs) == 0 {
			return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"ids"}}
		}

		podsCnt += len(data.PodIDs)
		if podsCnt > common.BKMaxWriteOpLimit {
			return ccErr.RawErrorInfo{ErrCode: common.CCErrCommXXExceedLimit, Args: []interface{}{
				"pods", common.BKMaxWriteOpLimit}}
		}
	}

	return ccErr.RawErrorInfo{}
}

// DeletePodsByIDsOption delete pods by ids option
type DeletePodsByIDsOption struct {
	// PodIDs delete pod id array
	PodIDs []int64 `json:"ids"`
}

// Validate delete pods by id option
func (d *DeletePodsByIDsOption) Validate() ccErr.RawErrorInfo {
	if len(d.PodIDs) == 0 {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"ids"}}
	}

	if len(d.PodIDs) > common.BKMaxWriteOpLimit {
		return ccErr.RawErrorInfo{ErrCode: common.CCErrCommXXExceedLimit, Args: []interface{}{
			"ids", common.BKMaxWriteOpLimit}}
	}

	return ccErr.RawErrorInfo{}
}
