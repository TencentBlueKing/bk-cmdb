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
	// parameters condition and biz_property_filter cannot be set at the same time, The field condition is not
	// maintained later, it is recommended to use the field biz_property_filter.
	Condition         mapstr.MapStr             `json:"condition"`
	BizPropertyFilter *querybuilder.QueryFilter `json:"biz_property_filter,omitempty"`
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
