/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
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
	"net/http"

	"configcenter/src/auth/meta"
)

func (ps *parseStream) adminRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.ConfigAdmin()

	return ps
}

var ConfigAdminConfigs = []AuthConfig{
	{
		Name:           "findConfigAdmin",
		Description:    "查询配置管理",
		Pattern:        "/api/v3/admin/find/system/config_admin",
		HTTPMethod:     http.MethodGet,
		ResourceType:   meta.ConfigAdmin,
		ResourceAction: meta.SkipAction,
	}, {
		Name:           "updateConfigAdmin",
		Description:    "更新配置管理",
		Pattern:        "/api/v3/admin/update/system/config_admin",
		HTTPMethod:     http.MethodPut,
		ResourceType:   meta.ConfigAdmin,
		ResourceAction: meta.Update,
	},
}

func (ps *parseStream) ConfigAdmin() *parseStream {
	return ParseStreamWithFramework(ps, ConfigAdminConfigs)
}
