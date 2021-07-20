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

var ProcessInstanceIAMResourceType = meta.ProcessServiceInstance

var ProcessInstanceAuthConfigs = []AuthConfig{
	{
		Name:           "createProcessInstances",
		Description:    "创建进程实例",
		Pattern:        "/api/v3/create/proc/process_instance",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   ProcessInstanceIAMResourceType,
		ResourceAction: meta.Update,
	}, {
		Name:           "updateProcessInstances",
		Description:    "更新进程实例",
		Pattern:        "/api/v3/update/proc/process_instance",
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   ProcessInstanceIAMResourceType,
		ResourceAction: meta.Update,
	}, {
		Name:           "updateProcessInstancesByIDs",
		Description:    "根据进程ID批量更新进程实例",
		Pattern:        "/api/v3/update/proc/process_instance/by_ids",
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   ProcessInstanceIAMResourceType,
		ResourceAction: meta.Update,
	}, {
		Name:           "deleteProcessInstance",
		Description:    "删除进程实例",
		Pattern:        "/api/v3/delete/proc/process_instance",
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   ProcessInstanceIAMResourceType,
		ResourceAction: meta.Update,
	}, {
		Name:           "listProcessInstances",
		Description:    "查找进程实例",
		Pattern:        "/api/v3/findmany/proc/process_instance",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   ProcessInstanceIAMResourceType,
		ResourceAction: meta.Find,
	}, {
		Name:           "listProcessRelatedInfo",
		Description:    "点分五位查询进程实例相关的信息",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/proc/process_related_info/biz/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       6,
		ResourceType:   ProcessInstanceIAMResourceType,
		ResourceAction: meta.Find,
	}, {
		Name:           "listProcessInstancesNameIDsInModule",
		Description:    "查询模块下的进程名和对应的进程ID列表",
		Pattern:        "/api/v3/findmany/proc/process_instance/name_ids",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   ProcessInstanceIAMResourceType,
		ResourceAction: meta.Find,
	}, {
		Name:           "listProcessInstancesDetails",
		Description:    "查询某业务下进程ID对应的进程详情",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/proc/process_instance/detail/biz/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		BizIndex:       7,
		ResourceType:   ProcessInstanceIAMResourceType,
		ResourceAction: meta.Find,
	}, {
		Name:           "listProcessInstancesDetailsByIDs",
		Description:    "根据进程ID列表批量查询这些进程的详情及关系",
		Pattern:        "/api/v3/findmany/proc/process_instance/detail/by_ids",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   ProcessInstanceIAMResourceType,
		ResourceAction: meta.Find,
	},
}

func (ps *parseStream) ProcessInstance() *parseStream {
	return ParseStreamWithFramework(ps, ProcessInstanceAuthConfigs)
}
