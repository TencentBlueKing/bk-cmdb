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
	"encoding/json"
	"fmt"

	"configcenter/pkg/filter"
	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/table"
)

// WorkLoadSpecFieldsDescriptor workLoad spec's fields descriptors.
// TODO remove this when kube attribute api supports workload types
var WorkLoadSpecFieldsDescriptor = table.FieldsDescriptors{
	{Field: KubeNameField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: NamespaceField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: LabelsField, Type: enumor.MapString, IsRequired: false, IsEditable: true},
	{Field: SelectorField, Type: enumor.Object, IsRequired: false, IsEditable: true},
	{Field: ReplicasField, Type: enumor.Numeric, IsRequired: true, IsEditable: true},
	{Field: StrategyTypeField, Type: enumor.String, IsRequired: false, IsEditable: true},
	{Field: MinReadySecondsField, Type: enumor.Numeric, IsRequired: false, IsEditable: true},
	{Field: RollingUpdateStrategyField, Type: enumor.Object, IsRequired: false, IsEditable: true},
}

// WorkLoadRefDescriptor  the description used when other resources refer to the workload.
var WorkLoadRefDescriptor = table.FieldsDescriptors{
	{Field: RefField, Type: enumor.Object, IsRequired: true, IsEditable: false},
	{Field: RefKindField, Type: enumor.String, IsRequired: true, IsEditable: false},
	{Field: RefIDField, Type: enumor.Numeric, IsRequired: false, IsEditable: false},
	{Field: RefNameField, Type: enumor.String, IsRequired: false, IsEditable: false},
}

// LabelSelectorOperator a label selector operator is the set of operators that can be used in a selector requirement.
type LabelSelectorOperator string

const (
	// LabelSelectorOpIn in operator for label selector
	LabelSelectorOpIn LabelSelectorOperator = "In"
	// LabelSelectorOpNotIn not in operator for label selector
	LabelSelectorOpNotIn LabelSelectorOperator = "NotIn"
	// LabelSelectorOpExists exists operator for label selector
	LabelSelectorOpExists LabelSelectorOperator = "Exists"
	// LabelSelectorOpDoesNotExist not exists operator for label selector
	LabelSelectorOpDoesNotExist LabelSelectorOperator = "DoesNotExist"
)

const (
	// WlUpdateLimit limit on the number of workload updates
	WlUpdateLimit = 200
	// WlDeleteLimit limit on the number of workload delete
	WlDeleteLimit = 200
	// WlCreateLimit limit on the number of workload create
	WlCreateLimit = 200
	// WlQueryLimit limit on the number of workload query
	WlQueryLimit = 500
)

// Type represents the stored type of IntOrString.
type Type int64

const (
	// IntType the IntOrString holds an int.
	IntType = 0
	// StringType the IntOrString holds a string.
	StringType = 1
)

// WorkloadInterface defines the workload data common operation.
type WorkloadInterface interface {
	ValidateCreate() errors.RawErrorInfo
	ValidateUpdate() errors.RawErrorInfo
	GetWorkloadBase() WorkloadBase
	SetWorkloadBase(wl WorkloadBase)
	BuildUpdateData(user string) (map[string]interface{}, error)
}

// Reference store pod-related workload related information
type Reference struct {
	// Kind workload kind
	Kind WorkloadType `json:"kind" bson:"kind"`

	// Name workload name
	Name string `json:"name" bson:"name"`

	// ID workload id in cc
	ID int64 `json:"id" bson:"id"`
}

// WorkloadSpec describes the common attributes of workload,
// it is used by the structure below it.
type WorkloadSpec struct {
	NamespaceSpec `json:",inline" bson:",inline"`
	Ref           Reference `json:"ref" bson:"ref"`
}

