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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	types "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/scene_server/validator"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"
)

func (ps *ProcServer) SearchTemplateVersion(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	ownerID := req.PathParameter(common.BKOwnerIDField)
	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if nil != err {
		blog.Errorf("search template version failed! derr: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	templateIDStr := req.PathParameter(common.BKTemlateIDField)
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if nil != err {
		blog.Errorf("search template version failed! derr: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	params := types.MapStr{}
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("decode request body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input := new(meta.QueryCondition)

	conditon := types.MapStr{common.BKOwnerIDField: ownerID, common.BKAppIDField: appID, common.BKTemlateIDField: templateID}
	status, ok := params[common.BKStatusField]
	if ok {
		conditon[common.BKStatusField] = status
	}
	input.Condition = conditon
	input.Limit.Limit = common.BKNoLimit

	ret, err := ps.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDTempVersion, input)
	if err != nil {
		blog.Errorf("SearchTemplateVersion SearchObjects http do error,err:%s,input:%+v,query:%+v,rid:%s", err.Error(), params, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
	}
	if !ret.Result {
		blog.Errorf("SearchTemplateVersion  SearchObjects http response error,err code:%d,err msg:%s,input:%+v,query:%+v,rid:%s", ret.Code, ret.ErrMsg, params, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(ret.Code, ret.ErrMsg)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(ret.Data.Info))
}

func (ps *ProcServer) CreateTemplateVersion(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr
	user := srvData.user

	ownerID := req.PathParameter(common.BKOwnerIDField)
	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if nil != err {
		blog.Errorf("create config template failed! derr: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	templateIDStr := req.PathParameter(common.BKTemlateIDField)
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if nil != err {
		blog.Errorf("create config template failed! derr: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	var params meta.TemplateVersion

	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("create config version failed! decode request body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input := types.MapStr{common.BKAppIDField: appID,
		common.BKOperatorField:    user,
		common.BKTemlateIDField:   templateID,
		common.BKContentField:     params.Content,
		common.BKStatusField:      params.Status,
		common.BKDescriptionField: params.Description}
	valid := validator.NewValidMap(ownerID, common.BKInnerObjIDTempVersion, srvData.header, ps.Engine)
	if err := valid.ValidMap(input, common.ValidCreate, 0); err != nil {
		blog.Errorf("fail to valid input parameters. err:%s,input:%+v,params:%+v,rid:%s", err.Error(), params, input, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommFieldNotValid)})
		return
	}

	input[common.CreateTimeField] = time.Now().UTC()

	ret, err := ps.CoreAPI.CoreService().Instance().CreateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDTempVersion, &meta.CreateModelInstance{Data: input})
	if err != nil {
		blog.Errorf("CreateTemplateVersion CreateObject http do error,err:%s,input:%+v,query:%+v,rid:%s", err.Error(), params, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
	}
	if !ret.Result {
		blog.Errorf("CreateTemplateVersion  CreateObject http response error,err code:%d,err msg:%s,input:%+v,query:%+v,rid:%s", ret.Code, ret.ErrMsg, params, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(ret.Code, ret.ErrMsg)})
		return
	}

	versionID := ret.Data.Created.ID

	//only one online status
	if params.Status == common.TemplateStatusOnline {
		data := types.MapStr{common.BKStatusField: common.TemplateStatusHistory}
		condition := types.MapStr{
			common.BKTemlateIDField: types.MapStr{common.BKDBEQ: versionID},
			common.BKOwnerIDField:   ownerID,
			common.BKStatusField:    common.TemplateStatusOnline}
		input := &meta.UpdateOption{Condition: condition, Data: data}
		ret, err := ps.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDTempVersion, input)
		if err != nil {
			blog.Errorf("CreateTemplateVersion UpdateObject http do error,err:%s,input:%+v,query:%+v,rid:%s", err.Error(), params, input, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		}
		if !ret.Result {
			blog.Errorf("CreateTemplateVersion  UpdateObject http response error,err code:%d,err msg:%s,input:%+v,query:%+v,rid:%s", ret.Code, ret.ErrMsg, params, input, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(ret.Code, ret.ErrMsg)})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) UpdateTemplateVersion(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)

	defErr := srvData.ccErr
	ownerID := srvData.ownerID

	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if nil != err {
		blog.Errorf("update config template failed! derr: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	templateIDStr := req.PathParameter(common.BKTemlateIDField)
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if nil != err {
		blog.Errorf("update config template failed! derr: %v,appIDStr:%s,rid:%s", appIDStr, err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}
	versionIDStr := req.PathParameter(common.BKVersionIDField)

	versionID, err := strconv.ParseInt(versionIDStr, 10, 64)
	if nil != err {
		blog.Errorf("update config template failed! derr: %v,versionIDStr:%s,rid:%s", err, versionIDStr, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}
	var params meta.TemplateVersion
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("create config version failed! decode request body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	condition := types.MapStr{
		common.BKAppIDField:     appID,
		common.BKTemlateIDField: templateID,
		common.BKVersionIDField: versionID}
	data := types.NewFromStruct(params, "field")

	input := &meta.UpdateOption{Condition: condition, Data: data}

	ret, err := ps.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDTempVersion, input)
	if err != nil {
		blog.Errorf("UpdateTemplateVersion UpdateObject http do error,err:%s,input:%+v,query:%+v,rid:%s", err.Error(), params, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
	}
	if !ret.Result {
		blog.Errorf("UpdateTemplateVersion  UpdateObject http response error,err code:%d,err msg:%s,input:%+v,query:%+v,rid:%s", ret.Code, ret.ErrMsg, params, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(ret.Code, ret.ErrMsg)})
		return
	}

	//only one online status
	if params.Status == common.TemplateStatusOnline {
		data := types.MapStr{common.BKStatusField: common.TemplateStatusHistory}
		condition := types.MapStr{
			common.BKTemlateIDField: templateID,
			common.BKVersionIDField: types.MapStr{common.BKDBNE: versionID},
			common.BKOwnerIDField:   ownerID,
			common.BKStatusField:    common.TemplateStatusOnline}

		input := &meta.UpdateOption{Condition: condition, Data: data}
		ret, err := ps.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDTempVersion, input)
		if err != nil {
			blog.Errorf("UpdateTemplateVersion UpdateObject http do error,err:%s,input:%+v,query:%+v,rid:%s", err.Error(), params, input, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		}
		if !ret.Result {
			blog.Errorf("UpdateTemplateVersion  UpdateObject http response error,err code:%d,err msg:%s,input:%+v,query:%+v,rid:%s", ret.Code, ret.ErrMsg, params, input, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(ret.Code, ret.ErrMsg)})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}
