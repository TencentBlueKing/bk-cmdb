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

package middleware

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
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
	BK_APP_SEARCH    = "biz/search"
	BK_SET_SEARCH    = "set/search"
	BK_MODULE_SEARCH = "module/search"
	BK_INST_SEARCH   = "inst/search"
	BK_HOSTS_SEARCH  = "hosts/search"
	BK_HOSTS_SNAP    = "hosts/snapshot"
	BK_HOSTS_HIS     = "hosts/history"
	BK_TOPO_MODEL    = "topo/model"
)

//bk topo
const (
	BK_INSTS  string = "inst"
	BK_TOPO   string = "topo"
	BK_SET    string = "set"
	BK_MODULE string = "module"
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
	resPatternRegexp    = regexp.MustCompile(resPattern)
	objectPatternRegexp = regexp.MustCompile(objectPattern)
)

//host update string
const BK_HOST_UPDATE string = "hosts/batch"

//host manage search pattern
const searchPattern = `(hosts/[\w]+/[\d]+)|(topo/internal/[\w]+/[\w]+)|(topo/inst/[\w]+/[\d]+)|(object/host/inst/[\d]+)`

var searchPatternRegexp = regexp.MustCompile(searchPattern)

var BK_CC_SYSCONFIG = []string{"event", "model", "audit"}

var BK_CC_MODEL_PRE = []string{"object", "objects"}

var BK_CC_AUDIT_PRE = []string{"audit"}

var BK_CC_EVENT_PRE = []string{"event"}

//ValidResAccess valid resource access privilege
func ValidResAccess(pathArr []string, c *gin.Context) bool {
	var userName string
	session := sessions.Default(c)
	role := session.Get("role")
	pathStr := c.Request.URL.Path
	method := c.Request.Method

	//admin have full privilege
	if nil != role {
		irole := role.(string)
		if "1" == irole {
			return true
		}
	}
	iuserName := session.Get("userName")
	if nil == iuserName {
		blog.Error("user name error")
		return false
	}

	userName = iuserName.(string)

	//index page or static page
	if 0 == len(pathArr) || "" == pathArr[1] || "static" == pathArr[1] {
		return true
	}

	//valid privilege url must match session
	if strings.Contains(pathStr, BK_CC_PRIVI_PATTERN) {
		if pathArr[len(pathArr)-1] == userName {
			return true
		}
		blog.Error("privilege user name error")
		return false
	}
	//search classfication return true
	if strings.Contains(pathStr, BK_CC_CLASSIFIC) && method == common.HTTPSelectPost {
		return true
	}

	//search object attr  return true
	if strings.Contains(pathStr, BK_CC_OBJECT_ATTR) && method == common.HTTPSelectPost {
		return true
	}

	//usercustom return true
	if strings.Contains(pathStr, BK_CC_USER_CUSTOM) {
		return true
	}

	//objectatt group return true
	if strings.Contains(pathStr, BK_OBJECT_ATT_GROUP) {
		return true
	}

	//favorites return true
	if strings.Contains(pathStr, BK_CC_HOST_FAVORITES) {
		return true
	}

	//search object return true
	if objectPatternRegexp.MatchString(pathStr) {
		return true
	}

	//valid resource config
	if resPatternRegexp.MatchString(pathStr) {
		blog.Debug("valid resource config: %v", pathStr)
		sysPrivi := session.Get("sysPrivi")
		return validSysConfigPrivi(sysPrivi, BK_CC_RESOURCE)

	}

	//valid inst  privilege  op
	if strings.Contains(pathStr, BK_INSTS) && !strings.Contains(pathStr, BK_TOPO) {
		est := c.GetHeader(common.BKAppIDField)
		if "" == est {
			modelPrivi := session.Get("modelPrivi").(string)
			if 0 == len(modelPrivi) {
				blog.Error("get model privilege json error")
				return false
			}
			return validModelConfigPrivi(modelPrivi, method, pathArr)
		}

	}

	if len(pathArr) > 3 {
		//valid system config exclude resource
		path3 := pathArr[3]
		if util.InArray(path3, BK_CC_MODEL_PRE) {
			//valid model config privilege
			sysPrivi := session.Get("sysPrivi")
			return validSysConfigPrivi(sysPrivi, BK_CC_MODEL)
		}
		if util.InArray(path3, BK_CC_EVENT_PRE) {
			//valid event config privilege
			sysPrivi := session.Get("sysPrivi")
			return validSysConfigPrivi(sysPrivi, BK_CC_EVENT)
		}
		if util.InArray(path3, BK_CC_AUDIT_PRE) {
			//valid event config privilege
			sysPrivi := session.Get("sysPrivi")
			return validSysConfigPrivi(sysPrivi, BK_CC_AUDIT)
		}

	}

	//biz  search privilege, return true
	if strings.Contains(pathStr, BK_APP_SEARCH) || strings.Contains(pathStr, BK_SET_SEARCH) || strings.Contains(pathStr, BK_MODULE_SEARCH) || strings.Contains(pathStr, BK_INST_SEARCH) || strings.Contains(pathStr, BK_HOSTS_SEARCH) {
		return true
	}
	if strings.Contains(pathStr, BK_HOSTS_SNAP) || strings.Contains(pathStr, BK_HOSTS_HIS) {
		return true
	}
	if strings.Contains(pathStr, BK_TOPO_MODEL) {
		return true
	}
	if searchPatternRegexp.MatchString(pathStr) {
		return true
	}

	//valid biz operaiton privilege
	return validAppConfigPrivi(c, method, pathStr)

}

