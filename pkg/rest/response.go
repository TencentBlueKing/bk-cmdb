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

	"github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/i18n"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
)

// Renderer interface for managing response payloads.
type Renderer interface {
	Render(w http.ResponseWriter) error
}

// APIErrorResp response for api request error
type APIErrorResp struct {
	HTTPCode int             `json:"-"`     // http response status code
	Error    *cerr.RespError `json:"error"` // response error
}

// APISuccessResp response for api request success
type APISuccessResp struct {
	Data any `json:"data"` // response data
}

// Render chi render interface implementation for error response
func (e *APIErrorResp) Render(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(e.HTTPCode)

	return json.MarshalWrite(w, e)
}

// Render chi render interface implementation for success response
func (e *APISuccessResp) Render(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	return json.MarshalWrite(w, e)
}

// APIOK 正常返回
func APIOK(data any) Renderer {
	return &APISuccessResp{
		Data: data,
	}
}

// APIError 错误返回
func APIError(kt *kit.Kit, err error) Renderer {

	respErr := i18n.RespError(kt, err)

	return &APIErrorResp{
		HTTPCode: cerr.GetHTTPStatus(respErr.Code),
		Error:    respErr,
	}
}

// APIErrorWithStatus returns an API response with the given error and HTTP status code.
func APIErrorWithStatus(kt *kit.Kit, err error, statusCode int) Renderer {
	respErr := i18n.RespError(kt, err)

	if statusCode == 0 {
		statusCode = cerr.GetHTTPStatus(respErr.Code)
	}

	return &APIErrorResp{
		HTTPCode: statusCode,
		Error:    respErr,
	}
}
