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

	"configcenter/src/common/http/httpclient"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginResult struct {
	Message string
	Code    string
	Result  bool
	Data    interface{}
}

var APIAddr func() string
var sLoginURL string
var check_url string

//ValidLogin   valid the user login status
func ValidLogin(params ...string) gin.HandlerFunc {
	loginURL := params[0]
	appCode := params[1]
	site := params[2]

	check_url = params[3]
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
		//		blog.Info("login page:%v", loginPage)
		pathArr := strings.Split(c.Request.URL.Path, "/")
		path1 := pathArr[1]

		switch path1 {
		case "healthz", "metrics":
			c.Next()
			return
		}

		if isAuthed(c, isMultiOwner, skipLogin, defaultlanguage) {
			//valid resource acess privilege
			ok := ValidResAccess(pathArr, c)
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

//isAuthed check user is authed
func isAuthed(c *gin.Context, isMultiOwner bool, skipLogin, defaultlanguage string) bool {
	if "1" == skipLogin {
		session := sessions.Default(c)

		cookieLanuage, err := c.Cookie(common.BKHTTPCookieLanugageKey)
		if "" == cookieLanuage || nil != err {
			c.SetCookie(common.BKHTTPCookieLanugageKey, defaultlanguage, 0, "/", "", false, false)
			session.Set("language", defaultlanguage)

		} else if cookieLanuage != session.Get("lanugage") {
			session.Set("language", cookieLanuage)
		}

		session.Set("userName", "admin")
		session.Set("role", "1")
		session.Set("owner_uin", "0")
		session.Set("skiplogin", "1")
		session.Save()
		return true
	}
	session := sessions.Default(c)
	cc_token := session.Get("bk_token")
	if nil == cc_token {
		return loginUser(c, isMultiOwner)
	}
	//	blog.Info("valid user login session token %s", cc_token)
	bk_token, err := c.Cookie("bk_token")
	//	blog.Info("valid user login cookie token %s", bk_token)
	if nil != err || bk_token != cc_token {
		return loginUser(c, isMultiOwner)
	}
	return true

}

//loginUser  user login
func loginUser(c *gin.Context, isMultiOwner bool) bool {
	bk_token, err := c.Cookie("bk_token")
	if nil != err {
		return false
	}
	if nil != err || 0 == len(bk_token) {
		return false
	}
	loginURL := check_url + bk_token
	httpCli := httpclient.NewHttpClient()
	httpCli.SetTimeOut(30 * time.Second)
	blog.Info("get user info cond: %s", string(loginURL))
	loginResult, err := httpCli.GET(loginURL, nil, nil)

	if nil != err {
		blog.Error("get user info return error: %v", err)
		return false
	}
	blog.Info("get user info return: %s", string(loginResult))
	var resultData LoginResult
	err = json.Unmarshal([]byte(loginResult), &resultData)
	if nil != err {
		blog.Error("get user info json error: %v", err)
		return false
	}
	userInfo, ok := resultData.Data.(map[string]interface{})
	if false == ok {
		blog.Error("get user info decode error: %v", err)
		return false
	}
	userName, ok := userInfo["username"]
	if false == ok {
		blog.Error("get user info username error: %v", err)
		return false
	}
	chName, ok := userInfo["chname"]
	if false == ok {
		blog.Error("get user info chname error: %v", err)
		return false
	}
	phone, ok := userInfo["phone"]
	if false == ok {
		blog.Error("get user info phone error: %v", err)
		return false
	}
	email, ok := userInfo["email"]
	if false == ok {
		blog.Error("get user info email error: %v", err)
		return false
	}
	role, ok := userInfo["role"]
	if false == ok {
		blog.Error("get user info role error: %v", err)
		return false
	}

	language, ok := userInfo["language"]
	if false == ok {
		blog.Error("get language info role error: %v", err)
	}
	ownerID := common.BKDefaultOwnerID
	if true == isMultiOwner {
		ownerID, ok = userInfo["owner_uin"].(string)
		if false == ok {
			blog.Error("get owner_uin info role error: %v", err)
			return false
		}
	}

	cookielanguage, _ := c.Cookie("blueking_language")
	session := sessions.Default(c)
	session.Set("userName", userName)
	session.Set("chName", chName)
	session.Set("phone", phone)
	session.Set("email", email)
	session.Set("role", role)
	session.Set("bk_token", bk_token)
	session.Set("owner_uin", ownerID)
	if "" != cookielanguage {
		session.Set("language", cookielanguage)
	} else {
		session.Set("language", language)
	}
	session.Save()
	return true
}
