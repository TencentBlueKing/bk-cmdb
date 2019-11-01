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

	"configcenter/src/auth/meta"
)

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
		Name:           "deleteServiceInstancePattern",
		Description:    "删除服务实例",
		Pattern:        "/api/v3/deletemany/proc/service_instance",
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Delete,
	}, {
		Name:           "diffServiceInstanceWithTemplatePattern",
		Description:    "对比服务实例与模板",
		Pattern:        "/api/v3/find/proc/service_instance/difference",
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
		Name:           "addServiceInstanceLabelsPattern",
		Description:    "服务实例添加label",
		Pattern:        "/api/v3/createmany/proc/service_instance/labels",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Find,
	}, {
		Name:           "removeServiceInstanceLabelsPattern",
		Description:    "服务实例删除label",
		Pattern:        "/api/v3/deletemany/proc/service_instance/labels",
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceInstance,
		ResourceAction: meta.Find,
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

func (ps *parseStream) ServiceInstance() *parseStream {
	return ParseStreamWithFramework(ps, ServiceInstanceAuthConfigs)
}
