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

package tool

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/types"
)

// service singleton global object
var service *Service

// lock ensure concurrency safety
var lock sync.Mutex

// Service dynamically adjust the structure of the runtime configuration
type Service struct {
	limiter *Limiter
	log     *LogService
}

// GetService get a service struct
func GetService() *Service {
	lock.Lock()
	defer lock.Unlock()
	if service == nil {
		service = &Service{
			limiter: NewLimiter(),
			log:     NewLogService(),
		}
	}
	return service
}

// GetLimiter get limiter
func GetLimiter() *Limiter {
	return GetService().limiter
}

// ServeHTTP method of dynamically adjusting runtime configuration
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	op := r.FormValue(common.SettingsOperation)

	switch op {
	case common.SettingsLimiter:
		s.settingsLimiter(w, r)

	case common.SettingsLog:
		loglevel := r.FormValue(common.SettingsLogLevel)
		if err := s.log.ChangeLogLevel(loglevel); err != nil {
			blog.Errorf("update log level failed err: %v", err)
			fmt.Fprintf(w, "update log level failed err: %v", err)
			return
		}
		fmt.Fprintln(w, "success!")

	default:
		blog.Errorf("adjust operation error, can't find the relevant operation to adjust.")
		fmt.Fprintln(w, "adjust operation error, can't find the relevant operation to adjust!")
	}
}

func (s *Service) settingsLimiter(w http.ResponseWriter, r *http.Request) {
	if common.GetIdentification() != types.CC_MODULE_APISERVER {
		blog.Errorf("adjust limiter operation error, can't find the relevant operation action function " +
			"to adjust.")
		fmt.Fprintln(w, "adjust limiter operation error, can't find the relevant operation action function "+
			"to adjust.")
		return
	}
	action := r.FormValue(common.SettingsAction)
	switch action {
	case common.SettingsGetAction:
		limiterRuleNames := new(meta.LimiterRuleNames)
		if err := json.NewDecoder(r.Body).Decode(&limiterRuleNames); nil != err {
			blog.Errorf("get limiter rule failed with decode body err: %v", err)
			fmt.Fprintf(w, "get limiter rule failed with decode body err: %v", err)
			return
		}
		rules, err := s.limiter.GetRules(limiterRuleNames.RuleNames)
		if err != nil {
			blog.Errorf("get limiter rule failed err: %v", err)
			fmt.Fprintf(w, "get limiter rule failed err: %v", err)
			return
		}
		data, err := json.Marshal(rules)
		if err != nil {
			blog.Errorf("get limiter rule failed with json marshal, err: %v", err)
			fmt.Fprintf(w, "get limiter rule failed with json marshal, err: %v", err)
			return
		}
		fmt.Fprintf(w, "%s", data)

	case common.SettingsAddAction:
		limiterRule := new(meta.LimiterRule)
		if err := json.NewDecoder(r.Body).Decode(&limiterRule); nil != err {
			blog.Errorf("add limiter rule failed with decode body err: %v", err)
			fmt.Fprintf(w, "add limiter rule failed with decode body err: %v", err)
			return
		}
		if err := s.limiter.AddRule(limiterRule); err != nil {
			blog.Errorf("add limiter rule failed err: %v", err)
			fmt.Fprintf(w, "add limiter rule failed err: %v", err)
			return
		}
		fmt.Fprintln(w, "success!")

	case common.SettingsDeleteAction:
		limiterRuleNames := new(meta.LimiterRuleNames)
		if err := json.NewDecoder(r.Body).Decode(&limiterRuleNames); nil != err {
			blog.Errorf("delete limiter rule failed with decode body err: %v", err)
			fmt.Fprintf(w, "delete limiter rule failed with decode body err: %v", err)
			return
		}
		if err := s.limiter.DelRules(limiterRuleNames.RuleNames); err != nil {
			blog.Errorf("delete limiter rule failed err: %v", err)
			fmt.Fprintf(w, "delete limiter rule failed err: %v", err)
			return
		}
		fmt.Fprintln(w, "success!")

	case common.SettingsGetAllAction:
		rules, err := s.limiter.GetAllRules()
		if err != nil {
			blog.Errorf("get all limiter rules error, err: %v", err)
			fmt.Fprintf(w, "get all limiter rules error, err: %v", err)
			return
		}
		data, err := json.Marshal(rules)
		if err != nil {
			blog.Errorf("get all limiter failed with json marshal, err: %v", err)
			fmt.Fprintf(w, "get all limiter failed with json marshal, err: %v", err)
			return
		}
		fmt.Fprintf(w, "%s", data)

	default:
		blog.Errorf("adjust limiter operation error, can't find the relevant operation action function " +
			"to adjust.")
		fmt.Fprintln(w, "adjust limiter operation error, can't find the relevant operation action function "+
			"to adjust.")
	}
}
