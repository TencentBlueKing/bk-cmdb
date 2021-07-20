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
	"net/http"

	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
)

func (e *eventServer) Query(ctx context.Context, ownerID string, appID string, h http.Header, dat metadata.ParamSubscriptionSearch) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/subscribe/search/%s/%s"

	err = e.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath, ownerID, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (e *eventServer) Ping(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/subscribe/ping"

	err = e.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (e *eventServer) Telnet(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/subscribe/telnet"

	err = e.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (e *eventServer) Subscribe(ctx context.Context, ownerID string, appID string, h http.Header, subscription *metadata.Subscription) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/subscribe/%s/%s"

	err = e.client.Post().
		WithContext(ctx).
		Body(subscription).
		SubResourcef(subPath, ownerID, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (e *eventServer) UnSubscribe(ctx context.Context, ownerID string, appID string, subscribeID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/subscribe/%s/%s/%s"

	err = e.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, ownerID, appID, subscribeID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (e *eventServer) Rebook(ctx context.Context, ownerID string, appID string, subscribeID string, h http.Header, subscription *metadata.Subscription) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/subscribe/%s/%s/%s"

	err = e.client.Put().
		WithContext(ctx).
		Body(subscription).
		SubResourcef(subPath, ownerID, appID, subscribeID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (e *eventServer) Watch(ctx context.Context, h http.Header, opts *watch.WatchEventOptions) (resp []*watch.WatchEventDetail, err error) {
	response := new(metadata.WatchEventResp)
	err = e.client.Post().
		WithContext(ctx).
		Body(opts).
		SubResourcef("/watch/resource/%s", opts.Resource).
		WithHeaders(h).
		Do().
		Into(response)

	if err != nil {
		return nil, err
	}

	if err = response.CCError(); err != nil {
		return nil, err
	}

	if response.Data == nil {
		return make([]*watch.WatchEventDetail, 0), nil
	}
	return response.Data.Events, nil
}
