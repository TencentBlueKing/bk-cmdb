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

package client

import (
	"configcenter/src/apimachinery/util"
	cc "configcenter/src/common/backbone/configcenter"
)

// GseConnConfig connect to gse config
type GseConnConfig struct {
	Endpoints []string
	TLSConf   *util.TLSClientConfig
}

// NewGseConnConfig new GseConnConfig struct
func NewGseConnConfig(prefix string) (*GseConnConfig, error) {
	endpoints, err := cc.StringSlice(prefix + ".endpoints")
	if err != nil {
		return nil, err
	}

	tlsConfig := &util.TLSClientConfig{}
	if tlsConfig.InsecureSkipVerify, err = cc.Bool(prefix + ".insecureSkipVerify"); err != nil {
		return nil, err
	}

	if tlsConfig.CertFile, err = cc.String(prefix + ".certFile"); err != nil {
		return nil, err
	}

	if tlsConfig.KeyFile, err = cc.String(prefix + ".keyFile"); err != nil {
		return nil, err
	}

	if tlsConfig.CAFile, err = cc.String(prefix + ".caFile"); err != nil {
		return nil, err
	}

	if val, err := cc.String(prefix + ".password"); err == nil {
		tlsConfig.Password = val
	}

	return &GseConnConfig{
		Endpoints: endpoints,
		TLSConf:   tlsConfig,
	}, nil
}
