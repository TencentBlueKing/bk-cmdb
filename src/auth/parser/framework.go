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

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/tidwall/gjson"
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
	Description      string
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

		var bizID int64
		var err error
		if item.BizIDGetter != nil {
			bizID, err = item.BizIDGetter(request, item)
			if err != nil {
				blog.Warnf("get business id in metadata failed, name: %s, err: %v, rid: %s", item.Name, err, request.Rid)
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
	bizID = gjson.GetBytes(request.Body, common.BKAppIDField).Int()
	if bizID != 0 {
		return bizID, nil
	}

	bizID, err = metadata.BizIDFromMetadata(request.Metadata)
	if err != nil {
		blog.Warnf("get business id from metadata failed, name: %s, err: %v, rid: %s", config.Name, err, request.Rid)
		return 0, err
	}
	return bizID, nil
}

var (
	BizIDRegex = regexp.MustCompile("bk_biz_id/([0-9]+)")
)

func BizIDFromURLGetter(request *RequestContext, config AuthConfig) (bizID int64, err error) {
	match := BizIDRegex.FindStringSubmatch(request.URI)
	if len(match) == 0 {
		return 0, fmt.Errorf("url: %s not match regex: %s", request.URI, BizIDRegex)
	}
	bizID, err = util.GetInt64ByInterface(match[1])
	if err != nil {
		blog.Warnf("get business id from request path failed, name: %s, err: %v, rid: %s", config.Name, err, request.Rid)
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
