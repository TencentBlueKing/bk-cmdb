/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	"encoding/json"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/querybuilder"

	"go.mongodb.org/mongo-driver/bson"
)

// AuditQueryResponse query audit logs response
type AuditQueryResponse struct {
	BaseResp `json:",inline"`
	Data     *AuditQueryResult `json:"data"`
}

// AuditQueryResult query audit logs result
type AuditQueryResult struct {
	Count int64      `json:"count"`
	Info  []AuditLog `json:"info"`
}

// CreateAuditLogParam TODO
type CreateAuditLogParam struct {
	Data []AuditLog `json:"data"`
}

// AuditQueryInput TODO
type AuditQueryInput struct {
	Condition AuditQueryCondition `json:"condition"`
	Page      BasePage            `json:"page,omitempty"`
}

// Validate validates the input param
func (input *AuditQueryInput) Validate() errors.RawErrorInfo {
	if input.Page.Limit <= 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"limit"},
		}
	}

	if input.Page.Limit > common.BKAuditLogPageLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommPageLimitIsExceeded,
		}
	}

	if len(input.Condition.OperationTime.Start) == 0 && len(input.Condition.OperationTime.End) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKOperationTimeField},
		}
	}

	return errors.RawErrorInfo{}
}

// AuditQueryCondition TODO
type AuditQueryCondition struct {
	AuditType    AuditType       `json:"audit_type"`
	User         string          `json:"user"`
	ResourceType ResourceType    `json:"resource_type" `
	Action       []ActionType    `json:"action"`
	OperateFrom  OperateFromType `json:"operate_from"`
	BizID        int64           `json:"bk_biz_id"`
	ResourceID   interface{}     `json:"resource_id"`
	// ResourceName filters audit logs by resource name, such as instance name, host ip etc., support fuzzy query
	ResourceName string `json:"resource_name"`
	// OperationTime is an array of start time and end time, filters audit logs between them
	OperationTime OperationTimeCondition `json:"operation_time"`
	// Category is used by front end, filters audit logs as business(business resource related to business), resource(instance resource not related to business), host or other category
	Category string `json:"category"`
	// ObjID is used for instance and instance association audit log filter like host deletion history
	ObjID string `json:"bk_obj_id"`
	// FuzzyQuery is used for searching for resource name using regex, use accurate query by default when it is not set
	FuzzyQuery bool `json:"fuzzy_query"`
	// Condition is used for new way to search audit log by user or resource_name
	Condition []querybuilder.AtomRule `json:"condition"`
}

// Validate is a AuditQueryCondition validator to validate user resource_name condition whether exist at the same time
func (a *AuditQueryCondition) Validate() error {
	if (len(a.User) != 0 || len(a.ResourceName) != 0) && len(a.Condition) != 0 {
		return errors.New(common.CCErrCommParamsInvalid, "condition, user and resource_name cannot exist at the same time")
	}
	return nil
}

// OperationTimeCondition TODO
type OperationTimeCondition struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// InstAuditQueryInput search instance audit query condition
type InstAuditQueryInput struct {
	Condition  InstAuditCondition `json:"condition"`
	Page       BasePage           `json:"page,omitempty"`
	WithDetail bool               `json:"with_detail"`
}

// Validate validates the input param
func (input *InstAuditQueryInput) Validate() errors.RawErrorInfo {

	if input.Page.Limit > common.BKAuditLogPageLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommPageLimitIsExceeded,
		}
	}

	if len(input.Condition.ResourceType) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{common.BKResourceTypeField},
		}
	}

	if len(input.Condition.ObjID) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{common.BKObjIDField},
		}
	}

	if IsCommon(input.Condition.ObjID) {
		return errors.RawErrorInfo{}
	}

	// 前端目前只允许查看主机、业务、自定义模型的变更记录，因此在此限制objid不为host和biz时报错
	// front-end only allow to see change record of host, biz, custom object
	if input.Condition.ObjID != common.BKInnerObjIDApp && input.Condition.ObjID != common.BKInnerObjIDHost &&
		input.Condition.ObjID != common.BKInnerObjIDBizSet {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{common.BKObjIDField},
		}
	}

	return errors.RawErrorInfo{}
}

// InstAuditCondition instance audit condition
type InstAuditCondition struct {
	User          string                 `json:"user"`
	BizID         int64                  `json:"bk_biz_id"`
	ObjID         string                 `json:"bk_obj_id"`
	ResourceName  string                 `json:"resource_name"`
	ResourceID    interface{}            `json:"resource_id"`
	ResourceType  ResourceType           `json:"resource_type" `
	Action        []ActionType           `json:"action"`
	OperationTime OperationTimeCondition `json:"operation_time"`
	// ID is an audit record's id
	ID []int64 `json:"id"`
}

