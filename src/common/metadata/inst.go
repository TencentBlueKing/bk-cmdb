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
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
)

// SetInst contains partial fields of a real set
type SetInst struct {
	BizID         int64  `bson:"bk_biz_id" json:"bk_biz_id" mapstructure:"bk_biz_id"`
	SetID         int64  `bson:"bk_set_id" json:"bk_set_id" mapstructure:"bk_set_id"`
	SetName       string `bson:"bk_set_name" json:"bk_set_name" mapstructure:"bk_set_name"`
	SetStatus     string `bson:"bk_service_status" json:"bk_service_status" mapstructure:"bk_service_status"`
	SetEnv        string `bson:"bk_set_env" json:"bk_set_env" mapstructure:"bk_set_env"`
	SetTemplateID int64  `bson:"set_template_id" json:"set_template_id" mapstructure:"set_template_id"`
	ParentID      int64  `bson:"bk_parent_id" json:"bk_parent_id" mapstructure:"bk_parent_id"`

	Creator         string `field:"creator" json:"creator,omitempty" bson:"creator" mapstructure:"creator"`
	CreateTime      Time   `field:"create_time" json:"create_time,omitempty" bson:"create_time" mapstructure:"create_time"`
	LastTime        Time   `field:"last_time" json:"last_time,omitempty" bson:"last_time" mapstructure:"last_time"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account" mapstructure:"bk_supplier_account"`
}

// ModuleInst contains partial fields of a real module
type ModuleInst struct {
	BizID             int64  `bson:"bk_biz_id" json:"bk_biz_id" field:"bk_biz_id" mapstructure:"bk_biz_id"`
	SetID             int64  `bson:"bk_set_id" json:"bk_set_id" field:"bk_set_id" mapstructure:"bk_set_id"`
	ModuleID          int64  `bson:"bk_module_id" json:"bk_module_id" field:"bk_module_id" mapstructure:"bk_module_id"`
	ModuleName        string `bson:"bk_module_name" json:"bk_module_name" field:"bk_module_name" mapstructure:"bk_module_name"`
	SupplierAccount   string `bson:"bk_supplier_account" json:"bk_supplier_account" field:"bk_supplier_account" mapstructure:"bk_supplier_account"`
	ServiceCategoryID int64  `bson:"service_category_id" json:"service_category_id" field:"service_category_id" mapstructure:"service_category_id"`
	ServiceTemplateID int64  `bson:"service_template_id" json:"service_template_id" field:"service_template_id" mapstructure:"service_template_id"`
	ParentID          int64  `bson:"bk_parent_id" json:"bk_parent_id" field:"bk_parent_id" mapstructure:"bk_parent_id"`
	SetTemplateID     int64  `bson:"set_template_id" json:"set_template_id" field:"set_template_id" mapstructure:"set_template_id"`
	Default           int64  `bson:"default" json:"default" field:"default" mapstructure:"default"`
	HostApplyEnabled  bool   `bson:"host_apply_enabled" json:"host_apply_enabled" field:"host_apply_enabled" mapstructure:"host_apply_enabled"`
	Creator           string `bson:"creator" json:"creator" field:"creator" mapstructure:"creator"`
	CreateTime        Time   `bson:"create_time" json:"create_time"`
	LastTime          Time   `bson:"last_time" json:"last_time"`
}

// BizInst TODO
type BizInst struct {
	BizID           int64  `bson:"bk_biz_id" mapstructure:"bk_biz_id"`
	BizName         string `bson:"bk_biz_name" mapstructure:"bk_biz_name"`
	SupplierAccount string `bson:"bk_supplier_account" mapstructure:"bk_supplier_account"`
}

// BizBasicInfo TODO
type BizBasicInfo struct {
	BizID   int64  `bson:"bk_biz_id" json:"bk_biz_id" field:"bk_biz_id" mapstructure:"bk_biz_id"`
	BizName string `bson:"bk_biz_name" json:"bk_biz_name" field:"bk_biz_name" mapstructure:"bk_biz_name"`
}

// BizSetInst biz set structure with pre-defined properties
type BizSetInst struct {
	BizSetID         int64       `json:"bk_biz_set_id" bson:"bk_biz_set_id"`
	BizSetName       string      `json:"bk_biz_set_name" bson:"bk_biz_set_name"`
	Description      string      `json:"description" bson:"description"`
	BizSetMaintainer string      `json:"bk_biz_maintainer" bson:"bk_biz_maintainer"`
	CreateTime       Time        `json:"create_time" bson:"create_time"`
	LastTime         Time        `json:"last_time" bson:"last_time"`
	SupplierAccount  string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	Scope            BizSetScope `json:"bk_scope" bson:"bk_scope"`
	Default          int64       `json:"default" bson:"default"`
}

