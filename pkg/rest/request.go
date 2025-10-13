/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

package rest

import (
	"context"
	"net/http"

	"github.com/TencentBlueKing/bk-cmdb/pkg/rest/codec"
	"github.com/TencentBlueKing/bk-cmdb/pkg/validator"
)

// decodeReq ...
func decodeReq[T any](r *http.Request) (*T, error) {
	in := new(T)

	// http.Request 直接返回
	if _, ok := any(in).(*http.Request); ok {
		return any(r).(*T), nil
	}

	// 空值不需要反序列化
	if _, ok := any(in).(*EmptyReq); ok {
		return in, nil
	}

	in, err := codec.Decode[T](r)
	if err != nil {
		return nil, err
	}

	return in, nil
}

// validate 参数校验
func validateReq(ctx context.Context, req any) error {
	// http.Request 直接返回
	if _, ok := req.(*http.Request); ok {
		return nil
	}

	// 空值不需要校验
	if _, ok := req.(*EmptyReq); ok {
		return nil
	}

	return validator.Struct(ctx, req)
}
