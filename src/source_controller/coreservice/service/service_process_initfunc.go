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

func (s *coreService) initProcess() {
	// service category
	s.addAction(http.MethodPost, "/create/process/service_category", s.CreateServiceCategory, nil)
	s.addAction(http.MethodGet, "/find/process/service_category/{service_category_id}", s.GetServiceCategory, nil)
	s.addAction(http.MethodGet, "/find/process/default_service_category", s.GetDefaultServiceCategory, nil)
	s.addAction(http.MethodPost, "/findmany/process/service_category", s.ListServiceCategories, nil)
	s.addAction(http.MethodPut, "/update/process/service_category/{service_category_id}", s.UpdateServiceCategory, nil)
	s.addAction(http.MethodDelete, "/delete/process/service_category/{service_category_id}", s.DeleteServiceCategory, nil)

	// service template
	s.addAction(http.MethodPost, "/create/process/service_template", s.CreateServiceTemplate, nil)
	s.addAction(http.MethodGet, "/find/process/service_template/{service_template_id}", s.GetServiceTemplate, nil)
	s.addAction(http.MethodGet, "/find/process/service_template/{service_template_id}/with_statistics", s.GetServiceTemplateWithStatistics, nil)
	s.addAction(http.MethodPost, "/findmany/process/service_template/detail/bk_biz_id/{bk_biz_id}", s.ListServiceTemplateDetail, nil)
	s.addAction(http.MethodPost, "/findmany/process/service_template", s.ListServiceTemplates, nil)
	s.addAction(http.MethodPut, "/update/process/service_template/{service_template_id}", s.UpdateServiceTemplate, nil)
	s.addAction(http.MethodDelete, "/delete/process/service_template/{service_template_id}", s.DeleteServiceTemplate, nil)

	// service instance
	s.addAction(http.MethodPost, "/create/process/service_instance", s.CreateServiceInstance, nil)
	s.addAction(http.MethodGet, "/find/process/service_instance/{service_instance_id}", s.GetServiceInstance, nil)
	s.addAction(http.MethodPost, "/findmany/process/service_instance", s.ListServiceInstances, nil)
	s.addAction(http.MethodPut, "/update/process/service_instance/{service_instance_id}", s.UpdateServiceInstance, nil)
	s.addAction(http.MethodDelete, "/delete/process/service_instance", s.DeleteServiceInstance, nil)
	s.addAction(http.MethodPost, "/update/process/service_instance_name/{service_instance_id}", s.ReconstructServiceInstanceName, nil)
	s.addAction(http.MethodPost, "/findmany/process/service_instance/details", s.ListServiceInstanceDetail, nil)

	// process template
	s.addAction(http.MethodPost, "/create/process/process_template", s.CreateProcessTemplate, nil)
	s.addAction(http.MethodGet, "/find/process/process_template/{process_template_id}", s.GetProcessTemplate, nil)
	s.addAction(http.MethodPost, "/findmany/process/process_template", s.ListProcessTemplates, nil)
	s.addAction(http.MethodPut, "/update/process/process_template/{process_template_id}", s.UpdateProcessTemplate, nil)
	s.addAction(http.MethodDelete, "/delete/process/process_template/{process_template_id}", s.DeleteProcessTemplate, nil)
	s.addAction(http.MethodPost, "/delete/process/process_template", s.BatchDeleteProcessTemplate, nil)

	// process instance relation
	s.addAction(http.MethodPost, "/create/process/process_instance_relation", s.CreateProcessInstanceRelation, nil)
	s.addAction(http.MethodGet, "/find/process/process_instance_relation/{process_instance_id}", s.GetProcessInstanceRelation, nil)
	s.addAction(http.MethodPost, "/findmany/process/process_instance_relation", s.ListProcessInstanceRelation, nil)
	s.addAction(http.MethodPut, "/update/process/process_instance_relation/{process_instance_id}", s.UpdateProcessInstanceRelation, nil)
	s.addAction(http.MethodDelete, "/delete/process/process_instance_relation", s.DeleteProcessInstanceRelation, nil)

	s.addAction(http.MethodGet, "/find/process/business_default_set_module_info/{bk_biz_id}", s.GetBusinessDefaultSetModuleInfo, nil)
	s.addAction(http.MethodDelete, "/delete/process/module_bound_template/{bk_module_id}", s.RemoveTemplateBindingOnModule, nil)
	s.addAction(http.MethodPost, "/findmany/process/proc2module", s.GetProc2Module, nil)
}
