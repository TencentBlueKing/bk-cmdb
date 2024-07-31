/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package cmdb

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/thirdparty/apigw/apigwutil"

	"github.com/tidwall/gjson"
)

// Client returns cmdb api gateway restful client
func (c *cmdb) Client() rest.ClientInterface {
	return c.service.Client
}

// SetApiGWAuthHeader set authorization header by api gateway config
func (c *cmdb) SetApiGWAuthHeader(header http.Header) http.Header {
	return apigwutil.SetApiGWAuthHeader(c.service.Config, header)
}

// Proxy cmdb api gateway request
func (c *cmdb) Proxy(req *http.Request, rw http.ResponseWriter) {
	resp := make(json.RawMessage, 0)

	body := []byte("")
	if req.Body != nil {
		var err error
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			blog.Errorf("read cmdb api gateway request body failed, err: %v", err)
			rw.Write([]byte(err.Error()))
			return
		}
	}

	result := c.service.Client.Verb(rest.VerbType(req.Method)).
		WithContext(req.Context()).
		Body(body).
		SubResourcef(req.URL.Path).
		WithParamsFromURL(req.URL).
		WithHeaders(c.SetApiGWAuthHeader(req.Header)).
		Do()

	if err := result.Into(&resp); err != nil {
		blog.Errorf("proxy cmdb api gateway request failed, err: %v", err)
		rw.Write([]byte(err.Error()))
		return
	}

	// parse api gateway response format to cmdb inner response format
	if gjson.GetBytes(resp, common.BkAPIErrorCode).Exists() {
		buf := bytes.NewBuffer([]byte{'{'})

		gjson.ParseBytes(resp).ForEach(func(key, value gjson.Result) bool {
			keyStr := key.String()
			switch keyStr {
			case common.BkAPIErrorCode:
				keyStr = common.HTTPBKAPIErrorCode
			case common.BkAPIErrorMessage:
				keyStr = common.HTTPBKAPIErrorMessage
			}

			buf.WriteByte('"')
			buf.WriteString(keyStr)
			buf.WriteString(`":`)
			buf.WriteString(value.Raw)
			buf.WriteByte(',')
			return true
		})

		buf.WriteByte('}')
		resp = buf.Bytes()
	}

	rw.WriteHeader(result.StatusCode)
	util.CopyHeader(result.Header, rw.Header())
	rw.Write(resp)
}
