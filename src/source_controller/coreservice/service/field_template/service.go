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

package fieldtmpl

import (
	"net/http"

	"configcenter/src/common/http/rest"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/source_controller/coreservice/service/capability"
)

type service struct {
	core core.Core
}

// InitFieldTemplate init field template service
func InitFieldTemplate(c *capability.Capability) {
	s := &service{
		core: c.Core,
	}

	// field template
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/field_template",
		Handler: s.ListFieldTemplate})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/field_template",
		Handler: s.CreateFieldTemplate})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/update/field_template/bind/object",
		Handler: s.FieldTemplateBindObject})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/update/field_template/unbind/object",
		Handler: s.FieldTemplateUnbindObject})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/field_template",
		Handler: s.DeleteFieldTemplate})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/field_template",
		Handler: s.UpdateFieldTemplate})

	// field template attribute
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/field_template/attribute",
		Handler: s.ListFieldTemplateAttr})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path: "/createmany/field_template/{bk_template_id}/attribute", Handler: s.CreateFieldTemplateAttrs})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodDelete,
		Path: "/delete/field_template/{bk_template_id}/attributes", Handler: s.DeleteFieldTemplateAttrs})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/field_template/{bk_template_id}/attributes",
		Handler: s.UpdateFieldTemplateAttrs})

	// field template unique
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/field_template/unique",
		Handler: s.ListFieldTemplateUnique})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/field_template/{bk_template_id}/unique",
		Handler: s.CreateFieldTemplateUniques})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/field_template/{bk_template_id}/uniques",
		Handler: s.DeleteFieldTemplateUniques})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/field_template/{bk_template_id}/uniques",
		Handler: s.UpdateFieldTemplateUniques})

	// field template relation
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/field_template/object/relation",
		Handler: s.ListObjFieldTmplRel})

	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/field_template/simplify/by_unique_template_id",
		Handler: s.FindFieldTmplSimplifyByUnique})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/field_template/simplify/by_attr_template_id",
		Handler: s.FindFieldTmplSimplifyByAttr})
}
