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
	"configcenter/src/common"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"regexp"
	"strconv"
)

var SetTemplateAuthConfigs = []AuthConfig{
	{
		Name:                  "createSetTemplatePattern",
		Description:           "创建集群模板",
		Regex:                 regexp.MustCompile(`^/api/v3/create/topo/set_template/bk_biz_id/([0-9]+)$`),
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.SetTemplate,
		ResourceAction:        meta.Create,
	}, {
		Name:                  "updateSetTemplate",
		Description:           "更新集群模板",
		Regex:                 regexp.MustCompile(`^/api/v3/update/topo/set_template/([0-9]+)/bk_biz_id/([0-9]+)$`),
		HTTPMethod:            http.MethodPut,
		RequiredBizInMetadata: true,
		ResourceType:          meta.SetTemplate,
		ResourceAction:        meta.Update,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			templateID := gjson.GetBytes(request.Body, common.BKFieldID).Int()
			if templateID <= 0 {
				return nil, errors.New("invalid set template")
			}
			return []int64{templateID}, nil
		},
	}, {
		Name:                  "getSetTemplate",
		Description:           "获取集群模板",
		Regex:                 regexp.MustCompile(`^/api/v3/find/topo/set_template/([0-9]+)/bk_biz_id/([0-9]+)$`),
		HTTPMethod:            http.MethodGet,
		RequiredBizInMetadata: true,
		ResourceType:          meta.SetTemplate,
		ResourceAction:        meta.FindMany,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			ss := re.FindStringSubmatch(request.URI)
			if len(ss) < 1 {
				return nil, errors.New("getSetTemplate regex match nothing")
			}
			id, err := strconv.ParseInt(ss[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("getSetTemplate regex parse match to int failed, err: %+v", err)
			}
			return []int64{id}, nil
		},
	}, {
		Name:                  "listSetTemplate",
		Description:           "列表查询集群模板",
		Regex:                 regexp.MustCompile(`^/api/v3/findmany/topo/set_template/bk_biz_id/([0-9]+)$`),
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.SetTemplate,
		ResourceAction:        meta.FindMany,
	}, {
		Name:                  "listSetTemplateWeb",
		Description:           "列表查询集群模板",
		Regex:                 regexp.MustCompile(`^/api/v3/findmany/topo/set_template/bk_biz_id/([0-9]+)/web$`),
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.SetTemplate,
	}, {
		Name:                  "deleteSetTemplatePattern",
		Description:           "删除集群模板",
		Pattern:               "/api/v3/deletemany/topo/set_template/bk_biz_id/([0-9]+)",
		HTTPMethod:            http.MethodDelete,
		RequiredBizInMetadata: true,
		ResourceType:          meta.SetTemplate,
		ResourceAction:        meta.Delete,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			// TODO
			return []int64{}, nil
		},
	},
}

func (ps *parseStream) SetTemplate() *parseStream {
	resources, err := MatchAndGenerateBizInURLIAMResource(ServiceTemplateAuthConfigs, ps.RequestCtx)
	if err != nil {
		ps.err = err
	}
	if resources != nil {
		ps.Attribute.Resources = resources
	}
	return ps
}
