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
	"configcenter/src/common/querybuilder"
	"fmt"
)

type SearchInstResult struct {
	BaseResp `json:",inline"`
	Data     InstResult `json:"data"`
}

type AppBasicInfoResult struct {
	BaseResp
	Data BizBasicInfo `json:"data"`
}

type CreateModelResult struct {
	BaseResp `json:",inline"`
	Data     Object `json:"data"`
}
type SearchModelResult struct {
	BaseResp `json:",inline"`
	Data     []Object `json:"data"`
}

type SearchInnterAppTopoResult struct {
	BaseResp `json:",inline"`
	Data     InnterAppTopo
}

type MainlineObjectTopoResult struct {
	BaseResp `json:",inline"`
	Data     []MainlineObjectTopo `json:"data"`
}

type CommonInstTopo struct {
	InstNameAsst
	Count    int            `json:"count"`
	Children []InstNameAsst `json:"children"`
}

type CommonInstTopoV2 struct {
	Prev []*CommonInstTopo `json:"prev"`
	Next []*CommonInstTopo `json:"next"`
	Curr interface{}       `json:"curr"`
}
type SearchAssociationTopoResult struct {
	BaseResp `json:",inline"`
	Data     []CommonInstTopoV2 `json:"data"`
}

type SearchTopoResult struct {
	BaseResp `json:",inline"`
	Data     []*CommonInstTopo `json:"data"`
}

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
	BizSetPropertyFilter *querybuilder.QueryFilter `json:"bk_biz_set_filter"`
}

// Validate validates preview business set info conditions format.
func (option *PreviewBusinessSetRequest) Validate() error {
	if option.BizSetPropertyFilter == nil {
		return fmt.Errorf("condition must be set")
	}
	op := &querybuilder.RuleOption{
		NeedSameSliceElementType: true,
		MaxSliceElementsCount:    querybuilder.DefaultMaxSliceElementsCount,
		MaxConditionOrRulesCount: querybuilder.DefaultMaxConditionOrRulesCount,
	}

	if key, err := option.BizSetPropertyFilter.Validate(op); err != nil {
		return fmt.Errorf("bizPropertyFilter is illegal,biz.property.%s", key)
	}

	if option.BizSetPropertyFilter.GetDeep() > common.BizSetConditionMaxDeep {
		return fmt.Errorf("exceed max query condition deepth: %d", common.BizSetConditionMaxDeep)
	}
	return nil
}

// QueryBusinessSetRequest query business set by query builder
type QueryBusinessSetRequest struct {
	BizSetPropertyFilter *querybuilder.QueryFilter `json:"bk_biz_set_filter"`
	Fields               []string                  `json:"fields"`

	// Page Limit must be set less than 500
	Page BasePage `json:"page"`
}

// Validate validates query business set info conditions format.
func (option *QueryBusinessSetRequest) Validate() error {

	if option.BizSetPropertyFilter == nil {
		return fmt.Errorf("query params biz set filter must be set")
	}

	op := &querybuilder.RuleOption{
		NeedSameSliceElementType: true,
		MaxSliceElementsCount:    querybuilder.DefaultMaxSliceElementsCount,
		MaxConditionOrRulesCount: querybuilder.DefaultMaxConditionOrRulesCount,
	}

	if _, err := option.BizSetPropertyFilter.Validate(op); err != nil {
		return fmt.Errorf("biz set property filter is illegal, err: %s", err.Error())
	}

	if option.BizSetPropertyFilter.GetDeep() > common.BizSetConditionMaxDeep {
		return fmt.Errorf("exceed max filter condition deepth: %d", common.BizSetConditionMaxDeep)
	}
	if option.Page.IsIllegalWithCount() {
		return fmt.Errorf("limit is illegal, limit:%d", option.Page.Limit)
	}
	return nil
}

type BriefBizRelations struct {
	Business   interface{} `json:"bk_biz_id"`
	SrcInstID  interface{} `json:"src_id"`
	DestInstID interface{} `json:"dest_id"`
}

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

func (o *GetBriefBizRelationOptions) Validate() errors.RawErrorInfo {
	if len(o.SrcBizObj) == 0 || o.SrcBizObj == common.BKInnerObjIDHost {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"src_biz_obj"},
		}
	}

	if len(o.SrcInstIDs) == 0 || len(o.SrcInstIDs) > 200 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"src_inst_ids"},
		}
	}

	if len(o.DestBizObj) == 0 || o.DestBizObj == common.BKInnerObjIDHost {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"dest_biz_obj is host"},
		}
	}

	if o.SrcBizObj == o.DestBizObj {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"src_inst_ids or dest_biz_obj is same "},
		}
	}

	if len(o.Page.Sort) != 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"page.sort should not set"},
		}
	}

	if err := o.Page.ValidateLimit(500); err != nil {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"page.limit"},
		}
	}

	return errors.RawErrorInfo{ErrCode: 0}
}

type UpdateBusinessStatusOption struct {
	BizName string `json:"bk_biz_name" mapstructure:"bk_biz_name"`
}

type SearchResourceDirParams struct {
	Fields    []string      `json:"fields"`
	Page      BasePage      `json:"page"`
	Condition mapstr.MapStr `json:"condition"`
	IsFuzzy   bool          `json:"is_fuzzy"`
}

type SearchResourceDirResult struct {
	BizID      int64  `json:"bk_biz_id"`
	ModuleID   int64  `json:"bk_module_id"`
	ModuleName string `json:"bk_module_name"`
	SetID      int64  `json:"bk_set_id"`
	HostCount  int64  `json:"host_count"`
}

type SearchBriefBizTopoOption struct {
	BizID        int64    `json:"bk_biz_id"`
	SetFields    []string `json:"set_fields"`
	ModuleFields []string `json:"module_fields"`
	HostFields   []string `json:"host_fields"`
}

// DeleteBizParam delete archived business parameter
type DeleteBizParam struct {
	BizID []int64 `json:"bk_biz_id""`
}

// Validate validates the input param
func (o *SearchBriefBizTopoOption) Validate() (rawError errors.RawErrorInfo) {
	if len(o.SetFields) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"set_fields"},
		}
	}

	if len(o.ModuleFields) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"module_fields"},
		}
	}

	if len(o.HostFields) == 0 {
		return errors.RawErrorInfo{
			ErrCode: common.CCErrCommParamsInvalid,
			Args:    []interface{}{"host_fields"},
		}
	}

	return errors.RawErrorInfo{}
}

type SetTopo struct {
	Set         map[string]interface{} `json:"set"`
	ModuleTopos []*ModuleTopo          `json:"modules"`
}

type ModuleTopo struct {
	Module map[string]interface{}   `json:"module"`
	Hosts  []map[string]interface{} `json:"hosts"`
}

type SearchBriefBizTopoResult struct {
	BaseResp `json:",inline"`
	Data     []*SetTopo
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
