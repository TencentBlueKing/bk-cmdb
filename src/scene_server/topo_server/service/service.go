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

package service

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"

	apiutil "configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/httpserver"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/logics"
	"configcenter/src/scene_server/topo_server/app/options"
	"configcenter/src/scene_server/topo_server/core"
	"configcenter/src/scene_server/topo_server/core/types"
)

// Service topo service
type Service struct {
	*options.Config
	*backbone.Engine
	*logics.Logics

	actions []action
	core    core.Core
}

// WebService the web service
func (s *Service) WebService(filter restful.FilterFunction) *restful.WebService {

	// set core
	s.core = core.New(s.CoreAPI)

	// init service actions
	s.initService()

	ws := new(restful.WebService)
	ws.Path("/topo/v3").Filter(filter).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)

	innerActions := s.Actions()

	for _, actionItem := range innerActions {
		switch actionItem.Verb {
		case http.MethodPost:
			ws.Route(ws.POST(actionItem.Path).To(actionItem.Handler))
		case http.MethodDelete:
			ws.Route(ws.DELETE(actionItem.Path).To(actionItem.Handler))
		case http.MethodPut:
			ws.Route(ws.PUT(actionItem.Path).To(actionItem.Handler))
		case http.MethodGet:
			ws.Route(ws.GET(actionItem.Path).To(actionItem.Handler))
		default:
			blog.Errorf(" the url (%s), the http method (%s) is not supported", actionItem.Path, actionItem.Verb)
		}
	}

	return ws
}

func (s *Service) createAPIRspStr(errcode int, info interface{}) (string, error) {
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

	data, err := json.Marshal(rsp)
	return string(data), err
}

func (s *Service) sendResponse(resp *restful.Response, dataMsg interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	if rsp, rspErr := s.createAPIRspStr(common.CCSuccess, dataMsg); nil == rspErr {
		io.WriteString(resp, rsp)
	}
}

// Actions return the all actions
func (s *Service) Actions() []*httpserver.Action {

	var httpactions []*httpserver.Action
	for _, a := range s.actions {

		func(act action) {

			httpactions = append(httpactions, &httpserver.Action{Verb: act.Method, Path: act.Path, Handler: func(req *restful.Request, resp *restful.Response) {

				ownerID := util.GetActionOnwerID(req)
				user := util.GetActionUser(req)

				// get the language
				language := util.GetActionLanguage(req)

				defLang := s.Language.CreateDefaultCCLanguageIf(language)

				// get the error info by the language
				defErr := s.CCErr.CreateDefaultCCErrorIf(language)

				value, err := ioutil.ReadAll(req.Request.Body)
				if err != nil {
					blog.Errorf("read http request body failed, error:%s", err.Error())
					errStr := defErr.Error(common.CCErrCommHTTPReadBodyFailed)
					respData, _ := s.createAPIRspStr(common.CCErrCommHTTPReadBodyFailed, errStr)
					s.sendResponse(resp, respData)
					return
				}

				mData := frtypes.MapStr{}
				if err := json.Unmarshal(value, &mData); nil != err {
					blog.Errorf("failed to unmarshal the data, error %s", err.Error())
					errStr := defErr.Error(common.CCErrCommJSONUnmarshalFailed)
					respData, _ := s.createAPIRspStr(common.CCErrCommJSONUnmarshalFailed, errStr)
					s.sendResponse(resp, respData)
					return
				}

				data, dataErr := act.HandlerFunc(types.LogicParams{
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
						respData, _ := s.createAPIRspStr(common.CCSystemBusy, dataErr.Error())
						s.sendResponse(resp, respData)
					case errors.CCErrorCoder:
						respData, _ := s.createAPIRspStr(e.GetCode(), dataErr.Error())
						s.sendResponse(resp, respData)
					}
					return
				}

				s.sendResponse(resp, data)

			}})
		}(a)

	}
	return httpactions
}
