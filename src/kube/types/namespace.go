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
	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/filter"
	"configcenter/src/storage/dal/table"
)

// NamespaceSpecFieldsDescriptor namespace spec's fields descriptors.
var NamespaceSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: LabelsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: ClusterUIDField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: ResourceQuotasField, Type: enumor.Array, IsRequired: false, IsEditable: true},
}

// ScopeSelectorOperator a scope selector operator is the set of operators
// that can be used in a scope selector requirement.
type ScopeSelectorOperator string

const (
	// ScopeSelectorOpIn in operator for scope selector
	ScopeSelectorOpIn ScopeSelectorOperator = "In"
	// ScopeSelectorOpNotIn not in operator for scope selector
	ScopeSelectorOpNotIn ScopeSelectorOperator = "NotIn"
	// ScopeSelectorOpExists exists operator for scope selector
	ScopeSelectorOpExists ScopeSelectorOperator = "Exists"
	// ScopeSelectorOpDoesNotExist not exists operator for scope selector
	ScopeSelectorOpDoesNotExist ScopeSelectorOperator = "DoesNotExist"
)

// ResourceQuotaScope defines a filter that must match each object tracked by a quota
type ResourceQuotaScope string

const (
	// ResourceQuotaScopeTerminating match all pod objects where spec.activeDeadlineSeconds >=0
	ResourceQuotaScopeTerminating ResourceQuotaScope = "Terminating"
	// ResourceQuotaScopeNotTerminating match all pod objects where spec.activeDeadlineSeconds is nil
	ResourceQuotaScopeNotTerminating ResourceQuotaScope = "NotTerminating"
	// ResourceQuotaScopeBestEffort match all pod objects that have best effort quality of service
	ResourceQuotaScopeBestEffort ResourceQuotaScope = "BestEffort"
	// ResourceQuotaScopeNotBestEffort match all pod objects that do not have best effort quality of service
	ResourceQuotaScopeNotBestEffort ResourceQuotaScope = "NotBestEffort"
	// ResourceQuotaScopePriorityClass match all pod objects that have priority class mentioned
	ResourceQuotaScopePriorityClass ResourceQuotaScope = "PriorityClass"
	// ResourceQuotaScopeCrossNamespacePodAffinity match all pod objects that have cross-namespace pod
	// (anti)affinity mentioned.
	ResourceQuotaScopeCrossNamespacePodAffinity ResourceQuotaScope = "CrossNamespacePodAffinity"
)

var (
	// NsUpdateLimit limit on the number of namespace updates
	NsUpdateLimit = 200
	// NsDeleteLimit limit on the number of namespace delete
	NsDeleteLimit = 200
	// NsCreateLimit limit on the number of namespace create
	NsCreateLimit = 200
)

