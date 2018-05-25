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
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

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
	if strings.Contains(pathStr, BK_INST_SEARCH_OWNER) && strings.Contains(pathStr, BK_OBJECT_PLAT) {
		return true
	}
	if searchPatternRegexp.MatchString(pathStr) {
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
			//common inst op valid
			modelPrivi := session.Get("modelPrivi").(string)
			if 0 == len(modelPrivi) {
				blog.Error("get model privilege json error")
				return false
			}
			return validModelConfigPrivi(modelPrivi, method, pathArr)
		} else {
			//mainline inst op valid
			var objName string
			var mainLineObjIDArr []string
			if method == common.HTTPCreate {
				objName = pathArr[len(pathArr)-1]
			} else {
				objName = pathArr[len(pathArr)-2]
			}
			mainLineObjIDStr := session.Get("mainLineObjID").(string)
			err := json.Unmarshal([]byte(mainLineObjIDStr), &mainLineObjIDArr)
			if nil != err {
				blog.Error("get main line object id array false")
				return false
			}
			if util.InStrArr(mainLineObjIDArr, objName) {
				goto appvalid
			}

		}

		blog.Error("valid inst error")
		return false

	}

	//valid inst import privilege
	if strings.Contains(pathStr, BK_INSTSI) && !strings.Contains(pathStr, BK_IMPORT) {
		est := c.GetHeader(common.BKAppIDField)
		if "" == est {
			modelPrivi := session.Get("modelPrivi").(string)
			if 0 == len(modelPrivi) {
				blog.Error("get model privilege json error")
				return false
			}
			return validInstsOpPrivi(modelPrivi, method, pathArr)
		}
		blog.Error("valid inst error")
		return false

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
			return true
			//			sysPrivi := session.Get("sysPrivi")
			//			return validSysConfigPrivi(sysPrivi, BK_CC_AUDIT)
		}

	}

	//valid biz operaiton privilege
appvalid:
	return validAppConfigPrivi(c, method, pathStr)

}
