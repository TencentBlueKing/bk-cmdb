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

package httpclient

import (
	"io/ioutil"
	"net/http/httputil"

	"net/url"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin"
)

//rest-api http请求代理
func ReqForward(req *restful.Request, url, method string) (string, error) {
	// blog.Infof("forward %s with header %v", url, req.Request.Header)
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return err.Error(), err
	}
	httpcli := NewHttpClient()
	for key := range req.Request.Header {
		httpcli.SetHeader(key, req.Request.Header.Get(key))
	}
	httpcli.SetHeader("Content-Type", "application/json")
	httpcli.SetHeader("Accept", "application/json")

	reply, err := httpcli.Request(url, method, req.Request.Header, body)
	if err != nil {
		return err.Error(), err
	}

	return string(reply), err
}

//rest-api请求转发
func ProxyRestHttp(req *restful.Request, resp *restful.Response, addr string) {
	u, err := url.Parse(addr)
	if err == nil {
		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ServeHTTP(resp.ResponseWriter, req.Request)
	} else {
		resp.ResponseWriter.Write([]byte(err.Error()))
	}
}

//rest-api请求转发
func ReqHttp(req *restful.Request, url, method string, body []byte) (string, error) {
	// blog.Infof("forward %s with header %v", url, req.Request.Header)
	httpcli := NewHttpClient()
	httpcli.SetHeader("Content-Type", "application/json")
	httpcli.SetHeader("Accept", "application/json")
	for key := range req.Request.Header {
		httpcli.SetHeader(key, req.Request.Header.Get(key))
	}
	reply, err := httpcli.Request(url, method, req.Request.Header, body)
	if err != nil {
		return err.Error(), err
	}

	return string(reply), err
}

//porxy http
func ProxyHttp(c *gin.Context, addr string) {

	u, err := url.Parse(addr)
	if err == nil {
		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ServeHTTP(c.Writer, c.Request)
	} else {
		c.Writer.Write([]byte(err.Error()))
	}
}
