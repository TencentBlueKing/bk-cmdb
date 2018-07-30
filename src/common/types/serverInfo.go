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

// zk path
const (
	CC_SERV_BASEPATH      = "/cc/services/endpoints"
	CC_SERVCONF_BASEPATH  = "/cc/services/config"
	CC_SERVERROR_BASEPATH = "/cc/services/errors"
	CC_SERVLANG_BASEPATH  = "/cc/services/language"
)

// cc modules
const (
	CC_MODULE_AUDITCONTROLLER  = "auditcontroller"
	CC_MODULE_HOSTCONTROLLER   = "hostcontroller"
	CC_MODULE_OBJECTCONTROLLER = "objectcontroller"
	CC_MODULE_PROCCONTROLLER   = "proccontroller"
	CC_MODULE_DATACOLLECTION   = "datacollection"
	CC_MODULE_HOST             = "host"
	CC_MODULE_MIGRATE          = "migrate"
	CC_MODULE_PROC             = "proc"
	CC_MODULE_TOPO             = "topo"
	CC_MODULE_APISERVER        = "apiserver"
	CC_MODULE_WEBSERVER        = "webserver"
	CC_MODULE_EVENTSERVER      = "eventserver"
	GSE_MODULE_PROCSERVER      = "gseprocserver"
)

// AllModule all cc module
var AllModule = map[string]bool{
	CC_MODULE_AUDITCONTROLLER:  true,
	CC_MODULE_HOSTCONTROLLER:   true,
	CC_MODULE_OBJECTCONTROLLER: true,
	CC_MODULE_PROCCONTROLLER:   true,
	CC_MODULE_DATACOLLECTION:   true,
	CC_MODULE_HOST:             true,
	CC_MODULE_MIGRATE:          true,
	CC_MODULE_PROC:             true,
	CC_MODULE_TOPO:             true,
	CC_MODULE_APISERVER:        true,
	CC_MODULE_WEBSERVER:        true,
	CC_MODULE_EVENTSERVER:      true,
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
