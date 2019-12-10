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

package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/web_server/middleware/user"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
)

const BkAccountUrl = "site.bk_account_url"

type userResult struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Code    string      `json:"code"`
	Result  bool        `json:"result"`
}

type userDataResult struct {
	Message string      `json:"bk_error_msg"`
	Data    interface{} `json:"data"`
	Code    string      `json:"bk_error_code"`
	Result  bool        `json:"result"`
}

var getUserFailData = userDataResult{
	Result:  false,
	Message: "get user list false",
	Code:    "",
	Data:    nil,
}

// GetUserList get user list
func (s *Service) GetUserList(c *gin.Context) {
	user := user.NewUser(*s.Config, s.Engine, s.CacheCli, s.VersionPlg)
	code, data := user.GetUserList(c)
	c.JSON(code, data)
	return
}

func (s *Service) UpdateUserLanguage(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	session := sessions.Default(c)
	language := c.Param("language")

	session.Set("language", language)
	err := session.Save()

	if nil != err {
		blog.Errorf("user update language error:%s, rid: %s", err.Error(), rid)
		c.JSON(200, userDataResult{
			Result:  false,
			Message: "user update language error",
			Code:    fmt.Sprintf("%d", common.CCErrCommHTTPDoRequestFailed),
			Data:    nil,
		})
		return
	}

	c.SetCookie("blueking_language", language, 0, "/", "", false, true)

	c.JSON(200, userDataResult{
		Result:  true,
		Message: "",
		Code:    "00",
		Data:    nil,
	})
	return
}

func (s *Service) UserInfo(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	session := sessions.Default(c)
	resultData := metadata.LoginUserInfoResult{}
	resultData.Result = true
	uin, ok := session.Get(common.WEBSessionUinKey).(string)
	if ok {
		resultData.Data.UserName = uin
	}
	name, ok := session.Get(common.WEBSessionChineseNameKey).(string)
	if ok {
		resultData.Data.ChName = name
	}
	ownerUin, ok := session.Get(common.WEBSessionOwnerUinKey).(string)
	if ok {
		resultData.Data.OnwerUin = ownerUin
	}
	strOwnerUinList, ok := session.Get(common.WEBSessionOwnerUinListeKey).(string)
	if ok {
		ownerUinList := make([]metadata.LoginUserInfoOwnerUinList, 0)
		err := json.Unmarshal([]byte(strOwnerUinList), &ownerUinList)
		if nil != err {
			blog.Errorf("[UserInfo] json unmarshal error:%s, rid: %s", err.Error(), rid)
		} else {
			resultData.Data.OwnerUinArr = ownerUinList
		}
	}
	avatarUrl, ok := session.Get(common.WEBSessionAvatarUrlKey).(string)
	if ok {
		resultData.Data.AvatarUrl = avatarUrl
	}
	iultiSupplier, ok := session.Get(common.WEBSessionMultiSupplierKey).(string)
	if ok && common.LoginSystemMultiSupplierTrue == iultiSupplier {
		resultData.Data.MultiSupplier = true // true
	} else {
		resultData.Data.MultiSupplier = false // true
	}

	c.JSON(200, resultData)
	return
}

func (s *Service) UpdateSupplier(c *gin.Context) {

	rid := util.GetHTTPCCRequestID(c.Request.Header)
	session := sessions.Default(c)

	strOwnerUinList, ok := session.Get(common.WEBSessionOwnerUinListeKey).(string)
	if "" == strOwnerUinList {
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommNotFound,
			ErrMsg: "not found",
		})
		return
	}
	ownerUinList := make([]metadata.LoginUserInfoOwnerUinList, 0)
	if ok {
		err := json.Unmarshal([]byte(strOwnerUinList), &ownerUinList)
		if nil != err {
			blog.Errorf("[UserInfo] json unmarshal error:%s, rid: %s", err.Error(), rid)
			c.JSON(http.StatusBadRequest, metadata.BaseResp{
				Result: false,
				Code:   common.CCErrCommJSONUnmarshalFailed,
				ErrMsg: "json unmarshal error",
			})
			return
		}
	}

	ownerID := c.Param("id")
	var supplier *metadata.LoginUserInfoOwnerUinList
	for idx, row := range ownerUinList {
		if row.OwnerID == ownerID {
			supplier = &ownerUinList[idx]
		}
	}

	if nil == supplier {
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommNotFound,
			ErrMsg: "not found",
		})
		return
	}
	session.Set(common.WEBSessionOwnerUinKey, supplier.OwnerID)
	session.Set(common.WEBSessionRoleKey, strconv.FormatInt(supplier.Role, 10))
	session.Set(common.WEBSessionSupplierID, strconv.FormatInt(supplier.SupplierID, 10))
	if err := session.Save(); err != nil {
		blog.Errorf("save session failed, err: %+v, rid: %s", err, rid)
	}
	ret := metadata.LoginChangeSupplierResult{}
	ret.Result = true
	ret.Data.ID = ownerID

	c.JSON(200, ret)
	return
}

