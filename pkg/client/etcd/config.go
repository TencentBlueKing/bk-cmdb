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

package etcd

import (
	"errors"

	"github.com/spf13/pflag"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
)

// Config is etcd config.
type Config struct {
	// Endpoints is the etcd endpoints.
	Endpoints []string
	// Username is the etcd username for authentication.
	Username string
	// Password is the etcd password for authentication.
	Password string
	// TLS is the etcd tls config.
	TLS *config.TLSConfig
}

// Validate etcd config.
func (c *Config) Validate() error {
	if len(c.Endpoints) == 0 {
		return errors.New("etcd endpoints is not set")
	}

	return nil
}

// AddFlags adds etcd flags to flag set.
func (c *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&c.Endpoints, "etcd-endpoints", c.Endpoints, "etcd endpoints")
	fs.StringVar(&c.Username, "etcd-username", c.Username, "etcd username")
	fs.StringVar(&c.Password, "etcd-password", c.Password, "etcd password")

	c.TLS = new(config.TLSConfig)
	fs.BoolVar(&c.TLS.InsecureSkipVerify, "etcd-skip-verify", true, "etcd skip tls verification flag")
	fs.StringVar(&c.TLS.CAFile, "etcd-ca", "", "etcd tls ca file path")
	fs.StringVar(&c.TLS.CertFile, "etcd-cert", "", "etcd tls cert file path")
}
