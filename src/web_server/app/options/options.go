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
	"configcenter/src/common/core/cc/config"
	"configcenter/src/storage/dal/redis"

	"github.com/spf13/pflag"
)

//ServerOption define option of server in flags
type ServerOption struct {
	ServConf *config.CCAPIConfig
}

//NewServerOption create a ServerOption object
func NewServerOption() *ServerOption {
	s := ServerOption{
		ServConf: config.NewCCAPIConfig(),
	}

	return &s
}

//AddFlags add flags
func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ServConf.AddrPort, "addrport", "", "The ip address and port for the serve on")
	fs.StringVar(&s.ServConf.RegDiscover, "regdiscv", "",
		"hosts of register and discover server. e.g: 127.0.0.1:2379")
	fs.StringVar(&s.ServConf.RdUser, "rduser", "",
		"user name for authentication in register and discover")
	fs.StringVar(&s.ServConf.RdPassword, "rdpwd", "",
		"password for authentication in register and discover")
	fs.StringVar(&s.ServConf.RdCertFile, "rdcert", "", "cert file in register and discover")
	fs.StringVar(&s.ServConf.RdKeyFile, "rdkey", "", "key file in register and discover")
	fs.StringVar(&s.ServConf.RdCaFile, "rdca", "", "CA file in register and discover")
	fs.StringVar(&s.ServConf.RegisterIP, "register-ip", "",
		"the ip address registered on register and discover, it can be domain")
	fs.StringVar(&s.ServConf.ExConfig, "config", "", "The config path. e.g conf/ccapi.conf")
}

type Session struct {
	Name            string
	DefaultLanguage string
	MultipleOwner   string
}

type Site struct {
	AccountUrl      string
	DomainUrl       string
	HttpsDomainUrl  string
	HtmlRoot        string
	ResourcesPath   string
	BkLoginUrl      string
	BkHttpsLoginUrl string
	AppCode         string
	CheckUrl        string
	// available value: internal, iam
	AuthScheme string
	// available value: off, on
	FullTextSearch string
	PaasDomainUrl  string
	HelpDocUrl     string
}

type Config struct {
	Site         Site
	Session      Session
	Redis        redis.Config
	Version      string
	AgentAppUrl  string
	LoginUrl     string
	LoginVersion string
	ConfigMap    map[string]string
	AuthCenter   AppInfo
	DisableOperationStatistic bool
}

type AppInfo struct {
	AppCode string `json:"appCode"`
	URL     string `json:"url"`
}
