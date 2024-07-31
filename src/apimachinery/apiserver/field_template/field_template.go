/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package fieldtmpl

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// ListFieldTemplate list field templates
func (t template) ListFieldTemplate(ctx context.Context, h http.Header, opt *metadata.CommonQueryOption) (
	*metadata.FieldTemplateInfo, errors.CCErrorCoder) {

	resp := new(metadata.ListFieldTemplateResp)

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/findmany/field_template").
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// ListObjFieldTmplRel list field template and object relations.
func (t template) ListObjFieldTmplRel(ctx context.Context, h http.Header, opt *metadata.ListObjFieldTmplRelOption) (
	*metadata.ObjFieldTmplRelInfo, errors.CCErrorCoder) {

	resp := new(metadata.ListObjFieldTmplRelResp)

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/findmany/field_template/object/relation").
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// CountFieldTemplateObj count field templates related objects
func (t template) CountFieldTemplateObj(ctx context.Context, h http.Header, opt *metadata.CountFieldTmplResOption) (
	[]metadata.FieldTmplResCount, errors.CCErrorCoder) {

	resp := new(metadata.CountFieldTemplateObjResp)

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/findmany/field_template/object/count").
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
