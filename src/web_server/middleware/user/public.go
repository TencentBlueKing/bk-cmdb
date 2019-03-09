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
	"plugin"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/web_server/app/options"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/user/plugins"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
	redis "gopkg.in/redis.v5"
)

type publicUser struct {
	config   options.Config
	engine   *backbone.Engine
	cacheCli *redis.Client
	loginPlg *plugin.Plugin
}

// LoginUser  user login
func (m *publicUser) LoginUser(c *gin.Context) bool {

	isMultiOwner := false
	loginSucc := false
	var userInfo *metadata.LoginUserInfo
	multipleOwner := m.config.Session.MultipleOwner
	if common.LoginSystemMultiSupplierTrue == multipleOwner {
		isMultiOwner = true
	}

	if nil == m.loginPlg {
		user := plugins.CurrentPlugin(c, m.config.LoginVersion)
		userInfo, loginSucc = user.LoginUser(c, m.config.ConfigMap, isMultiOwner)
	} else {

		loginUserFunc, err := m.loginPlg.Lookup("LoginUser")

		if nil != err {
			blog.Error("look login func error")
			return false

		}
		userInfo, loginSucc = loginUserFunc.(func(c *gin.Context, config map[string]string, isMultiOwner bool) (user *metadata.LoginUserInfo, loginSucc bool))(c, m.config.ConfigMap, isMultiOwner)

	}

	if !loginSucc {
		return false
	}
	if true == isMultiOwner || true == userInfo.MultiSupplier {
		ownerM := NewOwnerManager(userInfo.UserName, userInfo.OnwerUin, userInfo.Language)
		ownerM.CacheCli = m.cacheCli
		ownerM.Engine = m.engine
		ownerM.SetHttpHeader(common.BKHTTPSupplierID, strconv.FormatInt(userInfo.SupplierID, 10))
		err := ownerM.InitOwner()
		if nil != err {
			blog.Errorf("InitOwner error: %v", err)
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
	session.Set(common.WEBSessionSupplierID, strconv.FormatInt(userInfo.SupplierID, 10))
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
	var err error
	var userList []*metadata.LoginSystemUserInfo
	rspBody := metadata.LonginSystemUserListResult{}
	rspBody.Result = false
	if nil == m.loginPlg {
		user := plugins.CurrentPlugin(c, m.config.LoginVersion)
		userList, err = user.GetUserList(c, m.config.ConfigMap)
	} else {

		getUserListFunc, err := m.loginPlg.Lookup("GetUserList")

		if nil != err {
			blog.Error("look get user list error")
			rspBody.Code = common.CCErrCommHTTPDoRequestFailed
			rspBody.ErrMsg = err.Error()
			return 200, rspBody

		}
		userList, err = getUserListFunc.(func(c *gin.Context, config map[string]string) ([]*metadata.LoginSystemUserInfo, error))(c, m.config.ConfigMap)

	}
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

	params := new(metadata.LogoutRequestParams)
	err := json.NewDecoder(c.Request.Body).Decode(params)
	if nil != err || (common.LogoutHTTPSchemeHTTP != params.HTTPScheme && common.LogoutHTTPSchemeHTTPS != params.HTTPScheme) {
		params.HTTPScheme, err = c.Cookie(common.LogoutHTTPSchemeCookieKey)
		if nil != err || (common.LogoutHTTPSchemeHTTP != params.HTTPScheme && common.LogoutHTTPSchemeHTTPS != params.HTTPScheme) {
			params.HTTPScheme = common.LogoutHTTPSchemeHTTP
		}
	}

	if nil == m.loginPlg {
		user := plugins.CurrentPlugin(c, m.config.LoginVersion)
		return user.GetLoginUrl(c, m.config.ConfigMap, params)
	} else {

		getLoginUrlFunc, err := m.loginPlg.Lookup("GetLoginUrl")

		if nil != err {
			blog.Error("look get url func error")
			return ""

		}
		return getLoginUrlFunc.(func(c *gin.Context, config map[string]string, input *metadata.LogoutRequestParams) string)(c, m.config.ConfigMap, params)

	}

}
