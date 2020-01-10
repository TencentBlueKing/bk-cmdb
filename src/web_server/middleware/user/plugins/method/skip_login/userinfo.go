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

package skip_login

import (
	"encoding/json"
	"fmt"

	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/metadata"
	commonutil "configcenter/src/common/util"
	"configcenter/src/thirdpartyclient/esbserver"
	"configcenter/src/thirdpartyclient/esbserver/esbutil"
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/user/plugins/manager"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
)

func init() {
	plugin := &metadata.LoginPluginInfo{
		Name:       "skip login system",
		Version:    "skip-login",
		HandleFunc: &user{},
	}
	manager.RegisterPlugin(plugin)
}

type loginResultData struct {
	UserName string `json:"username"`
	ChName   string `json:"chname"`
	Phone    string `json:"Phone"`
	Email    string `json:"email"`
	Role     string `json:"role"`
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

	session := sessions.Default(c)

	cookieOwnerID, err := c.Cookie(common.BKHTTPOwnerID)
	if "" == cookieOwnerID || nil != err {
		c.SetCookie(common.BKHTTPOwnerID, common.BKDefaultOwnerID, 0, "/", "", false, false)
		session.Set(common.WEBSessionOwnerUinKey, cookieOwnerID)
	} else if cookieOwnerID != session.Get(common.WEBSessionOwnerUinKey) {
		session.Set(common.WEBSessionOwnerUinKey, cookieOwnerID)
	}

	user = &metadata.LoginUserInfo{
		UserName: "admin",
		ChName:   "admin",
		Phone:    "",
		Email:    "blueking",
		Role:     "",
		BkToken:  "",
		OnwerUin: "0",
		IsOwner:  false,
		Language: webCommon.GetLanguageByHTTPRequest(c),
	}
	return user, true
}

func (m *user) getEsbClient(config map[string]string) (esbserver.EsbClientInterface, error) {
	esbAddr, addrOk := config["esb.addr"]
	esbAppCode, appCodeOk := config["esb.appCode"]
	esbAppSecret, appSecretOk := config["esb.appSecret"]
	if addrOk == false || appCodeOk == false || appSecretOk == false {
		return nil, fmt.Errorf("esb config not found or incomplete, %+v", config)
	}
	tlsConfig, err := util.NewTLSClientConfigFromConfig("esb", config)
	if err != nil {
		return nil, fmt.Errorf("parse esb tls config failed, config: %+v, err: %+v", config, err)
	}
	apiMachineryConfig := &util.APIMachineryConfig{
		QPS:       1000,
		Burst:     1000,
		TLSConfig: &tlsConfig,
	}
	defaultCfg := &esbutil.EsbConfig{
		Addrs:     esbAddr,
		AppCode:   esbAppCode,
		AppSecret: esbAppSecret,
	}
	esbSrv, err := esbserver.NewEsb(apiMachineryConfig, nil, defaultCfg, nil)
	if err != nil {
		return nil, fmt.Errorf("create esb client failed. err: %v", err)
	}
	return esbSrv, nil
}

// GetUserList get user list from paas
func (m *user) GetUserList(c *gin.Context, config map[string]string) ([]*metadata.LoginSystemUserInfo, error) {
	rid := commonutil.GetHTTPCCRequestID(c.Request.Header)
	accountURL, ok := config["site.bk_account_url"]
	if !ok {
		// try to use esb user list api
		esbClient, err := m.getEsbClient(config)
		if err != nil {
			blog.Warnf("get esb client failed, err: %+v, rid: %s", err, rid)
			return nil, fmt.Errorf("config site.bk_account_url not found")
		}
		result, err := esbClient.User().ListUsers(c.Request.Context(), c.Request.Header)
		if err != nil {
			blog.Warnf("get users by esb client failed, http failed, err: %+v, rid: %s", err, rid)
			return nil, fmt.Errorf("get users by esb client failed, http failed")
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
	session := sessions.Default(c)
	skipLogin := session.Get(webCommon.IsSkipLogin)
	skipLogins, ok := skipLogin.(string)
	if ok && "1" == skipLogins {
		blog.V(5).Infof("use skip login flag: %v, rid: %s", skipLogin, rid)
		adminData := []*metadata.LoginSystemUserInfo{
			{
				CnName: "admin",
				EnName: "admin",
			},
		}

		return adminData, nil
	}

	token := session.Get(common.HTTPCookieBKToken)
	getURL := fmt.Sprintf(accountURL, token)
	httpClient := httpclient.NewHttpClient()

	if err := httpClient.SetTlsNoVerity(); err != nil {
		blog.Warnf("httpClient.SetTlsNoVerity failed, err: %s, rid: %s", err.Error(), rid)
	}
	reply, err := httpClient.GET(getURL, nil, nil)

	if nil != err {
		blog.Errorf("get user list error：%v, rid: %s", err, rid)
		return nil, fmt.Errorf("http do error:%s", err.Error())
	}
	blog.V(5).Infof("get user list url: %s, return：%s, rid: %s", getURL, reply, rid)
	var result userListResult
	err = json.Unmarshal([]byte(reply), &result)
	if nil != err || false == result.Result {
		blog.Errorf("get user list error：%v, error code:%s, error message: %s, rid: %s", err, result.Code, result.Message, rid)
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

func (m *user) GetLoginUrl(c *gin.Context, config map[string]string, input *metadata.LogoutRequestParams) string {
	var ok bool
	var loginURL string
	var siteURL string

	if common.LogoutHTTPSchemeHTTPS == input.HTTPScheme {
		loginURL, ok = config["site.bk_https_login_url"]
	} else {
		loginURL, ok = config["site.bk_login_url"]
	}
	if !ok {
		loginURL = ""
	}
	if common.LogoutHTTPSchemeHTTPS == input.HTTPScheme {
		siteURL, ok = config["site.https_domain_url"]
	} else {
		siteURL, ok = config["site.domain_url"]
	}
	if !ok {
		siteURL = ""
	}

	appCode, ok := config["site.app_code"]
	if !ok {
		appCode = ""
	}
	loginURL = fmt.Sprintf(loginURL, appCode, fmt.Sprintf("%s%s", siteURL, c.Request.URL.String()))
	return loginURL
}
