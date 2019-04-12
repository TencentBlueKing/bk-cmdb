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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	restful "github.com/emicklei/go-restful"
)

func checkHTTPAuth(req *restful.Request, defErr errors.DefaultCCErrorIf) (int, string) {
	util.SetActionOwerIDAndAccount(req)
	ownerId, user := util.GetActionOnwerIDAndUser(req)
	if "" == ownerId {
		return common.CCErrCommNotAuthItem, defErr.Errorf(common.CCErrCommNotAuthItem, "owner_id").Error()
	}
	if "" == user {
		return common.CCErrCommNotAuthItem, defErr.Errorf(common.CCErrCommNotAuthItem, "user").Error()
	}

	return common.CCSuccess, ""

}

func AllGlobalFilter(errFunc func() errors.CCErrorIf) func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		generateHttpHeaderRID(req, resp)

		whilteListSuffix := strings.Split(common.URLFilterWhiteListSuffix, common.URLFilterWhiteListSepareteChar)
		for _, url := range whilteListSuffix {
			if strings.HasSuffix(req.Request.URL.Path, url) {
				fchain.ProcessFilter(req, resp)
				return
			}
		}
		language := util.GetActionLanguage(req)
		defErr := errFunc().CreateDefaultCCErrorIf(language)

		errNO, errMsg := checkHTTPAuth(req, defErr)

		if common.CCSuccess != errNO {
			resp.WriteHeader(http.StatusInternalServerError)
			rsp, _ := createAPIRspStr(errNO, errMsg)
			io.WriteString(resp, rsp)
			return
		}

		if 1 < len(fchain.Filters) {
			fchain.ProcessFilter(req, resp)
			return
		}

		fchain.ProcessFilter(req, resp)
		return
	}
}

func HTTPRequestIDFilter(errFunc func() errors.CCErrorIf) func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		generateHttpHeaderRID(req, resp)
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

func generateHttpHeaderRID(req *restful.Request, resp *restful.Response) {
	cid := util.GetHTTPCCRequestID(req.Request.Header)
	if "" == cid {
		cid = getHTTPOtherRequestID(req.Request.Header)
		if cid == "" {
			cid = util.GenerateRID()
		}
		req.Request.Header.Set(common.BKHTTPCCRequestID, cid)
	}
	// todo support esb request id

	resp.Header().Set(common.BKHTTPCCRequestID, cid)
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
func getHTTPOtherRequestID(header http.Header) string {
	rid := header.Get(common.BKHTTPOtherRequestID)
	return rid
}
