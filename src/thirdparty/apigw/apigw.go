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
	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/apimachinery/util"
	"configcenter/src/thirdparty/apigw/apigwutil"
	"configcenter/src/thirdparty/apigw/cmdb"
	"configcenter/src/thirdparty/apigw/gse"

	"github.com/prometheus/client_golang/prometheus"
)

// ClientSet is the api gateway client set
type ClientSet interface {
	Gse() gse.ClientI
	Cmdb() cmdb.ClientI
}

type clientSet struct {
	gse  gse.ClientI
	cmdb cmdb.ClientI

	options *apigwutil.ApiGWOptions
}

// NewClientSet new api gateway client set
func NewClientSet(config *apigwutil.ApiGWConfig, metric prometheus.Registerer) (ClientSet, error) {
	apiMachineryConfig := &util.APIMachineryConfig{
		QPS:       2000,
		Burst:     2000,
		TLSConfig: config.TLSConfig,
	}

	client, err := util.NewClient(apiMachineryConfig.TLSConfig)
	if err != nil {
		return nil, err
	}

	flowControl := flowctrl.NewRateLimiter(apiMachineryConfig.QPS, apiMachineryConfig.Burst)

	options := &apigwutil.ApiGWOptions{
		Config: config,
		Capability: util.Capability{
			Client:     client,
			Throttle:   flowControl,
			MetricOpts: util.MetricOption{Register: metric},
		},
	}

	if options.Auth, err = apigwutil.GenDefaultAuthHeader(config); err != nil {
		return nil, err
	}

	return &clientSet{options: options}, nil
}

// Gse returns gse client
func (c *clientSet) Gse() gse.ClientI {
	if c.gse == nil {
		c.gse = gse.NewClient(c.options)
	}
	return c.gse
}

// Cmdb returns cmdb client
func (c *clientSet) Cmdb() cmdb.ClientI {
	if c.cmdb == nil {
		c.cmdb = cmdb.NewClient(c.options)
	}
	return c.cmdb
}