//validSysConfigPrivi valid system access privilege
func validSysConfigPrivi(sysPrivi interface{}, config string) bool {
	if nil != sysPrivi {
		ssysPrivi := sysPrivi.(string)
		if 0 == len(ssysPrivi) {
			blog.Error("no system config privilege")
			return false
		}
		var sysPriObj []string
		err := json.Unmarshal([]byte(ssysPrivi), &sysPriObj)
		if nil != err {
			blog.Error("no system config privilege not json")
			return false
		}
		if util.InArray(config, sysPriObj) {
			return true
		}
	}
	blog.Error("system privilege not pass")
	return false
}

//validModelConfigPrivi valid model inst privilege
func validModelConfigPrivi(modelPrivi string, method string, pathArr []string) bool {

	var mPrivi map[string][]string
	var objName string
	err := json.Unmarshal([]byte(modelPrivi), &mPrivi)
	if nil != err {
		blog.Error("get model privilege json error")
		return false
	}
	if method == common.HTTPCreate {
		objName = pathArr[len(pathArr)-1]
	} else {
		objName = pathArr[len(pathArr)-2]
	}

	priviArr, ok := mPrivi[objName]
	if false == ok {
		blog.Error("get object privilege  error")
		return false
	}
	//valid update privilege
	if method == common.HTTPUpdate && util.InArray(BK_CC_UPDATE, priviArr) {
		return true
	}

	//valid delete privilege
	if method == common.HTTPDelete && util.InArray(BK_CC_DELETE, priviArr) {
		return true
	}

	//valid create privilege
	if method == common.HTTPCreate && util.InArray(BK_CC_CREATE, priviArr) && !util.InArray(BK_CC_SEARCH, pathArr) {
		return true
	}

	//valid search privilege
	if method == common.HTTPSelectPost && util.InArray(BK_CC_SEARCH, priviArr) && util.InArray(BK_CC_SEARCH, pathArr) {
		return true
	}
	blog.Error("modle privilege valid not pass")
	return false
}

//validAppConfigPrivi valid app privilege
func validAppConfigPrivi(c *gin.Context, method, pathStr string) bool {

	//validate host update privilege
	if strings.Contains(pathStr, BK_HOST_UPDATE) && method == common.HTTPUpdate {
		return validAppAccessPrivi(c, BK_CC_HOSTUPDATE)
	}

	//validate host trans privilege
	if strings.Contains(pathStr, BK_HOST_TRANS) {
		return validAppAccessPrivi(c, BK_CC_HOSTTRANS)
	}

	//validate topo update privilege
	if strings.Contains(pathStr, BK_SET) || strings.Contains(pathStr, BK_MODULE) || strings.Contains(pathStr, BK_INSTS) || strings.Contains(pathStr, BK_TOPO) {
		return validAppAccessPrivi(c, BK_CC_TOPOUPDATE)
	}

	//validate user customer api privilege
	if strings.Contains(pathStr, BK_USER_API_S) {
		return validAppAccessPrivi(c, BK_CC_CUSTOMAPI)
	}
	//validate process config privilege
	if strings.Contains(pathStr, BK_PROC_S) {
		return validAppAccessPrivi(c, BK_CC_PROCCONFIG)
	}

	return true
}

//validate app access privilege
func validAppAccessPrivi(c *gin.Context, appResource string) bool {
	session := sessions.Default(c)
	appID := c.Request.Header.Get(common.BKAppIDField)
	if "" == appID {
		blog.Error("no app id in header")
		return false
	}
	userPriviAppStr, ok := session.Get("userPriviApp").(string)
	if false == ok {
		blog.Error("get user privilege from session error")
		return false
	}

	rolePrivilege, ok := session.Get("rolePrivilege").(string)
	if false == ok {
		blog.Error("get role privilege from session error")
		return false
	}

	//valid opearion under biz
	var userPriviApp, rolePrivi map[string][]string
	err := json.Unmarshal([]byte(userPriviAppStr), &userPriviApp)
	if nil != err {
		blog.Error("user privi app json error")
		return false
	}
	appRole, ok := userPriviApp[appID]
	if false == ok {
		blog.Error("no user privi app ")
		return false
	}
	//maintainer role pass the valid
	if util.InArray(BK_CC_MAINTAINERS, appRole) {
		return true
	}

	err = json.Unmarshal([]byte(rolePrivilege), &rolePrivi)
	if nil != err {
		blog.Error("role privi json error ")
		return false
	}
	priviArr := make([]string, 0)
	for _, role := range appRole {
		privi := rolePrivi[role]
		for _, p := range privi {
			priviArr = append(priviArr, p)
		}
	}
	if util.Contains(priviArr, appResource) {
		return true
	}
	blog.Error("valid user app privilege false")
	return false
}
