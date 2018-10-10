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

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/web_server/application/middleware/auth"
	"configcenter/src/web_server/application/middleware/user"
	webCommon "configcenter/src/web_server/common"
)

var APIAddr func() string
var sLoginURL string
var checkUrl string

//ValidLogin   valid the user login status
func ValidLogin(skipLogin, defaultlanguage string) gin.HandlerFunc {

	return func(c *gin.Context) {
		//url := site + c.Request.URL.Path
		//loginPage := fmt.Sprintf(loginURL, appCode, url)
		pathArr := strings.Split(c.Request.URL.Path, "/")
		path1 := pathArr[1]

		switch path1 {
		case "healthz", "metrics":
			c.Next()
			return
		}

		if isAuthed(c, skipLogin, defaultlanguage) {
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
				url := APIAddr() //apiSite
				if "" == url {
					blog.Errorf("get api server address error ")
				}
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
				user := user.NewUser()
				c.Redirect(302, user.GetLoginUrl(c))
			}

		}
	}

}

// IsAuthed check user is authed
func isAuthed(c *gin.Context, skipLogin, defaultlanguage string) bool {
	if "1" == skipLogin {
		session := sessions.Default(c)

		cookieLanuage, err := c.Cookie(common.BKHTTPCookieLanugageKey)
		if "" == cookieLanuage || nil != err {
			c.SetCookie(common.BKHTTPCookieLanugageKey, defaultlanguage, 0, "/", "", false, false)
			session.Set(common.WEBSessionLanguageKey, defaultlanguage)
		} else if cookieLanuage != session.Get(common.WEBSessionLanguageKey) {
			session.Set(common.WEBSessionLanguageKey, cookieLanuage)
		}

		cookieOwnerID, err := c.Cookie(common.BKHTTPOwnerID)
		if "" == cookieOwnerID || nil != err {
			c.SetCookie(common.BKHTTPOwnerID, common.BKDefaultOwnerID, 0, "/", "", false, false)
			session.Set(common.WEBSessionOwnerUinKey, cookieOwnerID)
		} else if cookieOwnerID != session.Get(common.WEBSessionOwnerUinKey) {
			session.Set(common.WEBSessionOwnerUinKey, cookieOwnerID)
			ownerMan := user.NewOwnerManager("admin", cookieOwnerID, cookieLanuage)
			if err := ownerMan.InitOwner(); nil != err {
				blog.Errorf("init owner fail %s", err.Error())
				return true
			}
		}

		blog.Info("skip login, cookieLanuage: %s, cookieOwnerID: %s", cookieLanuage, cookieOwnerID)

		session.Set(common.WEBSessionUinKey, "admin")
		session.Set(common.WEBSessionRoleKey, "1")
		session.Set(webCommon.IsSkipLogin, "1")
		session.Save()
		return true
	}
	session := sessions.Default(c)
	cc_token := session.Get(common.HTTPCookieBKToken)
	user := user.NewUser()
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
	blog.Info("valid user login session token %s, cookie token %s", cc_token, bk_token)
	if nil != err || bk_token != cc_token {
		return user.LoginUser(c)
	}
	return true

}
