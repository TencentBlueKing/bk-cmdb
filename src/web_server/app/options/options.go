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

// Package options TODO
package options

import (
	"configcenter/src/common"
	"configcenter/src/common/core/cc/config"
	"configcenter/src/storage/dal/redis"

	"github.com/spf13/pflag"
)

// ServerOption define option of server in flags
type ServerOption struct {
	ServConf         *config.CCAPIConfig
	DeploymentMethod common.DeploymentMethod
}

// NewServerOption create a ServerOption object
func NewServerOption() *ServerOption {
	s := ServerOption{
		ServConf:         config.NewCCAPIConfig(),
		DeploymentMethod: common.OpenSourceDeployment,
	}

	return &s
}

// AddFlags add flags
func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	s.ServConf.AddFlags(fs, "")
	fs.Var(&s.DeploymentMethod, "deployment-method", "The deployment method, supported value: open_source, blueking")
}

// Session TODO
type Session struct {
	Name            string
	DefaultLanguage string
	MultipleOwner   string
}

// Site TODO
type Site struct {
	AccountUrl      string
	DomainUrl       string
	HttpsDomainUrl  string
	HtmlRoot        string
	ResourcesPath   string
	BkLoginUrl      string
	BkHttpsLoginUrl string
	AppCode         string
	// available value: internal, iam
	AuthScheme string
	// available value: off, on
	FullTextSearch string
	PaasDomainUrl  string
	// BkDomain it is used to set the unified international language of Blue Whale.
	// this variable is returned to the front-end through configuration in the background.
	// the corresponding front-end variable is: cookieDomain.
	BkDomain string
	// BkComponentApiUrl is the blueking component api url, and is also the esb domain url
	// the corresponding front-end variable is: componentApiUrl.
	BkComponentApiUrl string
	HelpDocUrl        string
	// BkSharedResUrl is the blueking shared resource url
	BkSharedResUrl string
	// BkSharedResUrl is the blueking shared resource url
	BkApiUrlTmpl string
}

// Config TODO
type Config struct {
	Site                  Site
	Session               Session
	Redis                 redis.Config
	Version               string
	AgentAppUrl           string
	LoginUrl              string
	LoginVersion          string
	ConfigMap             map[string]string
	DeploymentMethod      common.DeploymentMethod
	EnableNotification    bool
	EnableMultiTenantMode bool
}
