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
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	types "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"
)

func (ps *ProcServer) SearchTemplateVersion(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	ownerID := req.PathParameter(common.BKOwnerIDField)
	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if nil != err {
		blog.Errorf("search template version failed! derr: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	templateIDStr := req.PathParameter(common.BKTemlateIDField)
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if nil != err {
		blog.Errorf("search template version failed! derr: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	params := types.MapStr{}
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	var input meta.QueryInput

	conditon := types.MapStr{common.BKOwnerIDField: ownerID, common.BKAppIDField: appID, common.BKTemlateIDField: templateID}
	status, ok := params[common.BKStatusField]
	if ok {
		conditon[common.BKStatusField] = status
	}
	input.Condition = conditon
	input.Fields = ""
	input.Start = 0
	input.Limit = common.BKNoLimit
	input.Sort = ""

	ret, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDTempVersion, req.Request.Header, &input)
	if err != nil || !ret.Result {
		blog.Errorf("query config template failed by processcontroll. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(ret.Data.Info))
}

func (ps *ProcServer) CreateTemplateVersion(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)
	pHeader := req.Request.Header
	user := util.GetUser(pHeader)
	ownerID := req.PathParameter(common.BKOwnerIDField)
	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if nil != err {
		blog.Errorf("create config template failed! derr: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	templateIDStr := req.PathParameter(common.BKTemlateIDField)
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if nil != err {
		blog.Errorf("create config template failed! derr: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	var params meta.TemplateVersion

	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("create config version failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input := types.MapStr{common.BKAppIDField: appID,
		common.BKOperatorField:    user,
		common.BKTemlateIDField:   templateID,
		common.BKContentField:     params.Content,
		common.BKStatusField:      params.Status,
		common.BKDescriptionField: params.Description}
	valid := validator.NewValidMap(ownerID, common.BKInnerObjIDTempVersion, pHeader, ps.Engine)
	if err := valid.ValidMap(input, common.ValidCreate, 0); err != nil {
		blog.Errorf("fail to valid input parameters. err:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommFieldNotValid)})
		return
	}

	input[common.CreateTimeField] = time.Now().UTC()

	ret, err := ps.CoreAPI.ObjectController().Instance().CreateObject(context.Background(), common.BKInnerObjIDTempVersion, pHeader, input)
	if nil != err || !ret.Result {
		blog.Errorf("create template version failed by  input :%v, return:%v, err: %v", input, ret, err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateTemplateFail)})
		return
	}

	versionID, err := ret.Data.Int64(common.BKVersionIDField)
	if nil != err {
		blog.Errorf("create template version failed by  input :%v, return:%v, err: %v", input, ret, err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateTemplateFail)})
		return
	}

	//only one online status
	if params.Status == common.TemplateStatusOnline {
		data := types.MapStr{common.BKStatusField: common.TemplateStatusHistory}
		condition := types.MapStr{
			common.BKTemlateIDField: types.MapStr{common.BKDBEQ: versionID},
			common.BKOwnerIDField:   ownerID,
			common.BKStatusField:    common.TemplateStatusOnline}
		input := types.MapStr{"condition": condition, "data": data}
		ret, err := ps.CoreAPI.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDTempVersion, pHeader, input)
		if nil != err || !ret.Result {
			blog.Errorf("create config template failed by  input :%v, return:%v, err: %v", input, ret, err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateTemplateFail)})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) UpdateTemplateVersion(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)
	pHeader := req.Request.Header
	ownerID := req.PathParameter(common.BKOwnerIDField)

	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if nil != err {
		blog.Errorf("update config template failed! derr: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	templateIDStr := req.PathParameter(common.BKTemlateIDField)
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if nil != err {
		blog.Errorf("update config template failed! derr: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}
	versionIDStr := req.PathParameter(common.BKVersionIDField)

	versionID, err := strconv.ParseInt(versionIDStr, 10, 64)
	if nil != err {
		blog.Errorf("update config template failed! derr: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}
	var params meta.TemplateVersion

	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("create config version failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	condition := types.MapStr{
		common.BKAppIDField:     appID,
		common.BKTemlateIDField: templateID,
		common.BKVersionIDField: versionID}
	input := types.MapStr{"condition": condition, "data": params}

	ret, err := ps.CoreAPI.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDTempVersion, pHeader, input)
	if nil != err || !ret.Result {
		blog.Errorf("create config template failed by  input :%v, return:%v, err: %v", input, ret, err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateTemplateFail)})
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

		input := types.MapStr{"condition": condition, "data": data}
		ret, err := ps.CoreAPI.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDTempVersion, pHeader, input)
		if nil != err || !ret.Result {
			blog.Errorf("update template version failed by  input :%v, return:%v, err: %v", input, ret, err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateTemplateFail)})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}
