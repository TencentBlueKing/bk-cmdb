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

package service

import (
	"encoding/json"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/common/version"
	webCommon "configcenter/src/web_server/common"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
)

// Index html file
func (s *Service) Index(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)
	session := sessions.Default(c)
	role := session.Get("role")
	userName, _ := session.Get(common.WEBSessionUinKey).(string)
	language := webCommon.GetLanguageByHTTPRequest(c)

	if s.Config.Site.AuthScheme == "internal" {
		userPriviApp, rolePrivilege, modelPrivi, sysPrivi, mainLineObjIDArr := s.Logics.GetUserAppPri(ctx, userName, common.BKDefaultOwnerID, language)

		var strUserPriveApp, strRolePrivilege, strModelPrivi, strSysPrivi, mainLineObjIDStr string
		if nil == userPriviApp {
			strUserPriveApp = ""
		} else {
			cstrUserPriveApp, _ := json.Marshal(userPriviApp)
			strUserPriveApp = string(cstrUserPriveApp)
		}

		if nil == rolePrivilege {
			strRolePrivilege = ""
		} else {
			cstrRolePrivilege, _ := json.Marshal(rolePrivilege)
			strRolePrivilege = string(cstrRolePrivilege)
		}
		if nil == modelPrivi {
			strModelPrivi = ""
		} else {
			cstrModelPrivi, err := json.Marshal(modelPrivi)
			if err != nil {
				blog.Errorf("marshal model privilege failed, model privilege: %+v, err: %v, rid: %s", modelPrivi, err, rid)
			}
			strModelPrivi = string(cstrModelPrivi)
		}
		if nil == sysPrivi {
			strSysPrivi = ""
		} else {
			cstrSysPrivi, err := json.Marshal(sysPrivi)
			if err != nil {
				blog.Errorf("marshal system privilege failed, info: %+v, err: %v, rid: %s", sysPrivi, err, rid)
			}
			strSysPrivi = string(cstrSysPrivi)
		}

		mainLineObjIDB, err := json.Marshal(mainLineObjIDArr)
		if err != nil {
			blog.Errorf("marshal mainline failed, info: %+v, err: %v, rid: %s", mainLineObjIDArr, err, rid)
		}
		mainLineObjIDStr = string(mainLineObjIDB)

		session.Set("userPriviApp", string(strUserPriveApp))
		session.Set("rolePrivilege", string(strRolePrivilege))
		session.Set("modelPrivi", string(strModelPrivi))
		session.Set("sysPrivi", string(strSysPrivi))
		session.Set("mainLineObjID", string(mainLineObjIDStr))
		if err := session.Save(); err != nil {
			blog.Errorf("save session failed, err: %+v, rid: %s", err, rid)
		}

		// set cookie
		appIDArr := make([]string, 0)
		for key := range userPriviApp {
			appIDArr = append(appIDArr, strconv.FormatInt(key, 10))
		}
		appIDStr := strings.Join(appIDArr, "-")
		c.SetCookie("bk_privi_biz_id", appIDStr, 24*60*60, "", "", false, false)
	}

	c.HTML(200, "index.html", gin.H{
		"site":           s.Config.Site.DomainUrl,
		"version":        s.Config.Version,
		"ccversion":      version.CCVersion,
		"authscheme":     s.Config.Site.AuthScheme,
		"fullTextSearch": s.Config.Site.FullTextSearch,
		"role":           role,
		"curl":           s.Config.LoginUrl,
		"userName":       userName,
		"agentAppUrl":    s.Config.AgentAppUrl,
		"authCenter":     s.Config.AuthCenter,
	})
}
