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

	"configcenter/src/ac/meta"
	"configcenter/src/common/watch"
)

func (ps *parseStream) eventRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.watch()

	return ps
}

var (
	watchResourceRegexp = regexp.MustCompile(`^/api/v3/event/watch/resource/\S+/?$`)
)

func (ps *parseStream) watch() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// watch resource.
	if ps.hitRegexp(watchResourceRegexp, http.MethodPost) {
		resource := ps.RequestCtx.Elements[5]
		if len(resource) == 0 {
			ps.err = fmt.Errorf("watch event resource, but got empty resource: %s", ps.RequestCtx.Elements[5])
			return ps
		}

		if resource == string(watch.HostIdentifier) {
			// redirect host identity resource to host resource in iam.
			resource = string(watch.Host)
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.EventWatch,
					Action: meta.Action(resource),
				},
			},
		}
		return ps
	}

	return ps
}
