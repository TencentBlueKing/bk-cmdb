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

package user

import (
	"encoding/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/web_server/app/options"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/user/plugins"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type publicUser struct {
	config options.Config
}

// LoginUser  user login
func (m *publicUser) LoginUser(c *gin.Context) bool {

	user := plugins.CurrentPlugin(c)
	isMultiOwner := false
	multipleOwner := m.config.Session.MultipleOwner
	if common.LoginSystemMultiSupplierTrue == multipleOwner {
		isMultiOwner = true
	}
	config := make(map[string]string)
	config["site.check_url"] = m.config.Site.CheckUrl
	userInfo, loginSucc := user.LoginUser(c, config, isMultiOwner)
	if !loginSucc {
		return false
	}

	if true == isMultiOwner || true == userInfo.MultiSupplier {
		err := NewOwnerManager(userInfo.UserName, userInfo.OnwerUin, userInfo.Language).InitOwner()
		if nil != err {
			blog.Error("InitOwner error: %v", err)
			return false
		}
	}
	strOwnerUinlist := []byte("")
	if 0 != len(userInfo.OwnerUinArr) {
		strOwnerUinlist, _ = json.Marshal(userInfo.OwnerUinArr)
	}

	cookielanguage, _ := c.Cookie("blueking_language")
	session := sessions.Default(c)

	session.Set(common.WEBSessionUinKey, userInfo.UserName)
	session.Set(common.WEBSessionChineseNameKey, userInfo.ChName)
	session.Set(common.WEBSessionPhoneKey, userInfo.Phone)
	session.Set(common.WEBSessionEmailKey, userInfo.Email)
	session.Set(common.WEBSessionRoleKey, userInfo.Role)
	session.Set(common.HTTPCookieBKToken, userInfo.BkToken)
	session.Set(common.WEBSessionOwnerUinKey, userInfo.OnwerUin)
	session.Set(common.WEBSessionAvatarUrlKey, userInfo.AvatarUrl)
	session.Set(common.WEBSessionOwnerUinListeKey, string(strOwnerUinlist))
	if userInfo.MultiSupplier {
		session.Set(common.WEBSessionMultiSupplierKey, common.LoginSystemMultiSupplierTrue)
	} else {
		session.Set(common.WEBSessionMultiSupplierKey, common.LoginSystemMultiSupplierFalse)
	}

	session.Set(webCommon.IsSkipLogin, "0")
	if "" != cookielanguage {
		session.Set(common.WEBSessionLanguageKey, cookielanguage)
	} else {
		session.Set(common.WEBSessionLanguageKey, userInfo.Language)
	}
	session.Save()
	return true
}

// GetUserList get user list from paas
func (m *publicUser) GetUserList(c *gin.Context) (int, interface{}) {

	user := plugins.CurrentPlugin(c)
	config := make(map[string]string)
	config["site.bk_account_url"] = m.config.Site.AccountUrl
	userList, err := user.GetUserList(c, config)
	rspBody := metadata.LonginSystemUserListResult{}
	if nil != err {
		rspBody.Code = common.CCErrCommHTTPDoRequestFailed
		rspBody.ErrMsg = err.Error()
		rspBody.Result = false
	}
	rspBody.Result = true
	rspBody.Data = userList
	return 200, rspBody
}

func (m *publicUser) GetLoginUrl(c *gin.Context) string {
	config := make(map[string]string)
	config["site.bk_login_url"] = m.config.Site.BkLoginUrl
	config["site.domain_url"] = m.config.Site.DomainUrl
	config["site.app_code"] = m.config.Site.AppCode
	params := new(metadata.LogoutRequestParams)
	err := json.NewDecoder(c.Request.Body).Decode(params)
	if nil != err || (common.LogoutHTTPSchemeHTTP != params.HTTPScheme && common.LogoutHTTPSchemeHTTPS != params.HTTPScheme) {
		params.HTTPScheme, err = c.Cookie(common.LogoutHTTPSchemeCookieKey)
		if nil != err || (common.LogoutHTTPSchemeHTTP != params.HTTPScheme && common.LogoutHTTPSchemeHTTPS != params.HTTPScheme) {
			params.HTTPScheme = common.LogoutHTTPSchemeHTTP
		}
	}
	user := plugins.CurrentPlugin(c)
	return user.GetLoginUrl(c, config, params)

}
