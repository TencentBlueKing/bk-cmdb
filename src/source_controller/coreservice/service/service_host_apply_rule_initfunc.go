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
	"net/http"

	"configcenter/src/common/http/rest"

	"github.com/emicklei/go-restful"
)

func (s *coreService) initHostApplyRule(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/host_apply_rule/bk_biz_id/{bk_biz_id}/", Handler: s.CreateHostApplyRule})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/host_apply_rule/{host_apply_rule_id}/bk_biz_id/{bk_biz_id}/", Handler: s.UpdateHostApplyRule})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/host_apply_rule/bk_biz_id/{bk_biz_id}/", Handler: s.DeleteHostApplyRule})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/host_apply_rule/{host_apply_rule_id}/bk_biz_id/{bk_biz_id}/", Handler: s.GetHostApplyRule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/host_apply_rule/bk_biz_id/{bk_biz_id}/", Handler: s.ListHostApplyRule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/updatemany/host_apply_rule/bk_biz_id/{bk_biz_id}", Handler: s.BatchUpdateHostApplyRule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/host_apply_plan/bk_biz_id/{bk_biz_id}/", Handler: s.GenerateApplyPlan})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/modules/bk_biz_id/{bk_biz_id}/host_apply_rule_related", Handler: s.SearchRuleRelatedModules})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/host/bk_biz_id/{bk_biz_id}/update_by_host_apply", Handler: s.UpdateHostByHostApplyRule})

	utility.AddToRestfulWebService(web)
}