// AuditLog TODO
type AuditLog struct {
	ID int64 `json:"id" bson:"id"`
	// AuditType is a high level abstract of the resource managed by this cmdb.
	// Each kind of concept, resource must belongs to one of the resource type.
	AuditType AuditType `json:"audit_type" bson:"audit_type"`
	// the supplier account that this resource belongs to.
	SupplierAccount string `json:"bk_supplier_account" bson:"bk_supplier_account"`
	// name of the one who triggered this operation.
	User string `json:"user" bson:"user"`
	// the operated resource by the user
	ResourceType ResourceType `json:"resource_type" bson:"resource_type"`
	// ActionType represent the user's operation type, like CUD etc.
	Action ActionType `json:"action" bson:"action"`
	// OperateFrom describe which form does this audit come from.
	OperateFrom OperateFromType `json:"operate_from" bson:"operate_from"`
	// OperationDetail describe the details information by a user.
	// Note: when the ResourceType relevant to Business, then the business id field must
	// be bk_biz_id, otherwise the user can not search this operation log with business id.
	OperationDetail DetailFactory `json:"operation_detail" bson:"operation_detail"`
	// OperationTime is the time that user do the operation.
	OperationTime Time `json:"operation_time" bson:"operation_time"`
	// the business id of the resource if it belongs to a business.
	BusinessID int64 `json:"bk_biz_id,omitempty" bson:"bk_biz_id,omitempty"`
	// ResourceID is the id of the resource instance. which is a unique id, dynamic grouping id is string type.
	// for service instance audit log,
	ResourceID interface{} `json:"resource_id" bson:"resource_id"`
	// ResourceName is the name of the resource, such as a switch model has a name "switch"
	ResourceName string `json:"resource_name" bson:"resource_name"`
	// AppCode is the app code of the system where the request comes from
	AppCode string `json:"code,omitempty" bson:"code,omitempty"`
	// RequestID is the request id of the request
	RequestID string `json:"rid,omitempty" bson:"rid,omitempty"`
}

type bsonAuditLog struct {
	ID              int64           `json:"id" bson:"id"`
	AuditType       AuditType       `json:"audit_type" bson:"audit_type"`
	SupplierAccount string          `json:"bk_supplier_account" bson:"bk_supplier_account"`
	User            string          `json:"user" bson:"user"`
	ResourceType    ResourceType    `json:"resource_type" bson:"resource_type"`
	Action          ActionType      `json:"action" bson:"action"`
	OperateFrom     OperateFromType `json:"operate_from" bson:"operate_from"`
	OperationTime   Time            `json:"operation_time" bson:"operation_time"`
	OperationDetail bson.Raw        `json:"operation_detail" bson:"operation_detail"`
	BusinessID      int64           `json:"bk_biz_id" bson:"bk_biz_id"`
	ResourceID      interface{}     `json:"resource_id" bson:"resource_id"`
	ResourceName    string          `json:"resource_name" bson:"resource_name"`
	AppCode         string          `json:"code" bson:"code"`
	RequestID       string          `json:"rid" bson:"rid"`
}

type jsonAuditLog struct {
	ID              int64           `json:"id" bson:"id"`
	AuditType       AuditType       `json:"audit_type" bson:"audit_type"`
	SupplierAccount string          `json:"bk_supplier_account" bson:"bk_supplier_account"`
	User            string          `json:"user" bson:"user"`
	ResourceType    ResourceType    `json:"resource_type" bson:"resource_type"`
	Action          ActionType      `json:"action" bson:"action"`
	OperateFrom     OperateFromType `json:"operate_from" bson:"operate_from"`
	OperationTime   Time            `json:"operation_time" bson:"operation_time"`
	OperationDetail json.RawMessage `json:"operation_detail" bson:"operation_detail"`
	BusinessID      int64           `json:"bk_biz_id" bson:"bk_biz_id"`
	ResourceID      interface{}     `json:"resource_id" bson:"resource_id"`
	ResourceName    string          `json:"resource_name" bson:"resource_name"`
	AppCode         string          `json:"code" bson:"code"`
	RequestID       string          `json:"rid" bson:"rid"`
}

// DetailFactory TODO
type DetailFactory interface {
	WithName() string
}

