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

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/table"
)

// PodFields merge the fields of the cluster and the details corresponding to the fields together.
var PodFields = table.MergeFields(CommonSpecFieldsDescriptor, BizIDDescriptor, HostIDDescriptor,
	ClusterBaseRefDescriptor, NodeBaseRefDescriptor, NamespaceBaseRefDescriptor,
	WorkLoadRefDescriptor, PodSpecFieldsDescriptor)

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
	{Field: OperatorField, Type: enumor.Array, IsRequired: true, IsEditable: true},
}

// PodBaseRefDescriptor the description used when other resources refer to the pod.
var PodBaseRefDescriptor = table.FieldsDescriptors{
	{Field: BKPodIDField, Type: enumor.Numeric, IsRequired: true, IsEditable: false},
}

// KubeTopoIDDescriptor is the fields descriptors for kube topo resource ids.
var KubeTopoIDDescriptor = table.FieldsDescriptors{
	{Field: BKClusterIDFiled, Type: enumor.Numeric, IsRequired: false, IsEditable: false},
	{Field: BKNamespaceIDField, Type: enumor.Numeric, IsRequired: false, IsEditable: false},
	{Field: RefField, Type: enumor.Object, IsRequired: false, IsEditable: false},
	{Field: RefKindField, Type: enumor.String, IsRequired: false, IsEditable: false},
	{Field: RefIDField, Type: enumor.Numeric, IsRequired: false, IsEditable: false},
}

// ContainerFields merge the fields of the cluster and the details corresponding to the fields together.
var ContainerFields = table.MergeFields(CommonSpecFieldsDescriptor, PodBaseRefDescriptor, BizIDDescriptor,
	KubeTopoIDDescriptor, ContainerSpecFieldsDescriptor)

// ContainerSpecFieldsDescriptor container spec's fields descriptors.
var ContainerSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: ContainerUIDField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: ImageField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: PortsField, Type: enumor.Array, IsRequired: false, IsEditable: true},
	{Field: HostPortsField, Type: enumor.Array, IsRequired: false, IsEditable: true},
	{Field: ArgsField, Type: enumor.Array, IsRequired: false, IsEditable: true},
	{Field: StartedField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: RequestsField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: LimitsField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: LivenessField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: EnvironmentField, Type: enumor.Array, IsRequired: false, IsEditable: true},
	{Field: MountsField, Type: enumor.Array, IsRequired: false, IsEditable: true},
}

const (
	// podQueryLimit limit on the number of pod query.
	podQueryLimit = 500
	// createPodsLimit the maximum number of pods to be created at one time.
	createPodsLimit = 200
	// containerQueryLimit limit on the number of container query
	containerQueryLimit = 500
)

// PodQueryOption pod query request
type PodQueryOption struct {
	BizID  int64              `json:"bk_biz_id"`
	Filter *filter.Expression `json:"filter"`
	Fields []string           `json:"fields,omitempty"`
	Page   metadata.BasePage  `json:"page,omitempty"`
}

// Validate validate PodQueryOption
func (p *PodQueryOption) Validate() ccErr.RawErrorInfo {
	if p.BizID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKAppIDField},
		}
	}

	if err := p.Page.ValidateWithEnableCount(false, podQueryLimit); err.ErrCode != 0 {
		return err
	}

	if p.Filter == nil {
		return ccErr.RawErrorInfo{}
	}

	op := filter.NewDefaultExprOpt(PodFields.FieldsType())
	op.MaxRulesDepth = 4
	if err := p.Filter.Validate(op); err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{err.Error()},
		}
	}
	return ccErr.RawErrorInfo{}
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
	Operator      *[]string          `json:"operator,omitempty" bson:"operator"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// createValidate validate the PodBaseFields
func (option *Pod) createValidate() ccErr.RawErrorInfo {

	if option == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"pod"},
		}
	}

	if option.Name == nil || *option.Name == "" {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"pod name"},
		}
	}

	if option.Operator == nil || len(*option.Operator) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"pod operator"},
		}
	}

	if err := ValidateCreate(*option, PodFields); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// Container container details
type Container struct {
	// cc的自增主键
	ID              int64            `json:"id,omitempty" bson:"id"`
	PodID           int64            `json:"bk_pod_id,omitempty" bson:"bk_pod_id"`
	BizID           int64            `json:"bk_biz_id,omitempty" bson:"bk_biz_id"`
	ClusterID       int64            `json:"bk_cluster_id,omitempty" bson:"bk_cluster_id"`
	NamespaceID     int64            `json:"bk_namespace_id,omitempty" bson:"bk_namespace_id"`
	Ref             *Reference       `json:"ref,omitempty" bson:"ref"`
	SupplierAccount string           `json:"bk_supplier_account" bson:"bk_supplier_account"`
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
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// validateCreate validate the ContainerBaseFields
func (option *Container) validateCreate() ccErr.RawErrorInfo {

	if option == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"container"},
		}
	}

	if option.Name == nil || *option.Name == "" {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"container name"},
		}
	}

	if option.ContainerID == nil || *option.ContainerID == "" {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"container_uid"},
		}
	}

	if err := ValidateCreate(*option, ContainerFields); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// SysSpec the relationship information related to the container
// that stores the cc, all types share this structure.
type SysSpec struct {
	SupplierAccount string `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
	WorkloadSpec    `json:",inline" bson:",inline"`
	HostID          int64 `json:"bk_host_id,omitempty" bson:"bk_host_id"`
	NodeID          int64 `json:"bk_node_id,omitempty" bson:"bk_node_id"`
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
	Spec       SpecSimpleInfo `json:"spec"`
	HostID     int64          `json:"bk_host_id"`
	Pod        `json:",inline"`
	Containers []Container `json:"containers"`
}

