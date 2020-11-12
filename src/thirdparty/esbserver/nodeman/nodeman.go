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
package nodeman

import (
	"context"
	"net/http"

	"configcenter/src/thirdparty/esbserver/esbutil"
)

func (p *nodeman) SearchPackage(ctx context.Context, h http.Header, processname string) (resp *SearchPluginPackageResult, err error) {
	resp = new(SearchPluginPackageResult)
	subPath := "/v2/nodeman/%s/package/"
	err = p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, processname).
		WithParams(esbutil.GetEsbQueryParameters(p.config.GetConfig(), h)).
		WithHeaders(h).
		Peek().
		Do().
		Into(resp)
	return
}

func (p *nodeman) SearchProcess(ctx context.Context, h http.Header, processname string) (resp *SearchPluginProcessResult, err error) {
	resp = new(SearchPluginProcessResult)
	subPath := "/v2/nodeman/process/%s/"
	err = p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, processname).
		WithParams(esbutil.GetEsbQueryParameters(p.config.GetConfig(), h)).
		WithHeaders(h).
		Peek().
		Do().
		Into(resp)
	return
}

func (p *nodeman) SearchProcessInfo(ctx context.Context, h http.Header, processname string) (resp *SearchPluginProcessInfoResult, err error) {
	resp = new(SearchPluginProcessInfoResult)
	subPath := "/v2/nodeman/process_info/%s/"
	err = p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, processname).
		WithParams(esbutil.GetEsbQueryParameters(p.config.GetConfig(), h)).
		WithHeaders(h).
		Peek().
		Do().
		Into(resp)
	return
}

func (p *nodeman) UpgradePlugin(ctx context.Context, h http.Header, bizID string, data *UpgradePluginRequest) (resp *UpgradePluginResult, err error) {
	resp = new(UpgradePluginResult)
	subPath := "/v2/nodeman/%s/tasks/"

	params := esbUpgradePluginParams{
		EsbCommParams:        esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		UpgradePluginRequest: data,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Peek().
		Do().
		Into(resp)
	return
}

func (p *nodeman) SearchTask(ctx context.Context, h http.Header, bizID int64, taskID int64) (resp *SearchTaskResult, err error) {
	resp = new(SearchTaskResult)
	subPath := "/v2/nodeman/%d/tasks/%d/"
	err = p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, bizID, taskID).
		WithParams(esbutil.GetEsbQueryParameters(p.config.GetConfig(), h)).
		WithHeaders(h).
		Peek().
		Do().
		Into(resp)
	return
}

func (p *nodeman) SearchPluginHost(ctx context.Context, h http.Header, processname string) (resp *SearchPluginHostResult, err error) {
	resp = new(SearchPluginHostResult)
	subPath := "/v2/nodeman/0/host_status/get_host/"
	err = p.client.Get().
		WithContext(ctx).
		SubResourcef(subPath).
		WithParams(esbutil.GetEsbQueryParameters(p.config.GetConfig(), h)).
		WithParam("name", processname).
		WithHeaders(h).
		Peek().
		Do().
		Into(resp)
	return
}
