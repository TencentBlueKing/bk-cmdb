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

// Package skip defines skip login method
package skip

import (
	"fmt"

	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/web_server/app/options"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/user/plugins/manager"
	"configcenter/src/web_server/middleware/user/plugins/method"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	plugin := &metadata.LoginPluginInfo{
		Name:       "skip login system",
		Version:    common.BKSkipLoginPluginVersion,
		HandleFunc: &user{},
	}
	manager.RegisterPlugin(plugin)
}

type user struct{}

// LoginUser user login
func (m *user) LoginUser(c *gin.Context, config options.Config, isMultiOwner bool) (user *metadata.LoginUserInfo,
	loginSucc bool) {

	rid := httpheader.GetRid(c.Request.Header)
	session := sessions.Default(c)
	tenantID, err := method.SetTenantFromCookie(c, config, session)
	if err != nil {
		blog.Errorf("set cookie failed, err: %v, rid: %s,", err, rid)
		return nil, false
	}

	user = &metadata.LoginUserInfo{
		UserName:  "admin",
		ChName:    "admin",
		BkToken:   "",
		TenantUin: tenantID,
		Language:  webCommon.GetLanguageByHTTPRequest(c),
	}
	return user, true
}

// GetLoginUrl get login url
func (m *user) GetLoginUrl(c *gin.Context, config map[string]string, input *metadata.LogoutRequestParams) string {
	var loginURL string
	var siteURL string
	var appCode string
	var err error
	if common.LogoutHTTPSchemeHTTPS == input.HTTPScheme {
		loginURL, err = cc.String("webServer.site.bkHttpsLoginUrl")
	} else {
		loginURL, err = cc.String("webServer.site.bkLoginUrl")
	}
	if err != nil {
		loginURL = ""
	}

	if common.LogoutHTTPSchemeHTTPS == input.HTTPScheme {
		siteURL, err = cc.String("webServer.site.httpsDomainUrl")
	} else {
		siteURL, err = cc.String("webServer.site.domainUrl")
	}
	if err != nil {
		siteURL = ""
	}

	appCode, err = cc.String("webServer.site.appCode")
	if err != nil {
		appCode = ""
	}
	loginURL = fmt.Sprintf(loginURL, appCode, fmt.Sprintf("%s%s", siteURL, c.Request.URL.String()))
	return loginURL
}

// GetUserList get user list
func (m *user) GetUserList(c *gin.Context, opts *metadata.GetUserListOptions) ([]*metadata.LoginSystemUserInfo,
	*errors.RawErrorInfo) {
	return []*metadata.LoginSystemUserInfo{
		{
			CnName: "admin",
			EnName: "admin",
		},
	}, nil
}
