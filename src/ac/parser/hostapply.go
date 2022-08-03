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

package parser

import (
	"net/http"
	"regexp"

	"configcenter/src/ac/meta"
)

// HostApplyAuthConfigs TODO
var HostApplyAuthConfigs = []AuthConfig{
	{
		Name:           "CreateHostApplyRuleRegex",
		Description:    "添加主机属性自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/create/host_apply_rule/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       5,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.Update,
	}, {
		Name:           "UpdateHostApplyRuleRegex",
		Description:    "更新主机属性自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/update/host_apply_rule/([0-9]+)/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       6,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.Update,
	}, {
		Name:           "DeleteModuleHostApplyRuleRegex",
		Description:    "删除模块场景下主机属性自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/host/deletemany/module/host_apply_rule/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       7,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.Delete,
	}, {
		Name:           "GetHostApplyRuleRegex",
		Description:    "获取主机属性自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/find/host_apply_rule/([0-9]+)/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodGet,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       6,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.DefaultHostApply,
	}, {
		Name:           "ListHostApplyRuleRegex",
		Description:    "列表查询主机属性自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/host_apply_rule/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       5,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.DefaultHostApply,
	}, {
		Name:           "FindmanyModuleHostApplyTaskStatus",
		Description:    "查询模块场景下主机自动应用任务状态",
		Regex:          regexp.MustCompile(`^/api/v3/host/findmany/module/host_apply_plan/status`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.SkipAction,
	}, {
		Name:           "BatchUpdateOrCreateHostApplyRuleRegex",
		Description:    "批量创建/更新主机属性自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/createmany/host_apply_rule/bk_biz_id/([0-9]+)/batch_create_or_update/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       5,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.Update,
	}, {
		Name:           "PreviewModuleApplyHostApplyRulePattern",
		Description:    "预览模块主机属性自动应用",
		Pattern:        "/api/v3/host/createmany/module/host_apply_plan/preview",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.DefaultHostApply,
	}, {
		Name:           "RunHostApplyRuleOnModuleRegex",
		Description:    "模块场景下执行主机属性自动应用",
		Regex:          regexp.MustCompile(`^/api/v3/host/updatemany/module/host_apply_plan/run`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.Update,
	}, {
		Name:           "FindHostRelatedHostApplyRuleRegex",
		Description:    "查询主机关联的主机属性自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/host_apply_rule/bk_biz_id/([0-9]+)/host_related_rules/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       5,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.DefaultHostApply,
	}, {
		Name:           "GetTemplateHostApplyStatusPattern",
		Description:    "获取模块所属服务模板是否开启了主机自动应用",
		Pattern:        "/api/v3/host/find/service_template/host_apply_status",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.DefaultHostApply,
	}, {
		Name:           "GetServiceTemplateHostApplyRulePattern",
		Description:    "获取服务模板主机自动应用规则",
		Pattern:        "/api/v3/host/findmany/service_template/host_apply_rule",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.DefaultHostApply,
	}, {
		Name:           "PreviewServiceTemplateApplyHostApplyRulePattern",
		Description:    "预览服务模版主机属性自动应用",
		Pattern:        "/api/v3/host/createmany/service_template/host_apply_plan/preview",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.DefaultHostApply,
	}, {
		Name:           "GetModuleInvalidHostCountPattern",
		Description:    "获取模块下的主机自动应用失效主机数量",
		Pattern:        "/api/v3/host/findmany/module/host_apply_plan/invalid_host_count",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.DefaultHostApply,
	}, {
		Name:           "GetServiceTemplateInvalidHostCountPattern",
		Description:    "获取服务模版下的主机自动应用失效主机数量",
		Pattern:        "/api/v3/host/findmany/service_template/host_apply_plan/invalid_host_count",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.DefaultHostApply,
	}, {
		Name:           "GetServiceTemplateHostApplyRuleCountPattern",
		Description:    "获取服务模板配置的规则数",
		Pattern:        "/api/v3/host/findmany/service_template/host_apply_rule_count",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.DefaultHostApply,
	}, {
		Name:           "GetModuleFinalRulesPattern",
		Description:    "获取模块的最终属性自动应用规则",
		Pattern:        "/api/v3/host/findmany/module/get_module_final_rules",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.DefaultHostApply,
	},
}

// HostApply TODO
func (ps *parseStream) HostApply() *parseStream {
	return ParseStreamWithFramework(ps, HostApplyAuthConfigs)
}
