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
	"encoding/json/v2"
	"net/http"

	ccError "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
)

// Renderer interface for managing response payloads.
type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

// APIResponse response for api request
type APIResponse struct {
	HTTPCode int                `json:"-"`               // http response status code
	Error    *ccError.RespError `json:"error,omitempty"` // response error
	Data     any                `json:"data,omitempty"`  // response data
}

// Render chi render interface implementation
func (e *APIResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(e.HTTPCode)

	return json.MarshalWrite(w, e)
}

// APIOK 正常返回
func APIOK(data any) Renderer {
	return &APIResponse{
		HTTPCode: http.StatusOK,
		Data:     data,
	}
}

// APIError 错误返回
func APIError(err *ccError.RespError) Renderer {
	err = ccError.GetDefaultErrorManager().ConvToRespError(err)

	if err == nil {
		err = &ccError.RespError{
			Code:    ccError.UNKNOWN,
			Message: "unknown error",
			System:  "",
		}
	}

	return &APIResponse{
		HTTPCode: ccError.GetHTTPStatus(err.Code),
		Error:    err,
	}
}