// UnmarshalJSON TODO
func (auditLog *AuditLog) UnmarshalJSON(data []byte) error {
	audit := jsonAuditLog{}
	if err := json.Unmarshal(data, &audit); err != nil {
		return err
	}
	auditLog.ID = audit.ID
	auditLog.AuditType = audit.AuditType
	auditLog.SupplierAccount = audit.SupplierAccount
	auditLog.User = audit.User
	auditLog.ResourceType = audit.ResourceType
	auditLog.Action = audit.Action
	auditLog.OperateFrom = audit.OperateFrom
	auditLog.OperationTime = audit.OperationTime
	auditLog.BusinessID = audit.BusinessID
	auditLog.ResourceID = audit.ResourceID
	auditLog.ResourceName = audit.ResourceName
	auditLog.AppCode = audit.AppCode
	auditLog.RequestID = audit.RequestID

	if audit.OperationDetail == nil {
		return nil
	}

	if audit.Action == AuditTransferHostModule || audit.Action == AuditAssignHost || audit.Action == AuditUnassignHost {
		operationDetail := new(HostTransferOpDetail)
		if err := json.Unmarshal(audit.OperationDetail, &operationDetail); err != nil {
			return err
		}
		auditLog.OperationDetail = operationDetail
		return nil
	}

	switch audit.ResourceType {
	case BusinessRes, BizSetRes, SetRes, ModuleRes, ProcessRes, HostRes, CloudAreaRes, ModelInstanceRes,
		MainlineInstanceRes, ResourceDirRes:
		operationDetail := new(InstanceOpDetail)
		if err := json.Unmarshal(audit.OperationDetail, &operationDetail); err != nil {
			return err
		}
		auditLog.OperationDetail = operationDetail
	case InstanceAssociationRes:
		operationDetail := new(InstanceAssociationOpDetail)
		if err := json.Unmarshal(audit.OperationDetail, &operationDetail); err != nil {
			return err
		}
		auditLog.OperationDetail = operationDetail
	case ModelAttributeRes, ModelAttributeGroupRes:
		operationDetail := new(ModelAttrOpDetail)
		if err := json.Unmarshal(audit.OperationDetail, &operationDetail); err != nil {
			return err
		}
		auditLog.OperationDetail = operationDetail
	case ServiceInstanceRes:
		operationDetail := new(ServiceInstanceOpDetail)
		if err := json.Unmarshal(audit.OperationDetail, &operationDetail); err != nil {
			return err
		}
		auditLog.OperationDetail = operationDetail
	default:
		operationDetail := new(BasicOpDetail)
		if err := json.Unmarshal(audit.OperationDetail, &operationDetail); err != nil {
			return err
		}
		auditLog.OperationDetail = operationDetail
	}
	return nil
}

// UnmarshalBSON TODO
func (auditLog *AuditLog) UnmarshalBSON(data []byte) error {
	audit := bsonAuditLog{}
	if err := bson.Unmarshal(data, &audit); err != nil {
		return err
	}
	auditLog.ID = audit.ID
	auditLog.AuditType = audit.AuditType
	auditLog.SupplierAccount = audit.SupplierAccount
	auditLog.User = audit.User
	auditLog.ResourceType = audit.ResourceType
	auditLog.Action = audit.Action
	auditLog.OperateFrom = audit.OperateFrom
	auditLog.OperationTime = audit.OperationTime
	auditLog.BusinessID = audit.BusinessID
	auditLog.ResourceID = audit.ResourceID
	auditLog.ResourceName = audit.ResourceName
	auditLog.AppCode = audit.AppCode
	auditLog.RequestID = audit.RequestID

	if audit.OperationDetail == nil {
		return nil
	}

	if audit.Action == AuditTransferHostModule || audit.Action == AuditAssignHost || audit.Action == AuditUnassignHost {
		operationDetail := new(HostTransferOpDetail)
		if err := bson.Unmarshal(audit.OperationDetail, &operationDetail); err != nil {
			return err
		}
		auditLog.OperationDetail = operationDetail
		return nil
	}

	switch audit.ResourceType {
	case BusinessRes, BizSetRes, SetRes, ModuleRes, ProcessRes, HostRes, CloudAreaRes, ModelInstanceRes,
		MainlineInstanceRes, ResourceDirRes:
		operationDetail := new(InstanceOpDetail)
		if err := bson.Unmarshal(audit.OperationDetail, &operationDetail); err != nil {
			return err
		}
		auditLog.OperationDetail = operationDetail
	case InstanceAssociationRes:
		operationDetail := new(InstanceAssociationOpDetail)
		if err := bson.Unmarshal(audit.OperationDetail, &operationDetail); err != nil {
			return err
		}
		auditLog.OperationDetail = operationDetail
	case ModelAttributeRes, ModelAttributeGroupRes:
		operationDetail := new(ModelAttrOpDetail)
		if err := bson.Unmarshal(audit.OperationDetail, &operationDetail); err != nil {
			return err
		}
		auditLog.OperationDetail = operationDetail
	case ServiceInstanceRes:
		operationDetail := new(ServiceInstanceOpDetail)
		if err := bson.Unmarshal(audit.OperationDetail, &operationDetail); err != nil {
			return err
		}
		auditLog.OperationDetail = operationDetail
	default:
		operationDetail := new(BasicOpDetail)
		if err := bson.Unmarshal(audit.OperationDetail, &operationDetail); err != nil {
			return err
		}
		auditLog.OperationDetail = operationDetail
	}
	return nil
}

