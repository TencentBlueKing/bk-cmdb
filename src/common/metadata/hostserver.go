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
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
)

// DeleteHostBatchOpt TODO
type DeleteHostBatchOpt struct {
	HostID string `json:"bk_host_id"`
}

// HostInstanceProperties TODO
type HostInstanceProperties struct {
	PropertyID    string      `json:"bk_property_id"`
	PropertyName  string      `json:"bk_property_name"`
	PropertyValue interface{} `json:"bk_property_value"`
}

// HostInstancePropertiesResult TODO
type HostInstancePropertiesResult struct {
	BaseResp `json:",inline"`
	Data     []HostInstanceProperties `json:"data"`
}

// HostInputType TODO
type HostInputType string

const (
	// ExcelType TODO
	ExcelType HostInputType = "excel"
	// CollectType TODO
	CollectType HostInputType = "collect"
)

// HostList TODO
type HostList struct {
	ApplicationID int64                            `json:"bk_biz_id"`
	ModuleID      int64                            `json:"bk_module_id"`
	HostInfo      map[int64]map[string]interface{} `json:"host_info"`
	InputType     HostInputType                    `json:"input_type"`
}

// AddHostToResourcePoolHostList TODO
type AddHostToResourcePoolHostList struct {
	HostInfo  []map[string]interface{} `json:"host_info"`
	Directory int64                    `json:"directory"`
}

// AddHostToResourcePoolResult TODO
type AddHostToResourcePoolResult struct {
	Success []AddOneHostToResourcePoolResult `json:"success,omitempty"`
	Error   []AddOneHostToResourcePoolResult `json:"error,omitempty"`
}

// AddOneHostToResourcePoolResult TODO
type AddOneHostToResourcePoolResult struct {
	Index    int    `json:"index"`
	HostID   int64  `json:"bk_host_id,omitempty"`
	ErrorMsg string `json:"error_message,omitempty"`
}

// AddHostFromAgentHostList TODO
type AddHostFromAgentHostList struct {
	HostInfo map[string]interface{} `json:"host_info"`
}

// HostSyncList TODO
type HostSyncList struct {
	ApplicationID int64                            `json:"bk_biz_id"`
	HostInfo      map[int64]map[string]interface{} `json:"host_info"`
	ModuleID      []int64                          `json:"bk_module_id"`
	InputType     HostInputType                    `json:"input_type"`
}

// HostsModuleRelation TODO
type HostsModuleRelation struct {
	ApplicationID int64   `json:"bk_biz_id"`
	HostID        []int64 `json:"bk_host_id"`
	ModuleID      []int64 `json:"bk_module_id"`
	IsIncrement   bool    `json:"is_increment"`
	// DisableAutoCreateSvcInst disable auto create service instance when transfer to a module with process in template
	DisableAutoCreateSvcInst bool `json:"disable_auto_create"`

	// DisableTransferHostAutoApply when this flag is true, it means that the user specifies not to automatically apply
	// the host in the host transfer scenario.
	DisableTransferHostAutoApply bool
}

// HostModuleConfig TODO
type HostModuleConfig struct {
	ApplicationID int64   `json:"bk_biz_id" bson:"bk_biz_id"`
	HostID        []int64 `json:"bk_host_id" bson:"bk_host_id"`
	ModuleID      []int64 `json:"bk_module_id" bson:"bk_module_id"`
}

// RemoveHostsFromModuleOption TODO
type RemoveHostsFromModuleOption struct {
	ApplicationID int64 `json:"bk_biz_id"`
	HostID        int64 `json:"bk_host_id"`
	ModuleID      int64 `json:"bk_module_id"`
}

// HostToAppModule TODO
type HostToAppModule struct {
	Ips         []string `json:"ips"`
	HostName    []string `json:"bk_host_name"`
	ModuleName  string   `json:"bk_module_name"`
	SetName     string   `json:"bk_set_name"`
	AppName     string   `json:"bk_biz_name"`
	OsType      string   `json:"bk_os_type"`
	OwnerID     string   `json:"bk_supplier_account"`
	PlatID      int64    `json:"bk_cloud_id"`
	IsIncrement bool     `json:"is_increment"`
}

// HostCommonSearch TODO
type HostCommonSearch struct {
	AppID     int64             `json:"bk_biz_id,omitempty"`
	Ip        IPInfo            `json:"ip"`
	Condition []SearchCondition `json:"condition"`
	Page      BasePage          `json:"page"`
	Pattern   string            `json:"pattern,omitempty"`
}

// SetCommonSearch TODO
type SetCommonSearch struct {
	AppID     int64             `json:"bk_biz_id,omitempty"`
	Condition []SearchCondition `json:"condition"`
	Page      BasePage          `json:"page"`
}

