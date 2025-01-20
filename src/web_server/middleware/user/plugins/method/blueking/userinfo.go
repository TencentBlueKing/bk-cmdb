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

// Package blueking defines user login method in blueking system
package blueking

import (
	"fmt"
	"strings"

	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	apigwcli "configcenter/src/common/resource/apigw"
	"configcenter/src/common/resource/esb"
	"configcenter/src/web_server/middleware/user/plugins/manager"

	"github.com/gin-gonic/gin"
)

func init() {
	plugin := &metadata.LoginPluginInfo{
		Name:       "blueking login system",
		Version:    common.BKBluekingLoginPluginVersion,
		HandleFunc: &user{},
	}
	manager.RegisterPlugin(plugin) // ("blueking login system", "self", "")
}

type user struct{}

// LoginUser user login
func (m *user) LoginUser(c *gin.Context, config map[string]string, isMultiOwner bool) (user *metadata.LoginUserInfo,
	loginSucc bool) {
	rid := httpheader.GetRid(c.Request.Header)

	bkTokens := getBkTokens(c)
	if len(bkTokens) == 0 {
		blog.Infof("LoginUser failed, bk_token empty, rid: %s", rid)
		return nil, false
	}

	for _, bkToken := range bkTokens {
		userInfo, err := apigwcli.Client().Login().GetUserByToken(c.Request.Context(), c.Request.Header, bkToken)
		if err != nil {
			blog.Errorf("get user info by token %s failed, err: %v, rid: %s", bkToken, err, rid)
			continue
		}

		user = &metadata.LoginUserInfo{
			UserName:  userInfo.Username,
			ChName:    userInfo.DisplayName,
			BkToken:   bkToken,
			TenantUin: userInfo.TenantID,
			IsTenant:  false,
			Language:  userInfo.Language,
		}
		break
	}
	if user == nil {
		return nil, false
	}
	return user, true
}

// getBkTokens get the values of the bk_token in the cookie
func getBkTokens(c *gin.Context) (bkTokens []string) {
	cookies := c.Request.Cookies()
	if len(cookies) == 0 {
		return bkTokens
	}
	for i := len(cookies) - 1; i >= 0; i-- {
		if cookies[i].Name == common.HTTPCookieBKToken {
			bkTokens = append(bkTokens, cookies[i].Value)
		}
	}
	return bkTokens
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
func (m *user) GetUserList(c *gin.Context, params map[string]string) ([]*metadata.LoginSystemUserInfo,
	*errors.RawErrorInfo) {
	rid := httpheader.GetRid(c.Request.Header)
	query := c.Request.URL.Query()
	for key, values := range query {
		params[key] = strings.Join(values, ";")
	}

	// try to use esb user list api
	result, err := esb.EsbClient().User().ListUsers(c.Request.Context(), c.Request.Header, params)
	if err != nil {
		blog.Errorf("get users by esb client failed, http failed, err: %+v, rid: %s", err, rid)
		return nil, &errors.RawErrorInfo{
			ErrCode: common.CCErrCommHTTPDoRequestFailed,
		}
	}

	if !result.Result {
		blog.Errorf("request esb, get user list failed, err: %v, rid: %s", result.Message, result.EsbRequestID)
		return nil, &errors.RawErrorInfo{
			ErrCode: common.CCErrCommHTTPDoRequestFailed,
		}
	}

	users := make([]*metadata.LoginSystemUserInfo, 0)
	for _, userInfo := range result.Data {
		user := &metadata.LoginSystemUserInfo{
			CnName: userInfo.DisplayName,
			EnName: userInfo.Username,
		}
		users = append(users, user)
	}

	return users, nil
}
