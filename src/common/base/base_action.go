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

package base

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/errors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/emicklei/go-restful"
)

// BaseAction 所有Action 的基类
type BaseAction struct {
	CC *api.APIResource
}

// CreateAction 执行Action 初始化
func (cli *BaseAction) CreateAction() error {
	cli.CC = api.NewAPIResource()
	return nil
}

// IsSuccess check the response
func (cli *BaseAction) IsSuccess(rst []byte) (*api.APIRsp, bool) {

	var rstRes api.APIRsp
	if jserr := json.Unmarshal(rst, &rstRes); nil != jserr {
		blog.Error("can not unmarshal the result , error: %s", jserr.Error())
		return &rstRes, false
	}

	if rstRes.Code != common.CCSuccess {
		return &rstRes, false
	}

	return &rstRes, true

}

// CallResponseEx  execute and deal with the response
func (cli *BaseAction) CallResponseEx(callfunc func() (httpcode int, reply interface{}, err error), resp *restful.Response) {

	httpcode, reply, err := callfunc()

	if nil == err {
		switch r := reply.(type) {
		case nil:
			cli.ResponseSuccessData(common.CCSuccessStr, resp)
		case string:
			cli.ResponseRspString(resp, r, common.CCSuccess, "")
		default:
			cli.ResponseSuccessData(reply, resp)
		}

	} else {
		switch e := err.(type) {
		case nil:
			blog.Error("the error is nil")
		case errors.CCErrorCoder:
			cli.ResponseFailedEx(httpcode, e.GetCode(), e.Error(), resp)
		}
	}
}

// CallResponse execute and deal with the response
func (cli *BaseAction) CallResponse(callfunc func() (string, error), resp *restful.Response) {

	reply, err := callfunc()

	if nil == err {
		cli.ResponseRspString(resp, reply, common.CCSuccess, "")
	} else {
		blog.Error("request failed, error:%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_DO, err.Error(), resp)
	}

}

// ResponseRspString  parse the api.APIRsp struct
func (cli *BaseAction) ResponseRspString(resp *restful.Response, data string, errno int, err string) {

	if errno != common.CCSuccess {

		blog.Error("request object att information failed, error information is %v", err)
		cli.ResponseFailed(errno, err, resp)

	} else {

		var rst api.APIRsp
		jserr := json.Unmarshal([]byte(data), &rst)
		if nil == jserr {
			cli.Response(&rst, resp)
			return
		}

		blog.Warnf("unmarshal the json failed, error:%s, data: %v", jserr.Error(), data)
		cli.ResponseSuccess(data, resp)
	}
}

// Response execute the response
func (cli *BaseAction) Response(rst *api.APIRsp, resp *restful.Response) {

	if rst.Result {
		cli.ResponseSuccess(rst.Data, resp)
	} else {
		cli.ResponseFailed(rst.Code, rst.Message, resp)
	}
}

// ResponseFailedEx deal with the http  code and response
func (cli *BaseAction) ResponseFailedEx(httpcode int, errno int, errmsg interface{}, resp *restful.Response) {

	rsp, rsperr := cli.CC.CreateAPIRspStr(errno, errmsg)
	if nil != rsperr {
		blog.Error("create response failed, error information is %v", rsperr)
	} else {
		// TODO: 暂时不设置 resp.WriteHeader(httpcode)
		io.WriteString(resp, rsp)
	}

}

// ResponseFailed deal with the http  code and response
func (cli *BaseAction) ResponseFailed(errno int, errmsg interface{}, resp *restful.Response) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusBadRequest)
	rsp, rsperr := cli.CC.CreateAPIRspStr(errno, errmsg)
	if nil != rsperr {
		blog.Error("create response failed, error information is %v", rsperr)
	} else {
		io.WriteString(resp, rsp)
	}

}

// ResponseFailedWithData deal with the failed response
func (cli *BaseAction) ResponseFailedWithData(errno int, errmsg, data interface{}, resp *restful.Response) {

	var rst api.APIRsp
	rst.Code = errno
	rst.Result = false
	rst.Message = errmsg
	rst.Data = data
	rsp, jserr := json.Marshal(rst)
	if nil == jserr {
		io.WriteString(resp, string(rsp))
		return

	}

	blog.Error("unmarshal the json failed, error:%s, data: %v", jserr.Error(), data)
	cli.ResponseFailed(common.CC_Err_Comm_http_DO, jserr.Error(), resp)

}

// ResponseSuccess deal with the success response
func (cli *BaseAction) ResponseSuccess(datamsg interface{}, resp *restful.Response) {
	resp.Header().Set("Content-Type", "application/json")
	if rsp, rsperr := cli.CC.CreateAPIRspStr(common.CCSuccess, datamsg); nil == rsperr {
		io.WriteString(resp, rsp)
	} else {
		blog.Error("fail to create response for add object att, error information is %v", rsperr)
	}

}

// ResponseSuccessData only response the success data
func (cli *BaseAction) ResponseSuccessData(datamsg interface{}, resp *restful.Response) bool {

	if rsp, rsperr := cli.CC.CreateAPIRspStr(common.CCSuccess, datamsg); nil == rsperr {
		io.WriteString(resp, rsp)
	} else {
		blog.Error("fail to create response for add object att, error information is %v", rsperr)
		return false
	}

	return true

}

// ResponseNative return Native data
func (cli *BaseAction) ResponseNative(datamsg string, resp *restful.Response) {
	io.WriteString(resp, datamsg)
}

// GetParams parse the params
func (cli *BaseAction) GetParams(cc *api.APIResource, keyvalues *map[string]string, key string, result interface{}, resp *restful.Response) error {

	val, ok := (*keyvalues)[key]
	if !ok {
		err := fmt.Errorf("the key word '%s' was lost", key)
		if rsp, rsperr := cc.CreateAPIRspStr(common.CC_Err_Comm_http_DO, err.Error()); nil != rsperr {
			blog.Error("create a response failed, error:%s", rsperr.Error())
		} else {
			io.WriteString(resp, rsp)
			return err
		}
	}

	objtype := reflect.TypeOf(result)
	objtype = objtype.Elem()
	switch objtype.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32:
		tmpval, tmperr := strconv.ParseInt(val, 10, 64)
		if nil != tmperr {
			if rsp, rsperr := cc.CreateAPIRspStr(common.CC_Err_Comm_http_DO, tmperr.Error()); nil != rsperr {
				blog.Error("create a response failed, error: %s", rsperr.Error())
			} else {
				io.WriteString(resp, rsp)
				return tmperr
			}
		}
		resulttmp := reflect.ValueOf(result)
		resulttmp.Elem().SetInt(tmpval)
	case reflect.String:

		resulttmp := reflect.ValueOf(result)
		if rawval, rawerr := url.QueryUnescape(val); nil != rawerr {
			blog.Error("url decode failed, error:%s", rawerr.Error())
			if rsp, rsperr := cc.CreateAPIRspStr(common.CC_Err_Comm_http_DO, rawerr.Error()); nil != rsperr {
				blog.Error("create a response failed, error:%s", rsperr.Error())
			} else {
				io.WriteString(resp, rsp)
				return rawerr
			}
		} else {
			resulttmp.Elem().SetString(rawval)
		}

	default:
		err := fmt.Errorf("the key word '%s' data type was not supported", key)
		if rsp, rsperr := cc.CreateAPIRspStr(common.CC_Err_Comm_http_DO, err.Error()); nil != rsperr {
			blog.Error("create a response failed, error:%s", rsperr.Error())
		} else {
			io.WriteString(resp, rsp)
			return err
		}
	}

	return nil
}
