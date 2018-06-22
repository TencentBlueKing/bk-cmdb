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

package module

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type ModuleInterface interface {
	GetHostModulesIDs(ctx context.Context, h http.Header, dat *metadata.ModuleHostConfigParams) (resp *metadata.GetHostModuleIDsResult, err error)
	AddModuleHostConfig(ctx context.Context, h http.Header, dat *metadata.ModuleHostConfigParams) (resp *metadata.BaseResp, err error)
	DelModuleHostConfig(ctx context.Context, h http.Header, dat *metadata.ModuleHostConfigParams) (resp *metadata.BaseResp, err error)
	DelDefaultModuleHostConfig(ctx context.Context, h http.Header, dat *metadata.ModuleHostConfigParams) (resp *metadata.BaseResp, err error)
	MoveHost2ResourcePool(ctx context.Context, h http.Header, dat *metadata.ParamData) (resp *metadata.BaseResp, err error)
	AssignHostToApp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.BaseResp, err error)
	GetModulesHostConfig(ctx context.Context, h http.Header, dat map[string][]int64) (resp *metadata.HostConfig, err error)
}

func NewModuleInterface(client rest.ClientInterface) ModuleInterface {
	return &mod{client: client}
}

type mod struct {
	client rest.ClientInterface
}
