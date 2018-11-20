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
	"strings"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/web_server/app/options"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/auth"
	"configcenter/src/web_server/middleware/user"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
	redis "gopkg.in/redis.v5"
)

var sLoginURL string
var checkUrl string

var Engine *backbone.Engine
var CacheCli *redis.Client

//ValidLogin   valid the user login status
func ValidLogin(config options.Config, disc discovery.DiscoveryInterface) gin.HandlerFunc {

	return func(c *gin.Context) {
		pathArr := strings.Split(c.Request.URL.Path, "/")
		path1 := pathArr[1]

		switch path1 {
		case "healthz", "metrics":
			c.Next()
			return
		}

		if isAuthed(c, config) {
			//valid resource acess privilege
			auth := auth.NewAuth()
			ok := auth.ValidResAccess(pathArr, c)
			if false == ok {
				c.JSON(403, gin.H{
					"status": "access forbidden",
				})
				return
			}
			//http request header add user
			session := sessions.Default(c)
			userName, _ := session.Get(common.WEBSessionUinKey).(string)
			language, _ := session.Get(common.WEBSessionLanguageKey).(string)
			ownerID, _ := session.Get(common.WEBSessionOwnerUinKey).(string)
			c.Request.Header.Add(common.BKHTTPHeaderUser, userName)
			c.Request.Header.Add(common.BKHTTPLanguage, language)
			c.Request.Header.Add(common.BKHTTPOwnerID, ownerID)

			if path1 == "api" {
				servers, err := disc.ApiServer().GetServers()
				if nil != err || 0 == len(servers) {
					blog.Fatal("api server addr not right")
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
				return
			} else {
				user := user.NewUser(config, Engine, CacheCli)
				url := user.GetLoginUrl(c)
				c.Redirect(302, url)
			}

		}
	}

}

// IsAuthed check user is authed
func isAuthed(c *gin.Context, config options.Config) bool {
	if "1" == config.Session.Skip {
		session := sessions.Default(c)
		cookieLanuage, err := c.Cookie(common.BKHTTPCookieLanugageKey)
		if "" == cookieLanuage || nil != err {
			c.SetCookie(common.BKHTTPCookieLanugageKey, config.Session.DefaultLanguage, 0, "/", "", false, false)
			session.Set(common.WEBSessionLanguageKey, config.Session.DefaultLanguage)
		} else if cookieLanuage != session.Get(common.WEBSessionLanguageKey) {
			session.Set(common.WEBSessionLanguageKey, cookieLanuage)
		}

		cookieOwnerID, err := c.Cookie(common.BKHTTPOwnerID)
		if "" == cookieOwnerID || nil != err {
			c.SetCookie(common.BKHTTPOwnerID, common.BKDefaultOwnerID, 0, "/", "", false, false)
			session.Set(common.WEBSessionOwnerUinKey, cookieOwnerID)
		} else if cookieOwnerID != session.Get(common.WEBSessionOwnerUinKey) {
			session.Set(common.WEBSessionOwnerUinKey, cookieOwnerID)
		}

		blog.V(5).Infof("skip login, cookieLanuage: %s, cookieOwnerID: %s", cookieLanuage, cookieOwnerID)
		session.Set(common.WEBSessionUinKey, "admin")
		session.Set(common.WEBSessionRoleKey, "1")
		session.Set(webCommon.IsSkipLogin, "1")
		session.Save()
		return true
	}
	session := sessions.Default(c)
	cc_token := session.Get(common.HTTPCookieBKToken)
	user := user.NewUser(config, Engine, CacheCli)
	if nil == cc_token {
		return user.LoginUser(c)
	}
	userName, ok := session.Get(common.WEBSessionUinKey).(string)

	if !ok || "" == userName {
		return user.LoginUser(c)
	}
	ownerID, ok := session.Get(common.WEBSessionOwnerUinKey).(string)
	if !ok || "" == ownerID {
		return user.LoginUser(c)
	}

	bk_token, err := c.Cookie(common.HTTPCookieBKToken)
	blog.Infof("valid user login session token %s, cookie token %s", cc_token, bk_token)
	if nil != err || bk_token != cc_token {
		return user.LoginUser(c)
	}
	return true

}
