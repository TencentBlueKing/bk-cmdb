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

package skip_login

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/user/plugins/manager"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
)

func init() {
	plugin := &metadata.LoginPluginInfo{
		Name:       "skip login system",
		Version:    "skip-login",
		HandleFunc: &user{},
	}
	manager.RegisterPlugin(plugin)
}

type user struct {
}

// LoginUser  user login
func (m *user) LoginUser(c *gin.Context, config map[string]string, isMultiOwner bool) (user *metadata.LoginUserInfo, loginSucc bool) {

	session := sessions.Default(c)

	cookieOwnerID, err := c.Cookie(common.BKHTTPOwnerID)
	if "" == cookieOwnerID || nil != err {
		c.SetCookie(common.BKHTTPOwnerID, common.BKDefaultOwnerID, 0, "/", "", false, false)
		session.Set(common.WEBSessionOwnerUinKey, cookieOwnerID)
	} else if cookieOwnerID != session.Get(common.WEBSessionOwnerUinKey) {
		session.Set(common.WEBSessionOwnerUinKey, cookieOwnerID)
	}

	user = &metadata.LoginUserInfo{
		UserName: "admin",
		ChName:   "admin",
		Phone:    "",
		Email:    "blueking",
		Role:     "",
		BkToken:  "",
		OnwerUin: "0",
		IsOwner:  false,
		Language: webCommon.GetLanguageByHTTPRequest(c),
	}
	return user, true
}

func (m *user) GetLoginUrl(c *gin.Context, config map[string]string, input *metadata.LogoutRequestParams) string {
	var ok bool
	var loginURL string
	var siteURL string

	if common.LogoutHTTPSchemeHTTPS == input.HTTPScheme {
		loginURL, ok = config["site.bk_https_login_url"]
	} else {
		loginURL, ok = config["site.bk_login_url"]
	}
	if !ok {
		loginURL = ""
	}
	if common.LogoutHTTPSchemeHTTPS == input.HTTPScheme {
		siteURL, ok = config["site.https_domain_url"]
	} else {
		siteURL, ok = config["site.domain_url"]
	}
	if !ok {
		siteURL = ""
	}

	appCode, ok := config["site.app_code"]
	if !ok {
		appCode = ""
	}
	loginURL = fmt.Sprintf(loginURL, appCode, fmt.Sprintf("%s%s", siteURL, c.Request.URL.String()))
	return loginURL
}
