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

package options

import (
	"configcenter/src/ac/iam"
	"configcenter/src/common/auth"
	"configcenter/src/common/core/cc/config"
	"configcenter/src/scene_server/event_server/sync/hostidentifier"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/gse/client"

	"github.com/spf13/pflag"
)

// ServerOption is options of server.
type ServerOption struct {
	// ServConf is CC API config.
	ServConf *config.CCAPIConfig
}

// NewServerOption creates a new ServerOption object.
func NewServerOption() *ServerOption {
	return &ServerOption{
		ServConf: config.NewCCAPIConfig(),
	}
}

// AddFlags add flags to server options.
func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ServConf.AddrPort, "addrport", "127.0.0.1:60009", "The ip address and port for the serve on")
	fs.StringVar(&s.ServConf.RegDiscover, "regdiscv", "", "hosts of register and discover server. e.g: 127.0.0.1:2181")
	fs.StringVar(&s.ServConf.ExConfig, "config", "", "The config path. e.g conf/api.conf")
	fs.StringVar(&s.ServConf.RegisterIP, "register-ip", "", "the ip address registered on zookeeper, it can be domain")
	fs.Var(auth.EnableAuthFlag, "enable-auth", "The auth center enable status, true for enabled, false for disabled")
}

// Config is configs for event server.
type Config struct {
	// MongoDB is mongodb configs.
	MongoDB mongo.Config

	// Redis is cc redis configs.
	Redis redis.Config

	// Auth is auth config
	Auth iam.AuthConfig

	// IdentifierConf host identifier config
	IdentifierConf *hostidentifier.HostIdentifierConf

	// TaskConf gse taskServer connection config
	TaskConf *client.GseConnConfig

	// ApiConf gse apiServer connection config
	ApiConf *client.GseConnConfig
}
