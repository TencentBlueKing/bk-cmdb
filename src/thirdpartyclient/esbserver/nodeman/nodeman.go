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

	"configcenter/src/common/metadata"
	"configcenter/src/thirdpartyclient/esbserver/esbutil"
)

func (p *nodeman) OperateProcess(ctx context.Context, h http.Header, data *metadata.GseProcRequest) (resp *metadata.EsbResponse, err error) {
	resp = new(metadata.EsbResponse)
	subPath := "/v2/gse/operate_proc/"
	type esbParams struct {
		*esbutil.EsbCommParams
		*metadata.GseProcRequest `"json:inline"`
	}
	params := &esbParams{
		EsbCommParams:  esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseProcRequest: data,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (p *nodeman) QueryProcOperateResult(ctx context.Context, h http.Header, taskID string) (resp *metadata.GseProcessOperateTaskResult, err error) {
	resp = new(metadata.GseProcessOperateTaskResult)
	subPath := "/v2/gse/get_proc_operate_result/"
	type esbParams struct {
		*esbutil.EsbCommParams
		TaskID string `json:"task_id"`
	}
	params := &esbParams{
		EsbCommParams: esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		TaskID:        taskID,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (p *nodeman) QueryProcStatus(ctx context.Context, h http.Header, data *metadata.GseProcRequest) (resp *metadata.EsbResponse, err error) {
	resp = new(metadata.EsbResponse)
	subPath := "/v2/gse/get_proc_status/"
	type esbParams struct {
		*esbutil.EsbCommParams
		*metadata.GseProcRequest `"json:inline"`
	}
	params := &esbParams{
		EsbCommParams:  esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseProcRequest: data,
	}
	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (p *nodeman) RegisterProcInfo(ctx context.Context, h http.Header, data *metadata.GseProcRequest) (resp *metadata.EsbResponse, err error) {
	resp = new(metadata.EsbResponse)
	subPath := "/v2/gse/register_proc_info/"
	type esbParams struct {
		*esbutil.EsbCommParams
		*metadata.GseProcRequest `"json:inline"`
	}
	params := &esbParams{
		EsbCommParams:  esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseProcRequest: data,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (p *nodeman) UnRegisterProcInfo(ctx context.Context, h http.Header, data *metadata.GseProcRequest) (resp *metadata.EsbResponse, err error) {
	resp = new(metadata.EsbResponse)
	subPath := "/v2/gse/unregister_proc_info/"
	type esbParams struct {
		*esbutil.EsbCommParams
		*metadata.GseProcRequest `"json:inline"`
	}
	params := &esbParams{
		EsbCommParams:  esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseProcRequest: data,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

func (p *nodeman) SearchPackage(ctx context.Context, h http.Header, processname string) (resp *SearchPluginPackageResult, err error) {
	resp = new(SearchPluginPackageResult)
	subPath := "/api/c/self-service-api/nodeman_api/2/tasks/2"
	type esbParams struct {
		*esbutil.EsbCommParams
		*metadata.GseProcRequest `"json:inline"`
	}
	params := &esbParams{
		EsbCommParams:  esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseProcRequest: data,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return nil, nil
}
func (p *nodeman) SearchProcess(ctx context.Context, h http.Header, processname string) (resp *SearchPluginProcessResult, err error) {
	return nil, nil
}
func (p *nodeman) SearchProcessInfo(ctx context.Context, h http.Header, processname string) (resp *SearchPluginProcessInfoResult, err error) {
	return nil, nil
}
func (p *nodeman) UpgradePlugin(ctx context.Context, h http.Header, bizID string, data UpgradePluginRequest) (resp *UpgradePluginResult, err error) {
	return nil, nil
}
func (p *nodeman) SearchTask(ctx context.Context, h http.Header, bizID string, taskID string) (resp *SearchTaskResult, err error) {
	resp = new(SearchTaskResult)
	subPath := fmt.Sprintf("/api/c/self-service-api/nodeman_api/%s/tasks/%s", bizID, taskID)
	type esbParams struct {
		*esbutil.EsbCommParams
		*metadata.GseProcRequest `"json:inline"`
	}
	params := &esbParams{
		EsbCommParams:  esbutil.GetEsbRequestParams(p.config.GetConfig(), h),
		GseProcRequest: data,
	}

	err = p.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return nil, nil
}
func (p *nodeman) SearchPluginHost(ctx context.Context, h http.Header, bizID string, processname string) (resp *SearchPluginHostResult, err error) {
	return nil, nil
}
