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
	"net/http"
	"net/http/httputil"
	"net/url"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"

	"github.com/gin-gonic/gin"
)

type queryProxyResp struct {
	metadata.BaseResp `json:",inline"`
	Data              interface{} `json:"data" mapstructure:"data"`
}

// ProxyRequest to proxy third-party api request
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
		c.JSON(http.StatusOK, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommParamsIsInvalid,
			ErrMsg: defErr.CCErrorf(common.CCErrCommParamsIsInvalid, "method/target/target_url").Error(),
		})
		return
	}

	url, err := url.Parse(fmt.Sprintf("%s%s", s.getTargetDomainUrl(target), target_url))
	if err != nil {
		blog.Errorf("parse url failed, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusOK, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommParamsIsInvalid,
			ErrMsg: err.Error(),
		})
		return
	}

	director := func(req *http.Request) {
		req.URL = url
	}
	proxy := &httputil.ReverseProxy{Director: director}
	c.Request.Method = "POST"
	proxy.ServeHTTP(c.Writer, c.Request)
	return
}

func (s *Service) getTargetDomainUrl(target string) string {
	switch target {
	case "usermanage":
		return s.Config.Site.PaasDomainUrl
	default:
		return ""
	}
}
