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

	"configcenter/src/auth/meta"
)

var ProcessInstanceAuthConfigs = []AuthConfig{
	{
		Name:                  "createProcessInstances",
		Description:           "创建进程实例",
		Pattern:               "/create/proc/process_instance",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.Process,
		ResourceAction:        meta.Create,
	}, {
		Name:                  "updateProcessInstances",
		Description:           "更新进程实例",
		Pattern:               "/update/proc/process_instance",
		HTTPMethod:            http.MethodPut,
		RequiredBizInMetadata: true,
		ResourceType:          meta.Process,
		ResourceAction:        meta.Update,
	}, {
		Name:                  "deleteProcessInstance",
		Description:           "删除进程实例",
		Pattern:               "/delete/proc/process_instance",
		HTTPMethod:            http.MethodDelete,
		RequiredBizInMetadata: true,
		ResourceType:          meta.Process,
		ResourceAction:        meta.Delete,
	}, {
		Name:                  "listProcessInstances",
		Description:           "查找进程实例",
		Pattern:               "/findmany/proc/process_instance",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.Process,
		ResourceAction:        meta.Find,
	},
}

func (ps *parseStream) ProcessInstance() *parseStream {
	resources, err := MatchAndGenerateIAMResource(ProcessInstanceAuthConfigs, ps.RequestCtx)
	if err != nil {
		ps.err = err
	}
	if resources != nil {
		ps.Attribute.Resources = resources
	}
	return ps
}
