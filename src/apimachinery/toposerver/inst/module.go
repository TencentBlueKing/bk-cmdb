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

package inst

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
)

// CreateModule TODO
func (t *instanceClient) CreateModule(ctx context.Context, appID, setID int64, h http.Header,
	dat map[string]interface{}) (mapstr.MapStr, errors.CCErrorCoder) {

	resp := new(metadata.CreateInstResult)
	subPath := "/module/%d/%d"

	err := t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath, appID, setID).
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

// DeleteModule TODO
func (t *instanceClient) DeleteModule(ctx context.Context, appID, setID, moduleID int64,
	h http.Header) errors.CCErrorCoder {

	resp := new(metadata.Response)
	subPath := "/module/%d/%d/%d"

	err := t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, appID, setID, moduleID).
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

// UpdateModule TODO
func (t *instanceClient) UpdateModule(ctx context.Context, appID, setID, moduleID int64, h http.Header,
	dat map[string]interface{}) errors.CCErrorCoder {

	resp := new(metadata.Response)
	subPath := "/module/%d/%d/%d"

	err := t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath, appID, setID, moduleID).
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

// SearchModule TODO
func (t *instanceClient) SearchModule(ctx context.Context, ownerID string, appID, setID int64, h http.Header,
	s *params.SearchParams) (*metadata.InstResult, errors.CCErrorCoder) {

	resp := new(metadata.SearchInstResult)
	subPath := "/module/search/%s/%d/%d"

	err := t.client.Post().
		WithContext(ctx).
		Body(s).
		SubResourcef(subPath, ownerID, appID, setID).
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

// SearchModuleByCondition TODO
func (t *instanceClient) SearchModuleByCondition(ctx context.Context, appID string, h http.Header, s *params.SearchParams) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := "/findmany/module/biz/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(s).
		SubResourcef(subPath, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchModuleBatch TODO
func (t *instanceClient) SearchModuleBatch(ctx context.Context, appID string, h http.Header, s *metadata.SearchInstBatchOption) (resp *metadata.MapArrayResponse, err error) {
	resp = new(metadata.MapArrayResponse)
	subPath := "/findmany/module/bk_biz_id/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(s).
		SubResourcef(subPath, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchModuleWithRelation TODO
func (t *instanceClient) SearchModuleWithRelation(ctx context.Context, appID string, h http.Header, dat map[string]interface{}) (resp *metadata.ResponseInstData, err error) {
	resp = new(metadata.ResponseInstData)
	subPath := "/findmany/module/with_relation/biz/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
