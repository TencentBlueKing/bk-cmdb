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

package apigwutil

import (
	"errors"
	"sync"

	"configcenter/src/apimachinery/util"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
)

// ApiGWBaseResponse api gateway base response
type ApiGWBaseResponse struct {
	Result  bool   `json:"result"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ApiGWConfig api gateway config
type ApiGWConfig struct {
	Address   []string
	AppCode   string
	AppSecret string
	Username  string
	TLSConfig *util.TLSClientConfig
}

// ApiGWDiscovery api gateway discovery struct
type ApiGWDiscovery struct {
	Servers []string
	index   int
	sync.Mutex
}

// GetServers get api gateway server
func (s *ApiGWDiscovery) GetServers() ([]string, error) {
	s.Lock()
	defer s.Unlock()

	num := len(s.Servers)
	if num == 0 {
		return []string{}, errors.New("oops, there is no server can be used")
	}

	if s.index < num-1 {
		s.index = s.index + 1
		return append(s.Servers[s.index-1:], s.Servers[:s.index-1]...), nil
	}

	s.index = 0
	return append(s.Servers[num-1:], s.Servers[:num-1]...), nil
}

// GetServersChan get api gateway chan
func (s *ApiGWDiscovery) GetServersChan() chan []string {
	return nil
}

// ParseApiGWConfig parse api gateway config
func ParseApiGWConfig(path string) (*ApiGWConfig, error) {
	address, err := cc.StringSlice(path + ".address")
	if err != nil {
		blog.Errorf("get api gateway address config error, err: %v", err)
		return nil, err
	}

	appCode, err := cc.String(path + ".appCode")
	if err != nil {
		blog.Errorf("get api gateway appCode config error, err: %v", err)
		return nil, err
	}

	appSecret, err := cc.String(path + ".appSecret")
	if err != nil {
		blog.Errorf("get api gateway appSecret config error, err: %v", err)
		return nil, err
	}

	username, err := cc.String(path + ".username")
	if err != nil {
		blog.Errorf("get api gateway username config error, err: %v", err)
		return nil, err
	}

	tlsConfig, err := util.NewTLSClientConfigFromConfig(path + ".tls")
	if err != nil {
		blog.Errorf("get api gateway tls config error, err: %v", err)
		return nil, err
	}

	apiGWConfig := &ApiGWConfig{
		Address:   address,
		AppCode:   appCode,
		AppSecret: appSecret,
		Username:  username,
		TLSConfig: &tlsConfig,
	}
	return apiGWConfig, nil
}
