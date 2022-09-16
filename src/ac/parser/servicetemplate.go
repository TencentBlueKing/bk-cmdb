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
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
)

// ServiceTemplateAuthConfigs TODO
var ServiceTemplateAuthConfigs = []AuthConfig{
	{
		Name:           "createServiceTemplatePattern",
		Description:    "创建服务模板",
		Pattern:        "/api/v3/create/proc/service_template",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.Create,
	}, {
		Name:           "createServiceTemplateAllInfo",
		Description:    "创建服务模板（全量信息）",
		Pattern:        "/api/v3/create/proc/service_template/all_info",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.Create,
	}, {
		Name:           "updateServiceTemplate",
		Description:    "更新服务模板",
		Pattern:        "/api/v3/update/proc/service_template",
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.Update,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			val, err := request.getValueFromBody(common.BKFieldID)
			if err != nil {
				return nil, err
			}
			templateID := val.Int()
			if templateID <= 0 {
				return nil, errors.New("invalid service template")
			}
			return []int64{templateID}, nil
		},
	}, {
		Name:           "updateServiceTemplateAllInfo",
		Description:    "更新服务模板（全量信息）",
		Pattern:        "/api/v3/update/proc/service_template/all_info",
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.Update,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) ([]int64, error) {
			val, err := request.getValueFromBody(common.BKFieldID)
			if err != nil {
				return nil, err
			}

			templateID := val.Int()
			if templateID <= 0 {
				return nil, errors.New("invalid service template id")
			}
			return []int64{templateID}, nil
		},
	}, {
		Name:           "getServiceTemplate",
		Description:    "获取服务模板",
		Regex:          regexp.MustCompile(`^/api/v3/find/proc/service_template/([0-9]+)$`),
		HTTPMethod:     http.MethodGet,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.Find,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			subMatch := re.FindStringSubmatch(request.URI)
			for _, subStr := range subMatch {
				if strings.Contains(subStr, "api") {
					continue
				}
				id, err := strconv.ParseInt(subStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("parse template id to int64 failed, err: %s", err)
				}
				return []int64{id}, nil
			}
			blog.Errorf("unexpected error: this code shouldn't be reached, rid: %s", request.Rid)
			return nil, errors.New("unexpected error: this code shouldn't be reached")
		},
	}, {
		Name:           "getServiceTemplateDetail",
		Description:    "获取服务模板详情",
		Regex:          regexp.MustCompile(`^/api/v3/find/proc/service_template/([0-9]+)/detail$`),
		HTTPMethod:     http.MethodGet,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.Find,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			subMatch := re.FindStringSubmatch(request.URI)
			for _, subStr := range subMatch {
				if strings.Contains(subStr, "api") {
					continue
				}
				id, err := strconv.ParseInt(subStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("parse template id to int64 failed, err: %s", err)
				}
				return []int64{id}, nil
			}
			blog.Errorf("unexpected error: this code shouldn't be reached, rid: %s", request.Rid)
			return nil, errors.New("unexpected error: this code shouldn't be reached")
		},
	}, {
		Name:           "getServiceTemplateAllInfo",
		Description:    "获取服务模板详情（全量信息）",
		Pattern:        "/api/v3/find/proc/service_template/all_info",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.Find,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) ([]int64, error) {
			val, err := request.getValueFromBody(common.BKFieldID)
			if err != nil {
				return nil, err
			}

			templateID := val.Int()
			if templateID <= 0 {
				return nil, errors.New("invalid service template id")
			}
			return []int64{templateID}, nil
		},
	}, {
		Name:           "listServiceTemplatePattern",
		Description:    "查询服务模板",
		Pattern:        "/api/v3/findmany/proc/service_template",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.FindMany,
	}, {
		Name:         "unbindServiceTemplateOnModule",
		Description:  "解绑模块的服务模板",
		Pattern:      "/api/v3/delete/proc/template_binding_on_module",
		HTTPMethod:   http.MethodDelete,
		BizIDGetter:  DefaultBizIDGetter,
		ResourceType: meta.ProcessServiceTemplate,
		// authorization should implements in scene server
		// TODO: implement authorization on scene server
		ResourceAction: meta.SkipAction,
	}, {
		Name:           "deleteServiceTemplatePattern",
		Description:    "删除服务模板",
		Pattern:        "/api/v3/delete/proc/service_template",
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.Delete,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			val, err := request.getValueFromBody(common.BKServiceTemplateIDField)
			if err != nil {
				return nil, err
			}
			templateID := val.Int()
			if templateID <= 0 {
				return nil, errors.New("invalid service template")
			}
			return []int64{templateID}, nil
		},
	}, {
		Name:           "FindServiceTemplateCountInfo",
		Description:    "查询服务模版的计数信息",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/proc/service_template/count_info/biz/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       7,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "GetServiceTemplateSyncStatus",
		Description:    "查询服务模版的同步状态",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/proc/service_template/sync_status/biz/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       7,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.FindMany,
	}, {
		// get service template sync status by biz set regex, authorize by biz set access permission, **only for ui**
		Name:           "GetServiceTemplateSyncStatusByBizSetRegexp",
		Description:    "查询业务集中服务模版的同步状态",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/proc/biz_set/[0-9]+/biz/[0-9]+/service_template/sync_status/?$`),
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta.BizSet,
		ResourceAction: meta.AccessBizSet,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			if len(request.Elements) != 10 {
				return nil, fmt.Errorf("get invalid url elements length %d", len(request.Elements))
			}

			bizSetID, err := strconv.ParseInt(request.Elements[5], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("get invalid business set id %s, err: %v", request.Elements[5], err)
			}
			return []int64{bizSetID}, nil
		},
	}, {
		Name:        "updatemanyServiceTemplateHostApplyEnableStatus",
		Description: "更新服务模板主机自动应用状态",
		// NOCC:tosa/linelength(忽略长度)
		Regex:          regexp.MustCompile(`^/api/v3/updatemany/proc/service_template/host_apply_enable_status/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       7,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.Update,
	},
	{
		Name:           "searchRuleRelatedServiceTemplates",
		Description:    "根据配置字段搜索服务模板",
		Pattern:        "/api/v3/find/proc/service_template/host_apply_rule_related",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "updateServiceTemplateAttribute",
		Description:    "更新服务模板配置字段",
		Pattern:        "/api/v3/update/proc/service_template/attribute",
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.Update,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) ([]int64, error) {
			val, err := request.getValueFromBody(common.BKFieldID)
			if err != nil {
				return nil, err
			}

			templateID := val.Int()
			if templateID <= 0 {
				return nil, errors.New("invalid service template id")
			}
			return []int64{templateID}, nil
		},
	}, {
		Name:           "deleteServiceTemplateAttribute",
		Description:    "删除服务模板配置字段",
		Pattern:        "/api/v3/delete/proc/service_template/attribute",
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.Update,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) ([]int64, error) {
			val, err := request.getValueFromBody(common.BKFieldID)
			if err != nil {
				return nil, err
			}

			templateID := val.Int()
			if templateID <= 0 {
				return nil, errors.New("invalid service template id")
			}
			return []int64{templateID}, nil
		},
	}, {
		Name:           "listServiceTemplateAttribute",
		Description:    "查询服务模板配置字段",
		Pattern:        "/api/v3/findmany/proc/service_template/attribute",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.FindMany,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) ([]int64, error) {
			val, err := request.getValueFromBody(common.BKFieldID)
			if err != nil {
				return nil, err
			}

			templateID := val.Int()
			if templateID <= 0 {
				return nil, errors.New("invalid service template id")
			}
			return []int64{templateID}, nil
		},
	},
}

// ServiceTemplate TODO
func (ps *parseStream) ServiceTemplate() *parseStream {
	return ParseStreamWithFramework(ps, ServiceTemplateAuthConfigs)
}
