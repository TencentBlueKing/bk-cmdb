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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"

	"github.com/emicklei/go-restful/v3"
)

// RefreshTenant refresh tenant info
func (s *service) RefreshTenant(req *restful.Request, resp *restful.Response) {

	rid := httpheader.GetRid(req.Request.Header)
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(req.Request.Header))
	allTenants, err := s.clientSet.CoreService().Tenant().RefreshTenants(req.Request.Context(), req.Request.Header)
	if err != nil {
		blog.Errorf("refresh tenant info failed, rid: %s", rid)
		resp.WriteError(http.StatusInternalServerError,
			&metadata.RespError{Msg: defErr.Errorf(common.CCErrCommAddTenantErr, "refresh tenant info failed")})
		return
	}

	tenant.SetTenant(allTenants)

	resp.WriteEntity(metadata.NewSuccessResp(allTenants))
}
