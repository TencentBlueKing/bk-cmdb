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

// ListFieldTemplateUnique list field template uniques
func (t template) ListFieldTemplateUnique(ctx context.Context, h http.Header, opt *metadata.CommonQueryOption) (
	*metadata.FieldTemplateUniqueInfo, errors.CCErrorCoder) {

	resp := new(metadata.ListFieldTmplUniqueResp)

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/findmany/field_template/unique").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// CreateFieldTemplateUniques create field template uniques
func (t template) CreateFieldTemplateUniques(ctx context.Context, h http.Header, templateID int64,
	opt []metadata.FieldTemplateUnique) (*metadata.RspIDs, errors.CCErrorCoder) {

	resp := new(metadata.CreateBatchResult)

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/createmany/field_template/%d/unique", templateID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// DeleteFieldTemplateUnique delete field template unique
func (t template) DeleteFieldTemplateUnique(ctx context.Context, h http.Header, templateID int64,
	opt *metadata.DeleteOption) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)

	err := t.client.Delete().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/delete/field_template/%d/uniques", templateID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// UpdateFieldTemplateUniques update field template uniques
func (t template) UpdateFieldTemplateUniques(ctx context.Context, h http.Header, templateID int64,
	opt []metadata.FieldTemplateUnique) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)

	err := t.client.Put().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/update/field_template/%d/uniques", templateID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}
