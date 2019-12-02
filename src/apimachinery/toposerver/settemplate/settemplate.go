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
	"fmt"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

func (st *SetTemplate) CreateSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.CreateSetTemplateOption) (*metadata.SetTemplateResult, errors.CCErrorCoder) {
	ret := new(metadata.SetTemplateResult)
	subPath := fmt.Sprintf("/create/topo/set_template/bk_biz_id/%d/", bizID)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("CreateSetTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return ret, nil
}

func (st *SetTemplate) UpdateSetTemplate(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.UpdateSetTemplateOption) (*metadata.SetTemplateResult, errors.CCErrorCoder) {
	ret := new(metadata.SetTemplateResult)
	subPath := fmt.Sprintf("/update/topo/set_template/%d/bk_biz_id/%d/", setTemplateID, bizID)

	err := st.client.Put().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("UpdateSetTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return ret, nil
}

func (st *SetTemplate) DeleteSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.DeleteSetTemplateOption) errors.CCErrorCoder {
	ret := struct {
		metadata.BaseResp
	}{}
	subPath := fmt.Sprintf("/deletemany/topo/set_template/bk_biz_id/%d/", bizID)

	err := st.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("DeleteSetTemplate failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return nil
}

func (st *SetTemplate) GetSetTemplate(ctx context.Context, header http.Header, bizID int64, setTemplateID int64) (*metadata.SetTemplateResult, errors.CCErrorCoder) {
	ret := &metadata.SetTemplateResult{}
	subPath := fmt.Sprintf("/find/topo/set_template/%d/bk_biz_id/%d/", setTemplateID, bizID)

	err := st.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("GetSetTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}

	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return ret, nil
}

func (st *SetTemplate) ListSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.ListSetTemplateOption) (*metadata.MultipleSetTemplateResult, errors.CCErrorCoder) {
	ret := &metadata.ListSetTemplateResult{}
	subPath := fmt.Sprintf("/findmany/topo/set_template/bk_biz_id/%d/", bizID)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("ListSetTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}

	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (st *SetTemplate) ListSetTemplateWeb(ctx context.Context, header http.Header, bizID int64, option metadata.ListSetTemplateOption) (*metadata.MultipleSetTemplateResult, errors.CCErrorCoder) {
	ret := &metadata.ListSetTemplateResult{}
	subPath := fmt.Sprintf("/findmany/topo/set_template/bk_biz_id/%d/web", bizID)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("ListSetTemplateWeb failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}

	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (st *SetTemplate) ListSetTplRelatedSvcTpl(ctx context.Context, header http.Header, bizID int64, setTemplateID int64) ([]metadata.ServiceTemplate, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data []metadata.ServiceTemplate `json:"data"`
	}{}
	subPath := fmt.Sprintf("/findmany/topo/set_template/%d/bk_biz_id/%d/service_templates", setTemplateID, bizID)

	err := st.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("ListSetTplRelatedSvcTpl failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (st *SetTemplate) ListSetTplRelatedSetsWeb(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.ListSetByTemplateOption) (*metadata.InstDataInfo, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data metadata.InstDataInfo `json:"data"`
	}{}
	subPath := fmt.Sprintf("/findmany/topo/set_template/%d/bk_biz_id/%d/sets/web", setTemplateID, bizID)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("ListSetTplRelatedSetsWeb failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}

func (st *SetTemplate) DiffSetTplWithInst(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.DiffSetTplWithInstOption) (*metadata.SetTplDiffResult, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data metadata.SetTplDiffResult `json:"data"`
	}{}
	subPath := fmt.Sprintf("/findmany/topo/set_template/%d/bk_biz_id/%d/diff_with_instances", setTemplateID, bizID)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("DiffSetTplWithInst failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return &ret.Data, nil
}