// CreatePodsRsp the response message
// body of the created pod result to the user.
type CreatePodsRsp struct {
	metadata.BaseResp
	Data metadata.RspIDs `json:"data"`
}

// PodsInfoArray create pods option
type PodsInfoArray struct {
	BizID int64      `json:"bk_biz_id"`
	Pods  []PodsInfo `json:"pods"`
}

// CreatePodsOption create pods option
type CreatePodsOption struct {
	Data []PodsInfoArray `json:"data"`
}

// Validate validate the CreatePodsOption
func (option *CreatePodsOption) Validate() ccErr.RawErrorInfo {
	if len(option.Data) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{errors.New("data")},
		}
	}
	var podsLen int
	for _, data := range option.Data {
		podsLen += len(data.Pods)
	}
	if podsLen > createPodsLimit {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"pods", createPodsLimit},
		}
	}

	for _, data := range option.Data {
		for _, pod := range data.Pods {
			if err := pod.Spec.validate(); err != nil {
				return ccErr.RawErrorInfo{
					ErrCode: common.CCErrCommParamsIsInvalid,
					Args:    []interface{}{err.Error()},
				}
			}
			if pod.HostID == 0 {
				return ccErr.RawErrorInfo{
					ErrCode: common.CCErrCommParamsNeedSet,
					Args:    []interface{}{errors.New("host id")},
				}
			}

			if err := pod.createValidate(); err.ErrCode != 0 {
				return err
			}

			for _, container := range pod.Containers {
				if err := container.validateCreate(); err.ErrCode != 0 {
					return err
				}
			}
		}
	}
	return ccErr.RawErrorInfo{}
}

// ContainerQueryOption container query request
type ContainerQueryOption struct {
	BizID  int64              `json:"bk_biz_id"`
	PodID  int64              `json:"bk_pod_id"`
	Filter *filter.Expression `json:"filter"`
	Fields []string           `json:"fields,omitempty"`
	Page   metadata.BasePage  `json:"page,omitempty"`
}

// Validate validate ContainerQueryOption
func (p *ContainerQueryOption) Validate() ccErr.RawErrorInfo {
	if p.BizID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKAppIDField},
		}
	}

	if p.PodID < 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{BKPodIDField},
		}
	}

	if err := p.Page.ValidateWithEnableCount(false, containerQueryLimit); err.ErrCode != 0 {
		return err
	}

	if p.Filter == nil {
		return ccErr.RawErrorInfo{}
	}

	op := filter.NewDefaultExprOpt(ContainerFields.FieldsType())
	if err := p.Filter.Validate(op); err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{err.Error()},
		}
	}
	return ccErr.RawErrorInfo{}
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

// GetContainerByTopoOption query container by topo request
type GetContainerByTopoOption struct {
	BizID           int64              `json:"bk_biz_id"`
	Nodes           []NodeMsg          `json:"bk_kube_nodes"`
	PodFilter       *filter.Expression `json:"pod_filter"`
	ContainerFilter *filter.Expression `json:"container_filter"`
	PodFields       []string           `json:"pod_fields"`
	ContainerFields []string           `json:"container_fields"`
	Page            metadata.BasePage  `json:"page"`
}

// NodeMsg kube node message
type NodeMsg struct {
	Kind string `json:"kind"`
	ID   int64  `json:"id"`
}

const arrLimit = 200

