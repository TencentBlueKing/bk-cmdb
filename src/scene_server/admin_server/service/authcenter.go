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
	"net/http"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *Service) InitAuthCenter(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(rHeader))

	bizs := []metadata.BizInst{}
	if err := s.db.Table(common.BKTableNameBaseApp).Find(condition.CreateCondition().Field("default").NotEq(1).ToMapStr()).All(s.ctx, &bizs); err != nil {
		blog.Errorf("init authcenter error: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommInitAuthcenterFailed, err.Error())})
		return
	}

	cls := []metadata.Classification{}
	if err := s.db.Table(common.BKTableNameObjClassifiction).Find(condition.CreateCondition().Field("ispre").NotEq(true).ToMapStr()).All(s.ctx, &cls); err != nil {
		blog.Errorf("init authcenter error: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommInitAuthcenterFailed, err.Error())})
		return
	}

	models := []metadata.Object{}
	if err := s.db.Table(common.BKTableNameObjDes).Find(condition.CreateCondition().Field("ispre").NotEq(true).ToMapStr()).All(s.ctx, &models); err != nil {
		blog.Errorf("init authcenter error: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommInitAuthcenterFailed, err.Error())})
		return
	}

	if err := s.authCenter.Init(s.ctx, meta.InitConfig{
		Bizs:            bizs,
		Models:          models,
		Classifications: cls,
	}); nil != err {
		blog.Errorf("init authcenter error: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommInitAuthcenterFailed, err.Error())})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp("init authcenter success"))
}
