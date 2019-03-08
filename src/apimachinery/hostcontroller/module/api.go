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

	"configcenter/src/common/metadata"
)

func (m *mod) GetHostModulesIDs(ctx context.Context, h http.Header, dat *metadata.ModuleHostConfigParams) (resp *metadata.GetHostModuleIDsResult, err error) {
	resp = new(metadata.GetHostModuleIDsResult)
	subPath := "/meta/hosts/modules/search"

	err = m.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *mod) AddModuleHostConfig(ctx context.Context, h http.Header, dat *metadata.ModuleHostConfigParams) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := "/meta/hosts/modules"

	err = m.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *mod) DelModuleHostConfig(ctx context.Context, h http.Header, dat *metadata.ModuleHostConfigParams) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := "/meta/hosts/modules"

	err = m.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *mod) DelDefaultModuleHostConfig(ctx context.Context, h http.Header, dat *metadata.ModuleHostConfigParams) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := "/meta/hosts/defaultmodules"

	err = m.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *mod) MoveHost2ResourcePool(ctx context.Context, h http.Header, dat *metadata.ParamData) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := "/meta/hosts/resource"

	err = m.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *mod) AssignHostToApp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := "/meta/hosts/assign"

	err = m.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (m *mod) GetModulesHostConfig(ctx context.Context, h http.Header, dat map[string][]int64) (resp *metadata.HostConfig, err error) {
	resp = new(metadata.HostConfig)
	subPath := "/meta/hosts/module/config/search"

	err = m.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
