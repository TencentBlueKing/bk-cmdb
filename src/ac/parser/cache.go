/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package parser

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"configcenter/src/ac/meta"
)

func (ps *parseStream) cacheRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.cacheTopology()

	return ps
}

var (
	findBizBriefTopologyRegexp = regexp.MustCompile(`^/api/v3/cache/find/cache/topo/brief/biz/[0-9]+/?$`)
	findBizTopoNodePathRegexp  = regexp.MustCompile(`^/api/v3/cache/find/cache/topo/node_path/biz/[0-9]+/?$`)
)

func (ps *parseStream) cacheTopology() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitRegexp(findBizBriefTopologyRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) < 8 {
			ps.err = errors.New("find biz brief topology, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[8], 10, 64)
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.BizTopology,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(findBizTopoNodePathRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.BizTopology,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	return ps
}
