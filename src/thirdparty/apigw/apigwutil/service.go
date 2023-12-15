/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package apigwutil

import (
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
)

// ApiGWOptions is the api gateway client options
type ApiGWOptions struct {
	Config     *ApiGWConfig
	Auth       string
	Capability util.Capability
}

// ApiGWSrv api gateway service
type ApiGWSrv struct {
	Client rest.ClientInterface
	Config *ApiGWConfig
	Auth   string
}

// NewApiGW new api gateway service
func NewApiGW(options *ApiGWOptions, apiName ApiName) *ApiGWSrv {
	capability := options.Capability
	capability.Discover = &ApiGWDiscovery{
		Servers: ReplaceApiName(options.Config.Address, apiName),
	}

	apigw := &ApiGWSrv{
		Client: rest.NewRESTClient(&capability, "/"),
		Config: options.Config,
		Auth:   options.Auth,
	}

	return apigw
}