// FindHostsBySrvTplOpt TODO
type FindHostsBySrvTplOpt struct {
	ServiceTemplateIDs []int64  `json:"bk_service_template_ids"`
	ModuleIDs          []int64  `json:"bk_module_ids"`
	Fields             []string `json:"fields"`
	Page               BasePage `json:"page"`
}

// Validate TODO
func (o *FindHostsBySrvTplOpt) Validate() (rawError errors.RawErrorInfo) {
	if len(o.ServiceTemplateIDs) == 0 || len(o.ServiceTemplateIDs) > common.BKMaxInstanceLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrArrayLengthWrong,
			Args:    []interface{}{"bk_service_template_ids", common.BKMaxInstanceLimit},
		}
	}

	if len(o.ModuleIDs) > common.BKMaxInstanceLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrArrayLengthWrong,
			Args:    []interface{}{"bk_module_ids", common.BKMaxInstanceLimit},
		}
	}

	if len(o.Fields) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"fields"},
		}
	}

	if o.Page.IsIllegal() {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"page.limit"},
		}
	}

	return errors.RawErrorInfo{}
}

// FindHostsBySetTplOpt TODO
type FindHostsBySetTplOpt struct {
	SetTemplateIDs []int64  `json:"bk_set_template_ids"`
	SetIDs         []int64  `json:"bk_set_ids"`
	Fields         []string `json:"fields"`
	Page           BasePage `json:"page"`
}

// Validate TODO
func (o *FindHostsBySetTplOpt) Validate() (rawError errors.RawErrorInfo) {
	if len(o.SetTemplateIDs) == 0 || len(o.SetTemplateIDs) > common.BKMaxInstanceLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrArrayLengthWrong,
			Args:    []interface{}{"bk_set_template_ids", common.BKMaxInstanceLimit},
		}
	}

	if len(o.SetIDs) > common.BKMaxInstanceLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrArrayLengthWrong,
			Args:    []interface{}{"bk_set_ids", common.BKMaxInstanceLimit},
		}
	}

	if len(o.Fields) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"fields"},
		}
	}

	if o.Page.IsIllegal() {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"page.limit"},
		}
	}

	return errors.RawErrorInfo{}
}

// FindHostsByTopoOpt TODO
type FindHostsByTopoOpt struct {
	ObjID  string   `json:"bk_obj_id"`
	InstID int64    `json:"bk_inst_id"`
	Fields []string `json:"fields"`
	Page   BasePage `json:"page"`
}

// Validate TODO
func (o *FindHostsByTopoOpt) Validate() (rawError errors.RawErrorInfo) {
	if o.ObjID == "" {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKObjIDField},
		}
	}

	if o.ObjID == common.BKInnerObjIDApp {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{common.BKObjIDField},
		}
	}

	if o.InstID <= 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{common.BKInstIDField},
		}
	}

	if len(o.Fields) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"fields"},
		}
	}

	if o.Page.IsIllegal() {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"page"},
		}
	}

	return errors.RawErrorInfo{}
}

// FindHostRelationWtihTopoOpt TODO
type FindHostRelationWtihTopoOpt struct {
	Business int64    `json:"bk_biz_id"`
	ObjID    string   `json:"bk_obj_id"`
	InstIDs  []int64  `json:"bk_inst_ids"`
	Fields   []string `json:"fields"`
	Page     BasePage `json:"page"`
}

// Validate TODO
func (f *FindHostRelationWtihTopoOpt) Validate() *errors.RawErrorInfo {
	if f.ObjID == "" {
		return &errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKObjIDField},
		}
	}

	if f.ObjID == common.BKInnerObjIDApp {
		return &errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{common.BKObjIDField},
		}
	}

	if f.Business <= 0 {
		return &errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{common.BKAppIDField},
		}
	}

	if len(f.InstIDs) <= 0 || len(f.InstIDs) > 50 {
		return &errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{common.BKInstIDField},
		}
	}

	if len(f.Fields) == 0 {
		return &errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"fields"},
		}
	}

	if f.Page.Limit > common.BKMaxInstanceLimit {
		return &errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"page.limit"},
		}
	}

	if f.Page.IsIllegal() {
		return &errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"page"},
		}
	}

	return nil
}

// FindModuleHostRelationParameter TODO
type FindModuleHostRelationParameter struct {
	ModuleIDS    []int64  `json:"bk_module_ids"`
	ModuleFields []string `json:"module_fields"`
	HostFields   []string `json:"host_fields"`
	Page         BasePage `json:"page"`
}

// Validate TODO
func (param FindModuleHostRelationParameter) Validate() errors.RawErrorInfo {
	if len(param.ModuleIDS) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_module_ids"},
		}
	}
	if len(param.ModuleIDS) > 200 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"bk_module_ids", 200},
		}
	}
	if param.Page.IsIllegal() {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"page"},
		}
	}
	if len(param.HostFields) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"host_fields"},
		}
	}
	if len(param.ModuleFields) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"module_fields"},
		}
	}
	return errors.RawErrorInfo{}
}

