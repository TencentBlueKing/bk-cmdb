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
		Name:           "DeleteHostApplyRuleRegex",
		Description:    "删除主机属性自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/deletemany/host_apply_rule/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       5,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.Update,
	}, {
		Name:           "GetHostApplyRuleRegex",
		Description:    "获取主机属性自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/find/host_apply_rule/([0-9]+)/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodGet,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       6,
		ResourceType:   meta.MainlineInstanceTopology,
		ResourceAction: meta.SkipAction,
	}, {
		Name:           "ListHostApplyRuleRegex",
		Description:    "列表查询主机属性自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/host_apply_rule/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       5,
		ResourceType:   meta.MainlineInstanceTopology,
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
		Name:           "PreviewApplyHostApplyRuleRegex",
		Description:    "预览主机属性自动应用",
		Regex:          regexp.MustCompile(`^/api/v3/createmany/host_apply_plan/bk_biz_id/([0-9]+)/preview/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       5,
		ResourceType:   meta.MainlineInstanceTopology,
		ResourceAction: meta.SkipAction,
	}, {
		Name:           "RunHostApplyRuleRegex",
		Description:    "执行主机属性自动应用",
		Regex:          regexp.MustCompile(`^/api/v3/updatemany/host_apply_plan/bk_biz_id/([0-9]+)/run/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       5,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.Update,
	}, {
		Name:           "FindHostRelatedHostApplyRuleRegex",
		Description:    "查询主机关联的主机属性自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/host_apply_rule/bk_biz_id/([0-9]+)/host_related_rules/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       5,
		ResourceType:   meta.MainlineInstanceTopology,
		ResourceAction: meta.SkipAction,
	},
}

func (ps *parseStream) HostApply() *parseStream {
	return ParseStreamWithFramework(ps, HostApplyAuthConfigs)
}
