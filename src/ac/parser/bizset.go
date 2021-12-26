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
	"configcenter/src/ac/meta"
	"net/http"
	"regexp"
)

var (
	createBusinessSetRegexp           = regexp.MustCompile(`^/api/v3/biz_set/[^\s/]+/?$`)
	findmanyBusinessSetRegexp         = `api/v3/findmany/biz_set`
	findReducedBusinessSetListPattern = `/api/v3/findmany/biz_set/with_reduced`
	previewBusinessSet                = `/api/v3/preview/biz_set`
)

func (ps *parseStream) businessSet() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create business set, this is not a normalize api.
	if ps.hitRegexp(createBusinessSetRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.BusinessSet,
					Action: meta.Create,
				},
			},
		}
		return ps
	}
	// find many business set list for the user with any business set resources
	if ps.hitPattern(findmanyBusinessSetRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.BusinessSet,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	// find reduced business set list for the user with any business set resources
	if ps.hitPattern(findReducedBusinessSetListPattern, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.BusinessSet,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}
	// preview business set
	if ps.hitPattern(previewBusinessSet, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.BusinessSet,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}
	return ps
}
