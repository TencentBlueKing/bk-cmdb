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

package inst

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/topo_service/logics"
	auditlogAPI "configcenter/src/source_controller/api/auditlog"
	"configcenter/src/source_controller/common/commondata"
)

var audit = &auditAction{}

// instAction
type auditAction struct {
	base.BaseAction
}

func init() {

	// register action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/audit/search", Params: nil, Handler: audit.Query})

	// create cc
	audit.CreateAction()
}

// Query query auditlog
func (cli *auditAction) Query(req *restful.Request, resp *restful.Response) {

	// get language
	language := util.GetActionLanguage(req)

	// get the error object by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		var dat commondata.ObjQueryInput
		err = json.Unmarshal([]byte(value), &dat)
		if err != nil {
			blog.Error("get audit input:%v error:%v", value, err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		//user := sencecommon.GetUserFromHeader(req)
		ownerID := common.BKDefaultOwnerID
		iConds := dat.Condition
		if nil == iConds {
			dat.Condition = common.KvMap{common.BKOwnerIDField: ownerID}
		} else {
			conds := iConds.(map[string]interface{})
			times, ok := conds[common.BKOpTimeField].([]interface{})
			if ok {
				if 2 != len(times) {
					blog.Error("search operation log input params times error, info: %v", times)
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommParamsInvalid)
				}

				conds[common.BKOpTimeField] = common.KvMap{"$gte": times[0], "$lte": times[1], commondata.CC_time_type_parse_flag: "1"}
				//delete(conds, "Time")
			}
			conds[common.BKOwnerIDField] = ownerID
			dat.Condition = conds
		}
		if 0 == dat.Limit {
			dat.Limit = common.BKDefaultLimit
		}

		client := auditlogAPI.NewClient(cli.CC.AuditCtrl(), req.Request.Header)
		ret, err := client.GetAuditlogs(dat)
		blog.Debug("search operation log  params: %v", dat)
		if nil != err {
			blog.Error("search operation log error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommParamsInvalid)
		}
		ret = logics.TranslateOpLanguage(ret, defLang)

		return http.StatusOK, ret, nil
	}, resp)

}
