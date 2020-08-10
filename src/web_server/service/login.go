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
	"strings"
	"time"

	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/web_server/middleware/user"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
)

// LogOutUser log out user
func (s *Service) LogOutUser(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	c.Request.URL.Path = ""
	userManger := user.NewUser(*s.Config, s.Engine, s.CacheCli)
	loginURL := userManger.GetLoginUrl(c)
	ret := metadata.LogoutResult{}
	ret.BaseResp.Result = true
	ret.Data.LogoutURL = loginURL
	c.JSON(200, ret)
	return
}

// Login html file
func (s *Service) Login(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{})
}

// LoginUser log in user
func (s *Service) LoginUser(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(c.Request.Header))
	userName := c.PostForm("username")
	password := c.PostForm("password")
	if userName == "" || password == "" {
		c.HTML(200, "login.html", gin.H{
			"error": defErr.CCError(common.CCErrWebNeedFillinUsernamePasswd).Error(),
		})
	}
	userInfo, err := cc.String("webServer.session.userInfo")
	if err != nil {
		c.HTML(200, "login.html", gin.H{
			"error": defErr.CCError(common.CCErrWebNoUsernamePasswd).Error(),
		})
		return
	}
	userInfos := strings.Split(userInfo, ",")
	for _, userInfo := range userInfos {
		userWithPassword := strings.Split(userInfo, ":")
		if len(userWithPassword) != 2 {
			blog.Errorf("user info config %s invalid, rid: %s", userInfo, rid)
			c.HTML(200, "login.html", gin.H{
				"error": defErr.CCError(common.CCErrWebUserinfoFormatWrong).Error(),
			})
			return
		}
		if userWithPassword[0] == userName && userWithPassword[1] == password {
			c.SetCookie(common.BKUser, userName, 24*60*60, "/", "", false, false)
			session := sessions.Default(c)
			session.Set(userName, time.Now().Unix())
			if err := session.Save(); err != nil {
				blog.Warnf("save session failed, err: %s, rid: %s", err.Error(), rid)
			}
			userManger := user.NewUser(*s.Config, s.Engine, s.CacheCli)
			userManger.LoginUser(c)
			var redirectURL string
			if c.Param("c_url") != "" {
				redirectURL = c.Param("c_url")
			} else {
				redirectURL = s.Config.Site.DomainUrl
			}
			c.Redirect(302, redirectURL)
			return
		}
	}
	c.HTML(200, "login.html", gin.H{
		"error": defErr.CCError(common.CCErrWebUsernamePasswdWrong).Error(),
	})
	return
}
