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
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"configcenter/src/ac/meta"
)

// ProcessInstanceIAMResourceType TODO
var ProcessInstanceIAMResourceType = meta.ProcessServiceInstance

// ProcessInstanceAuthConfigs TODO
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
		// list process instances by biz set regex, authorize by biz set access permission, **only for ui**
		Name:           "listProcessInstancesByBizSetRegexp",
		Description:    "查找业务集下的进程实例",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/proc/biz_set/[0-9]+/process_instance/?$`),
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta.BizSet,
		ResourceAction: meta.AccessBizSet,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			if len(request.Elements) != 7 {
				return nil, fmt.Errorf("get invalid url elements length %d", len(request.Elements))
			}

			bizSetID, err := strconv.ParseInt(request.Elements[5], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("get invalid business set id %s, err: %v", request.Elements[5], err)
			}
			return []int64{bizSetID}, nil
		},
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
		// list process name and id by biz set regex, authorize by biz set access permission, **only for ui**
		Name:           "listProcessInstancesNameIDsInModule",
		Description:    "查询业务集中模块下的进程名和对应的进程ID列表",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/proc/biz_set/[0-9]+/process_instance/name_ids/?$`),
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta.BizSet,
		ResourceAction: meta.AccessBizSet,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			if len(request.Elements) != 8 {
				return nil, fmt.Errorf("get invalid url elements length %d", len(request.Elements))
			}

			bizSetID, err := strconv.ParseInt(request.Elements[5], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("get invalid business set id %s, err: %v", request.Elements[5], err)
			}
			return []int64{bizSetID}, nil
		},
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
	}, {
		// list process instances by ids and biz set regex, authorize by biz set access permission, **only for ui**
		Name:           "listProcessInstancesDetailsByIDsAndBizSetRegexp",
		Description:    "根据进程ID列表批量查询这些进程的详情及关系",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/proc/biz_set/[0-9]+/process_instance/detail/by_ids/?$`),
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
	},
}

// ProcessInstance TODO
func (ps *parseStream) ProcessInstance() *parseStream {
	return ParseStreamWithFramework(ps, ProcessInstanceAuthConfigs)
}
