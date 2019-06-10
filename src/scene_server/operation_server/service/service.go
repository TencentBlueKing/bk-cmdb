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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/httpserver"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/operation_server/app/options"
	"configcenter/src/scene_server/operation_server/core"

	"github.com/emicklei/go-restful"
)

type Service struct {
	Engine      *backbone.Engine
	Config      options.Config
	Core        *core.Operation
	AuthManager *extensions.AuthManager
	Error       errors.CCErrorIf
	Language    language.CCLanguageIf
	actions     []action
}

type ServiceInterface interface {
	WebService() *restful.WebService
	SetConfig(cfg options.Config, engin *backbone.Engine, err errors.CCErrorIf, language language.CCLanguageIf) error
}

func New() ServiceInterface {
	return &Service{}
}

func (s *Service) SetConfig(cfg options.Config, engin *backbone.Engine, err errors.CCErrorIf, language language.CCLanguageIf) error {

	s.Config = cfg
	s.Engine = engin

	if nil != err {
		s.Error = err
	}

	if nil != language {
		s.Language = language
	}

<<<<<<< HEAD
	return nil
}
=======
	o.newOperationService(api)
	container := restful.NewContainer()
	container.Add(api)
>>>>>>> c7685d399... fix: operation crud bugs

// WebService the web service
func (s *Service) WebService() *restful.WebService {

	// init service actions
	s.initService()

<<<<<<< HEAD
	ws := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.Error
	}
	ws.Path("/operation/v3").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)

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
=======
func (o *OperationServer) newOperationService(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  o.Engine.CCErr,
		Language: o.Engine.Language,
	})

	// service category
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/operation/chart", Handler: o.CreateStatisticChart})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/operation/chart", Handler: o.DeleteStatisticChart})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/update/operation/chart", Handler: o.UpdateStatisticChart})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/search/operation/chart", Handler: o.SearchStatisticChart})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/search/operation/chart/data", Handler: o.SearchChartData})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/update/operation/chart/position", Handler: o.UpdateChartPosition})

	utility.AddToRestfulWebService(web)
>>>>>>> c7685d399... fix: operation crud bugs
}

func (s *Service) createAPIRspStr(errcode int, info interface{}) (string, error) {

	rsp := metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     nil,
	}

	if common.CCSuccess != errcode {
		rsp.Code = errcode
		rsp.Result = false
		rsp.ErrMsg = fmt.Sprintf("%v", info)
	} else {
		rsp.ErrMsg = common.CCSuccessStr
		rsp.Data = info
	}

	data, err := json.Marshal(rsp)
	return string(data), err
}

func (s *Service) createCompleteAPIRspStr(errcode int, errmsg string, info interface{}) (string, error) {

	rsp := metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     nil,
	}

	if common.CCSuccess != errcode {
		rsp.Code = errcode
		rsp.Result = false
		rsp.ErrMsg = errmsg
	} else {
		rsp.ErrMsg = common.CCSuccessStr
	}
	rsp.Data = info
	data, err := json.Marshal(rsp)
	return string(data), err
}

func (s *Service) sendResponse(resp *restful.Response, errorCode int, dataMsg interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	if rsp, rspErr := s.createAPIRspStr(errorCode, dataMsg); nil == rspErr {
		io.WriteString(resp, rsp)
	} else {
		blog.Errorf("failed to send response , error info is %s", rspErr.Error())
	}
}

func (s *Service) sendCompleteResponse(resp *restful.Response, errorCode int, errMsg string, info interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	rsp, rspErr := s.createCompleteAPIRspStr(errorCode, errMsg, info)
	if nil == rspErr {
		io.WriteString(resp, rsp)
		return
	}
	blog.Errorf("failed to send response , error info is %s", rspErr.Error())

}

func (s *Service) addAction(method string, path string, handlerFunc LogicFunc, handlerParseOriginDataFunc ParseOriginDataFunc) {
	actionObject := action{
		Method:                     method,
		Path:                       path,
		HandlerFunc:                handlerFunc,
		HandlerParseOriginDataFunc: handlerParseOriginDataFunc,
	}
	s.actions = append(s.actions, actionObject)
}

// Actions return the all actions
func (s *Service) Actions() []*httpserver.Action {

	var httpactions []*httpserver.Action
	for _, a := range s.actions {

		func(act action) {

			httpactions = append(httpactions, &httpserver.Action{Verb: act.Method, Path: act.Path, Handler: func(req *restful.Request, resp *restful.Response) {

				ownerID := util.GetOwnerID(req.Request.Header)
				user := util.GetUser(req.Request.Header)

				// get the language
				language := util.GetLanguage(req.Request.Header)

				defLang := s.Language.CreateDefaultCCLanguageIf(language)

				// get the error info by the language
				defErr := s.Error.CreateDefaultCCErrorIf(language)

				value, err := ioutil.ReadAll(req.Request.Body)
				if err != nil {
					blog.Errorf("read http request body failed, error:%s", err.Error())
					errStr := defErr.Error(common.CCErrCommHTTPReadBodyFailed)
					s.sendResponse(resp, common.CCErrCommHTTPReadBodyFailed, errStr)
					return
				}

				mData := mapstr.MapStr{}
				if nil == act.HandlerParseOriginDataFunc {
					if err := json.Unmarshal(value, &mData); nil != err && 0 != len(value) {
						blog.Errorf("failed to unmarshal the data, error %s", err.Error())
						errStr := defErr.Error(common.CCErrCommJSONUnmarshalFailed)
						s.sendResponse(resp, common.CCErrCommJSONUnmarshalFailed, errStr)
						return
					}
				} else {
					mData, err = act.HandlerParseOriginDataFunc(value)
					if nil != err {
						blog.Errorf("failed to unmarshal the data, error %s", err.Error())
						errStr := defErr.Error(common.CCErrCommJSONUnmarshalFailed)
						s.sendResponse(resp, common.CCErrCommJSONUnmarshalFailed, errStr)
						return
					}
				}

				data, dataErr := act.HandlerFunc(core.ContextParams{
					Context:         util.GetDBContext(context.Background(), req.Request.Header),
					Error:           defErr,
					Lang:            defLang,
					Header:          req.Request.Header,
					SupplierAccount: ownerID,
					ReqID:           util.GetHTTPCCRequestID(req.Request.Header),
					User:            user,
				},
					req.PathParameter,
					req.QueryParameter,
					mData)

				if nil != dataErr {
					switch e := dataErr.(type) {
					default:
						s.sendCompleteResponse(resp, common.CCSystemBusy, dataErr.Error(), data)
					case errors.CCErrorCoder:
						s.sendCompleteResponse(resp, e.GetCode(), dataErr.Error(), data)
					}
					return
				}

				s.sendResponse(resp, common.CCSuccess, data)

			}})
		}(a)

	}
	return httpactions
}