// ModuleHostRelation TODO
type ModuleHostRelation struct {
	Host    map[string]interface{}   `json:"host"`
	Modules []map[string]interface{} `json:"modules"`
}

// FindModuleHostRelationResult TODO
type FindModuleHostRelationResult struct {
	Count    int                  `json:"count"`
	Relation []ModuleHostRelation `json:"relation"`
}

// FindModuleHostRelationResp TODO
type FindModuleHostRelationResp struct {
	BaseResp `json:",inline"`
	Data     FindModuleHostRelationResult `json:"data"`
}

// ListHostsParameter TODO
type ListHostsParameter struct {
	SetIDs             []int64                   `json:"bk_set_ids"`
	SetCond            []ConditionItem           `json:"set_cond"`
	ModuleIDs          []int64                   `json:"bk_module_ids"`
	ModuleCond         []ConditionItem           `json:"module_cond"`
	HostPropertyFilter *querybuilder.QueryFilter `json:"host_property_filter"`
	Fields             []string                  `json:"fields"`
	Page               BasePage                  `json:"page"`
}

// Validate TODO
func (option ListHostsParameter) Validate() (string, error) {
	if key, err := option.Page.Validate(false); err != nil {
		return fmt.Sprintf("page.%s", key), err
	}

	if option.HostPropertyFilter != nil {
		if key, err := option.HostPropertyFilter.Validate(&querybuilder.RuleOption{NeedSameSliceElementType: true}); err != nil {
			return fmt.Sprintf("host_property_filter.%s", key), err
		}
		if option.HostPropertyFilter.GetDeep() > querybuilder.MaxDeep {
			return "host_property_filter.rules", fmt.Errorf("exceed max query condition deepth: %d", querybuilder.MaxDeep)
		}
	}

	if len(option.SetIDs) > 200 {
		return "bk_set_ids", fmt.Errorf("exceed max length: 200")
	}

	if len(option.ModuleIDs) > 500 {
		return "bk_module_ids", fmt.Errorf("exceed max length: 500")
	}

	return "", nil
}

// ListHostsWithNoBizParameter TODO
type ListHostsWithNoBizParameter struct {
	HostPropertyFilter *querybuilder.QueryFilter `json:"host_property_filter"`
	Fields             []string                  `json:"fields"`
	Page               BasePage                  `json:"page"`
}

// Validate TODO
func (option ListHostsWithNoBizParameter) Validate() (string, error) {
	if key, err := option.Page.Validate(false); err != nil {
		return fmt.Sprintf("page.%s", key), err
	}

	if option.HostPropertyFilter != nil {
		if key, err := option.HostPropertyFilter.Validate(&querybuilder.RuleOption{NeedSameSliceElementType: true}); err != nil {
			return fmt.Sprintf("host_property_filter.%s", key), err
		}
		if option.HostPropertyFilter.GetDeep() > querybuilder.MaxDeep {
			return "host_property_filter.rules", fmt.Errorf("exceed max query condition deepth: %d", querybuilder.MaxDeep)
		}
	}

	return "", nil
}

// ListBizHostsTopoParameter parameter for listing biz hosts with topology info
type ListBizHostsTopoParameter struct {
	SetPropertyFilter    *querybuilder.QueryFilter `json:"set_property_filter"`
	ModulePropertyFilter *querybuilder.QueryFilter `json:"module_property_filter"`
	HostPropertyFilter   *querybuilder.QueryFilter `json:"host_property_filter"`
	Fields               []string                  `json:"fields"`
	Page                 BasePage                  `json:"page"`
}

