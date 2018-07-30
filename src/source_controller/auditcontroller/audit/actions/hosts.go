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

package actions

import (
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/auditcontroller/audit/logics"
	"encoding/json"
	"io/ioutil"
	"strconv"

	restful "github.com/emicklei/go-restful"
)

var hostAudit *hostAuditAction = &hostAuditAction{}

// ObjectAction
type hostAuditAction struct {
	base.BaseAction
}

func init() {
	hostAudit.CreateAction()

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/host/{owner_id}/{biz_id}/{user}", Params: nil, Handler: hostAudit.AddHostLog})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/{owner_id}/{biz_id}/{user}", Params: nil, Handler: hostAudit.AddHostLogs})

	// set cc api resource
}

//主机操作日志
func (a *hostAuditAction) AddHostLog(req *restful.Request, resp *restful.Response) {
	type paramsStruct struct {
		Content interface{}            `json:"content"`
		OpDesc  string                 `json:"op_desc"`
		InnerIP string                 `json:"bk_host_innerip"`
		OpType  auditoplog.AuditOpType `json:"op_type"`
		HostID  int                    `json:"inst_id"`
	}
	language := util.GetActionLanguage(req)
	defErr := a.CC.Error.CreateDefaultCCErrorIf(language)

	ownerID := util.GetActionOnwerID(req)
	strAppID := req.PathParameter("biz_id")
	appID, _ := strconv.Atoi(strAppID)
	user := req.PathParameter("user")
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request boody error:%s", err.Error())
		hostAudit.ResponseFailed(common.CCErrCommHTTPReadBodyFailed, defErr.Error(common.CCErrCommHTTPReadBodyFailed).Error(), resp)
		return
	}
	params := paramsStruct{}
	err = json.Unmarshal([]byte(value), &params)
	if err != nil {
		blog.Error("json unmarshal failed:%v error:%v", string(value), err)
		hostAudit.ResponseFailed(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}

	logics.DB = appAudit.CC.InstCli
	err = logics.AddLogWithStr(appID, params.HostID, params.OpType, common.BKInnerObjIDHost, params.Content, params.InnerIP, params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("add host log error:%s", err.Error())
		hostAudit.ResponseFailed(common.CCErrCommDBInsertFailed, defErr.Error(common.CCErrCommDBInsertFailed).Error(), resp)
		return
	} else {
		hostAudit.ResponseSuccess(nil, resp)
		return
	}

}

//插入多行主机操作日志型操作
func (a *hostAuditAction) AddHostLogs(req *restful.Request, resp *restful.Response) {
	type paramsStruct struct {
		Content []auditoplog.AuditLogExt `json:"content"`
		OpDesc  string                   `json:"op_desc"`
		OpType  auditoplog.AuditOpType   `json:"op_type"`
	}
	language := util.GetActionLanguage(req)
	defErr := a.CC.Error.CreateDefaultCCErrorIf(language)

	ownerID := util.GetActionOnwerID(req)
	strAppID := req.PathParameter("biz_id")
	appID, _ := strconv.Atoi(strAppID)
	user := req.PathParameter("user")
	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Errorf("read http request boody error:%s", err.Error())
		hostAudit.ResponseFailed(common.CCErrCommHTTPReadBodyFailed, defErr.Error(common.CCErrCommHTTPReadBodyFailed).Error(), resp)
		return
	}

	params := paramsStruct{}
	err = json.Unmarshal([]byte(value), &params)
	if err != nil {
		blog.Error("json unmarshal failed,input:%v error:%v", string(value), err)
		hostAudit.ResponseFailed(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}

	logics.DB = appAudit.CC.InstCli
	err = logics.AddLogMultiWithExtKey(appID, params.OpType, common.BKInnerObjIDHost, params.Content, params.OpDesc, ownerID, user)
	if nil != err {
		blog.Errorf("add host log error:%s", err.Error())
		hostAudit.ResponseFailed(common.CCErrCommDBInsertFailed, defErr.Error(common.CCErrCommDBInsertFailed).Error(), resp)
		return
	} else {
		hostAudit.ResponseSuccess(nil, resp)
		return
	}

}
