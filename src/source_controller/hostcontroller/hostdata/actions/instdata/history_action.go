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
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	. "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	restful "github.com/emicklei/go-restful"
	"github.com/rs/xid"
)

var history *historyAction = &historyAction{}

// ObjectAction
type historyAction struct {
	base.BaseAction
}

//AddHistory add history
func (cli *historyAction) AddHistory(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	bodyData := new(HistoryContent)
	if err := json.NewDecoder(req.Request.Body).Decode(bodyData); err != nil {
		blog.Errorf("add history, but decode body failed, err: %v", err)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommHTTPReadBodyFailed).Error()})
		return
	}

	if bodyData.Content == "" {
		blog.Errorf("add history, but history content is empty.")
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommParamsNeedSet).Error()})
		return
	}

	data := make(map[string]interface{}, 4)
	data["content"] = bodyData.Content
	data["user"] = req.PathParameter("user")
	data[common.CreateTimeField] = time.Now()
	id := xid.New()
	data["id"] = id.String()

	_, err := history.CC.InstCli.Insert("cc_History", data)
	if nil != err {
		blog.Error("add history failed, err: %v, params:%v", err, data)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommDBInsertFailed).Error()})
		return
	}

	resp.WriteAsJson(HostFavorite{
		BaseResp: BaseResp{true, http.StatusOK, ""},
		Data:     ID{ID: id.String()},
	})

}

func (cli *historyAction) GetHistorys(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	start, err := strconv.Atoi(req.PathParameter("start"))
	if err != nil {
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommParamsIsInvalid).Error()})
		return
	}
	limit, err := strconv.Atoi(req.PathParameter("limit"))
	if err != nil {
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommParamsIsInvalid).Error()})
		return
	}

	conds := make(map[string]interface{}, 1)
	conds["user"] = req.PathParameter("user")
	fields := []string{"id", "content", common.CreateTimeField, "user"}
	var result []interface{}
	sort := "-" + common.LastTimeField
	err = history.CC.InstCli.GetMutilByCondition("cc_History", fields, conds, &result, sort, start, limit)
	if nil != err {
		blog.Error("query  history failed, err: %v, params: %v", err, conds)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommDBSelectFailed).Error()})
		return
	}

	nums, err := history.CC.InstCli.GetCntByCondition("cc_History", conds)
	if nil != err {
		blog.Error("query  history failed, err: %v, params:%v", err, conds)
		resp.WriteAsJson(BaseResp{Code: http.StatusBadRequest, ErrMsg: defErr.Error(common.CCErrCommDBInsertFailed).Error()})
		return
	}

	resp.WriteAsJson(GetHistoryResult{
		BaseResp: BaseResp{true, http.StatusOK, ""},
		Data: HistoryResult{
			Count: nums,
			Info:  result,
		},
	})
}

func init() {
	history.CC = api.NewAPIResource()
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/history/{user}", Params: nil, Handler: history.AddHistory})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/history/{user}/{start}/{limit}", Params: nil, Handler: history.GetHistorys})

}
