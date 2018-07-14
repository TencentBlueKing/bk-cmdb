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

package host

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

func (t *hostctrl) GetHostByID(ctx context.Context, hostID string, h http.Header) (resp *metadata.HostInstanceResult, err error) {
	resp = new(metadata.HostInstanceResult)
	subPath := fmt.Sprintf("/host/%s", hostID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *hostctrl) GetHosts(ctx context.Context, h http.Header, opt *metadata.QueryInput) (resp *metadata.GetHostsResult, err error) {
	resp = new(metadata.GetHostsResult)
	subPath := "/hosts/search"

	err = t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *hostctrl) AddHost(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/insts"

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *hostctrl) GetHostSnap(ctx context.Context, hostID string, h http.Header) (resp *metadata.GetHostSnapResult, err error) {
	resp = new(metadata.GetHostSnapResult)
	subPath := fmt.Sprintf("/host/snapshot/%s", hostID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
