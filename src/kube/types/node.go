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
	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/table"
)

const (
	maxDeleteNodeNum = 100
	maxCreateNodeNum = 100
	maxUpdateNodeNum = 100
)

// NodeFields merge the fields of the cluster and the details corresponding to the fields together.
var NodeFields = table.MergeFields(CommonSpecFieldsDescriptor, BizIDDescriptor, HostIDDescriptor,
	ClusterBaseRefDescriptor, NodeSpecFieldsDescriptor)

// NodeSpecFieldsDescriptor node spec's fields descriptors.
var NodeSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: RolesField, Type: enumor.Enum, IsRequired: false, IsEditable: true},
	{Field: LabelsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: TaintsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: UnschedulableField, Type: enumor.Boolean, IsRequired: false, IsEditable: true},
	{Field: InternalIPField, Type: enumor.Array, IsRequired: false, IsEditable: true},
	{Field: ExternalIPField, Type: enumor.Array, IsRequired: false, IsEditable: true},
	{Field: HostnameField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: RuntimeComponentField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: KubeProxyModeField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: PodCidrField, Type: enumor.String, IsRequired: false, IsEditable: true},
}

// NodeBaseRefDescriptor the description used when other resources refer to the node.
var NodeBaseRefDescriptor = table.FieldsDescriptors{
	{Field: NodeField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: BKNodeIDField, Type: enumor.Numeric, IsRequired: false, IsEditable: false},
}

// Node node structural description.
type Node struct {
	// ID cluster auto-increment ID in cc
	ID int64 `json:"id,omitempty" bson:"id"`

	// ClusterSpec cluster-related information in the node
	ClusterSpec `json:",inline" bson:",inline"`
	// SupplierAccount the supplier account that this resource belongs to.
	SupplierAccount string `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
	// HostID the node ID to which the host belongs
	HostID int64 `json:"bk_host_id,omitempty" bson:"bk_host_id"`

	// HasPod this field indicates whether there is a pod in the node.
	// if there is a pod, this field is true. If there is no pod, this
	// field is false. this field is false when node is created by default.
	HasPod           *bool                 `json:"has_pod,omitempty" bson:"has_pod"`
	Name             *string               `json:"name,omitempty" bson:"name"`
	Roles            *string               `json:"roles,omitempty" bson:"roles"`
	Labels           *enumor.MapStringType `json:"labels,omitempty" bson:"labels"`
	Taints           *enumor.MapStringType `json:"taints,omitempty" bson:"taints"`
	Unschedulable    *bool                 `json:"unschedulable,omitempty" bson:"unschedulable"`
	InternalIP       *[]string             `json:"internal_ip,omitempty" bson:"internal_ip"`
	ExternalIP       *[]string             `json:"external_ip,omitempty" bson:"external_ip"`
	HostName         *string               `json:"hostname,omitempty" bson:"hostname"`
	RuntimeComponent *string               `json:"runtime_component,omitempty" bson:"runtime_component"`
	KubeProxyMode    *string               `json:"kube_proxy_mode,omitempty" bson:"kube_proxy_mode"`
	PodCidr          *string               `json:"pod_cidr,omitempty" bson:"pod_cidr"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// IgnoredUpdateNodeFields  update fields that need to be ignored in node scenarios
var IgnoredUpdateNodeFields = []string{common.BKFieldID, common.BKAppIDField, ClusterUIDField,
	common.BKFieldName, common.BKOwnerIDField, BKClusterIDField, common.BKHostIDField, HasPodField}

// createValidate validate the NodeBaseFields
func (option *Node) createValidate() ccErr.RawErrorInfo {

	if option == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateCreate(*option, NodeFields); err.ErrCode != 0 {
		return err
	}

	if err := option.validateNodeIP(true); err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{err.Error()},
		}
	}
	return ccErr.RawErrorInfo{}
}

func (option *Node) validateNodeIP(isCreate bool) error {

	if isCreate && option.ExternalIP == nil && option.InternalIP == nil {
		return errors.New("external_ip and internal_ip cannot be null at the same time")
	}

	if (option.ExternalIP != nil && len(*option.ExternalIP) == 0) &&
		(option.InternalIP != nil && len(*option.InternalIP) == 0) {
		return errors.New("the length of external_ip and internal_ip cannot be 0 at the same time")
	}

	return nil
}

