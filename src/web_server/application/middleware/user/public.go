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
	"fmt"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/web_server/application/middleware/types"
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
	bk_token, err := c.Cookie("bk_token")
	if nil != err {
		return false
	}
	if nil != err || 0 == len(bk_token) {
		return false
	}
	loginURL := checkUrl + bk_token
	httpCli := httpclient.NewHttpClient()
	httpCli.SetTimeOut(30 * time.Second)
	loginResult, err := httpCli.GET(loginURL, nil, nil)

	if nil != err {
		blog.Error("get user info return error: %v", err)
		return false
	}
	blog.Infof("get user info cond %v, return: %s ", string(loginURL), string(loginResult))
	var resultData types.LoginResult
	err = json.Unmarshal([]byte(loginResult), &resultData)
	if nil != err {
		blog.Error("get user info json error: %v", err)
		return false
	}
	userInfo, ok := resultData.Data.(map[string]interface{})
	if false == ok {
		blog.Error("get user info decode error: %v", err)
		return false
	}
	userName, ok := userInfo["username"]
	if false == ok {
		blog.Error("get user info username error: %v", err)
		return false
	}
	chName, ok := userInfo["chname"]
	if false == ok {
		blog.Error("get user info chname error: %v", err)
		return false
	}
	phone, ok := userInfo["phone"]
	if false == ok {
		blog.Error("get user info phone error: %v", err)
		return false
	}
	email, ok := userInfo["email"]
	if false == ok {
		blog.Error("get user info email error: %v", err)
		return false
	}
	role, ok := userInfo["role"]
	if false == ok {
		blog.Error("get user info role error: %v", err)
		return false
	}

	language, ok := userInfo["language"]
	if false == ok {
		blog.Error("get language info role error: %v", err)
	}
	ownerID := common.BKDefaultOwnerID
	if true == isMultiOwner {
		ownerID, ok = userInfo["owner_uin"].(string)
		if false == ok {
			blog.Error("get owner_uin info role error: %v", err)
			return false
		}
		_, ok = userName.(string)
		if false == ok {
			blog.Error("get username info role error: %v", err)
			return false
		}
		_, ok = language.(string)
		if false == ok {
			blog.Error("get language info role error: %v", err)
			return false
		}

		err := NewOwnerManager(userName.(string), ownerID, language.(string)).InitOwner()
		if nil != err {
			blog.Error("InitOwner error: %v", err)
			return false
		}
	}

	cookielanguage, _ := c.Cookie("blueking_language")
	session := sessions.Default(c)
	session.Set("userName", userName)
	session.Set("chName", chName)
	session.Set("phone", phone)
	session.Set("email", email)
	session.Set("role", role)
	session.Set("bk_token", bk_token)
	session.Set("owner_uin", ownerID)
	session.Set(webCommon.IsSkipLogin, "0")
	if "" != cookielanguage {
		session.Set("language", cookielanguage)
	} else {
		session.Set("language", language)
	}
	session.Save()
	return true
}

// GetUserList get user list from paas
func (m *publicUser) GetUserList(c *gin.Context, accountURL string) (int, interface{}) {
	session := sessions.Default(c)
	skiplogin := session.Get(webCommon.IsSkipLogin)
	skiplogins, ok := skiplogin.(string)
	if ok && "1" == skiplogins {
		blog.Info("use skip login flag: %v", skiplogin)
		adminData := []map[string]interface{}{
			{
				"chinese_name": "admin",
				"english_name": "admin",
			},
		}

		return 200, map[string]interface{}{
			"result":        true,
			"bk_error_msg":  "get user list ok",
			"bk_error_code": "00",
			"data":          adminData,
		}
	}

	token := session.Get("bk_token")
	getURL := fmt.Sprintf(accountURL, token)
	httpClient := httpclient.NewHttpClient()

	reply, err := httpClient.GET(getURL, nil, nil)

	if nil != err {
		blog.Error("get user list error：%v", err)
		return 200, getUserFailData
	}
	blog.Info("get user list urlL: %s, return：%s", getURL, reply)
	var result userResult
	info := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(reply), &result)
	if nil != err || false == result.Result {
		return 200, getUserFailData
	}
	for _, user := range result.Data {
		cellData := make(map[string]interface{})
		cellData["chinese_name"] = user.Chname
		cellData["english_name"] = user.UserName
		info = append(info, cellData)
	}
	return 200, map[string]interface{}{
		"result":        true,
		"bk_error_msg":  "get user list ok",
		"bk_error_code": "00",
		"data":          info,
	}
}
