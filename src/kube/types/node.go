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

// NodeFields merge the fields of the cluster and the details corresponding to the fields together.
var NodeFields = table.MergeFields(NodeSpecFieldsDescriptor)

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

// Node node 结构描述
type Node struct {
	// ID cluster auto-increment ID in cc
	ID int64 `json:"id" bson:"id"`
	// BizID the business ID to which the cluster belongs
	BizID int64 `json:"bk_biz_id" bson:"bk_biz_id"`
	// NodeFields node base fields
	NodeBaseFields `json:",inline" bson:",inline"`
	// SupplierAccount the supplier account that this resource belongs to.
	SupplierAccount string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// NodeBaseFields node的基础属性字段描述
type NodeBaseFields struct {
	// HostID the node ID to which the host belongs
	HostID *int64 `json:"bk_host_id" bson:"bk_host_id"`
	// ClusterID the node ID to which the cluster belongs
	ClusterID *int64 `json:"bk_cluster_id" bson:"bk_cluster_id"`
	// ClusterUID the node ID to which the cluster belongs
	ClusterUID *string `json:"cluster_uid" bson:"cluster_uid"`
	// HasPod this field indicates whether there is a pod in the node.
	// if there is a pod, this field is true. If there is no pod, this
	// field is false. this field is false when node is created by default.
	HasPod           *bool                 `json:"has_pod" bson:"has_pod"`
	Name             *string               `json:"name" bson:"name"`
	Roles            *string               `json:"roles" bson:"roles"`
	Labels           *enumor.MapStringType `json:"labels" bson:"labels"`
	Taints           *enumor.MapStringType `json:"taints" bson:"taints"`
	Unschedulable    bool                  `json:"unschedulable" bson:"unschedulable"`
	InternalIP       *[]string             `json:"internal_ip" bson:"internal_ip"`
	ExternalIP       *[]string             `json:"external_ip" bson:"external_ip"`
	HostName         *string               `json:"hostname" bson:"hostname"`
	RuntimeComponent *string               `json:"runtime_component" bson:"runtime_component"`
	KubeProxyMode    *string               `json:"kube_proxy_mode" bson:"kube_proxy_mode"`
	PodCidr          *string               `json:"pod_cidr" bson:"pod_cidr"`
}

// UpdateValidate 校验更新node场景的参数有效性
func (option *NodeBaseFields) UpdateValidate() error {
	if option == nil {
		return errors.New("node information must be given")
	}
	typeOfOption := reflect.TypeOf(*option)
	valueOfOption := reflect.ValueOf(*option)
	for i := 0; i < typeOfOption.NumField(); i++ {
		fieldValue := valueOfOption.Field(i)
		//	1、查看每个变量是否为空指针。如果是空指针表示不更新此字段，直接跳过
		if fieldValue.IsNil() {
			continue
		}
		// 2、非空指针的变量获取对应的tag。
		tag := typeOfOption.Field(i).Tag.Get("json")
		// 3、根据tag获取是否是可编辑字段
		if !NodeFields.IsFieldEditableByField(tag) {
			return fmt.Errorf("field [%s] is a non-editable field", tag)
		}
	}
	return nil
}

// CreateNodesReq create node requests in batches.
type CreateNodesReq struct {
	Nodes []NodeBaseFields `json:"nodes"`
}

// ArrangeDeleteNodeOption cleaned up delete Node request
type ArrangeDeleteNodeOption struct {
	Option map[interface{}][]interface{} `json:"option"`
	Flag   bool                          `json:"flag"`
}

// CreateNodesResult create node results in batches.
type CreateNodesResult struct {
	metadata.BaseResp
	Info []int64 `json:"ids" bson:"ids"`
}

// QueryNodeReq query node by query builder
type QueryNodeReq struct {
	Filter     *querybuilder.QueryFilter `json:"filter"`
	ClusterID  int64                     `json:"bk_cluster_id"`
	HostID     int64                     `json:"bk_host_id"`
	ClusterUID int64                     `json:"cluster_uid"`
	Page       metadata.BasePage         `json:"page"`
	Fields     []string                  `json:"fields"`
}

// ResponseNode 查询node的回应
type ResponseNode struct {
	Data []Node `json:"node"`
}

// Validate validate the param QueryNodeReq
func (option *QueryNodeReq) Validate() ccErr.RawErrorInfo {
	op := &querybuilder.RuleOption{
		NeedSameSliceElementType: true,
		MaxSliceElementsCount:    querybuilder.DefaultMaxSliceElementsCount,
		MaxConditionOrRulesCount: querybuilder.DefaultMaxConditionOrRulesCount,
	}

	if err := option.Page.ValidateWithEnableCount(false, common.BKMaxLimitSize); err.ErrCode != 0 {
		return err
	}

	if option.ClusterUID > 0 && option.ClusterID > 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{errors.New("the param cluster_id and cluster_uid can only be filled in one")},
		}
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

// NodeCmdbOption node related information in cmdb.
type NodeCmdbOption struct {
	ClusterID int `json:"bk_cluster_id"`
	ID        int `json:"id"`
}

// NodeKubeOption information about the node itself.
type NodeKubeOption struct {
	ClusterUID string `json:"cluster_uid"`
	Name       string `json:"name"`
}

// OneUpdateNode update individual node details.
type OneUpdateNode struct {
	NodeCmdbFilter *NodeCmdbOption `json:"node_cmdb_filter"`
	NodeKubeFilter *NodeKubeOption `json:"node_kube_filter"`
	Data           *NodeBaseFields `json:"data"`
}

// UpdateNodeOption update node field option
type UpdateNodeOption struct {
	Nodes []OneUpdateNode `json:"nodes"`
}

// Validate check whether the request parameters for updating the node are legal.
func (option *UpdateNodeOption) Validate() error {

	if len(option.Nodes) == 0 {
		return errors.New("parameter cannot be empty")
	}

	for _, node := range option.Nodes {
		if node.NodeCmdbFilter == nil && node.NodeKubeFilter == nil {
			return errors.New("node_cmdb_filter and node_kube_filter cannot be empty at the same time")
		}
		if node.NodeCmdbFilter != nil && node.NodeKubeFilter != nil {
			return errors.New("node_cmdb_filter and node_kube_filter cannot be set at the same time")
		}
		if err := node.Data.UpdateValidate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate validate the NodeBaseFields
func (node *NodeBaseFields) Validate() error {
	if *node.HostID == 0 {
		return errors.New("host id must be set")
	}
	if *node.ClusterID == 0 {
		return errors.New("cluster id must be set")
	}
	//if err := ValidateString(*node.ClusterUID, StringSettings{}); err != nil {
	//	return err
	//}

	//if err := ValidateString(*node.Name, StringSettings{}); err != nil {
	//	return err
	//}
	return nil
}

// ValidateCreate validate the create nodes request
func (option *CreateNodesReq) ValidateCreate() error {

	if len(option.Nodes) == 0 {
		return errors.New("param must be set")
	}

	if len(option.Nodes) > 100 {
		return errors.New("the number of nodes created at one time does not exceed 100")
	}
	for _, node := range option.Nodes {
		if err := node.Validate(); err != nil {
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
