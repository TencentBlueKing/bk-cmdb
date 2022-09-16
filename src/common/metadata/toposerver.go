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
	"errors"
	"fmt"

	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/querybuilder"
)

// SearchInstResult TODO
type SearchInstResult struct {
	BaseResp `json:",inline"`
	Data     InstResult `json:"data"`
}

// AppBasicInfoResult TODO
type AppBasicInfoResult struct {
	BaseResp
	Data BizBasicInfo `json:"data"`
}

// CreateModelResult TODO
type CreateModelResult struct {
	BaseResp `json:",inline"`
	Data     Object `json:"data"`
}

// SearchModelResult TODO
type SearchModelResult struct {
	BaseResp `json:",inline"`
	Data     []Object `json:"data"`
}

// SearchInnterAppTopoResult TODO
type SearchInnterAppTopoResult struct {
	BaseResp `json:",inline"`
	Data     InnterAppTopo
}

// MainlineObjectTopoResult TODO
type MainlineObjectTopoResult struct {
	BaseResp `json:",inline"`
	Data     []MainlineObjectTopo `json:"data"`
}

// CommonInstTopo TODO
type CommonInstTopo struct {
	InstNameAsst
	Count    int            `json:"count"`
	Children []InstNameAsst `json:"children"`
}

// CommonInstTopoV2 TODO
type CommonInstTopoV2 struct {
	Prev []*CommonInstTopo `json:"prev"`
	Next []*CommonInstTopo `json:"next"`
	Curr interface{}       `json:"curr"`
}

// SearchAssociationTopoResult TODO
type SearchAssociationTopoResult struct {
	BaseResp `json:",inline"`
	Data     []CommonInstTopoV2 `json:"data"`
}

// SearchTopoResult TODO
type SearchTopoResult struct {
	BaseResp `json:",inline"`
	Data     []*CommonInstTopo `json:"data"`
}

// QueryBusinessRequest TODO
type QueryBusinessRequest struct {
	Fields []string `json:"fields"`
	Page   BasePage `json:"page"`
	// Deprecated: parameters condition and biz_property_filter cannot be set at the same time,
	// The field condition is not
	// maintained later, it is recommended to use the field biz_property_filter.
	Condition         mapstr.MapStr             `json:"condition"`
	BizPropertyFilter *querybuilder.QueryFilter `json:"biz_property_filter,omitempty"`
}

// PreviewBusinessSetRequest 此场景的预览用于用户创建业务集但是还没有点击创建时的中间阶段预览场景。由于此时还没有真正的创建业务集，所以
// 并没有生成业务集实例，只能用范围条件进行获取业务范围。
type PreviewBusinessSetRequest struct {
	BizSetPropertyFilter *querybuilder.QueryFilter `json:"bk_scope"`
	Filter               *querybuilder.QueryFilter `json:"filter"`
	Page                 BasePage                  `json:"page"`
}

// Validate validates preview business set info conditions format.
func (option *PreviewBusinessSetRequest) Validate(allowNoLimit bool) ccErr.RawErrorInfo {

	if errInfo, err := option.Page.Validate(allowNoLimit); err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{errInfo},
		}
	}

	if option.BizSetPropertyFilter == nil && option.Filter == nil {
		return ccErr.RawErrorInfo{}
	}
	op := &querybuilder.RuleOption{
		NeedSameSliceElementType: true,
		MaxSliceElementsCount:    querybuilder.DefaultMaxSliceElementsCount,
		MaxConditionOrRulesCount: querybuilder.DefaultMaxConditionOrRulesCount,
	}

	if option.BizSetPropertyFilter != nil {
		if key, err := option.BizSetPropertyFilter.Validate(op); err != nil {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{(key)},
			}
		}

		if option.BizSetPropertyFilter.GetDeep() > querybuilder.MaxDeep {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommXXExceedLimit,
				Args:    []interface{}{},
			}
		}
	}
	if option.Filter != nil {
		if key, err := option.Filter.Validate(op); err != nil {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsNeedSet,
				Args:    []interface{}{(key)},
			}
		}

		if option.Filter.GetDeep() > common.BizSetConditionMaxDeep {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommXXExceedLimit,
				Args:    []interface{}{},
			}
		}
	}

	return ccErr.RawErrorInfo{}
}

