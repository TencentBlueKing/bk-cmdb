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
	"configcenter/src/auth/meta"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"net/http"
	"regexp"
)

type AuthConfig struct {
	Name                  string
	Pattern               string
	Regex                 *regexp.Regexp
	HTTPMethod            string
	RequiredBizInMetadata bool
	ResourceType          meta.ResourceType
	ResourceAction        meta.Action
}

func (config *AuthConfig) Match(request *RequestContext) bool {
	if config.HTTPMethod != request.Method {
		return false
	}
	if config.Regex != nil && config.Regex.MatchString(request.URI) == false {
		return false
	}

	return config.Pattern == request.URI
}

var ServiceInstanceAuthConfigs = []AuthConfig{
	{
		Name:                  "createServiceInstancePattern",
		Pattern:               "/api/v3/create/proc/service_instance",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceInstance,
		ResourceAction:        meta.Create,
	}, {
		Name:                  "findServiceInstancePattern",
		Pattern:               "/api/v3/find/proc/service_instance",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceInstance,
		ResourceAction:        meta.Find,
	}, {
		Name:                  "deleteServiceInstancePattern",
		Pattern:               "/api/v3/deletemany/proc/service_instance",
		HTTPMethod:            http.MethodDelete,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceInstance,
		ResourceAction:        meta.Delete,
	}, {
		Name:                  "findServiceInstanceDifferencePattern",
		Pattern:               "/api/v3/find/proc/service_instance/difference",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceInstance,
		ResourceAction:        meta.Find,
	}, {
		Name:                  "syncServiceInstanceAccordingToServiceTemplate",
		Pattern:               "/api/v3/update/proc/service_instance/sync",
		HTTPMethod:            http.MethodPut,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceInstance,
		ResourceAction:        meta.Update,
	}, {
		Name:                  "listServiceInstanceWithHostPattern",
		Pattern:               "/api/v3/findmany/proc/service_instance/with_host",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceInstance,
		ResourceAction:        meta.Find,
	}, {
		Name:                  "addServiceInstanceLabelsPattern",
		Pattern:               "/api/v3/createmany/proc/service_instance/labels",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: false,
		ResourceType:          meta.ProcessServiceInstance,
		ResourceAction:        meta.Find,
	}, {
		Name:                  "removeServiceInstanceLabelsPattern",
		Pattern:               "/api/v3/deletemany/proc/service_instance/labels",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceInstance,
		ResourceAction:        meta.Find,
	}, {
		Name:                  "deleteProcessInstanceInServiceInstanceRegexp",
		Regex:                 regexp.MustCompile(`/api/v3/delete/proc/service_instance/[0-9]+/process/?$`),
		HTTPMethod:            http.MethodDelete,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessServiceInstance,
		ResourceAction:        meta.Delete,
	},
}

func (ps *parseStream) ServiceInstance() *parseStream {
	for _, item := range ServiceInstanceAuthConfigs {
		if item.Match(ps.RequestCtx) == false {
			continue
		}

		var businessID int64
		if item.RequiredBizInMetadata {
			bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
			if err != nil {
				blog.Warnf("get business id in metadata failed, err: %v, rid: %s", err, ps.RequestCtx.Rid)
				ps.err = err
				return ps
			}
			businessID = bizID
		}

		iamResource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Type:   item.ResourceType,
				Action: item.ResourceAction,
			},
			BusinessID: businessID,
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{iamResource}
	}
	return ps
}
