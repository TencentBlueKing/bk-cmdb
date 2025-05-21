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

package service

import (
	"net/http"

	"configcenter/pkg/tenant"
	"configcenter/pkg/tenant/logics"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"

	"github.com/emicklei/go-restful/v3"
)

// RefreshTenant refresh tenant info
func (s *cacheService) RefreshTenant(req *restful.Request, resp *restful.Response) {
	kit := rest.NewKitFromHeader(req.Request.Header, s.err)
	tenants, err := logics.GetAllTenantsFromDB(kit.Ctx, mongodb.Shard(kit.SysShardOpts()))
	if err != nil {
		blog.Errorf("refresh tenant info but get all tenants failed, err: %v, rid: %s", err)
		resp.WriteError(http.StatusOK, &metadata.RespError{Msg: err})
		return
	}

	tenant.SetTenant(tenants)
	resp.WriteEntity(metadata.NewSuccessResp(tenants))
}
