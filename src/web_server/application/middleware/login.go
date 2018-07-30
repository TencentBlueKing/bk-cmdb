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
	"fmt"
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
func ValidLogin(params ...string) gin.HandlerFunc {
	loginURL := params[0]
	appCode := params[1]
	site := params[2]

	checkUrl = params[3]
	skipLogin := params[5]
	multipleOwner := params[6]
	isMultiOwner := true
	defaultlanguage := params[7]

	if "0" == multipleOwner {
		isMultiOwner = false
	}
	return func(c *gin.Context) {
		url := site + c.Request.URL.Path
		loginPage := fmt.Sprintf(loginURL, appCode, url)
		pathArr := strings.Split(c.Request.URL.Path, "/")
		path1 := pathArr[1]

		switch path1 {
		case "healthz", "metrics":
			c.Next()
			return
		}

		if isAuthed(c, isMultiOwner, skipLogin, defaultlanguage) {
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
			userName, _ := session.Get("userName").(string)
			language, _ := session.Get("language").(string)
			ownerID, _ := session.Get("owner_uin").(string)
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
				c.Redirect(301, loginPage)
			}

		}
	}

}

// IsAuthed check user is authed
func isAuthed(c *gin.Context, isMultiOwner bool, skipLogin, defaultlanguage string) bool {
	if "1" == skipLogin {
		blog.Info("skip login")
		session := sessions.Default(c)

		cookieLanuage, err := c.Cookie(common.BKHTTPCookieLanugageKey)
		if "" == cookieLanuage || nil != err {
			c.SetCookie(common.BKHTTPCookieLanugageKey, defaultlanguage, 0, "/", "", false, false)
			session.Set("language", defaultlanguage)
		} else if cookieLanuage != session.Get("lanugage") {
			session.Set("language", cookieLanuage)
		}

		cookieOwnerID, err := c.Cookie(common.BKHTTPOwnerID)
		if "" == cookieOwnerID || nil != err {
			c.SetCookie(common.BKHTTPOwnerID, common.BKDefaultOwnerID, 0, "/", "", false, false)
			session.Set("owner_uin", cookieOwnerID)
		} else if cookieOwnerID != session.Get("owner_uin") {
			session.Set("owner_uin", cookieOwnerID)
			ownerMan := user.NewOwnerManager("admin", cookieOwnerID, cookieLanuage)
			if err := ownerMan.InitOwner(); nil != err {
				blog.Errorf("init owner fail %s", err.Error())
				return true
			}

		}
		session.Set("userName", "admin")
		session.Set("role", "1")
		session.Set(webCommon.IsSkipLogin, "1")
		session.Save()
		return true
	}
	session := sessions.Default(c)
	cc_token := session.Get("bk_token")
	user := user.NewUser()
	if nil == cc_token {
		return user.LoginUser(c, checkUrl, isMultiOwner)
	}
	bk_token, err := c.Cookie("bk_token")
	blog.Info("valid user login session token %s, cookie token %s", cc_token, bk_token)
	if nil != err || bk_token != cc_token {
		return user.LoginUser(c, checkUrl, isMultiOwner)
	}
	return true

}
