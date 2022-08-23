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

	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/storage/dal/table"
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

type ClusterID struct {
	ID  int64 `json:"id"`
	Uid int64 `json:"uid"`
}

type DeleteClusterOption struct {
	IDs  []int64 `json:"id"`
	Uids []int64 `json:"uid"`
}

const (
	maxDeleteClusterNum = 10
)

// Validate validate the  DeleteClusterOption
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

// ClusterBaseFields 创建集群请求字段
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

// ResponseCluster
type ResponseCluster struct {
	Data []Cluster `json:"cluster"`
}

// QueryClusterInfo query cluster response
type QueryClusterInfo struct {
	Info  []Cluster `json:"info"`
	Count int       `json:"count"`
}

// ValidateCreate 校验创建集群参数是否合法
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

// KubeResourceInfo 请求资源的种类和对应的资源ID
type KubeResourceInfo struct {
	Kind string `json:"kind"`
	ID   int64  `json:"id"`
}

// KubeTopoCountReq 计算资源节点主机或者pod数量的请求
type KubeTopoCountReq struct {
	ResourceInfos []KubeResourceInfo `json:"resource_info"`
}

// Validate validate the KubeTopoCountReq
func (option *KubeTopoCountReq) Validate() ccErr.RawErrorInfo {
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