// UpdateValidate verifying the validity of parameters for updating node scenarios
func (option *Node) updateValidate() ccErr.RawErrorInfo {

	if option == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateUpdate(*option, NodeFields); err.ErrCode != 0 {
		return err
	}
	return ccErr.RawErrorInfo{}
}

// BatchDeleteNodeOption delete nodes option.
type BatchDeleteNodeOption struct {
	IDs []int64 `json:"ids"`
}

// Validate validate the BatchDeleteNodeOption
func (option *BatchDeleteNodeOption) Validate() ccErr.RawErrorInfo {

	if len(option.IDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"ids"},
		}
	}

	if len(option.IDs) > maxDeleteNodeNum {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", maxDeleteClusterNum},
		}
	}
	return ccErr.RawErrorInfo{}
}

// OneNodeCreateOption node request parameter details.
type OneNodeCreateOption struct {
	// HostID the node ID to which the host belongs
	HostID int64 `json:"bk_host_id" bson:"bk_host_id"`
	// ClusterID the node ID to which the cluster belongs
	ClusterID int64 `json:"bk_cluster_id" bson:"bk_cluster_id"`
	Node      `json:",inline" bson:",inline"`
}

// validateCreate validate the OneNodeCreateOption
func (option *OneNodeCreateOption) validateCreate() ccErr.RawErrorInfo {

	if option.ClusterID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{BKClusterIDField},
		}
	}
	if err := option.createValidate(); err.ErrCode != 0 {
		return err
	}
	return ccErr.RawErrorInfo{}
}

// CreateNodesOption create node requests in batches.
type CreateNodesOption struct {
	Nodes []OneNodeCreateOption `json:"data"`
}

// CreateNodesRsp create the response
// message body of the node result to the user.
type CreateNodesRsp struct {
	metadata.BaseResp
	Data metadata.RspIDs `json:"data"`
}

// ValidateCreate validate the create nodes request
func (option *CreateNodesOption) ValidateCreate() ccErr.RawErrorInfo {

	if len(option.Nodes) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if len(option.Nodes) > maxCreateNodeNum {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"data", maxCreateNodeNum},
		}
	}

	for _, node := range option.Nodes {
		if err := node.validateCreate(); err.ErrCode != 0 {
			return err
		}
	}
	return ccErr.RawErrorInfo{}
}

// CreateNodesResult create node results in batches.
type CreateNodesResult struct {
	metadata.BaseResp
	Info []Node `json:"data" bson:"data"`
}

// QueryNodeOption query node by query builder
type QueryNodeOption struct {
	Filter *filter.Expression `json:"filter"`
	Page   metadata.BasePage  `json:"page"`
	Fields []string           `json:"fields"`
}

// Validate validate the param QueryNodeReq
func (option *QueryNodeOption) Validate() ccErr.RawErrorInfo {
	if err := option.Page.ValidateWithEnableCount(false, common.BKMaxLimitSize); err.ErrCode != 0 {
		return err
	}

	if option.Filter == nil {
		return ccErr.RawErrorInfo{}
	}

	op := filter.NewDefaultExprOpt(NodeFields.FieldsType())
	if err := option.Filter.Validate(op); err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{err.Error()},
		}
	}
	return ccErr.RawErrorInfo{}
}

// SearchNodeRsp query node's response.
type SearchNodeRsp struct {
	Data []Node `json:"node"`
}

// NodeKubeOption information about the node itself.
type NodeKubeOption struct {
	ClusterUID string `json:"cluster_uid"`
	Name       string `json:"name"`
}

// UpdateNodeInfo update individual node details.
type UpdateNodeInfo struct {
	NodeIDs []int64 `json:"ids"`
	Data    Node    `json:"node"`
}

// UpdateNodeOption update node field option
type UpdateNodeOption struct {
	IDs  []int64 `json:"ids"`
	Data Node    `json:"data"`
}

// UpdateValidate check whether the request parameters for updating the node are legal.
func (option *UpdateNodeOption) UpdateValidate() ccErr.RawErrorInfo {

	if len(option.IDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"ids"},
		}
	}
	if len(option.IDs) > maxUpdateNodeNum {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", maxUpdateNodeNum},
		}
	}

	if err := option.Data.updateValidate(); err.ErrCode != 0 {
		return err
	}

	if err := option.Data.validateNodeIP(false); err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"internal_ip or external_ip"},
		}
	}

	return ccErr.RawErrorInfo{}
}

// NodeCondition node condition for search host
type NodeCondition struct {
	Filter *filter.Expression `json:"filter"`
	Fields []string           `json:"fields"`
}
