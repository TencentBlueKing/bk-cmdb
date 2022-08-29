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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/criteria/enumor"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/filter"
	"configcenter/src/kube/orm"
	"configcenter/src/storage/dal/table"
)

// WorkLoadSpecFieldsDescriptor workLoad spec's fields descriptors.
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

// WorkloadI defines the workload data common operation.
type WorkloadI interface {
	ValidateCreate() errors.RawErrorInfo
	ValidateUpdate() errors.RawErrorInfo
	SetID(id int64)
	SetCreateTime(createTime int64)
	SetUpdateTime(updateTime int64)
	GetNamespaceSpec() NamespaceSpec
	SetNamespaceSpec(spec NamespaceSpec)
	SetSupplierAccount(supplierAccount string)
}

// Workload define the workload common struct.
type Workload struct {
	NamespaceSpec   `json:",inline" bson:",inline"`
	ID              *int64             `json:"id" bson:"id"`
	Name            *string            `json:"name" bson:"name"`
	Labels          *map[string]string `json:"labels" bson:"labels"`
	Selector        *LabelSelector     `json:"selector" bson:"selector"`
	Replicas        *int64             `json:"replicas" bson:"replicas"`
	MinReadySeconds *int64             `json:"min_ready_seconds" bson:"min_ready_seconds"`
	CreateTime      *int64             `json:"create_time" bson:"create_time"`
	UpdateTime      *int64             `json:"update_time" bson:"update_time"`
	SupplierAccount *string            `json:"bk_supplier_account" bson:"bk_supplier_account"`
}

// ValidateCreate validate create workload
func (w *Workload) ValidateCreate() errors.RawErrorInfo {
	if w.NamespaceID == nil && (w.ClusterUID == nil || w.Namespace == nil) {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{BKNamespaceIDField + " or <" + ClusterUIDField + " and " + NamespaceField + " >"},
		}
	}

	if w.Name == nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKFieldName},
		}
	}

	return errors.RawErrorInfo{}
}

// ValidateUpdate validate update workload
func (w *Workload) ValidateUpdate() errors.RawErrorInfo {
	// todo
	return errors.RawErrorInfo{}
}

// SetID set id
func (w *Workload) SetID(id int64) {
	w.ID = &id
}

// SetCreateTime set create time
func (w *Workload) SetCreateTime(createTime int64) {
	w.CreateTime = &createTime
}

// SetUpdateTime set update time
func (w *Workload) SetUpdateTime(updateTime int64) {
	w.UpdateTime = &updateTime
}

// SetSupplierAccount set supplierAccount
func (w *Workload) SetSupplierAccount(supplierAccount string) {
	w.SupplierAccount = &supplierAccount
}

// GetNamespaceSpec get namespace spec
func (w *Workload) GetNamespaceSpec() NamespaceSpec {
	return w.NamespaceSpec
}

// SetNamespaceSpec set namespace spec
func (w *Workload) SetNamespaceSpec(spec NamespaceSpec) {
	w.NamespaceSpec = spec
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

// WlIdentification workload unique identity to identify workload
type WlIdentification struct {
	ID     []int64    `json:"id"`
	Unique []WlUnique `json:"unique"`
}

// Count return namespace update data count
func (w *WlIdentification) Count() int {
	if len(w.ID) != 0 {
		return len(w.ID)
	}

	if len(w.Unique) != 0 {
		return len(w.Unique)
	}

	return 0
}

// Validate validate WlIdentification
func (w *WlIdentification) Validate() errors.RawErrorInfo {
	if len(w.ID) == 0 && len(w.Unique) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"id and unique"},
		}
	}

	if len(w.ID) != 0 && len(w.Unique) != 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"id and unique"},
		}
	}

	return errors.RawErrorInfo{}
}

