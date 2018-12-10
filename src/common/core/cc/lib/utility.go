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
 
package lib

import (
	"configcenter/src/common"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/blog"
	chttp "configcenter/src/common/http"
	"configcenter/src/common/http/httpclient"
	"io/ioutil"

	"github.com/emicklei/go-restful"
)

// request2sence proxy the requst for sence
func request2sence(req *restful.Request, host, uri, method string) (string, error) {
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		blog.Error("read http request body failed. err: %s", err.Error())
		err = chttp.InternalError(common.CC_Err_Comm_http_ReadReqBody, common.CC_Err_Comm_http_ReadReqBody_STR+err.Error())
		return err.Error(), err
	}

	a := api.GetAPIResource()

	url := host + "/ccapi/v1/" + uri //a.Conf.CCHost
	blog.Debug("do request to url(%s), method(%s), request:%s", url, method, string(body))

	httpcli := httpclient.NewHttpClient()
	httpcli.SetHeader("Content-Type", "application/json")
	httpcli.SetHeader("Accept", "application/json")
	if a.IsClientSSL() {
		httpcli.SetTlsVerityConfig(a.GetClientSSL())
	}

	reply, err := httpcli.Request(url, method, req.Request.Header, body)
	if err != nil {
		blog.Error("http request failed. err: %s", err.Error())
		err = chttp.InternalError(common.CC_Err_Comm_http_DO, common.CC_Err_Comm_http_DO_STR+err.Error())
		return err.Error(), err
	}
	blog.Debug("respone from url(%s), rsp: %s", url, string(reply))
	return string(reply), err
}
