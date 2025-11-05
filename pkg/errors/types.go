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
	// InvalidRequest invalid request
	InvalidRequest ErrorCode = "INVALID_REQUEST"
	// InvalidArgument invalid argument
	InvalidArgument ErrorCode = "INVALID_ARGUMENT"
	// StatusUnauthorized unauthenticated
	StatusUnauthorized ErrorCode = "UNAUTHORIZED"
	// NoPermission no permission
	NoPermission ErrorCode = "FORBIDDEN"
	// RateLimitExceed rate limit exceed
	RateLimitExceed ErrorCode = "TOO_MANY_REQUESTS"
	// Internal internal error
	Internal ErrorCode = "SERVER_ERROR"
	// NotFound not found
	NotFound ErrorCode = "NOT_FOUND"
	// MethodNotAllowed method not allowed
	MethodNotAllowed ErrorCode = "METHOD_NOT_ALLOWED"
	// Unknown  error
	Unknown ErrorCode = "UNKNOWN_ERROR"
)

// StatusCodeMap error code and status map
var (
	StatusCodeMap = map[ErrorCode]int{
		InvalidRequest:     http.StatusBadRequest,
		InvalidArgument:    http.StatusBadRequest,
		StatusUnauthorized: http.StatusUnauthorized,
		NoPermission:       http.StatusForbidden,
		RateLimitExceed:    http.StatusTooManyRequests,
		Internal:           http.StatusInternalServerError,
		NotFound:          http.StatusNotFound,
		MethodNotAllowed: http.StatusMethodNotAllowed,
		Unknown:            http.StatusInternalServerError,
	}

	errCodeStatusMap = map[int]ErrorCode{
		http.StatusBadRequest:          InvalidRequest,
		http.StatusUnauthorized:        StatusUnauthorized,
		http.StatusForbidden:           NoPermission,
		http.StatusTooManyRequests:     RateLimitExceed,
		http.StatusInternalServerError: Internal,
		http.StatusNotFound:            NotFound,
		http.StatusMethodNotAllowed:    MethodNotAllowed,
	}
)

// GetHTTPStatus get http status by error code
func GetHTTPStatus(code ErrorCode) int {
	if v, ok := StatusCodeMap[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

// GetErrCodeByHTTPStatus get http status by error code
func GetErrCodeByHTTPStatus(status int) ErrorCode {
	if code, ok := errCodeStatusMap[status]; ok {
		return code
	}
	return Unknown
}
