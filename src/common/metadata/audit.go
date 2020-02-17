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
	"configcenter/src/common/auditoplog"
)

type SaveAuditLogParams struct {
	ID      int64                  `json:"inst_id"`
	Model   string                 `json:"op_target"`
	Content interface{}            `json:"content"`
	ExtKey  string                 `json:"ext"`
	OpDesc  string                 `json:"op_desc"`
	OpType  auditoplog.AuditOpType `json:"op_type"`
	BizID   int64                  `json:"biz_id"`
}

// AuditQueryResult add single host log paramm
type AuditQueryResult struct {
	BaseResp `json:",inline"`
	Data     struct {
		Count int            `json:"count"`
		Info  []OperationLog `json:"info"`
	} `json:"data"`
}

type AuditLog struct {
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
}

type DetailFactory interface {
	WithName() string
}

type BasicOpDetail struct {
	// the business id of the resource if it belongs to a business.
	BusinessID int64 `json:"bk_biz_id" bson:"bk_biz_id"`
	// ResourceID is the id of the resource instance. which is a unique id.
	ResourceID int64 `json:"resource_id" bson:"resource_id"`
	// ResourceName is the name of the resource, such as a switch model has a name "switch"
	ResourceName string `json:"resource_name" bson:"resource_name"`
	// Details contains all the details information about a user's operation
	Details *BasicContent `json:"details" bson:"details"`
}

func (Op *BasicOpDetail) WithName() string {
	return "BasicDetail"
}

type AssociationOpDetail struct {
	AssociationID   string `json:"asso_id" bson:"asso_id"`
	SourceModel     string `json:"src_model" bson:"src_model"`
	TargetModel     string `json:"target_model" bson:"target_model"`
	SourceModelID   int64  `json:"src_model_id" bson:"src_model_id"`
	SourceModelName string `json:"src_model_name" bson:"src_model_name"`
	TargetModelID   int64  `json:"target_model_id" bson:"target_model_id"`
	TargetModelName int64  `json:"target_model_name" bson:"target_model_name"`
}

func (ao *AssociationOpDetail) WithName() string {
	return "AssociationOpDetail"
}

// Content contains the details information with in a user's operation.
// Generally, works for business, model, model instance etc.
type BasicContent struct {
	// the previous data before the operation
	PreData map[string]interface{} `json:"pre_data" bson:"pre_data"`
	// the current date being operated
	CurData map[string]interface{} `json:"cur_data" bson:"cur_data"`
	// data properties being operated, normally is a model's attributes.
	Properties []Property `json:"properties" bson:"properties"`
}

type Property struct {
	PropertyID   string `json:"bk_property_id" bson:"bk_property_id"`
	PropertyName string `json:"bk_property_name" bson:"bk_property_name"`
}

type AuditType string

const (
	// BusinessKind represent business itself's operation audit. such as you change a business maintainer, it's
	// audit belongs to this kind.
	BusinessType AuditType = "business"

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

	// EventType represent all the event related operation audit.
	EventPushType AuditType = "event"

	// CloudResource represent all the operation audit related with cloud, such as:
	// - cloud area
	// - cloud account
	// - cloud synchronize job
	// - others
	CloudResourceType AuditType = "cloud_resource"
)

type ResourceType string

const (
	// business related operation type
	BusinessRes        ResourceType = "business"
	ServiceTemplateRes ResourceType = "service_template"
	SetTemplateRes     ResourceType = "set_template"
	ServiceCategoryRes ResourceType = "service_category"
	DynamicGroupRes    ResourceType = "dynamic_group"
	ServiceInstanceRes ResourceType = "service_instance"
	SetRes             ResourceType = "set"
	ModuleRes          ResourceType = "module"

	// model related operation type
	ModelRes               ResourceType = "model"
	ModelInstanceRes       ResourceType = "model_instance"
	ModelAssociationRes    ResourceType = "model_association"
	InstanceAssociationRes ResourceType = "instance_association"
	ModelGroupRes          ResourceType = "model_group"
	ModelUniqueRes         ResourceType = "model_unique"

	AssociationKindRes ResourceType = "association_kind"
	CloudAccountRes    ResourceType = "cloud_account"
	CloudSyncTaskRes   ResourceType = "cloud_sync_task"
)

type OperateFromType string

const (
	// FromUser means this audit come from a user's operation, such as web.
	FromUser OperateFromType = "user"
	// FromDataCollection means this audit is created by data collection.
	FromDataCollection OperateFromType = "data_collection"
	// FromSynchronizer means this audit is created by the data synchronizer.
	FromSynchronizer OperateFromType = "synchronizer"
)

// ActionType defines all the user's operation type
type ActionType string

const (
	// create a resource
	AuditCreate ActionType = "create"
	// update a resource
	AuditUpdate ActionType = "update"
	// delete a resource
	AuditDelete ActionType = "delete"
	// transfer a host from to resource pool or
	// transfer host to a business.
	AuditTransferHost ActionType = "transfer_host"
)
