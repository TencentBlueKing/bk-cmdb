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

package gse

import (
	"context"
	"fmt"
	"net/http"

	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/apigw/apigwutil"
)

// ListAgentState list gse agent state
func (p *gse) ListAgentState(ctx context.Context, h http.Header, data *ListAgentStateRequest) (*ListAgentStateResp,
	error) {

	resp := new(ListAgentStateResp)
	subPath := "/api/v2/cluster/list_agent_state"

	err := p.service.Client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(httpheader.SetBkAuth(h, p.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp, nil
}

// AsyncPushFile push file to target host
func (p *gse) AsyncPushFile(ctx context.Context, h http.Header, data *AsyncPushFileRequest) (*AsyncPushFileResp,
	error) {

	resp := new(AsyncPushFileResp)
	subPath := "/api/v2/task/async_push_file"

	err := p.service.Client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(httpheader.SetBkAuth(h, p.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp, nil
}

// GetTransferFileResult get transfer file result
func (p *gse) GetTransferFileResult(ctx context.Context, h http.Header, data *GetTransferFileResultRequest) (
	*GetTransferFileResultResp, error) {

	resp := new(GetTransferFileResultResp)
	subPath := "/api/v2/task/async/get_transfer_file_result"

	err := p.service.Client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(httpheader.SetBkAuth(h, p.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp, nil
}

// ConfigAddStreamTo 添加数据路由入库配置
func (p *gse) ConfigAddStreamTo(ctx context.Context, h http.Header, data *metadata.GseConfigAddStreamToParams) (
	*metadata.GseConfigAddStreamToResult, error) {

	resp := new(AddStreamToResp)
	subPath := "/api/v2/data/add_streamto/"

	err := p.service.Client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(httpheader.SetBkAuth(h, p.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("gse config add streamto failed, code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data, nil
}

// ConfigUpdateStreamTo 修改数据入库配置信息
func (p *gse) ConfigUpdateStreamTo(ctx context.Context, h http.Header,
	data *metadata.GseConfigUpdateStreamToParams) error {

	resp := new(apigwutil.ApiGWBaseResponse)
	subPath := "/api/v2/data/update_streamto/"

	err := p.service.Client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(httpheader.SetBkAuth(h, p.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return fmt.Errorf("gse config update streamto failed, code: %d, message: %s", resp.Code, resp.Message)
	}

	return nil
}

// ConfigQueryStreamTo 查询数路由入库的配置
func (p *gse) ConfigQueryStreamTo(ctx context.Context, h http.Header, data *metadata.GseConfigQueryStreamToParams) (
	[]metadata.GseConfigAddStreamToParams, error) {

	resp := new(QueryStreamToResp)
	subPath := "/api/v2/data/query_streamto/"

	err := p.service.Client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(httpheader.SetBkAuth(h, p.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	// special error code for streamTo not exists
	if resp.Code == 14001 || resp.Code == 1014003 || resp.Code == 1014505 {
		return make([]metadata.GseConfigAddStreamToParams, 0), nil
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("gse config query streamto failed, code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data, nil
}

// ConfigAddRoute 添加数据路由
func (p *gse) ConfigAddRoute(ctx context.Context, h http.Header, data *metadata.GseConfigAddRouteParams) (
	*metadata.GseConfigAddRouteResult, error) {

	resp := new(AddRouteResp)
	subPath := "/api/v2/data/add_route/"

	err := p.service.Client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(httpheader.SetBkAuth(h, p.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("gse config add route failed, code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data, nil
}

// ConfigUpdateRoute 更新路由配置
func (p *gse) ConfigUpdateRoute(ctx context.Context, h http.Header, data *metadata.GseConfigUpdateRouteParams) error {

	resp := new(apigwutil.ApiGWBaseResponse)
	subPath := "/api/v2/data/update_route/"

	err := p.service.Client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(httpheader.SetBkAuth(h, p.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return fmt.Errorf("gse config update route failed, code: %d, message: %s", resp.Code, resp.Message)
	}

	return nil
}

// ConfigQueryRoute 查询数据路由配置信息
func (p *gse) ConfigQueryRoute(ctx context.Context, h http.Header, data *metadata.GseConfigQueryRouteParams) (
	[]metadata.GseConfigChannel, bool, error) {

	resp := new(QueryRouteResp)
	subPath := "/api/v2/data/query_route/"

	err := p.service.Client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(httpheader.SetBkAuth(h, p.service.Auth)).
		Do().
		Into(resp)

	if err != nil {
		return nil, false, err
	}

	// special error code for route not exists
	if resp.Code == 14001 || resp.Code == 1014003 || resp.Code == 1014505 {
		return nil, false, nil
	}

	if resp.Code != 0 {
		return nil, false, fmt.Errorf("gse config query route failed, code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data, true, nil
}
