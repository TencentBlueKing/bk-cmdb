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
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	types "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"

	"github.com/emicklei/go-restful"
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin/json"
)

func (ps *ProcServer) PreviewCfg(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if nil != err {
		blog.Errorf("search template version failed! derr: %v,appIDStr:%s,rid:%s", err, appIDStr, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}
	var params meta.FilePriviewMap

	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("decode request body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	instArr := strings.Split(params.Inst, ".")
	if 0 != len(instArr) {
		blog.Errorf("inst params error: %v,input:%+v,rid:%s", err, params, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	setName := instArr[0]
	moduleName := instArr[1]
	funIDStr := instArr[2]
	instIDStr := instArr[3]

	funcID, err := strconv.ParseInt(funIDStr, 10, 64)
	if 0 != len(instArr) {
		blog.Errorf("funcID params error: %v,input:%+v,rid:%s", err, params, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}
	instID, err := strconv.ParseInt(instIDStr, 10, 64)
	if 0 != len(instArr) {
		blog.Errorf("inst params error: %v,input:%+v,rid:%s", err, params, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	variables := srvData.lgc.NewVariables(srvData.ctx, appID)
	vars := variables.GetStandVariables(srvData.ctx, setName, moduleName, funcID, instID)
	tpl, err := pongo2.FromString(params.Content)
	if err != nil {
		blog.Errorf("content params error: %v,input:%+v,rid:%s", err, params, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	out, err := tpl.Execute(pongo2.Context(vars))
	if err != nil {
		blog.Errorf("content params error: %v,input:%+v,rid:%s", err, params, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	result := types.MapStr{"content": out}
	resp.WriteEntity(meta.NewSuccessResp(result))
}

func (ps *ProcServer) CreateCfg(req *restful.Request, resp *restful.Response) {

}

func (ps *ProcServer) PushCfg(req *restful.Request, resp *restful.Response) {

}

func (ps *ProcServer) GetRemoteCfg(req *restful.Request, resp *restful.Response) {

}

// If it is necessary
func (ps *ProcServer) DiffCfg(req *restful.Request, resp *restful.Response) {

}
