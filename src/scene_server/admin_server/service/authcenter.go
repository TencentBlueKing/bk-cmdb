/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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
	"encoding/json"
	"fmt"
	"net/http"

	"configcenter/pkg/tenant/logics"
	aciam "configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/iam"

	"github.com/emicklei/go-restful/v3"
)

// InitAuthCenter init auth resources on IAM
func (s *Service) InitAuthCenter(req *restful.Request, resp *restful.Response) {
	if !auth.EnableAuthorize() {
		_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

	rHeader := req.Request.Header
	rid := httpheader.GetRid(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))
	kit := rest.NewKitFromHeader(rHeader, s.CCErr)

	if !logics.ValidatePlatformTenantMode(kit.TenantID, s.Config.EnableMultiTenantMode) {
		blog.Errorf("non-system tenant register auth account, tenantID: %s, rid: %s", kit.TenantID, kit.Rid)
		result := &metadata.RespError{
			Msg:     fmt.Errorf("non-system tenant register auth account, tenantID: %s", kit.TenantID),
			ErrCode: common.CCErrAPICheckTenantInvalid,
		}
		_ = resp.WriteError(http.StatusOK, result)
		return
	}

	param := struct {
		Host string `json:"host"`
	}{}
	if err := json.NewDecoder(req.Request.Body).Decode(&param); err != nil {
		blog.Errorf("init iam failed with decode body err: %s, rid:%s", err.Error(), rid)
		_ = resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.CCError(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if param.Host == "" {
		blog.Errorf("init iam host not set, rid:%s", rid)
		_ = resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommParamsNeedSet, "host")})
		return
	}

	// 由于模型实例的编辑&删除拆分为实例级别, 需要先拿到当前已存在的模型, 再进行相应的IAM注册操作
	models, err := iam.GetCustomObjects(kit, s.db)
	if err != nil {
		blog.Errorf("init iam failed, collect notPre-models failed, err: %s, rid:%s", err.Error(), rid)
		_ = resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.CCError(common.CCErrCommDBSelectFailed)})
		return
	}

	opt := &aciam.RegisterIamOptions{
		Host:    param.Host,
		Objects: models,
	}
	if err := s.iam.Register(kit.Ctx, kit.Header, s.cache, opt, kit.Rid); err != nil {
		blog.Errorf("init iam failed, err: %+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.CCErrorf(common.CCErrCommInitAuthCenterFailed, err.Error()),
		}
		_ = resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}
