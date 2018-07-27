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

//主机操作日志
func (s *Service) AddHostLog(req *restful.Request, resp *restful.Response) {

	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ownerID := util.GetOwnerID(req.Request.Header)
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := util.GetInt64ByInterface(strAppID)
	if nil != err {
		blog.Errorf("AddHostLog json unmarshal error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditHostLogParams)
	if err := json.NewDecoder(req.Request.Body).Decode(params); nil != err {
		blog.Errorf("AddHostLog json unmarshal  error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogWithStr(appID, params.HostID, params.OpType, common.BKInnerObjIDHost, params.Content, params.InnerIP, params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("AddHostLog add host log error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

//插入多行主机操作日志型操作
func (s *Service) AddHostLogs(req *restful.Request, resp *restful.Response) {

	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ownerID := util.GetOwnerID(req.Request.Header)
	strAppID := req.PathParameter("biz_id")
	user := req.PathParameter("user")

	appID, err := util.GetInt64ByInterface(strAppID)
	if nil != err {
		blog.Errorf("AddHostLogs json unmarshal error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}

	params := new(metadata.AuditHostsLogParams)
	if err = json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Error("AddHostLogs json unmarshal failed, error:%v", err)
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	err = s.Logics.AddLogMultiWithExtKey(appID, params.OpType, common.BKInnerObjIDHost, params.Content, params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("AddHostLogs add host log error:%s", err.Error())
		resp.WriteError(http.StatusBadGateway, &metadata.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}
