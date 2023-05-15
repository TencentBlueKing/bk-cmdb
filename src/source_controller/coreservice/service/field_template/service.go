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
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/field_template/bind/object",
		Handler: s.FieldTemplateBindObject})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/field_template/unbind/object",
		Handler: s.FieldTemplateUnBindObject})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "find/field_template/tasks_status",
		Handler: s.FindFieldTemplateTasksStatus})
}
