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

// Package options defines sync server options.
package options

import (
	"configcenter/src/common/core/cc/config"

	"github.com/spf13/pflag"
)

// ServerOption defines sync server options in flags
type ServerOption struct {
	ServConf *config.CCAPIConfig
}

// NewServerOption new ServerOption
func NewServerOption() *ServerOption {
	s := ServerOption{
		ServConf: config.NewCCAPIConfig(),
	}

	return &s
}

// AddFlags add flags for ServerOption
func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ServConf.AddrPort, "addrport", "127.0.0.1:60015", "The ip address and port of the server")
	fs.StringVar(&s.ServConf.RegDiscover, "regdiscv", "", "Register and discover service address. e.g: 127.0.0.1:2181")
	fs.StringVar(&s.ServConf.ExConfig, "config", "", "Extra config path. e.g conf/api.conf")
	fs.StringVar(&s.ServConf.RegisterIP, "register-ip", "", "The ip address registered on zookeeper, it can be domain")
}
