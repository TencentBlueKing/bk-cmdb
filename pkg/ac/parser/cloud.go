/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
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
	meta2 "configcenter/pkg/ac/meta"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"configcenter/pkg/common"
)

func (ps *parseStream) cloudRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.cloudAccount().cloudResourceTask().CloudResourceDirectory()

	return ps
}

var cloudAccountConfigs = []AuthConfig{
	{
		Name:           "verifyCloudAccountPattern",
		Description:    "测试云账户连通性",
		Pattern:        "/api/v3/cloud/account/verify",
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta2.CloudAccount,
		ResourceAction: meta2.SkipAction,
	}, {
		Name:           "searchCloudAccountValidityPattern",
		Description:    "查询云账户有效性",
		Pattern:        "/api/v3/findmany/cloud/account/validity",
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta2.CloudAccount,
		ResourceAction: meta2.SkipAction,
	}, {
		Name:           "listCloudAccountPattern",
		Description:    "查询云账户",
		Pattern:        "/api/v3/findmany/cloud/account",
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta2.CloudAccount,
		ResourceAction: meta2.SkipAction,
	}, {
		Name:           "createCloudAccountPattern",
		Description:    "创建云账户",
		Pattern:        "/api/v3/create/cloud/account",
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta2.CloudAccount,
		ResourceAction: meta2.Create,
	}, {
		Name:           "updateCloudAccountRegex",
		Description:    "更新云账户",
		Regex:          regexp.MustCompile(`^/api/v3/update/cloud/account/([0-9]+)$`),
		HTTPMethod:     http.MethodPut,
		ResourceType:   meta2.CloudAccount,
		ResourceAction: meta2.Update,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			subMatch := re.FindStringSubmatch(request.URI)
			for _, subStr := range subMatch {
				if strings.Contains(subStr, "api") {
					continue
				}
				id, err := strconv.ParseInt(subStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("parse account id to int64 failed, err: %s", err)
				}
				return []int64{id}, nil
			}
			return nil, errors.New("unexpected error: this code shouldn't be reached")
		},
	}, {
		Name:           "deleteCloudAccountRegex",
		Description:    "删除云账户",
		Regex:          regexp.MustCompile(`^/api/v3/delete/cloud/account/([0-9]+)$`),
		HTTPMethod:     http.MethodDelete,
		ResourceType:   meta2.CloudAccount,
		ResourceAction: meta2.Delete,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			subMatch := re.FindStringSubmatch(request.URI)
			for _, subStr := range subMatch {
				if strings.Contains(subStr, "api") {
					continue
				}
				id, err := strconv.ParseInt(subStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("parse account id to int64 failed, err: %s", err)
				}
				return []int64{id}, nil
			}
			return nil, errors.New("unexpected error: this code shouldn't be reached")
		},
	},
}

