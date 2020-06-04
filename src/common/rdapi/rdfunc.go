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
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func checkHTTPAuth(req *restful.Request, defErr errors.DefaultCCErrorIf) (int, string) {
	util.SetOwnerIDAndAccount(req)
	if "" == util.GetOwnerID(req.Request.Header) {
		return common.CCErrCommNotAuthItem, defErr.Errorf(common.CCErrCommNotAuthItem, "owner_id").Error()
	}
	if "" == util.GetUser(req.Request.Header) {
		return common.CCErrCommNotAuthItem, defErr.Errorf(common.CCErrCommNotAuthItem, "user").Error()
	}

	return common.CCSuccess, ""

}

func AllGlobalFilter(errFunc func() errors.CCErrorIf) func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		defer func() {
			if fetalErr := recover(); fetalErr != nil {
				rid := util.GetHTTPCCRequestID(req.Request.Header)
				blog.Errorf("server panic, err:%#v, rid:%s, debug strace:%s", fetalErr, rid, debug.Stack())
				ccErrTip := errFunc().CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header)).Errorf(common.CCErrCommInternalServerError, common.GetIdentification())
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
		language := util.GetLanguage(req.Request.Header)
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

func RequestLogFilter() func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		header := req.Request.Header
		body, _ := util.PeekRequest(req.Request)
		blog.Infof("code: %s, user: %s, rip: %s, uri: %s, body: %s, rid: %s",
			header.Get("Bk-App-Code"), header.Get("Bk_user"), header.Get("X-Real-Ip"),
			req.Request.RequestURI, body, util.GetHTTPCCRequestID(header))

		fchain.ProcessFilter(req, resp)
		return
	}
}

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

func GenerateHttpHeaderRID(req *http.Request, resp http.ResponseWriter) {
	cid := util.GetHTTPCCRequestID(req.Header)
	if "" == cid {
		cid = GetHTTPOtherRequestID(req.Header)
		if cid == "" {
			cid = util.GenerateRID()
		}
		req.Header.Set(common.BKHTTPCCRequestID, cid)
		resp.Header().Set(common.BKHTTPCCRequestID, cid)
	}

	return
}

func ServiceErrorHandler(err restful.ServiceError, req *restful.Request, resp *restful.Response) {
	blog.Errorf("HTTP ERROR: %v, HTTP MESSAGE: %v, RequestURI: %s %s", err.Code, err.Message, req.Request.Method, req.Request.RequestURI)
	ret := metadata.BaseResp{
		Result: false,
		Code:   -1,
		ErrMsg: fmt.Sprintf("HTTP ERROR: %v, HTTP MESSAGE: %v, RequestURI: %s %s", err.Code, err.Message, req.Request.Method, req.Request.RequestURI),
	}

	resp.WriteHeaderAndJson(err.Code, ret, "application/json")
}

// getHTTPOtherRequestID return other system request id from http header
func GetHTTPOtherRequestID(header http.Header) string {
	return header.Get(common.BKHTTPOtherRequestID)
}