// QueryBusinessSetRequest query business set by query builder
type QueryBusinessSetRequest struct {
	BizSetPropertyFilter *querybuilder.QueryFilter `json:"bk_biz_set_filter,omitempty"`
	// 非必填，只能用来查时间，且与Condition是与关系
	TimeCondition *TimeCondition `json:"time_condition,omitempty"`
	Fields        []string       `json:"fields,omitempty"`
	Page          BasePage       `json:"page,omitempty"`
}

// QueryBusinessSetResponse query business response
type QueryBusinessSetResponse struct {
	Info  []interface{} `json:"info"`
	Count int           `json:"count"`
}

// Validate validates query business set info conditions format.
func (option *QueryBusinessSetRequest) Validate(allowNoLimit bool) ccErr.RawErrorInfo {

	if errInfo, err := option.Page.Validate(allowNoLimit); err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{errInfo},
		}
	}

	if option.BizSetPropertyFilter == nil {
		return ccErr.RawErrorInfo{}
	}

	op := &querybuilder.RuleOption{
		NeedSameSliceElementType: true,
		MaxSliceElementsCount:    querybuilder.DefaultMaxSliceElementsCount,
		MaxConditionOrRulesCount: querybuilder.DefaultMaxConditionOrRulesCount,
	}

	if _, err := option.BizSetPropertyFilter.Validate(op); err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{fmt.Sprintf("biz set property filter is illegal, err: %s", err.Error())},
		}
	}
	if option.BizSetPropertyFilter.GetDeep() > querybuilder.MaxDeep {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args: []interface{}{
				fmt.Sprintf("exceed max filter condition deepth: %d", querybuilder.MaxDeep),
			},
		}
	}
	return ccErr.RawErrorInfo{}
}

// BriefBizRelations TODO
type BriefBizRelations struct {
	Business   interface{} `json:"bk_biz_id"`
	SrcInstID  interface{} `json:"src_id"`
	DestInstID interface{} `json:"dest_id"`
}

// GetBriefBizRelationOptions TODO
type GetBriefBizRelationOptions struct {
	// the source object, which should be one of the mainline object, except host model.
	SrcBizObj string `json:"src_biz_obj"`
	// instance ids belongs to the source object
	SrcInstIDs []int64 `json:"src_ids"`
	// the destination object which is directly associated with the source object, and should be one of the mainline
	// object, except host model.
	// SrcBizObj and DestBizObj can not be same.
	DestBizObj string   `json:"dest_biz_obj"`
	Page       BasePage `json:"page"`
}

// Validate validate get brief biz relation option.
func (o *GetBriefBizRelationOptions) Validate() ccErr.RawErrorInfo {
	if len(o.SrcBizObj) == 0 || o.SrcBizObj == common.BKInnerObjIDHost {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"src_biz_obj"},
		}
	}

	if len(o.SrcInstIDs) == 0 || len(o.SrcInstIDs) > 200 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"src_inst_ids"},
		}
	}

	if len(o.DestBizObj) == 0 || o.DestBizObj == common.BKInnerObjIDHost {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"dest_biz_obj is host"},
		}
	}

	if o.SrcBizObj == o.DestBizObj {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"src_inst_ids or dest_biz_obj is same "},
		}
	}

	if len(o.Page.Sort) != 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"page.sort should not set"},
		}
	}

	if err := o.Page.ValidateLimit(500); err != nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"page.limit"},
		}
	}

	return ccErr.RawErrorInfo{ErrCode: 0}
}

// UpdateBusinessStatusOption TODO
type UpdateBusinessStatusOption struct {
	BizName string `json:"bk_biz_name" mapstructure:"bk_biz_name"`
}

// SearchResourceDirParams TODO
type SearchResourceDirParams struct {
	Fields    []string      `json:"fields"`
	Page      BasePage      `json:"page"`
	Condition mapstr.MapStr `json:"condition"`
	IsFuzzy   bool          `json:"is_fuzzy"`
}

