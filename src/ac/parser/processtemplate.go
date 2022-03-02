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
	"strings"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
)

// NOTE: 进程模板增删改操作检查服务模板的编辑权限
var ProcessTemplateAuthConfigs = []AuthConfig{
	{
		Name:         "createProcessTemplateBatchPattern",
		Description:  "创建进程模板",
		Pattern:      "/api/v3/createmany/proc/proc_template",
		HTTPMethod:   http.MethodPost,
		BizIDGetter:  DefaultBizIDGetter,
		ResourceType: meta.ProcessTemplate,
		// ResourceAction:        meta.Create,
		ResourceAction: meta.SkipAction,
	}, {
		Name:         "updateProcessTemplatePattern",
		Description:  "更新进程模板",
		Pattern:      "/api/v3/update/proc/proc_template",
		HTTPMethod:   http.MethodPut,
		BizIDGetter:  DefaultBizIDGetter,
		ResourceType: meta.ProcessTemplate,
		// ResourceAction:        meta.Update,
		ResourceAction: meta.SkipAction,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			val, err := request.getValueFromBody(common.BKProcessTemplateIDField)
			if err != nil {
				return nil, err
			}
			procTemplateID := val.Int()
			if procTemplateID <= 0 {
				return nil, errors.New("invalid process template id")
			}
			return []int64{procTemplateID}, nil
		},
	}, {
		Name:         "deleteProcessTemplateBatchPattern",
		Description:  "删除进程模板",
		Pattern:      "/api/v3/deletemany/proc/proc_template",
		HTTPMethod:   http.MethodDelete,
		BizIDGetter:  DefaultBizIDGetter,
		ResourceType: meta.ProcessServiceTemplate,
		// ResourceAction:        meta.Delete,
		ResourceAction: meta.SkipAction,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			val, err := request.getValueFromBody("process_templates")
			if err != nil {
				return nil, err
			}
			procTemplateIDs := val.Array()
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
		Name:         "findProcessTemplateBatchPattern",
		Description:  "查找进程模板",
		Pattern:      "/api/v3/findmany/proc/proc_template",
		HTTPMethod:   http.MethodPost,
		BizIDGetter:  DefaultBizIDGetter,
		ResourceType: meta.ProcessTemplate,
		// authorization should implements in scene server
		ResourceAction: meta.SkipAction,
	}, {
		// search process template by biz set regex, authorize by biz set access permission, **only for ui**
		Name:           "findProcessTemplateBatchByBizSetRegexp",
		Description:    "查找业务集下的进程模板",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/proc/biz_set/[0-9]+/proc_template/?$`),
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta.BizSet,
		ResourceAction: meta.AccessBizSet,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			if len(request.Elements) != 7 {
				return nil, fmt.Errorf("get invalid url elements length %d", len(request.Elements))
			}

			bizSetID, err := strconv.ParseInt(request.Elements[5], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("get invalid business set id %s, err: %v", request.Elements[5], err)
			}
			return []int64{bizSetID}, nil
		},
	}, {
		Name:         "findProcessTemplateRegexp",
		Description:  "获取进程模板",
		Regex:        regexp.MustCompile(`/api/v3/find/proc/proc_template/id/([0-9]+)/?$`),
		HTTPMethod:   http.MethodPost,
		BizIDGetter:  DefaultBizIDGetter,
		ResourceType: meta.ProcessTemplate,
		// ResourceAction:        meta.Find,
		ResourceAction: meta.SkipAction,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			subMatch := re.FindStringSubmatch(request.URI)
			for _, subStr := range subMatch {
				if strings.Contains(subStr, "api") {
					continue
				}
				id, err := strconv.ParseInt(subStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("parse template id to int64 failed, err: %s", err)
				}
				return []int64{id}, nil
			}
			blog.Errorf("unexpected error: this code shouldn't be reached, rid: %s", request.Rid)
			return nil, errors.New("unexpected error: this code shouldn't be reached")
		},
	}, {
		// find process template by id and biz set regex, authorize by biz set access permission, **only for ui**
		Name:           "findProcessTemplateByBizSetRegexp",
		Description:    "查找业务集下的进程模板",
		Regex:          regexp.MustCompile(`^/api/v3/find/proc/biz_set/[0-9]+/proc_template/id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPost,
		ResourceType:   meta.BizSet,
		ResourceAction: meta.AccessBizSet,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			if len(request.Elements) != 9 {
				return nil, fmt.Errorf("get invalid url elements length %d", len(request.Elements))
			}

			procTempID, err := strconv.ParseInt(request.Elements[8], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("get invalid process template id %s, err: %v", request.Elements[6], err)
			}

			if procTempID <= 0 {
				return nil, fmt.Errorf("get invalid process template id %s, err: %v", request.Elements[6], err)
			}

			bizSetID, err := strconv.ParseInt(request.Elements[5], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("get invalid business set id %s, err: %v", request.Elements[5], err)
			}
			return []int64{bizSetID}, nil
		},
	},
}

func (ps *parseStream) ProcessTemplate() *parseStream {
	return ParseStreamWithFramework(ps, ProcessTemplateAuthConfigs)
}