// Validate validate if the parameter is valid for listing biz hosts with topology info
func (option ListBizHostsTopoParameter) Validate(errProxy errors.DefaultCCErrorIf) errors.CCErrorCoder {
	if option.Page.Limit == 0 {
		return errProxy.CCErrorf(common.CCErrCommParamsNeedSet, "page.limit")
	}
	if option.Page.Limit > common.BKMaxInstanceLimit {
		return errProxy.CCErrorf(common.CCErrCommXXExceedLimit, "page.limit", common.BKMaxInstanceLimit)
	}

	opt := &querybuilder.RuleOption{NeedSameSliceElementType: true}
	if option.HostPropertyFilter != nil {
		if key, err := option.HostPropertyFilter.Validate(opt); err != nil {
			blog.Errorf("valid host property filter failed, err: %v", err)
			return errProxy.CCErrorf(common.CCErrCommParamsInvalid, fmt.Sprintf("host_property_filter.%s", key))
		}
		if option.HostPropertyFilter.GetDeep() > querybuilder.MaxDeep {
			return errProxy.CCErrorf(common.CCErrCommXXExceedLimit, "host_property_filter.rules", querybuilder.MaxDeep)
		}
	}

	if option.SetPropertyFilter != nil {
		if key, err := option.SetPropertyFilter.Validate(opt); err != nil {
			blog.Errorf("valid set property filter failed, err: %v", err)
			return errProxy.CCErrorf(common.CCErrCommParamsInvalid, fmt.Sprintf("set_property_filter.%s", key))
		}
		if option.SetPropertyFilter.GetDeep() > querybuilder.MaxDeep {
			return errProxy.CCErrorf(common.CCErrCommXXExceedLimit, "set_property_filter.rules", querybuilder.MaxDeep)
		}
	}

	if option.ModulePropertyFilter != nil {
		if key, err := option.ModulePropertyFilter.Validate(opt); err != nil {
			blog.Errorf("valid module property filter failed, err: %v", err)
			return errProxy.CCErrorf(common.CCErrCommParamsInvalid, fmt.Sprintf("module_property_filter.%s", key))
		}
		if option.ModulePropertyFilter.GetDeep() > querybuilder.MaxDeep {
			return errProxy.CCErrorf(common.CCErrCommXXExceedLimit, "module_property_filter.rules", querybuilder.MaxDeep)
		}
	}

	return nil
}

// ListHostsDetailAndTopoOption TODO
type ListHostsDetailAndTopoOption struct {
	// return the host's topology with biz info.
	WithBiz            bool                      `json:"with_biz"`
	HostPropertyFilter *querybuilder.QueryFilter `json:"host_property_filter"`
	Fields             []string                  `json:"fields"`
	Page               BasePage                  `json:"page"`
}

// Validate TODO
func (option *ListHostsDetailAndTopoOption) Validate() *errors.RawErrorInfo {

	if key, err := option.Page.Validate(false); err != nil {
		return &errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"page." + key},
		}
	}

	if option.Page.Limit > common.BKMaxInstanceLimit {
		return &errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"page." + "limit"},
		}
	}

	if option.HostPropertyFilter != nil {
		if key, err := option.HostPropertyFilter.Validate(&querybuilder.RuleOption{NeedSameSliceElementType: true}); err != nil {
			return &errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{"host_property_filter." + key},
			}
		}
		if option.HostPropertyFilter.GetDeep() > querybuilder.MaxDeep {
			return &errors.RawErrorInfo{
				ErrCode: common.CCErrCommParamsInvalid,
				Args:    []interface{}{"host_property_filter exceeded max allowed deep"},
			}
		}
	}

	if len(option.Fields) == 0 {
		return &errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"fields"},
		}
	}

	return nil
}

// CountTopoNodeHostsOption TODO
type CountTopoNodeHostsOption struct {
	Nodes []TopoNode `json:"topo_nodes" mapstructure:"topo_nodes"`
}

// TimeRange TODO
type TimeRange struct {
	Start *time.Time
	End   *time.Time
}

// ListHosts TODO
type ListHosts struct {
	BizID              int64                     `json:"bk_biz_id,omitempty"`
	SetIDs             []int64                   `json:"bk_set_ids"`
	ModuleIDs          []int64                   `json:"bk_module_ids"`
	HostPropertyFilter *querybuilder.QueryFilter `json:"host_property_filter"`
	Fields             []string                  `json:"fields"`
	Page               BasePage                  `json:"page"`
}

// Validate whether ListHosts is valid
// errKey: invalid key
// er: detail reason why errKey in invalid
func (option ListHosts) Validate() (errKey string, err error) {
	if key, err := option.Page.Validate(false); err != nil {
		return fmt.Sprintf("page.%s", key), err
	}

	if option.HostPropertyFilter != nil {
		if key, err := option.HostPropertyFilter.Validate(&querybuilder.RuleOption{NeedSameSliceElementType: true}); err != nil {
			return fmt.Sprintf("host_property_filter.%s", key), err
		}
		if option.HostPropertyFilter.GetDeep() > querybuilder.MaxDeep {
			return "host_property_filter.rules", fmt.Errorf("exceed max query condition deepth: %d", querybuilder.MaxDeep)
		}
	}

	return "", nil
}

// GetHostPropertyFilter TODO
func (option ListHosts) GetHostPropertyFilter(ctx context.Context) (map[string]interface{}, error) {
	if option.HostPropertyFilter != nil {
		mgoFilter, key, err := option.HostPropertyFilter.ToMgo()
		if err != nil {
			return nil, fmt.Errorf("invalid key:host_property_filter.%s, err: %s", key, err)
		}
		return mgoFilter, nil
	}
	return make(map[string]interface{}), nil
}

