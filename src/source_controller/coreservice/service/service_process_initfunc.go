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

func (s *coreService) initProcess(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	// service category
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/process/service_category", Handler: s.CreateServiceCategory})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/process/service_category/{service_category_id}", Handler: s.GetServiceCategory})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/process/default_service_category", Handler: s.GetDefaultServiceCategory})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/process/service_category", Handler: s.ListServiceCategories})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/process/service_category/{service_category_id}", Handler: s.UpdateServiceCategory})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/process/service_category/{service_category_id}", Handler: s.DeleteServiceCategory})

	// service template
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/process/service_template", Handler: s.CreateServiceTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/process/service_template/{service_template_id}", Handler: s.GetServiceTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/process/service_template/{service_template_id}/with_statistics", Handler: s.GetServiceTemplateWithStatistics})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/process/service_template/detail/bk_biz_id/{bk_biz_id}", Handler: s.ListServiceTemplateDetail})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/process/service_template", Handler: s.ListServiceTemplates})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/process/service_template/{service_template_id}", Handler: s.UpdateServiceTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/process/service_template/{service_template_id}", Handler: s.DeleteServiceTemplate})

	// service instance
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/process/service_instance", Handler: s.CreateServiceInstance})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/process/service_instance", Handler: s.CreateServiceInstances})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/process/service_instance/{service_instance_id}", Handler: s.GetServiceInstance})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/process/service_instance", Handler: s.ListServiceInstances})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/process/service_instance/biz/{bk_biz_id}", Handler: s.UpdateServiceInstances})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/process/service_instance", Handler: s.DeleteServiceInstance})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/update/process/service_instance_name", Handler: s.ConstructServiceInstanceName})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/update/process/service_instance_name/{service_instance_id}", Handler: s.ReconstructServiceInstanceName})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/process/service_instance/details", Handler: s.ListServiceInstanceDetail})

	// process template
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/process/process_template", Handler: s.CreateProcessTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/process/process_template/{process_template_id}", Handler: s.GetProcessTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/process/process_template", Handler: s.ListProcessTemplates})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/process/process_template/{process_template_id}", Handler: s.UpdateProcessTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/process/process_template/{process_template_id}", Handler: s.DeleteProcessTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/delete/process/process_template", Handler: s.BatchDeleteProcessTemplate})

	// process instance relation
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/process/process_instance_relation", Handler: s.CreateProcessInstanceRelation})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/process/process_instance_relation", Handler: s.CreateProcessInstanceRelations})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/process/process_instance_relation/{process_instance_id}", Handler: s.GetProcessInstanceRelation})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/process/process_instance_relation", Handler: s.ListProcessInstanceRelation})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/process/process_instance_relation/{process_instance_id}", Handler: s.UpdateProcessInstanceRelation})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/process/process_instance_relation", Handler: s.DeleteProcessInstanceRelation})

	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/process/business_default_set_module_info/{bk_biz_id}", Handler: s.GetBusinessDefaultSetModuleInfo})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/process/module_bound_template/{bk_module_id}", Handler: s.RemoveTemplateBindingOnModule})

	utility.AddToRestfulWebService(web)
}
