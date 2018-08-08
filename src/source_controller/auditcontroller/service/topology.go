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
	"net/http"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// AddAppLog app操作日志
func (s *Service) AddAppLog(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ownerID := util.GetOwnerID(req.Request.Header)
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := util.GetInt64ByInterface(strAppID)
	if nil != err {
		blog.Errorf("AddAppLog json unmarshal error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditAppParams)
	if err = json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Errorf("AddAppLog json unmarshal failed error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	err = s.Logics.AddLogWithStr(appID, appID, params.OpType, common.BKInnerObjIDApp, params.Content, "", params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("AddAppLog add application log error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

//操作日志
func (s *Service) AddSetLog(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ownerID := util.GetOwnerID(req.Request.Header)
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := util.GetInt64ByInterface(strAppID)
	if nil != err {
		blog.Errorf("AddSetLog json unmarshal error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditSetParams)
	if err = json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Errorf("AddSetLog json unmarshal failed,error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogWithStr(appID, params.SetID, params.OpType, common.BKInnerObjIDSet, params.Content, "", params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("AddSetLog add application log error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

//插入多行主机操作日志型操作
func (s *Service) AddSetLogs(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ownerID := util.GetOwnerID(req.Request.Header)
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := util.GetInt64ByInterface(strAppID)
	if nil != err {
		blog.Errorf("AddSetLogs json unmarshal error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditSetsParams)
	if err = json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Error("AddSetLogs json unmarshal failed,error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogMulti(appID, params.OpType, common.BKInnerObjIDSet, params.Content, params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("AddSetLogs add set log error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

// AddModuleLog 操作日志
func (s *Service) AddModuleLog(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ownerID := util.GetOwnerID(req.Request.Header)
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := util.GetInt64ByInterface(strAppID)
	if nil != err {
		blog.Errorf("AddModuleLog json unmarshal error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditModuleParams)
	if err = json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Error("AddModuleLog json unmarshal failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogWithStr(appID, params.ModuleID, params.OpType, common.BKInnerObjIDModule, params.Content, "", params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("AddModuleLog add module log error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

// AddModuleLogs 插入多行主机操作日志型操作
func (s *Service) AddModuleLogs(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ownerID := util.GetOwnerID(req.Request.Header)
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := util.GetInt64ByInterface(strAppID)
	if nil != err {
		blog.Errorf("AddModuleLogs json unmarshal error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditModulesParams)
	if err = json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Errorf("AddModuleLogs json unmarshal failed, error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogMulti(appID, params.OpType, common.BKInnerObjIDModule, params.Content, params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("add module log error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}
