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
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"configcenter/src/auth/meta"
)

func (ps *parseStream) cloudRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.CloudAccount()

	return ps
}

var CloudAccountConfigs = []AuthConfig{
	{
		Name:           "verifyCloudAccountPattern",
		Description:    "测试云账户连通性",
		Pattern:        "/api/v3/cloud/account/verify",
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta.CloudAccount,
		ResourceAction: meta.SkipAction,
	}, {
		Name:           "listCloudAccountPattern",
		Description:    "查询云账户",
		Pattern:        "/api/v3/findmany/cloud/account",
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta.CloudAccount,
		ResourceAction: meta.Find,
	}, {
		Name:           "createCloudAccountPattern",
		Description:    "创建云账户",
		Pattern:        "/api/v3/create/cloud/account",
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta.CloudAccount,
		ResourceAction: meta.Create,
	}, {
		Name:           "updateCloudAccountRegex",
		Description:    "更新云账户",
		Regex:          regexp.MustCompile(`^/api/v3/update/cloud/account/([0-9]+)$`),
		HTTPMethod:     http.MethodPut,
		ResourceType:   meta.CloudAccount,
		ResourceAction: meta.Update,
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
		ResourceType:   meta.CloudAccount,
		ResourceAction: meta.Delete,
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

func (ps *parseStream) CloudAccount() *parseStream {
	return ParseStreamWithFramework(ps, CloudAccountConfigs)
}
