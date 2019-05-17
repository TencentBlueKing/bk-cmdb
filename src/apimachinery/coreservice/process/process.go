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

package process

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type ProcessInterface interface {
	// service category
	CreateServiceCategory(ctx context.Context, h http.Header, category *metadata.ServiceCategory) (resp *metadata.ServiceCategory, err error)
	GetServiceCategory(ctx context.Context, h http.Header, categoryID int64) (resp *metadata.ServiceCategory, err error)
	UpdateServiceCategory(ctx context.Context, h http.Header, categoryID int64, category metadata.ServiceCategory) (resp *metadata.ServiceCategory, err error)
	ListServiceCategories(ctx context.Context, h http.Header, bizID int64, withStatistics bool) (resp *metadata.ServiceCategoryWithStatistics, err error)
	DeleteServiceCategory(ctx context.Context, h http.Header, categoryID int64) error

	// service template
	CreateServiceTemplate(ctx context.Context, h http.Header, template metadata.ServiceTemplate) (resp *metadata.ServiceTemplate, err error)
	GetServiceTemplate(ctx context.Context, h http.Header, templateID int64) (resp *metadata.ServiceTemplate, err error)
	UpdateServiceTemplate(ctx context.Context, h http.Header, templateID int64, template metadata.ServiceTemplate) (resp *metadata.ServiceTemplate, err error)
	ListServiceTemplates(ctx context.Context, h http.Header, bizID int64, categoryID int64) (resp *metadata.MultipleServiceTemplate, err error)
	DeleteServiceTemplate(ctx context.Context, h http.Header, serviceTemplateID int64) error

	// process template
	CreateProcessTemplate(ctx context.Context, h http.Header, template metadata.ProcessTemplate) (resp *metadata.ProcessTemplate, err error)
	GetProcessTemplate(ctx context.Context, h http.Header, templateID int64) (resp *metadata.ProcessTemplate, err error)
	UpdateProcessTemplate(ctx context.Context, h http.Header, templateID int64, template metadata.ProcessTemplate) (resp *metadata.ProcessTemplate, err error)
	ListProcessTemplates(ctx context.Context, h http.Header, bizID int64, serviceTemplateID int64) (resp *metadata.MultipleProcessTemplate, err error)
	DeleteProcessTemplate(ctx context.Context, h http.Header, processTemplateID int64) error
	BatchDeleteProcessTemplate(ctx context.Context, h http.Header, processTemplateIDs []int64) error

	// service instance
	CreateServiceInstance(ctx context.Context, h http.Header, template metadata.ServiceInstance) (resp *metadata.ServiceInstance, err error)
	GetServiceInstance(ctx context.Context, h http.Header, templateID int64) (resp *metadata.ServiceInstance, err error)
	UpdateServiceInstance(ctx context.Context, h http.Header, templateID int64, template metadata.ServiceInstance) (resp *metadata.ServiceInstance, err error)
	ListServiceInstance(ctx context.Context, h http.Header, bizID int64, serviceTemplateID int64, hostID int64) (resp *metadata.MultipleServiceInstance, err error)
	DeleteServiceInstance(ctx context.Context, h http.Header, processTemplateID int64) error

	// service instance relation
	CreateProcessInstanceRelation(ctx context.Context, h http.Header, relation metadata.ProcessInstanceRelation) (resp *metadata.ProcessInstanceRelation, err error)
	GetProcessInstanceRelation(ctx context.Context, h http.Header, processID int64) (resp *metadata.ProcessInstanceRelation, err error)
	UpdateProcessInstanceRelation(ctx context.Context, h http.Header, processID int64, template metadata.ProcessInstanceRelation) (resp *metadata.ProcessInstanceRelation, err error)
	ListProcessInstanceRelation(ctx context.Context, h http.Header, bizID int64, serviceInstanceID int64, hostID int64) (resp *metadata.MultipleProcessInstanceRelation, err error)
	DeleteProcessInstanceRelation(ctx context.Context, h http.Header, processID int64) error
}

func NewProcessInterfaceClient(client rest.ClientInterface) ProcessInterface {
	return &process{client: client}
}

type process struct {
	client rest.ClientInterface
}
