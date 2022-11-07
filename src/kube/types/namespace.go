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
	"configcenter/src/common/util"
	"configcenter/src/filter"
	"configcenter/src/storage/dal/table"
)

// NamespaceFields merge the fields of the namespace and the details corresponding to the fields together.
var NamespaceFields = table.MergeFields(CommonSpecFieldsDescriptor, BizIDDescriptor, NamespaceSpecFieldsDescriptor)

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

const (
	// NsUpdateLimit limit on the number of namespace updates
	NsUpdateLimit = 200
	// NsDeleteLimit limit on the number of namespace delete
	NsDeleteLimit = 200
	// NsCreateLimit limit on the number of namespace create
	NsCreateLimit = 200
	// NsQueryLimit limit on the number of namespace query
	NsQueryLimit = 500
)

// Namespace define the namespace struct.
type Namespace struct {
	ClusterSpec     `json:",inline" bson:",inline"`
	ID              int64              `json:"id,omitempty" bson:"id"`
	Name            string             `json:"name,omitempty" bson:"name"`
	Labels          *map[string]string `json:"labels,omitempty" bson:"labels"`
	ResourceQuotas  *[]ResourceQuota   `json:"resource_quotas,omitempty" bson:"resource_quotas"`
	SupplierAccount string             `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
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

// NsUpdateReq update namespace request
type NsUpdateReq struct {
	IDs  []int64    `json:"ids"`
	Data *Namespace `json:"data"`
}

// Validate validate update namespace request
func (ns *NsUpdateReq) Validate() errors.RawErrorInfo {
	if len(ns.IDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"ids"},
		}
	}

	if len(ns.IDs) > NsUpdateLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", NsUpdateLimit},
		}
	}

	if ns.Data == nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"data"},
		}
	}

	if err := ValidateUpdate(*ns.Data, NamespaceFields); err.ErrCode != 0 {
		return err
	}
	return errors.RawErrorInfo{}
}

// NsDeleteReq delete namespace request
type NsDeleteReq struct {
	IDs []int64 `json:"ids"`
}

// Validate validate NsDeleteReq
func (ns *NsDeleteReq) Validate() errors.RawErrorInfo {
	if len(ns.IDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"ids"},
		}
	}

	if len(ns.IDs) > NsDeleteLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", NsDeleteLimit},
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
		if err := ValidateCreate(data, NamespaceFields); err.ErrCode != 0 {
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
	if ns.ClusterUID != "" && ns.ClusterID != 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrorTopoIdentificationIllegal,
		}
	}

	if err := ns.Page.ValidateWithEnableCount(false, NsQueryLimit); err.ErrCode != 0 {
		return err
	}

	op := filter.NewDefaultExprOpt(NamespaceFields.FieldsType())
	if err := ns.Filter.Validate(op); err != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{err.Error()},
		}
	}
	return errors.RawErrorInfo{}
}

// BuildCond build query namespace condition
func (ns *NsQueryReq) BuildCond(bizID int64, supplierAccount string) (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}
	cond = util.SetQueryOwner(cond, supplierAccount)

	if ns.ClusterID != 0 {
		cond[BKClusterIDFiled] = ns.ClusterID
	}
	if ns.ClusterUID != "" {
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

// NsInstResp namespace instance response
type NsInstResp struct {
	metadata.BaseResp `json:",inline"`
	Data              NsDataResp `json:"data"`
}

// NsDataResp namespace data
type NsDataResp struct {
	Data []Namespace `json:"data"`
}
