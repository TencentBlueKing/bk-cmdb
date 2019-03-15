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
	"regexp"
	"strings"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func ParseAttribute(req *restful.Request) (*meta.AuthAttribute, error) {
	body, err := util.PeekRequest(req.Request)
	if err != nil {
		return nil, err
	}

	meta := struct {
		Metadata metadata.Metadata `json:"metadata"`
	}{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &meta); err != nil {
			return nil, err
		}
	}

	elements, err := urlParse(req.Request.URL.Path)
	if err != nil {
		return nil, err
	}

	requestContext := &RequestContext{
		Header:   req.Request.Header,
		Method:   req.Request.Method,
		URI:      req.Request.URL.Path,
		Elements: elements,
		Body:     body,
		Metadata: meta.Metadata,
	}

	stream, err := newParseStream(requestContext)
	if err != nil {
		return nil, err
	}

	return stream.Parse()
}

// ParseCommonInfo get common info from req, aims at avoiding too much repeat code
func ParseCommonInfo(req *restful.Request) (*meta.CommonInfo, error) {
	commonInfo := new(meta.CommonInfo)

	userInfo, err := ParseUserInfo(req)
	if err != nil {
		return nil, err
	}
	commonInfo.User = *userInfo

	apiVersion, err := ParseAPIVersion(req)
	if err != nil {
		return nil, err
	}
	commonInfo.APIVersion = apiVersion

	return commonInfo, nil
}

func ParseUserInfo(req *restful.Request) (*meta.UserInfo, error) {
	userInfo := new(meta.UserInfo)
	user := req.Request.Header.Get(common.BKHTTPHeaderUser)
	if len(user) == 0 {
		return nil, errors.New("miss BK_User in your request header")
	}
	userInfo.UserName = user
	supplierID := req.Request.Header.Get(common.BKHTTPSupplierID)
	if len(supplierID) == 0 {
		return nil, errors.New("miss bk_supplier_id in your request header")
	}
	userInfo.SupplierAccount = supplierID
	return userInfo, nil
}

func ParseAPIVersion(req *restful.Request) (string, error) {
	elements, err := urlParse(req.Request.URL.Path)
	if err != nil {
		return "", err
	}
	version := elements[1]
	if version != "v3" {
		return "", fmt.Errorf("unsupported api version: %s", version)
	}
	return version, nil
}

// url example: /api/v3/create/model
var urlRegex = regexp.MustCompile(`^/api/([^/]+)(/[^/]+)+/?$`)

func urlParse(url string) (elements []string, err error) {
	if !urlRegex.MatchString(url) {
		return nil, errors.New("invalid url format")
	}

	return strings.Split(url, "/")[1:], nil
}

func filterAction(action string) (meta.Action, error) {
	switch action {
	case "find":
		return meta.Find, nil
	case "findMany":
		return meta.FindMany, nil

	case "create":
		return meta.Create, nil
	case "createMany":
		return meta.CreateMany, nil

	case "update":
		return meta.Update, nil
	case "updateMany":
		return meta.UpdateMany, nil

	case "delete":
		return meta.Delete, nil
	case "deleteMany":
		return meta.DeleteMany, nil

	default:
		return meta.Unknown, fmt.Errorf("unsupported action %s", action)
	}
}