// MarshalBSON TODO
func (auditLog AuditLog) MarshalBSON() ([]byte, error) {
	audit := bsonAuditLog{}
	audit.ID = auditLog.ID
	audit.AuditType = auditLog.AuditType
	audit.SupplierAccount = auditLog.SupplierAccount
	audit.User = auditLog.User
	audit.ResourceType = auditLog.ResourceType
	audit.Action = auditLog.Action
	audit.OperateFrom = auditLog.OperateFrom
	audit.OperationTime = auditLog.OperationTime
	audit.BusinessID = auditLog.BusinessID
	audit.ResourceID = auditLog.ResourceID
	audit.ResourceName = auditLog.ResourceName
	audit.AppCode = auditLog.AppCode
	audit.RequestID = auditLog.RequestID
	var err error
	switch val := auditLog.OperationDetail.(type) {
	default:
		audit.OperationDetail, err = bson.Marshal(val)
		if err != nil {
			return []byte{}, err
		}
	}
	return bson.Marshal(audit)
}

// BasicOpDetail TODO
type BasicOpDetail struct {
	// Details contains all the details information about a user's operation
	Details *BasicContent `json:"details" bson:"details"`
}

// WithName TODO
func (op *BasicOpDetail) WithName() string {
	return "BasicDetail"
}

// ModelAttrOpDetail TODO
type ModelAttrOpDetail struct {
	BasicOpDetail `bson:",inline"`
	// BkObjID the attribute object ID
	BkObjID string `json:"bk_obj_id" bson:"bk_obj_id"`
	// BkObjName the attribute object name
	BkObjName string `json:"bk_obj_name" bson:"bk_obj_name"`
}

// WithName TODO
func (op *ModelAttrOpDetail) WithName() string {
	return "ModelAttrDetail"
}

// InstanceOpDetail TODO
type InstanceOpDetail struct {
	BasicOpDetail `bson:",inline"`
	// BkObjID the object ID of the instance's model
	ModelID string `json:"bk_obj_id" bson:"bk_obj_id"`
}

// WithName TODO
func (op *InstanceOpDetail) WithName() string {
	return "InstanceOpDetail"
}

// HostTransferOpDetail TODO
type HostTransferOpDetail struct {
	// PreData the previous biz topology of the host before transfer
	PreData HostBizTopo `json:"pre_data" bson:"pre_data"`
	// CurData the current biz topology of the host before transfer
	CurData HostBizTopo `json:"cur_data" bson:"cur_data"`
}

// HostBizTopo TODO
type HostBizTopo struct {
	BizID   int64  `json:"bk_biz_id" bson:"bk_biz_id"`
	BizName string `json:"bk_biz_name" bson:"bk_biz_name"`
	Set     []Topo `json:"set" bson:"set"`
}

// WithName TODO
func (op *HostTransferOpDetail) WithName() string {
	return "HostTransferOpDetail"
}

// AssociationOpDetail TODO
type AssociationOpDetail struct {
	// AssociationID association id between two object
	AssociationID string `json:"asst_id" bson:"asst_id"`
	// AssociationKind association kind id
	AssociationKind string `json:"asst_kind" bson:"asst_kind"`
}

// InstanceAssociationOpDetail TODO
type InstanceAssociationOpDetail struct {
	AssociationOpDetail `bson:",inline"`
	// SourceModelID the source instance's object ID
	SourceModelID string `json:"src_obj_id" bson:"src_obj_id"`
	// TargetModelID the target instance's object ID
	TargetModelID string `json:"dest_obj_id" bson:"dest_obj_id"`
	// TargetInstanceID the target instance ID
	TargetInstanceID int64 `json:"dest_inst_id" bson:"dest_inst_id"`
	// TargetInstanceID the target instance name
	TargetInstanceName string `json:"dest_inst_name" bson:"dest_inst_name"`
}

