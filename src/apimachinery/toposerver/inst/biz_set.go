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

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// CreateBizSet TODO
func (t *instanceClient) CreateBizSet(ctx context.Context, h http.Header, opt metadata.CreateBizSetRequest) (
	int64, errors.CCErrorCoder) {

	resp := new(metadata.CreateBizSetResponse)
	subPath := "/create/biz_set"

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return 0, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return 0, err
	}

	return resp.Data, nil
}

// UpdateBizSet TODO
func (t *instanceClient) UpdateBizSet(ctx context.Context, h http.Header,
	opt metadata.UpdateBizSetOption) errors.CCErrorCoder {

	resp := new(metadata.Response)
	subPath := "/updatemany/biz_set"

	err := t.client.Put().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	return resp.CCError()
}

// DeleteBizSet TODO
func (t *instanceClient) DeleteBizSet(ctx context.Context, h http.Header,
	opt metadata.DeleteBizSetOption) errors.CCErrorCoder {

	resp := new(metadata.Response)
	subPath := "/deletemany/biz_set"

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	return resp.CCError()
}

// FindBizInBizSet TODO
func (t *instanceClient) FindBizInBizSet(ctx context.Context, h http.Header, opt *metadata.FindBizInBizSetOption) (
	*metadata.InstResult, errors.CCErrorCoder) {

	resp := new(metadata.QueryInstResult)
	subPath := "/find/biz_set/biz_list"

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
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

// FindBizSetTopo TODO
func (t *instanceClient) FindBizSetTopo(ctx context.Context, h http.Header, opt *metadata.FindBizSetTopoOption) (
	[]mapstr.MapStr, errors.CCErrorCoder) {

	resp := new(metadata.MapArrayResponse)
	subPath := "/find/biz_set/topo_path"

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
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

// SearchBusinessSet TODO
func (t *instanceClient) SearchBusinessSet(ctx context.Context, h http.Header, opt *metadata.QueryBusinessSetRequest) (
	*metadata.InstResult, errors.CCErrorCoder) {

	resp := new(metadata.QueryInstResult)
	subPath := "/findmany/biz_set"

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
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

// SearchModuleInBizSet search modules in one biz set
func (t *instanceClient) SearchModuleInBizSet(ctx context.Context, bizSetID int64, bizID int64, setID int64,
	h http.Header) (*metadata.InstResult, errors.CCErrorCoder) {

	resp := new(metadata.QueryInstResult)
	subPath := "/findmany/module/biz_set/%d/biz/%d/set/%d"

	err := t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, bizSetID, bizID, setID).
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

// SearchBizSetTopoPath search biz set topo path
func (t *instanceClient) SearchBizSetTopoPath(ctx context.Context, bizSetID int64, bizID int64, h http.Header,
	data mapstr.MapStr) (*metadata.TopoPathResult, errors.CCErrorCoder) {

	resp := new(metadata.TopoPathResp)
	subPath := "/find/topopath/biz_set/%d/biz/%d"

	err := t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizSetID, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	if err != nil {
		blog.Errorf("search topo path failed, err: %v", err)
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}
