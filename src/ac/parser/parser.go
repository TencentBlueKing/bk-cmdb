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

package parser

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func ParseAttribute(req *restful.Request, engine *backbone.Engine) (*meta.AuthAttribute, error) {
	elements, err := urlParse(req.Request.URL.Path)
	if err != nil {
		return nil, err
	}

	requestContext := &RequestContext{
		Rid:      util.GetHTTPCCRequestID(req.Request.Header),
		Header:   req.Request.Header,
		Method:   req.Request.Method,
		URI:      req.Request.URL.Path,
		Elements: elements,
		getBody: func() (body []byte, err error) {
			body, err = util.PeekRequest(req.Request)
			if err != nil {
				return nil, err
			}
			return
		},
	}

	stream, err := newParseStream(requestContext, engine)
	if err != nil {
		return nil, err
	}

	return stream.Parse()
}

// ParseCommonInfo get common info from req, aims at avoiding too much repeat code
func ParseCommonInfo(requestHeader *http.Header) (*meta.CommonInfo, error) {
	commonInfo := new(meta.CommonInfo)

	userInfo, err := ParseUserInfo(requestHeader)
	if err != nil {
		return nil, err
	}
	commonInfo.User = *userInfo

	return commonInfo, nil
}

func ParseUserInfo(requestHeader *http.Header) (*meta.UserInfo, error) {
	userInfo := new(meta.UserInfo)
	user := requestHeader.Get(common.BKHTTPHeaderUser)
	if len(user) == 0 {
		return nil, errors.New("parse user info failed, miss BK_User in your request header")
	}
	userInfo.UserName = user
	supplierID := requestHeader.Get(common.BKHTTPOwnerID)
	if len(supplierID) == 0 {
		return nil, errors.New("parse user info failed, miss bk_supplier_id in your request header")
	}
	userInfo.SupplierAccount = supplierID
	return userInfo, nil
}

// url example: /api/v3/create/model
var urlRegex = regexp.MustCompile(`^/api/([^/]+)(/[^/]+)+/?$`)

func urlParse(url string) (elements []string, err error) {
	if !urlRegex.MatchString(url) {
		return nil, fmt.Errorf("invalid url format, url=%s", url)
	}

	return strings.Split(url, "/")[1:], nil
}
