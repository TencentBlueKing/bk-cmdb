/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
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

// ContainerAuthConfigs auth configs related to container
var ContainerAuthConfigs = []AuthConfig{
	{
		Name:           "createPod",
		Description:    "创建Pod",
		Regex:          regexp.MustCompile(`^/api/v3/create/container/bk_biz_id/([0-9])+/pod$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.Pod,
		ResourceAction: meta.Create,
	},
	{
		Name:           "createManyPod",
		Description:    "创建多个Pod",
		Regex:          regexp.MustCompile(`^/api/v3/createmany/container/bk_biz_id/([0-9])+/pod$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.Pod,
		ResourceAction: meta.CreateMany,
	},
	{
		Name:           "updatePod",
		Description:    "更新Pod",
		Regex:          regexp.MustCompile(`^/api/v3/update/container/bk_biz_id/([0-9])+/pod$`),
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.Pod,
		ResourceAction: meta.Update,
	},
	{
		Name:           "deletePod",
		Description:    "删除Pod",
		Regex:          regexp.MustCompile(`^/api/v3/delete/container/bk_biz_id/([0-9])+/pod$`),
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.Pod,
		ResourceAction: meta.Delete,
	},
	{
		Name:           "findManyPod",
		Description:    "查询多个Pod",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/container/bk_biz_id/([0-9])+/pod$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.Pod,
		ResourceAction: meta.FindMany,
	},
}

func (ps *parseStream) Pod() *parseStream {
	return ParseStreamWithFramework(ps, ContainerAuthConfigs)
}

func (ps *parseStream) containerRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.Pod()

	return ps
}
