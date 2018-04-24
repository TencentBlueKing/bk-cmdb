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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"github.com/emicklei/go-restful"
	//"github.com/tidwall/gjson"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func (cli *instAction) InstSearch(req *restful.Request, resp *restful.Response) {
	blog.Info("search all insts by object")

	// get language
	language := util.GetActionLanguage(req)

	// get error info by language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// logics
	cli.CallResponseEx(func() (int, interface{}, error) {
		ownerID := req.PathParameter("owner_id")
		objID := req.PathParameter("obj_id")

		searchParams := make(map[string]interface{})

		value, readErr := ioutil.ReadAll(req.Request.Body)

		if nil != readErr {
			blog.Error("failed to read the body , error info is %s", readErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		if 0 != len(value) {

			var js params.SearchParams
			err := json.Unmarshal([]byte(value), &js)
			if nil != err {
				blog.Error("failed to unmarshal the data[%s], error is %s", value, err.Error())
				return http.StatusBadRequest, "", defErr.Error(common.CCErrCommJSONUnmarshalFailed)
			}

			condition := params.ParseAppSearchParams(js.Condition)

			condition[common.BKOwnerIDField] = ownerID
			condition[common.BKObjIDField] = objID

			page := js.Page

			searchParams["condition"] = condition
			searchParams["fields"] = strings.Join(js.Fields, ",")
			searchParams["start"] = page["start"]
			searchParams["limit"] = page["limit"]
			searchParams["sort"] = page["sort"]

		} else {
			condition := make(map[string]interface{}, 0)
			condition[common.BKOwnerIDField] = map[string]interface{}{
				common.BKDBIN: []string{"", ownerID},
			}
			condition[common.BKObjIDField] = objID
			searchParams["condition"] = condition
			searchParams["fields"] = ""
			searchParams["start"] = 0
			searchParams["limit"] = common.BKDefaultLimit
			searchParams["sort"] = ""

		}

		//search
		innerObject := getInnerNameByObjectID(objID)
		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + innerObject + "/search"
		inputJSON, jsErr := json.Marshal(searchParams)

		if nil != jsErr {
			blog.Error("failed to marshal the data[%+v], error info is %s", searchParams, jsErr.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrCommJSONMarshalFailed)
		}

		objRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, inputJSON)
		blog.Debug("search url(%s) inst params: %s", sURL, string(inputJSON))

		if nil != err {
			blog.Error("failed to select the insts, error info is %s", err.Error())
			return http.StatusInternalServerError, "", defErr.Error(common.CCErrTopoInstSelectFailed)
		}

		retStr, retStrErr := cli.getInstDetails(req, objID, ownerID, objRes, map[string]interface{}{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  "",
		})

		if common.CCSuccess != retStrErr {
			return http.StatusInternalServerError, "", defErr.Error(retStrErr)
		}

		return http.StatusOK, retStr["data"], nil
	}, resp)
}
