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

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/topo_server/logics"
	"configcenter/src/scene_server/topo_server/service/capability"
)

type service struct {
	clientSet apimachinery.ClientSetInterface
	logics    logics.Logics
	auth      *extensions.AuthManager
}

// InitFieldTemplate init field template service
func InitFieldTemplate(c *capability.Capability) {
	s := &service{
		clientSet: c.ClientSet,
		logics:    c.Logics,
		auth:      c.AuthManager,
	}

	// field template
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/field_template",
		Handler: s.ListFieldTemplate})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/field_template",
		Handler: s.CreateFieldTemplate})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/field_template/{id}",
		Handler: s.FindFieldTemplateByID})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/update/field_template/bind/object",
		Handler: s.FieldTemplateBindObject})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/update/field_template/unbind/object",
		Handler: s.FieldTemplateUnbindObject})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/field_template",
		Handler: s.DeleteFieldTemplate})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/field_template/clone",
		Handler: s.CloneFieldTemplate})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/field_template",
		Handler: s.UpdateFieldTemplate})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/field_template/info",
		Handler: s.UpdateFieldTemplateInfo})

	// field template attribute
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/field_template/attribute",
		Handler: s.ListFieldTemplateAttr})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/field_template/attribute/count",
		Handler: s.CountFieldTemplateAttr})

	// field template unique
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/field_template/unique",
		Handler: s.ListFieldTemplateUnique})

	// field template sync to object
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/update/topo/field_template/sync",
		Handler: s.SyncFieldTemplateInfoToObjects})
	// field template relation
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/field_template/object/relation",
		Handler: s.ListObjFieldTmplRel})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/field_template/by_object",
		Handler: s.ListFieldTmplByObj})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/object/by_field_template",
		Handler: s.ListObjByFieldTmpl})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/field_template/object/count",
		Handler: s.CountFieldTemplateObj})

	// field template task
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/sync/field_template/object/task",
		Handler: s.SyncFieldTemplateToObjectTask})

	// compare field template with object
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/field_template/attribute/difference",
		Handler: s.CompareFieldTemplateAttr})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/field_template/unique/difference",
		Handler: s.CompareFieldTemplateUnique})

	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/field_template/tasks_status",
		Handler: s.ListFieldTemplateTasksStatus})

	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/field_template/sync/status",
		Handler: s.ListFieldTemplateSyncStatus})

	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/field_template/simplify/by_unique_template_id",
		Handler: s.ListFieldTmplByUniqueTmplIDForUI})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/field_template/simplify/by_attr_template_id",
		Handler: s.ListFieldTmplByObjectTmplIDForUI})

	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/field_template/model/status",
		Handler: s.ListFieldTemplateModelStatus})
}
