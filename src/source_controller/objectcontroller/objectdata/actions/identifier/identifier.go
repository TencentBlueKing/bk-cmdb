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
	"io/ioutil"
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/common/commondata"
	"configcenter/src/source_controller/common/instdata"

	"github.com/emicklei/go-restful"
)

var obj = &identifierAction{}

// identifierAction
type identifierAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/identifier/{obj_type}/search", Params: nil, Handler: obj.SearchIdentifier})
	// set cc api interface
	obj.CreateAction()
}

//search object
func (cli *identifierAction) SearchIdentifier(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		pathParams := req.PathParameters()
		objType := pathParams["obj_type"]
		instdata.DataH = cli.CC.InstCli

		value, err := ioutil.ReadAll(req.Request.Body)
		var dat commondata.ObjQueryInput
		err = json.Unmarshal([]byte(value), &dat)
		if err != nil {
			blog.Error("get object type:%s,input:%v error:%v", string(objType), value, err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)

		}
		//dat.ConvTime()
		fields := dat.Fields
		condition := dat.Condition

		skip := dat.Start
		limit := dat.Limit
		sort := dat.Sort
		fieldArr := strings.Split(fields, ",")
		result := make([]map[string]interface{}, 0)
		count, err := instdata.GetCntByCondition(objType, condition)
		if err != nil {
			blog.Error("get object type:%s,input:%v error:%v", objType, string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectSelectInstFailed)

		}
		err = instdata.GetObjectByCondition(defLang, objType, fieldArr, condition, &result, sort, skip, limit)
		if err != nil {
			blog.Error("get object type:%s,input:%v error:%v", string(objType), string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectSelectInstFailed)
		}
		info := make(map[string]interface{})
		info["count"] = count
		info["info"] = result
		return http.StatusOK, info, nil
	}, resp)
}
