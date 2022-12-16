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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	apiutil "configcenter/src/apimachinery/util"
	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/metadata"
	"configcenter/src/common/resource/esb"
	"configcenter/src/common/util"
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

type user struct{}

// LoginUser user login
func (m *user) LoginUser(c *gin.Context, config map[string]string, isMultiOwner bool) (user *metadata.LoginUserInfo,
	loginSucc bool) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)

	bkTokens := getBkTokens(c)
	if len(bkTokens) == 0 {
		blog.Infof("LoginUser failed, bk_token empty, rid: %s", rid)
		return nil, false
	}

	checkUrl, err := cc.String("webServer.site.checkUrl")
	if err != nil {
		blog.Errorf("get login url config item not found, rid: %s", rid)
		return nil, false
	}

	var resultData loginResult
	httpCli := httpclient.NewHttpClient()
	httpCli.SetTimeOut(30 * time.Second)

	tlsConf, err := apiutil.NewTLSClientConfigFromConfig("webServer.site.paas.tls")
	if err != nil {
		blog.Errorf("get tls config error, err: %v, rid: %s", err, rid)
		return nil, false
	}

	if err := m.setTLSConf(&tlsConf, httpCli, rid); err != nil {
		blog.Errorf("set tls config error, err: %v, rid: %s", err, rid)
		return nil, false
	}

	for _, bkToken := range bkTokens {
		loginURL := checkUrl + bkToken
		loginResultByteArr, err := httpCli.GET(loginURL, nil, nil)
		if err != nil {
			blog.Errorf("get user info return error: %v, rid: %s", err, rid)
			continue
		}
		blog.V(3).Infof("get user info url: %s, result: %s, rid: %s", loginURL, string(loginResultByteArr), rid)

		err = json.Unmarshal(loginResultByteArr, &resultData)
		if err != nil {
			blog.Errorf("fail to unmarshal data: %s, error: %v, rid: %s", string(loginResultByteArr), err, rid)
			return nil, false
		}

		if !resultData.Result {
			blog.Errorf("get user info return error, error code: %s, error message: %s, rid: %s", resultData.Code,
				resultData.Message, rid)
			continue
		}

		user = setUser(resultData, bkToken)
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

// setUser get userInfo from resultData
func setUser(resultData loginResult, bkToken string) (user *metadata.LoginUserInfo) {
	userDetail := resultData.Data
	if len(userDetail.OwnerUin) == 0 {
		userDetail.OwnerUin = common.BKDefaultOwnerID
	}

	user = &metadata.LoginUserInfo{
		UserName: userDetail.UserName,
		ChName:   userDetail.ChName,
		Phone:    userDetail.Phone,
		Email:    userDetail.Email,
		Role:     userDetail.Role,
		BkToken:  bkToken,
		OnwerUin: userDetail.OwnerUin,
		IsOwner:  false,
		Language: userDetail.Language,
	}
	return user
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
	rid := util.GetHTTPCCRequestID(c.Request.Header)
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

func (m *user) setTLSConf(tlsConf *apiutil.TLSClientConfig, httpCli *httpclient.HttpClient, rid string) error {
	if tlsConf != nil && len(tlsConf.CAFile) != 0 && len(tlsConf.CertFile) != 0 && len(tlsConf.KeyFile) != 0 {
		if err := httpCli.SetTLSVerify(tlsConf); err != nil {
			return err
		}
		return nil
	}

	if err := httpCli.SetTlsNoVerity(); err != nil {
		blog.Warnf("httpCli.SetTlsNoVerity failed, err: %+v, rid: %s", err, rid)
	}

	return nil
}
