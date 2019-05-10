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

	restful "github.com/emicklei/go-restful"
	"github.com/rs/xid"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

var (
	srvNames map[string]int = make(map[string]int) // servername mapping
)

type addrSrv interface {
	GetServer(servType string) (string, error)
}

func GetRdAddrSrvHandle(typeSrv string, addrSrv api.AddrSrv) func() string {
	srvNames[typeSrv] = 1
	return func() string {
		url, err := addrSrv.GetServer(typeSrv)
		blog.V(3).Infof("GetRdAddrSrvHandle  get %s url:%s", typeSrv, url)
		if nil != err {
			blog.Errorf("get %s addr from service discovery module error: %s", typeSrv, err.Error())
			return ""
		}
		if "" == url {
			blog.Errorf("get %s addr from service discovery module,no available service found", typeSrv)
			return ""
		}
		return url
	}

}

func FilterRdAddrSrv(typeSrv string) func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	srvNames[typeSrv] = 1
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		cli := api.NewAPIResource()

		url, err := cli.AddrSrv.GetServer(typeSrv)
		language := util.GetActionLanguage(req)
		if nil == cli.Error {
			rsp, _ := createAPIRspStr(common.CCErrCommConfMissItem, "config file is missing err item")
			io.WriteString(resp, rsp)
			return
		}
		blog.V(3).Infof("FilterRdAddrSrv %s url:%s", typeSrv, url)

		defErr := cli.Error.CreateDefaultCCErrorIf(language)
		if nil != err {
			blog.Errorf("get %s addr from service discovery module error: %s", typeSrv, err.Error())
			resp.WriteHeader(http.StatusInternalServerError)
			rsp, rsperr := createAPIRspStr(common.CCErrCommRelyOnServerAddressFailed, defErr.Errorf(common.CCErrCommRelyOnServerAddressFailed, typeSrv).Error())
			if nil != rsperr {
				blog.Errorf("create response failed, error information is %v", rsperr)
			} else {
				// TODO: 暂时不设置 resp.WriteHeader(httpcode)
				io.WriteString(resp, rsp)
			}
			return

		} else if "" == url {
			blog.Errorf("get %s addr from service discovery module,no available service found", typeSrv)
			resp.WriteHeader(http.StatusInternalServerError)
			rsp, rsperr := createAPIRspStr(common.CCErrCommRelyOnServerAddressFailed, defErr.Errorf(common.CCErrCommRelyOnServerAddressFailed, typeSrv).Error())
			if nil != rsperr {
				blog.Errorf("create response failed, error information is %v", rsperr)
			} else {
				// TODO: 暂时不设置 resp.WriteHeader(httpcode)
				io.WriteString(resp, rsp)
			}
			return
		}
		fchain.ProcessFilter(req, resp)
		return
	}
}

func FilterRdAddrSrvs(typeSrvs ...string) func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	for _, typeSrv := range typeSrvs {
		srvNames[typeSrv] = 1
	}

	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		cli := api.NewAPIResource()
		language := util.GetActionLanguage(req)
		if nil == cli.Error {
			rsp, _ := createAPIRspStr(common.CCErrCommConfMissItem, "config file is missing err item")
			io.WriteString(resp, rsp)
			return
		}
		defErr := cli.Error.CreateDefaultCCErrorIf(language)

		for _, typeSrv := range typeSrvs {
			url, err := cli.AddrSrv.GetServer(typeSrv)
			blog.V(3).Infof("FilterRdAddrSrv %s url:%s", typeSrv, url)
			if nil != err {
				blog.Errorf("get %s addr from service discovery module error: %s", typeSrv, err.Error())
				resp.WriteHeader(http.StatusInternalServerError)
				rsp, rsperr := createAPIRspStr(common.CCErrCommRelyOnServerAddressFailed, defErr.Errorf(common.CCErrCommRelyOnServerAddressFailed, typeSrv).Error())
				if nil != rsperr {
					blog.Errorf("create response failed, error information is %v", rsperr)
				} else {
					// TODO: 暂时不设置 resp.WriteHeader(httpcode)
					io.WriteString(resp, rsp)
				}
				return

			} else if "" == url {
				blog.Errorf("get %s addr from service discovery module,no available service found", typeSrv)
				resp.WriteHeader(http.StatusInternalServerError)
				rsp, rsperr := createAPIRspStr(common.CCErrCommRelyOnServerAddressFailed, defErr.Errorf(common.CCErrCommRelyOnServerAddressFailed, typeSrv).Error())
				if nil != rsperr {
					blog.Errorf("create response failed, error information is %v", rsperr)
				} else {
					// TODO: 暂时不设置 resp.WriteHeader(httpcode)
					io.WriteString(resp, rsp)
				}
				return
			}
		}
		fchain.ProcessFilter(req, resp)

		return
	}

}

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
		defer func() {
			if fetalErr := recover(); fetalErr != nil {
				rid := util.GetHTTPCCRequestID(req.Request.Header)
				blog.Errorf("server panic, err:%#v, rid:%s, debug strace:%s", fetalErr, rid, debug.Stack())
				// 3.2.x version. identity is not enabled. tempporarily returned to null
				ccErrTip := errFunc().CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header)).Errorf(common.CCErrCommInternalServerError, "")
				respErrInfo := &metadata.RespError{Msg: ccErrTip}
				io.WriteString(resp, respErrInfo.Error())
			}

		}()
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

func GlobalFilter(typeSrvs ...string) func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		cli := api.NewAPIResource()
		language := util.GetActionLanguage(req)
		if nil == cli.Error {
			rsp, _ := createAPIRspStr(common.CCErrCommConfMissItem, "config file is missing err item")
			io.WriteString(resp, rsp)
			return
		}
		defErr := cli.Error.CreateDefaultCCErrorIf(language)

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

		for _, typeSrv := range typeSrvs {
			url, err := cli.AddrSrv.GetServer(typeSrv)
			blog.V(3).Infof("GlobalFilter %s url:%s", typeSrv, url)
			if nil != err {
				blog.Errorf("get %s addr from service discovery module error: %s", typeSrv, err.Error())
				resp.WriteHeader(http.StatusInternalServerError)
				rsp, rsperr := createAPIRspStr(common.CCErrCommRelyOnServerAddressFailed, defErr.Errorf(common.CCErrCommRelyOnServerAddressFailed, typeSrv).Error())
				if nil != rsperr {
					blog.Errorf("create response failed, error information is %v", rsperr)
				} else {
					// TODO: 暂时不设置 resp.WriteHeader(httpcode)
					io.WriteString(resp, rsp)
				}
				return

			} else if "" == url {
				blog.Errorf("get %s addr from service discovery module,no available service found", typeSrv)
				resp.WriteHeader(http.StatusInternalServerError)
				rsp, rsperr := createAPIRspStr(common.CCErrCommRelyOnServerAddressFailed, defErr.Errorf(common.CCErrCommRelyOnServerAddressFailed, typeSrv).Error())
				if nil != rsperr {
					blog.Errorf("create response failed, error information is %v", rsperr)
				} else {
					// TODO: 暂时不设置 resp.WriteHeader(httpcode)
					io.WriteString(resp, rsp)
				}
				return
			}
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
	unused := "0000"
	cid := util.GetHTTPCCRequestID(req.Request.Header)
	if "" == cid {
		cid = generateRID(unused)
		req.Request.Header.Set(common.BKHTTPCCRequestID, cid)
	}
	// todo support esb request id

	resp.Header().Set(common.BKHTTPCCRequestID, cid)
}

func generateRID(unused string) string {
	id := xid.New()
	return fmt.Sprintf("cc%s%s", unused, id.String())
}
