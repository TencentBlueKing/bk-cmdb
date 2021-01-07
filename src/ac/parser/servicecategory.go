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

	"configcenter/src/ac/meta"
)

// utility.AddHandler(rest.Action{Verb: , Path: , Handler: ps.UpdateServiceCategory})
var ServiceCategoryAuthConfigs = []AuthConfig{
	{
		Name:           "findmanyServiceCategoryPattern",
		Description:    "list 服务分类",
		Pattern:        "/api/v3/findmany/proc/service_category",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceCategory,
		ResourceAction: meta.Find,
	}, {
		Name:           "findmanyServiceCategoryPattern",
		Description:    "list 服务分类(含引用统计)",
		Pattern:        "/api/v3/findmany/proc/service_category/with_statistics",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceCategory,
		ResourceAction: meta.Find,
	}, {
		Name:           "createServiceCategoryPattern",
		Description:    "创建服务分类",
		Pattern:        "/api/v3/create/proc/service_category",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceCategory,
		ResourceAction: meta.Create,
	}, {
		Name:           "deleteServiceCategoryPattern",
		Description:    "修改服务分类",
		Pattern:        "/api/v3/update/proc/service_category",
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceCategory,
		ResourceAction: meta.Update,
	}, {
		Name:           "deleteServiceCategoryPattern",
		Description:    "删除服务分类",
		Pattern:        "/api/v3/delete/proc/service_category",
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceCategory,
		ResourceAction: meta.Delete,
	},
}

func (ps *parseStream) ServiceCategory() *parseStream {
	return ParseStreamWithFramework(ps, ServiceCategoryAuthConfigs)
}
