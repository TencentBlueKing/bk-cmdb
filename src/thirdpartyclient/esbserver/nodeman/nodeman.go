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
	"fmt"
	"net/http"

	"configcenter/src/thirdpartyclient/esbserver/esbutil"
)

func (p *nodeman) SearchPackage(ctx context.Context, h http.Header, processname string) (resp *SearchPluginPackageResult, err error) {
	resp = new(SearchPluginPackageResult)
	subPath := fmt.Sprintf("/%s/package/?%s", processname, esbutil.GetEsbQueryParameters(p.config.GetConfig(), h))
	err = p.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return nil, nil
}
func (p *nodeman) SearchProcess(ctx context.Context, h http.Header, processname string) (resp *SearchPluginProcessResult, err error) {
	resp = new(SearchPluginProcessResult)
	subPath := fmt.Sprintf("/process/%s/?%s", processname, esbutil.GetEsbQueryParameters(p.config.GetConfig(), h))
	err = p.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return resp, nil
}
func (p *nodeman) SearchProcessInfo(ctx context.Context, h http.Header, processname string) (resp *SearchPluginProcessInfoResult, err error) {
	resp = new(SearchPluginProcessInfoResult)
	subPath := fmt.Sprintf("/process_info/%s/?%s", processname, esbutil.GetEsbQueryParameters(p.config.GetConfig(), h))
	err = p.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return resp, nil
}
func (p *nodeman) UpgradePlugin(ctx context.Context, h http.Header, bizID string, data *UpgradePluginRequest) (resp *UpgradePluginResult, err error) {
	resp = new(UpgradePluginResult)
	subPath := fmt.Sprintf("/%s/tasks/?%s", bizID, esbutil.GetEsbQueryParameters(p.config.GetConfig(), h))

	params := struct {
		*esbutil.EsbCommParams
		*UpgradePluginRequest
	}{
		EsbCommParams:        esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		UpgradePluginRequest: data,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return resp, nil
}
func (p *nodeman) SearchTask(ctx context.Context, h http.Header, bizID string, taskID string) (resp *SearchTaskResult, err error) {
	resp = new(SearchTaskResult)
	subPath := fmt.Sprintf("/%s/tasks/%s/?%s", bizID, taskID, esbutil.GetEsbQueryParameters(p.config.GetConfig(), h))
	err = p.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return resp, nil
}
func (p *nodeman) SearchPluginHost(ctx context.Context, h http.Header, processname string) (resp *SearchPluginHostResult, err error) {
	resp = new(SearchPluginHostResult)
	subPath := fmt.Sprintf("/0/host_status/get_host/?name=%s&%s", processname, esbutil.GetEsbQueryParameters(p.config.GetConfig(), h))
	err = p.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return resp, nil
}
