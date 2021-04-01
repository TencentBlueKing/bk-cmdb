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

func (p *gse) OperateProcess(ctx context.Context, h http.Header, data *metadata.GseProcRequest) (resp *metadata.EsbResponse, err error) {
	resp = new(metadata.EsbResponse)
	subPath := "/v2/gse/operate_proc/"
	params := &esbGseProcParams{
		EsbCommParams:  esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseProcRequest: data,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (p *gse) QueryProcOperateResult(ctx context.Context, h http.Header, taskID string) (resp *metadata.GseProcessOperateTaskResult, err error) {
	resp = new(metadata.GseProcessOperateTaskResult)
	subPath := "/v2/gse/get_proc_operate_result/"
	params := &esbTaskIDParams{
		EsbCommParams: esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		TaskID:        taskID,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (p *gse) QueryProcStatus(ctx context.Context, h http.Header, data *metadata.GseProcRequest) (resp *metadata.EsbResponse, err error) {
	resp = new(metadata.EsbResponse)
	subPath := "/v2/gse/get_proc_status/"
	params := &esbGseProcParams{
		EsbCommParams:  esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseProcRequest: data,
	}
	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (p *gse) RegisterProcInfo(ctx context.Context, h http.Header, data *metadata.GseProcRequest) (resp *metadata.EsbResponse, err error) {
	resp = new(metadata.EsbResponse)
	subPath := "/v2/gse/register_proc_info/"
	params := &esbGseProcParams{
		EsbCommParams:  esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseProcRequest: data,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (p *gse) UnRegisterProcInfo(ctx context.Context, h http.Header, data *metadata.GseProcRequest) (resp *metadata.EsbResponse, err error) {
	resp = new(metadata.EsbResponse)
	subPath := "/v2/gse/unregister_proc_info/"
	params := &esbGseProcParams{
		EsbCommParams:  esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseProcRequest: data,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

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
		WithHeaders(h).
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

func (p *gse) ConfigUpdateStreamTo(ctx context.Context, h http.Header, data *metadata.GseConfigUpdateStreamToParams) error {
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
		WithHeaders(h).
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

func (p *gse) ConfigDeleteStreamTo(ctx context.Context, h http.Header, data *metadata.GseConfigDeleteStreamToParams) error {
	resp := new(metadata.EsbBaseResponse)
	subPath := "/v2/gse/config_delete_streamto/"
	params := &esbGseConfigDeleteStreamToParams{
		EsbCommParams:                 esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseConfigDeleteStreamToParams: data,
	}

	err := p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return err
	}
	if !resp.Result || resp.Code != 0 {
		return fmt.Errorf("gse config delete streamto failed, code: %d, message: %s", resp.Code, resp.Message)
	}
	return nil
}

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
		WithHeaders(h).
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
		WithHeaders(h).
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
		WithHeaders(h).
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

func (p *gse) ConfigDeleteRoute(ctx context.Context, h http.Header, data *metadata.GseConfigDeleteRouteParams) error {
	resp := new(metadata.EsbBaseResponse)
	subPath := "/v2/gse/config_delete_route/"
	params := &esbGseConfigDeleteRouteParams{
		EsbCommParams:              esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseConfigDeleteRouteParams: data,
	}

	err := p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return err
	}
	if !resp.Result || resp.Code != 0 {
		return fmt.Errorf("gse config delete route failed, code: %d, message: %s", resp.Code, resp.Message)
	}
	return nil
}

func (p *gse) ConfigQueryRoute(ctx context.Context, h http.Header, data *metadata.GseConfigQueryRouteParams) (
	[]metadata.GseConfigChannel, error) {

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
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}
	if !resp.Result || resp.Code != 0 {
		return nil, fmt.Errorf("gse config query route failed, code: %d, message: %s", resp.Code, resp.Message)
	}
	return resp.Data, nil
}
