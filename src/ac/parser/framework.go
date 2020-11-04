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
	"regexp"

	"configcenter/src/ac/meta"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
)

type InstanceIDGetter func(request *RequestContext, re *regexp.Regexp) ([]int64, error)
type BizIDGetter func(request *RequestContext, config AuthConfig) (bizID int64, err error)

type AuthConfig struct {
	Name             string
	Pattern          string
	Regex            *regexp.Regexp
	HTTPMethod       string
	ResourceType     meta.ResourceType
	ResourceAction   meta.Action
	InstanceIDGetter InstanceIDGetter
	BizIDGetter      BizIDGetter
	// BizIndex is the index in the request uri elements, used when the BizIDGetter get bizID from url
	BizIndex    int
	Description string
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
		blog.V(4).Infof("match method:%s, pattern:%s, regex:%s", item.HTTPMethod, item.Pattern, item.Regex)

		var bizID int64
		var err error
		if item.BizIDGetter != nil {
			bizID, err = item.BizIDGetter(request, item)
			if err != nil {
				blog.Warnf("get business id failed, name: %s, err: %v, rid: %s", item.Name, err, request.Rid)
				return nil, err
			}
		}

		iamResources := make([]meta.ResourceAttribute, 0)
		if item.InstanceIDGetter == nil {
			iamResource := meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   item.ResourceType,
					Action: item.ResourceAction,
				},
				BusinessID: bizID,
			}
			iamResources = append(iamResources, iamResource)
		} else {
			ids, err := item.InstanceIDGetter(request, item.Regex)
			if err != nil {
				blog.Warnf("get business id failed, name: %s, err: %v, rid: %s", item.Name, err, request.Rid)
				return nil, err
			}
			for _, id := range ids {
				iamResource := meta.ResourceAttribute{
					Basic: meta.Basic{
						Type:       item.ResourceType,
						Action:     item.ResourceAction,
						InstanceID: id,
					},
					BusinessID: bizID,
				}
				iamResources = append(iamResources, iamResource)
			}
		}
		return iamResources, nil
	}
	return nil, nil
}

func DefaultBizIDGetter(request *RequestContext, config AuthConfig) (bizID int64, err error) {
	bizID, err = request.getBizIDFromBody()
	if err != nil {
		return
	}
	return
}

func BizIDFromURLGetter(request *RequestContext, config AuthConfig) (bizID int64, err error) {

	if len(request.Elements) <= config.BizIndex {
		blog.Errorf("invalid BizIndex:%d for uri:%s", config.BizIndex, request.URI)
		return 0, fmt.Errorf("invalid BizIndex:%d for uri:%s", config.BizIndex, request.URI)
	}

	bizIDStr := request.Elements[config.BizIndex]
	bizID, err = util.GetInt64ByInterface(bizIDStr)
	if err != nil {
		blog.Errorf("get business id from request path failed, name: %s, BizIndex:%d, uri:%s, err: %v, rid: %s", config.Name, config.BizIndex, request.URI, err, request.Rid)
		return 0, err
	}

	return bizID, nil
}

func ParseStreamWithFramework(ps *parseStream, authConfigs []AuthConfig) *parseStream {
	resources, err := MatchAndGenerateIAMResource(authConfigs, ps.RequestCtx)
	if err != nil {
		ps.err = err
	}
	if resources != nil {
		ps.Attribute.Resources = resources
	}
	blog.V(7).Infof("ParseStreamWithFramework result: %s", resources)
	return ps
}
