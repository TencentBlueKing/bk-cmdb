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

	"configcenter/src/auth"
	"configcenter/src/common"
	"github.com/emicklei/go-restful"
)

func ParseAttribute(req *restful.Request) (*auth.Attribute, error) {
	attr := new(auth.Attribute)
	user := req.Request.Header.Get(common.BKHTTPHeaderUser)
	if len(user) == 0 {
		return nil, errors.New("miss BK_User in your request header")
	}

	elements, err := urlParse(req.Request.URL.Path)
	if err != nil {
		return nil, err
	}

	version := elements[1]
	if version != "v3" {
		return nil, fmt.Errorf("unsupported api version: %s", version)
	}
	attr.APIVersion = version

	return nil, nil
}

// url example: /api/v3/create/model
var urlRegex = regexp.MustCompile(`^/api/([^/]+)/([^/]+)/([^/]+)/(.*)$`)

func urlParse(url string) (elements []string, err error) {
	if !urlRegex.MatchString(url) {
		return nil, errors.New("invalid url format")
	}

	return strings.Split(url, "/")[1:], nil
}

func filterAction(action string) (auth.Action, error) {
	switch action {
	case "find":
		return auth.Find, nil
	case "findMany":
		return auth.FindMany, nil

	case "create":
		return auth.Create, nil
	case "createMany":
		return auth.CreateMany, nil

	case "update":
		return auth.Update, nil
	case "updateMany":
		return auth.UpdateMany, nil

	case "delete":
		return auth.Delete, nil
	case "deleteMany":
		return auth.DeleteMany, nil

	default:
		return auth.Unknown, fmt.Errorf("unsupported action %s", action)
	}
}