// CloudInst TODO
type CloudInst struct {
	CloudID   int64  `bson:"bk_cloud_id" json:"bk_cloud_id"`
	CloudName string `bson:"bk_cloud_name" json:"bk_cloud_name"`
}

// ProcessInst TODO
type ProcessInst struct {
	ProcessID       int64  `json:"bk_process_id" bson:"bk_process_id"`               // 进程名称
	ProcessName     string `json:"bk_process_name" bson:"bk_process_name"`           // 进程名称
	BindIP          string `json:"bind_ip" bson:"bind_ip"`                           // 绑定IP, 枚举: [{ID: "1", Name: "127.0.0.1"}, {ID: "2", Name: "0.0.0.0"}, {ID: "3", Name: "第一内网IP"}, {ID: "4", Name: "第一外网IP"}]
	PORT            string `json:"port" bson:"port"`                                 // 端口, 单个端口："8080", 多个连续端口："8080-8089", 多个不连续端口："8080-8089,8199"
	PROTOCOL        string `json:"protocol" bson:"protocol"`                         // 协议, 枚举: [{ID: "1", Name: "TCP"}, {ID: "2", Name: "UDP"}],
	FuncName        string `json:"bk_func_name" bson:"bk_func_name"`                 // 功能名称
	StartParamRegex string `json:"bk_start_param_regex" bson:"bk_start_param_regex"` // 启动参数匹配规则
}

// HostIdentifier TODO
type HostIdentifier struct {
	HostID          int64                       `json:"bk_host_id" bson:"bk_host_id"`
	CloudID         int64                       `json:"bk_cloud_id" bson:"bk_cloud_id"`
	InnerIP         StringArrayToString         `json:"bk_host_innerip" bson:"bk_host_innerip"`
	OSType          string                      `json:"bk_os_type" bson:"bk_os_type"`
	SupplierAccount string                      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	HostIdentModule map[string]*HostIdentModule `json:"associations" bson:"associations"`
	Process         []HostIdentProcess          `json:"process" bson:"process"`
}

// HostIdentProcess TODO
type HostIdentProcess struct {
	ProcessID   int64  `json:"bk_process_id" bson:"bk_process_id"`     // 进程名称
	ProcessName string `json:"bk_process_name" bson:"bk_process_name"` // 进程名称
	// deprecated  后续的版本会被废弃掉
	BindIP string `json:"bind_ip" bson:"bind_ip"` // 绑定IP, 枚举: [{ID: "1", Name: "127.0.0.1"}, {ID: "2", Name: "0.0.0.0"}, {ID: "3", Name: "第一内网IP"}, {ID: "4", Name: "第一外网IP"}]
	// deprecated  后续的版本会被废弃掉
	Port string `json:"port" bson:"port"` // 端口, 单个端口："8080", 多个连续端口："8080-8089", 多个不连续端口："8080-8089,8199"
	// deprecated  后续的版本会被废弃掉
	Protocol        string `json:"protocol" bson:"protocol"`                         // 协议, 枚举: [{ID: "1", Name: "TCP"}, {ID: "2", Name: "UDP"}],
	FuncName        string `json:"bk_func_name" bson:"bk_func_name"`                 // 功能名称
	StartParamRegex string `json:"bk_start_param_regex" bson:"bk_start_param_regex"` // 启动参数匹配规则
	// deprecated  后续的版本会被废弃掉
	PortEnable bool `field:"bk_enable_port" json:"bk_enable_port" bson:"bk_enable_port"`
	// BindInfo 进程绑定信息
	BindInfo []ProcBindInfo `field:"bind_info" json:"bind_info" bson:"bind_info"`
}

// HostIdentModule HostIdentifier module define
type HostIdentModule struct {
	BizID    int64  `json:"bk_biz_id"`    // 业务ID
	SetID    int64  `json:"bk_set_id"`    // 所属集群(bk_set_id)
	ModuleID int64  `json:"bk_module_id"` // 所属模块(bk_module_id)
	Layer    *Layer `json:"layer"`        // 自定义层级
}

// Layer TODO
type Layer struct {
	InstID   int64  `json:"bk_inst_id"`
	InstName string `json:"bk_inst_name"`
	ObjID    string `json:"bk_obj_id"`
	Child    *Layer `json:"child"`
}

// MainlineInstInfo TODO
type MainlineInstInfo struct {
	InstID   int64  `json:"bk_inst_id" bson:"bk_inst_id"`
	InstName string `json:"bk_inst_name" bson:"bk_inst_name"`
	ObjID    string `json:"bk_obj_id" bson:"bk_obj_id"`
	ParentID int64  `json:"bk_parent_id" bson:"bk_parent_id"`
}

