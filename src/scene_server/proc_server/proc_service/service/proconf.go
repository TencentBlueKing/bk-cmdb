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
	"context"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"
)

func (ps *ProcServer) CreateConfigTemp(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	reqParam := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&reqParam); err != nil {
		blog.Errorf("create config template failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	reqParam[common.BKConfTempIdField] = "0" //use unified id generation method, uri will provide

	ret, err := ps.CoreAPI.ProcController().CreateConfTemp(context.Background(), req.Request.Header, reqParam)
	if err != nil || (err == nil && !ret.Result) {
		blog.Errorf("create config template failed by processcontroll. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateProcConf)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) DeleteConfigTemp(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	reqParam := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&reqParam); err != nil {
		blog.Errorf("delete config template failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	//logic here

	ret, err := ps.CoreAPI.ProcController().DeleteConfTemp(context.Background(), req.Request.Header, reqParam)
	if err != nil || (err == nil && !ret.Result) {
		blog.Errorf("delete config template failed by processcontroll. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteProcConf)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) UpdateConfigTemp(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	reqParam := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&reqParam); err != nil {
		blog.Errorf("update config template failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	//logic here

	ret, err := ps.CoreAPI.ProcController().UpdateConfTemp(context.Background(), req.Request.Header, reqParam)
	if err != nil || (err == nil && !ret.Result) {
		blog.Errorf("update config template failed by processcontroll. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcUpdateProcConf)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) QueryConfigTemp(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	reqParam := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&reqParam); err != nil {
		blog.Errorf("query config template failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	ret, err := ps.CoreAPI.ProcController().QueryConfTemp(context.Background(), req.Request.Header, reqParam)
	if err != nil || (err == nil && !ret.Result) {
		blog.Errorf("query config template failed by processcontroll. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(ret.Data))
}
