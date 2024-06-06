// Package gse TODO
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
package gse

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/esbserver/esbutil"
)

// ConfigAddStreamTo TODO
func (p *gse) ConfigAddStreamTo(ctx context.Context, h http.Header, data *metadata.GseConfigAddStreamToParams) (
	*metadata.GseConfigAddStreamToResult, error) {

	resp := new(metadata.GseConfigAddStreamToResp)
	subPath := "/v2/gse/config_add_streamto/"
	params := &esbGseConfigAddStreamToParams{
		EsbCommParams:              esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseConfigAddStreamToParams: data,
	}

	err := p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}
	if !resp.Result || resp.Code != 0 {
		return nil, fmt.Errorf("gse config add streamto failed, code: %d, message: %s", resp.Code, resp.Message)
	}
	return resp.Data, nil
}

// ConfigUpdateStreamTo TODO
func (p *gse) ConfigUpdateStreamTo(ctx context.Context, h http.Header,
	data *metadata.GseConfigUpdateStreamToParams) error {
	resp := new(metadata.EsbBaseResponse)
	subPath := "/v2/gse/config_update_streamto/"
	params := &esbGseConfigUpdateStreamToParams{
		EsbCommParams:                 esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseConfigUpdateStreamToParams: data,
	}

	err := p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Do().
		Into(resp)

	if err != nil {
		return err
	}
	if !resp.Result || resp.Code != 0 {
		return fmt.Errorf("gse config update streamto failed, code: %d, message: %s", resp.Code, resp.Message)
	}
	return nil
}

// ConfigQueryStreamTo TODO
func (p *gse) ConfigQueryStreamTo(ctx context.Context, h http.Header, data *metadata.GseConfigQueryStreamToParams) (
	[]metadata.GseConfigAddStreamToParams, error) {

	resp := new(metadata.GseConfigQueryStreamToResp)
	subPath := "/v2/gse/config_query_streamto/"
	params := &esbGseConfigQueryStreamToParams{
		EsbCommParams:                esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseConfigQueryStreamToParams: data,
	}

	err := p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}
	if !resp.Result || resp.Code != 0 {
		return nil, fmt.Errorf("gse config query streamto failed, code: %d, message: %s", resp.Code, resp.Message)
	}
	return resp.Data, nil
}

// ConfigAddRoute TODO
func (p *gse) ConfigAddRoute(ctx context.Context, h http.Header, data *metadata.GseConfigAddRouteParams) (
	*metadata.GseConfigAddRouteResult, error) {

	resp := new(metadata.GseConfigAddRouteResp)
	subPath := "/v2/gse/config_add_route/"
	params := &esbGseConfigAddRouteParams{
		EsbCommParams:           esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseConfigAddRouteParams: data,
	}

	err := p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}
	if !resp.Result || resp.Code != 0 {
		return nil, fmt.Errorf("gse config add route failed, code: %d, message: %s", resp.Code, resp.Message)
	}
	return resp.Data, nil
}

// ConfigUpdateRoute TODO
func (p *gse) ConfigUpdateRoute(ctx context.Context, h http.Header, data *metadata.GseConfigUpdateRouteParams) error {
	resp := new(metadata.EsbBaseResponse)
	subPath := "/v2/gse/config_update_route/"
	params := &esbGseConfigUpdateRouteParams{
		EsbCommParams:              esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseConfigUpdateRouteParams: data,
	}

	err := p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Do().
		Into(resp)

	if err != nil {
		return err
	}
	if !resp.Result || resp.Code != 0 {
		return fmt.Errorf("gse config update route failed, code: %d, message: %s", resp.Code, resp.Message)
	}
	return nil
}

// ConfigQueryRoute TODO
func (p *gse) ConfigQueryRoute(ctx context.Context, h http.Header, data *metadata.GseConfigQueryRouteParams) (
	[]metadata.GseConfigChannel, bool, error) {

	resp := new(metadata.GseConfigQueryRouteResp)
	subPath := "/v2/gse/config_query_route/"
	params := &esbGseConfigQueryRouteParams{
		EsbCommParams:             esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseConfigQueryRouteParams: data,
	}

	err := p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(esbutil.SetEsbAuthHeader(p.config.GetConfig(), h)).
		Do().
		Into(resp)

	if err != nil {
		return nil, false, err
	}

	// special error code for route not exists
	if resp.Code == 14001 || resp.Code == 1014003 || resp.Code == 1014505 {
		return nil, false, nil
	}

	if !resp.Result || resp.Code != 0 {
		return nil, false, fmt.Errorf("gse config query route failed, code: %d, message: %s", resp.Code, resp.Message)
	}
	return resp.Data, true, nil
}