// WithName TODO
func (ao *InstanceAssociationOpDetail) WithName() string {
	return "InstanceAssociationOpDetail"
}

// ModelAssociationOpDetail TODO
type ModelAssociationOpDetail struct {
	AssociationOpDetail `bson:",inline"`
	// AssociationName the alias name of the association
	AssociationName string `json:"asst_name" bson:"asst_name"`
	// Mapping the association mapping, defines which kind of association can be used
	Mapping AssociationMapping `json:"mapping" bson:"mapping"`
	// OnDelete the association on delete action
	OnDelete AssociationOnDeleteAction `json:"on_delete" bson:"on_delete"`
	// IsPre describe whether this association is a pre-defined association or not
	IsPre *bool `json:"is_pre" bson:"is_pre"`
	// TargetModelID the target object ID
	TargetModelID string `json:"dest_obj_id" bson:"dest_obj_id"`
	// TargetModelID the target object name
	TargetModelName int64 `json:"dest_obj_name" bson:"dest_obj_name"`
}

// WithName TODO
func (ao *ModelAssociationOpDetail) WithName() string {
	return "ModelAssociationOpDetail"
}

// ServiceInstanceOpDetail service instance operation detail
type ServiceInstanceOpDetail struct {
	BasicOpDetail `json:",inline" bson:",inline"`
	HostID        int64 `json:"bk_host_id" bson:"bk_host_id"`
	// Processes operation detail of processes in service instance
	Processes []SvcInstProOpDetail `json:"processes,omitempty" bson:"processes,omitempty"`
}

// SvcInstProOpDetail process operation detail, aggregated in service instance operation detail
type SvcInstProOpDetail struct {
	// Action process operation action, can be one of CUD
	Action ActionType `json:"action" bson:"action"`
	// ProcessIDs operated process ids
	ProcessIDs int64 `json:"bk_process_ids" bson:"bk_process_ids"`
	// ProcessNames operated process names, used by ui
	ProcessNames  string `json:"bk_process_names" bson:"bk_process_names"`
	BasicOpDetail `json:",inline" bson:",inline"`
}

// WithName returns the service instance operation detail name
func (op *ServiceInstanceOpDetail) WithName() string {
	return "ServiceInstanceOpDetail"
}

// BasicContent TODO
// Content contains the details information with in a user's operation.
// Generally, works for business, model, model instance etc.
type BasicContent struct {
	// PreData the previous data before the deletion or updating operation
	PreData map[string]interface{} `json:"pre_data" bson:"pre_data"`
	// CurData the current date after the creation operation
	CurData map[string]interface{} `json:"cur_data" bson:"cur_data"`
	// UpdateFields the data that user uses to update the pre data, might not be the actual changed data
	UpdateFields map[string]interface{} `json:"update_fields" bson:"update_fields"`
}

// AuditType TODO
type AuditType string

const (
	// BusinessType represent business itself's operation audit. such as you change a business maintainer, it's
	// audit belongs to this kind.
	BusinessType AuditType = "business"

	// BizSetType represent operation audit of biz set itself
	BizSetType AuditType = "biz_set"

	// BusinessResourceType TODO
	// Business resource include resources as follows:
	// - service template
	// - set template
	// - service category
	// - dynamic group
	// - main line instance, such as user-defined topology level, set, module etc.
	// - service instance.
	// - others.
	//
	// Note: host does not belong to business resource, it's a independent resource kind.
	BusinessResourceType AuditType = "business_resource"

	// HostType represent all the host related resource's operation audit.
	HostType AuditType = "host"

	// ModelType represent all the operation audit related with model in the system
	ModelType AuditType = "model"

	// ModelInstanceType represent all the operation audit related with model instance in the system,
	// and the instance association is included.
	ModelInstanceType AuditType = "model_instance"

	// AssociationKindType represent all the association kind operation audit.
	AssociationKindType AuditType = "association_kind"

	// EventPushType represent all the event related operation audit.
	EventPushType AuditType = "event"

	// CloudResourceType represent all the operation audit related with cloud, such as:
	// - cloud area
	// - cloud account
	// - cloud synchronize job
	// - others
	CloudResourceType AuditType = "cloud_resource"

	// DynamicGroupType is dynamic grouping audit type.
	DynamicGroupType AuditType = "dynamic_grouping"

	// PlatFormSettingType is platform audit type
	PlatFormSettingType AuditType = "platform_setting"
)

