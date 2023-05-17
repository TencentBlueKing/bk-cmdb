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

//
//func (ps *parseStream) fieldTemplate1() *parseStream {
//	if ps.shouldReturn() {
//		return ps
//	}
//
//	if ps.hitRegexp(findFieldTemplateRegexp, http.MethodGet) {
//		if len(ps.RequestCtx.Elements) != 5 {
//			ps.err = errors.New("find field template failed, got invalid url")
//			return ps
//		}
//		ps.Attribute.Resources = []meta.ResourceAttribute{
//			{
//				// todo: 待后续调整
//				Basic: meta.Basic{
//					// Type:   ,
//					Action: meta.SkipAction,
//				},
//			},
//		}
//	}
//
//	// 字段模版绑定模型的操作
//	if ps.hitPattern(fieldTemplateBindObjPattern, http.MethodPost) {
//
//		val, err := ps.RequestCtx.getValueFromBody(common.BKFieldID)
//		if err != nil {
//			ps.err = err
//			return ps
//		}
//		// todo：待补充字段组合模版的编辑权限
//		// 获取字段模版ID
//		// id := val.Int()
//
//		val, err = ps.RequestCtx.getValueFromBody("bk_obj_ids")
//		if err != nil {
//			ps.err = err
//			return ps
//		}
//
//		models := val.Array()
//		modelObjs := make([]string, 0)
//		for _, modelID := range models {
//			idStr := modelID.String()
//			if idStr == "" {
//				ps.err = errors.New("invalid process template id")
//				return ps
//			}
//			modelObjs = append(modelObjs, idStr)
//		}
//		cond := mapstr.MapStr{
//			common.BKObjIDField: mapstr.MapStr{
//				common.BKDBIN: modelObjs,
//			},
//		}
//		res, err := ps.searchModels(cond)
//		if err != nil {
//			ps.err = err
//			return ps
//		}
//		if len(res) == 0 {
//			ps.err = fmt.Errorf("model [%+v] not found", cond)
//			return ps
//		}
//
//		for _, model := range res {
//			ps.Attribute.Resources = append(ps.Attribute.Resources, meta.ResourceAttribute{
//				Basic: meta.Basic{
//					Type:   meta.Model,
//					Action: meta.Update,
//				},
//				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
//			})
//		}
//
//		return ps
//	}
//
//	// 字段模版解除绑定模型的操作
//	if ps.hitPattern(fieldTemplateUnBindObjPattern, http.MethodPost) {
//		val, err := ps.RequestCtx.getValueFromBody(common.BKObjIDField)
//		if err != nil {
//			ps.err = err
//			return ps
//		}
//		obj := val.String()
//		if obj == "" {
//			ps.err = errors.New("obj must be set")
//			return ps
//		}
//
//		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: obj})
//		if err != nil {
//			ps.err = err
//			return ps
//		}
//		ps.Attribute.Resources = []meta.ResourceAttribute{
//			{
//				Basic: meta.Basic{
//					Type:       meta.Model,
//					Action:     meta.Update,
//					InstanceID: model.ID,
//				},
//			},
//		}
//	}
//
//	// 字段模版解除绑定模型的操作
//	if ps.hitPattern(findFieldTemplateTasksStatusPattern, http.MethodPost) {
//		ps.Attribute.Resources = []meta.ResourceAttribute{
//			{
//				Basic: meta.Basic{
//					//Type:       meta.Model,
//					Action: meta.SkipAction,
//				},
//			},
//		}
//	}
//
//	return ps
//}

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
		Description:    "查询字段模版简要信息",
		Regex:          regexp.MustCompile(`^/api/v3/find/field_template/[0-9]+/?$`),
		HTTPMethod:     http.MethodGet,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "FieldTemplateBindingModel",
		Description:    "字段模版绑定模型",
		Pattern:        "/api/v3/field_template/bind/object",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "FieldTemplateUnBindingModel",
		Description:    "字段模版解除绑定模型",
		Pattern:        "/api/v3/field_template/unbind/object",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "FieldTemplateUnBindingModel",
		Description:    "查询字段模版任务状态",
		Pattern:        "/api/v3/task/find/field_template/tasks_status",
		HTTPMethod:     http.MethodGet,
		ResourceAction: meta.SkipAction,
	},
}

func (ps *parseStream) fieldTemplate() *parseStream {
	return ParseStreamWithFramework(ps, FieldTemplateAuthConfigs)
}
