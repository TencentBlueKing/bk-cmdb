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

package instdata

import (
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/rs/xid"
)

var history = &historyAction{}

// ObjectAction
type historyAction struct {
	base.BaseAction
}

//AddHistory add history
func (cli *historyAction) AddHistory(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		value, _ := ioutil.ReadAll(req.Request.Body)
		bodyData := make(map[string]interface{})
		err := json.Unmarshal([]byte(value), &bodyData)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)

		}
		content, _ := bodyData["content"].(string)

		data := make(map[string]interface{}, 4)
		data["content"] = content

		if "" == data["content"] {
			blog.Error("param content could not be empty")
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, "content")
		}
		data["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
		data[common.CreateTimeField] = time.Now()
		id := xid.New()
		data["id"] = id.String()

		data = util.SetModOwner(data, ownerID)
		_, err = history.CC.InstCli.Insert("cc_History", data)
		if nil != err {
			blog.Error("Create  history fail, error information is %s, params:%v", err.Error(), data)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBInsertFailed)
		}
		data = make(map[string]interface{}, 1)
		data["id"] = id

		return http.StatusOK, data, nil
	}, resp)
}

//GetHistorys get historys
func (cli *historyAction) GetHistorys(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		start, _ := strconv.Atoi(req.PathParameter("start"))
		limit, _ := strconv.Atoi(req.PathParameter("limit"))

		conds := make(map[string]interface{}, 1)
		conds["user"] = req.PathParameter("user") //libraries.GetOperateUser(req)
		fields := []string{"id", "content", common.CreateTimeField, "user"}
		//GetMutilByCondition(cName string, fields []string, s interface{}, result interface{}, sort string, skip, limit int) error
		var result []interface{}
		sort := "-" + common.LastTimeField
		conds = util.SetModOwner(conds, ownerID)
		err := history.CC.InstCli.GetMutilByCondition("cc_History", fields, conds, &result, sort, start, limit)
		if nil != err {
			blog.Error("query  history fail, error information is %s, params:%v", err.Error(), conds)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBSelectFailed)
		}

		nums, err := history.CC.InstCli.GetCntByCondition("cc_History", conds)
		if nil != err {
			blog.Error("query  history fail, error information is %s, params:%v", err.Error(), conds)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBInsertFailed)
		}
		data := make(map[string]interface{}, 2)
		data["count"] = nums
		data["info"] = result

		return http.StatusOK, data, nil
	}, resp)
}

func init() {
	history.CC = api.NewAPIResource()
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/history/{user}", Params: nil, Handler: history.AddHistory})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/history/{user}/{start}/{limit}", Params: nil, Handler: history.GetHistorys})

}
