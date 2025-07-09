/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package eventserver

import (
	"context"
	"net/http"

	"configcenter/src/common/watch"
)

// Watch event
func (e *eventServer) Watch(ctx context.Context, h http.Header,
	opts *watch.WatchEventOptions) (resp []*watch.WatchEventDetail, err error) {
	response := new(watch.WatchEventResp)
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