// BuildUpdateFilter build update filter
func (w *WlIdentification) BuildUpdateFilter(bizID int64, supplierAccount string) map[string]interface{} {
	var filter map[string]interface{}
	if len(w.ID) != 0 {
		filter = map[string]interface{}{
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: supplierAccount,
			common.BKFieldID: map[string]interface{}{
				common.BKDBIN: w.ID,
			},
		}
	}

	if len(w.Unique) != 0 {
		orCond := make([]map[string]interface{}, 0)
		for _, unique := range w.Unique {
			cond := map[string]interface{}{
				ClusterUIDField:    unique.ClusterUID,
				NamespaceField:     unique.Namespace,
				common.BKFieldName: unique.Name,
			}
			orCond = append(orCond, cond)
		}
		filter = map[string]interface{}{
			common.BKAppIDField:   bizID,
			common.BKOwnerIDField: supplierAccount,
			common.BKDBOR:         orCond,
		}
	}
	return filter
}

// WlUnique workload unique identification
type WlUnique struct {
	ClusterUID string `json:"cluster_uid" bson:"cluster_uid"`
	Namespace  string `json:"namespace" bson:"namespace"`
	Name       string `json:"name" bson:"name"`
	ID         int64  `json:"id" bson:"id"`
}

// Validate validate WlUnique
func (wl *WlUnique) Validate() errors.RawErrorInfo {
	if wl.Name != "" && wl.ClusterUID != "" && wl.Namespace != "" && wl.ID != 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"data"},
		}
	}

	if wl.Name == "" && wl.ClusterUID == "" && wl.Namespace != "" && wl.ID == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"data"},
		}
	}

	if wl.ID == 0 && (wl.Namespace == "" || wl.ClusterUID == "" || wl.Name == "") {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"data"},
		}
	}

	return errors.RawErrorInfo{}
}

type jsonWlUpdateReq struct {
	Data json.RawMessage `json:"data"`
}

// WlUpdateReq defines the workload update request common operation.
type WlUpdateReq struct {
	Kind WorkloadType    `json:"kind"`
	Data []WlUpdateDataI `json:"data"`
}

