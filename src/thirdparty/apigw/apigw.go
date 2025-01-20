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
	"configcenter/src/thirdparty/apigw/login"
	"configcenter/src/thirdparty/apigw/notice"

	"github.com/prometheus/client_golang/prometheus"
)

// ClientSet is the api gateway client set
type ClientSet interface {
	Gse() gse.ClientI
	Cmdb() cmdb.ClientI
	Notice() notice.ClientI
	Login() login.ClientI
}

type clientSet struct {
	gse    gse.ClientI
	cmdb   cmdb.ClientI
	notice notice.ClientI
	login  login.ClientI
}

// NewClientSet new api gateway client set
func NewClientSet(config *apigwutil.ApiGWConfig, metric prometheus.Registerer, neededClients []ClientType) (ClientSet,
	error) {

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

	cs := new(clientSet)

	neededCliMap := make(map[ClientType]struct{})
	for _, neededClient := range neededClients {
		neededCliMap[neededClient] = struct{}{}
	}

	if _, exists := neededCliMap[Gse]; exists {
		cs.gse, err = gse.NewClient(options)
		if err != nil {
			return nil, err
		}
	}

	if _, exists := neededCliMap[Cmdb]; exists {
		cs.cmdb, err = cmdb.NewClient(options)
		if err != nil {
			return nil, err
		}
	}

	if _, exists := neededCliMap[Notice]; exists {
		cs.notice, err = notice.NewClient(options)
		if err != nil {
			return nil, err
		}
	}

	if _, exists := neededCliMap[Login]; exists {
		cs.login, err = login.NewClient(options)
		if err != nil {
			return nil, err
		}
	}

	return cs, nil
}

// Gse returns gse client
func (c *clientSet) Gse() gse.ClientI {
	return c.gse
}

// Cmdb returns cmdb client
func (c *clientSet) Cmdb() cmdb.ClientI {
	return c.cmdb
}

// Notice returns bk-notice client
func (c *clientSet) Notice() notice.ClientI {
	return c.notice
}

// Login returns bk-login client
func (c *clientSet) Login() login.ClientI {
	return c.login
}

// ClientType is the api gateway client type, used to specify which client is needed
type ClientType string

const (
	// Gse is the gse client type
	Gse ClientType = "gse"
	// Cmdb is the cmdb client type
	Cmdb ClientType = "cmdb"
	// Notice is the notice client type
	Notice ClientType = "notice"
	// Login is the login client type
	Login ClientType = "login"
)