// IPInfo TODO
// ip search info
type IPInfo struct {
	Data  []string `json:"data"`
	Exact int64    `json:"exact"`
	Flag  string   `json:"flag"`
}

// SearchCondition TODO
// search condition
type SearchCondition struct {
	Fields    []string        `json:"fields"`
	Condition []ConditionItem `json:"condition"`
	ObjectID  string          `json:"bk_obj_id"`
	// 非必填，只能用来查时间，且与Condition是与关系
	TimeCondition *TimeCondition `json:"time_condition,omitempty"`
}

// SearchHost TODO
type SearchHost struct {
	Count int             `json:"count"`
	Info  []mapstr.MapStr `json:"info"`
}

// ListHostResp list host response
type ListHostResp struct {
	BaseResp `json:",inline"`
	Data     *ListHostResult `json:"data"`
}

// ListHostResult TODO
type ListHostResult struct {
	Count int                      `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

// HostTopoResult TODO
type HostTopoResult struct {
	Count int        `json:"count"`
	Info  []HostTopo `json:"info"`
}

// HostTopo TODO
type HostTopo struct {
	Host map[string]interface{} `json:"host"`
	Topo []Topo                 `json:"topo"`
}

// Topo TODO
type Topo struct {
	SetID   int64    `json:"bk_set_id" bson:"bk_set_id"`
	SetName string   `json:"bk_set_name" bson:"bk_set_name"`
	Module  []Module `json:"module" bson:"module"`
}

// HostDetailWithTopo TODO
type HostDetailWithTopo struct {
	Host map[string]interface{} `json:"host"`
	Topo []*HostTopoNode        `json:"topo"`
}

// HostTopoNode TODO
type HostTopoNode struct {
	Instance *NodeInstance   `json:"inst"`
	Children []*HostTopoNode `json:"children"`
}

// NodeInstance TODO
type NodeInstance struct {
	Object   string      `json:"obj"`
	InstName interface{} `json:"name"`
	InstID   interface{} `json:"id"`
}

// Module TODO
type Module struct {
	ModuleID   int64  `json:"bk_module_id" bson:"bk_module_id" mapstructure:"bk_module_id"`
	ModuleName string `json:"bk_module_name" bson:"bk_module_name" mapstructure:"bk_module_name"`
}

// ExtractHostIDs TODO
func (sh SearchHost) ExtractHostIDs() *[]int64 {
	hostIDArray := make([]int64, 0)
	for _, h := range sh.Info {
		if _, exist := h["host"]; exist == false {
			blog.ErrorJSON("unexpected error, host: %s don't have host field.", h)
			continue
		}
		hostID, exist := h["host"].(mapstr.MapStr)[common.BKHostIDField]
		if exist == false {
			blog.ErrorJSON("unexpected error, host: %s don't have host.bk_host_id field.", h)
			continue
		}
		id, err := util.GetInt64ByInterface(hostID)
		if err != nil {
			blog.ErrorJSON("unexpected error, host: %s host.bk_host_id field is not integer.", h)
			continue
		}
		hostIDArray = append(hostIDArray, id)
	}
	return &hostIDArray
}

// SearchHostResult TODO
type SearchHostResult struct {
	BaseResp `json:",inline"`
	Data     *SearchHost `json:"data"`
}

// HostCloneInputParams TODO
type HostCloneInputParams struct {
	OrgIP  string `json:"bk_org_ip"`
	DstIP  string `json:"bk_dst_ip"`
	AppID  int64  `json:"bk_biz_id"`
	PlatID int64  `json:"bk_cloud_id"`
}

// SetHostConfigParams TODO
type SetHostConfigParams struct {
	ApplicationID int64 `json:"bk_biz_id"`
	SetID         int64 `json:"bk_set_id"`
	ModuleID      int64 `json:"bk_module_id"`
}

// CloneHostPropertyParams TODO
type CloneHostPropertyParams struct {
	AppID int64 `json:"bk_biz_id"`
	// source and destination host inner ip
	OrgIP string `json:"bk_org_ip"`
	DstIP string `json:"bk_dst_ip"`
	// source and destination host id
	OrgID   int64 `json:"bk_org_id"`
	DstID   int64 `json:"bk_dst_id"`
	CloudID int64 `json:"bk_cloud_id"`
}

// TransferHostAcrossBusinessParameter Transfer host across business request parameter
type TransferHostAcrossBusinessParameter struct {
	SrcAppID    int64   `json:"src_bk_biz_id"`
	DstAppID    int64   `json:"dst_bk_biz_id"`
	HostID      []int64 `json:"bk_host_id"`
	DstModuleID int64   `json:"bk_module_id"`
}

// TransferResourceHostAcrossBusinessParam Transfer hosts across business request parameter.
type TransferResourceHostAcrossBusinessParam struct {

	// ResourceSrcHosts source host list.
	ResourceSrcHosts []TransferResourceParam `json:"resource_hosts"`

	// DstAppID destination biz.
	DstAppID int64 `json:"dst_bk_biz_id"`

	// DstModuleID destination module.
	DstModuleID int64 `json:"dst_bk_module_id"`
}

// TransferResourceParam src biz ids and host list.
type TransferResourceParam struct {

	// SrcAppId src biz id.
	SrcAppId int64 `json:"src_bk_biz_id"`

	// HostIDs hosts to be transferred.
	HostIDs []int64 `json:"src_bk_host_ids"`
}

// HostModuleRelationParameter get host and module  relation parameter
type HostModuleRelationParameter struct {
	AppID  int64   `json:"bk_biz_id"`
	HostID []int64 `json:"bk_host_id"`
}

// Validate TODO
func (h *HostModuleRelationParameter) Validate() (rawError errors.RawErrorInfo) {
	if len(h.HostID) == 0 || len(h.HostID) > common.BKMaxInstanceLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrArrayLengthWrong,
			Args:    []interface{}{"bk_host_id", common.BKMaxInstanceLimit},
		}
	}

	return errors.RawErrorInfo{}
}

// DeleteHostFromBizParameter delete host from business
type DeleteHostFromBizParameter struct {
	AppID     int64   `json:"bk_biz_id"`
	HostIDArr []int64 `json:"bk_host_ids"`
}

// CloudAreaSearchParam search cloud area parameter
type CloudAreaSearchParam struct {
	SearchCloudOption `json:",inline"`
	SyncTaskIDs       bool `json:"sync_task_ids"`
}

// CloudAreaHostCount cloud area host count param
type CloudAreaHostCount struct {
	CloudIDs []int64 `json:"bk_cloud_ids"`
}

// Validate TODO
func (c *CloudAreaHostCount) Validate() (rawError errors.RawErrorInfo) {
	maxLimit := 50
	if len(c.CloudIDs) == 0 || len(c.CloudIDs) > maxLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrArrayLengthWrong,
			Args:    []interface{}{"bk_cloud_ids", maxLimit},
		}
	}

	return errors.RawErrorInfo{}
}

// CloudAreaHostCountResult TODO
type CloudAreaHostCountResult struct {
	BaseResp `json:",inline"`
	Data     []CloudAreaHostCountElem `json:"data"`
}

// CloudAreaHostCountElem TODO
type CloudAreaHostCountElem struct {
	CloudID   int64 `json:"bk_cloud_id"`
	HostCount int64 `json:"host_count"`
}

// CreateManyCloudAreaResult TODO
type CreateManyCloudAreaResult struct {
	BaseResp `json:",inline"`
	Data     []CreateManyCloudAreaElem `json:"data"`
}

// CreateManyCloudAreaElem TODO
type CreateManyCloudAreaElem struct {
	CloudID int64  `json:"bk_cloud_id"`
	ErrMsg  string `json:"err_msg"`
}

// TopoNode TODO
type TopoNode struct {
	ObjectID   string `field:"bk_obj_id" json:"bk_obj_id" mapstructure:"bk_obj_id"`
	InstanceID int64  `field:"bk_inst_id" json:"bk_inst_id" mapstructure:"bk_inst_id"`
}

// Key TODO
func (node TopoNode) Key() string {
	return fmt.Sprintf("%s:%d", node.ObjectID, node.InstanceID)
}

// String 用于打印
func (node TopoNode) String() string {
	return fmt.Sprintf("%s:%d", node.ObjectID, node.InstanceID)
}

// TopoNodeHostCount TODO
type TopoNodeHostCount struct {
	Node      TopoNode `field:"topo_node" json:"topo_node" mapstructure:"topo_node"`
	HostCount int      `field:"host_count" json:"host_count" mapstructure:"host_count"`
}

// TransferHostWithAutoClearServiceInstanceOption TODO
type TransferHostWithAutoClearServiceInstanceOption struct {
	HostIDs []int64 `json:"bk_host_ids"`

	RemoveFromModules []int64 `json:"remove_from_modules,omitempty"`
	AddToModules      []int64 `json:"add_to_modules,omitempty"`
	// 主机从 RemoveFromModules 移除后如果不再属于其它模块， 默认转移到空闲机模块
	// DefaultInternalModule 支持调整这种默认行为，可设置成待回收模块或者故障机模块
	DefaultInternalModule int64 `json:"default_internal_module,omitempty"`
	// IsRemoveFromAll if set, remove host from all of its current modules, if not, use remove_from_modules
	IsRemoveFromAll bool `json:"is_remove_from_all"`

	Options TransferOptions `json:"options,omitempty"`
}

// TransferOptions TODO
type TransferOptions struct {
	ServiceInstanceOptions ServiceInstanceOptions `json:"service_instance_options"`

	// HostApplyConflictResolvers update the attribute value of the host with the host as the dimension.
	HostApplyConflictResolvers []HostApplyConflictResolver `json:"host_apply_conflict_resolvers"`

	// HostApplyTransPropertyRule update attributes with the dimension of the rule, which is used to update the host
	// attribute value in the host transfer scenario。
	HostApplyTransPropertyRule HostApplyTransRules `json:"host_apply_trans_rule"`
}

// HostTransferPlan TODO
type HostTransferPlan struct {
	HostID                  int64            `field:"bk_host_id" json:"bk_host_id"`
	FinalModules            []int64          `field:"final_modules" json:"final_modules"`
	ToRemoveFromModules     []int64          `field:"to_remove_from_modules" json:"to_remove_from_modules"`
	ToAddToModules          []int64          `field:"to_add_to_modules" json:"to_add_to_modules"`
	IsTransferToInnerModule bool             `field:"is_transfer_to_inner_module" json:"is_transfer_to_inner_module"`
	HostApplyPlan           OneHostApplyPlan `field:"host_apply_plan" json:"host_apply_plan" mapstructure:"host_apply_plan"`
}

// HostTransferResult transfer host result, contains the transfer status and message
type HostTransferResult struct {
	HostID  int64  `json:"bk_host_id"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RemoveFromModuleInfo TODO
type RemoveFromModuleInfo struct {
	ModuleID         int64             `field:"bk_module_id" json:"bk_module_id"`
	ServiceInstances []ServiceInstance `field:"service_instances" json:"service_instances"`
}

// AddToModuleInfo TODO
type AddToModuleInfo struct {
	ModuleID        int64                  `field:"bk_module_id" json:"bk_module_id"`
	ServiceTemplate *ServiceTemplateDetail `field:"service_template" json:"service_template"`
}

// HostTransferPreview TODO
type HostTransferPreview struct {
	HostID              int64                  `field:"bk_host_id" json:"bk_host_id"`
	FinalModules        []int64                `field:"final_modules" json:"final_modules"`
	ToRemoveFromModules []RemoveFromModuleInfo `field:"to_remove_from_modules" json:"to_remove_from_modules"`
	ToAddToModules      []AddToModuleInfo      `field:"to_add_to_modules" json:"to_add_to_modules"`
	HostApplyPlan       OneHostApplyPlan       `field:"host_apply_plan" json:"host_apply_plan"`
}

// UpdateHostCloudAreaFieldOption TODO
type UpdateHostCloudAreaFieldOption struct {
	BizID   int64   `field:"bk_biz_id" json:"bk_biz_id" mapstructure:"bk_biz_id"`
	HostIDs []int64 `field:"bk_host_ids" json:"bk_host_ids" mapstructure:"bk_host_ids"`
	CloudID int64   `field:"bk_cloud_id" json:"bk_cloud_id" mapstructure:"bk_cloud_id"`
}

// UpdateHostPropertyBatchParameter batch update host property parameter
type UpdateHostPropertyBatchParameter struct {
	Update []UpdateHostProperty `json:"update"`
}

// UpdateHostProperty host property parameter
type UpdateHostProperty struct {
	HostID     int64                  `json:"bk_host_id"`
	Properties map[string]interface{} `json:"properties"`
}

// HostIDArray hostID array struct
type HostIDArray struct {
	HostIDs []int64 `field:"bk_host_ids" json:"bk_host_ids" mapstructure:"bk_host_ids"`
}

type customTopoFilter struct {
	ObjectID string                    `json:"bk_obj_id"`
	Filter   *querybuilder.QueryFilter `json:"filter"`
}

// FindHostTotalTopo find host total topo parameter
type FindHostTotalTopo struct {
	MainlinePropertyFilter []customTopoFilter        `json:"mainline_property_filters"`
	SetPropertyFilter      *querybuilder.QueryFilter `json:"set_property_filter"`
	ModulePropertyFilter   *querybuilder.QueryFilter `json:"module_property_filter"`
	HostPropertyFilter     *querybuilder.QueryFilter `json:"host_property_filter"`
	Fields                 []string                  `json:"fields"`
	Page                   BasePage                  `json:"page"`
}

// Validate validate FindHostTotalTopo params whether correct
func (f *FindHostTotalTopo) Validate(errProxy errors.DefaultCCErrorIf) errors.CCErrorCoder {

	if f.Page.Limit <= 0 {
		return errProxy.CCErrorf(common.CCErrCommParamsNeedSet, "page.limit")
	}
	if f.Page.Limit > common.BKMaxInstanceLimit {
		return errProxy.CCErrorf(common.CCErrCommXXExceedLimit, "page.limit", common.BKMaxInstanceLimit)
	}

	opt := &querybuilder.RuleOption{NeedSameSliceElementType: true}
	for _, objFilter := range f.MainlinePropertyFilter {

		if objFilter.Filter == nil || len(objFilter.ObjectID) == 0 {
			blog.Errorf("get object filter failed, filter is empty or object ID didn't provide")
			return errProxy.CCErrorf(common.CCErrCommParamsInvalid, "mainline_property_filters")
		}

		if key, err := objFilter.Filter.Validate(opt); err != nil {
			return errProxy.CCErrorf(common.CCErrCommParamsInvalid, fmt.Sprintf("%s of %s", key, objFilter.ObjectID))
		}

		if objFilter.Filter.GetDeep() > querybuilder.MaxDeep {
			return errProxy.CCErrorf(common.CCErrCommXXExceedLimit,
				fmt.Sprintf("filter.rule of %s", objFilter.ObjectID), querybuilder.MaxDeep)
		}
	}

	if f.SetPropertyFilter != nil {
		if key, err := f.SetPropertyFilter.Validate(opt); err != nil {
			blog.Errorf("valid set property filter failed, err: %v", err)
			return errProxy.CCErrorf(common.CCErrCommParamsInvalid, fmt.Sprintf("set_property_filter.%s", key))
		}
		if f.SetPropertyFilter.GetDeep() > querybuilder.MaxDeep {
			return errProxy.CCErrorf(common.CCErrCommXXExceedLimit, "set_property_filter.rules",
				querybuilder.MaxDeep)
		}
	}

	if f.ModulePropertyFilter != nil {
		if key, err := f.ModulePropertyFilter.Validate(opt); err != nil {
			blog.Errorf("valid module property filter failed, err: %v", err)
			return errProxy.CCErrorf(common.CCErrCommParamsInvalid, fmt.Sprintf("module_property_filter.%s", key))
		}
		if f.ModulePropertyFilter.GetDeep() > querybuilder.MaxDeep {
			return errProxy.CCErrorf(common.CCErrCommXXExceedLimit, "module_property_filter.rules",
				querybuilder.MaxDeep)
		}
	}

	if f.HostPropertyFilter != nil {
		if key, err := f.HostPropertyFilter.Validate(opt); err != nil {
			return errProxy.CCErrorf(common.CCErrCommParamsInvalid, fmt.Sprintf("%s of %s", key,
				common.BKInnerObjIDHost))
		}

		if f.HostPropertyFilter.GetDeep() > querybuilder.MaxDeep {
			return errProxy.CCErrorf(common.CCErrCommXXExceedLimit, fmt.Sprintf("filter.rule of %s",
				common.BKInnerObjIDHost), querybuilder.MaxDeep)
		}
	}

	return nil
}

// HostMainlineTopoResult result of host mainline topo
type HostMainlineTopoResult struct {
	Count int                  `json:"count"`
	Info  []HostDetailWithTopo `json:"info"`
}

// Validate validate hostIDs length
func (h *HostIDArray) Validate() (rawError errors.RawErrorInfo) {
	if len(h.HostIDs) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsIsInvalid,
			Args:    []interface{}{"bk_host_ids"},
		}
	}

	if len(h.HostIDs) > common.BKMaxSyncIdentifierLimit {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommXXExceedLimit,
			Args:    []interface{}{"bk_host_ids", common.BKMaxSyncIdentifierLimit},
		}
	}
	return errors.RawErrorInfo{}
}

