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
	subPath := "/findmany/topo/set_template/%d/bk_biz_id/%d/service_templates"

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

func (p *setTemplate) ListSetTplRelatedSvcTpl(ctx context.Context, header http.Header, bizID int64, setTemplateID int64) ([]metadata.ServiceTemplate, errors.CCErrorCoder) {
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
		blog.Errorf("ListSetTplRelatedSvcTpl failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret.Data, nil
}

func (p *setTemplate) UpdateSetTemplateSyncStatus(ctx context.Context, header http.Header, setID int64, syncStatus metadata.SetTemplateSyncStatus) errors.CCErrorCoder {
	ret := struct {
		metadata.BaseResp
	}{}
	subPath := "/update/topo/set_template_sync_status/bk_set_id/%d"

	err := p.client.Put().
		WithContext(ctx).
		Body(syncStatus).
		SubResourcef(subPath, setID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("UpdateSetTemplateSyncStatus failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

func (p *setTemplate) DeleteSetTemplateSyncStatus(ctx context.Context, header http.Header, bizID int64, setIDs []int64) errors.CCErrorCoder {
	ret := struct {
		metadata.BaseResp
	}{}
	subPath := "/deletemany/topo/set_template_sync_status/bk_biz_id/%d"

	option := metadata.DeleteSetTemplateSyncStatusOption{
		SetIDs: setIDs,
		BizID:  bizID,
	}
	err := p.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("DeleteSetTemplateSyncStatus failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}

func (p *setTemplate) ListSetTemplateSyncStatus(ctx context.Context, header http.Header, bizID int64, option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data metadata.MultipleSetTemplateSyncStatus
	}{}
	subPath := "/findmany/topo/set_template_sync_status/bk_biz_id/%d"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("ListSetTemplateSyncStatus failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}

func (p *setTemplate) ListSetTemplateSyncHistory(ctx context.Context, header http.Header, bizID int64, option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data metadata.MultipleSetTemplateSyncStatus
	}{}
	subPath := "/findmany/topo/set_template_sync_history/bk_biz_id/%d"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("ListSetTemplateSyncHistory failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}

func (p *setTemplate) ModifySetTemplateSyncStatus(ctx context.Context, header http.Header, setID int64, syncStatus metadata.SyncStatus) errors.CCErrorCoder {
	ret := struct {
		metadata.BaseResp
	}{}
	subPath := "/update/topo/set_template_sync_status/bk_set_id/%d/status/%s"

	err := p.client.Put().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, setID, syncStatus).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("UpdateSetTemplateSyncStatus failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}

	return ret.CCError()
}
