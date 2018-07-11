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
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/auditcontroller/audit/logics"
	"configcenter/src/source_controller/common/commondata"
	"encoding/json"
	"io/ioutil"

	restful "github.com/emicklei/go-restful"
)

var queryAudit *queryAuditAction = &queryAuditAction{}

// ObjectAction
type queryAuditAction struct {
	base.BaseAction
}

func init() {
	queryAudit.CreateAction()

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/search", Params: nil, Handler: queryAudit.Get})
	// set cc api resource
}

func (q *queryAuditAction) Get(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := q.CC.Error.CreateDefaultCCErrorIf(language)
	ownerID := util.GetActionOnwerID(req)
	value, err := ioutil.ReadAll(req.Request.Body)

	if err != nil {
		blog.Errorf("read http request boody error:%s", err.Error())
		q.ResponseFailed(common.CCErrCommHTTPReadBodyFailed, defErr.Error(common.CCErrCommHTTPReadBodyFailed).Error(), resp)
		return
	}
	var dat commondata.ObjQueryInput
	err = json.Unmarshal([]byte(value), &dat)
	if err != nil {
		blog.Error("json unmarshal failed,input:%v error:%v", string(value), err)
		q.ResponseFailed(common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}
	dat.Condition = util.SetModOwner(dat.Condition, ownerID)
	logics.DB = appAudit.CC.InstCli
	rows, cnt, err := logics.Search(dat)
	if nil != err {
		blog.Error("get data from data  error:%s", err.Error())
		q.ResponseFailed(common.CCErrCommDBSelectFailed, defErr.Error(common.CCErrCommDBSelectFailed).Error(), resp)
		return
	}
	data := common.KvMap{"info": rows, "count": cnt}
	queryAudit.ResponseSuccess(data, resp)
}
