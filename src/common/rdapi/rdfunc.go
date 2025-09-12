/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package rdapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful/v3"
)

func checkHTTPAuth(req *restful.Request, defErr errors.DefaultCCErrorIf) (int, string) {
	if httpheader.GetTenantID(req.Request.Header) == "" {
		return common.CCErrCommNotAuthItem, defErr.Errorf(common.CCErrCommNotAuthItem, "tenant_id").Error()
	}
	if httpheader.GetUser(req.Request.Header) == "" {
		return common.CCErrCommNotAuthItem, defErr.Errorf(common.CCErrCommNotAuthItem, "user").Error()
	}

	return common.CCSuccess, ""

}

// AllGlobalFilter TODO
func AllGlobalFilter(errFunc func() errors.CCErrorIf) func(req *restful.Request, resp *restful.Response,
	fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		defer func() {
			if fetalErr := recover(); fetalErr != nil {
				rid := httpheader.GetRid(req.Request.Header)
				blog.Errorf("server panic, err: %v, rid: %s, debug strace: %s", fetalErr, rid, debug.Stack())
				ccErrTip := errFunc().CreateDefaultCCErrorIf(httpheader.GetLanguage(req.Request.Header)).
					Errorf(common.CCErrCommInternalServerError, common.GetIdentification())
				respErrInfo := &metadata.RespError{Msg: ccErrTip}
				io.WriteString(resp, respErrInfo.Error())
			}

		}()

		GenerateHttpHeaderRID(req.Request, resp.ResponseWriter)

		whiteListSuffix := strings.Split(common.URLFilterWhiteListSuffix, common.URLFilterWhiteListSepareteChar)
		for _, url := range whiteListSuffix {
			if strings.HasSuffix(req.Request.URL.Path, url) {
				fchain.ProcessFilter(req, resp)
				return
			}
		}
		language := httpheader.GetLanguage(req.Request.Header)
		defErr := errFunc().CreateDefaultCCErrorIf(language)

		errNO, errMsg := checkHTTPAuth(req, defErr)

		if common.CCSuccess != errNO {
			resp.WriteHeader(http.StatusInternalServerError)
			rsp, _ := createAPIRspStr(errNO, errMsg)
			io.WriteString(resp, rsp)
			return
		}

		fchain.ProcessFilter(req, resp)
		return
	}
}

// RequestLogFilter TODO
func RequestLogFilter() func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		header := req.Request.Header
		body, _ := util.PeekRequest(req.Request)
		blog.Infof("code: %s, user: %s, rip: %s, uri: %s, body: %s, rid: %s",
			httpheader.GetAppCode(header), httpheader.GetUser(header), httpheader.GetReqRealIP(header),
			req.Request.RequestURI, util.FormatHttpBody(req.Request.URL.Path, body), httpheader.GetRid(header))

		fchain.ProcessFilter(req, resp)
		return
	}
}

// HTTPRequestIDFilter TODO
func HTTPRequestIDFilter() func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		GenerateHttpHeaderRID(req.Request, resp.ResponseWriter)
		if 1 < len(fchain.Filters) {
			fchain.ProcessFilter(req, resp)
			return
		}

		fchain.ProcessFilter(req, resp)
		return
	}
}

func createAPIRspStr(errcode int, info string) (string, error) {

	var rsp metadata.Response

	if 0 != errcode {
		rsp.Result = false
		rsp.Code = errcode
		rsp.ErrMsg = info
	} else {
		rsp.Data = info
	}

	s, err := json.Marshal(rsp)

	return string(s), err
}

// GenerateHttpHeaderRID generate http header request id
func GenerateHttpHeaderRID(req *http.Request, resp http.ResponseWriter) {
	rid := httpheader.GetRid(req.Header)
	if rid == "" {
		rid = util.GenerateRID()
		httpheader.SetRid(req.Header, rid)
	}
	httpheader.SetRid(resp.Header(), rid)
}

// ServiceErrorHandler TODO
func ServiceErrorHandler(err restful.ServiceError, req *restful.Request, resp *restful.Response) {
	blog.Errorf("HTTP ERROR: %v, HTTP MESSAGE: %v, RequestURI: %s %s", err.Code, err.Message, req.Request.Method,
		req.Request.RequestURI)
	ret := metadata.BaseResp{
		Result: false,
		Code:   -1,
		ErrMsg: fmt.Sprintf("HTTP ERROR: %v, HTTP MESSAGE: %v, RequestURI: %s %s", err.Code, err.Message,
			req.Request.Method, req.Request.RequestURI),
	}

	resp.WriteHeaderAndJson(err.Code, ret, "application/json")
}
