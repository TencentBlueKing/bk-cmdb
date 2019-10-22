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
)

// zk path
const (
	CCSvrBasePath      = "/cc/services/endpoints"
	CCSvrConfBasePath  = "/cc/services/config"
	CCSvrErrorBasePath = "/cc/services/errors"
	CCSvrLangBasePath  = "/cc/services/language"

	ccSrvUINodeName       = "ui"
	ccSrvAPINodeName      = "api"
	ccSrvSenceNodeName    = "sence"
	ccSrvResourceNodeName = "sence"
	ccSrvCommonNodeName   = "common"
	CCSvrUIBasePath       = CCSvrBasePath + "/" + ccSrvUINodeName
	CCSvrAPIBasePath      = CCSvrBasePath + "/" + ccSrvAPINodeName
	CCSvrSceneBasePath    = CCSvrBasePath + "/" + ccSrvSenceNodeName
	CCSvrResourceBasePath = CCSvrBasePath + "/" + ccSrvResourceNodeName
	CCSvrCommonBasePath   = CCSvrBasePath + "/" + ccSrvCommonNodeName

	CC_DISCOVERY_PREFIX = "cc_"
)

// cc modules
var (
	CCModuleDataCollection = SvrModuleInfo{Name: "datacollection", Layer: Sence}
	CCModuleHost           = SvrModuleInfo{Name: "host", Layer: Sence}
	CCModuleMigrate        = SvrModuleInfo{Name: "migrate", Layer: Sence}
	CCModuleProc           = SvrModuleInfo{Name: "proc", Layer: Sence}
	CCModuleTop            = SvrModuleInfo{Name: "topo", Layer: Sence}
	CCModuleAPIServer      = SvrModuleInfo{Name: "apiserver", Layer: API}
	CCModuleWebServer      = SvrModuleInfo{Name: "webserver", Layer: UI}
	CCModuleEventServer    = SvrModuleInfo{Name: "eventserver", Layer: Sence}
	CCModuleCoreService    = SvrModuleInfo{Name: "coreservice", Layer: Resource}
	GSEModuleProcServer    = SvrModuleInfo{Name: "gseprocserver", Layer: Sence}
	CCModuleTXC            = SvrModuleInfo{Name: "txc", Layer: Common}
	// CCModuleSynchronize multiple cmdb synchronize data server
	CCModuleSynchronize = SvrModuleInfo{Name: "sync", Layer: Sence}
	CCModuleOperation   = SvrModuleInfo{Name: "operation", Layer: Sence}
	CCModuleTask        = SvrModuleInfo{Name: "task", Layer: Sence}
)

// LayerModuleMap all cc module
// 根据节点前缀发下节点下面的服务， 这样保证依赖的服务可以自动发现，无需配置。
var LayerModuleMap = map[Layer][]SvrModuleInfo{
	// UI 层需要发现节点
	UI: []SvrModuleInfo{
		CCModuleAPIServer,
	},
	// API 层需要发现节点
	API: []SvrModuleInfo{
		CCModuleHost,
		CCModuleDataCollection,
		CCModuleProc,
		CCModuleTop,
		CCModuleEventServer,
		CCModuleSynchronize,
		CCModuleOperation,
		CCModuleTask,
		CCModuleTXC,
	},
	// Sence 层需要发现节点
	Sence: []SvrModuleInfo{
		CCModuleTXC,
		CCModuleTask,
		CCModuleCoreService,
	},
	// Resource 层需要发现节点
	Resource: []SvrModuleInfo{
		CCModuleTXC,
		CCModuleTask,
	},
}

// Layer curent layer name
type Layer int64

const (
	// UI webserver,  ui 层需要发现节点
	UI Layer = iota + 1
	// API  apiserver, api 层需要发现节点
	API
	// Sence layer. sence 层需要发现节点
	Sence
	// Resource controller layer. controller 层需要发现节点
	Resource
	// Common for common, Can be found by all services
	// reserved text
	Common
)

func (l Layer) String() string {
	switch l {
	case UI:
		return ccSrvUINodeName
	case API:
		return ccSrvAPINodeName
	case Sence:
		return ccSrvSenceNodeName
	case Resource:
		return ccSrvResourceNodeName
	case Common:
		return ccSrvCommonNodeName

	}
	return ""
}

// SvrModuleInfo service module information
type SvrModuleInfo struct {
	Name  string
	Layer Layer
}

func (smi *SvrModuleInfo) String() string {
	return smi.Layer.String() + "/" + smi.Name
}

// cc functionality define
const (
	CCFunctionalityServicediscover = "servicediscover"
	CCFunctionalityMongo           = "mongo"
	CCFunctionalityRedis           = "redis"
)

// ServerInfo define base server information
type ServerInfo struct {
	IP       string `json:"ip"`
	Port     uint   `json:"port"`
	HostName string `json:"hostname"`
	Scheme   string `json:"scheme"`
	Version  string `json:"version"`
	Pid      int    `json:"pid"`
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

// MigrateServInfo migrate server information
type MigrateServInfo struct {
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
func (s *ServerInfo) Address() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%s://%s:%d", s.Scheme, s.IP, s.Port)
}

func (s *ServerInfo) Instance() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}
