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
	"plugin"
	"strings"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/web_server/app/options"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/user"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
	"gopkg.in/redis.v5"
)

var Engine *backbone.Engine
var CacheCli *redis.Client
var LoginPlg *plugin.Plugin

// ValidLogin valid the user login status
func ValidLogin(config options.Config, disc discovery.DiscoveryInterface) gin.HandlerFunc {

	return func(c *gin.Context) {
		rid := util.GetHTTPCCRequestID(c.Request.Header)
		pathArr := strings.Split(c.Request.URL.Path, "/")
		path1 := pathArr[1]

		switch path1 {
		case "healthz", "metrics":
			c.Next()
			return
		}

		if isAuthed(c, config) {
			// http request header add user
			session := sessions.Default(c)
			userName, _ := session.Get(common.WEBSessionUinKey).(string)
			ownerID, _ := session.Get(common.WEBSessionOwnerUinKey).(string)
			supplierID, _ := session.Get(common.WEBSessionSupplierID).(string)
			language := webCommon.GetLanguageByHTTPRequest(c)
			c.Request.Header.Add(common.BKHTTPHeaderUser, userName)
			c.Request.Header.Add(common.BKHTTPLanguage, language)
			c.Request.Header.Add(common.BKHTTPOwnerID, ownerID)
			c.Request.Header.Add(common.BKHTTPSupplierID, supplierID)

			if path1 == "api" {
				servers, err := disc.ApiServer().GetServers()
				if nil != err || 0 == len(servers) {
					blog.Errorf("no api server can be used. err: %v, rid: %s", err, rid)
					c.JSON(503, gin.H{
						"status": "no api server can be used.",
					})
					c.Abort()
					return
				}
				url := servers[0]
				httpclient.ProxyHttp(c, url)

			} else {
				c.Next()
			}
		} else {
			if path1 == "api" {
				c.JSON(401, gin.H{
					"status": "log out",
				})
				c.Abort()
				return
			} else {
				user := user.NewUser(config, Engine, CacheCli, LoginPlg)
				url := user.GetLoginUrl(c)
				c.Redirect(302, url)
				c.Abort()
			}

		}
	}

}

// IsAuthed check user is authed
func isAuthed(c *gin.Context, config options.Config) bool {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	if "1" == config.Session.Skip {
		session := sessions.Default(c)

		cookieOwnerID, err := c.Cookie(common.BKHTTPOwnerID)
		if "" == cookieOwnerID || nil != err {
			c.SetCookie(common.BKHTTPOwnerID, common.BKDefaultOwnerID, 0, "/", "", false, false)
			session.Set(common.WEBSessionOwnerUinKey, cookieOwnerID)
		} else if cookieOwnerID != session.Get(common.WEBSessionOwnerUinKey) {
			session.Set(common.WEBSessionOwnerUinKey, cookieOwnerID)
		}
		session.Set(common.WEBSessionSupplierID, "0")

		blog.V(5).Infof("skip login, cookie language: %s, cookieOwnerID: %s, rid: %s", webCommon.GetLanguageByHTTPRequest(c), cookieOwnerID, rid)
		session.Set(common.WEBSessionUinKey, "admin")

		session.Set(common.WEBSessionRoleKey, "1")
		session.Set(webCommon.IsSkipLogin, "1")
		if err := session.Save(); err != nil {
			blog.Warnf("save session failed, err: %s, rid: %s", err.Error(), rid)
		}
		return true
	}
	user := user.NewUser(config, Engine, CacheCli, LoginPlg)
	session := sessions.Default(c)

	// check bk_token
	ccToken := session.Get(common.HTTPCookieBKToken)
	if ccToken == nil {
		blog.Errorf("session key %s not found, rid: %s", common.HTTPCookieBKToken, rid)
		return user.LoginUser(c)
	}

	// check username
	userName, ok := session.Get(common.WEBSessionUinKey).(string)
	if !ok || "" == userName {
		return user.LoginUser(c)
	}

	// check owner_uin
	ownerID, ok := session.Get(common.WEBSessionOwnerUinKey).(string)
	if !ok || "" == ownerID {
		return user.LoginUser(c)
	}

	// check supplier_id
	supplierID, ok := session.Get(common.WEBSessionSupplierID).(string)
	if !ok || "" == supplierID {
		return user.LoginUser(c)
	}

	bkTokenName := common.HTTPCookieBKToken
	if nil != LoginPlg {
		bkPluginTokenName, err := LoginPlg.Lookup("BKTokenName")
		if nil == err {
			bkTokenName = *bkPluginTokenName.(*string)
		}
	}
	bkToken, err := c.Cookie(bkTokenName)
	blog.V(5).Infof("valid user login session token %s, cookie token %s, rid: %s", ccToken, bkToken, rid)
	if nil != err || bkToken != ccToken {
		return user.LoginUser(c)
	}
	return true

}
