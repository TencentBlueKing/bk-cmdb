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

// Package config defines common config info.
package config

import (
	"errors"
	"fmt"

	"github.com/rs/xid"
)

// ServiceName represents the service name.
type ServiceName string

const (
	// AdminServer is the name of admin server.
	AdminServer ServiceName = "admin-server"
	// ApiServer is the name of api server.
	ApiServer ServiceName = "api-server"
	// CoreServer is the name of core server.
	CoreServer ServiceName = "core-server"
	// Collector is the name of data collector.
	Collector ServiceName = "collector"
	// Governancer is the name of data governancer.
	Governancer ServiceName = "governancer"
	// AuthServer is the name of auth server.
	AuthServer ServiceName = "auth-server"
	// TaskServer is the name of task server.
	TaskServer ServiceName = "task-server"
	// CDCServer is the name of cdc server.
	CDCServer ServiceName = "cdc-server"
	// WebServer is the name of web server.
	WebServer ServiceName = "web-server"
)

// ServerInfo is the basic server info.
type ServerInfo struct {
	// Name is the service name.
	Name ServiceName `json:"name"`
	// IP is the service listen ip.
	IP string `json:"ip"`
	// Port is the service listen port.
	Port uint `json:"port"`
	// RegisterIP is the service ip used for registration.
	RegisterIP string `json:"register_ip"`
	// Scheme is the http scheme of the service.
	Scheme string `json:"scheme"`
	// UUID is used to distinguish which service is master.
	UUID string `json:"uuid"`
	// Environment is the server's environment, servers can only discover other servers in the same environment.
	Environment string `json:"env"`
}

// Validate server info.
func (s *ServerInfo) Validate() error {
	if len(s.Name) == 0 {
		return errors.New("service name cannot be empty")
	}

	if len(s.IP) == 0 || s.IP == "0.0.0.0" {
		return errors.New("service ip is invalid")
	}

	if len(s.RegisterIP) == 0 {
		s.RegisterIP = s.IP
	}

	if s.RegisterIP == "0.0.0.0" {
		return errors.New("service register ip is invalid")
	}

	if len(s.Scheme) == 0 {
		s.Scheme = "http"
	}

	if s.Port == 0 || s.Port > 65535 {
		return errors.New("service port is invalid")
	}

	if len(s.UUID) == 0 {
		s.UUID = xid.New().String()
	}

	return nil
}

// RegisterAddress get register address of the server.
func (s *ServerInfo) RegisterAddress() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%s://%s:%d", s.Scheme, s.RegisterIP, s.Port)
}

// Instance get the instance identifier of the server.
func (s *ServerInfo) Instance() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}
