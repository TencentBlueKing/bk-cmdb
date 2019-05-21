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
package synchronizeserver

import (
	"context"
	"net/http"

	//"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

func (s *synchronize) Find(ctx context.Context, h http.Header, input *metadata.SynchronizeFindInfoParameter) (resp *metadata.ResponseInstData, err error) {
	resp = new(metadata.ResponseInstData)
	subPath := "/search"

	err = s.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}
