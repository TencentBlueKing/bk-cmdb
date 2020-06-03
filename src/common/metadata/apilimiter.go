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

package metadata

import (
	"fmt"
	"regexp"

	"configcenter/src/common/util"
)

// LimiterRule is a rule for api limiter
type LimiterRule struct {
	RuleName string `json:"rulename"`
	AppCode  string `json:"appcode"`
	User     string `json:"user"`
	IP       string `json:"ip"`
	Method   string `json:"method"`
	Url      string `json:"url"`
	Limit    int64  `json:"limit"`
	TTL      int64  `json:"ttl"`
	DenyAll  bool   `json:"denyall"`
}

// Verify to check the fields of LimiterRule
func (r LimiterRule) Verify() error {
	if r.RuleName == "" {
		return fmt.Errorf("rulename must be set")
	}
	if r.AppCode == "" && r.User == "" && r.IP == "" && r.Url == "" && r.Method == "" {
		return fmt.Errorf("one of appcode, user, ip, url, method must be set")
	}
	if r.Method != "" {
		if util.Normalize(r.Method) != "POST" && util.Normalize(r.Method) != "GET" && util.Normalize(r.Method) != "PUT" && util.Normalize(r.Method) != "DELETE" {
			return fmt.Errorf("method must be one of POST,GET,PUT,DELETE")
		}
	}
	if r.Url != "" {
		if _, err := regexp.Compile(r.Url); err != nil {
			return fmt.Errorf("url is not a valid regular expression，%s", err.Error())
		}
	}
	if !r.DenyAll {
		if r.Limit <= 0 || r.TTL <= 0 {
			return fmt.Errorf("both limit and ttl must be set and bigger than 0 when denyall is false")
		}
	}
	return nil
}