// CountHostCPUReq count host cpu num request
type CountHostCPUReq struct {
	BizID int64     `json:"bk_biz_id,omitempty"`
	Page  *BasePage `json:"page,omitempty"`
}

// Validate CountHostCPUReq
func (c *CountHostCPUReq) Validate() errors.RawErrorInfo {
	if c.BizID == 0 && c.Page == nil {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"bk_biz_id or page"}}
	}

	if c.BizID != 0 && c.Page != nil {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid, Args: []interface{}{"bk_biz_id and page"}}
	}

	if c.BizID < 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsInvalid, Args: []interface{}{"bk_biz_id"}}
	}

	if c.BizID > 0 {
		return errors.RawErrorInfo{}
	}

	if c.Page.Limit == 0 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommParamsNeedSet, Args: []interface{}{"page.limit"}}
	}

	if c.Page.Limit > 10 {
		return errors.RawErrorInfo{ErrCode: common.CCErrCommPageLimitIsExceeded, Args: []interface{}{"page.limit", 10}}
	}

	return errors.RawErrorInfo{}
}

// BizHostCpuCount host cpu count in biz
type BizHostCpuCount struct {
	BizID          int64 `json:"bk_biz_id"`
	HostCount      int64 `json:"host_count"`
	CpuCount       int64 `json:"cpu_count"`
	NoCpuHostCount int64 `json:"no_cpu_host_count"`
}
