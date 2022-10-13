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
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/filter"
	"configcenter/src/storage/dal/table"
)

const (
	maxDeleteNodeNum = 100
	maxCreateNodeNum = 100
)

// NodeFields merge the fields of the cluster and the details corresponding to the fields together.
var NodeFields = table.MergeFields(CommonSpecFieldsDescriptor, BizIDDescriptor, NodeSpecFieldsDescriptor)

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

// Node node structural description.
type Node struct {
	// ID cluster auto-increment ID in cc
	ID int64 `json:"id,omitempty" bson:"id"`
	// BizID the business ID to which the cluster belongs
	BizID int64 `json:"bk_biz_id,omitempty" bson:"bk_biz_id"`
	// SupplierAccount the supplier account that this resource belongs to.
	SupplierAccount string `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
	// HostID the node ID to which the host belongs
	HostID int64 `json:"bk_host_id,omitempty" bson:"bk_host_id"`
	// ClusterID the node ID to which the cluster belongs
	ClusterID int64 `json:"bk_cluster_id,omitempty" bson:"bk_cluster_id"`
	// ClusterUID the node ID to which the cluster belongs
	ClusterUID string `json:"cluster_uid,omitempty" bson:"cluster_uid"`

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
	common.BKFieldName, common.BKOwnerIDField, BKClusterIDFiled, common.BKHostIDField, HasPodField}

// createValidate validate the NodeBaseFields
func (option *Node) createValidate() error {

	if option == nil {
		return errors.New("node information must be given")
	}

	// get a list of required fields.
	typeOfOption := reflect.TypeOf(*option)
	valueOfOption := reflect.ValueOf(*option)
	for i := 0; i < typeOfOption.NumField(); i++ {

		tag, flag := getFieldTag(typeOfOption, i)
		if flag {
			continue
		}

		if !NodeFields.IsFieldRequiredByField(tag) {
			continue
		}

		if err := isRequiredField(tag, valueOfOption, i); err != nil {
			return err
		}
	}

	if err := option.validateNodeIP(true); err != nil {
		return err
	}
	return nil
}

func (option *Node) validateNodeIP(isCreate bool) error {

	var (
		bExternalIP, bInternalIP bool
	)

	if isCreate && option.ExternalIP == nil && option.InternalIP == nil {
		return errors.New("external_ip and internal_ip cannot be null at the same time")
	}

	if option.ExternalIP != nil && len(*option.ExternalIP) == 0 {
		bExternalIP = true
	}
	if option.InternalIP != nil && len(*option.InternalIP) == 0 {
		bInternalIP = true
	}

	if bExternalIP && bInternalIP {
		return errors.New("external_ip and internal_ip cannot be null at the same time")
	}
	return nil
}

// UpdateValidate verifying the validity of parameters for updating node scenarios
func (option *Node) updateValidate() error {

	if option == nil {
		return errors.New("node information must be given")
	}

	typeOfOption := reflect.TypeOf(*option)
	valueOfOption := reflect.ValueOf(*option)
	for i := 0; i < typeOfOption.NumField(); i++ {
		tag, flag := getFieldTag(typeOfOption, i)
		if flag {
			continue
		}

		if flag := isEditableField(tag, valueOfOption, i); flag {
			continue
		}

		// get whether it is an editable field based on tag
		if !NodeFields.IsFieldEditableByField(tag) {
			return fmt.Errorf("field [%s] is a non-editable field", tag)
		}
	}
	return nil
}

// BatchDeleteNodeOption delete nodes option.
type BatchDeleteNodeOption struct {
	IDs []int64 `json:"ids"`
}

// Validate validate the BatchDeleteNodeOption
func (option *BatchDeleteNodeOption) Validate() error {

	if len(option.IDs) == 0 {
		return errors.New("node ids must be set")
	}

	if len(option.IDs) > maxDeleteNodeNum {
		return fmt.Errorf("the maximum number of nodes to be deleted is not allowed to exceed %d",
			maxDeleteClusterNum)
	}
	return nil
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
func (option *OneNodeCreateOption) validateCreate() error {

	if option.ClusterID == 0 {
		return errors.New("cluster id must be set")
	}
	if err := option.createValidate(); err != nil {
		return err
	}
	return nil
}

// CreateNodesOption create node requests in batches.
type CreateNodesOption struct {
	Nodes []OneNodeCreateOption `json:"data"`
}

// ValidateCreate validate the create nodes request
func (option *CreateNodesOption) ValidateCreate() error {

	if len(option.Nodes) == 0 {
		return errors.New("param must be set")
	}

	if len(option.Nodes) > maxCreateNodeNum {
		return fmt.Errorf("the number of nodes created at one time does not exceed %d", maxCreateNodeNum)
	}

	for _, node := range option.Nodes {
		if err := node.validateCreate(); err != nil {
			return err
		}
	}
	return nil
}

// CreateNodesResult create node results in batches.
type CreateNodesResult struct {
	metadata.BaseResp
	Info []Node `json:"data" bson:"data"`
}

// QueryNodeOption query node by query builder
type QueryNodeOption struct {
	Filter    *querybuilder.QueryFilter `json:"filter"`
	ClusterID int64                     `json:"bk_cluster_id"`
	HostID    int64                     `json:"bk_host_id"`
	Page      metadata.BasePage         `json:"page"`
	Fields    []string                  `json:"fields"`
}

// Validate validate the param QueryNodeReq
func (option *QueryNodeOption) Validate() ccErr.RawErrorInfo {
	op := &querybuilder.RuleOption{
		NeedSameSliceElementType: true,
		MaxSliceElementsCount:    querybuilder.DefaultMaxSliceElementsCount,
		MaxConditionOrRulesCount: querybuilder.DefaultMaxConditionOrRulesCount,
	}

	if err := option.Page.ValidateWithEnableCount(false, common.BKMaxLimitSize); err.ErrCode != 0 {
		return err
	}

	if option.Filter == nil {
		return ccErr.RawErrorInfo{}
	}

	if invalidKey, err := option.Filter.Validate(op); err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{fmt.Errorf("conditions.%s, err: %s", invalidKey, err.Error())},
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
	Nodes []UpdateNodeInfo `json:"data"`
}

// Validate check whether the request parameters for updating the node are legal.
func (option *UpdateNodeOption) Validate() error {

	if len(option.Nodes) == 0 {
		return errors.New("parameter cannot be empty")
	}

	for _, node := range option.Nodes {
		if len(node.NodeIDs) == 0 {
			return errors.New("node_ids must be set")
		}
		if err := node.Data.validateNodeIP(false); err != nil {
			return err
		}
		if err := node.Data.updateValidate(); err != nil {
			return err
		}
	}
	return nil
}

// SearchHostReq search host request
type SearchHostReq struct {
	BizID       int64                    `json:"bk_biz_id"`
	ClusterID   int64                    `json:"bk_cluster_id"`
	Folder      bool                     `json:"folder"`
	NamespaceID int64                    `json:"bk_namespace_id"`
	WorkloadID  int64                    `json:"bk_workload_id"`
	WlKind      WorkloadType             `json:"kind"`
	NodeCond    *NodeCond                `json:"node_cond"`
	Ip          metadata.IPInfo          `json:"ip"`
	HostCond    metadata.SearchCondition `json:"host_condition"`
	Page        metadata.BasePage        `json:"page"`
}

// NodeCond node condition for search host
type NodeCond struct {
	Filter *filter.Expression `json:"filter"`
	Fields []string           `json:"fields"`
}
