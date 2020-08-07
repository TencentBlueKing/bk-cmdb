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

package cache

import (
	"context"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/cache/topo_tree"
)

func (b *baseCache) SearchTopologyTree(ctx context.Context, h http.Header, opt *topo_tree.SearchOption) ([]topo_tree.Topology, error) {
	type Topo struct {
		metadata.BaseResp `json:",inline"`
		Data              []topo_tree.Topology `json:"data"`
	}

	resp := new(Topo)

	err := b.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/find/cache/topotree").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !resp.Result {
		return nil, errors.New(resp.Code, resp.ErrMsg)
	}

	return resp.Data, nil
}

func (b *baseCache) SearchHostWithInnerIP(ctx context.Context, h http.Header, opt *metadata.SearchHostWithInnerIPOption) (jsonString string, err error) {

	resp, err := b.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/find/cache/host/with_inner_ip").
		WithHeaders(h).
		Do().
		IntoJsonString()

	if err != nil {
		return "", errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !resp.Result {
		return "", errors.New(resp.Code, resp.ErrMsg)
	}

	return resp.Data, nil
}

func (b *baseCache) SearchHostWithHostID(ctx context.Context, h http.Header, opt *metadata.SearchHostWithIDOption) (jsonString string, err error) {

	resp, err := b.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/find/cache/host/with_host_id").
		WithHeaders(h).
		Do().
		IntoJsonString()

	if err != nil {
		return "", errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !resp.Result {
		return "", errors.New(resp.Code, resp.ErrMsg)
	}

	return resp.Data, nil
}

func (b *baseCache) ListHostWithHostID(ctx context.Context, h http.Header, opt *metadata.ListHostWithIDOption) (jsonString string, err error) {

	resp, err := b.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/findmany/cache/host/with_host_id").
		WithHeaders(h).
		Do().
		IntoJsonString()

	if err != nil {
		return "", errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !resp.Result {
		return "", errors.New(resp.Code, resp.ErrMsg)
	}

	return resp.Data, nil
}

func (b *baseCache) SearchBusiness(ctx context.Context, h http.Header, bizID int64) (string, error) {
	resp, err := b.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef("/find/cache/biz/%d", bizID).
		WithHeaders(h).
		Do().
		IntoJsonString()
	if err != nil {
		return "", errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}
	if !resp.Result {
		return "", errors.New(resp.Code, resp.ErrMsg)
	}
	return resp.Data, nil
}

func (b *baseCache) SearchSet(ctx context.Context, h http.Header, setID int64) (string, error) {
	resp, err := b.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef("/find/cache/set/%d", setID).
		WithHeaders(h).
		Do().
		IntoJsonString()
	if err != nil {
		return "", errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}
	if !resp.Result {
		return "", errors.New(resp.Code, resp.ErrMsg)
	}
	return resp.Data, nil
}

func (b *baseCache) SearchModule(ctx context.Context, h http.Header, moduleID int64) (string, error) {
	resp, err := b.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef("/find/cache/module/%d", moduleID).
		WithHeaders(h).
		Do().
		IntoJsonString()
	if err != nil {
		return "", errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}
	if !resp.Result {
		return "", errors.New(resp.Code, resp.ErrMsg)
	}
	return resp.Data, nil
}

func (b *baseCache) SearchCustomLayer(ctx context.Context, h http.Header, objID string, instID int64) (string, error) {
	resp, err := b.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef("/find/cache/%s/%d", objID, instID).
		WithHeaders(h).
		Do().
		IntoJsonString()
	if err != nil {
		return "", errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}
	if !resp.Result {
		return "", errors.New(resp.Code, resp.ErrMsg)
	}
	return resp.Data, nil
}
