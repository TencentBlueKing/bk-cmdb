// Package esbutil TODO
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
package esbutil

import (
	"net/http"

	httpheader "configcenter/src/common/http/header"
	"configcenter/src/thirdparty/apigw/apigwutil"
)

// EsbConfig TODO
type EsbConfig struct {
	Addrs     string
	AppCode   string
	AppSecret string
}

// EsbCommParams is esb common parameters
type EsbCommParams struct {
	TenantID string `json:"tenant_id"`
}

// SetEsbAuthHeader set esb authorization header
func SetEsbAuthHeader(esbConfig EsbConfig, header http.Header) http.Header {
	appConf := apigwutil.AppAuthConfig{
		AppCode:   esbConfig.AppCode,
		AppSecret: esbConfig.AppSecret,
	}
	return apigwutil.SetAuthHeader(appConf, header)
}

// GetEsbRequestParams get esb request parameters
func GetEsbRequestParams(esbConfig EsbConfig, header http.Header) *EsbCommParams {
	return &EsbCommParams{
		TenantID: httpheader.GetTenantID(header),
	}
}
