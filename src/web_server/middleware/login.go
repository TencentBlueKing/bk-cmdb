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
	"net/http"
	"strings"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/resource/esb"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/web_server/app/options"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Engine TODO
var Engine *backbone.Engine

// CacheCli TODO
var CacheCli redis.Client

const (
	message          string = "message"
	inaccessibleCode int    = 1302403
)

// ValidLogin valid the user login status
func ValidLogin(config options.Config, disc discovery.DiscoveryInterface) gin.HandlerFunc {

	return func(c *gin.Context) {
		rid := util.GetHTTPCCRequestID(c.Request.Header)
		pathArr := strings.Split(c.Request.URL.Path, "/")
		path1 := pathArr[1]

		// 删除 Accept-Encoding 避免返回值被压缩
		c.Request.Header.Del("Accept-Encoding")

		switch path1 {
		case "healthz", "metrics", "login", "static", "is_login":
			c.Next()
			return
		}

		if isAuthed(c, config) {
			handleAuthedReq(c, config, path1, disc, rid)
			return
		}

		if path1 == "api" {
			c.JSON(401, gin.H{
				"status": "log out",
			})
			c.Abort()
			return
		}

		user := user.NewUser(config, Engine, CacheCli)
		url := user.GetLoginUrl(c)
		c.Redirect(302, url)
		c.Abort()
	}

}

func handleAuthedReq(c *gin.Context, config options.Config, path1 string, disc discovery.DiscoveryInterface,
	rid string) {

	// http request header add user
	session := sessions.Default(c)
	userName, _ := session.Get(common.WEBSessionUinKey).(string)
	ownerID, _ := session.Get(common.WEBSessionOwnerUinKey).(string)
	bkToken, _ := session.Get(common.HTTPCookieBKToken).(string)
	bkTicket, _ := session.Get(common.HTTPCookieBKTicket).(string)
	language := webCommon.GetLanguageByHTTPRequest(c)
	c.Request.Header.Add(common.BKHTTPHeaderUser, userName)
	c.Request.Header.Add(common.BKHTTPLanguage, language)
	c.Request.Header.Add(common.BKHTTPOwnerID, ownerID)
	c.Request.Header.Add(common.HTTPCookieBKToken, bkToken)
	c.Request.Header.Add(common.HTTPCookieBKTicket, bkTicket)

	if config.LoginVersion == common.BKBluekingLoginPluginVersion {
		resp, err := esb.EsbClient().LoginSrv().GetUser(c.Request.Context(), c.Request.Header)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"status": fmt.Sprintf("get user from bk-login failed, err: %v", err)})
			c.Abort()
			return
		}

		if resp.Code == inaccessibleCode {
			data := gin.H{
				message: resp.Message,
			}
			c.HTML(http.StatusOK, webCommon.InaccessibleHtml, data)
			c.Abort()
			return
		}
	}

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
	return
}

// isAuthed check user is authed
func isAuthed(c *gin.Context, config options.Config) bool {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	user := user.NewUser(config, Engine, CacheCli)
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

	bkTokenName := common.HTTPCookieBKToken
	bkToken, err := c.Cookie(bkTokenName)
	blog.V(5).Infof("valid user login session token %s, cookie token %s, rid: %s", ccToken, bkToken, rid)
	if nil != err || bkToken != ccToken {
		return user.LoginUser(c)
	}
	return true

}
