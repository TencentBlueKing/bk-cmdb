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

	"configcenter/src/common/util"
)

type EsbConfig struct {
	Addrs     string
	AppCode   string
	AppSecret string
}

type EsbCommParams struct {
	AppCode    string `json:"bk_app_code"`
	AppSecret  string `json:"bk_app_secret"`
	UserName   string `json:"bk_username"`
	SupplierID string `json:"bk_supplier_id"`
}

func GetEsbRequestParams(esbConfig EsbConfig, header http.Header) *EsbCommParams {
	return &EsbCommParams{
		AppCode:    esbConfig.AppCode,
		AppSecret:  esbConfig.AppSecret,
		UserName:   util.GetUser(header),
		SupplierID: util.GetOwnerID(header),
	}
}

func GetEsbQueryParameters(esbConfig EsbConfig, header http.Header) map[string]string {
	return map[string]string{
		"bk_app_code":   esbConfig.AppCode,
		"bk_app_secret": esbConfig.AppSecret,
		"bk_username":   util.GetUser(header),
	}
}
