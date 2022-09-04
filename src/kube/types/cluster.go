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
	"configcenter/src/storage/dal/table"
)

const (
	maxDeleteClusterNum = 10
	maxDeleteNodeNum    = 100
	maxUpdateClusterNum = 10
)

// ClusterFields merge the fields of the cluster and the details corresponding to the fields together.
var ClusterFields = table.MergeFields(ClusterFieldsDescriptor)

// ClusterFieldsDescriptor cluster's fields descriptors.
var ClusterFieldsDescriptor = table.MergeFieldDescriptors(
	table.FieldsDescriptors{
		{Field: BKIDField, Type: enumor.Numeric, IsRequired: true, IsEditable: false},
		{Field: BKBizIDField, Type: enumor.Numeric, IsRequired: true, IsEditable: false},
		{Field: BKSupplierAccountField, Type: enumor.String, IsRequired: true, IsEditable: false},
		{Field: CreatorField, Type: enumor.String, IsRequired: true, IsEditable: false},
		{Field: ModifierField, Type: enumor.String, IsRequired: true, IsEditable: true},
		{Field: CreateTimeField, Type: enumor.Numeric, IsRequired: true, IsEditable: false},
		{Field: LastTimeField, Type: enumor.Numeric, IsRequired: true, IsEditable: true},
	},
	table.MergeFieldDescriptors(ClusterSpecFieldsDescriptor),
)

// ClusterSpecFieldsDescriptor cluster spec's fields descriptors.
var ClusterSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: SchedulingEngineField, Type: enumor.String, IsRequired: false, IsEditable: false},
	{Field: UidField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: XidField, Type: enumor.String, IsRequired: false, IsEditable: false},
	{Field: VersionField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: NetworkTypeField, Type: enumor.Enum, IsRequired: false, IsEditable: true},
	{Field: RegionField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: VpcField, Type: enumor.String, IsRequired: false, IsEditable: false},
	{Field: NetworkField, Type: enumor.String, IsRequired: false, IsEditable: false},
	{Field: TypeField, Type: enumor.String, IsRequired: false, IsEditable: true},
}

