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
	"configcenter/pkg/filter"
	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/table"
)

const (
	maxDeleteClusterNum = 10
	maxUpdateClusterNum = 10
)

// ClusterFields merge the fields of the cluster and the details corresponding to the fields together.
var ClusterFields = table.MergeFields(CommonSpecFieldsDescriptor, BizIDDescriptor, ClusterSpecFieldsDescriptor)

// ClusterSpecFieldsDescriptor cluster spec's fields descriptors.
var ClusterSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: true},
	{Field: SchedulingEngineField, Type: enumor.String, IsRequired: false, IsEditable: false},
	{Field: UidField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: XidField, Type: enumor.String, IsRequired: false, IsEditable: false},
	{Field: VersionField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: ClusterEnvironmentField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: NetworkTypeField, Type: enumor.Enum, IsRequired: false, IsEditable: true},
	{Field: RegionField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: VpcField, Type: enumor.String, IsRequired: false, IsEditable: false},
	{Field: NetworkField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: TypeField, Type: enumor.String, IsRequired: true, IsEditable: true},
	{Field: ProjectNameField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: ProjectIDField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: ProjectCodeField, Type: enumor.String, IsRequired: false, IsEditable: true},
}

// ClusterBaseRefDescriptor the description used when other resources refer to the cluster.
var ClusterBaseRefDescriptor = table.FieldsDescriptors{
	{Field: ClusterUIDField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: BKClusterIDFiled, Type: enumor.Numeric, IsRequired: false, IsEditable: false},
}

// ClusterSpec describes the common attributes of cluster, it is used by the structure below it.
type ClusterSpec struct {
	// BizID business id in cc
	BizID int64 `json:"bk_biz_id,omitempty" bson:"bk_biz_id"`

	// ClusterID cluster id in cc
	ClusterID int64 `json:"bk_cluster_id,omitempty" bson:"bk_cluster_id"`

	// ClusterUID cluster id in third party platform
	ClusterUID string `json:"cluster_uid,omitempty" bson:"cluster_uid"`
}