// SearchIdentifierParam defines the param
type SearchIdentifierParam struct {
	IP   IPParam `json:"ip"`
	Page BasePage
}

// SearchHostIdentifierParam 查询主机身份的条件
type SearchHostIdentifierParam struct {
	HostIDs []int64 `json:"host_ids"`
}

// IPParam TODO
type IPParam struct {
	Data    []string `json:"data"`
	CloudID *int64   `json:"bk_cloud_id"`
}

// SearchHostIdentifierResult TODO
type SearchHostIdentifierResult struct {
	BaseResp `json:",inline"`
	Data     SearchHostIdentifierData `json:"data"`
}

// SearchHostIdentifierData host identifier detail
type SearchHostIdentifierData struct {
	Count int              `json:"count"`
	Info  []HostIdentifier `json:"info"`
}

// SearchInstsNamesOption search instances names option
type SearchInstsNamesOption struct {
	ObjID string `json:"bk_obj_id"`
	BizID int64  `json:"bk_biz_id"`
	Name  string `json:"name"`
}

// UpdateBizPropertyBatchParameter batch update business properties parameter
/* e.g.:
{
    "properties":{
      "bk_biz_developer":"developer",
      "bk_biz_maintainer": "maintainer",
      "bk_biz_name":"biz_test",
      "bk_biz_productor": "productor",
      "bk_biz_tester":"tester",
      "operator": "operator"
    },
    "condition": {
        "bk_biz_id": {"$in": [3,4]}"
    }
}
*/
type UpdateBizPropertyBatchParameter struct {
	// Properties is business property keys and values to be updated
	Properties map[string]interface{} `json:"properties"`
	// Condition is business property update condition
	Condition map[string]interface{} `json:"condition"`
}

// ObjsForSearchName TODO
var ObjsForSearchName = map[string]bool{
	common.BKInnerObjIDSet:    true,
	common.BKInnerObjIDModule: true,
}

// Validate verify the SearchInstsNamesOption
func (o *SearchInstsNamesOption) Validate() (rawError errors.RawErrorInfo) {
	if _, ok := ObjsForSearchName[o.ObjID]; !ok {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"bk_obj_id"},
		}
	}

	if o.BizID <= 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"bk_biz_id"},
		}
	}

	if o.Name == "" {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"name"},
		}
	}

	return errors.RawErrorInfo{}
}

// CreateManyCommInst parameters for creating multiple instances
type CreateManyCommInst struct {
	ObjID   string          `json:"bk_obj_id"`
	Details []mapstr.MapStr `json:"details"`
}

// CreateManyCommInstResult result of creating multiple instances
type CreateManyCommInstResult struct {
	BaseResp `json:",inline"`
	Data     CreateManyCommInstResultDetail `json:"data"`
}

// CreateManyCommInstResultDetail details of CreateManyInstancesResult
type CreateManyCommInstResultDetail struct {
	SuccessCreated map[int64]int64  `json:"success_created"`
	Error          map[int64]string `json:"error_msg"`
}

// NewManyCommInstResultDetail TODO
func NewManyCommInstResultDetail() *CreateManyCommInstResultDetail {
	return &CreateManyCommInstResultDetail{
		SuccessCreated: make(map[int64]int64, 0),
		Error:          make(map[int64]string, 0),
	}
}

// InstBatchInfo info struct of batch create or update inst
type InstBatchInfo struct {
	// BatchInfo batch info
	// map[rownumber]map[property_id][date]
	BatchInfo map[int64]mapstr.MapStr `json:"BatchInfo"`
	InputType string                  `json:"input_type"`
}

// GetInstID get inst id by objid
func GetInstID(objID string, data mapstr.MapStr) (int64, error) {
	return data.Int64(GetInstIDFieldByObjID(objID))
}

// GetInstName get inst name by objid
func GetInstName(objID string, data mapstr.MapStr) (string, error) {
	return data.String(GetInstNameFieldName(objID))
}

// GetBizID get biz id by inst data
func GetBizID(data mapstr.MapStr) (int64, error) {
	return data.Int64(common.BKAppIDField)
}

// GetParentID get parent id by inst data
func GetParentID(data mapstr.MapStr) (int64, error) {
	return data.Int64(common.BKParentIDField)
}

type opcondition struct {
	InstID []int64 `json:"inst_ids"`
}

type deleteCondition struct {
	opcondition `json:",inline"`
}

type updateCondition struct {
	InstID   int64                  `json:"inst_id"`
	InstInfo map[string]interface{} `json:"datas"`
}

// OpCondition the condition operation
type OpCondition struct {
	Delete deleteCondition   `json:"delete"`
	Update []updateCondition `json:"update"`
}