// SearchResourceDirResult TODO
type SearchResourceDirResult struct {
	BizID      int64  `json:"bk_biz_id"`
	ModuleID   int64  `json:"bk_module_id"`
	ModuleName string `json:"bk_module_name"`
	SetID      int64  `json:"bk_set_id"`
	HostCount  int64  `json:"host_count"`
}

// SearchBriefBizTopoOption TODO
type SearchBriefBizTopoOption struct {
	BizID        int64    `json:"bk_biz_id"`
	SetFields    []string `json:"set_fields"`
	ModuleFields []string `json:"module_fields"`
	HostFields   []string `json:"host_fields"`
}

// DeleteBizParam delete archived business parameter
type DeleteBizParam struct {
	BizID []int64 `json:"bk_biz_id"`
}

// Validate validates the input param
func (o *SearchBriefBizTopoOption) Validate() (rawError ccErr.RawErrorInfo) {
	if len(o.SetFields) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"set_fields"},
		}
	}

	if len(o.ModuleFields) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"module_fields"},
		}
	}

	if len(o.HostFields) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"host_fields"},
		}
	}

	return ccErr.RawErrorInfo{}
}

// SetTopo TODO
type SetTopo struct {
	Set         map[string]interface{} `json:"set"`
	ModuleTopos []*ModuleTopo          `json:"modules"`
}

// ModuleTopo TODO
type ModuleTopo struct {
	Module map[string]interface{}   `json:"module"`
	Hosts  []map[string]interface{} `json:"hosts"`
}

// SearchBriefBizTopoResult TODO
type SearchBriefBizTopoResult struct {
	BaseResp `json:",inline"`
	Data     []*SetTopo
}

// BizSetScope defines the scope of biz in biz set, can be all biz or specific biz that matches the filter
type BizSetScope struct {
	MatchAll bool                      `json:"match_all" bson:"match_all"`
	Filter   *querybuilder.QueryFilter `json:"filter" bson:"filter,omitempty"`
}

// BizSetScopeField specific conditions of business scope.
type BizSetScopeField struct {
	Field string
	Value interface{}
}

// BizSetScopeParamsInfo  field info
type BizSetScopeParamsInfo struct {
	Operator  querybuilder.Operator
	FieldInfo []BizSetScopeField
}

// Validate 用于创建和更新场景下的对于业务集scope参数的校验，校验scope仅存在两层且子条件是与的关系，返回其中包含的字段用于后续校验
func (scope *BizSetScope) Validate() (*BizSetScopeParamsInfo, error) {

	fieldInfo := new(BizSetScopeParamsInfo)
	if scope.MatchAll {
		if scope.Filter != nil {
			return nil, errors.New("when match_all is true, params filter can not be set")
		}
		return nil, nil
	}
	if scope.Filter == nil {
		return nil, errors.New("when match_all is false, params filter must be set")
	}

	option := &querybuilder.RuleOption{
		NeedSameSliceElementType: true,
		MaxSliceElementsCount:    querybuilder.DefaultMaxSliceElementsCount,
		MaxConditionOrRulesCount: querybuilder.DefaultMaxConditionOrRulesCount,
	}

	if invalidKey, err := scope.Filter.Validate(option); err != nil {
		return nil, fmt.Errorf("conditions.%s, err: %s", invalidKey, err.Error())
	}

	if scope.Filter.GetDeep() != common.BizSetConditionMaxDeep {
		return nil, fmt.Errorf("scope condition depth must be equal to: %d", common.BizSetConditionMaxDeep)
	}
	if scope.Filter.Rule == nil {
		return fieldInfo, nil
	}

	// 此场景下仅支持 CombinedRule 类型
	if _, ok := scope.Filter.Rule.(querybuilder.CombinedRule); !ok {
		return nil, errors.New("query filter must be combined rules")
	}
	qf := scope.Filter.Rule.(querybuilder.CombinedRule)
	if qf.Condition != querybuilder.ConditionAnd {
		return nil, errors.New("scope condition must be and")
	}

	// 由于只支持2层，所以可以直接获取rules中的field字段
	for _, rule := range qf.Rules {
		if _, ok := rule.(querybuilder.AtomRule); !ok {
			return nil, errors.New("rule type must be AtomRule")
		}
		r := rule.(querybuilder.AtomRule)
		// 仅支持 equal 和 in 操作符
		if r.Operator != querybuilder.OperatorIn && r.Operator != querybuilder.OperatorEqual {
			return nil, errors.New("scope operator must be equal or in")
		}
		fieldInfo.Operator = r.Operator
		fieldInfo.FieldInfo = append(fieldInfo.FieldInfo, BizSetScopeField{
			Value: r.Value,
			Field: r.Field,
		})
	}
	return fieldInfo, nil
}

