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
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *Service) InitAuthCenter(req *restful.Request, resp *restful.Response) {
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
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.CCError(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if param.Host == "" {
		blog.Errorf("init iam host not set, rid:%s", rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.CCErrorf(common.CCErrCommParamsNeedSet, "host")})
		return
	}

	// 由于模型实例的编辑&删除拆分为实例级别, 需要先拿到当前已存在的模型, 再进行相应的IAM注册操作
	models, err := s.GetCustomObjects(rHeader)
	if err != nil {
		blog.Errorf("init iam failed, collect notPre-models failed, err: %s, rid:%s", err.Error(), rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.CCError(common.CCErrCommDBSelectFailed)})
		return
	}

	if err := s.iam.RegisterSystem(s.ctx, param.Host, models); err != nil {
		blog.Errorf("init iam failed, err: %+v, rid: %s", err, rid)
		result := &metadata.RespError{
			Msg: defErr.CCErrorf(common.CCErrCommInitAuthCenterFailed, err.Error()),
		}
		_ = resp.WriteError(http.StatusInternalServerError, result)
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}

// GetCustomObjects get objects which are custom.
func (s *Service) GetCustomObjects(header http.Header) ([]metadata.Object, error) {
	resp, err := s.CoreAPI.CoreService().Model().ReadModel(s.ctx, header, &metadata.QueryCondition{
		Fields: []string{common.BKObjIDField, common.BKObjNameField, common.BKFieldID},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKIsPre: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get custom models failed, err: %+v", err)
	}
	if len(resp.Data.Info) == 0 {
		blog.Info("get custom models, custom models not found")
	}

	objects := make([]metadata.Object, 0)
	for _, item := range resp.Data.Info {
		objects = append(objects, item.Spec)
	}

	return objects, nil
}