// WorkloadBase define the workload common struct, k8s workload attributes are placed in their respective structures,
// except for very public variables, please do not put them in.
type WorkloadBase struct {
	NamespaceSpec   `json:",inline" bson:",inline"`
	ID              int64  `json:"id,omitempty" bson:"id"`
	Name            string `json:"name,omitempty" bson:"name"`
	SupplierAccount string `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
	// Revision record this app's revision information
	table.Revision `json:",inline" bson:",inline"`
}

// LabelSelector a label selector is a label query over a set of resources.
// the result of matchLabels and matchExpressions are ANDed. An empty label
// selector matches all objects. A null label selector matches no objects.
type LabelSelector struct {
	// MatchLabels is a map of {key,value} pairs.
	MatchLabels map[string]string `json:"match_labels" bson:"match_labels"`
	// MatchExpressions is a list of label selector requirements. The requirements are ANDed.
	MatchExpressions []LabelSelectorRequirement `json:"match_expressions" bson:"match_expressions"`
}

// LabelSelectorRequirement a label selector requirement is a selector that contains values, a key,
// and an operator that relates the key and values.
type LabelSelectorRequirement struct {
	// key is the label key that the selector applies to.
	Key string `json:"key" bson:"key"`
	// operator represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator LabelSelectorOperator `json:"operator" bson:"operator"`
	// Values is an array of string values. If the operator is In or NotIn,
	// values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty.
	Values []string `json:"values" bson:"values"`
}

// IntOrString is a type that can hold an int32 or a string.
type IntOrString struct {
	Type   Type   `json:"type" bson:"type"`
	IntVal int32  `json:"int_val" bson:"int_val"`
	StrVal string `json:"str_val" bson:"str_val"`
}

type jsonWlData struct {
	IDs  []int64         `json:"ids"`
	Data json.RawMessage `json:"data"`
}

// WlUpdateOption defines the workload update request common operation.
type WlUpdateOption struct {
	Kind WorkloadType      `json:"kind"`
	IDs  []int64           `json:"ids"`
	Data WorkloadInterface `json:"data"`
}

// Validate validate WlCommonUpdate
func (w *WlUpdateOption) Validate() errors.RawErrorInfo {
	if len(w.IDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"ids"},
		}
	}

	if len(w.IDs) >= WlUpdateLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"data", WlUpdateLimit},
		}
	}

	if w.Data == nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"data"},
		}
	}

	if err := w.Data.ValidateUpdate(); err.ErrCode != 0 {
		return err
	}

	return errors.RawErrorInfo{}
}

// UnmarshalJSON unmarshal WlUpdateReq
func (w *WlUpdateOption) UnmarshalJSON(data []byte) error {
	kind := w.Kind
	var err error
	if err = kind.Validate(); err != nil {
		return err
	}

	req := new(jsonWlData)
	if err = json.Unmarshal(data, req); err != nil {
		return err
	}
	w.IDs = req.IDs

	if len(req.Data) == 0 {
		return nil
	}

	w.Data, err = kind.NewInst()
	if err != nil {
		return err
	}
	if err = json.Unmarshal(req.Data, &w.Data); err != nil {
		return err
	}

	return nil
}

// WlDeleteOption workload delete request
type WlDeleteOption struct {
	IDs []int64 `json:"ids"`
}

// Validate validate WlDeleteOption
func (ns *WlDeleteOption) Validate() errors.RawErrorInfo {
	if len(ns.IDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"ids"},
		}
	}

	if len(ns.IDs) > WlDeleteLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"ids", WlDeleteLimit},
		}
	}

	return errors.RawErrorInfo{}
}

// WlDataResp workload data
type WlDataResp struct {
	Kind WorkloadType        `json:"kind"`
	Info []WorkloadInterface `json:"info"`
}

type jsonWlInfo struct {
	Info json.RawMessage `json:"info"`
}

// UnmarshalJSON unmarshal WlDataResp
// NOCC:golint/fnsize(workload类型会不断增多)
func (w *WlDataResp) UnmarshalJSON(data []byte) error {
	kind := w.Kind
	wlData := new(jsonWlInfo)
	if err := json.Unmarshal(data, wlData); err != nil {
		return err
	}

	if err := kind.Validate(); err != nil {
		return err
	}

	if wlData.Info == nil {
		return nil
	}

	switch kind {
	case KubeDeployment:
		array := make([]*Deployment, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubeStatefulSet:
		array := make([]*StatefulSet, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubeDaemonSet:
		array := make([]*DaemonSet, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubeGameDeployment:
		array := make([]*GameDeployment, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubeGameStatefulSet:
		array := make([]*GameStatefulSet, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubeCronJob:
		array := make([]*CronJob, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubeJob:
		array := make([]*Job, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	case KubePodWorkload:
		array := make([]*PodsWorkload, 0)
		if err := json.Unmarshal(wlData.Info, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Info = append(w.Info, data)
		}

	default:
		return fmt.Errorf("can not support this workload type: %v", kind)
	}
	return nil
}

// WlInstResp workload instance response
type WlInstResp struct {
	metadata.BaseResp `json:",inline"`
	Data              WlDataResp `json:"data"`
}

// WlCreateOption create workload request
type WlCreateOption struct {
	Kind WorkloadType        `json:"kind"`
	Data []WorkloadInterface `json:"data"`
}

// UnmarshalJSON unmarshal WlCreateOption
// NOCC:golint/fnsize(workload类型会不断增多)
func (w *WlCreateOption) UnmarshalJSON(data []byte) error {
	kind := w.Kind
	req := new(jsonWlData)
	if err := json.Unmarshal(data, req); err != nil {
		return err
	}

	if len(req.Data) == 0 {
		return nil
	}

	if err := kind.Validate(); err != nil {
		return err
	}

	switch kind {
	case KubeDeployment:
		array := make([]*Deployment, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeStatefulSet:
		array := make([]*StatefulSet, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeDaemonSet:
		array := make([]*DaemonSet, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeGameDeployment:
		array := make([]*GameDeployment, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeGameStatefulSet:
		array := make([]*GameStatefulSet, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeCronJob:
		array := make([]*CronJob, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeJob:
		array := make([]*Job, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubePodWorkload:
		array := make([]*PodsWorkload, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	default:
		return fmt.Errorf("can not support this workload type: %v", kind)
	}
	return nil
}

// Validate validate WlCreateOption
func (ns *WlCreateOption) Validate() errors.RawErrorInfo {
	if len(ns.Data) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if len(ns.Data) > WlCreateLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"data", WlCreateLimit},
		}
	}

	for _, data := range ns.Data {
		if err := data.ValidateCreate(); err.ErrCode != 0 {
			return err
		}
	}

	return errors.RawErrorInfo{}
}

// WlCreateResp create workload response
type WlCreateResp struct {
	metadata.BaseResp `json:",inline"`
	Data              metadata.RspIDs `json:"data"`
}

var wlIgnoreField = []string{
	common.BKAppIDField, BKClusterIDField, ClusterUIDField, BKNamespaceIDField, NamespaceField, common.BKFieldName,
	common.BKFieldID, common.CreateTimeField,
}

// WlQueryOption workload query request
type WlQueryOption struct {
	Filter *filter.Expression `json:"filter"`
	Fields []string           `json:"fields,omitempty"`
	Page   metadata.BasePage  `json:"page,omitempty"`
}

// Validate validate WlQueryReq
func (wl *WlQueryOption) Validate(kind WorkloadType) errors.RawErrorInfo {
	if err := wl.Page.ValidateWithEnableCount(false, WlQueryLimit); err.ErrCode != 0 {
		return err
	}

	fields, err := kind.Fields()
	if err != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{KindField},
		}
	}

	if wl.Filter == nil {
		return errors.RawErrorInfo{}
	}

	op := filter.NewDefaultExprOpt(fields.FieldsType())
	if err := wl.Filter.Validate(op); err != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{err.Error()},
		}
	}
	return errors.RawErrorInfo{}
}

// BuildCond build query workload condition
func (wl *WlQueryOption) BuildCond(bizID int64) (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}
	if wl.Filter != nil {
		filterCond, err := wl.Filter.ToMgo()
		if err != nil {
			return nil, err
		}
		cond = mapstr.MapStr{common.BKDBAND: []mapstr.MapStr{cond, filterCond}}
	}
	return cond, nil
}
