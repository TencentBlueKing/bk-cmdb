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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	types "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"
)

func (ps *ProcServer) CreateTemplate(req *restful.Request, resp *restful.Response) {
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

	input := types.MapStr{}
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("create config template failed! decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	input[common.BKAppIDField] = appID
	valid := validator.NewValidMap(ownerID, common.BKInnerObjIDConfigTemp, pHeader, ps.Engine)
	if err := valid.ValidMap(input, common.ValidCreate, 0); err != nil {
		blog.Errorf("fail to valid input parameters. err:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommFieldNotValid)})
		return
	}

	tempFields, err := ps.Logics.GetTemplateAttributes(ownerID, pHeader)
	if err != nil {
		blog.Errorf("create config template  err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateTemplateFail)})
		return
	}
	ret, err := ps.CoreAPI.ObjectController().Instance().CreateObject(context.Background(), common.BKInnerObjIDConfigTemp, pHeader, input)
	if nil != err || !ret.Result {
		blog.Errorf("create config template failed by  input :%v, return:%v, err: %v", input, ret, err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateTemplateFail)})
		return
	}

	templateID, err := ret.Data.Int64(common.BKTemlateIDField)
	if nil != err {
		blog.Errorf("create config template failed by  err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateTemplateFail)})
	}

	curData, err := ps.GetTemplateInstanceDetails(pHeader, ownerID, templateID)
	if nil != err {
		blog.Errorf("create config template failed by curData:%v, err: %v, tempID:%v", curData, err, templateID)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateTemplateFail)})
	}

	logContent := auditoplog.AuditLogExt{
		ID: templateID,
		Content: meta.Content{
			PreData: types.MapStr{},
			CurData: curData,
			Headers: tempFields,
		},
	}

	logs := types.MapStr{common.BKContentField: logContent, common.BKOpDescField: "create template", common.BKOpTypeField: auditoplog.AuditOpTypeAdd}
	result, err := ps.CoreAPI.AuditController().AddProcLog(context.Background(), common.BKDefaultOwnerID, appIDStr, user, pHeader, logs)
	if err != nil || !result.Result {
		blog.Errorf("create config template failed, but [%s] audit failed, err: %v, %v", templateID, err, result.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) DeleteTemplate(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)
	pHeader := req.Request.Header
	user := util.GetUser(pHeader)
	var logContent auditoplog.AuditLogExt

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
		blog.Errorf("intput params err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsInvalid)})
		return
	}

	tempFields, err := ps.Logics.GetTemplateAttributes(ownerID, pHeader)
	if nil != err {
		blog.Errorf("delete config template  err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}
	logger := ps.Logics.NewTemplate(pHeader, ownerID)
	if err := logger.WithPrevious(templateID, tempFields); err != nil {
		blog.Errorf("delete template, but get temp[%d] pre data for audit failed, err: %v", templateID, err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteInstanceModel)})
		return
	}

	input := types.MapStr{
		common.BKOwnerIDField:   ownerID,
		common.BKAppIDField:     appID,
		common.BKTemlateIDField: templateID}

	ret, err := ps.CoreAPI.ObjectController().Instance().DelObject(context.Background(), common.BKInnerObjIDConfigTemp, pHeader, input)
	if err != nil || !ret.Result {
		blog.Errorf("delete config template failed by  intput :%v, return:%v, err: %v", input, ret, err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}

	logContent = *logger.AuditLog(templateID)
	logs := types.MapStr{common.BKContentField: logContent, common.BKOpDescField: "delete template", common.BKOpTypeField: auditoplog.AuditOpTypeDel}
	result, err := ps.CoreAPI.AuditController().AddProcLog(context.Background(), common.BKDefaultOwnerID, appIDStr, user, pHeader, logs)
	if err != nil || !result.Result {
		blog.Errorf("delete config template failed, but [%d] audit failed, err: %v, %v", templateID, err, result.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) UpdateTemplate(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)
	pHeader := req.Request.Header
	user := util.GetUser(pHeader)
	rid := util.GetHTTPCCRequestID(pHeader)
	var logContent auditoplog.AuditLogExt

	ownerID := req.PathParameter(common.BKOwnerIDField)
	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if nil != err {
		blog.Errorf("create config template failed! derr: %v,rid:%s", err, rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	templateIDStr := req.PathParameter(common.BKTemlateIDField)
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)

	if nil != err {
		blog.Errorf("intput params err: %v,rid:%s", err, rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsInvalid)})
		return
	}

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("update config template failed! decode request body err: %v,rid:%s", err, rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	valid := validator.NewValidMap(ownerID, common.BKInnerObjIDConfigTemp, req.Request.Header, ps.Engine)
	if err := valid.ValidMap(input, common.ValidUpdate, int64(templateID)); err != nil {
		blog.Errorf("fail to valid input parameters. err:%s,rid:%s", err.Error(), rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommFieldNotValid)})
		return
	}

	tempFields, err := ps.GetTemplateAttributes(ownerID, pHeader)
	if nil != err {
		blog.Errorf("delete config template  err: %v,rid:%s", err, rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}
	logger := ps.Logics.NewTemplate(pHeader, ownerID)
	if err := logger.WithPrevious(templateID, tempFields); err != nil {
		blog.Errorf("delete template, but get temp[%v] pre data for audit failed, err: %v,rid:%s", templateID, err, rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteInstanceModel)})
		return
	}

	data := types.MapStr{
		"data": input,
		"condition": map[string]interface{}{
			common.BKTemlateIDField: templateID,
			common.BKAppIDField:     appID,
			common.BKOwnerIDField:   ownerID,
		},
	}
	ret, err := ps.CoreAPI.ObjectController().Instance().UpdateObject(context.Background(), common.BKInnerObjIDConfigTemp, req.Request.Header, data)
	if nil != err || !ret.Result {
		blog.Errorf("update config template failed by processcontroll. err: %v, errcode: %d, errmsg: %s,rid:%s", err, ret.Code, ret.ErrMsg, rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcUpdateTemplateFail)})
		return
	}

	if err := logger.WithCurrent(templateID); err != nil {
		blog.Errorf("delete config template, but get current host data failed, err: %v,rid:%s", err, rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}

	logContent = *logger.AuditLog(templateID)
	logs := common.KvMap{common.BKContentField: logContent, common.BKOpDescField: "update template", common.BKOpTypeField: auditoplog.AuditOpTypeModify}
	result, err := ps.CoreAPI.AuditController().AddProcLog(context.Background(), common.BKDefaultOwnerID, appIDStr, user, pHeader, logs)
	if nil != err || !result.Result {
		blog.Errorf("delete config template failed, but add template[%v] audit failed, err: %v, %v,rid:%s", templateID, err, result.ErrMsg, rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) SearchTemplate(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	var params meta.SearchParams
	var input meta.QueryInput
	var err error
	var ok bool
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input.Condition = params.Condition
	input.Fields = strings.Join(params.Fields, ",")
	input.Start, err = util.GetIntByInterface(params.Page["start"])
	if nil != err {
		blog.Errorf("request body query condition format error start not integer, input:%v", params.Page["start"])
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, "start")})
		return
	}
	input.Limit, err = util.GetIntByInterface(params.Page["limit"])
	if nil != err {
		blog.Errorf("request body query condition format error limit not integer, input:%v", params.Page["limit"])
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, "limit")})
		return
	}
	input.Sort, ok = params.Page["sort"].(string)
	if false == ok {
		input.Sort = ""
	}

	ret, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDConfigTemp, req.Request.Header, &input)
	if err != nil || (err == nil && !ret.Result) {
		blog.Errorf("query config template failed by processcontroll. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(ret.Data))
}
