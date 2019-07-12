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
	"net/http"

	"configcenter/src/auth/meta"
	"github.com/tidwall/gjson"
)

var ServiceCategoryAuthConfigs = []AuthConfig{
	{
		Name:                  "findmanyServiceCategoryPattern",
		Description:           "list 服务分类",
		Pattern:               "/api/v3/findmany/proc/service_category",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceCategory,
		ResourceAction:        meta.FindMany,
	}, {
		Name:                  "createServiceCategoryPattern",
		Description:           "创建服务分类",
		Pattern:               "/api/v3/create/proc/service_category",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceCategory,
		ResourceAction:        meta.Create,
	}, {
		Name:                  "deleteServiceCategoryPattern",
		Description:           "删除服务分类",
		Pattern:               "/api/v3/delete/proc/service_category",
		HTTPMethod:            http.MethodDelete,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceCategory,
		ResourceAction:        meta.Delete,
		InstanceIDGetter: func(request *RequestContext) ([]int64, error) {
			categoryID := gjson.GetBytes(request.Body, "id").Int()
			if categoryID <= 0 {
				return nil, errors.New("invalid category id")
			}
			return []int64{categoryID}, nil
		},
	},
}

func (ps *parseStream) ServiceCategory() *parseStream {
	resources, err := MatchAndGenerateIAMResource(ServiceCategoryAuthConfigs, ps.RequestCtx)
	if err != nil {
		ps.err = err
	}
	if resources != nil {
		ps.Attribute.Resources = resources
	}
	return ps
}
