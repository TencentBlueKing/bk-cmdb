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
	"regexp"

	"configcenter/src/auth/meta"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

type InstanceIDGetter func(request *RequestContext, re *regexp.Regexp) ([]int64, error)

type AuthConfig struct {
	Name                  string
	Pattern               string
	Regex                 *regexp.Regexp
	HTTPMethod            string
	RequiredBizInMetadata bool
	ResourceType          meta.ResourceType
	ResourceAction        meta.Action
	InstanceIDGetter      InstanceIDGetter
	Description           string
}

func (config *AuthConfig) Match(request *RequestContext) bool {
	if config.HTTPMethod != request.Method {
		return false
	}
	if config.Regex != nil {
		return config.Regex.MatchString(request.URI)
	}

	return config.Pattern == request.URI
}

func MatchAndGenerateIAMResource(authConfigs []AuthConfig, request *RequestContext) ([]meta.ResourceAttribute, error) {
	for _, item := range authConfigs {
		if item.Match(request) == false {
			continue
		}

		var businessID int64
		if item.RequiredBizInMetadata {
			bizID, err := metadata.BizIDFromMetadata(request.Metadata)
			if err != nil {
				blog.Warnf("get business id in metadata failed, name: %s, err: %v, rid: %s", item.Name, err, request.Rid)
				return nil, err
			}
			businessID = bizID
		}

		iamResources := make([]meta.ResourceAttribute, 0)
		if item.InstanceIDGetter == nil {
			iamResource := meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   item.ResourceType,
					Action: item.ResourceAction,
				},
				BusinessID: businessID,
			}
			iamResources = append(iamResources, iamResource)
		} else {
			ids, err := item.InstanceIDGetter(request, item.Regex)
			if err != nil {
				blog.Warnf("get business id in metadata failed, name: %s, err: %v, rid: %s", item.Name, err, request.Rid)
				return nil, err
			}
			for _, id := range ids {
				iamResource := meta.ResourceAttribute{
					Basic: meta.Basic{
						Type:       item.ResourceType,
						Action:     item.ResourceAction,
						InstanceID: id,
					},
					BusinessID: businessID,
				}
				iamResources = append(iamResources, iamResource)
			}
		}
		return iamResources, nil
	}
	return nil, nil
}
