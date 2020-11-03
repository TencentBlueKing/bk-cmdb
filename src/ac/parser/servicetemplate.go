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
		Name:           "listServiceTemplatePattern",
		Description:    "查询服务模板",
		Pattern:        "/api/v3/findmany/proc/service_template",
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    DefaultBizIDGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "listServiceTemplateDetailPattern",
		Description:    "查询服务模板详情",
		Pattern:        "/api/v3/findmany/proc/service_template/with_detail",
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
	},
}

func (ps *parseStream) ServiceTemplate() *parseStream {
	return ParseStreamWithFramework(ps, ServiceTemplateAuthConfigs)
}
