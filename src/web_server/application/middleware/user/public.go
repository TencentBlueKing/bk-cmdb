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
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/metadata"
	"configcenter/src/web_server/application/middleware/user/plugins"
	webCommon "configcenter/src/web_server/common"
)

type publicUser struct {
}

var getUserFailData = map[string]interface{}{
	"result":        false,
	"bk_error_msg":  "get user list false",
	"bk_error_code": "",
	"data":          nil,
}

type userResult struct {
	Message string     `json:"message"`
	Data    []userInfo `json:"data"`
	Code    string     `json:"code"`
	Result  bool       `json:"result"`
}

type userInfo struct {
	UserName string `json:"username"`
	QQ       string `json:"qq"`
	Role     string `json:"role"`
	Language string `json:"language"`
	Phone    string `json:"phone"`
	WxUserid string `json:"wx_userid"`
	Email    string `json:"email"`
	Chname   string `json:"chname"`
	TimeZone string `json:"time_zone"`
}

// LoginUser  user login
func (m *publicUser) LoginUser(c *gin.Context, checkUrl string, isMultiOwner bool) bool {

	ccapi := api.NewAPIResource()
	config, _ := ccapi.ParseConfig()
	user := plugins.CurrentPlugin(c)
	userInfo, loginSucc := user.LoginUser(c, config, isMultiOwner)
	if !loginSucc {
		return false
	}

	if true == isMultiOwner {
		err := NewOwnerManager(userInfo.UserName, userInfo.OnwerUin, userInfo.Language).InitOwner()
		if nil != err {
			blog.Error("InitOwner error: %v", err)
			return false
		}
	}

	cookielanguage, _ := c.Cookie("blueking_language")
	session := sessions.Default(c)
	session.Set("userName", userInfo.UserName)
	session.Set("chName", userInfo.ChName)
	session.Set("phone", userInfo.Phone)
	session.Set("email", userInfo.Email)
	session.Set("role", userInfo.Role)
	session.Set("bk_token", userInfo.BkToken)
	session.Set("owner_uin", userInfo.OnwerUin)
	session.Set(webCommon.IsSkipLogin, "0")
	if "" != cookielanguage {
		session.Set("language", cookielanguage)
	} else {
		session.Set("language", userInfo.Language)
	}
	session.Save()
	return true
}

// GetUserList get user list from paas
func (m *publicUser) GetUserList(c *gin.Context, accountURL string) (int, interface{}) {

	ccapi := api.NewAPIResource()
	config, _ := ccapi.ParseConfig()
	user := plugins.CurrentPlugin(c)
	userList, err := user.GetUserList(c, config)
	rspBody := metadata.LonginSystemUserListResult{}
	if nil != err {
		rspBody.Code = common.CCErrCommHTTPDoRequestFailed
		rspBody.ErrMsg = err.Error()
		rspBody.Result = false
	}
	rspBody.Data = userList
	return 200, rspBody
}
