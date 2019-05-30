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
	"fmt"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

func (p *process) CreateServiceCategory(ctx context.Context, h http.Header, category *metadata.ServiceCategory) (resp *metadata.ServiceCategory, err error) {
	ret := new(metadata.OneServiceCategoryResult)
	subPath := "/create/process/service_category"

	err = p.client.Post().
		WithContext(ctx).
		Body(category).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) GetServiceCategory(ctx context.Context, h http.Header, categoryID int64) (resp *metadata.ServiceCategory, err error) {
	ret := new(metadata.OneServiceCategoryResult)
	subPath := fmt.Sprintf("/find/process/service_category/%d", categoryID)

	err = p.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) UpdateServiceCategory(ctx context.Context, h http.Header, categoryID int64, category *metadata.ServiceCategory) (resp *metadata.ServiceCategory, err error) {
	ret := new(metadata.OneServiceCategoryResult)
	subPath := fmt.Sprintf("/update/process/service_category/%d", categoryID)

	err = p.client.Post().
		WithContext(ctx).
		Body(category).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) DeleteServiceCategory(ctx context.Context, h http.Header, categoryID int64) error {
	ret := new(metadata.OneServiceCategoryResult)
	subPath := fmt.Sprintf("/delete/process/service_category/%d", categoryID)

	err := p.client.Delete().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return err
	}
	if ret.Result == false || ret.Code != 0 {
		return errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return nil
}

