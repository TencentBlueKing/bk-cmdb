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
	"time"

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
	"configcenter/src/source_controller/coreservice/app/options"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/source_controller/coreservice/core/association"
	"configcenter/src/source_controller/coreservice/core/datasynchronize"
	"configcenter/src/source_controller/coreservice/core/instances"
	"configcenter/src/source_controller/coreservice/core/model"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/remote"

	"github.com/emicklei/go-restful"
)

// CoreServiceInterface the topo service methods used to init
type CoreServiceInterface interface {
	WebService() *restful.WebService
	SetConfig(cfg options.Config, engin *backbone.Engine, err errors.CCErrorIf, language language.CCLanguageIf) error
}

// New ceate topo servcie instance
func New() CoreServiceInterface {
	return &coreService{}
}

// coreService topo service
type coreService struct {
	engin    *backbone.Engine
	language language.CCLanguageIf
	err      errors.CCErrorIf
	actions  []action
	cfg      options.Config
	core     core.Core
}

func (s *coreService) SetConfig(cfg options.Config, engin *backbone.Engine, err errors.CCErrorIf, language language.CCLanguageIf) error {

	s.cfg = cfg
	s.engin = engin

	if nil != err {
		s.err = err
	}

	if nil != language {
		s.language = language
	}

	var db dal.DB
	var dbErr error
	if cfg.Mongo.Transaction == "enable" {
		blog.Infof("connecting to transaction manager")
		db, dbErr = remote.NewWithDiscover(engin.ServiceManageInterface.TMServer().GetServers, cfg.Mongo)
		if dbErr != nil {
			blog.Errorf("failed to connect the txc server, error info is %s", dbErr.Error())
			return dbErr
		}
	} else {
		db, dbErr = local.NewMgo(cfg.Mongo.BuildURI(), time.Minute)
		if dbErr != nil {
			blog.Errorf("failed to connect the remote server(%s), error info is %s", cfg.Mongo.BuildURI(), dbErr.Error())
			return dbErr
		}
	}
	// connect the remote mongodb

	s.core = core.New(model.New(db, s), instances.New(db, s), association.New(db, s), datasynchronize.New(db, s))
	return nil
}

// WebService the web service
func (s *coreService) WebService() *restful.WebService {

	// init service actions
	s.initService()

	ws := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.err
	}
	ws.Path("/api/v3").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)

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

func (s *coreService) createAPIRspStr(errcode int, info interface{}) (string, error) {

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

func (s *coreService) createCompleteAPIRspStr(errcode int, errmsg string, info interface{}) (string, error) {

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

func (s *coreService) sendResponse(resp *restful.Response, errorCode int, dataMsg interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	if rsp, rspErr := s.createAPIRspStr(errorCode, dataMsg); nil == rspErr {
		io.WriteString(resp, rsp)
	} else {
		blog.Errorf("failed to send response , error info is %s", rspErr.Error())
	}
}

func (s *coreService) sendCompleteResponse(resp *restful.Response, errorCode int, errMsg string, info interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	rsp, rspErr := s.createCompleteAPIRspStr(errorCode, errMsg, info)
	if nil == rspErr {
		io.WriteString(resp, rsp)
		return
	}
	blog.Errorf("failed to send response , error info is %s", rspErr.Error())

}

// Actions return the all actions
func (s *coreService) Actions() []*httpserver.Action {

	var httpactions []*httpserver.Action
	for _, a := range s.actions {

		func(act action) {

			httpactions = append(httpactions, &httpserver.Action{Verb: act.Method, Path: act.Path, Handler: func(req *restful.Request, resp *restful.Response) {

				ownerID := util.GetActionOnwerID(req)
				user := util.GetActionUser(req)

				// get the language
				language := util.GetActionLanguage(req)

				defLang := s.language.CreateDefaultCCLanguageIf(language)

				// get the error info by the language
				defErr := s.err.CreateDefaultCCErrorIf(language)

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
