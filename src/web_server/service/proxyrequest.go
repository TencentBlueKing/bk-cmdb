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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	webcom "configcenter/src/web_server/common"

	"github.com/gin-gonic/gin"
)

type queryProxy struct {
	Url  string                 `json:"url"`
	Args map[string]interface{} `json:"args"`
}

type queryProxyResp struct {
	metadata.BaseResp `json:",inline"`
	Data              interface{} `json:"data" mapstructure:"data"`
}

// ProxyRequest to proxy third-party api request
func (s *Service) ProxyRequest(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	webcom.SetProxyHeader(c)

	query := new(queryProxy)
	if err := c.BindJSON(query); err != nil {
		blog.Errorf("unmarshal failed, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusOK, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommJSONUnmarshalFailed,
			ErrMsg: err.Error(),
		})
		return
	}

	if len(query.Url) == 0 {
		c.JSON(http.StatusOK, metadata.BaseResp{Result: true})
		return
	}

	args := make([]string, 0)
	for k, v := range query.Args {
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	token, err := c.Cookie("bk_token")
	if err != nil {
		blog.Errorf("get bk_token failed, token: %s, err: %v, rid: %s", token, err, rid)
		c.JSON(http.StatusOK, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommParamsIsInvalid,
			ErrMsg: err.Error(),
		})
		return
	}
	rsp, err := http.Get(fmt.Sprintf("%s?%s&%s", query.Url, token, strings.Join(args, "&")))
	if err != nil {
		return
	}

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		blog.Errorf("read response body failed, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusOK, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommHTTPReadBodyFailed,
			ErrMsg: err.Error(),
		})
		return
	}

	result := new(queryProxyResp)
	if err := json.Unmarshal([]byte(string(body)), result); err != nil {
		blog.Errorf("unmarshal failed, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusOK, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommJSONUnmarshalFailed,
			ErrMsg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
	return
}
