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
	{Field: NetworkTypeField, Type: enumor.Enum, IsRequired: false, IsEditable: true},
	{Field: RegionField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: VpcField, Type: enumor.String, IsRequired: false, IsEditable: false},
	{Field: NetworkField, Type: enumor.String, IsRequired: false, IsEditable: false},
	{Field: TypeField, Type: enumor.String, IsRequired: false, IsEditable: true},
}

// Cluster container cluster table structure
type Cluster struct {
	// ID cluster auto-increment ID in cc
	ID int64 `json:"id" bson:"id"`
	// BizID the business ID to which the cluster belongs
	BizID int64 `json:"bk_biz_id" bson:"bk_biz_id"`
	// SupplierAccount the supplier account that this resource belongs to.
	SupplierAccount string `json:"bk_supplier_account" bson:"bk_supplier_account"`
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
	// NetWork global routing network address (container overlay network) For example: ["1.1.1.0/21"]
	NetWork *[]string `json:"network,omitempty" bson:"network"`
	// Type cluster network type, e.g. public clusters, private clusters, etc.
	Type *string `json:"type,omitempty" bson:"type"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// IgnoredUpdateClusterFields update the fields that need to be ignored in the cluster scenario.
var IgnoredUpdateClusterFields = []string{common.BKFieldID, common.BKOwnerIDField, BKBizIDField, ClusterUIDField}

// CreateClusterResult create cluster result.
type CreateClusterResult struct {
	metadata.BaseResp
	Info *Cluster `json:"data"`
}

// DeleteClusterOption delete cluster result.
type DeleteClusterOption struct {
	IDs []int64 `json:"ids"`
}

// Validate validate the DeleteClusterOption
func (option *DeleteClusterOption) Validate() error {

	if len(option.IDs) == 0 {
		return errors.New("cluster ids must be set")
	}

	if len(option.IDs) > maxDeleteClusterNum {
		return fmt.Errorf("the maximum number of clusters to be deleted is not allowed to exceed %d",
			maxDeleteClusterNum)
	}
	return nil
}

// CreateValidate check whether the parameters for creating a cluster are legal.
func (option *Cluster) CreateValidate() error {

	if option == nil {
		return errors.New("cluster information must be given")
	}

	typeOfOption := reflect.TypeOf(*option)
	valueOfOption := reflect.ValueOf(*option)
	for i := 0; i < typeOfOption.NumField(); i++ {
		tag, flag := getFieldTag(typeOfOption, i)
		if flag {
			continue
		}

		if !ClusterFields.IsFieldRequiredByField(tag) {
			continue
		}

		if err := isRequiredField(tag, valueOfOption, i); err != nil {
			return err
		}
	}
	return nil
}

// UpdateValidate verifying the validity of parameters for updating node scenarios
func (option *Cluster) updateValidate() error {

	if option == nil {
		return errors.New("cluster information must be given")
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
		if !ClusterFields.IsFieldEditableByField(tag) {
			return fmt.Errorf("field [%s] is a non-editable field", tag)
		}
	}
	return nil
}

// QueryClusterOption query cluster by query builder
type QueryClusterOption struct {
	Filter *querybuilder.QueryFilter `json:"filter"`
	Page   metadata.BasePage         `json:"page"`
	Fields []string                  `json:"fields"`
}

// Validate validate the QueryClusterOption
func (option *QueryClusterOption) Validate() ccErr.RawErrorInfo {
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

// ResponseCluster query the response of the cluster.
type ResponseCluster struct {
	Data []Cluster `json:"cluster"`
}

// OneUpdateCluster update individual cluster information.
type OneUpdateCluster struct {
	ID   int64   `json:"id"`
	Data Cluster `json:"cluster"`
}

// UpdateClusterOption update cluster request。
type UpdateClusterOption struct {
	Clusters []OneUpdateCluster `json:"data"`
}

// Validate validate the UpdateClusterOption
func (option *UpdateClusterOption) Validate() error {

	if option == nil {
		return errors.New("cluster information must be given")
	}

	if len(option.Clusters) == 0 {
		return errors.New("the params for updating the cluster must be set")
	}
	if len(option.Clusters) > maxUpdateClusterNum {
		return fmt.Errorf("the number of update clusters cannot exceed %d at a time", maxUpdateClusterNum)
	}

	for _, one := range option.Clusters {
		if one.ID == 0 {
			return errors.New("id cannot be empty")
		}
		if err := one.Data.updateValidate(); err != nil {
			return err
		}
	}

	return nil
}
