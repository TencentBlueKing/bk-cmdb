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

package settemplate

import (
	"context"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// CreateSetTemplate TODO
func (p *setTemplate) CreateSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.CreateSetTemplateOption) (metadata.SetTemplate, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp `json:",inline"`
		Data              metadata.SetTemplate `json:"data"`
	}{}
	subPath := "/create/topo/set_template/bk_biz_id/%d/"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("CreateSetTemplate failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}

// UpdateSetTemplate TODO
func (p *setTemplate) UpdateSetTemplate(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.UpdateSetTemplateOption) (metadata.SetTemplate, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp `json:",inline"`
		Data              metadata.SetTemplate `json:"data"`
	}{}
	subPath := "/update/topo/set_template/%d/bk_biz_id/%d/"

	err := p.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, setTemplateID, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("UpdateSetTemplate failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}

// DeleteSetTemplate TODO
func (p *setTemplate) DeleteSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.DeleteSetTemplateOption) errors.CCErrorCoder {
	ret := struct {
		metadata.BaseResp `json:",inline"`
	}{}
	subPath := "/deletemany/topo/set_template/bk_biz_id/%d/"

	err := p.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("DeleteSetTemplate failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

// GetSetTemplate TODO
func (p *setTemplate) GetSetTemplate(ctx context.Context, header http.Header, bizID int64, setTemplateID int64) (metadata.SetTemplate, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp `json:",inline"`
		Data              metadata.SetTemplate `json:"data"`
	}{}
	subPath := "/find/topo/set_template/%d/bk_biz_id/%d/"

	err := p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, setTemplateID, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("GetSetTemplate failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}

// ListSetTemplate TODO
func (p *setTemplate) ListSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.ListSetTemplateOption) (*metadata.MultipleSetTemplateResult, errors.CCErrorCoder) {
	ret := metadata.ListSetTemplateResult{}
	subPath := "/findmany/topo/set_template/bk_biz_id/%d/"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("ListSetTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return &ret.Data, nil
}

// CountSetTplInstances TODO
func (p *setTemplate) CountSetTplInstances(ctx context.Context, header http.Header, bizID int64, option metadata.CountSetTplInstOption) (map[int64]int64, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data []metadata.CountSetTplInstItem `json:"data"`
	}{}
	subPath := "/findmany/topo/set_template/count_instances/bk_biz_id/%d/"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("CountSetTplInstances failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	data := make(map[int64]int64)
	for _, item := range ret.Data {
		data[item.SetTemplateID] = item.SetInstanceCount
	}

	return data, nil
}

// ListSetServiceTemplateRelations get relations of SetTemplate <==> ServiceTemplate
func (p *setTemplate) ListSetServiceTemplateRelations(ctx context.Context, header http.Header, bizID int64, setTemplateID int64) ([]metadata.SetServiceTemplateRelation, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data []metadata.SetServiceTemplateRelation `json:"data"`
	}{}
	subPath := "/findmany/topo/set_template/%d/bk_biz_id/%d/service_templates_relations"

	err := p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, setTemplateID, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("ListSetServiceTemplateRelations failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret.Data, nil
}

// ListSetTplRelatedSvcTpl search related about set template and service template
func (p *setTemplate) ListSetTplRelatedSvcTpl(ctx context.Context, header http.Header, bizID int64,
	setTemplateID int64) ([]metadata.ServiceTemplate, errors.CCErrorCoder) {

	ret := struct {
		metadata.BaseResp
		Data []metadata.ServiceTemplate `json:"data"`
	}{}
	subPath := "/findmany/topo/set_template/%d/bk_biz_id/%d/service_templates"

	err := p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, setTemplateID, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := ret.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return ret.Data, nil
}

// CreateSetTemplateAttribute create set template attribute
func (p *setTemplate) CreateSetTemplateAttribute(ctx context.Context, h http.Header,
	option *metadata.CreateSetTempAttrsOption) ([]int64, errors.CCErrorCoder) {

	ret := new(metadata.CreateBatchResult)
	subPath := "/create/set_template/attribute"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret.Data.IDs, nil
}

// UpdateSetTemplateAttribute update set template attribute
func (p *setTemplate) UpdateSetTemplateAttribute(ctx context.Context, h http.Header,
	option *metadata.UpdateSetTempAttrOption) errors.CCErrorCoder {

	ret := new(metadata.BaseResp)
	subPath := "/update/set_template/attribute"

	err := p.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

// DeleteSetTemplateAttribute delete set template attribute
func (p *setTemplate) DeleteSetTemplateAttribute(ctx context.Context, h http.Header,
	option *metadata.DeleteSetTempAttrOption) errors.CCErrorCoder {

	ret := new(metadata.BaseResp)
	subPath := "/delete/set_template/attribute"

	err := p.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

// ListSetTemplateAttribute list set Template Attribute
func (p *setTemplate) ListSetTemplateAttribute(ctx context.Context, h http.Header,
	option *metadata.ListSetTempAttrOption) (*metadata.SetTempAttrData, errors.CCErrorCoder) {

	ret := new(metadata.SetTemplateAttributeResult)
	subPath := "/findmany/set_template/attribute"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret.Data, nil
}