// ResourceType TODO
type ResourceType string

const (
	// BusinessRes related operation type
	BusinessRes ResourceType = "business"
	// BizSetRes TODO
	BizSetRes ResourceType = "biz_set"
	// ServiceTemplateRes TODO
	ServiceTemplateRes ResourceType = "service_template"
	// SetTemplateRes TODO
	SetTemplateRes ResourceType = "set_template"
	// ServiceCategoryRes TODO
	ServiceCategoryRes ResourceType = "service_category"
	// DynamicGroupRes TODO
	DynamicGroupRes ResourceType = "dynamic_group"
	// ServiceInstanceRes TODO
	ServiceInstanceRes ResourceType = "service_instance"
	// ServiceInstanceLabelRes TODO
	ServiceInstanceLabelRes ResourceType = "service_instance_label"
	// SetRes TODO
	SetRes ResourceType = "set"
	// ModuleRes TODO
	ModuleRes ResourceType = "module"
	// ProcessRes TODO
	ProcessRes ResourceType = "process"
	// HostApplyRes TODO
	HostApplyRes ResourceType = "host_apply"
	// CustomFieldRes TODO
	CustomFieldRes ResourceType = "custom_field"

	// ModelRes TODO
	// model related operation type
	ModelRes ResourceType = "model"
	// ModelInstanceRes TODO
	ModelInstanceRes ResourceType = "model_instance"
	// MainlineInstanceRes TODO
	MainlineInstanceRes ResourceType = "mainline_instance"
	// ModelAssociationRes TODO
	ModelAssociationRes ResourceType = "model_association"
	// ModelAttributeRes TODO
	ModelAttributeRes ResourceType = "model_attribute"
	// ModelAttributeGroupRes TODO
	ModelAttributeGroupRes ResourceType = "model_attribute_group"
	// ModelClassificationRes TODO
	ModelClassificationRes ResourceType = "model_classification"
	// InstanceAssociationRes TODO
	InstanceAssociationRes ResourceType = "instance_association"
	// ModelGroupRes TODO
	ModelGroupRes ResourceType = "model_group"
	// ModelUniqueRes TODO
	ModelUniqueRes ResourceType = "model_unique"
	// ResourceDirectoryRes TODO
	ResourceDirectoryRes ResourceType = "resource_directory"

	// AssociationKindRes TODO
	AssociationKindRes ResourceType = "association_kind"
	// EventPushRes TODO
	EventPushRes ResourceType = "event"
	// CloudAreaRes TODO
	CloudAreaRes ResourceType = "cloud_area"
	// CloudAccountRes TODO
	CloudAccountRes ResourceType = "cloud_account"
	// CloudSyncTaskRes TODO
	CloudSyncTaskRes ResourceType = "cloud_sync_task"

	// HostRes TODO
	// host related operation type
	HostRes ResourceType = "host"

	// ResourceDirRes TODO
	ResourceDirRes ResourceType = "resource_directory"

	// PlatFormSettingRes platform related operation type
	PlatFormSettingRes ResourceType = "platform_setting"
)

// OperateFromType TODO
type OperateFromType string

const (
	// FromCCSystem means this audit come from cc system operation, such as upgrader.
	FromCCSystem OperateFromType = "cc_system"
	// FromUser means this audit come from a user's operation, such as web.
	FromUser OperateFromType = "user"
	// FromDataCollection means this audit is created by data collection.
	FromDataCollection OperateFromType = "data_collection"
	// FromSynchronizer means this audit is created by the data synchronizer.
	FromSynchronizer OperateFromType = "synchronizer"
	// FromCloudSync means this audit is created by cloud sync.
	FromCloudSync OperateFromType = "cloud_sync"
)

// ActionType defines all the user's operation type
type ActionType string

const (
	// AuditCreate TODO
	// create a resource
	AuditCreate ActionType = "create"
	// AuditUpdate TODO
	// update a resource
	AuditUpdate ActionType = "update"
	// AuditDelete TODO
	// delete a resource
	AuditDelete ActionType = "delete"
	// AuditAssignHost TODO
	// transfer a host from resource pool to biz
	AuditAssignHost ActionType = "assign_host"
	// AuditUnassignHost TODO
	// transfer a host from biz to resource pool
	AuditUnassignHost ActionType = "unassign_host"
	// AuditTransferHostModule TODO
	// transfer host to another module
	AuditTransferHostModule ActionType = "transfer_host_module"
	// AuditArchive TODO
	// archive a resource
	AuditArchive ActionType = "archive"
	// AuditRecover TODO
	// recover a resource
	AuditRecover ActionType = "recover"
	// AuditPause TODO
	// pause an object
	AuditPause ActionType = "stop"
	// AuditResume TODO
	// resume using an object
	AuditResume ActionType = "resume"
)

