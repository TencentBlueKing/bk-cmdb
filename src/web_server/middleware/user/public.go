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
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/web_server/app/options"
	"configcenter/src/web_server/middleware/user/plugins"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
)

type publicUser struct {
	config   options.Config
	engine   *backbone.Engine
	cacheCli redis.Client
}

// LoginUser  user login
func (m *publicUser) LoginUser(c *gin.Context) bool {
	rid := util.GetHTTPCCRequestID(c.Request.Header)

	isMultiOwner := false
	loginSuccess := false
	var userInfo *metadata.LoginUserInfo
	multipleOwner := m.config.Session.MultipleOwner
	if common.LoginSystemMultiSupplierTrue == multipleOwner {
		isMultiOwner = true
	}

	user := plugins.CurrentPlugin(c, m.config.LoginVersion)
	userInfo, loginSuccess = user.LoginUser(c, m.config.ConfigMap, isMultiOwner)

	if !loginSuccess {
		blog.Infof("login user with plugin failed, rid: %s", rid)
		return false
	}
	if true == isMultiOwner || true == userInfo.MultiSupplier {
		ownerM := NewOwnerManager(userInfo.UserName, userInfo.OnwerUin, userInfo.Language)
		ownerM.CacheCli = m.cacheCli
		ownerM.Engine = m.engine
		// 初始化失败，不影响登录
		_, err := ownerM.InitOwner()
		if nil != err {
			blog.ErrorJSON("init onwer resource pool failed, err:%s, user:%s, rid: %s", err, userInfo, rid)
		}
	}
	strOwnerUinList := []byte("")
	if 0 != len(userInfo.OwnerUinArr) {
		strOwnerUinList, _ = json.Marshal(userInfo.OwnerUinArr)
	}

	session := sessions.Default(c)

	session.Set(common.WEBSessionUinKey, userInfo.UserName)
	session.Set(common.WEBSessionChineseNameKey, userInfo.ChName)
	session.Set(common.WEBSessionPhoneKey, userInfo.Phone)
	session.Set(common.WEBSessionEmailKey, userInfo.Email)
	session.Set(common.WEBSessionRoleKey, userInfo.Role)
	session.Set(common.HTTPCookieBKToken, userInfo.BkToken)
	session.Set(common.WEBSessionOwnerUinKey, userInfo.OnwerUin)
	session.Set(common.WEBSessionAvatarUrlKey, userInfo.AvatarUrl)
	session.Set(common.WEBSessionOwnerUinListeKey, string(strOwnerUinList))
	if userInfo.MultiSupplier {
		session.Set(common.WEBSessionMultiSupplierKey, common.LoginSystemMultiSupplierTrue)
	} else {
		session.Set(common.WEBSessionMultiSupplierKey, common.LoginSystemMultiSupplierFalse)
	}

	if err := session.Save(); err != nil {
		blog.Warnf("save session failed, err: %s, rid: %s", err.Error(), rid)
	}
	return true
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

	user := plugins.CurrentPlugin(c, m.config.LoginVersion)
	return user.GetLoginUrl(c, m.config.ConfigMap, params)

}

func (m *publicUser) GetUserList(c *gin.Context) ([]*metadata.LoginSystemUserInfo, *errors.RawErrorInfo) {
	user := plugins.CurrentPlugin(c, m.config.LoginVersion)
	return user.GetUserList(c, m.config.ConfigMap)
}
