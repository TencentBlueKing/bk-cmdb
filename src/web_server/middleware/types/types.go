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
	"regexp"
)

const BK_CC_MAINTAINERS string = "bk_biz_maintainer"

const (
	BK_CC_EVENT    string = "event"
	BK_CC_MODEL    string = "model"
	BK_CC_AUDIT    string = "audit"
	BK_CC_RESOURCE string = "resource"
)

const (
	BK_CC_UPDATE     string = "update"
	BK_CC_DELETE     string = "delete"
	BK_CC_CREATE     string = "create"
	BK_CC_SEARCH     string = "search"
	BK_CC_HOSTUPDATE string = "hostupdate"
	BK_CC_HOSTTRANS  string = "hosttrans"
	BK_CC_TOPOUPDATE string = "topoupdate"
	BK_CC_CUSTOMAPI  string = "customapi"
	BK_CC_PROCCONFIG string = "proconfig"
)

//bk search
const (
	BK_APP_SEARCH                    = "biz/search"
	BK_SET_SEARCH                    = "set/search"
	BK_MODULE_SEARCH                 = "module/search"
	BK_INST_SEARCH                   = "inst/search"
	BK_HOSTS_SEARCH                  = "hosts/search"
	BK_HOSTS_SNAP                    = "hosts/snapshot"
	BK_HOSTS_HIS                     = "hosts/history"
	BK_TOPO_MODEL                    = "topo/model"
	BK_INST_ASSOCIATION_TOPO_SEARCH  = "inst/association/topo/search"
	BK_INST_ASSOCIATION_OWNER_SEARCH = "inst/association/search/owner"
)

//bk topo
const (
	BK_INSTS  string = "inst"
	BK_TOPO   string = "topo"
	BK_SET    string = "set"
	BK_MODULE string = "module"
)

const (
	BK_INSTSI string = "insts"
	BK_IMPORT string = "import"
)

const (
	BK_INST_SEARCH_OWNER string = "inst/search/owner"
	BK_OBJECT_PLAT       string = "object/plat"
)

const BK_CC_CLASSIFIC string = "object/classifications"

const BK_CC_OBJECT_ATTR string = "object/attr/search"

const BK_CC_PRIVI_PATTERN string = "topo/privilege/user/detail"

const BK_OBJECT_ATT_GROUP = "objectatt/group/property/owner"

const BK_CC_USER_CUSTOM string = "usercustom"

const BK_CC_HOST_FAVORITES string = "favorites"

//proc manage pattern
const BK_PROC_S string = "proc"

//user api pattern
const BK_USER_API_S string = "userapi"

//host trans pattern
const BK_HOST_TRANS = "hosts/modules"

//system config privilege pattern
const (
	resPattern    = `(hosts/import|export)|(hosts/modules/resource/idle)`
	objectPattern = `object/classification/[a-z0-9A-Z]+/objects$`
)

//system config privilege regexp
var (
	ResPatternRegexp    = regexp.MustCompile(resPattern)
	ObjectPatternRegexp = regexp.MustCompile(objectPattern)
)

//host update string
const BK_HOST_UPDATE string = "hosts/batch"

//host manage search pattern
const searchPattern = `(hosts/[\w]+/[\d]+)|(topo/internal/[\w]+/[\w]+)|(topo/inst/[\w]+/[\d]+)|(object/host/inst/[\d]+)`

var SearchPatternRegexp = regexp.MustCompile(searchPattern)

var BK_CC_SYSCONFIG = []string{"event", "model", "audit"}

var BK_CC_MODEL_PRE = []string{"object", "objects"}

var BK_CC_AUDIT_PRE = []string{"audit"}

var BK_CC_EVENT_PRE = []string{"event"}

type LoginResult struct {
	Message string
	Code    string
	Result  bool
	Data    interface{}
}

type RolePriResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"bk_error_code"`
	Message interface{} `json:"bk_error_msg"`
	Data    []string    `json:"data"`
}

type RoleAppResult struct {
	Result  bool                     `json:"result"`
	Code    int                      `json:"bk_error_code"`
	Message interface{}              `json:"bk_error_msg"`
	Data    []map[string]interface{} `json:"data"`
}

type SearchAppResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"bk_error_code"`
	Message interface{} `json:"bk_error_msg"`
	Data    AppResult   `json:"data"`
}

type AppResult struct {
	Count int                      `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}
