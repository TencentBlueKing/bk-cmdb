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
		Name:           "ListFieldTemplateAttrCount",
		Description:    "查询字段模板属性数量",
		Pattern:        "/api/v3/findmany/field_template/attribute/count",
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
		Name:           "CompareFieldTemplateSyncStatus",
		Description:    "对比模型和模版的差异状态",
		Pattern:        "/api/v3/find/field_template/sync/status",
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
		Name:           "FindFieldTemplateTaskStatus",
		Description:    "查询字段组合模版任务状态",
		Pattern:        "/api/v3/find/field_template/tasks_status",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "FindFieldTemplateSimplifyByUnique",
		Description:    "根据模型唯一校验上的模版ID查找对应字段模版的简要信息",
		Pattern:        "/api/v3/find/field_template/simplify/by_unique_template_id",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "FindFieldTemplateSimplifyByObjectAttr",
		Description:    "根据模型属性上的模版ID查找对应字段模版的简要信息",
		Pattern:        "/api/v3/find/field_template/simplify/by_attr_template_id",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "ListFieldTemplateUnique",
		Description:    "查询字段模板唯一校验列表",
		Pattern:        "/api/v3/findmany/field_template/unique",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "CreateFieldTemplate",
		Description:    "创建字段模版",
		Pattern:        "/api/v3/create/field_template",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "UpdateFieldTemplate",
		Description:    "更新字段模版",
		Pattern:        "/api/v3/update/field_template",
		HTTPMethod:     http.MethodPut,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "DeleteFieldTemplate",
		Description:    "删除字段模版",
		Pattern:        "/api/v3/delete/field_template",
		HTTPMethod:     http.MethodDelete,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "CloneFieldTemplate",
		Description:    "克隆字段模版",
		Pattern:        "/api/v3/create/field_template/clone",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "UpdateFieldTemplateInfo",
		Description:    "更新字段模版基础信息",
		Pattern:        "/api/v3/update/field_template/info",
		HTTPMethod:     http.MethodPut,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "SyncFieldTemplateInfoToObjects",
		Description:    "字段组合模版的属性和唯一校验同步到模型",
		Pattern:        "/api/v3/update/topo/field_template/sync",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "FindFieldTemplateModelStatus",
		Description:    "查询字段组合模版模型同步状态",
		Pattern:        "/api/v3/find/field_template/model/status",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
}

func (ps *parseStream) fieldTemplate() *parseStream {
	return ParseStreamWithFramework(ps, FieldTemplateAuthConfigs)
}
