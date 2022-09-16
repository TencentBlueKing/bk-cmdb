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
var NodeFields = table.MergeFields(NodeSpecFieldsDescriptor)

// NodeSpecFieldsDescriptor node spec's fields descriptors.
var NodeSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, IsRequired: true, IsEditable: false},
	{Field: RolesField, IsRequired: false, IsEditable: true},
	{Field: LabelsField, IsRequired: false, IsEditable: true},
	{Field: TaintsField, IsRequired: false, IsEditable: true},
	{Field: UnschedulableField, IsRequired: false, IsEditable: true},
	{Field: InternalIPField, IsRequired: false, IsEditable: true},
	{Field: ExternalIPField, IsRequired: false, IsEditable: true},
	{Field: HostnameField, IsRequired: false, IsEditable: true},
	{Field: RuntimeComponentField, IsRequired: false, IsEditable: true},
	{Field: KubeProxyModeField, IsRequired: false, IsEditable: true},
	{Field: PodCidrField, IsRequired: false, IsEditable: true},
}

// Node node structural description.
type Node struct {
	// ID cluster auto-increment ID in cc
	ID int64 `json:"id,omitempty" bson:"id"`
	// BizID the business ID to which the cluster belongs
	BizID int64 `json:"bk_biz_id" bson:"bk_biz_id"`
	// SupplierAccount the supplier account that this resource belongs to.
	SupplierAccount string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	// HostID the node ID to which the host belongs
	HostID int64 `json:"bk_host_id,omitempty" bson:"bk_host_id"`
	// ClusterID the node ID to which the cluster belongs
	ClusterID int64 `json:"bk_cluster_id,omitempty" bson:"bk_cluster_id"`
	// ClusterUID the node ID to which the cluster belongs
	ClusterUID string `json:"cluster_uid" bson:"cluster_uid"`
	// NodeFields node base fields
	NodeBaseFields `json:",inline" bson:",inline"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

func initNodeFieldsType() {
	typeOfCat := reflect.TypeOf(NodeBaseFields{})
	valueOf := reflect.ValueOf(NodeBaseFields{})
	for i := 0; i < typeOfCat.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		for _, descripor := range NodeSpecFieldsDescriptor {
			fieldType := typeOfCat.Field(i)
			tag := fieldType.Tag.Get("json")
			if descripor.Field == tag {
				descripor.Type = enumor.GetFieldType(valueOf.Field(i).Type().String())
			}
		}
	}
}

// NodeBaseFields node's basic attribute field description.
type NodeBaseFields struct {
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
}

// CreateValidate validate the NodeBaseFields
func (option *NodeBaseFields) CreateValidate() error {

	if option == nil {
		return errors.New("node information must be given")
	}

	// get a list of required fields.

	typeOfOption := reflect.TypeOf(*option)
	valueOfOption := reflect.ValueOf(*option)
	for i := 0; i < typeOfOption.NumField(); i++ {
		tag := typeOfOption.Field(i).Tag.Get("json")
		if !ClusterFields.IsFieldRequiredByField(tag) {
			continue
		}
		if NodeFields.IsFieldRequiredByField(tag) {
			fieldValue := valueOfOption.Field(i)
			if fieldValue.IsNil() {
				return fmt.Errorf("required fields cannot be empty, %s", tag)
			}
		}
	}

	if err := option.validateNodeIP(); err != nil {
		return err
	}
	return nil
}

func (option *NodeBaseFields) validateNodeIP() error {
	if option.ExternalIP == nil && option.InternalIP == nil {
		return errors.New("external_ip and internal_ip cannot be null at the same time")
	}
	var (
		bExternalIP, bInternalIP bool
	)
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
func (option *NodeBaseFields) updateValidate() error {

	if option == nil {
		return errors.New("node information must be given")
	}

	typeOfOption := reflect.TypeOf(*option)
	valueOfOption := reflect.ValueOf(*option)
	for i := 0; i < typeOfOption.NumField(); i++ {
		fieldValue := valueOfOption.Field(i)
		//	1、check each variable for a null pointer.
		//	if it is a null pointer, it means that
		//	this field will not be updated, skip it directly.
		if fieldValue.IsNil() {
			continue
		}
		// 2、a variable with a non-null pointer gets the corresponding tag.
		tag := typeOfOption.Field(i).Tag.Get("json")
		// 3、get whether it is an editable field based on tag
		if !NodeFields.IsFieldEditableByField(tag) {
			return fmt.Errorf("field [%s] is a non-editable field", tag)
		}
	}
	return nil
}

// OneDeleteNodeOption delete node by id of cmdb.
type OneDeleteNodeOption struct {
	ClusterID int64   `json:"bk_cluster_id"`
	IDs       []int64 `json:"ids"`
}

// BatchDeleteNodeOption delete nodes option.
type BatchDeleteNodeOption struct {
	Nodes []OneDeleteNodeOption `json:"nodes"`
}

// Validate validate the BatchDeleteNodeOption
func (option *BatchDeleteNodeOption) Validate() error {

	if len(option.Nodes) == 0 {
		return errors.New("params must be set")
	}

	if len(option.Nodes) > maxDeleteNodeNum {
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
	ClusterID      int64 `json:"bk_cluster_id" bson:"bk_cluster_id"`
	NodeBaseFields `json:",inline" bson:",inline"`
}

// ValidateCreate validate the OneNodeCreateOption
func (option *OneNodeCreateOption) ValidateCreate() error {
	if option.HostID == 0 {
		return errors.New("host id must be set")
	}
	if option.ClusterID == 0 {
		return errors.New("cluster id must be set")
	}
	if err := option.CreateValidate(); err != nil {
		return err
	}
	return nil
}

// CreateNodesOption create node requests in batches.
type CreateNodesOption struct {
	Nodes []OneNodeCreateOption `json:"nodes"`
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
		if err := node.ValidateCreate(); err != nil {
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
	NodeIDs []int64         `json:"ids"`
	Data    *NodeBaseFields `json:"data"`
}

// UpdateNodeOption update node field option
type UpdateNodeOption struct {
	Nodes []UpdateNodeInfo `json:"nodes"`
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