// Cluster container cluster table structure
type Cluster struct {
	// ID cluster auto-increment ID in cc
	ID int64 `json:"id" bson:"id"`
	// BizID the business ID to which the cluster belongs
	BizID int64 `json:"bk_biz_id" bson:"bk_biz_id"`
	// TenantID the supplier account that this resource belongs to.
	TenantID string `json:"tenant_id" bson:"tenant_id"`
	// Name cluster name.
	Name *string `json:"name,omitempty" bson:"name"`
	// SchedulingEngine scheduling engines, such as k8s, tke, etc.
	SchedulingEngine *string `json:"scheduling_engine,omitempty" bson:"scheduling_engine"`
	// Uid ID of the cluster itself
	Uid *string `json:"uid,omitempty" bson:"uid"`
	// Xid The underlying cluster ID it depends on
	Xid *string `json:"xid,omitempty" bson:"xid"`
	// Version cluster version
	Version *string `json:"version,omitempty" bson:"version"`
	// NetworkType network type, such as overlay or underlay
	NetworkType *string `json:"network_type,omitempty" bson:"network_type"`
	// Region the region where the cluster is located
	Region *string `json:"region,omitempty" bson:"region"`
	// Vpc vpc network
	Vpc *string `json:"vpc,omitempty" bson:"vpc"`
	// Environment cluster environment
	Environment *string `json:"environment,omitempty" bson:"environment"`
	// NetWork global routing network address (container overlay network) For example: ["1.1.1.0/21"]
	NetWork *[]string `json:"network,omitempty" bson:"network"`
	// Type cluster network type, e.g. INDEPENDENT_CLUSTER, SHARE_CLUSTER etc.
	Type *ClusterType `json:"type,omitempty" bson:"type"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// ClusterType defines the type for cluster type enum
type ClusterType string

// Validate the ClusterType
func (c ClusterType) Validate() ccErr.RawErrorInfo {
	switch c {
	case IndependentClusterType, SharedClusterType:
	default:
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{TypeField},
		}
	}

	return ccErr.RawErrorInfo{}
}

const (
	SharedClusterType      ClusterType = "SHARE_CLUSTER"
	IndependentClusterType ClusterType = "INDEPENDENT_CLUSTER"
)

// IgnoredUpdateClusterFields update the fields that need to be ignored in the cluster scenario.
var IgnoredUpdateClusterFields = []string{common.BKFieldID, common.TenantID, BKBizIDField, ClusterUIDField}

// CreateClusterResult create cluster result for internal call.
type CreateClusterResult struct {
	metadata.BaseResp
	Info *Cluster `json:"data"`
}

// CreateClusterRsp create cluster result for external call.
type CreateClusterRsp struct {
	metadata.BaseResp
	Data metadata.RspID `json:"data"`
}

// DeleteClusterOption delete cluster result.
type DeleteClusterOption struct {
	BizID int64 `json:"bk_biz_id"`
	DeleteClusterByIDsOption
}

// Validate validate the DeleteClusterOption
func (option *DeleteClusterOption) Validate() ccErr.RawErrorInfo {
	if option.BizID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKAppIDField},
		}
	}

	return option.DeleteClusterByIDsOption.Validate()
}

// DeleteClusterByIDsOption delete cluster by ids result.
type DeleteClusterByIDsOption struct {
	IDs []int64 `json:"ids"`
}

// Validate the DeleteClusterByIDsOption
func (option *DeleteClusterByIDsOption) Validate() ccErr.RawErrorInfo {
	if len(option.IDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"ids"},
		}
	}

	if len(option.IDs) > maxDeleteClusterNum {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", maxDeleteClusterNum},
		}
	}
	return ccErr.RawErrorInfo{}
}

// ValidateCreate check whether the parameters for creating a cluster are legal.
func (option *Cluster) ValidateCreate() ccErr.RawErrorInfo {
	if option.BizID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKAppIDField},
		}
	}

	if option == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateCreate(*option, ClusterFields); err.ErrCode != 0 {
		return err
	}

	if option.Type == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{TypeField},
		}
	}

	if rawErr := option.Type.Validate(); rawErr.ErrCode != 0 {
		return rawErr
	}

	return ccErr.RawErrorInfo{}
}

// validateUpdate verifying the validity of parameters for updating node scenarios
func (option *Cluster) validateUpdate() ccErr.RawErrorInfo {

	if option == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"cluster"},
		}
	}

	if err := ValidateUpdate(*option, ClusterFields); err.ErrCode != 0 {
		return err
	}
	return ccErr.RawErrorInfo{}
}

// QueryClusterOption query cluster by query builder
type QueryClusterOption struct {
	BizID  int64              `json:"bk_biz_id"`
	Filter *filter.Expression `json:"filter"`
	Page   metadata.BasePage  `json:"page"`
	Fields []string           `json:"fields"`
}

// Validate the QueryClusterOption
func (option *QueryClusterOption) Validate() ccErr.RawErrorInfo {
	if option.BizID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKAppIDField},
		}
	}

	if err := option.Page.ValidateWithEnableCount(false, common.BKMaxLimitSize); err.ErrCode != 0 {
		return err
	}

	if option.Filter == nil {
		return ccErr.RawErrorInfo{}
	}

	op := filter.NewDefaultExprOpt(ClusterFields.FieldsType())
	if err := option.Filter.Validate(op); err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{err.Error()},
		}
	}
	return ccErr.RawErrorInfo{}
}

// ResponseCluster query the response of the cluster.
type ResponseCluster struct {
	Data []Cluster `json:"cluster"`
}

// UpdateClusterOption update cluster request.
type UpdateClusterOption struct {
	BizID int64 `json:"bk_biz_id"`
	UpdateClusterByIDsOption
}

// Validate the UpdateClusterOption
func (option *UpdateClusterOption) Validate() ccErr.RawErrorInfo {
	if option == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if option.BizID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKAppIDField},
		}
	}
	return option.UpdateClusterByIDsOption.Validate()
}

// UpdateClusterByIDsOption update cluster by ids request。
type UpdateClusterByIDsOption struct {
	IDs  []int64 `json:"ids"`
	Data Cluster `json:"data"`
}

// Validate the UpdateClusterByIDsOption
func (option *UpdateClusterByIDsOption) Validate() ccErr.RawErrorInfo {
	if option == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if len(option.IDs) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"ids"},
		}
	}
	if len(option.IDs) > maxUpdateClusterNum {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", maxUpdateClusterNum},
		}
	}

	if err := option.Data.validateUpdate(); err.ErrCode != 0 {
		return err
	}

	return ccErr.RawErrorInfo{}
}

// UpdateClusterTypeOpt update cluster type option
type UpdateClusterTypeOpt struct {
	BizID int64       `json:"bk_biz_id"`
	ID    int64       `json:"id"`
	Type  ClusterType `json:"type"`
}

// Validate the UpdateClusterTypeOpt
func (option UpdateClusterTypeOpt) Validate() ccErr.RawErrorInfo {
	if option.BizID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKAppIDField},
		}
	}

	if option.ID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKFieldID},
		}
	}

	if rawErr := option.Type.Validate(); rawErr.ErrCode != 0 {
		return rawErr
	}

	return ccErr.RawErrorInfo{}
}

// NsSharedClusterRel shared cluster and biz relationship by namespace dimension.
type NsSharedClusterRel struct {
	NamespaceID int64  `json:"bk_namespace_id" bson:"bk_namespace_id"`
	ClusterID   int64  `json:"bk_cluster_id" bson:"bk_cluster_id"`
	BizID       int64  `json:"bk_biz_id" bson:"bk_biz_id"`
	AsstBizID   int64  `json:"bk_asst_biz_id" bson:"bk_asst_biz_id"`
	TenantID    string `json:"tenant_id" bson:"tenant_id"`
}

// NsSharedClusterRelData namespace and shared cluster relation data.
type NsSharedClusterRelData struct {
	Count uint64               `json:"count"`
	Info  []NsSharedClusterRel `json:"info"`
}
