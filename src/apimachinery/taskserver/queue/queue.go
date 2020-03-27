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
package queue

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

func (tq *taskQueue) Post(ctx context.Context, header http.Header, path string, data interface{}) (resp *metadata.Response, err error) {
	resp = &metadata.Response{}

	err = tq.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(path).
		WithHeaders(header).
		Do().
		Into(resp)

	return
}