// Validate validate GetContainerByTopoOption
func (p *GetContainerByTopoOption) Validate() ccErr.RawErrorInfo {
	if p.BizID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKAppIDField},
		}
	}

	if len(p.Nodes) > arrLimit {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"bk_kube_nodes", arrLimit},
		}
	}

	for _, nodeMsg := range p.Nodes {
		if !IsKubeResourceKind(nodeMsg.Kind) {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{"non-kube objects", nodeMsg.Kind},
			}
		}
	}

	if p.PodFilter != nil {
		op := filter.NewDefaultExprOpt(PodFields.FieldsType())
		op.MaxRulesDepth = 4
		if err := p.PodFilter.Validate(op); err != nil {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{err.Error()},
			}
		}
	}

	if p.ContainerFilter != nil {
		op := filter.NewDefaultExprOpt(ContainerFields.FieldsType())
		if err := p.ContainerFilter.Validate(op); err != nil {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{err.Error()},
			}
		}
	}

	if len(p.PodFields) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"pod_fields"},
		}
	}

	if len(p.ContainerFields) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"container_fields"},
		}
	}

	if err := p.Page.ValidateWithEnableCount(false, containerQueryLimit); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// ParseCond parse pod and container condition
func (p *GetContainerByTopoOption) ParseCond() (map[string]interface{}, map[string]interface{}, error) {
	// if pod filter is not set, add node condition to container condition and use this condition to filter container
	if p.PodFilter == nil {
		containerCond, err := p.MergeKubeTopoFilterRules(p.ContainerFilter)
		if err != nil {
			return nil, nil, err
		}
		return nil, containerCond, nil
	}

	// container has no filter while pod has filter, use pod condition to filter container
	if p.ContainerFilter == nil {
		podCond, err := p.MergeKubeTopoFilterRules(p.PodFilter)
		if err != nil {
			return nil, nil, err
		}
		return podCond, nil, nil
	}

	// container has filter, change pod cond to the cond of the inner object "pod" of the container aggregation result
	podCond, err := filtertools.GenAtomFilter("pod", filter.Object, p.PodFilter).ToMgo()
	if err != nil {
		return nil, nil, err
	}

	containerCond, err := p.MergeKubeTopoFilterRules(p.ContainerFilter)
	if err != nil {
		return nil, nil, err
	}
	return podCond, containerCond, nil
}

// MergeKubeTopoFilterRules merge kube topo filter rules to condition
func (p *GetContainerByTopoOption) MergeKubeTopoFilterRules(filterExpr *filter.Expression) (map[string]interface{},
	error) {

	rules := []filter.RuleFactory{filtertools.GenAtomFilter(common.BKAppIDField, filter.Equal, p.BizID)}

	nodeMap := make(map[string][]int64)
	for _, node := range p.Nodes {
		nodeMap[node.Kind] = append(nodeMap[node.Kind], node.ID)
	}

	nodeRules := make([]filter.RuleFactory, 0)
	for kind, ids := range nodeMap {
		switch kind {
		case KubeCluster:
			nodeRules = append(nodeRules, filtertools.GenAtomFilter(BKClusterIDFiled, filter.In, ids))
		case KubeNamespace:
			nodeRules = append(nodeRules, filtertools.GenAtomFilter(BKNamespaceIDField, filter.In, ids))
		default:
			kindRule := filtertools.GenAtomFilter(KindField, filter.Equal, kind)
			idsRule := filtertools.GenAtomFilter(common.BKFieldID, filter.In, ids)
			andRule, err := filtertools.And(kindRule, idsRule)
			if err != nil {
				return nil, err
			}
			nodeRules = append(nodeRules, filtertools.GenAtomFilter(RefField, filter.Object, andRule))
		}
	}

	if len(nodeRules) != 0 {
		rule := &filter.Expression{RuleFactory: &filter.CombinedRule{Condition: filter.Or, Rules: nodeRules}}
		rules = append(rules, rule)
	}

	if filterExpr != nil {
		rules = append(rules, filterExpr)
	}

	andCond, err := filtertools.And(rules...)
	if err != nil {
		return nil, err
	}

	cond, err := andCond.ToMgo()
	if err != nil {
		return nil, err
	}

	return cond, nil
}

// ContainerWithTopo container with topo message
type ContainerWithTopo struct {
	Container mapstr.MapStr `json:"container"`
	Pod       mapstr.MapStr `json:"pod"`
	Topo      Topo          `json:"topo"`
}

// Topo container topo message
type Topo struct {
	BizID        int64        `json:"bk_biz_id"`
	ClusterID    int64        `json:"bk_cluster_id"`
	NamespaceID  int64        `json:"bk_namespace_id"`
	WorkloadID   int64        `json:"bk_workload_id"`
	WorkloadType WorkloadType `json:"workload_type"`
	HostID       int64        `json:"bk_host_id"`
}

// GetContainerByPodOption get container by pod option
type GetContainerByPodOption struct {
	PodCond       map[string]interface{} `json:"pod_cond"`
	ContainerCond map[string]interface{} `json:"container_cond"`
	Fields        []string               `json:"fields"`
	Page          metadata.BasePage      `json:"page"`
}

// GetContainerByPodResp get container by pod response
type GetContainerByPodResp struct {
	Info  []mapstr.MapStr `json:"info"`
	Count int64           `json:"count"`
}
