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
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"configcenter/src/ac/meta"
)

func (ps *parseStream) cloudRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.CloudResourceDirectory()

	return ps
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
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ResourcePoolDirectory,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// 创建主机池目录
	if ps.hitPattern(createCloudResourceDirectoryPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ResourcePoolDirectory,
					Action: meta.Create,
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

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ResourcePoolDirectory,
					Action:     meta.Update,
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

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ResourcePoolDirectory,
					Action:     meta.Delete,
					InstanceID: dirID,
				},
			},
		}
		return ps
	}

	return ps
}