func (p *process) ListServiceCategories(ctx context.Context, h http.Header, bizID int64, withStatistics bool) (resp *metadata.MultipleServiceCategory, err error) {
	ret := new(metadata.MultipleServiceCategoryResult)
	subPath := "/findmany/process/service_category"

	input := map[string]interface{}{
		"bk_biz_id":       bizID,
		"with_statistics": withStatistics,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

/*
	service template api
*/
func (p *process) CreateServiceTemplate(ctx context.Context, h http.Header, template *metadata.ServiceTemplate) (resp *metadata.ServiceTemplate, err error) {
	ret := new(metadata.OneServiceTemplateResult)
	subPath := "/create/process/service_template"

	err = p.client.Post().
		WithContext(ctx).
		Body(template).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) GetServiceTemplate(ctx context.Context, h http.Header, templateID int64) (resp *metadata.ServiceTemplate, err error) {
	ret := new(metadata.OneServiceTemplateResult)
	subPath := fmt.Sprintf("/find/process/service_template/%d", templateID)

	err = p.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) UpdateServiceTemplate(ctx context.Context, h http.Header, templateID int64, template *metadata.ServiceTemplate) (resp *metadata.ServiceTemplate, err error) {
	ret := new(metadata.OneServiceTemplateResult)
	subPath := fmt.Sprintf("/update/process/service_template/%d", templateID)

	err = p.client.Post().
		WithContext(ctx).
		Body(template).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) DeleteServiceTemplate(ctx context.Context, h http.Header, templateID int64) error {
	ret := new(metadata.OneServiceTemplateResult)
	subPath := fmt.Sprintf("/delete/process/service_template/%d", templateID)

	err := p.client.Delete().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return err
	}
	if ret.Result == false || ret.Code != 0 {
		return errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return nil
}

func (p *process) ListServiceTemplates(ctx context.Context, h http.Header, bizID int64, categoryID int64) (resp *metadata.MultipleServiceTemplate, err error) {
	ret := new(metadata.MultipleServiceTemplateResult)
	subPath := "/findmany/process/service_template"

	input := map[string]interface{}{
		"bk_biz_id":           bizID,
		"service_category_id": categoryID,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

/*
	process template api
*/
func (p *process) CreateProcessTemplate(ctx context.Context, h http.Header, template *metadata.ProcessTemplate) (resp *metadata.ProcessTemplate, err error) {
	ret := new(metadata.OneProcessTemplateResult)
	subPath := "/create/process/process_template"

	err = p.client.Post().
		WithContext(ctx).
		Body(template).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) GetProcessTemplate(ctx context.Context, h http.Header, templateID int64) (resp *metadata.ProcessTemplate, err error) {
	ret := new(metadata.OneProcessTemplateResult)
	subPath := fmt.Sprintf("/find/process/process_template/%d", templateID)

	err = p.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) UpdateProcessTemplate(ctx context.Context, h http.Header, templateID int64, template *metadata.ProcessTemplate) (resp *metadata.ProcessTemplate, err error) {
	ret := new(metadata.OneProcessTemplateResult)
	subPath := fmt.Sprintf("/update/process/process_template/%d", templateID)

	err = p.client.Put().
		WithContext(ctx).
		Body(template).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) DeleteProcessTemplate(ctx context.Context, h http.Header, templateID int64) error {
	ret := new(metadata.OneProcessTemplateResult)
	subPath := fmt.Sprintf("/delete/process/process_template/%d", templateID)

	err := p.client.Delete().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return err
	}
	if ret.Result == false || ret.Code != 0 {
		return errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return nil
}

func (p *process) DeleteProcessTemplateBatch(ctx context.Context, h http.Header, templateIDs []int64) error {
	ret := new(metadata.OneProcessTemplateResult)
	subPath := "/delete/process/process_template"

	input := map[string]interface{}{
		"process_template_ids": templateIDs,
	}

	err := p.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return err
	}
	if ret.Result == false || ret.Code != 0 {
		return errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return nil
}

func (p *process) ListProcessTemplates(ctx context.Context, h http.Header, option *metadata.ListProcessTemplatesOption) (resp *metadata.MultipleProcessTemplate, err error) {
	ret := new(metadata.MultipleProcessTemplateResult)
	subPath := "/findmany/process/process_template"

	err = p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

/*
	service instance api
*/
func (p *process) CreateServiceInstance(ctx context.Context, h http.Header, instance *metadata.ServiceInstance) (resp *metadata.ServiceInstance, err error) {
	ret := new(metadata.OneServiceInstanceResult)
	subPath := "/create/process/service_instance"

	err = p.client.Post().
		WithContext(ctx).
		Body(instance).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) GetServiceInstance(ctx context.Context, h http.Header, instanceID int64) (resp *metadata.ServiceInstance, err error) {
	ret := new(metadata.OneServiceInstanceResult)
	subPath := fmt.Sprintf("/find/process/service_instance/%d", instanceID)

	err = p.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) UpdateServiceInstance(ctx context.Context, h http.Header, instanceID int64, instance *metadata.ServiceInstance) (resp *metadata.ServiceInstance, err error) {
	ret := new(metadata.OneServiceInstanceResult)
	subPath := fmt.Sprintf("/update/process/service_instance/%d", instanceID)

	err = p.client.Post().
		WithContext(ctx).
		Body(instance).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) DeleteServiceInstance(ctx context.Context, h http.Header, instanceID int64) error {
	ret := new(metadata.OneServiceInstanceResult)
	subPath := fmt.Sprintf("/delete/process/service_instance/%d", instanceID)

	err := p.client.Delete().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return err
	}
	if ret.Result == false || ret.Code != 0 {
		return errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return nil
}

func (p *process) ListServiceInstance(ctx context.Context, h http.Header, option *metadata.ListServiceInstanceOption) (resp *metadata.MultipleServiceInstance, err error) {
	ret := new(metadata.MultipleServiceInstanceResult)
	subPath := "/findmany/process/service_instance"

	err = p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

/*
	process instance relation api
*/
func (p *process) CreateProcessInstanceRelation(ctx context.Context, h http.Header, instance *metadata.ProcessInstanceRelation) (resp *metadata.ProcessInstanceRelation, err error) {
	ret := new(metadata.OneProcessInstanceRelationResult)
	subPath := "/create/process/process_instance_relation"

	err = p.client.Post().
		WithContext(ctx).
		Body(instance).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) GetProcessInstanceRelation(ctx context.Context, h http.Header, processID int64) (resp *metadata.ProcessInstanceRelation, err error) {
	ret := new(metadata.OneProcessInstanceRelationResult)
	subPath := fmt.Sprintf("/find/process/process_instance_relation/%d", processID)

	err = p.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) UpdateProcessInstanceRelation(ctx context.Context, h http.Header, instanceID int64, instance *metadata.ProcessInstanceRelation) (resp *metadata.ProcessInstanceRelation, err error) {
	ret := new(metadata.OneProcessInstanceRelationResult)
	subPath := fmt.Sprintf("/update/process/process_instance_relation/%d", instanceID)

	err = p.client.Post().
		WithContext(ctx).
		Body(instance).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (p *process) DeleteProcessInstanceRelation(ctx context.Context, h http.Header, processID int64) error {
	ret := new(metadata.OneProcessInstanceRelationResult)
	subPath := fmt.Sprintf("/delete/process/process_instance_relation/%d", processID)

	err := p.client.Delete().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return err
	}
	if ret.Result == false || ret.Code != 0 {
		return errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return nil
}

func (p *process) ListProcessInstanceRelation(ctx context.Context, h http.Header, option *metadata.ListProcessInstanceRelationOption) (resp *metadata.MultipleProcessInstanceRelation, err error) {
	ret := new(metadata.MultipleProcessInstanceRelationResult)
	subPath := "/findmany/process/process_instance_relation"

	err = p.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, err
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}
