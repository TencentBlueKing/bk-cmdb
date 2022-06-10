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

package apigw

import (
	"fmt"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/thirdparty/apigw/apigwutil"

	"github.com/prometheus/client_golang/prometheus"
)

var AuthKey = "x-bkapi-authorization"

// ApiGWSrv api gateway service
type ApiGWSrv struct {
	Client rest.ClientInterface
	Auth   string
}

// NewApiGW new a api gateway client
func NewApiGW(config *apigwutil.ApiGWConfig, reg prometheus.Registerer) (*ApiGWSrv, error) {

	apiMachineryConfig := &util.APIMachineryConfig{
		QPS:       2000,
		Burst:     2000,
		TLSConfig: config.TLSConfig,
	}

	client, err := util.NewClient(apiMachineryConfig.TLSConfig)
	if nil != err {
		return nil, err
	}

	flowControl := flowctrl.NewRateLimiter(apiMachineryConfig.QPS, apiMachineryConfig.Burst)

	esbCapability := &util.Capability{
		Client: client,
		Discover: &apigwutil.ApiGWDiscovery{
			Servers: config.Address,
		},
		Throttle:   flowControl,
		MetricOpts: util.MetricOption{Register: reg},
	}

	apigw := &ApiGWSrv{
		Client: rest.NewRESTClient(esbCapability, "/"),
		Auth: fmt.Sprintf(`{"bk_username": "%s", "bk_app_code": "%s", "bk_app_secret": "%s"}`, config.Username,
			config.AppCode, config.AppSecret),
	}
	return apigw, nil
}
