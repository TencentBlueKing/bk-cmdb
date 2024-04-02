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

	"configcenter/src/thirdparty/apigw"
)

// ListAgentState list gse agent state
func (p *gse) ListAgentState(ctx context.Context, h http.Header, data *ListAgentStateRequest) (*ListAgentStateResp,
	error) {

	h.Set(apigw.AuthKey, p.service.Auth)
	resp := new(ListAgentStateResp)
	subPath := "/api/v2/cluster/list_agent_state"

	err := p.service.Client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
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

	h.Set(apigw.AuthKey, p.service.Auth)
	resp := new(AsyncPushFileResp)
	subPath := "/api/v2/task/async_push_file"

	err := p.service.Client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
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

	h.Set(apigw.AuthKey, p.service.Auth)
	resp := new(GetTransferFileResultResp)
	subPath := "/api/v2/task/async/get_transfer_file_result"

	err := p.service.Client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
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
