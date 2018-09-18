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

package self

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/metadata"
	"configcenter/src/web_server/application/middleware/user/plugins/register"
	webCommon "configcenter/src/web_server/common"
)

func init() {
	plugin := &metadata.LoginPluginInfo{
		Name:       "blueking login system",
		Version:    common.BKDefaultLoginUserPluginVersion,
		HandleFunc: &user{},
	}
	register.RegisterPlugin(plugin) //("blueking login system", "self", "")
}

type loginResultData struct {
	UserName string `json:"username"`
	ChName   string `json:"chname"`
	Phone    string `json:"Phone"`
	Email    string `json:"email"`
	Role     int64  `json:"role"`
	Language string `json:"language"`
	OwnerUin string `json:"owner_uin"`
}

type loginResult struct {
	Message string
	Code    string
	Result  bool
	Data    *loginResultData
}

type userListResult struct {
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

type user struct {
}

// LoginUser  user login
func (m *user) LoginUser(c *gin.Context, config map[string]string, isMultiOwner bool) (user *metadata.LoginUserInfo, loginSucc bool) {
	bk_token, err := c.Cookie("bk_token")
	if nil != err {
		return nil, false
	}
	if nil != err || 0 == len(bk_token) {
		return nil, false
	}
	checkUrl, ok := config["site.bk_login_url"]
	if !ok {
		blog.Errorf("get login url config item not found")
		return nil, false
	}
	loginURL := checkUrl + bk_token
	httpCli := httpclient.NewHttpClient()
	httpCli.SetTimeOut(30 * time.Second)
	loginResultByteArr, err := httpCli.GET(loginURL, nil, nil)

	if nil != err {
		blog.Errorf("get user info return error: %v", err)
		return nil, false
	}
	blog.V(3).Infof("get user info cond %v, return: %s ", string(loginURL), string(loginResultByteArr))
	var resultData loginResult
	err = json.Unmarshal(loginResultByteArr, &resultData)
	if nil != err {
		blog.Errorf("get user info json error: %v", err)
		return nil, false
	}

	if !resultData.Result {
		blog.Errorf("get user info return error , error code: %s, error message: %s", resultData.Code, resultData.Message)
		return nil, false
	}

	userDetail := resultData.Data
	user = &metadata.LoginUserInfo{
		UserName: userDetail.UserName,
		ChName:   userDetail.ChName,
		Phone:    userDetail.Phone,
		Email:    userDetail.Email,
		Role:     userDetail.Role,
		BkToken:  bk_token,
		OnwerUin: userDetail.OwnerUin,
		IsOwner:  false,
		Language: userDetail.Language,
	}
	return user, true
}

// GetUserList get user list from paas
func (m *user) GetUserList(c *gin.Context, config map[string]string) ([]*metadata.LoginSystemUserInfo, error) {

	accountURL, ok := config["site.bk_account_url"]
	if !ok {
		return nil, fmt.Errorf("config site.bk_account_url not found")
	}
	session := sessions.Default(c)
	skiplogin := session.Get(webCommon.IsSkipLogin)
	skiplogins, ok := skiplogin.(string)
	if ok && "1" == skiplogins {
		blog.Info("use skip login flag: %v", skiplogin)
		adminData := []*metadata.LoginSystemUserInfo{
			&metadata.LoginSystemUserInfo{
				CnName: "admin",
				EnName: "admin",
			},
		}

		return adminData, nil
	}

	token := session.Get("bk_token")
	getURL := fmt.Sprintf(accountURL, token)
	httpClient := httpclient.NewHttpClient()

	reply, err := httpClient.GET(getURL, nil, nil)

	if nil != err {
		blog.Errorf("get user list error：%v", err)
		return nil, fmt.Errorf("http do error:%s", err.Error())
	}
	blog.Info("get user list url: %s, return：%s", getURL, reply)
	var result userListResult
	err = json.Unmarshal([]byte(reply), &result)
	if nil != err || false == result.Result {
		blog.Errorf("get user list error：%v, error code:%d, error messsage:", err, result.Code, result.Message)
		return nil, fmt.Errorf("get user list reply error")
	}
	userListArr := make([]*metadata.LoginSystemUserInfo, 0)
	for _, user := range result.Data {
		cellData := make(map[string]interface{})
		cellData["chinese_name"] = user.Chname
		cellData["english_name"] = user.UserName
		userListArr = append(userListArr, &metadata.LoginSystemUserInfo{
			CnName: user.Chname,
			EnName: user.UserName,
		})
	}

	return userListArr, nil
}