// Namespace define the namespace struct.
type Namespace struct {
	ClusterSpec     `json:",inline" bson:",inline"`
	ID              *int64             `json:"id" bson:"id"`
	Name            *string            `json:"name" bson:"name"`
	Labels          *map[string]string `json:"labels" bson:"labels"`
	ResourceQuotas  *[]ResourceQuota   `json:"resource_quotas" bson:"resource_quotas"`
	CreateTime      *int64             `json:"create_time" bson:"create_time"`
	UpdateTime      *int64             `json:"update_time" bson:"update_time"`
	SupplierAccount *string            `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// ValidateCreate validate create namespace
func (ns *Namespace) ValidateCreate() errors.RawErrorInfo {
	if ns.ClusterUID == nil && ns.ClusterID == nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{ClusterUIDField + " or " + BKClusterIDFiled},
		}
	}

	if ns.Name == nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKFieldName},
		}
	}

	return errors.RawErrorInfo{}
}

// ValidateUpdate validate update namespace
func (ns *Namespace) ValidateUpdate() errors.RawErrorInfo {
	// todo
	return errors.RawErrorInfo{}
}

// ResourceQuota defines the desired hard limits to enforce for Quota.
type ResourceQuota struct {
	Hard          map[string]string    `json:"hard" bson:"hard"`
	Scopes        []ResourceQuotaScope `json:"scopes" bson:"scopes"`
	ScopeSelector *ScopeSelector       `json:"scope_selector" bson:"scope_selector"`
}

// ScopeSelector a scope selector represents the AND of the selectors represented
// by the scoped-resource selector requirements.
type ScopeSelector struct {
	// MatchExpressions a list of scope selector requirements by scope of the resources.
	MatchExpressions []ScopedResourceSelectorRequirement `json:"match_expressions" bson:"match_expressions"`
}

// ScopedResourceSelectorRequirement a scoped-resource selector requirement is a selector that
// contains values, a scope name, and an operator that relates the scope name and values.
type ScopedResourceSelectorRequirement struct {
	// ScopeName The name of the scope that the selector applies to.
	ScopeName ResourceQuotaScope `json:"scope_name" bson:"scope_name"`
	// Represents a scope's relationship to a set of values.
	// Valid operators are In, NotIn, Exists, DoesNotExist.
	Operator ScopeSelectorOperator `json:"operator" bson:"operator"`
	// Values An array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty.
	Values []string `json:"values" bson:"values"`
}

// NsUpdateReq update namespace request struct
type NsUpdateReq struct {
	Data []NsUpdateData `json:"data"`
}

// Validate validate namespace update request data
func (ns *NsUpdateReq) Validate() errors.RawErrorInfo {
	if len(ns.Data) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	sum := 0
	for _, data := range ns.Data {
		if err := data.Validate(); err.ErrCode != 0 {
			return err
		}

		sum += data.Count()
		if sum > NsUpdateLimit {
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommXXExceedLimit,
				Args:    []interface{}{"data", NsUpdateLimit},
			}
		}
	}

	return errors.RawErrorInfo{}
}

// NsUpdateData update namespace struct
type NsUpdateData struct {
	ID     []int64    `json:"id"`
	Unique []NsUnique `json:"unique"`
	Info   *Namespace `json:"info"`
}

// Count return namespace update data count
func (ns *NsUpdateData) Count() int {
	if len(ns.ID) != 0 {
		return len(ns.ID)
	}

	if len(ns.Unique) != 0 {
		return len(ns.Unique)
	}

	return 0
}

// Validate validate namespace update data
func (ns *NsUpdateData) Validate() errors.RawErrorInfo {
	if len(ns.ID) == 0 && len(ns.Unique) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"id and unique"},
		}
	}

	if len(ns.ID) != 0 && len(ns.Unique) != 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"id and unique"},
		}
	}

	if ns.Info == nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"info"},
		}
	}

	if err := ns.Info.ValidateUpdate(); err.ErrCode != 0 {
		return err
	}
	return errors.RawErrorInfo{}
}

// NsUnique namespace unique identification
type NsUnique struct {
	ClusterUID string `json:"cluster_uid" bson:"cluster_uid"`
	Name       string `json:"name" bson:"name"`
	ID         int64  `json:"id" bson:"id"`
}

// Validate validate NsUnique
func (ns *NsUnique) Validate() errors.RawErrorInfo {
	if ns.Name != "" && ns.ClusterUID != "" && ns.ID != 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"data"},
		}
	}

	if ns.Name == "" && ns.ClusterUID == "" && ns.ID == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"data"},
		}
	}

	if ns.ID == 0 && (ns.ClusterUID == "" || ns.Name == "") {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"data"},
		}
	}

	return errors.RawErrorInfo{}
}

// NsDeleteReq delete namespace request
type NsDeleteReq struct {
	Data []NsUnique `json:"data"`
}

// Validate validate NsDeleteReq
func (ns *NsDeleteReq) Validate() errors.RawErrorInfo {
	if len(ns.Data) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if len(ns.Data) > NsDeleteLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"data", NsDeleteLimit},
		}
	}

	for _, data := range ns.Data {
		if err := data.Validate(); err.ErrCode != 0 {
			return err
		}
	}

	return errors.RawErrorInfo{}
}

// NsCreateReq create namespace request
type NsCreateReq struct {
	Data []Namespace `json:"data"`
}

// Validate validate NsCreateReq
func (ns *NsCreateReq) Validate() errors.RawErrorInfo {
	if len(ns.Data) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if len(ns.Data) > NsCreateLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"data", NsCreateLimit},
		}
	}

	for _, data := range ns.Data {
		if err := data.ValidateCreate(); err.ErrCode != 0 {
			return err
		}
	}

	return errors.RawErrorInfo{}
}

// NsCreateResp create namespace response
type NsCreateResp struct {
	metadata.BaseResp `json:",inline"`
	Data              NsCreateRespData `json:"data"`
}

// NsCreateRespData create namespace response data
type NsCreateRespData struct {
	IDs []int64 `json:"ids"`
}

// NsQueryReq namespace query request
type NsQueryReq struct {
	ClusterSpec `json:",inline" bson:",inline"`
	Filter      *filter.Expression `json:"filter"`
	Fields      []string           `json:"fields,omitempty"`
	Page        metadata.BasePage  `json:"page,omitempty"`
}

// Validate validate NsQueryReq
func (ns *NsQueryReq) Validate() errors.RawErrorInfo {
	if ns.ClusterUID != nil && ns.ClusterID != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrorTopoIdentificationIllegal,
		}
	}

	if errInfo, err := ns.Page.Validate(false); err != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{errInfo},
		}
	}

	// todo validate Filter
	return errors.RawErrorInfo{}
}

// BuildCond build query namespace condition
func (ns *NsQueryReq) BuildCond(bizID int64, supplierAccount string) (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}
	if supplierAccount != "" {
		cond[common.BkSupplierAccount] = supplierAccount
	}
	if ns.ClusterID != nil {
		cond[BKClusterIDFiled] = ns.ClusterID
	}
	if ns.ClusterUID != nil {
		cond[ClusterUIDField] = ns.ClusterUID
	}

	if ns.Filter != nil {
		filterCond, err := ns.Filter.ToMgo()
		if err != nil {
			return nil, err
		}
		cond = mapstr.MapStr{common.BKDBAND: []mapstr.MapStr{cond, filterCond}}
	}
	return cond, nil
}
