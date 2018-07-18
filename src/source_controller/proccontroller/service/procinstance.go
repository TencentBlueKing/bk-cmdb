/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
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

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (ps *ProctrlServer) CreateProcInstanceModel(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)

	reqParam := make([]meta.ProcInstanceModel, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&reqParam); err != nil {
		blog.Errorf("create process instance model failed, decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.Infof("will create process instance model: %+v", reqParam)
	if err := ps.DbInstance.InsertMuti(common.BKTableNameProcInstanceModel, reqParam); err != nil {
		blog.Errorf("create process instance model failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateInstanceModel)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProctrlServer) GetProcInstanceModel(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)

	reqParam := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&reqParam); err != nil {
		blog.Errorf("get process instance model failed, decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.Infof("will get process instance model. condition: %v", reqParam)
	var result []interface{}
	if err := ps.DbInstance.GetMutilByCondition(common.BKTableNameProcInstanceModel, []string{}, reqParam, &result, "", 0, 0); err != nil {
		blog.Errorf("get process instance model failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetInstanceModel)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(result))
}

func (ps *ProctrlServer) DeleteProcInstanceModel(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.Core.CCErr.CreateDefaultCCErrorIf(language)

	reqParam := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&reqParam); err != nil {
		blog.Errorf("delete process instance model failed, decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	blog.Infof("will delete process instance model. condition: %+v", reqParam)
	if err := ps.DbInstance.DelByCondition(common.BKTableNameProcInstanceModel, reqParam); err != nil {
		blog.Errorf("delete process instance model failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteInstanceModel)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}
