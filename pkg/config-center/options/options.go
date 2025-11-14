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

// Package options define common app runtime option.
package options

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/TencentBlueKing/bk-cmdb/pkg/client/etcd"
	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// Options define common server runtime options.
type Options struct {
	Log    *log.HandlerOptions
	Server *config.ServerInfo
	Etcd   *etcd.Config
}

// NewOptions returns initialized Options
func NewOptions(service config.ServiceName) *Options {
	return &Options{
		Log: log.NewHandlerOptions(),
		Server: &config.ServerInfo{
			Name: service,
		},
		Etcd: new(etcd.Config),
	}
}

// Validate validates the common server options.
func (o *Options) Validate() error {
	if err := o.Log.Validate(); err != nil {
		return fmt.Errorf("log handler options validate failed: %w", err)
	}

	if err := o.Server.Validate(); err != nil {
		return fmt.Errorf("server info validate failed: %w", err)
	}

	if err := o.Etcd.Validate(); err != nil {
		return fmt.Errorf("etcd config validate failed: %w", err)
	}
	return nil
}

// AddFlags adds flags to fs and binds them to options.
func (o *Options) AddFlags(fs *pflag.FlagSet, isRPC bool) {
	o.Log.AddFlags(fs)
	o.Server.AddFlags(fs, isRPC)
	o.Etcd.AddFlags(fs)
}
