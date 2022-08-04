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

	"configcenter/src/ac/meta"
)

var (
	searchServiceInstancesBySetTemplateRuleRegex = regexp.MustCompile(`^/api/v3/findmany/proc/service/set_template/list_service_instance/biz/([0-9]+)/?$`)
)

// ServiceInstanceAuthConfigs TODO
var ServiceInstanceAuthConfigs = []AuthConfig{
	{
		Name:           "createServiceInstancePattern",
		Description:    "创建服务实例",
		Pattern:        "/api/v3/create/proc/service_instance",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Create,
	}, {
		Name:           "createServiceInstancePreviewPattern",
		Description:    "创建服务实例预览",
		Pattern:        "/api/v3/create/proc/service_instance/preview",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Create,
	}, {
		Name:           "searchHostWithNoServiceInstancePattern",
		Description:    "获取无服务实例的主机",
		Pattern:        "/api/v3/findmany/proc/host/with_no_service_instance",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostInstance,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "findServiceInstancePattern",
		Description:    "list 服务实例",
		Pattern:        "/api/v3/findmany/proc/service_instance",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "findServiceInstanceWebPattern",
		Description:    "list 服务实例",
		Pattern:        "/api/v3/findmany/proc/web/service_instance",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.FindMany,
	}, {
		// search service instance by biz set regex, authorize by biz set access permission, **only for ui**
		Name:           "findServiceInstanceWebByBizSetRegexp",
		Description:    "UI查询业务集下的服务实例",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/proc/web/biz_set/[0-9]+/service_instance/?$`),
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta.BizSet,
		ResourceAction: meta.AccessBizSet,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			if len(request.Elements) != 8 {
				return nil, fmt.Errorf("get invalid url elements length %d", len(request.Elements))
			}

			bizSetID, err := strconv.ParseInt(request.Elements[6], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("get invalid business set id %s, err: %v", request.Elements[6], err)
			}
			return []int64{bizSetID}, nil
		},
	},
	{
		Name:           "findServiceInstanceDetailsPattern",
		Description:    "list 服务实例详情",
		Pattern:        "/api/v3/findmany/proc/service_instance/details",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "findServiceInstanceByHostPattern",
		Description:    "根据主机服务实例",
		Pattern:        "/api/v3/findmany/proc/service_instance/with_host",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "findServiceInstanceByHostWebPattern",
		Description:    "根据主机服务实例-frontend",
		Pattern:        "/api/v3/findmany/proc/web/service_instance/with_host",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.FindMany,
	}, {
		// search service instance by biz set regex, authorize by biz set access permission, **only for ui**
		Name:           "uiFindServiceInstanceByHostAndBizSetRegexp",
		Description:    "根据主机查询业务集下的服务实例-frontend",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/proc/web/biz_set/[0-9]+/service_instance/with_host/?$`),
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta.BizSet,
		ResourceAction: meta.AccessBizSet,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			if len(request.Elements) != 9 {
				return nil, fmt.Errorf("get invalid url elements length %d", len(request.Elements))
			}

			bizSetID, err := strconv.ParseInt(request.Elements[6], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("get invalid business set id %s, err: %v", request.Elements[6], err)
			}
			return []int64{bizSetID}, nil
		},
	}, {
		Name:           "updateServiceInstances",
		Description:    "更新某业务下的服务实例",
		Regex:          regexp.MustCompile(`^/api/v3/updatemany/proc/service_instance/biz/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       6,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.UpdateMany,
	}, {
		Name:        "updateServiceTemplateHostApplyEnableStatus",
		Description: "更新服务模板主机自动应用状态",
		// NOCC:tosa/linelength(ignore length)
		Regex:          regexp.MustCompile(`^/api/v3/updatemany/proc/service_template/host_apply_enable_status/biz/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       7,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.Update,
	}, {
		Name:           "findmanyServiceTemplateHostApplyTaskStatus",
		Description:    "查看服务模板场景下主机自动应用任务状态",
		Regex:          regexp.MustCompile(`/api/v3/findmany/proc/service_template/host_apply_plan/status`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.SkipAction,
	}, {
		Name:           "deleteServiceTemplateHostApplyRule",
		Description:    "删除服务模板场景下主机自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/deletemany/proc/service_template/host_apply_rule/biz/([0-9]+)/?$`),
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       7,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.Delete,
	}, {
		Name:           "updateServiceTemplateHostApplyRule",
		Description:    "编辑服务模板场景下主机自动应用规则",
		Regex:          regexp.MustCompile(`^/api/v3/updatemany/proc/service_template/host_apply_plan/run`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.HostApply,
		ResourceAction: meta.Update,
	}, {
		Name:           "deleteServiceInstancePattern",
		Description:    "删除服务实例",
		Pattern:        "/api/v3/deletemany/proc/service_instance",
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Delete,
	}, {
		Name:           "diffServiceInstanceWithTemplatePattern",
		Description:    "对比服务实例与模板差异涉及到的进程列表",
		Pattern:        "/api/v3/find/proc/service_template/general_difference",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Find,
	}, {
		Name:           "diffServiceInstanceWithTemplatePattern",
		Description:    "对比服务实例与模板差异涉及到的服务实例",
		Pattern:        "/api/v3/find/proc/difference/service_instances",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Find,
	}, {
		Name:           "diffServiceInstanceWithTemplatePattern",
		Description:    "对比服务实例与模板涉及到的服务实例详细信息",
		Pattern:        "/api/v3/find/proc/service_instance/difference_detail",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Find,
	}, {
		Name:           "syncServiceInstanceAccordingToServiceTemplate",
		Description:    "用服务模板更新服务实例",
		Pattern:        "/api/v3/update/proc/service_instance/sync",
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Update,
	}, {
		Name:           "findServiceTemplateSyncStatus",
		Description:    "获取服务模板同步状态",
		Regex:          regexp.MustCompile(`/api/v3/findmany/proc/service_template_sync_status/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       6,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.Find,
	}, {
		Name:           "listServiceInstanceWithHostPattern",
		Description:    "根据主机查询服务实例",
		Pattern:        "/api/v3/findmany/proc/service_instance/with_host",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "aggregationServiceInstanceLabels",
		Description:    "聚合服务实例labels",
		Pattern:        "/api/v3/findmany/proc/service_instance/labels/aggregation",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.FindMany,
	}, {
		// aggregate service instance labels by biz set regex, authorize by biz set access permission, **only for ui**
		Name:        "aggregationServiceInstanceLabelsByBizSetRegexp",
		Description: "聚合业务集下的服务实例标签",
		Regex: regexp.MustCompile(
			`^/api/v3/findmany/proc/biz_set/[0-9]+/service_instance/labels/aggregation/?$`),
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta.BizSet,
		ResourceAction: meta.AccessBizSet,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			if len(request.Elements) != 9 {
				return nil, fmt.Errorf("get invalid url elements length %d", len(request.Elements))
			}

			bizSetID, err := strconv.ParseInt(request.Elements[5], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("get invalid business set id %s, err: %v", request.Elements[5], err)
			}
			return []int64{bizSetID}, nil
		},
	}, {
		Name:           "addServiceInstanceLabelsPattern",
		Description:    "服务实例添加label",
		Pattern:        "/api/v3/createmany/proc/service_instance/labels",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Update,
	}, {
		Name:           "removeServiceInstanceLabelsPattern",
		Description:    "服务实例删除label",
		Pattern:        "/api/v3/deletemany/proc/service_instance/labels",
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Delete,
	}, {
		Name:           "updateServiceInstanceLabelsPattern",
		Description:    "服务实例更新label",
		Pattern:        "/api/v3/updatemany/proc/service_instance/labels",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Update,
	}, {
		Name:           "deleteProcessInstanceInServiceInstanceRegexp",
		Description:    "删除进程实例",
		Regex:          regexp.MustCompile(`/api/v3/delete/proc/service_instance/[0-9]+/process/?$`),
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Delete,
	}, {
		Name:         "deleteServiceInstancePreviewRegexp",
		Description:  "删除服务实例预览",
		Regex:        regexp.MustCompile(`/api/v3/deletemany/proc/service_instance/preview/?$`),
		HTTPMethod:   http.MethodPost,
		BizIDGetter:  DefaultBizIDGetter,
		ResourceType: meta.ProcessServiceInstance,
		// ResourceAction:        meta.Find,
		ResourceAction: meta.SkipAction,
	},
}

// ServiceInstance TODO
func (ps *parseStream) ServiceInstance() *parseStream {

	if ps.hitRegexp(searchServiceInstancesBySetTemplateRuleRegex, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 9 {
			ps.err = errors.New("search serviceInstance by setTemplate, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[8], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("search serviceInstance by setTemplate, but got invalid business id %s", ps.RequestCtx.Elements[8])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.FindMany,
				},
			},
		}

		return ps
	}

	return ParseStreamWithFramework(ps, ServiceInstanceAuthConfigs)
}
