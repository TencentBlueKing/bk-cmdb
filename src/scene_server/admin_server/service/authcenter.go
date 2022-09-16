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

package service

import (
	"encoding/json"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
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
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))
	if !auth.EnableAuthorize() {
		blog.Warnf("received iam initialization request, but auth not enabled, rid: %s", rid)
		_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
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
	models, err := iam.GetCustomObjects(s.ctx, s.db)
	if err != nil {
		blog.Errorf("init iam failed, collect notPre-models failed, err: %s, rid:%s", err.Error(), rid)
		_ = resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.CCError(common.CCErrCommDBSelectFailed)})
		return
	}

	if err := s.iam.Register(s.ctx, param.Host, models, rid); err != nil {
		blog.Errorf("init iam failed, err: %+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.CCErrorf(common.CCErrCommInitAuthCenterFailed, err.Error()),
		}
		_ = resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}

// RegisterAuthAccount TODO
/**
	此接口用于在全新环境部署时，需要进行调用，目的是注册CMDB平台信息到IAM接口
	需要传入host参数，是authserver的ip:port信息，如下通过curl调用：
	curl -X POST -H 'Content-Type:application/json' -H 'BK_USER:migrate' -H 'HTTP_BLUEKING_SUPPLIER_ID:0'
	--data '{"host": "http://cmdb-auth.service.consul:'$BK_CMDB_AUTH_PORT'"}'
    http://cmdb-admin.service.consul:$BK_CMDB_ADMIN_PORT/migrate/v3/authcenter/register
*/
// RegisterAuthAccount register auth account to iam
func (s *Service) RegisterAuthAccount(req *restful.Request, resp *restful.Response) {
	if !auth.EnableAuthorize() {
		_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

	rHeader := req.Request.Header
	rid := util.GetHTTPCCRequestID(rHeader)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))
	if !auth.EnableAuthorize() {
		blog.Warnf("received iam register request, but auth not enabled, rid: %s", rid)
		_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

	param := struct {
		Host string `json:"host"`
	}{}
	if err := json.NewDecoder(req.Request.Body).Decode(&param); err != nil {
		blog.Errorf("register iam failed with decode body err: %v, rid:%s", err, rid)
		_ = resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.CCError(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if param.Host == "" {
		blog.Errorf("register iam host not set, rid:%s", rid)
		_ = resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommParamsNeedSet, "host")})
		return
	}

	if err := s.iam.RegisterToIAM(s.ctx, param.Host); err != nil {
		blog.Errorf("register cmdb to iam failed, err: %v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: err,
		}
		_ = resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}
