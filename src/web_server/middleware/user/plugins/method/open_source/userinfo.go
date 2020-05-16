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

package open_source

import (
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/user/plugins/manager"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
)

func init() {
	plugin := &metadata.LoginPluginInfo{
		Name:       "open source system",
		Version:    common.BKOpenSourceLoginPluginVersion,
		HandleFunc: &user{},
	}
	manager.RegisterPlugin(plugin)
}

type user struct{}

// LoginUser  user login
func (m *user) LoginUser(c *gin.Context, config map[string]string, isMultiOwner bool) (*metadata.LoginUserInfo, bool) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	session := sessions.Default(c)

	cookieOwnerID, err := c.Cookie(common.BKHTTPOwnerID)
	if "" == cookieOwnerID || nil != err {
		c.SetCookie(common.BKHTTPOwnerID, common.BKDefaultOwnerID, 0, "/", "", false, false)
		session.Set(common.WEBSessionOwnerUinKey, cookieOwnerID)
	} else if cookieOwnerID != session.Get(common.WEBSessionOwnerUinKey) {
		session.Set(common.WEBSessionOwnerUinKey, cookieOwnerID)
	}
	if err := session.Save(); err != nil {
		blog.Warnf("save session failed, err: %s, rid: %s", err.Error(), rid)
	}

	cookieUser, err := c.Cookie(common.BKUser)
	if "" == cookieUser || nil != err {
		blog.Errorf("login user not found, rid: %s", rid)
		return nil, false
	}

	loginTime, ok := session.Get(cookieUser).(int64)
	if !ok {
		blog.Errorf("login time not int64, rid: %s", rid)
		return nil, false
	}
	if time.Now().Unix()-loginTime < 24*60*60 {
		return &metadata.LoginUserInfo{
			UserName: cookieUser,
			ChName:   cookieUser,
			Phone:    "",
			Email:    "blueking",
			Role:     "",
			BkToken:  "",
			OnwerUin: "0",
			IsOwner:  false,
			Language: webCommon.GetLanguageByHTTPRequest(c),
		}, true
	}

	return nil, false
}

func (m *user) GetLoginUrl(c *gin.Context, config map[string]string, input *metadata.LogoutRequestParams) string {
	var ok bool
	var siteURL string

	if common.LogoutHTTPSchemeHTTPS == input.HTTPScheme {
		siteURL, ok = config["site.https_domain_url"]
	} else {
		siteURL, ok = config["site.domain_url"]
	}
	if !ok {
		siteURL = ""
	}
	return fmt.Sprintf("%s/login?c_url=%s%s", siteURL, siteURL, c.Request.URL.String())
}
