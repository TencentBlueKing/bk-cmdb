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
)

func (s *coreService) initHostApplyRule() {
	s.addAction(http.MethodPost, "/create/host_apply_rule/bk_biz_id/{bk_biz_id}/", s.CreateHostApplyRule, nil)
	s.addAction(http.MethodPut, "/update/host_apply_rule/{host_apply_rule_id}/bk_biz_id/{bk_biz_id}/", s.UpdateHostApplyRule, nil)
	s.addAction(http.MethodDelete, "/deletemany/host_apply_rule/bk_biz_id/{bk_biz_id}/", s.DeleteHostApplyRule, nil)
	s.addAction(http.MethodGet, "/find/host_apply_rule/{host_apply_rule_id}/bk_biz_id/{bk_biz_id}/", s.GetHostApplyRule, nil)
	s.addAction(http.MethodPost, "/findmany/host_apply_rule/bk_biz_id/{bk_biz_id}/", s.ListHostApplyRule, nil)
	s.addAction(http.MethodPost, "/updatemany/host_apply_rule/bk_biz_id/{bk_biz_id}", s.BatchUpdateHostApplyRule, nil)
	s.addAction(http.MethodPost, "/findmany/host_apply_plan/bk_biz_id/{bk_biz_id}/", s.GenerateApplyPlan, nil)
	s.addAction(http.MethodPost, "/findmany/modules/bk_biz_id/{bk_biz_id}/host_apply_rule_related", s.SearchRuleRelatedModules, nil)
}
