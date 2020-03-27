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

package inst

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

func (t *instanceClient) QueryAudit(ctx context.Context, ownerID string, h http.Header, input *metadata.QueryInput) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/app/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
func (t *instanceClient) QueryAuditLog(ctx context.Context, h http.Header, input *metadata.QueryInput) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/audit/search"

	err = t.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) GetInternalModule(ctx context.Context, ownerID, appID string, h http.Header) (resp *metadata.SearchInnterAppTopoResult, err error) {
	resp = new(metadata.SearchInnterAppTopoResult)
	subPath := "/topo/internal/%s/%s"

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, ownerID, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
