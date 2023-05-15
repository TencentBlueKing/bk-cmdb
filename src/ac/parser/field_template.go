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
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

var (
	findFieldTemplateRegexp = regexp.MustCompile(`^/api/v3/find/field_template/[0-9]+/?$`)
)

const (
	fieldTemplateBindObjPattern         = `/api/v3/field_template/bind/object`
	fieldTemplateUnBindObjPattern       = `/api/v3/field_template/unbind/object`
	findFieldTemplateTasksStatusPattern = `/api/v3/find/field_template/tasks_status`
)

func (ps *parseStream) fieldTemplate() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitRegexp(findFieldTemplateRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("find field template failed, got invalid url")
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				// todo: 待后续调整
				Basic: meta.Basic{
					// Type:   ,
					Action: meta.SkipAction,
				},
			},
		}
	}

	// 字段模版绑定模型的操作
	if ps.hitPattern(fieldTemplateBindObjPattern, http.MethodPost) {

		val, err := ps.RequestCtx.getValueFromBody(common.BKFieldID)
		if err != nil {
			ps.err = err
			return ps
		}
		// todo：待补充字段组合模版的编辑权限
		// 获取字段模版ID
		// id := val.Int()

		val, err = ps.RequestCtx.getValueFromBody("bk_obj_ids")
		if err != nil {
			ps.err = err
			return ps
		}

		models := val.Array()
		modelObjs := make([]string, 0)
		for _, modelID := range models {
			idStr := modelID.String()
			if idStr == "" {
				ps.err = errors.New("invalid process template id")
				return ps
			}
			modelObjs = append(modelObjs, idStr)
		}
		cond := mapstr.MapStr{
			common.BKObjIDField: mapstr.MapStr{
				common.BKDBIN: modelObjs,
			},
		}
		res, err := ps.searchModels(cond)
		if err != nil {
			ps.err = err
			return ps
		}
		if len(res) == 0 {
			ps.err = fmt.Errorf("model [%+v] not found", cond)
			return ps
		}

		for _, model := range res {
			ps.Attribute.Resources = append(ps.Attribute.Resources, meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.Update,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			})
		}

		return ps
	}

	// 字段模版解除绑定模型的操作
	if ps.hitPattern(fieldTemplateUnBindObjPattern, http.MethodPost) {
		val, err := ps.RequestCtx.getValueFromBody(common.BKObjIDField)
		if err != nil {
			ps.err = err
			return ps
		}
		obj := val.String()
		if obj == "" {
			ps.err = errors.New("obj must be set")
			return ps
		}

		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: obj})
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Model,
					Action:     meta.Update,
					InstanceID: model.ID,
				},
			},
		}
	}

	// 字段模版解除绑定模型的操作
	if ps.hitPattern(findFieldTemplateTasksStatusPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					//Type:       meta.Model,
					Action: meta.SkipAction,
				},
			},
		}
	}

	return ps
}
