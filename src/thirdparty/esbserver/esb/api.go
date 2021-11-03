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
package esb

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/esbserver/esbutil"
)

// GetApiPublicKey get api publicKey from esb
func (e *esb) GetApiPublicKey(ctx context.Context, h http.Header) (resp *metadata.EsbResponse, err error) {
	resp = new(metadata.EsbResponse)
	subPath := "/v2/esb/get_api_public_key/"
	params := esbutil.GetEsbRequestParams(e.config.GetConfig(), h)

	err = e.client.Get().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}
