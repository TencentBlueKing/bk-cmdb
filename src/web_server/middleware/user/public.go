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
	"net/http"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/web_server/app/options"
	"configcenter/src/web_server/middleware/user/plugins"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
	"gopkg.in/redis.v5"
)

type publicUser struct {
	config   options.Config
	engine   *backbone.Engine
	cacheCli *redis.Client
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
		ownerM.SetHttpHeader(common.BKHTTPSupplierID, strconv.FormatInt(userInfo.SupplierID, 10))
		err := ownerM.InitOwner()
		if nil != err {
			blog.Errorf("InitOwner error: %v, rid: %s", err, rid)
			return false
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
	session.Set(common.WEBSessionSupplierID, strconv.FormatInt(userInfo.SupplierID, 10))
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

// GetUserList get user list from PaaS
func (m *publicUser) GetUserList(c *gin.Context) (int, interface{}) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	var err error
	var userList []*metadata.LoginSystemUserInfo
	rspBody := metadata.LonginSystemUserListResult{}
	rspBody.Result = true
	query := c.Request.URL.Query()
	params := make(map[string]string)
	for key, values := range query {
		params[key] = strings.Join(values, ";")
	}
	user := plugins.CurrentPlugin(c, m.config.LoginVersion)
	userList, err = user.GetUserList(c, m.config.ConfigMap, params)
	if nil != err {
		blog.Error("GetUserList failed, err: %+v, rid: %s", err, rid)
		rspBody.Code = common.CCErrCommHTTPDoRequestFailed
		rspBody.ErrMsg = err.Error()
		rspBody.Result = false
		return http.StatusInternalServerError, rspBody
	}
	rspBody.Result = true
	rspBody.Data = userList
	return http.StatusOK, rspBody
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

// GetDepartment get department info from PaaS
func (m *publicUser) GetDepartment(c *gin.Context) (int, interface{}) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	var err error
	var departments *metadata.DepartmentData
	rspBody := metadata.DepartmentResult{}
	rspBody.Result = true
	user := plugins.CurrentPlugin(c, m.config.LoginVersion)
	departments, err = user.GetDepartment(c, m.config.ConfigMap)
	if nil != err {
		blog.Error("GetDepartment failed, err: %+v, rid: %s", err, rid)
		rspBody.Code = common.CCErrCommHTTPDoRequestFailed
		rspBody.ErrMsg = err.Error()
		rspBody.Result = false
		return http.StatusInternalServerError, rspBody
	}
	rspBody.Result = true
	rspBody.Data = departments
	return http.StatusOK, rspBody
}

// GetDepartmentProfile get department profile from PaaS
func (m *publicUser) GetDepartmentProfile(c *gin.Context) (int, interface{}) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	var err error
	var departmentprofile *metadata.DepartmentProfileData
	rspBody := metadata.DepartmentProfileResult{}
	rspBody.Result = true
	user := plugins.CurrentPlugin(c, m.config.LoginVersion)
	departmentprofile, err = user.GetDepartmentProfile(c, m.config.ConfigMap)
	if nil != err {
		blog.Error("GetDepartmentProfile failed, err: %+v, rid: %s", err, rid)
		rspBody.Code = common.CCErrCommHTTPDoRequestFailed
		rspBody.ErrMsg = err.Error()
		rspBody.Result = false
		return http.StatusInternalServerError, rspBody
	}
	rspBody.Result = true
	rspBody.Data = departmentprofile
	return http.StatusOK, rspBody
}
