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

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

func (p *process) CreateServiceCategory(ctx context.Context, h http.Header, category *metadata.ServiceCategory) (*metadata.ServiceCategory, errors.CCErrorCoder) {
	ret := new(metadata.OneServiceCategoryResult)
	subPath := "/create/process/service_category"

	err := p.client.Post().
		WithContext(ctx).
		Body(category).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("CreateServiceCategory failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) GetServiceCategory(ctx context.Context, h http.Header, categoryID int64) (*metadata.ServiceCategory, errors.CCErrorCoder) {
	ret := new(metadata.OneServiceCategoryResult)
	subPath := "/find/process/service_category/%d"

	err := p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, categoryID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("GetServiceCategory failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) GetDefaultServiceCategory(ctx context.Context, h http.Header) (*metadata.ServiceCategory, errors.CCErrorCoder) {
	ret := new(metadata.OneServiceCategoryResult)
	subPath := "/find/process/default_service_category"

	err := p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("GetDefaultServiceCategory failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) UpdateServiceCategory(ctx context.Context, h http.Header, categoryID int64, category *metadata.ServiceCategory) (*metadata.ServiceCategory, errors.CCErrorCoder) {
	ret := new(metadata.OneServiceCategoryResult)
	subPath := "/update/process/service_category/%d"

	err := p.client.Put().
		WithContext(ctx).
		Body(category).
		SubResourcef(subPath, categoryID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("UpdateServiceCategory failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) DeleteServiceCategory(ctx context.Context, h http.Header, categoryID int64) errors.CCErrorCoder {
	ret := new(metadata.OneServiceCategoryResult)
	subPath := "/delete/process/service_category/%d"

	err := p.client.Delete().
		WithContext(ctx).
		SubResourcef(subPath, categoryID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("DeleteServiceCategory failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

func (p *process) ListServiceCategories(ctx context.Context, h http.Header, option metadata.ListServiceCategoriesOption) (*metadata.MultipleServiceCategoryWithStatistics, errors.CCErrorCoder) {
	ret := new(metadata.MultipleServiceCategoryWithStatisticsResult)
	subPath := "/findmany/process/service_category"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("ListServiceCategories failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

/*
	service template api
*/
func (p *process) CreateServiceTemplate(ctx context.Context, h http.Header, template *metadata.ServiceTemplate) (*metadata.ServiceTemplate, errors.CCErrorCoder) {
	ret := new(metadata.OneServiceTemplateResult)
	subPath := "/create/process/service_template"

	err := p.client.Post().
		WithContext(ctx).
		Body(template).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("CreateServiceTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) ListServiceTemplateDetail(ctx context.Context, h http.Header, bizID int64, templateIDs ...int64) (metadata.MultipleServiceTemplateDetail, errors.CCErrorCoder) {
	ret := new(metadata.MultipleServiceTemplateDetailResult)
	subPath := "/findmany/process/service_template/detail/bk_biz_id/%d"

	option := map[string]interface{}{
		"service_template_ids": templateIDs,
	}
	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("ListServiceTemplateDetail failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}

func (p *process) GetServiceTemplateWithStatistics(ctx context.Context, h http.Header, templateID int64) (*metadata.ServiceTemplateWithStatistics, errors.CCErrorCoder) {
	ret := new(metadata.OneServiceTemplateWithStatisticsResult)
	subPath := "/find/process/service_template/%d/with_statistics"

	err := p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, templateID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("GetServiceTemplateDetail failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) GetServiceTemplate(ctx context.Context, h http.Header, templateID int64) (*metadata.ServiceTemplate, errors.CCErrorCoder) {
	ret := new(metadata.OneServiceTemplateResult)
	subPath := "/find/process/service_template/%d"

	err := p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, templateID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("GetServiceTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) UpdateServiceTemplate(ctx context.Context, h http.Header, templateID int64, template *metadata.ServiceTemplate) (*metadata.ServiceTemplate, errors.CCErrorCoder) {
	ret := new(metadata.OneServiceTemplateResult)
	subPath := "/update/process/service_template/%d"

	err := p.client.Put().
		WithContext(ctx).
		Body(template).
		SubResourcef(subPath, templateID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("UpdateServiceTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) DeleteServiceTemplate(ctx context.Context, h http.Header, templateID int64) errors.CCErrorCoder {
	ret := new(metadata.OneServiceTemplateResult)
	subPath := "/delete/process/service_template/%d"

	err := p.client.Delete().
		WithContext(ctx).
		SubResourcef(subPath, templateID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("DeleteServiceTemplate failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

func (p *process) ListServiceTemplates(ctx context.Context, h http.Header, option *metadata.ListServiceTemplateOption) (*metadata.MultipleServiceTemplate, errors.CCErrorCoder) {
	ret := new(metadata.MultipleServiceTemplateResult)
	subPath := "/findmany/process/service_template"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("ListServiceTemplates failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) CreateProcessTemplate(ctx context.Context, h http.Header, template *metadata.ProcessTemplate) (*metadata.ProcessTemplate, errors.CCErrorCoder) {
	ret := new(metadata.OneProcessTemplateResult)
	subPath := "/create/process/process_template"

	err := p.client.Post().
		WithContext(ctx).
		Body(template).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("CreateProcessTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) GetProcessTemplate(ctx context.Context, h http.Header, templateID int64) (*metadata.ProcessTemplate, errors.CCErrorCoder) {
	ret := new(metadata.OneProcessTemplateResult)
	subPath := "/find/process/process_template/%d"

	err := p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, templateID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("GetProcessTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) UpdateProcessTemplate(ctx context.Context, h http.Header, templateID int64, property map[string]interface{}) (*metadata.ProcessTemplate, errors.CCErrorCoder) {
	ret := new(metadata.OneProcessTemplateResult)
	subPath := "/update/process/process_template/%d"

	err := p.client.Put().
		WithContext(ctx).
		Body(property).
		SubResourcef(subPath, templateID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("UpdateProcessTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) DeleteProcessTemplate(ctx context.Context, h http.Header, templateID int64) errors.CCErrorCoder {
	ret := new(metadata.OneProcessTemplateResult)
	subPath := "/delete/process/process_template/%d"

	err := p.client.Delete().
		WithContext(ctx).
		SubResourcef(subPath, templateID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("DeleteProcessTemplate failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

func (p *process) DeleteProcessTemplateBatch(ctx context.Context, h http.Header, templateIDs []int64) errors.CCErrorCoder {
	ret := new(metadata.OneProcessTemplateResult)
	subPath := "/delete/process/process_template"

	input := map[string]interface{}{
		"process_template_ids": templateIDs,
	}

	err := p.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("DeleteProcessTemplateBatch failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

func (p *process) ListProcessTemplates(ctx context.Context, h http.Header, option *metadata.ListProcessTemplatesOption) (*metadata.MultipleProcessTemplate, errors.CCErrorCoder) {
	ret := new(metadata.MultipleProcessTemplateResult)
	subPath := "/findmany/process/process_template"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("ListProcessTemplates failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

/*
	service instance api
*/
func (p *process) CreateServiceInstance(ctx context.Context, h http.Header, instance *metadata.ServiceInstance) (*metadata.ServiceInstance, errors.CCErrorCoder) {
	ret := new(metadata.OneServiceInstanceResult)
	subPath := "/create/process/service_instance"

	err := p.client.Post().
		WithContext(ctx).
		Body(instance).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("CreateServiceInstance failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) CreateServiceInstances(ctx context.Context, h http.Header, instances []*metadata.ServiceInstance) ([]*metadata.ServiceInstance, errors.CCErrorCoder) {
	ret := new(metadata.ManyServiceInstanceResult)
	subPath := "/createmany/process/service_instance"

	err := p.client.Post().
		WithContext(ctx).
		Body(instances).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("CreateServiceInstances failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret.Data, nil
}

func (p *process) GetServiceInstance(ctx context.Context, h http.Header, instanceID int64) (*metadata.ServiceInstance, errors.CCErrorCoder) {
	ret := new(metadata.OneServiceInstanceResult)
	subPath := "/find/process/service_instance/%d"

	err := p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, instanceID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("GetServiceInstance failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) UpdateServiceInstances(ctx context.Context, h http.Header, bizID int64, option *metadata.UpdateServiceInstanceOption) errors.CCErrorCoder {
	ret := new(metadata.OneServiceInstanceResult)
	subPath := "/updatemany/process/service_instance/biz/%d"

	err := p.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("UpdateServiceInstances failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

func (p *process) DeleteServiceInstance(ctx context.Context, h http.Header, option *metadata.CoreDeleteServiceInstanceOption) errors.CCErrorCoder {
	ret := new(metadata.OneServiceInstanceResult)
	subPath := "/delete/process/service_instance"

	err := p.client.Delete().
		Body(option).
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("DeleteServiceInstance failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

func (p *process) ListServiceInstance(ctx context.Context, h http.Header, option *metadata.ListServiceInstanceOption) (*metadata.MultipleServiceInstance, errors.CCErrorCoder) {
	ret := new(metadata.MultipleServiceInstanceResult)
	subPath := "/findmany/process/service_instance"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("ListServiceInstance failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) ListServiceInstanceDetail(ctx context.Context, h http.Header, option *metadata.ListServiceInstanceDetailOption) (*metadata.MultipleServiceInstanceDetail, errors.CCErrorCoder) {
	ret := new(metadata.MultipleServiceInstanceDetailResult)
	subPath := "/findmany/process/service_instance/details"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("ListServiceInstanceDetail failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

/*
	process instance relation api
*/
func (p *process) CreateProcessInstanceRelation(ctx context.Context, h http.Header, relation *metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {
	ret := new(metadata.OneProcessInstanceRelationResult)
	subPath := "/create/process/process_instance_relation"

	err := p.client.Post().
		WithContext(ctx).
		Body(relation).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) CreateProcessInstanceRelations(ctx context.Context, h http.Header, relations []*metadata.ProcessInstanceRelation) ([]*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {
	ret := new(metadata.ManyProcessInstanceRelationResult)
	subPath := "/createmany/process/process_instance_relation"

	err := p.client.Post().
		WithContext(ctx).
		Body(relations).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("CreateProcessInstanceRelations failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret.Data, nil
}

func (p *process) GetProcessInstanceRelation(ctx context.Context, h http.Header, processID int64) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {
	ret := new(metadata.OneProcessInstanceRelationResult)
	subPath := "/find/process/process_instance_relation/%d"

	err := p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, processID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("GetProcessInstanceRelation failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) UpdateProcessInstanceRelation(ctx context.Context, h http.Header, instanceID int64, instance *metadata.ProcessInstanceRelation) (*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {
	ret := new(metadata.OneProcessInstanceRelationResult)
	subPath := "/update/process/process_instance_relation/%d"

	err := p.client.Put().
		WithContext(ctx).
		Body(instance).
		SubResourcef(subPath, instanceID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("UpdateProcessInstanceRelation failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) DeleteProcessInstanceRelation(ctx context.Context, h http.Header, option metadata.DeleteProcessInstanceRelationOption) errors.CCErrorCoder {
	ret := new(metadata.OneProcessInstanceRelationResult)
	subPath := "/delete/process/process_instance_relation"

	err := p.client.Delete().
		WithContext(ctx).
		SubResourcef(subPath).
		Body(option).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("DeleteProcessInstanceRelation failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

func (p *process) ListProcessInstanceRelation(ctx context.Context, h http.Header, option *metadata.ListProcessInstanceRelationOption) (*metadata.MultipleProcessInstanceRelation, errors.CCErrorCoder) {
	ret := new(metadata.MultipleProcessInstanceRelationResult)
	subPath := "/findmany/process/process_instance_relation"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("ListProcessInstanceRelation failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

func (p *process) GetBusinessDefaultSetModuleInfo(ctx context.Context, h http.Header, bizID int64) (metadata.BusinessDefaultSetModuleInfo, errors.CCErrorCoder) {
	ret := new(metadata.BusinessDefaultSetModuleInfoResult)
	subPath := "/find/process/business_default_set_module_info/%d"

	emptyInfo := metadata.BusinessDefaultSetModuleInfo{}
	err := p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("GetBusinessDefaultSetModuleInfo failed, http request failed, err: %+v", err)
		return emptyInfo, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return emptyInfo, ret.CCError()
	}

	return ret.Data, nil
}

func (p *process) RemoveTemplateBindingOnModule(ctx context.Context, h http.Header, moduleID int64) (*metadata.RemoveTemplateBoundOnModuleResult, errors.CCErrorCoder) {
	ret := new(metadata.RemoveTemplateBoundOnModuleResult)
	subPath := "/delete/process/module_bound_template/%d"

	err := p.client.Delete().
		WithContext(ctx).
		SubResourcef(subPath, moduleID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("GetBusinessDefaultSetModuleInfo failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return nil, nil
}

func (p *process) ConstructServiceInstanceName(ctx context.Context, h http.Header, params *metadata.SrvInstNameParams) errors.CCErrorCoder {
	ret := new(metadata.RemoveTemplateBoundOnModuleResult)
	subPath := "/update/process/service_instance_name"

	err := p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("ReconstructServiceInstanceName failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

func (p *process) ReconstructServiceInstanceName(ctx context.Context, h http.Header, instanceID int64) errors.CCErrorCoder {
	ret := new(metadata.RemoveTemplateBoundOnModuleResult)
	subPath := "/update/process/service_instance_name/%d"

	err := p.client.Post().
		WithContext(ctx).
		SubResourcef(subPath, instanceID).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("ReconstructServiceInstanceName failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}
