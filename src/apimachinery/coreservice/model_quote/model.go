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

package modelquote

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// ListModelQuoteRelation list model quote relationships.
func (q quote) ListModelQuoteRelation(ctx context.Context, h http.Header, req *metadata.CommonQueryOption) (
	*metadata.ListModelQuoteRelRes, errors.CCErrorCoder) {

	resp := new(metadata.ListModelQuoteRelResp)

	err := q.client.Post().
		WithContext(ctx).
		Body(req).
		SubResourcef("/list/model/quote/relation").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// CreateModelQuoteRelation create model quote relationships.
func (q quote) CreateModelQuoteRelation(ctx context.Context, h http.Header,
	data []metadata.ModelQuoteRelation) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)

	err := q.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef("/createmany/model/quote/relation").
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

// DeleteModelQuoteRelation delete model quote relationships.
func (q quote) DeleteModelQuoteRelation(ctx context.Context, h http.Header,
	req *metadata.CommonFilterOption) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)

	err := q.client.Delete().
		WithContext(ctx).
		Body(req).
		SubResourcef("/deletemany/model/quote/relation").
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
