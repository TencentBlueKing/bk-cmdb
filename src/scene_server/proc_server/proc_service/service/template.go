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

	input := types.MapStr{}
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("create config template failed! decode request body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	input[common.BKAppIDField] = appID
	valid := validator.NewValidMap(ownerID, common.BKInnerObjIDConfigTemp, srvData.header, ps.Engine)
	if err := valid.ValidMap(input, common.ValidCreate, 0); err != nil {
		blog.Errorf("fail to valid input parameters. err:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommFieldNotValid)})
		return
	}

	tempFields, err := srvData.lgc.GetTemplateAttributes(srvData.ctx, ownerID)
	if err != nil {
		blog.Errorf("create config template  err: %v,appIDStr:%s,rid:%s", appIDStr, err, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateTemplateFail)})
		return
	}
	ret, err := ps.CoreAPI.CoreService().Instance().CreateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDConfigTemp, &meta.CreateModelInstance{Data: input})
	if nil != err {
		blog.Errorf("CreateTemplate CreateObject http do error.  err:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !ret.Result {
		blog.Errorf("CreateTemplate CreateObject http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", ret.Code, ret.Result, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.New(ret.Code, ret.ErrMsg)})
		return
	}

	templateID := int64(ret.Data.Created.ID)
	blog.Errorf("create config template failed by  err: %v,appID:%v,inst info:%+v,rid:%s", err, appID, ret.Data, srvData.rid)
	resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcCreateTemplateFail)})

	curData, err := srvData.lgc.GetTemplateInstanceDetails(srvData.ctx, templateID)
	if nil != err {
		blog.Errorf("create config template failed by curData:%v, err: %v, tempID:%v,rid:%s", curData, err, templateID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: err})
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
	result, err := ps.CoreAPI.AuditController().AddProcLog(srvData.ctx, common.BKDefaultOwnerID, appIDStr, user, srvData.header, logs)
	if err != nil || !result.Result {
		blog.Errorf("create config template failed, but [%s] audit failed, err: %v, %v,rid:%s", templateID, err, result.ErrMsg, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) DeleteTemplate(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr
	user := srvData.user
	var logContent auditoplog.AuditLogExt

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
		blog.Errorf("intput params err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsInvalid)})
		return
	}

	tempFields, err := srvData.lgc.GetTemplateAttributes(srvData.ctx, ownerID)
	if nil != err {
		blog.Errorf("delete config template  err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}
	logger := srvData.lgc.NewTemplate()
	if err := logger.WithPrevious(srvData.ctx, templateID, tempFields); err != nil {
		blog.Errorf("delete template, but get temp[%d] pre data for audit failed, err: %v,rid:%s", templateID, err, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteInstanceModel)})
		return
	}

	input := &meta.DeleteOption{Condition: types.MapStr{
		common.BKOwnerIDField:   ownerID,
		common.BKAppIDField:     appID,
		common.BKTemlateIDField: templateID},
	}

	ret, err := ps.CoreAPI.CoreService().Instance().DeleteInstance(srvData.ctx, srvData.header, common.BKInnerObjIDConfigTemp, input)
	if nil != err {
		blog.Errorf("DeleteTemplate DelObject http do error.  err:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !ret.Result {
		blog.Errorf("DeleteTemplate DelObject http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", ret.Code, ret.Result, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.New(ret.Code, ret.ErrMsg)})
		return
	}

	logContent = *logger.AuditLog(templateID)
	logs := types.MapStr{common.BKContentField: logContent, common.BKOpDescField: "delete template", common.BKOpTypeField: auditoplog.AuditOpTypeDel}
	result, err := ps.CoreAPI.AuditController().AddProcLog(srvData.ctx, srvData.ownerID, appIDStr, user, srvData.header, logs)
	if err != nil || !result.Result {
		blog.Errorf("delete config template failed, but [%d] audit failed, err: %v, %v,input:%+v,rid:%s", templateID, err, result.ErrMsg, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) UpdateTemplate(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr
	user := srvData.user
	var logContent auditoplog.AuditLogExt

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
		blog.Errorf("intput params err: %v,templateID:%v,rid:%s", err, templateIDStr, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsInvalid)})
		return
	}

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("update config template failed! decode request body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	valid := validator.NewValidMap(ownerID, common.BKInnerObjIDConfigTemp, req.Request.Header, ps.Engine)
	if err := valid.ValidMap(input, common.ValidUpdate, int64(templateID)); err != nil {
		blog.Errorf("fail to valid input parameters. err:%s,rid:%s", err.Error(), srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommFieldNotValid)})
		return
	}

	tempFields, err := srvData.lgc.GetTemplateAttributes(srvData.ctx, ownerID)
	if nil != err {
		blog.Errorf("delete config template  err: %v,templateID:%v,rid:%s", err, templateID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}
	logger := srvData.lgc.NewTemplate()
	if err := logger.WithPrevious(srvData.ctx, templateID, tempFields); err != nil {
		blog.Errorf("delete template, but get temp[%v] pre data for audit failed, err: %v,rid:%s", templateID, err, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteInstanceModel)})
		return
	}

	data := &meta.UpdateOption{
		Data: input,
		Condition: map[string]interface{}{
			common.BKTemlateIDField: templateID,
			common.BKAppIDField:     appID,
			common.BKOwnerIDField:   ownerID,
		},
	}
	ret, err := ps.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDConfigTemp, data)
	if nil != err {
		blog.Errorf("DeleteTemplate DelObject http do error.  err:%s, input:%+v,rid:%s", err.Error(), data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !ret.Result {
		blog.Errorf("DeleteTemplate DelObject http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", ret.Code, ret.Result, data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.New(ret.Code, ret.ErrMsg)})
		return
	}

	if err := logger.WithCurrent(srvData.ctx, templateID); err != nil {
		blog.Errorf("delete config template, but get current host data failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}

	logContent = *logger.AuditLog(templateID)
	logs := common.KvMap{common.BKContentField: logContent, common.BKOpDescField: "update template", common.BKOpTypeField: auditoplog.AuditOpTypeModify}
	result, err := ps.CoreAPI.AuditController().AddProcLog(srvData.ctx, srvData.ownerID, appIDStr, user, srvData.header, logs)
	if nil != err || !result.Result {
		blog.Errorf("delete config template failed, but add template[%s] audit failed, err: %v, %v,input:%+v,rid:%s", templateID, err, result.ErrMsg, data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcDeleteTemplateFail)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

func (ps *ProcServer) SearchTemplate(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	var params meta.SearchParams
	input := new(meta.QueryCondition)
	var err error
	var ok bool
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("decode request body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input.Condition = params.Condition
	input.Fields = params.Fields
	input.Limit.Offset, err = util.GetInt64ByInterface(params.Page["start"])
	if nil != err {
		blog.Errorf("request body query condition format error start not integer, input:%+v,rid:%s", params.Page["start"], srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, "start")})
		return
	}
	input.Limit.Limit, err = util.GetInt64ByInterface(params.Page["limit"])
	if nil != err {
		blog.Errorf("request body query condition format error limit not integer, input:%+v,rid:%s", params.Page["limit"], srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, "limit")})
		return
	}
	sort, ok := params.Page["sort"].(string)
	if false == ok {
		input.SortArr = nil
	} else {
		input.SortArr = meta.NewSearchSortParse().String(sort).ToSearchSortArr()
	}

	ret, err := ps.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDConfigTemp, input)
	if nil != err {
		blog.Errorf("SearchTemplate SearchObjects http do error.  err:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !ret.Result {
		blog.Errorf("SearchTemplate SearchObjects http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", ret.Code, ret.Result, input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.New(ret.Code, ret.ErrMsg)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(ret.Data))
}

func (ps *ProcServer) GetTemplateGroup(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	appIDStr := req.PathParameter(common.BKAppIDField)
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if nil != err {
		blog.Errorf("get config template group failed! derr: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	var input meta.QueryInput
	input.Condition = types.MapStr{common.BKAppIDField: appID}
	input.Fields = ""
	input.Start = 0
	input.Limit = common.BKNoLimit
	input.Sort = ""

	ret, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(srvData.ctx, common.BKInnerObjIDConfigTemp, req.Request.Header, &input)
	if err != nil || (err == nil && !ret.Result) {
		blog.Errorf("query config template failed by processcontroll. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	result := make([]string, 0)
	for _, data := range ret.Data.Info {
		group, ok := data[common.BKGroupField].(string)
		if false == ok {
			continue
		}
		if !util.InArray(group, result) {
			result = append(result, group)
		}

	}

	resp.WriteEntity(meta.NewSuccessResp(result))
}
