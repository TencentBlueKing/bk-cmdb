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
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"configcenter/src/apimachinery/util"
	"configcenter/src/common/blog"

	"github.com/emicklei/go-restful/v3"
	"github.com/gin-gonic/gin"
)

// ReqForward TODO
// rest-api http请求代理
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

// ProxyRestHttp TODO
// rest-api请求转发
func ProxyRestHttp(req *restful.Request, resp *restful.Response, addr string) {
	u, err := url.Parse(addr)
	if err == nil {
		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ServeHTTP(resp.ResponseWriter, req.Request)
	} else {
		resp.ResponseWriter.Write([]byte(err.Error()))
	}
}

// ReqHttp TODO
// rest-api请求转发
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

// ProxyHttp TODO
// porxy http
func ProxyHttp(c *gin.Context, addr string) {
	tlsConf, err := util.GetClientTLSConfig("webServer.tls")
	if err != nil {
		c.Writer.Write([]byte(err.Error()))
		blog.Errorf("get webServer.tls config error, err: %v", err)
		return
	}

	u, err := url.Parse(addr)
	if err == nil {
		proxy := httputil.NewSingleHostReverseProxy(u)
		if tlsConf != nil {
			proxy.Transport = &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig:       tlsConf,
			}

		}
		proxy.ServeHTTP(c.Writer, c.Request)
	} else {
		c.Writer.Write([]byte(err.Error()))
	}
}
