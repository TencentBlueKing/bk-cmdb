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

func (s *coreService) initSetTemplate() {
	s.addAction(http.MethodPost, "/create/topo/set_template/bk_biz_id/{bk_biz_id}/", s.CreateSetTemplate, nil)
	s.addAction(http.MethodPut, "/update/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/", s.UpdateSetTemplate, nil)
	s.addAction(http.MethodDelete, "/deletemany/topo/set_template/bk_biz_id/{bk_biz_id}/", s.DeleteSetTemplate, nil)
	s.addAction(http.MethodGet, "/find/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/", s.GetSetTemplate, nil)
	s.addAction(http.MethodPost, "/findmany/topo/set_template/bk_biz_id/{bk_biz_id}/", s.ListSetTemplate, nil)
	s.addAction(http.MethodPost, "/findmany/topo/set_template/count_instances/bk_biz_id/{bk_biz_id}/", s.CountSetTplInstances, nil)
	s.addAction(http.MethodGet, "/findmany/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/service_templates", s.ListSetTplRelatedSvcTpl, nil)
	s.addAction(http.MethodPut, "/update/topo/set_template_sync_status/bk_set_id/{bk_set_id}", s.UpdateSetTemplateSyncStatus, nil)
	s.addAction(http.MethodPost, "/findmany/topo/set_template_sync_status/bk_biz_id/{bk_biz_id}", s.ListSetTemplateSyncStatus, nil)
	s.addAction(http.MethodPost, "/findmany/topo/set_template_sync_history/bk_biz_id/{bk_biz_id}", s.ListSetTemplateSyncHistory, nil)
	s.addAction(http.MethodDelete, "/deletemany/topo/set_template_sync_status/bk_biz_id/{bk_biz_id}", s.DeleteSetTemplateSyncStatus, nil)
}
