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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"configcenter/src/auth/meta"
)

var SetTemplateAuthConfigs = []AuthConfig{
	{
		Name:           "CreateSetTemplateRegex",
		Description:    "创建集群模板",
		Regex:          regexp.MustCompile(`^/api/v3/create/topo/set_template/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.SetTemplate,
		ResourceAction: meta.Create,
	}, {
		Name:           "UpdateSetTemplateRegex",
		Description:    "更新集群模板",
		Regex:          regexp.MustCompile(`^/api/v3/update/topo/set_template/([0-9]+)/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPut,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.SetTemplate,
		ResourceAction: meta.Update,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			ss := re.FindStringSubmatch(request.URI)
			if len(ss) < 2 {
				return nil, errors.New("UpdateSetTemplateRegex regex doesn't match")
			}
			id, err := strconv.ParseInt(ss[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("UpdateSetTemplateRegex regex parse match to int failed, err: %+v", err)
			}
			return []int64{id}, nil
		},
	}, {
		Name:           "DeleteSetTemplateRegex",
		Description:    "删除集群模板",
		Regex:          regexp.MustCompile(`^/api/v3/deletemany/topo/set_template/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.SetTemplate,
		ResourceAction: meta.Delete,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			data := &struct {
				SetTemplateIDs []int64 `json:"set_template_ids" mapstructure:"set_template_ids"`
			}{}
			if err := json.Unmarshal(request.Body, data); err != nil {
				return nil, fmt.Errorf("unmarshal failed, err: %+v", err)
			}
			return data.SetTemplateIDs, nil
		},
	}, {
		Name:           "GetSetTemplateRegex",
		Description:    "获取集群模板",
		Regex:          regexp.MustCompile(`^/api/v3/find/topo/set_template/([0-9]+)/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodGet,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.SetTemplate,
		ResourceAction: meta.Find,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			ss := re.FindStringSubmatch(request.URI)
			if len(ss) < 2 {
				return nil, errors.New("getSetTemplate regex match nothing")
			}
			id, err := strconv.ParseInt(ss[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("getSetTemplate regex parse match to int failed, err: %+v", err)
			}
			return []int64{id}, nil
		},
	}, {
		Name:           "ListSetTemplateRegex",
		Description:    "列表查询集群模板",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/topo/set_template/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.SetTemplate,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "ListSetTemplateWebRegex",
		Description:    "列表查询集群模板-Web",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/topo/set_template/bk_biz_id/([0-9]+)/web/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.SetTemplate,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "ListSetTplRelatedSvcTplRegex",
		Description:    "查询集群模板关联的服务模板列表",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/topo/set_template/([0-9]+)/bk_biz_id/([0-9]+)/service_templates/?$`),
		HTTPMethod:     http.MethodGet,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "ListSetTplRelatedSvcTplRegex",
		Description:    "查询集群模板关联的服务模板列表-Web",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/topo/set_template/([0-9]+)/bk_biz_id/([0-9]+)/service_templates/with_statistics/?$`),
		HTTPMethod:     http.MethodGet,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "ListSetTplRelatedSetsWebRegex",
		Description:    "查询集群模板关联的集群列表-Web",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/topo/set_template/([0-9]+)/bk_biz_id/([0-9]+)/sets/web/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.ProcessServiceTemplate,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "DiffSetTplWithInstRegex",
		Description:    "对比集群模板与集群之间的差异",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/topo/set_template/([0-9]+)/bk_biz_id/([0-9]+)/diff_with_instances/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.ModelSet,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "SyncSetTplToInstRegex",
		Description:    "用集群模板同步集群",
		Regex:          regexp.MustCompile(`^/api/v3/updatemany/topo/set_template/([0-9]+)/bk_biz_id/([0-9]+)/sync_to_instances/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.ModelSet,
		ResourceAction: meta.UpdateMany,
		InstanceIDGetter: func(request *RequestContext, re *regexp.Regexp) (int64s []int64, e error) {
			data := &struct {
				SetIDs []int64 `json:"bk_set_ids" mapstructure:"bk_set_ids"`
			}{}
			if err := json.Unmarshal(request.Body, data); err != nil {
				return nil, fmt.Errorf("unmarshal failed, err: %+v", err)
			}
			return data.SetIDs, nil
		},
	}, {
		Name:           "GetSetSyncStatusRegex",
		Description:    "获取集群同步状态",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/topo/set_template/([0-9]+)/bk_biz_id/([0-9]+)/instances_sync_status/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.ModelSet,
		ResourceAction: meta.FindMany,
	}, {
		Name:           "ListSetTemplateSyncStatusRegex",
		Description:    "获取集群同步状态",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/topo/set_template_sync_status/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    BizIDFromURLGetter,
		ResourceType:   meta.ModelSet,
		ResourceAction: meta.FindMany,
	}, {
		Name:             "ListSetTemplateSyncHistoryRegex",
		Description:      "集群模板的同步历史记录",
		Regex:            regexp.MustCompile(`^/api/v3/findmany/topo/set_template_sync_history/bk_biz_id/([0-9]+)/?$`),
		HTTPMethod:       http.MethodPost,
		BizIDGetter:      BizIDFromURLGetter,
		ResourceType:     meta.ModelSet,
		ResourceAction:   meta.FindMany,
		InstanceIDGetter: nil,
	},
}

func (ps *parseStream) SetTemplate() *parseStream {
	return ParseStreamWithFramework(ps, SetTemplateAuthConfigs)
}
