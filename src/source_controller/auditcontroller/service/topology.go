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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// AddAppLog app操作日志
func (s *Service) AddAppLog(req *restful.Request, resp *restful.Response) {
	type paramsStruct struct {
		Content string                 `json:"content"`
		OpDesc  string                 `json:"op_desc"`
		OpType  auditoplog.AuditOpType `json:"op_type"`
	}

	language := util.GetActionLanguage(req)
	defErr := a.CC.Error.CreateDefaultCCErrorIf(language)

	ownerID := req.PathParameter("owner_id")
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := strconv.Atoi(strAppID)
	if nil != err {
		blog.Errorf("AddAppLog json unmarshal error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.ErrorF(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditAppParams)
	if err = json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Error("AddAppLog json unmarshal failed,input:%v error:%v", string(value), err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	err = s.Logics.AddLogWithStr(appID, appID, params.OpType, common.BKInnerObjIDApp, params.Content, "", params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("AddAppLog add application log error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	} else {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

}

//操作日志
func (s *Service) AddSetLog(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := a.CC.Error.CreateDefaultCCErrorIf(language)

	ownerID := req.PathParameter("owner_id")
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := strconv.Atoi(strAppID)
	if nil != err {
		blog.Errorf("AddSetLog json unmarshal error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.ErrorF(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditSetParams)
	if err = json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Error("AddSetLog json unmarshal failed,input:%v error:%v", string(value), err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogWithStr(appID, params.SetID, params.OpType, common.BKInnerObjIDSet, params.Content, "", params.OpDesc, ownerID, user)
	if nil != err {
		blog.Error("AddSetLog json unmarshal failed,input:%v error:%v", string(value), err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	} else {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

}

//插入多行主机操作日志型操作
func (s *Service) AddSetLogs(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := a.CC.Error.CreateDefaultCCErrorIf(language)

	ownerID := req.PathParameter("owner_id")
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := strconv.Atoi(strAppID)
	if nil != err {
		blog.Errorf("AddSetLogs json unmarshal error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.ErrorF(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditSetsParams)
	err = json.Unmarshal([]byte(value), &params)
	if err = json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Error("AddSetLogs json unmarshal failed,input:%v error:%v", string(value), err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogMulti(appID, params.OpType, common.BKInnerObjIDSet, params.Content, params.OpDesc, ownerID, user)
	if nil != err {
		blog.Error("AddSetLogs json unmarshal failed,input:%v error:%v", string(value), err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	} else {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

}

// AddModuleLog 操作日志
func (s *Service) AddModuleLog(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := a.CC.Error.CreateDefaultCCErrorIf(language)

	ownerID := req.PathParameter("owner_id")
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := strconv.Atoi(strAppID)
	if nil != err {
		blog.Errorf("AddModuleLog json unmarshal error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.ErrorF(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditModuleParams)
	if err = json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Error("AddModuleLog json unmarshal failed,input:%v error:%v", string(value), err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogWithStr(appID, params.ModuleID, params.OpType, common.BKInnerObjIDModule, params.Content, "", params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("AddModuleLog add module log error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	} else {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

}

// AddModuleLogs 插入多行主机操作日志型操作
func (s *Service) AddModuleLogs(req *restful.Request, resp *restful.Response) {
	type paramsStruct struct {
		Content []auditoplog.AuditLogContext `json:"content"`
		OpDesc  string                       `json:"op_desc"`
		OpType  auditoplog.AuditOpType       `json:"op_type"`
	}
	language := util.GetActionLanguage(req)
	defErr := a.CC.Error.CreateDefaultCCErrorIf(language)

	ownerID := req.PathParameter("owner_id")
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := strconv.Atoi(strAppID)
	if nil != err {
		blog.Errorf("AddModuleLogs json unmarshal error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.ErrorF(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditModulesParams)
	if err = json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Error("AddModuleLogs json unmarshal failed,input:%v error:%v", string(value), err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogMulti(appID, params.OpType, common.BKInnerObjIDModule, params.Content, params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("add module log error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	} else {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

}
