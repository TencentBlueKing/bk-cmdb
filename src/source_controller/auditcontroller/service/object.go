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

// AddLog 操作日志
func (s *Service) AddObjectLog(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ownerID := util.GetOwnerID(req.Request.Header)
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")
	appID, err := util.GetInt64ByInterface(strAppID)
	if nil != err {
		blog.Errorf("AddObjectLog json unmarshal error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditObjParams) //paramsStruct{}
	if json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Errorf("AddObjectLog json unmarshal failed,  error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogWithStr(appID, params.InstID, params.OpType, params.OpTarget, params.Content, "", params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("AddObjectLog add module log error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	} else {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

}

// AddLogs 插入多行主机操作日志型操作
func (s *Service) AddObjectLogs(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ownerID := util.GetOwnerID(req.Request.Header)
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := util.GetInt64ByInterface(strAppID)
	if nil != err {
		blog.Errorf("AddObjectLogs json unmarshal error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditObjsParams)
	if json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Errorf("AddObjectLogs json unmarshal failed,  error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogMulti(appID, params.OpType, params.OpTarget, params.Content, params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("AddObjectLogs add module log error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	} else {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

}

// AddProcLog 操作日志
func (s *Service) AddProcLog(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ownerID := util.GetOwnerID(req.Request.Header)
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := util.GetInt64ByInterface(strAppID)
	if nil != err {
		blog.Errorf("AddProcLog json unmarshal error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}
	params := new(metadata.AuditProcParams)
	if json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Errorf("AddProcLog json unmarshal failed,  error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogWithStr(appID, params.ProcID, params.OpType, common.BKInnerObjIDProc, params.Content, "", params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("AddProcLog json unmarshal failed,input:%v error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

// AddProcLogs  插入多行主机操作日志型操作
func (s *Service) AddProcLogs(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ownerID := util.GetOwnerID(req.Request.Header)
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := util.GetInt64ByInterface(strAppID)
	if nil != err {
		blog.Errorf("AddProcLogs json unmarshal error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditProcsParams)
	if json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Errorf("AddProcLogs json unmarshal failed,  error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogMulti(appID, params.OpType, common.BKInnerObjIDProc, params.Content, params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("AddProcLogs json unmarshal failed,  error:%v", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}