// GetAuditTypeByObjID TODO
func GetAuditTypeByObjID(objID string, isMainline bool) AuditType {
	switch objID {
	case common.BKInnerObjIDBizSet:
		return BizSetType
	case common.BKInnerObjIDApp:
		return BusinessType
	case common.BKInnerObjIDSet:
		return BusinessResourceType
	case common.BKInnerObjIDModule:
		return BusinessResourceType
	case common.BKInnerObjIDObject:
		return ModelInstanceType
	case common.BKInnerObjIDHost:
		return HostType
	case common.BKInnerObjIDProc:
		return BusinessResourceType
	case common.BKInnerObjIDPlat:
		return CloudResourceType
	default:
		if isMainline {
			return BusinessResourceType
		}
		return ModelInstanceType
	}
}

// GetResourceTypeByObjID TODO
func GetResourceTypeByObjID(objID string, isMainline bool) ResourceType {
	switch objID {
	case common.BKInnerObjIDBizSet:
		return BizSetRes
	case common.BKInnerObjIDApp:
		return BusinessRes
	case common.BKInnerObjIDSet:
		return SetRes
	case common.BKInnerObjIDModule:
		return ModuleRes
	case common.BKInnerObjIDObject:
		return ModelInstanceRes
	case common.BKInnerObjIDHost:
		return HostRes
	case common.BKInnerObjIDProc:
		return ProcessRes
	case common.BKInnerObjIDPlat:
		return CloudAreaRes
	default:
		if isMainline {
			return MainlineInstanceRes
		}
		return ModelInstanceRes
	}
}

// GetAuditTypesByCategory return audittype array by check category
func GetAuditTypesByCategory(category string) []AuditType {
	switch category {
	case "business":
		return []AuditType{BusinessResourceType, DynamicGroupType}
	case "resource":
		return []AuditType{BusinessType, BizSetType, ModelInstanceType, CloudResourceType}
	case "host":
		return []AuditType{HostType}
	case "other":
		return []AuditType{ModelType, AssociationKindType, EventPushType, DynamicGroupType, PlatFormSettingType}
	}
	return []AuditType{}
}

// GetAuditDict get audit dict according to language type
func GetAuditDict(languageType common.LanguageType) []resourceTypeInfo {
	switch languageType {
	case common.Chinese:
		return auditDict
	case common.English:
		return auditEnDict
	}
	return nil
}

// 注意：记得在auditEnDict中添加对应的英文
var auditDict = []resourceTypeInfo{
	{
		ID:   ModuleRes,
		Name: "模块",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   SetRes,
		Name: "集群",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   MainlineInstanceRes,
		Name: "节点",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   ServiceInstanceRes,
		Name: "服务实例",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   HostRes,
		Name: "主机",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
			actionInfoMap[AuditAssignHost],
			actionInfoMap[AuditUnassignHost],
			actionInfoMap[AuditTransferHostModule],
		},
	},
	{
		ID:   BusinessRes,
		Name: "业务",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditArchive],
			actionInfoMap[AuditRecover],
		},
	},
	{
		ID:   BizSetRes,
		Name: "业务集",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   CloudAreaRes,
		Name: "云区域",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   ModelInstanceRes,
		Name: "模型实例",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   InstanceAssociationRes,
		Name: "实例关联",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   ResourceDirectoryRes,
		Name: "资源池目录",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   ModelGroupRes,
		Name: "模型分组",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
			actionInfoMap[AuditPause],
			actionInfoMap[AuditResume],
		},
	},
	{
		ID:   ModelRes,
		Name: "模型",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   ModelAttributeRes,
		Name: "模型字段",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   ModelAttributeGroupRes,
		Name: "模型字段分组",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   CloudAccountRes,
		Name: "云账户",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   CloudSyncTaskRes,
		Name: "云资源同步任务",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   DynamicGroupRes,
		Name: "动态分组",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
	{
		ID:   PlatFormSettingRes,
		Name: "平台管理",
		Operations: []actionTypeInfo{
			actionInfoMap[AuditCreate],
			actionInfoMap[AuditUpdate],
			actionInfoMap[AuditDelete],
		},
	},
}

