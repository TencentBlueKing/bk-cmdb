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

// ErrorCode error code type
type ErrorCode string

const (
	// INVALID_REQUEST invalid request
	INVALID_REQUEST ErrorCode = "BAD_REQUEST"
	// UNAUTHENTICATED unauthenticated
	UNAUTHENTICATED ErrorCode = "UNAUTHORIZED"
	// NO_PERMISSION no permission
	NO_PERMISSION ErrorCode = "FORBIDDEN"
	// RATELIMIT_EXCEED rate limit exceed
	RATELIMIT_EXCEED ErrorCode = "TOO_MANY_REQUESTS"
	// INTERNAL internal error
	INTERNAL ErrorCode = "SERVER_ERROR"
	// NOT_FOUND not found
	NOT_FOUND ErrorCode = "NOT_FOUND"
	// METHOD_NOT_ALLOWED method not allowed
	METHOD_NOT_ALLOWED ErrorCode = "METHOD_NOT_ALLOWED"
	// UNKNOWN unknown error
	UNKNOWN ErrorCode = "UNKNOWN_ERROR"
)

// StatusCodeMap error code and status map
var (
	StatusCodeMap = map[ErrorCode]int{
		INVALID_REQUEST:    http.StatusBadRequest,
		UNAUTHENTICATED:    http.StatusUnauthorized,
		NO_PERMISSION:      http.StatusForbidden,
		RATELIMIT_EXCEED:   http.StatusTooManyRequests,
		INTERNAL:           http.StatusInternalServerError,
		NOT_FOUND:          http.StatusNotFound,
		METHOD_NOT_ALLOWED: http.StatusMethodNotAllowed,
		UNKNOWN:            http.StatusInternalServerError,
	}

	errCodeStatusMap = map[int]ErrorCode{
		http.StatusBadRequest:          INVALID_REQUEST,
		http.StatusUnauthorized:        UNAUTHENTICATED,
		http.StatusForbidden:           NO_PERMISSION,
		http.StatusTooManyRequests:     RATELIMIT_EXCEED,
		http.StatusInternalServerError: INTERNAL,
		http.StatusNotFound:            NOT_FOUND,
		http.StatusMethodNotAllowed:    METHOD_NOT_ALLOWED,
	}
)

// GetHTTPStatus get http status by error code
func GetHTTPStatus(code ErrorCode) int {
	if v, ok := StatusCodeMap[code]; ok {
		return v
	}
	return http.StatusBadRequest
}

// GetErrCodeByHTTPStatus get http status by error code
func GetErrCodeByHTTPStatus(status int) ErrorCode {
	if code, ok := errCodeStatusMap[status]; ok {
		return code
	}
	return UNKNOWN
}
