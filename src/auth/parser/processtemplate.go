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
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"

	"github.com/tidwall/gjson"
)

var ProcessTemplateAuthConfigs = []AuthConfig{
	{
		Name:                  "createProcessTemplateBatchPattern",
		Description:           "创建进程模板",
		Pattern:               "/api/v3/createmany/proc/proc_template",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessTemplate,
		ResourceAction:        meta.Create,
	}, {
		Name:                  "updateProcessTemplatePattern",
		Description:           "更新进程模板",
		Pattern:               "/api/v3/update/proc/proc_template",
		HTTPMethod:            http.MethodPut,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessTemplate,
		ResourceAction:        meta.Update,
		InstanceIDGetter: func(request *RequestContext, config AuthConfig) (int64s []int64, e error) {
			procTemplateID := gjson.GetBytes(request.Body, common.BKProcessTemplateIDField).Int()
			if procTemplateID <= 0 {
				return nil, errors.New("invalid process template id")
			}
			return []int64{procTemplateID}, nil
		},
	}, {
		Name:                  "deleteProcessTemplateBatchPattern",
		Description:           "删除进程模板",
		Pattern:               "/api/v3/deletemany/proc/proc_template",
		HTTPMethod:            http.MethodDelete,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessTemplate,
		ResourceAction:        meta.Delete,
		InstanceIDGetter: func(request *RequestContext, config AuthConfig) (int64s []int64, e error) {
			procTemplateIDs := gjson.GetBytes(request.Body, "process_templates").Array()
			ids := make([]int64, 0)
			for _, procTemplateID := range procTemplateIDs {
				id := procTemplateID.Int()
				if id <= 0 {
					return nil, errors.New("invalid process template id")
				}
				ids = append(ids, id)
			}
			return ids, nil
		},
	}, {
		Name:                  "findProcessTemplateBatchPattern",
		Description:           "查找进程模板",
		Pattern:               "/api/v3/findmany/proc/proc_template",
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessTemplate,
		// authorization should implements in scene server
		ResourceAction: meta.SkipAction,
	}, {
		Name:                  "findProcessTemplateRegexp",
		Description:           "获取进程模板",
		Regex:                 regexp.MustCompile(`/api/v3/find/proc/proc_template/id/([0-9]+)/?$`),
		HTTPMethod:            http.MethodPost,
		RequiredBizInMetadata: true,
		ResourceType:          meta.ProcessTemplate,
		ResourceAction:        meta.Find,
		InstanceIDGetter: func(request *RequestContext, config AuthConfig) (int64s []int64, e error) {
			subMatch := config.Regex.FindStringSubmatch(request.URI)
			for _, subStr := range subMatch {
				id, err := strconv.ParseInt(subStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("parse template id to int64 failed, err: %s", err)
				}
				return []int64{id}, nil
			}
			blog.Errorf("unexpected error: this code shouldn't be reached, rid: %s", request.Rid)
			return nil, errors.New("unexpected error: this code shouldn't be reached")
		},
	},
}

func (ps *parseStream) ProcessTemplate() *parseStream {
	resources, err := MatchAndGenerateIAMResource(ProcessTemplateAuthConfigs, ps.RequestCtx)
	if err != nil {
		ps.err = err
	}
	if resources != nil {
		ps.Attribute.Resources = resources
	}
	return ps
}
