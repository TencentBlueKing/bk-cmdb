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

package cerr

import (
	"net/http"
)

// ErrorCodeType error code type
type ErrorCodeType string

const (
	// INVALID_REQUEST invalid request
	INVALID_REQUEST ErrorCodeType = "INVALID_REQUEST"
	// UNAUTHENTICATED unauthenticated
	UNAUTHENTICATED ErrorCodeType = "UNAUTHENTICATED"
	// NO_PERMISSION no permission
	NO_PERMISSION ErrorCodeType = "NO_PERMISSION"
	// RATELIMIT_EXCEED rate limit exceed
	RATELIMIT_EXCEED ErrorCodeType = "RATELIMIT_EXCEED"
	// INTERNAL internal error
	INTERNAL ErrorCodeType = "INTERNAL"
	// UNKNOWN unknown error
	UNKNOWN ErrorCodeType = "UNKNOWN"
)

// StatusCodeMap error code and status map
var (
	StatusCodeMap = map[ErrorCodeType]int{
		INVALID_REQUEST:  http.StatusBadRequest,
		UNAUTHENTICATED:  http.StatusUnauthorized,
		NO_PERMISSION:    http.StatusForbidden,
		RATELIMIT_EXCEED: http.StatusTooManyRequests,
		INTERNAL:         http.StatusInternalServerError,
		UNKNOWN:          http.StatusInternalServerError,
	}
)

// GetHTTPStatus get http status by error code
func GetHTTPStatus(code ErrorCodeType) int {
	if v, ok := StatusCodeMap[code]; ok {
		return v
	}
	return http.StatusBadRequest
}
