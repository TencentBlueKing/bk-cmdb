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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/backbone/setting"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

// SetLimitRule adjust current limit rules
func (s *service) SetLimitRule(req *restful.Request, resp *restful.Response) {
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))
	action := req.QueryParameter(setting.SettingsAction)

	switch setting.ActionType(action) {
	case setting.Get:
		limiterRuleNames := new(meta.LimiterRuleNames)
		if err := json.NewDecoder(req.Request.Body).Decode(&limiterRuleNames); nil != err {
			blog.Errorf("get limiter rule failed with decode body err: %v", err)
			resp.WriteError(http.StatusBadRequest,
				&meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
			return
		}
		rules, err := s.limiter.GetRules(limiterRuleNames.RuleNames)
		if err != nil {
			blog.Errorf("get limiter rule failed err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
			return
		}
		resp.WriteEntity(meta.NewSuccessResp(rules))
		return

	case setting.Add:
		limiterRule := new(meta.LimiterRule)
		if err := json.NewDecoder(req.Request.Body).Decode(&limiterRule); nil != err {
			blog.Errorf("add limiter rule failed with decode body err: %v", err)
			resp.WriteError(http.StatusBadRequest,
				&meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
			return
		}
		if err := s.limiter.AddRule(limiterRule); err != nil {
			blog.Errorf("add limiter rule failed err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
			return
		}
		resp.WriteEntity(meta.NewSuccessResp(nil))
		return

	case setting.Delete:
		limiterRuleNames := new(meta.LimiterRuleNames)
		if err := json.NewDecoder(req.Request.Body).Decode(&limiterRuleNames); nil != err {
			blog.Errorf("delete limiter rule failed with decode body err: %v", err)
			resp.WriteError(http.StatusBadRequest,
				&meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
			return
		}
		if err := s.limiter.DelRules(limiterRuleNames.RuleNames); err != nil {
			blog.Errorf("delete limiter rule failed err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
			return
		}
		resp.WriteEntity(meta.NewSuccessResp(nil))
		return

	case setting.GetAll:
		rules, err := s.limiter.GetAllRules()
		if err != nil {
			blog.Errorf("get all limiter rules error, err: %v", err)
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: err})
			return
		}
		resp.WriteEntity(meta.NewSuccessResp(rules))
		return

	default:
		blog.Errorf("adjust limiter operation error, can't find the relevant operation action function " +
			"to adjust.")
		resp.WriteError(http.StatusBadRequest,
			&meta.RespError{Msg: defErr.Error(common.CCErrAPINoCurrentLimitingOperation)})
	}
}
