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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/web_server/middleware/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type userDataResult struct {
	Message string      `json:"bk_error_msg"`
	Data    interface{} `json:"data"`
	Code    string      `json:"bk_error_code"`
	Result  bool        `json:"result"`
}

// GetUserList get user list
func (s *Service) GetUserList(c *gin.Context) {
	rid := httpheader.GetRid(c.Request.Header)
	rspBody := metadata.LonginSystemUserListResult{}

	userManger := user.NewUser(*s.Config, s.Engine, s.CacheCli)
	userList, rawErr := userManger.GetUserList(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(c.Request.Header))
	if rawErr != nil && rawErr.ErrCode != 0 {
		blog.Error("GetUserList failed, err: %v, rid: %s", rawErr.ToCCError(defErr), rid)
		rspBody.Code = rawErr.ErrCode
		rspBody.ErrMsg = rawErr.ToCCError(defErr).Error()
		rspBody.Result = false
		c.JSON(http.StatusInternalServerError, rspBody)
		return
	}

	rspBody.Result = true
	rspBody.Data = userList

	c.JSON(http.StatusOK, rspBody)
	return
}

// UpdateUserLanguage TODO
func (s *Service) UpdateUserLanguage(c *gin.Context) {
	rid := httpheader.GetRid(c.Request.Header)
	session := sessions.Default(c)
	language := c.Param("language")

	session.Set("language", language)
	err := session.Save()

	if nil != err {
		blog.Errorf("user update language error: %v, rid: %s", err, rid)
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

// UserInfo set user info
func (s *Service) UserInfo(c *gin.Context) {
	rid := httpheader.GetRid(c.Request.Header)
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
	tenantUin, ok := session.Get(common.WEBSessionTenantUinKey).(string)
	if ok {
		resultData.Data.TenantUin = tenantUin
	}
	strTenantUinList, ok := session.Get(common.WEBSessionTenantUinListeKey).(string)
	if ok {
		tenantUinList := make([]metadata.LoginUserInfoTenantUinList, 0)
		err := json.Unmarshal([]byte(strTenantUinList), &tenantUinList)
		if nil != err {
			blog.Errorf("[UserInfo] json unmarshal error: %v, rid: %s", err, rid)
		} else {
			resultData.Data.TenantUinArr = tenantUinList
		}
	}
	avatarUrl, ok := session.Get(common.WEBSessionAvatarUrlKey).(string)
	if ok {
		resultData.Data.AvatarUrl = avatarUrl
	}
	iultiTenant, ok := session.Get(common.WEBSessionMultiTenantKey).(string)
	if ok && common.LoginSystemMultiTenantTrue == iultiTenant {
		resultData.Data.MultiTenant = true // true
	} else {
		resultData.Data.MultiTenant = false // true
	}

	c.JSON(200, resultData)
	return
}

// UpdateTenant update tenant info
func (s *Service) UpdateTenant(c *gin.Context) {

	rid := httpheader.GetRid(c.Request.Header)
	session := sessions.Default(c)

	strTenantUinList, ok := session.Get(common.WEBSessionTenantUinListeKey).(string)
	if strTenantUinList == "" {
		blog.Errorf("session not tenant info, rid: %s", rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommNotFound,
			ErrMsg: "not found",
		})
		return
	}
	tenantUinList := make([]metadata.LoginUserInfoTenantUinList, 0)
	if ok {
		err := json.Unmarshal([]byte(strTenantUinList), &tenantUinList)
		if nil != err {
			blog.Errorf("[UserInfo] json unmarshal error: %v, rid: %s", err, rid)
			c.JSON(http.StatusBadRequest, metadata.BaseResp{
				Result: false,
				Code:   common.CCErrCommJSONUnmarshalFailed,
				ErrMsg: "json unmarshal error",
			})
			return
		}
	}

	tenantID := c.Param("id")
	var tenant *metadata.LoginUserInfoTenantUinList
	for idx, row := range tenantUinList {
		if row.TenantID == tenantID {
			tenant = &tenantUinList[idx]
		}
	}

	if tenant == nil {
		blog.Errorf("session not tenant info, tenant: %s, tenant list: %+v, rid: %s", tenantID, tenantUinList, rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommNotFound,
			ErrMsg: "not found",
		})
		return
	}
	session.Set(common.WEBSessionTenantUinKey, tenant.TenantID)
	if err := session.Save(); err != nil {
		blog.Errorf("save session failed, err: %v, rid: %s", err, rid)
	}

	ret := metadata.LoginChangeTenantResult{}
	ret.Result = true
	ret.Data.TenantID = tenantID

	c.JSON(200, ret)
	return
}
