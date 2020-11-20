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

package secrets

import (
	"context"
	"errors"
	"net/http"

	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (s *secretsClient) GetCloudAccountSecretKey(ctx context.Context, header http.Header) (string, error) {
	util.CopyHeader(s.basicHeader, header)
	resp := new(metadata.SecretKeyResult)
	err := s.client.Get().
		SubResourcef(s.config.SecretKeyUrl).
		WithContext(ctx).
		WithHeaders(header).
		Do().Into(resp)

	if err != nil {
		return "", err
	}

	if resp.Code != 0 {
		return "", errors.New(resp.Message)
	}

	secretKey := resp.Data.Content.SecretKey
	if len(secretKey) != 16 && len(secretKey) != 24 && len(secretKey) != 32 {
		return "", errors.New("secret_key is invalid, it must be 128,192 or 256 bit")
	}

	return secretKey, nil
}