// 注意：记得在actionInfoEnMap中添加对应的英文
var actionInfoMap = map[ActionType]actionTypeInfo{
	AuditCreate:             {ID: AuditCreate, Name: "新增"},
	AuditUpdate:             {ID: AuditUpdate, Name: "修改"},
	AuditDelete:             {ID: AuditDelete, Name: "删除"},
	AuditAssignHost:         {ID: AuditAssignHost, Name: "分配到业务"},
	AuditUnassignHost:       {ID: AuditUnassignHost, Name: "归还到资源池"},
	AuditTransferHostModule: {ID: AuditTransferHostModule, Name: "转移模块"},
	AuditArchive:            {ID: AuditArchive, Name: "归档"},
	AuditRecover:            {ID: AuditRecover, Name: "恢复"},
	AuditPause:              {ID: AuditPause, Name: "停用"},
	AuditResume:             {ID: AuditResume, Name: "启用"},
}

type resourceTypeInfo struct {
	ID         ResourceType     `json:"id"`
	Name       string           `json:"name"`
	Operations []actionTypeInfo `json:"operations"`
}

type actionTypeInfo struct {
	ID   ActionType `json:"id"`
	Name string     `json:"name"`
}

// AuditDetailQueryInput TODO
type AuditDetailQueryInput struct {
	IDs []int64 `json:"id"`
}

// Validate validates the input param
func (input *AuditDetailQueryInput) Validate() errors.RawErrorInfo {
	if len(input.IDs) > common.BKAuditLogPageLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommPageLimitIsExceeded,
		}
	}

	if len(input.IDs) <= 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKFieldID},
		}
	}

	return errors.RawErrorInfo{}
}

var auditEnDict = []resourceTypeInfo{
	{
		ID:   ModuleRes,
		Name: "Module",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   SetRes,
		Name: "Set",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   MainlineInstanceRes,
		Name: "Node",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   ServiceInstanceRes,
		Name: "Service Instance",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   HostRes,
		Name: "Host",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
			actionInfoEnMap[AuditAssignHost],
			actionInfoEnMap[AuditUnassignHost],
			actionInfoEnMap[AuditTransferHostModule],
		},
	},
	{
		ID:   BusinessRes,
		Name: "Business",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditArchive],
			actionInfoEnMap[AuditRecover],
		},
	},
	{
		ID:   BizSetRes,
		Name: "Business Set",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   CloudAreaRes,
		Name: "Cloud Area",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   ModelInstanceRes,
		Name: "Model Instance",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   InstanceAssociationRes,
		Name: "Instance Association",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   ResourceDirectoryRes,
		Name: "Resource Directory",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   ModelGroupRes,
		Name: "Model Group",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
			actionInfoEnMap[AuditPause],
			actionInfoEnMap[AuditResume],
		},
	},
	{
		ID:   ModelRes,
		Name: "Model",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   ModelAttributeRes,
		Name: "Model Attribute",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   ModelAttributeGroupRes,
		Name: "Model Attribute Group",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   CloudAccountRes,
		Name: "Cloud Account",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   CloudSyncTaskRes,
		Name: "Cloud Synchronous Task",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   DynamicGroupRes,
		Name: "Dynamic Group",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
	{
		ID:   PlatFormSettingRes,
		Name: "Platform Management",
		Operations: []actionTypeInfo{
			actionInfoEnMap[AuditCreate],
			actionInfoEnMap[AuditUpdate],
			actionInfoEnMap[AuditDelete],
		},
	},
}

var actionInfoEnMap = map[ActionType]actionTypeInfo{
	AuditCreate:             {ID: AuditCreate, Name: "Create"},
	AuditUpdate:             {ID: AuditUpdate, Name: "Update"},
	AuditDelete:             {ID: AuditDelete, Name: "Delete"},
	AuditAssignHost:         {ID: AuditAssignHost, Name: "Assign to business"},
	AuditUnassignHost:       {ID: AuditUnassignHost, Name: "Return to the resource pool"},
	AuditTransferHostModule: {ID: AuditTransferHostModule, Name: "Transfer module"},
	AuditArchive:            {ID: AuditArchive, Name: "Archive"},
	AuditRecover:            {ID: AuditRecover, Name: "Recover"},
	AuditPause:              {ID: AuditPause, Name: "Pause"},
	AuditResume:             {ID: AuditResume, Name: "Resume"},
}
