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

package eventserver

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/core/cc/api"
	paraparse "configcenter/src/common/paraparse"
	"configcenter/src/scene_server/event_server/types"
)

func (e *eventServer) Query(ctx context.Context, ownerID string, appID string, h http.Header, dat paraparse.SubscribeCommonSearch) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/subscribe/search/%s/%s", ownerID, appID)

	err = e.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (e *eventServer) Ping(ctx context.Context, h http.Header, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/subscribe/ping"

	err = e.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (e *eventServer) Telnet(ctx context.Context, h http.Header, dat interface{}) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/subscribe/telnet"

	err = e.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (e *eventServer) Subscribe(ctx context.Context, ownerID string, appID string, h http.Header, subscription *types.Subscription) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/subscribe/%s/%s", ownerID, appID)

	err = e.client.Post().
		WithContext(ctx).
		Body(subscription).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (e *eventServer) UnSubscribe(ctx context.Context, ownerID string, appID string, subscribeID string, h http.Header) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/subscribe/%s/%s/%s", ownerID, appID, subscribeID)

	err = e.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (e *eventServer) Rebook(ctx context.Context, ownerID string, appID string, subscribeID string, h http.Header, subscription *types.Subscription) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/subscribe/%s/%s/%s", ownerID, appID, subscribeID)

	err = e.client.Put().
		WithContext(ctx).
		Body(subscription).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
