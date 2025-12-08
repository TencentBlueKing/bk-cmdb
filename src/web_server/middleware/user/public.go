/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package user

import (
	"encoding/json"

	"configcenter/pkg/tenant/tools"
	"configcenter/src/apimachinery/apiserver"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/web_server/app/options"
	"configcenter/src/web_server/middleware/user/plugins"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type publicUser struct {
	config   options.Config
	engine   *backbone.Engine
	cacheCli redis.Client
	apiCli   apiserver.ApiServerClientInterface
}

// LoginUser  user login
func (m *publicUser) LoginUser(c *gin.Context) bool {
	rid := httpheader.GetRid(c.Request.Header)

	isMultiOwner := false
	loginSuccess := false
	var userInfo *metadata.LoginUserInfo
	multipleOwner := m.config.Session.MultipleOwner
	if common.LoginSystemMultiTenantTrue == multipleOwner {
		isMultiOwner = true
	}

	user := plugins.CurrentPlugin(m.config.LoginVersion)
	userInfo, loginSuccess = user.LoginUser(c, m.config, isMultiOwner)

	if !loginSuccess {
		blog.Infof("login user with plugin failed, rid: %s", rid)
		return false
	}
	strOwnerUinList := []byte("")
	if 0 != len(userInfo.TenantUinArr) {
		strOwnerUinList, _ = json.Marshal(userInfo.TenantUinArr)
	}

	session := sessions.Default(c)

	tenantID, err := tools.ValidateDisableTenantMode(userInfo.TenantUin, m.config.EnableMultiTenantMode)
	if err != nil {
		blog.Errorf("get tenant id failed, err: %v", err)
		return false
	}

	session.Set(common.WEBSessionTenantUinKey, tenantID)
	session.Set(common.WEBSessionUinKey, userInfo.UserName)
	session.Set(common.WEBSessionChineseNameKey, userInfo.ChName)
	session.Set(common.HTTPCookieBKToken, userInfo.BkToken)
	session.Set(common.HTTPCookieBKTicket, userInfo.BkTicket)
	session.Set(common.WEBSessionAvatarUrlKey, userInfo.AvatarUrl)
	session.Set(common.WEBSessionTenantUinListeKey, string(strOwnerUinList))
	if userInfo.MultiTenant {
		session.Set(common.WEBSessionMultiTenantKey, common.LoginSystemMultiTenantTrue)
	} else {
		session.Set(common.WEBSessionMultiTenantKey, common.LoginSystemMultiTenantFalse)
	}
	session.Set(common.WEBSessionTimeZoneKey, userInfo.TimeZone)

	if err := session.Save(); err != nil {
		blog.Warnf("save session failed, err: %s, rid: %s", err.Error(), rid)
	}
	return true
}

// GetLoginUrl TODO
func (m *publicUser) GetLoginUrl(c *gin.Context) string {

	params := new(metadata.LogoutRequestParams)
	err := json.NewDecoder(c.Request.Body).Decode(params)
	if nil != err || (common.LogoutHTTPSchemeHTTP != params.HTTPScheme && common.LogoutHTTPSchemeHTTPS != params.HTTPScheme) {
		params.HTTPScheme, err = c.Cookie(common.LogoutHTTPSchemeCookieKey)
		if nil != err || (common.LogoutHTTPSchemeHTTP != params.HTTPScheme && common.LogoutHTTPSchemeHTTPS != params.HTTPScheme) {
			params.HTTPScheme = common.LogoutHTTPSchemeHTTP
		}
	}

	user := plugins.CurrentPlugin(m.config.LoginVersion)
	return user.GetLoginUrl(c, m.config.ConfigMap, params)
}

// GetUserList TODO
func (m *publicUser) GetUserList(c *gin.Context) ([]*metadata.LoginSystemUserInfo, *errors.RawErrorInfo) {
	user := plugins.CurrentPlugin(m.config.LoginVersion)
	return user.GetUserList(c, &metadata.GetUserListOptions{NeedAll: true})
}
