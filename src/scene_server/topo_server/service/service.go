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
	"configcenter/src/scene_server/topo_server/app/options"
	"configcenter/src/scene_server/topo_server/core"
	"configcenter/src/scene_server/topo_server/core/types"
	"configcenter/src/storage/dal"

	"github.com/emicklei/go-restful"
)

type Service struct {
	Engine      *backbone.Engine
	Txn         dal.DB
	Core        core.Core
	Config      options.Config
	AuthManager *extensions.AuthManager
	Error       errors.CCErrorIf
	Language    language.CCLanguageIf
	actions     []action
}

// WebService the web service
func (s *Service) WebService() *restful.WebService {

	// init service actions
	s.initService()

	ws := new(restful.WebService)

	getErrFunc := func() errors.CCErrorIf {
		return s.Error
	}
	// TODO: {version} need to replaced by v3
	ws.Path("/topo/{version}").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)

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
		if _, err := io.WriteString(resp, rsp); nil != err {
			blog.Errorf("failed to write string, error info is %s", err.Error())
		}
	} else {
		blog.Errorf("failed to send response , error info is %s", rspErr.Error())
	}
}

func (s *Service) sendCompleteResponse(resp *restful.Response, errorCode int, errMsg string, info interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	rsp, rspErr := s.createCompleteAPIRspStr(errorCode, errMsg, info)
	if nil == rspErr {
		if _, err := io.WriteString(resp, rsp); nil != err {
			blog.Errorf("it is failed to write some data, err:%s", err.Error())
		}
		return
	}
	blog.Errorf("failed to send response , error info is %s", rspErr.Error())

}

func (s *Service) addAction(method string, path string, handlerFunc LogicFunc, handlerParseOriginDataFunc ParseOriginDataFunc) {
	s.addActionEx(method, path, handlerFunc, handlerParseOriginDataFunc, false)
}

func (s *Service) addPublicAction(method string, path string, handlerFunc LogicFunc, handlerParseOriginDataFunc ParseOriginDataFunc) {
	s.addActionEx(method, path, handlerFunc, handlerParseOriginDataFunc, true)
}

func (s *Service) addActionEx(method string, path string, handlerFunc LogicFunc, handlerParseOriginDataFunc ParseOriginDataFunc, publicOnly bool) {
	actionObject := action{
		Method:                     method,
		Path:                       path,
		HandlerFunc:                handlerFunc,
		HandlerParseOriginDataFunc: handlerParseOriginDataFunc,
		PublicOnly:                 publicOnly,
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
				rid := util.GetHTTPCCRequestID(req.Request.Header)

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
				blog.V(9).Infof("request path: %s, body: %s, rid: %s", act.Path, value, rid)

				mData := mapstr.MapStr{}
				if nil == act.HandlerParseOriginDataFunc {
					if err = json.Unmarshal(value, &mData); nil != err && len(value) != 0 {
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

				ctx, _ := s.Engine.CCCtx.WithCancel()
				ctx = context.WithValue(ctx, common.ContextRequestIDField, rid)

				handlerContext := types.ContextParams{
					Context:         ctx,
					Err:             defErr,
					Lang:            defLang,
					MaxTopoLevel:    s.Config.BusinessTopoLevelMax,
					Header:          req.Request.Header,
					SupplierAccount: ownerID,
					User:            user,
					Engin:           s.Engine,
					ReqID:           rid,
				}

				// parse metadata for none public only handler
				if act.PublicOnly == false {
					md := new(MetaShell)
					if len(value) != 0 {
						if err := json.Unmarshal(value, md); err != nil {
							blog.Errorf("parse metadata from request failed, err: %v", err)
							s.sendResponse(resp, common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed))
							return
						}
					}
					handlerContext.MetaData = md.Metadata
				}

				data, dataErr := act.HandlerFunc(handlerContext, req.PathParameter, req.QueryParameter, mData)

				if dataErr == nil {
					s.sendResponse(resp, common.CCSuccess, data)
					return
				}

				switch e := dataErr.(type) {
				case errors.CCErrorCoder:
					s.sendCompleteResponse(resp, e.GetCode(), dataErr.Error(), data)
				default:
					s.sendCompleteResponse(resp, common.CCSystemBusy, dataErr.Error(), data)
				}
				return
			}})
		}(a)

	}
	return httpactions
}

type MetaShell struct {
	Metadata *metadata.Metadata `json:"metadata"`
}
