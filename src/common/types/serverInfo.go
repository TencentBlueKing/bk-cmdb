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

package types

import (
	"fmt"
	"os"

	"configcenter/src/common/core/cc/config"
	"configcenter/src/common/version"

	"github.com/rs/xid"
)

// register and discover key base
const (
	CCDiscoverBaseEndpoint = "/bk_cmdb/endpoints"
	CCDiscoverBaseConfig   = "/bk_cmdb/config"
	CCDiscoverBaseError    = "/bk_cmdb/errors"
	CCDiscoverBaseLang     = "/bk_cmdb/language"
	CCDiscoverBaseNotice   = "/bk_cmdb/notice"
	CCDiscoverBaseLimiter  = "/bk_cmdb/limiter"
	CCDiscoverBaseElection = "/bk_cmdb/election"
)

// cc modules
const (
	CCModuleDataCollection = "datacollection"
	CCModuleHost           = "hostserver"
	CCModuleAdmin          = "adminserver"
	CCModuleProc           = "procserver"
	CCModuleTopo           = "toposerver"
	CCModuleApi            = "apiserver"
	CCModuleWeb            = "webserver"
	CCModuleEvent          = "eventserver"
	CCModuleCoreService    = "coreservice"
	CCModuleSync           = "syncserver"
	CCModuleOperation      = "operationserver"
	CCModuleTask           = "taskserver"
	CCModuleCloud          = "cloudserver"
	CCModuleAuth           = "authserver"
	CCModuleCacheService   = "cacheservice"
)

// AllModule all cc module
var AllModule = map[string]bool{
	CCModuleDataCollection: true,
	CCModuleHost:           true,
	CCModuleAdmin:          true,
	CCModuleProc:           true,
	CCModuleTopo:           true,
	CCModuleApi:            true,
	CCModuleWeb:            true,
	CCModuleEvent:          true,
	CCModuleCoreService:    true,
	// CCModuleSync: true,
	CCModuleOperation:    true,
	CCModuleTask:         true,
	CCModuleCloud:        true,
	CCModuleAuth:         true,
	CCModuleCacheService: true,
}

// cc functionality define
const (
	CCFunctionalityServicediscover = "servicediscover"
	CCFunctionalityMongo           = "mongo"
	CCFunctionalityRedis           = "redis"
)

const (
	CCConfigureRedis  = "redis"
	CCConfigureMongo  = "mongodb"
	CCConfigureCommon = "common"
	CCConfigureExtra  = "extra"
)

// ServerInfo define base server information
type ServerInfo struct {
	IP         string `json:"ip"`
	Port       uint   `json:"port"`
	RegisterIP string `json:"registerip"`
	HostName   string `json:"hostname"`
	Scheme     string `json:"scheme"`
	Version    string `json:"version"`
	Pid        int    `json:"pid"`
	// UUID is used to distinguish which service is master in zookeeper
	UUID string `json:"uuid"`
}

// NewServerInfo new a ServerInfo object
func NewServerInfo(conf *config.CCAPIConfig) (*ServerInfo, error) {
	ip, err := conf.GetAddress()
	if err != nil {
		return nil, err
	}

	port, err := conf.GetPort()
	if err != nil {
		return nil, err
	}

	registerIP := conf.RegisterIP
	// if no registerIP is set, default to be the ip
	if registerIP == "" {
		registerIP = ip
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	info := &ServerInfo{
		IP:         ip,
		Port:       port,
		RegisterIP: registerIP,
		HostName:   hostname,
		Scheme:     "http",
		Version:    version.GetVersion(),
		Pid:        os.Getpid(),
		UUID:       xid.New().String(),
	}
	return info, nil
}

// APIServerServInfo apiserver informaiton
type APIServerServInfo struct {
	ServerInfo
}

// WebServerInfo web server information
type WebServerInfo struct {
	ServerInfo
}

// AuditControllerServInfo audit-controller server information
type AuditControllerServInfo struct {
	ServerInfo
}

// HostControllerServInfo host-controller server information
type HostControllerServInfo struct {
	ServerInfo
}

// MigrateControllerServInfo migrate-controller server information
type MigrateControllerServInfo struct {
	ServerInfo
}

// ObjectControllerServInfo object-controller server information
type ObjectControllerServInfo struct {
	ServerInfo
}

// ProcControllerServInfo proc-controller server information
type ProcControllerServInfo struct {
	ServerInfo
}

// DataCollectionServInfo data-conllection server information
type DataCollectionServInfo struct {
	ServerInfo
}

// HostServerInfo host server information
type HostServerInfo struct {
	ServerInfo
}

// AdminServInfo admin server information
type AdminServInfo struct {
	ServerInfo
}

// ProcServInfo proc server information
type ProcServInfo struct {
	ServerInfo
}

// TopoServInfo topo server information
type TopoServInfo struct {
	ServerInfo
}

// EventServInfo topo server information
type EventServInfo struct {
	ServerInfo
}

// Address convert struct to host address
func (s *ServerInfo) RegisterAddress() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%s://%s:%d", s.Scheme, s.RegisterIP, s.Port)
}

func (s *ServerInfo) Instance() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}