var cloudResourceTaskConfigs = []AuthConfig{
	{
		Name:           "getCloudAccountVpcRegex",
		Description:    "查询账户下的vpc数据",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/cloud/account/vpc/([0-9]+)$`),
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta2.CloudResourceTask,
		ResourceAction: meta2.SkipAction,
	}, {
		Name:           "listCloudResourceTaskPattern",
		Description:    "查询云资源同步任务",
		Pattern:        "/api/v3/findmany/cloud/sync/task",
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta2.CloudResourceTask,
		ResourceAction: meta2.SkipAction,
	}, {
		Name:           "createCloudResourceTaskPattern",
		Description:    "创建云资源同步任务",
		Pattern:        "/api/v3/create/cloud/sync/task",
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta2.CloudResourceTask,
		ResourceAction: meta2.Create,
	}, {
		Name:           "updateCloudResourceTaskRegex",
		Description:    "更新云资源同步任务",
		Regex:          regexp.MustCompile(`^/api/v3/update/cloud/sync/task/([0-9]+)$`),
		HTTPMethod:     http.MethodPut,
		ResourceType:   meta2.CloudResourceTask,
		ResourceAction: meta2.Update,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			subMatch := re.FindStringSubmatch(request.URI)
			for _, subStr := range subMatch {
				if strings.Contains(subStr, "api") {
					continue
				}
				id, err := strconv.ParseInt(subStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("parse account id to int64 failed, err: %s", err)
				}
				return []int64{id}, nil
			}
			return nil, errors.New("unexpected error: this code shouldn't be reached")
		},
	}, {
		Name:           "deleteCloudResourceTaskRegex",
		Description:    "删除云资源同步任务",
		Regex:          regexp.MustCompile(`^/api/v3/delete/cloud/sync/task/([0-9]+)$`),
		HTTPMethod:     http.MethodDelete,
		ResourceType:   meta2.CloudResourceTask,
		ResourceAction: meta2.Delete,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			subMatch := re.FindStringSubmatch(request.URI)
			for _, subStr := range subMatch {
				if strings.Contains(subStr, "api") {
					continue
				}
				id, err := strconv.ParseInt(subStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("parse account id to int64 failed, err: %s", err)
				}
				return []int64{id}, nil
			}
			return nil, errors.New("unexpected error: this code shouldn't be reached")
		},
	}, {
		Name:           "listCloudResourceTaskHistoryPattern",
		Description:    "查询云资源同步历史记录",
		Pattern:        "/api/v3/findmany/cloud/sync/history",
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta2.CloudResourceTask,
		ResourceAction: meta2.Find,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			val, err := request.getValueFromBody(common.BKCloudTaskID)
			if err != nil {
				return nil, err
			}
			taskID := val.Int()
			if taskID <= 0 {
				return nil, errors.New("invalid cloud sync task id")
			}
			return []int64{taskID}, nil
		},
	},
	{
		Name:           "listCloudResourceRegionPattern",
		Description:    "查询云资源同步地域信息",
		Pattern:        "/api/v3/findmany/cloud/sync/region",
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta2.CloudResourceTask,
		ResourceAction: meta2.SkipAction,
	},
}

func (ps *parseStream) cloudAccount() *parseStream {
	return ParseStreamWithFramework(ps, cloudAccountConfigs)
}

func (ps *parseStream) cloudResourceTask() *parseStream {
	return ParseStreamWithFramework(ps, cloudResourceTaskConfigs)
}

const (
	getCloudResourceDirectoryPattern    = "/api/v3/findmany/resource/directory"
	createCloudResourceDirectoryPattern = "/api/v3/create/resource/directory"
)

var (
	updateCloudResourceDirectoryRegexp = regexp.MustCompile(`^/api/v3/update/resource/directory/([0-9]+)$`)
	deleteCloudResourceDirectoryRegexp = regexp.MustCompile(`^/api/v3/delete/resource/directory/([0-9]+)$`)
)

// CloudResourceDirectory TODO
func (ps *parseStream) CloudResourceDirectory() *parseStream {

	if ps.shouldReturn() {
		return ps
	}

	// "查询主机池目录"
	if ps.hitPattern(getCloudResourceDirectoryPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			{
				Basic: meta2.Basic{
					Type:   meta2.ResourcePoolDirectory,
					Action: meta2.SkipAction,
				},
			},
		}
		return ps
	}

	// 创建主机池目录
	if ps.hitPattern(createCloudResourceDirectoryPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			{
				Basic: meta2.Basic{
					Type:   meta2.ResourcePoolDirectory,
					Action: meta2.Create,
				},
			},
		}
		return ps
	}

	// 更新主机池目录
	if ps.hitRegexp(updateCloudResourceDirectoryRegexp, http.MethodPut) {
		dirID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("parse resource dir id %s failed, err: %v", ps.RequestCtx.Elements[5], err)
			return ps
		}

		ps.Attribute.Resources = []meta2.ResourceAttribute{
			{
				Basic: meta2.Basic{
					Type:       meta2.ResourcePoolDirectory,
					Action:     meta2.Update,
					InstanceID: dirID,
				},
			},
		}
		return ps
	}

	// 删除主机池目录
	if ps.hitRegexp(deleteCloudResourceDirectoryRegexp, http.MethodDelete) {
		dirID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("parse resource dir id %s failed, err: %v", ps.RequestCtx.Elements[5], err)
			return ps
		}

		ps.Attribute.Resources = []meta2.ResourceAttribute{
			{
				Basic: meta2.Basic{
					Type:       meta2.ResourcePoolDirectory,
					Action:     meta2.Delete,
					InstanceID: dirID,
				},
			},
		}
		return ps
	}

	return ps
}