// UserDetail 用户信息查询接口（根据用户名查询用户的详细信息，比如中文名，电话之类的）
// Note: 该接口仅用于前端页面将用户英文名转成中文名，因此，任何出错的情况均直接忽略错误，并返回空
func (s *Service) UserDetail(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)

	result := struct {
		Users    []metadata.LoginSystemUserInfo `json:"users"`
		Count    int                            `json:"count"`
		Next     interface{}                    `json:"next"`
		Previous interface{}                    `json:"previous"`
	}{}
	result.Users = make([]metadata.LoginSystemUserInfo, 0)

	// construct request url
	esbUrl, ok := s.Config.ConfigMap["esb.addr"]
	if ok == false {
		blog.Errorf("UserDetail failed, esb.addr not set, rid: %s", rid)
		c.JSON(http.StatusOK, metadata.NewSuccessResponse(result))
		return
	}
	targetUrl, err := url.Parse(esbUrl)
	if err != nil {
		blog.Errorf("UserDetail failed, parse login url failed, err: %+v, rid: %s", err, rid)
		c.JSON(http.StatusOK, metadata.NewSuccessResponse(result))
		return
	}
	frontUrl, err := url.Parse(c.Request.RequestURI)
	if err != nil {
		blog.Errorf("UserDetail failed, parse front url failed, err: %+v, rid: %s", err, rid)
		c.JSON(http.StatusOK, metadata.NewSuccessResponse(result))
		return
	}
	targetUrl.Path = "/o/bk_user_manage/api/v2/profiles/"
	targetUrl.RawQuery = frontUrl.RawQuery
	requestUrl := targetUrl.String()
	blog.V(5).Infof("request user detail by url: %s, rid: %s", requestUrl, rid)

	// do http request
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	rq := &http.Request{
		Method: c.Request.Method,
		URL:    targetUrl,
		Header: c.Request.Header,
	}
	rq.Header.Add("Accept", "application/json; nested=true")
	response, err := client.Do(rq)
	if err != nil {
		blog.Errorf("UserDetail failed, http get failed, err: %+v, rid: %s", err, rid)
		c.JSON(http.StatusOK, metadata.NewSuccessResponse(result))
		return
	}

	// decode response body
	data, err := ioutil.ReadAll(response.Body)
	blog.V(5).Infof("response body: %s, rid: %s", data, rid)
	if err != nil {
		blog.Errorf("UserDetail failed, read response body failed, err: %+v, rid: %s", err, rid)
		c.JSON(http.StatusOK, metadata.NewSuccessResponse(result))
		return
	}
	if data == nil {
		blog.Errorf("UserDetail failed, response body empty, err: %+v, rid: %s", err, rid)
		c.JSON(http.StatusOK, metadata.NewSuccessResponse(result))
		return
	}
	type UserDetail struct {
		Username    string `json:"username"`
		Qq          string `json:"qq"`
		UserType    string `json:"user_type"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
	}
	type ResponseData struct {
		Count    int          `json:"count"`
		Next     interface{}  `json:"next"`
		Previous interface{}  `json:"previous"`
		Results  []UserDetail `json:"results"`
	}
	responseData := ResponseData{}
	if err = json.Unmarshal(data, &responseData); err != nil {
		blog.Errorf("UserDetail failed, decode response data into struct failed, data: %s, err: %+v, rid: %s", data, err, rid)
		c.JSON(http.StatusOK, metadata.NewSuccessResponse(result))
		return
	}
	for _, item := range responseData.Results {
		result.Users = append(result.Users, metadata.LoginSystemUserInfo{
			CnName: item.DisplayName,
			EnName: item.Username,
		})
	}
	result.Count = responseData.Count
	result.Previous = responseData.Previous
	result.Next = responseData.Next

	if blog.V(5) {
		blog.InfoJSON("login url: %s, rid: %s", targetUrl, rid)
	}
	c.JSON(http.StatusOK, metadata.NewSuccessResponse(result))
	return
}
