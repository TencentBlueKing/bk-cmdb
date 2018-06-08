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
 
package api

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/emicklei/go-restful"

	apiutil "configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/httpserver"
	"configcenter/src/common/language"
	frtypes "configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core"
	"configcenter/src/scene_server/topo_server/core/types"
)

// topoAPI the topo server api
type topoAPI struct {
	initFuncs []func()
	core      core.Core
	actions   []action
	err       errors.CCErrorIf
	lang      language.CCLanguageIf
}

func (cli *topoAPI) createAPIRspStr(errcode int, info interface{}) (string, error) {
	rsp := api.BKAPIRsp{
		Result:  true,
		Code:    0,
		Message: nil,
		Data:    nil,
	}

	if common.CCSuccess != errcode {
		rsp.Result = false
		rsp.Code = errcode
		rsp.Message = info
	} else {
		rsp.Message = common.CCSuccessStr
		rsp.Data = info
	}

	s, err := json.Marshal(rsp)

	return string(s), err
}

func (cli *topoAPI) sendResponse(resp *restful.Response, dataMsg interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	if rsp, rspErr := cli.createAPIRspStr(common.CCSuccess, dataMsg); nil == rspErr {
		io.WriteString(resp, rsp)
	}
}

// SetCore set the core instance
func (cli *topoAPI) SetCore(coreMgr core.Core) {

	// set core
	cli.core = coreMgr

	// init
	for _, targetInitFunc := range cli.initFuncs {
		targetInitFunc()
	}

}

// Actions return the all actions
func (cli *topoAPI) Actions() []*httpserver.Action {

	var httpactions []*httpserver.Action
	for _, a := range cli.actions {
		httpactions = append(httpactions, &httpserver.Action{Verb: a.Method, Path: a.Path, Handler: func(req *restful.Request, resp *restful.Response) {

			ownerID := util.GetActionOnwerID(req)
			user := util.GetActionUser(req)

			// get the language
			language := util.GetActionLanguage(req)

			defLang := cli.lang.CreateDefaultCCLanguageIf(language)

			// get the error info by the language
			defErr := cli.err.CreateDefaultCCErrorIf(language)

			value, err := ioutil.ReadAll(req.Request.Body)
			if err != nil {
				blog.Errorf("read http request body failed, error:%s", err.Error())
				errStr := defErr.Error(common.CCErrCommHTTPReadBodyFailed)
				respData, _ := cli.createAPIRspStr(common.CCErrCommHTTPReadBodyFailed, errStr)
				cli.sendResponse(resp, respData)
				return
			}

			mData := frtypes.MapStr{}
			if err := json.Unmarshal(value, &mData); nil != err {
				blog.Errorf("failed to unmarshal the data, error %s", err.Error())
				errStr := defErr.Error(common.CCErrCommJSONUnmarshalFailed)
				respData, _ := cli.createAPIRspStr(common.CCErrCommJSONUnmarshalFailed, errStr)
				cli.sendResponse(resp, respData)
				return
			}

			data, dataErr := a.HandlerFunc(types.LogicParams{
				Err:  defErr,
				Lang: defLang,
				Header: apiutil.Headers{
					Language: language,
					User:     user,
					OwnerID:  ownerID,
				},
			},
				req.PathParameter,
				req.QueryParameter,
				mData)

			if nil != dataErr {
				blog.Errorf("%s", dataErr.Error())
				switch e := dataErr.(type) {
				default:
					respData, _ := cli.createAPIRspStr(common.CCSystemBusy, dataErr.Error())
					cli.sendResponse(resp, respData)
				case errors.CCErrorCoder:
					respData, _ := cli.createAPIRspStr(e.GetCode(), dataErr.Error())
					cli.sendResponse(resp, respData)
				}
				return
			}

			cli.sendResponse(resp, data)

		}})
	}
	return httpactions
}