// CreateBizSetRequest biz set struct
type CreateBizSetRequest struct {
	BizSetAttr  map[string]interface{} `json:"bk_biz_set_attr"`
	BizSetScope *BizSetScope           `json:"bk_scope"`
}

// Validate validates create biz set params
func (op *CreateBizSetRequest) Validate() (*BizSetScopeParamsInfo, ccErr.RawErrorInfo) {
	if op.BizSetAttr == nil || op.BizSetScope == nil {
		return nil, ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"params is nil"},
		}
	}
	if name, ok := op.BizSetAttr[common.BKBizSetNameField]; !ok || name == "" {
		return nil, ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"biz set name must be set"},
		}
	}

	fieldInfo, err := op.BizSetScope.Validate()
	if err != nil {
		return nil, ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{err},
		}
	}

	return fieldInfo, ccErr.RawErrorInfo{}
}

// CreateBizSetResponse create biz set response, returns the created biz set id
type CreateBizSetResponse struct {
	BaseResp `json:",inline"`
	Data     int64 `json:"data"`
}

// GetHostAndSerInstCountResult 获取topo节点主机/服务实例数量结构
type GetHostAndSerInstCountResult struct {
	BaseResp `json:",inline"`
	Data     []*TopoNodeCount
}

// TopoNodeCount 获取topo节点主机/服务实例数量结果
type TopoNodeCount struct {
	ObjID                string `json:"bk_obj_id"`
	InstID               int64  `json:"bk_inst_id"`
	HostCount            int64  `json:"host_count"`
	ServiceInstanceCount int64  `json:"service_instance_count"`
}

// FindBizInBizSetOption find all biz id and name list in biz set option
type FindBizInBizSetOption struct {
	BizSetID int64                     `json:"bk_biz_set_id"`
	Filter   *querybuilder.QueryFilter `json:"filter,omitempty"`
	Fields   []string                  `json:"fields"`
	Page     BasePage                  `json:"page"`
}

// FindBizSetTopoOption find biz set topo nodes option
type FindBizSetTopoOption struct {
	BizSetID    int64  `json:"bk_biz_set_id"`
	ParentObjID string `json:"bk_parent_obj_id"`
	ParentID    int64  `json:"bk_parent_id"`
}

// Validate validate the input option, check if all required fields are set
func (opt *FindBizSetTopoOption) Validate() ccErr.RawErrorInfo {
	if opt.BizSetID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKBizSetIDField},
		}
	}

	if opt.ParentID == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{common.BKParentIDField},
		}
	}

	if len(opt.ParentObjID) == 0 {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"bk_parent_obj_id"},
		}
	}

	return ccErr.RawErrorInfo{}
}

// DeleteBizSetOption delete business set option
type DeleteBizSetOption struct {
	BizSetIDs []int64 `json:"bk_biz_set_ids"`
}

// UpdateBizSetOption update business set option
type UpdateBizSetOption struct {
	BizSetIDs []int64           `json:"bk_biz_set_ids"`
	Data      *UpdateBizSetData `json:"data"`
}

// UpdateBizSetData update business set data
type UpdateBizSetData struct {
	BizSetAttr mapstr.MapStr `json:"bk_biz_set_attr"`
	Scope      *BizSetScope  `json:"bk_scope"`
}
