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
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/esbserver/esbutil"
)

type GseClientInterface interface {
	OperateProcess(ctx context.Context, h http.Header, data *metadata.GseProcRequest) (resp *metadata.EsbResponse, err error)
	QueryProcOperateResult(ctx context.Context, h http.Header, taskid string) (resp *metadata.GseProcessOperateTaskResult, err error)
	QueryProcStatus(ctx context.Context, h http.Header, data *metadata.GseProcRequest) (resp *metadata.EsbResponse, err error)
	RegisterProcInfo(ctx context.Context, h http.Header, data *metadata.GseProcRequest) (resp *metadata.EsbResponse, err error)
	UnRegisterProcInfo(ctx context.Context, h http.Header, data *metadata.GseProcRequest) (resp *metadata.EsbResponse, err error)
	ConfigAddStreamTo(ctx context.Context, h http.Header, data *metadata.GseConfigAddStreamToParams) (
		*metadata.GseConfigAddStreamToResult, error)
	ConfigUpdateStreamTo(ctx context.Context, h http.Header, data *metadata.GseConfigUpdateStreamToParams) error
	ConfigDeleteStreamTo(ctx context.Context, h http.Header, data *metadata.GseConfigDeleteStreamToParams) error
	ConfigQueryStreamTo(ctx context.Context, h http.Header, data *metadata.GseConfigQueryStreamToParams) (
		[]metadata.GseConfigAddStreamToParams, error)
	ConfigAddRoute(ctx context.Context, h http.Header, data *metadata.GseConfigAddRouteParams) (
		resp *metadata.GseConfigAddRouteResult, err error)
	ConfigUpdateRoute(ctx context.Context, h http.Header, data *metadata.GseConfigUpdateRouteParams) error
	ConfigDeleteRoute(ctx context.Context, h http.Header, data *metadata.GseConfigDeleteRouteParams) error
	ConfigQueryRoute(ctx context.Context, h http.Header, data *metadata.GseConfigQueryRouteParams) (
		[]metadata.GseConfigChannel, error)
}

func NewGsecClientInterface(client rest.ClientInterface, config *esbutil.EsbConfigSrv) GseClientInterface {
	return &gse{
		client: client,
		config: config,
	}
}

type gse struct {
	config *esbutil.EsbConfigSrv
	client rest.ClientInterface
}

type esbGseProcParams struct {
	*esbutil.EsbCommParams
	*metadata.GseProcRequest `json:",inline"`
}

type esbTaskIDParams struct {
	*esbutil.EsbCommParams
	TaskID string `json:"task_id"`
}

type esbGseConfigAddStreamToParams struct {
	*esbutil.EsbCommParams
	*metadata.GseConfigAddStreamToParams `json:",inline"`
}

type esbGseConfigUpdateStreamToParams struct {
	*esbutil.EsbCommParams
	*metadata.GseConfigUpdateStreamToParams `json:",inline"`
}

type esbGseConfigDeleteStreamToParams struct {
	*esbutil.EsbCommParams
	*metadata.GseConfigDeleteStreamToParams `json:",inline"`
}

type esbGseConfigQueryStreamToParams struct {
	*esbutil.EsbCommParams
	*metadata.GseConfigQueryStreamToParams `json:",inline"`
}

type esbGseConfigAddRouteParams struct {
	*esbutil.EsbCommParams
	*metadata.GseConfigAddRouteParams `json:",inline"`
}

type esbGseConfigUpdateRouteParams struct {
	*esbutil.EsbCommParams
	*metadata.GseConfigUpdateRouteParams `json:",inline"`
}

type esbGseConfigDeleteRouteParams struct {
	*esbutil.EsbCommParams
	*metadata.GseConfigDeleteRouteParams `json:",inline"`
}

type esbGseConfigQueryRouteParams struct {
	*esbutil.EsbCommParams
	*metadata.GseConfigQueryRouteParams `json:",inline"`
}
