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

// Package opensource open-source login method, use configuration to define the user & pwd to login
package opensource

import (
	"fmt"
	"strings"
	"time"

	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/user/plugins/manager"
	"configcenter/src/web_server/middleware/user/plugins/method"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

// LoginUser user login
func (m *user) LoginUser(c *gin.Context, config map[string]string, isMultiOwner bool) (*metadata.LoginUserInfo, bool) {
	rid := httpheader.GetRid(c.Request.Header)
	session := sessions.Default(c)

	if err := method.SetCookie(c, session); err != nil {
		blog.Errorf("set cookie failed, err: %v, rid: %s", err, rid)
		return nil, false
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
			UserName:  cookieUser,
			ChName:    cookieUser,
			Phone:     "",
			Email:     "blueking",
			BkToken:   "",
			TenantUin: session.Get(common.WEBSessionTenantUinKey).(string),
			IsTenant:  false,
			Language:  webCommon.GetLanguageByHTTPRequest(c),
		}, true
	}

	return nil, false
}

// GetLoginUrl get login url
func (m *user) GetLoginUrl(c *gin.Context, config map[string]string, input *metadata.LogoutRequestParams) string {
	var siteURL string
	var err error
	if common.LogoutHTTPSchemeHTTPS == input.HTTPScheme {
		siteURL, err = cc.String("webServer.site.httpsDomainUrl")
	} else {
		siteURL, err = cc.String("webServer.site.domainUrl")
	}
	if err != nil {
		siteURL = ""
	}
	siteURL = strings.TrimRight(siteURL, "/")
	return fmt.Sprintf("%s/login?c_url=%s%s", siteURL, siteURL, c.Request.URL.String())
}

// GetUserList get user list
func (m *user) GetUserList(c *gin.Context, config map[string]string) ([]*metadata.LoginSystemUserInfo,
	*errors.RawErrorInfo) {
	rid := httpheader.GetRid(c.Request.Header)
	users := make([]*metadata.LoginSystemUserInfo, 0)
	userInfo, err := cc.String("webServer.session.userInfo")
	if err != nil {
		blog.Errorf("User name and password can't be found at webServer.session.userInfo in config file common.yaml, rid:%s",
			rid)
		return nil, &errors.RawErrorInfo{
			ErrCode: common.CCErrWebNoUsernamePasswd,
		}
	}
	userInfos := strings.Split(userInfo, ",")
	for _, userInfo := range userInfos {
		userPasswd := strings.Split(userInfo, ":")
		if len(userPasswd) != 2 {
			blog.Errorf("The format of user name and password are wrong, please check webServer.session.userInfo in config file common.yaml, rid:%s",
				rid)
			return nil, &errors.RawErrorInfo{
				ErrCode: common.CCErrWebUserinfoFormatWrong,
			}
		}
		user := &metadata.LoginSystemUserInfo{
			CnName: userPasswd[0],
			EnName: userPasswd[0],
		}
		users = append(users, user)
	}

	return users, nil
}
