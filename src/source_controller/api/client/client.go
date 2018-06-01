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

package client

import (
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client
type Client struct {
	Base    base.BaseLogic
	address string
}

// SetAddress
func (cli *Client) SetAddress(address string) {
	cli.address = address
}

func (cli *Client) GetAddress() string {
	return cli.address
}

func (cli *Client) GetRequestInfo(method string, data interface{}, url string) (interface{}, error) {
	byteData, _ := json.Marshal(data)
	var rst []byte
	var err error
	switch method {
	case http.MethodPost: //HTTPSelectPost, HTTPCreate
		rst, err = cli.Base.HttpCli.POST(url, nil, byteData)
	case common.HTTPSelectGet:
		rst, err = cli.Base.HttpCli.GET(url, nil, byteData)
	case common.HTTPUpdate:
		rst, err = cli.Base.HttpCli.PUT(url, nil, byteData)
	case common.HTTPDelete:
		rst, err = cli.Base.HttpCli.DELETE(url, nil, byteData)

	}
	blog.Debug("url:%s, input:%s", url, string(byteData))
	if nil != err {
		blog.Error("request failed, error:%v, log centent:%s", err, string(byteData))
		return "", Err_Request
	}

	var retObj common.APIRsp
	if jserr := json.Unmarshal(rst, &retObj); nil != jserr {

		blog.Error("can not unmarshal the result , error information is %v, log content:%s", jserr, string(byteData))
		return "", jserr
	}

	if retObj.Code != common.CCSuccess {
		return "", fmt.Errorf("%v", retObj.Message)
	}
	retData := retObj.Data

	var bkRetOjb api.BKAPIRsp
	if jserr := json.Unmarshal(rst, &bkRetOjb); nil != jserr {

		blog.Error("can not unmarshal the result , error information is %v, log content:%s", jserr, string(byteData))
		return "", jserr
	}

	if bkRetOjb.Code != common.CCSuccess {
		return "", fmt.Errorf("%v", bkRetOjb.Message)
	}
	if nil != retData {
		retData = bkRetOjb.Data
	}
	return retData, nil
}

//GetRequestInfoEx 获取http调用返回值， httpcode， 内容，error
func (cli *Client) GetRequestInfoEx(method string, data interface{}, url string) (int, *common.APIRsp, error) {
	byteData, _ := json.Marshal(data)
	var rst []byte
	var err error
	var code int
	switch method {
	case http.MethodPost: //HTTPSelectPost, HTTPCreate
		code, rst, err = cli.Base.HttpCli.POSTEx(url, nil, byteData)
	case common.HTTPSelectGet:
		code, rst, err = cli.Base.HttpCli.GETEx(url, nil, byteData)
	case common.HTTPUpdate:
		code, rst, err = cli.Base.HttpCli.PUTEx(url, nil, byteData)
	case common.HTTPDelete:
		code, rst, err = cli.Base.HttpCli.DELETEEx(url, nil, byteData)

	}
	blog.Debug(url)
	if nil != err {
		blog.Error("request failed, error:%v, log centent:%s", err, string(byteData))
		return 0, nil, Err_Request
	}
	var retObj common.APIRsp
	if jserr := json.Unmarshal(rst, &retObj); nil != jserr {

		blog.Error("can not unmarshal the result , error information is %v, log content:%s", jserr, string(byteData))
		return http.StatusBadGateway, nil, jserr
	}

	if retObj.Code != common.CCSuccess {
		return code, nil, fmt.Errorf("%v", retObj.Message)
	}

	var bkRetOjb api.BKAPIRsp
	if jserr := json.Unmarshal(rst, &bkRetOjb); nil != jserr {

		blog.Error("can not unmarshal the result , error information is %v, log content:%s", jserr, string(byteData))
		return http.StatusBadGateway, nil, jserr
	}

	if bkRetOjb.Code != common.CCSuccess {
		return bkRetOjb.Code, nil, fmt.Errorf("%v", bkRetOjb.Message)
	}
	if nil != bkRetOjb.Data {
		retObj.Data = bkRetOjb.Data
	}

	return code, &retObj, nil
}
