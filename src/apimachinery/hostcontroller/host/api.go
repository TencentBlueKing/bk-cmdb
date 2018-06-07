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

	"configcenter/src/apimachinery/util"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/common/commondata"
)

func (t *hostctrl) GetHostByID(ctx context.Context, hostID string, h util.Headers) (resp *metadata.Response, err error) {
	subPath := fmt.Sprintf("/host/%s", hostID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *hostctrl) GetHosts(ctx context.Context, h util.Headers, opt *commondata.ObjQueryInput) (resp *metadata.GetHostsResult, err error) {
	subPath := "/hosts/search"

	err = t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *hostctrl) AddHost(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/insts"

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}

func (t *hostctrl) GetHostSnap(ctx context.Context, hostID string, h util.Headers) (resp *metadata.GetHostSnapResult, err error) {
	subPath := fmt.Sprintf("/host/snapshot/%s", hostID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h.ToHeader()).
		Do().
		Into(resp)
	return
}
