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

package errors

import (
	"net/http"
	"sync"
)

// ErrorCodeType error code type
type ErrorCodeType string

const (
	// INVALID_ARGUMENT invalid argument
	INVALID_ARGUMENT ErrorCodeType = "INVALID_ARGUMENT"
	// INVALID_REQUEST invalid request
	INVALID_REQUEST ErrorCodeType = "INVALID_REQUEST"
	// OUT_OF_RANGE out of range
	OUT_OF_RANGE ErrorCodeType = "OUT_OF_RANGE"
	// FAILED_PRECONDITION failed precondition
	FAILED_PRECONDITION ErrorCodeType = "FAILED_PRECONDITION"
	// UNAUTHENTICATED unauthenticated
	UNAUTHENTICATED ErrorCodeType = "UNAUTHENTICATED"
	// IAM_NO_PERMISSION iam no permission
	IAM_NO_PERMISSION ErrorCodeType = "IAM_NO_PERMISSION"
	// NO_PERMISSION no permission
	NO_PERMISSION ErrorCodeType = "NO_PERMISSION"
	// NOT_FOUND not found
	NOT_FOUND ErrorCodeType = "NOT_FOUND"
	// ALREADY_EXISTS already exists
	ALREADY_EXISTS ErrorCodeType = "ALREADY_EXISTS"
	// ABORTED aborted
	ABORTED ErrorCodeType = "ABORTED"
	// RATELIMIT_EXCEED rate limit exceed
	RATELIMIT_EXCEED ErrorCodeType = "RATELIMIT_EXCEED"
	// RESOURCE_EXHAUSTED resource exhausted
	RESOURCE_EXHAUSTED ErrorCodeType = "RESOURCE_EXHAUSTED"
	// INTERNAL internal error
	INTERNAL ErrorCodeType = "INTERNAL"
	// UNKNOWN unknown error
	UNKNOWN ErrorCodeType = "UNKNOWN"
	// NOT_IMPLEMENTED api not implemented
	NOT_IMPLEMENTED ErrorCodeType = "NOT_IMPLEMENTED"
)

// StatusCodeMap error code and status map
var (
	StatusCodeMap = map[ErrorCodeType]int{
		INVALID_ARGUMENT:    http.StatusBadRequest,
		INVALID_REQUEST:     http.StatusBadRequest,
		OUT_OF_RANGE:        http.StatusBadRequest,
		FAILED_PRECONDITION: http.StatusBadRequest,
		UNAUTHENTICATED:     http.StatusUnauthorized,
		IAM_NO_PERMISSION:   http.StatusForbidden,
		NO_PERMISSION:       http.StatusForbidden,
		NOT_FOUND:           http.StatusNotFound,
		ALREADY_EXISTS:      http.StatusConflict,
		ABORTED:             http.StatusConflict,
		RATELIMIT_EXCEED:    http.StatusTooManyRequests,
		RESOURCE_EXHAUSTED:  http.StatusTooManyRequests,
		INTERNAL:            http.StatusInternalServerError,
		UNKNOWN:             http.StatusInternalServerError,
		NOT_IMPLEMENTED:     http.StatusNotImplemented,
	}
	statusMu sync.RWMutex
)

// GetHTTPStatus get http status by error code
func GetHTTPStatus(code ErrorCodeType) int {
	statusMu.RLock()
	defer statusMu.RUnlock()
	if v, ok := StatusCodeMap[code]; ok {
		return v
	}
	return http.StatusBadRequest
}

// RegisterHttpStatus register http status with error code
func RegisterHttpStatus(code ErrorCodeType, status int) {
	statusMu.Lock()
	defer statusMu.Unlock()
	StatusCodeMap[code] = status
}
