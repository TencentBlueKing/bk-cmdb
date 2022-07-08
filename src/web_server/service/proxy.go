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
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	apiutil "configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"

	"github.com/gin-gonic/gin"
)

// ProxyRequest to proxy third-party api request to solve cross domain issue.
func (s *Service) ProxyRequest(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	webCommon.SetProxyHeader(c)
	language := webCommon.GetLanguageByHTTPRequest(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	method := c.Param("method")
	target := c.Param("target")
	target_url := c.Param("target_url")

	if len(method) == 0 || len(target) == 0 || len(target_url) == 0 {
		blog.Errorf("path parameter must be filled in, rid: %s", rid)
		reply := getReturnStr(common.CCErrCommParamsIsInvalid,
			defErr.CCErrorf(common.CCErrCommParamsIsInvalid, "method/target/target_url").Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}

	reqDomain, err := s.getTargetDomainUrl(target)
	if err != nil {
		blog.Errorf("get target domain url failed, err: %v, rid: %s", err, rid)
		reply := getReturnStr(common.CCErrCommParamsIsInvalid, err.Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}

	url, err := url.Parse(fmt.Sprintf("%s%s", reqDomain, target_url))
	if err != nil {
		blog.Errorf("parse url failed, err: %v, rid: %s", err, rid)
		_, _ = c.Writer.Write([]byte(err.Error()))
		return
	}

	tlsConf, err := apiutil.GetClientTLSConfig("webServer.site.paas.tls")
	if err != nil {
		c.Writer.Write([]byte(err.Error()))
		blog.Errorf("get webServer.site.paas.tls config error, err: %v", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
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

	reqMethod, err := getRequestMethod(method)
	if err != nil {
		blog.Errorf("get true request method failed, err: %v, rid: %s", err, rid)
		reply := getReturnStr(common.CCErrCommParamsIsInvalid, err.Error(), nil)
		_, _ = c.Writer.Write([]byte(reply))
		return
	}

	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.URL.Path = target_url
		req.Method = reqMethod
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func (s *Service) getTargetDomainUrl(target string) (string, error) {
	switch target {
	case "usermanage":
		return s.Config.Site.PaasDomainUrl, nil
	default:
		return "", fmt.Errorf("did not support this target saas: %s", target)
	}
}

func getRequestMethod(urlMethod string) (string, error) {
	method := strings.ToUpper(urlMethod)
	switch method {
	case "GET", "POST", "PUT", "DELETE", "CONNECT", "HEAD", "OPTIONS", "TRACE":
		return method, nil
	default:
		return "", fmt.Errorf("did not support this proxy request method: %s", urlMethod)
	}
}
