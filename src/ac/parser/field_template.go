/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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
	"net/http"
	"regexp"

	"configcenter/src/ac/meta"
)

// FieldTemplateAuthConfigs field template related auth configs, skip all, authorize in topo-server.
var FieldTemplateAuthConfigs = []AuthConfig{
	{
		Name:           "ListFieldTemplate",
		Description:    "查询字段模板列表",
		Pattern:        "/api/v3/findmany/field_template",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "ListFieldTemplateAttr",
		Description:    "查询字段模板属性列表",
		Pattern:        "/api/v3/findmany/field_template/attribute",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "ListObjFieldTmplRel",
		Description:    "查询字段模板和模型关系列表",
		Pattern:        "/api/v3/findmany/field_template/object/relation",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "ListFieldTmplByObj",
		Description:    "根据模型查询字段模板列表",
		Pattern:        "/api/v3/findmany/field_template/by_object",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "ListObjByFieldTmpl",
		Description:    "根据字段模板查询模型列表",
		Pattern:        "/api/v3/findmany/object/by_field_template",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "CompareFieldTemplateAttr",
		Description:    "对比字段模板和模型中的字段",
		Pattern:        "/api/v3/find/field_template/attribute/difference",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "CompareFieldTemplateUnique",
		Description:    "对比字段模板和模型中的唯一校验",
		Pattern:        "/api/v3/find/field_template/unique/difference",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "QueryFieldTemplateBriefInformation",
		Description:    "查询字段组合模版简要信息",
		Regex:          regexp.MustCompile(`^/api/v3/find/field_template/[0-9]+/?$`),
		HTTPMethod:     http.MethodGet,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "FieldTemplateBindingModel",
		Description:    "字段组合模版绑定模型",
		Pattern:        "/api/v3/update/field_template/bind/object",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "FieldTemplateUnbindingModel",
		Description:    "字段组合模版解除绑定模型",
		Pattern:        "/api/v3/update/field_template/unbind/object",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "FieldTemplateUnbindingModel",
		Description:    "查询字段组合模版任务状态",
		Pattern:        "/api/v3/task/find/field_template/tasks_status",
		HTTPMethod:     http.MethodGet,
		ResourceAction: meta.SkipAction,
	},
}

func (ps *parseStream) fieldTemplate() *parseStream {
	return ParseStreamWithFramework(ps, FieldTemplateAuthConfigs)
}
