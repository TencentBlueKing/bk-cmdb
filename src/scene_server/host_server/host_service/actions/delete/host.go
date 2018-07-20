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

package delete

import (
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	sceneCommon "configcenter/src/scene_server/common"
	"configcenter/src/scene_server/host_server/host_service/logics"
	"configcenter/src/source_controller/api/auditlog"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"

	"fmt"
)

var host *hostAction = &hostAction{}

type hostAction struct {
	base.BaseAction
}

type AppResult struct {
	Result  bool        `json:result`
	Code    int         `json:code`
	Message interface{} `json:message`
	Data    DataInfo    `json:data`
}

type DataInfo struct {
	Count int                      `json:count`
	Info  []map[string]interface{} `json:info`
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/host/batch", Params: nil, Handler: host.DeleteHostBatch})
	// create CC object
	host.CreateAction()
}

//DeleteHostBatch batch delete host
func (cli *hostAction) DeleteHostBatch(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		objCtrl := cli.CC.ObjCtrl()
		hostCtrl := cli.CC.HostCtrl()

		ownerID, user := util.GetActionOnwerIDAndUser(req)
		//delete host
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read input error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)

		}
		js, err := simplejson.NewJson([]byte(value))
		if nil != err {
			blog.Error("params can not be decode error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)

		}
		data, err := js.Map()
		if nil != err {
			blog.Error("params can not be decode to map error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)

		}
		hostIDStr, ok := data["bk_host_id"].(string)
		if false == ok {
			blog.Error("params do not hava host id:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)

		}

		//gAppURL := "http://" + cli.CC.ObjCtrl + "/object/v1/insts/"+common.BKInnerObjIDApp+"/search"
		gAppURL := objCtrl + "/object/v1/insts/" + common.BKInnerObjIDApp + "/search"
		cond := make(map[string]interface{})
		cond[common.BKDefaultField] = 1
		cond[common.BKOwnerIDField] = ownerID
		condition := make(map[string]interface{})
		condition["condition"] = cond
		conditionStr, _ := json.Marshal(condition)
		appResult, err := httpcli.ReqHttp(req, gAppURL, common.HTTPSelectPost, []byte(conditionStr))
		if nil != err {
			blog.Error("request failed:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)

		}

		var appData AppResult
		err = json.Unmarshal([]byte(appResult), &appData)
		if nil != err {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)

		}
		appInfo := appData.Data.Info
		if len(appInfo) == 0 {
			blog.Error("not found failed:%s", appResult)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)

		}
		var appID int
		appCellInfo, ok := appInfo[0][common.BKAppIDField]
		if false == ok {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)

		}
		appID64, _ := appCellInfo.(float64)
		appID = int(appID64)
		hostIDArr := strings.Split(hostIDStr, ",")
		var iHostIDArr []int
		for _, i := range hostIDArr {
			iHostID, _ := strconv.Atoi(i)
			iHostIDArr = append(iHostIDArr, iHostID)
		}

		dMhConfigURL := hostCtrl + "/host/v1/meta/hosts/modules"
		hostFields, _ := logics.GetHostLogFields(req, ownerID, objCtrl)
		var logConents []auditoplog.AuditLogExt
		for _, hostID := range iHostIDArr {
			strHostID := fmt.Sprintf("%d", hostID)
			logObj := logics.NewHostLog(req, ownerID, strHostID, hostCtrl, objCtrl, hostFields)
			input := make(map[string]interface{})
			input[common.BKHostIDField] = hostID
			input[common.BKAppIDField] = appID
			inputJson, _ := json.Marshal(input)
			blog.Info("delete module host config batch url:%s", dMhConfigURL)
			blog.Info("delete module host config content:%s", string(inputJson))
			result, err := httpcli.ReqHttp(req, dMhConfigURL, common.HTTPDelete, []byte(inputJson))
			blog.Info("delete module host config return:%s", string(result))
			if nil != err {
				blog.Error("delete host batch fail:%v", err)
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostDeleteFail)

			}
			err = sceneCommon.DeleteInstAssociation(cli.CC.ObjCtrl(), req, hostID, ownerID, common.BKInnerObjIDHost, "")
			if nil != err {
				blog.Error("delete host batch fail:%v", err)
				return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostDeleteFail)

			}
			logContent, _ := logObj.GetHostLog(strHostID, true)

			logConents = append(logConents, auditoplog.AuditLogExt{ID: hostID, Content: logContent, ExtKey: logObj.GetInnerIP()})
		}

		hostCond := make(map[string]interface{})
		condInput := make(map[string]interface{})
		hostCond[common.BKDBIN] = iHostIDArr
		condInput[common.BKHostIDField] = hostCond

		dHostURL := objCtrl + "/object/v1/insts/host"

		inputJson, _ := json.Marshal(condInput)
		blog.Info("delete host batch url:%s", dHostURL)
		blog.Info("delete host batch content:%s", string(inputJson))
		_, err = httpcli.ReqHttp(req, dHostURL, common.HTTPDelete, []byte(inputJson))
		if nil != err {
			blog.Error("delete host batch fail:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostDeleteFail)

		}
		opClient := auditlog.NewClient(cli.CC.AuditCtrl(), req.Request.Header)
		opClient.AuditHostsLog(logConents, "delete host", ownerID, fmt.Sprintf("%d", appID), user, auditoplog.AuditOpTypeDel)

		return http.StatusOK, common.CCSuccessStr, nil
	}, resp)
}
