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

// Package tenantset is the tenant set package
package tenantset

import (
	"net/http"

	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/topo_server/service/capability"
)

type service struct {
	*capability.Capability
}

// InitTenantSet init tenant set service
func InitTenantSet(utility *rest.RestUtility, c *capability.Capability) {
	s := &service{
		Capability: c,
	}

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/tenant_set", Handler: s.ListTenantSet})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/tenant_set/{tenant_set_id}/tenant",
		Handler: s.ListTenantInTenantSet})
}
