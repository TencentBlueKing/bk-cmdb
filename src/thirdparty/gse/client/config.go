/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package client

import (
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/ssl"
)

// GseConnConfig connect to gse config
type GseConnConfig struct {
	Endpoints []string
	TLSConf   *ssl.TLSClientConfig
}

// NewGseConnConfig new GseConnConfig struct
func NewGseConnConfig(prefix string) (*GseConnConfig, error) {
	endpoints, err := cc.StringSlice(prefix + ".endpoints")
	if err != nil {
		return nil, err
	}

	tlsConfig := new(ssl.TLSClientConfig)
	if tlsConfig.InsecureSkipVerify, err = cc.Bool(prefix + ".insecureSkipVerify"); err != nil {
		blog.Errorf("get gse %v insecureSkipVerify config error, err: %v", prefix, err)
		return nil, err
	}

	if tlsConfig.CertFile, err = cc.String(prefix + ".certFile"); err != nil {
		blog.Errorf("get gse %v certFile config error, err: %v", prefix, err)
		return nil, err
	}

	if tlsConfig.KeyFile, err = cc.String(prefix + ".keyFile"); err != nil {
		blog.Errorf("get gse %v keyFile config error, err: %v", prefix, err)
		return nil, err
	}

	if tlsConfig.CAFile, err = cc.String(prefix + ".caFile"); err != nil {
		blog.Errorf("get gse %v caFile config error, err: %v", prefix, err)
		return nil, err
	}

	if cc.IsExist(prefix + ".password") {
		if tlsConfig.Password, err = cc.String(prefix + ".password"); err != nil {
			blog.Errorf("get gse %v password config error, err: %v", prefix, err)
			return nil, err
		}
	}

	return &GseConnConfig{
		Endpoints: endpoints,
		TLSConf:   tlsConfig,
	}, nil
}
