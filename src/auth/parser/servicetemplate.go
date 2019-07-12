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
	"configcenter/src/common"

	"github.com/tidwall/gjson"
)

var ServiceTemplateAuthConfigs = []AuthConfig{
	{
		Name:                  "createServiceTemplatePattern",
		Description:           "创建服务模板",
		Pattern:               "/process/v3/create/proc/service_template",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceTemplate,
		ResourceAction:        meta.Create,
	}, {
		Name:                  "listServiceTemplatePattern",
		Description:           "查询服务模板",
		Pattern:               "/process/v3/findmany/proc/service_template",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceTemplate,
		ResourceAction:        meta.FindMany,
	}, {
		Name:                  "listServiceTemplateDetailPattern",
		Description:           "查询服务模板详情",
		Pattern:               "/process/v3/findmany/proc/service_template/with_detail",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceTemplate,
		ResourceAction:        meta.FindMany,
	}, {
		Name:                  "deleteServiceTemplatePattern",
		Description:           "删除服务模板",
		Pattern:               "/process/v3/delete/proc/service_template",
		HTTPMethod:            http.MethodDelete,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceTemplate,
		ResourceAction:        meta.Delete,
		InstanceIDGetter: func(request *RequestContext) (int64s []int64, e error) {
			templateID := gjson.GetBytes(request.Body, common.BKServiceTemplateIDField).Int()
			if templateID <= 0 {
				return nil, errors.New("invalid service template")
			}
			return []int64{templateID}, nil
		},
	},
}

func (ps *parseStream) ServiceTemplate() *parseStream {
	resources, err := MatchAndGenerateIAMResource(ServiceTemplateAuthConfigs, ps.RequestCtx)
	if err != nil {
		ps.err = err
	}
	if resources != nil {
		ps.Attribute.Resources = resources
	}
	return ps
}