// UnmarshalJSON unmarshal WlUpdateReq
func (w *WlUpdateReq) UnmarshalJSON(data []byte) error {
	kind := w.Kind
	req := jsonWlUpdateReq{}
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	if req.Data == nil || !IsInnerWorkload(kind) {
		return nil
	}

	switch kind {
	case KubeDeployment:
		array := make([]*DeployUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeStatefulSet:
		array := make([]*StatefulSetUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeDaemonSet:
		array := make([]*DaemonSetUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeGameDeployment:
		array := make([]*GameDeployUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	case KubeGameStatefulSet:
		array := make([]*GameStatefulSetUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}

	default:
		array := make([]*WlUpdateData, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}
	}
	return nil
}

// Validate validate workload update request data
func (d *WlUpdateReq) Validate() errors.RawErrorInfo {
	if len(d.Data) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	sum := 0
	for _, data := range d.Data {
		if err := data.Validate(); err.ErrCode != 0 {
			return err
		}

		sum += data.Count()
		if sum > WlUpdateLimit {
			return errors.RawErrorInfo{
				ErrCode: common.CCErrCommXXExceedLimit,
				Args:    []interface{}{"data", WlUpdateLimit},
			}
		}
	}

	return errors.RawErrorInfo{}
}

// WlUpdateDataI defines the workload update data common operation.
type WlUpdateDataI interface {
	Validate() errors.RawErrorInfo
	Count() int
	BuildUpdateFilter(bizID int64, supplierAccount string) map[string]interface{}
	BuildUpdateData() (map[string]interface{}, error)
}

// WlUpdateData defines the workload update data common operation.
type WlUpdateData struct {
	WlIdentification `json:",inline"`
	Info             Workload `json:"info"`
}

// BuildUpdateData build workload update data
func (w *WlUpdateData) BuildUpdateData() (map[string]interface{}, error) {
	now := time.Now().Unix()
	w.Info.UpdateTime = &now
	opts := orm.NewFieldOptions().AddIgnoredFields(wlIgnoreField...)
	updateData, err := orm.GetUpdateFieldsWithOption(w.Info, opts)
	if err != nil {
		return nil, err
	}
	return updateData, err
}

// WlDeleteReq workload delete request
type WlDeleteReq struct {
	Data []WlUnique `json:"data"`
}

// Validate validate WlDeleteReq
func (ns *WlDeleteReq) Validate() errors.RawErrorInfo {
	if len(ns.Data) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if len(ns.Data) > WlDeleteLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"data", WlDeleteLimit},
		}
	}

	for _, data := range ns.Data {
		if err := data.Validate(); err.ErrCode != 0 {
			return err
		}
	}

	return errors.RawErrorInfo{}
}

type jsonWlCreateReq struct {
	Data json.RawMessage `json:"data"`
}

// WlCreateReq create workload request
type WlCreateReq struct {
	Kind WorkloadType `json:"kind"`
	Data []WorkloadI  `json:"data"`
}

// UnmarshalJSON unmarshal WlUpdateReq
func (w *WlCreateReq) UnmarshalJSON(data []byte) error {
	kind := w.Kind
	req := jsonWlCreateReq{}
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	if req.Data == nil || !IsInnerWorkload(kind) {
		return nil
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

	default:
		array := make([]*Workload, 0)
		if err := json.Unmarshal(req.Data, &array); err != nil {
			return err
		}
		for _, data := range array {
			w.Data = append(w.Data, data)
		}
	}
	return nil
}

// Validate validate WlCreateReq
func (ns *WlCreateReq) Validate() errors.RawErrorInfo {
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
	Data              WlCreateRespData `json:"data"`
}

// WlCreateRespData create workload response data
type WlCreateRespData struct {
	IDs []int64 `json:"ids"`
}

var wlIgnoreField = []string{
	common.BKAppIDField, BKClusterIDFiled, ClusterUIDField, BKNamespaceIDField, NamespaceField, common.BKFieldName,
	common.BKFieldID,
}

// IsInnerWorkload is inner workload type
func IsInnerWorkload(kind WorkloadType) bool {
	switch kind {
	case KubeDeployment, KubeStatefulSet, KubeDaemonSet,
		KubeGameStatefulSet, KubeGameDeployment, KubeCronJob,
		KubeJob, KubePodWorkload:
		return true
	default:
		return false
	}
}

// GetWorkloadTableName get workload table name
func GetWorkloadTableName(kind WorkloadType) (string, error) {
	switch kind {
	case KubeDeployment:
		return BKTableNameBaseDeployment, nil

	case KubeStatefulSet:
		return BKTableNameBaseStatefulSet, nil

	case KubeDaemonSet:
		return BKTableNameBaseDaemonSet, nil

	case KubeGameStatefulSet:
		return BKTableNameGameStatefulSet, nil

	case KubeGameDeployment:
		return BKTableNameGameDeployment, nil

	case KubeCronJob:
		return BKTableNameBaseCronJob, nil

	case KubeJob:
		return BKTableNameBaseJob, nil

	case KubePodWorkload:
		return BKTableNameBasePodWorkload, nil

	default:
		return "", fmt.Errorf("can not find table name, kind: %s", kind)
	}
}

// WlQueryReq workload query request
type WlQueryReq struct {
	NamespaceSpec `json:",inline" bson:",inline"`
	Filter        *filter.Expression `json:"filter"`
	Fields        []string           `json:"fields,omitempty"`
	Page          metadata.BasePage  `json:"page,omitempty"`
}

// Validate validate WlQueryReq
func (wl *WlQueryReq) Validate() errors.RawErrorInfo {
	if (wl.ClusterID != nil || wl.NamespaceID != nil) && (wl.ClusterUID != nil && wl.Namespace != nil) {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrorTopoIdentificationIllegal,
		}
	}

	if err := wl.Page.ValidateWithEnableCount(false, WlQueryLimit); err.ErrCode != 0 {
		return err
	}

	// todo validate Filter
	return errors.RawErrorInfo{}
}

// BuildCond build query workload condition
func (wl *WlQueryReq) BuildCond(bizID int64, supplierAccount string) (mapstr.MapStr, error) {
	cond := mapstr.MapStr{
		common.BKAppIDField:      bizID,
		common.BkSupplierAccount: supplierAccount,
	}

	if wl.ClusterID != nil {
		cond[BKClusterIDFiled] = wl.ClusterID
	}

	if wl.ClusterUID != nil {
		cond[ClusterUIDField] = wl.ClusterUID
	}

	if wl.NamespaceID != nil {
		cond[BKNamespaceIDField] = wl.NamespaceID
	}

	if wl.Namespace != nil {
		cond[NamespaceField] = wl.Namespace
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
