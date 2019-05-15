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
	s.addAction(http.MethodPost, "/list/process/service_category", s.ListServiceCategories, nil)
	s.addAction(http.MethodPut, "/update/process/service_category/{service_category_id}", s.UpdateServiceCategory, nil)
	s.addAction(http.MethodDelete, "/delete/process/service_category/{service_category_id}", s.DeleteServiceCategory, nil)
}
