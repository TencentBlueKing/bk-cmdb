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

package update

import (
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	scenecommon "configcenter/src/scene_server/common"
	"configcenter/src/scene_server/host_server/host_service/logics"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/auditlog"
	"configcenter/src/source_controller/api/metadata"
	sourceAPI "configcenter/src/source_controller/api/object"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	"configcenter/src/common/util"
	"fmt"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"

	"github.com/emicklei/go-restful"
)

var host *hostAction = &hostAction{}

type hostAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/host/batch", Params: nil, Handler: host.UpdateHostBatch})
	// create CC object
	host.CreateAction()
}

//UpdateHostBatch update host batch
func (cli *hostAction) UpdateHostBatch(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		user := util.GetActionUser(req)
		//update host
		hostCond := make(map[string]interface{})
		condInput := make(map[string]interface{})
		input := make(map[string]interface{})
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, err := simplejson.NewJson([]byte(value))
		data, _ := js.Map()
		hostIDStr, ok := data[common.BKHostIDField].(string)
		if false == ok {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostFeildValidFail)

		}
		delete(data, common.BKHostIDField)
		forward := &sourceAPI.ForwardParam{Header: req.Request.Header}
		valid := validator.NewValidMap(common.BKDefaultOwnerID, common.BKInnerObjIDHost, cli.CC.ObjCtrl(), forward, defErr)

		hostIDArr := strings.Split(hostIDStr, ",")
		var iHostIDArr []int
		logPreConents := make(map[int]auditoplog.AuditLogExt, 0)
		hostFields, _ := logics.GetHostLogFields(req, common.BKDefaultOwnerID, cli.CC.ObjCtrl())
		for _, i := range hostIDArr {

			iHostID, _ := strconv.Atoi(i)
			//加日志
			logObj := logics.NewHostLog(req, common.BKDefaultOwnerID, i, cli.CC.HostCtrl(), cli.CC.ObjCtrl(), hostFields)

			//validate
			_, err = valid.ValidMap(data, common.ValidUpdate, iHostID)
			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostFeildValidFail)

			}
			iHostIDArr = append(iHostIDArr, iHostID)

			logContent := logObj.GetPreHostData()
			logPreConents[iHostID] = auditoplog.AuditLogExt{ID: iHostID, Content: logContent, ExtKey: logObj.GetInnerIP()}

		}
		hostModuleConfig, _ := logics.GetConfigByCond(req, cli.CC.HostCtrl(), map[string]interface{}{common.BKHostIDField: iHostIDArr})

		hostCond[common.BKDBIN] = iHostIDArr
		condInput[common.BKHostIDField] = hostCond
		input["condition"] = condInput
		input["data"] = data

		uHostURL := cli.CC.ObjCtrl() + "/object/v1/insts/host"
		inputJson, _ := json.Marshal(input)
		blog.Info("update host batch url:%s", uHostURL)
		blog.Info("update host batch content:%s", string(inputJson))
		_, err = httpcli.ReqHttp(req, uHostURL, common.HTTPUpdate, []byte(inputJson))
		if nil != err {
			blog.Error("update host batch fail:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostUpdateFail)
		}
		appID := "0"
		if len(hostModuleConfig) > 0 {
			appID = fmt.Sprintf("%v", hostModuleConfig[0][common.BKAppIDField])

		}

		var logLastConents []auditoplog.AuditLogExt
		for _, i := range iHostIDArr {

			// set the inst association table
			if err := scenecommon.UpdateInstAssociation(cli.CC.ObjCtrl(), req, i, common.BKDefaultOwnerID, common.BKInnerObjIDHost, data); nil != err {
				blog.Errorf("failed to update the inst association, error info is %s ", err.Error())
			}
			//get change value
			logObj := logics.NewHostLog(req, common.BKDefaultOwnerID, fmt.Sprintf("%d", i), cli.CC.HostCtrl(), cli.CC.ObjCtrl(), hostFields)
			logContent := logObj.GetPreHostData()
			preLogContent, ok := logPreConents[i]
			//set change value to curdata
			logContent.CurData = logContent.PreData
			if ok {
				content, _ := preLogContent.Content.(*metadata.Content)
				logContent.PreData = content.PreData
			}
			logLastConents = append(logLastConents, auditoplog.AuditLogExt{ID: i, Content: logContent, ExtKey: preLogContent.ExtKey})

		}
		opClient := auditlog.NewClient(cli.CC.AuditCtrl(), req.Request.Header)
		opClient.AuditHostsLog(logLastConents, "update host", common.BKDefaultOwnerID, appID, user, auditoplog.AuditOpTypeModify)

		return http.StatusOK, common.CCSuccessStr, nil
	}, resp)

}