// Cluster container cluster table structure
type Cluster struct {
	// ID cluster auto-increment ID in cc
	ID *int64 `json:"id" bson:"id"`
	// BizID the business ID to which the cluster belongs
	BizID *int64 `json:"bk_biz_id" bson:"bk_biz_id"`
	// ClusterFields cluster base fields
	ClusterBaseFields `json:",inline" bson:",inline"`
	// SupplierAccount the supplier account that this resource belongs to.
	SupplierAccount *string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// CreateClusterResult 创建集群结果
type CreateClusterResult struct {
	metadata.BaseResp
	ID int64 `field:"id" json:"id" bson:"id"`
}

// CreateContainerResult 创建容器结果
type CreateContainerResult struct {
	metadata.BaseResp
	ID int64 `field:"id" json:"id" bson:"id"`
}

// CreatePodResult 创建pod结果
type CreatePodResult struct {
	metadata.BaseResp
	ID int64 `field:"id" json:"id" bson:"id"`
}

// DeleteClusterOption 删除集群的请求
type DeleteClusterOption struct {
	IDs  []int64 `json:"ids"`
	Uids []int64 `json:"uids"`
}

// DeleteNodeCmdbOption 通过cmdb的id进行删除node
type DeleteNodeCmdbOption struct {
	ClusterID int64 `json:"bk_cluster_id"`
	ID        int64 `json:"id"`
}

// DeleteNodeOption 通过原生id进行删除node
type DeleteNodeOption struct {
	ClusterUID string `json:"cluster_uid"`
	Name       string `json:"name"`
}

// BatchDeleteNodeOption 删除node请求
type BatchDeleteNodeOption struct {
	NodeCmdbIDs []DeleteNodeCmdbOption `json:"node_cmdb_ids"`
	NodeIDs     []DeleteNodeOption     `json:"node_ids"`
}

// Validate validate the BatchDeleteNodeOption
func (option *BatchDeleteNodeOption) Validate() error {
	if len(option.NodeIDs) > 0 && len(option.NodeCmdbIDs) > 0 {
		return errors.New("params cannot be set at the same time")
	}
	if len(option.NodeIDs) == 0 && len(option.NodeCmdbIDs) == 0 {
		return errors.New("params must be set")
	}
	if len(option.NodeIDs) > maxDeleteNodeNum || len(option.NodeCmdbIDs) > maxDeleteNodeNum {
		return fmt.Errorf("the maximum number of nodes to be deleted is not allowed to exceed %d",
			maxDeleteClusterNum)
	}
	return nil
}

// Validate validate the DeleteClusterOption
func (option *DeleteClusterOption) Validate() error {

	if len(option.IDs) > 0 && len(option.Uids) > 0 {
		return errors.New("cannot fill in the id and uid fields at the same time")
	}

	if len(option.IDs) == 0 && len(option.Uids) == 0 {
		return errors.New("cluster id or uid must be set at least one")
	}

	if len(option.IDs) > maxDeleteClusterNum || len(option.Uids) > maxDeleteClusterNum {
		return fmt.Errorf("the maximum number of clusters to be deleted is not allowed to exceed %d",
			maxDeleteClusterNum)
	}
	return nil
}

// ClusterBaseFields basic description fields for container clusters.
type ClusterBaseFields struct {
	Name *string `json:"name" bson:"name"`
	// SchedulingEngine scheduling engines, such as k8s, tke, etc.
	SchedulingEngine *string `json:"scheduling_engine" bson:"scheduling_engine"`
	// Uid ID of the cluster itself
	Uid *string `json:"uid" bson:"uid"`
	// Xid The underlying cluster ID it depends on
	Xid *string `json:"xid" bson:"xid"`
	// Version cluster version
	Version *string `json:"version" bson:"version"`
	// NetworkType network type, such as overlay or underlay
	NetworkType *string `json:"network_type" bson:"network_type"`
	// Region the region where the cluster is located
	Region *string `json:"region" bson:"region"`
	// Vpc vpc network
	Vpc *string `json:"vpc" bson:"vpc"`
	// NetWork global routing network address (container overlay network) For example: ["1.1.1.0/21"]
	NetWork *[]string `json:"network" bson:"network"`
	// Type cluster network type, e.g. public clusters, private clusters, etc.
	Type *string `json:"type" bson:"type"`
}

// QueryClusterReq query cluster by query builder
type QueryClusterReq struct {
	Filter *querybuilder.QueryFilter `json:"filter"`
	Page   metadata.BasePage         `json:"page"`
	Fields []string                  `json:"fields"`
}

// Validate validate the QueryClusterReq
func (option *QueryClusterReq) Validate() ccErr.RawErrorInfo {
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

// ResponseCluster  query the response of the container cluster.
type ResponseCluster struct {
	Data []Cluster `json:"cluster"`
}

// UpdateClusterOption update cluster request。
type UpdateClusterOption struct {
	Clusters []OneUpdateCluster `json:"clusters"`
}

// OneUpdateCluster update individual cluster information.
type OneUpdateCluster struct {
	ID   int64             `json:"id"`
	UID  string            `json:"uid"`
	Data ClusterBaseFields `json:"data"`
}

// Validate validate the UpdateClusterOption
func (option *UpdateClusterOption) Validate() error {

	if option == nil {
		return errors.New("cluster option is null")
	}

	if len(option.Clusters) == 0 {
		return errors.New("the params for updating the cluster must be set")
	}
	if len(option.Clusters) > maxUpdateClusterNum {
		return fmt.Errorf("the number of update clusters cannot exceed %d at a time", maxUpdateClusterNum)
	}

	for _, one := range option.Clusters {
		if one.UID == "" && one.ID == 0 {
			return errors.New("id and uid cannot be empty at the same time")
		}
		if one.UID != "" && one.ID != 0 {
			return errors.New("id and uid cannot be set at the same time")
		}
		if err := one.Data.UpdateValidate(); err != nil {
			return err
		}
	}
	return nil
}

// UpdateValidate verifying the validity of parameters for updating node scenarios
func (option *ClusterBaseFields) UpdateValidate() error {
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
		if !ClusterFields.IsFieldEditableByField(tag) {
			return fmt.Errorf("field [%s] is a non-editable field", tag)
		}
	}
	return nil
}

// ValidateCreate check whether the parameters for creating a cluster are legal.
func (option *ClusterBaseFields) ValidateCreate() error {

	if option.Name == nil || *option.Name == "" {
		return errors.New("name can not be empty")
	}
	//if err := ValidateString(*option.Name, StringSettings{}); err != nil {
	//	return err
	//}

	if option.Uid == nil || *option.Uid == "" {
		return errors.New("uid can not be empty")
	}

	//if err := ValidateString(*option.Uid, StringSettings{}); err != nil {
	//	return err
	//}
	return nil
}

// KubeResourceInfo the type of the requested resource and the corresponding resource ID.
// it should be noted that when the kind is folder, the host cannot be obtained through
// the pod table. In this case, the node table needs to be used to find the corresponding
// number of hosts. Since the node is only associated with the cluster, the id in this
// scenario needs to pass the corresponding clusterID.
type KubeResourceInfo struct {
	Kind string `json:"kind"`
	ID   int64  `json:"id"`
}

// KubeTopoCountOption calculate the number of hosts or pods under the container resource node.
type KubeTopoCountOption struct {
	ResourceInfos []KubeResourceInfo `json:"resource_info"`
}

// Validate validate the KubeTopoCountOption
func (option *KubeTopoCountOption) Validate() ccErr.RawErrorInfo {
	if len(option.ResourceInfos) > 100 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{errors.New("the requested array length exceeds the maximum value of 100")},
		}
	}
	for _, info := range option.ResourceInfos {
		if !IsKubeResourceKind(info.Kind) {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{errors.New("non-container resource objects\n")},
			}
		}
	}
	return ccErr.RawErrorInfo{}
}

// KubeTopoCountRsp 节点主机或者Pod数量的回应
type KubeTopoCountRsp struct {
	Kind  string `json:"kind"`
	ID    int64  `json:"id"`
	Count int64  `json:"count"`
}

// KubeTopoPathReq 获取容器拓扑路径请求
type KubeTopoPathReq struct {
	ReferenceObjID string            `json:"bk_reference_obj_id"`
	ReferenceID    int64             `json:"bk_reference_id"`
	Page           metadata.BasePage `json:"page"`
}

// Validate validate the KubeTopoPathReq
func (option *KubeTopoPathReq) Validate() ccErr.RawErrorInfo {

	if option.ReferenceID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{errors.New("bk_reference_id must be set")},
		}
	}

	// 判断下是都合法是否合法
	if !IsContainerTopoResource(option.ReferenceObjID) {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{errors.New("bk_reference_obj_id is illegal")},
		}
	}

	if err := option.Page.ValidateWithEnableCount(false, common.BKMaxLimitSize); err.ErrCode != 0 {
		return err
	}
	return ccErr.RawErrorInfo{}
}

// KubeObjectInfo 容器对象信息
type KubeObjectInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

// KubeTopoPathRsp 获取拓扑路径回应
type KubeTopoPathRsp struct {
	Info  []KubeObjectInfo `json:"info"`
	Count int              `json:"count"`
}
