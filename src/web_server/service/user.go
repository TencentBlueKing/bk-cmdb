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
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/web_server/middleware/user"
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
	user := user.NewUser(s.Config, s.Engine, s.CacheCli, s.VersionPlg)
	code, data := user.GetUserList(c)
	c.JSON(code, data)
	return
}

func (s *Service) UpdateUserLanguage(c *gin.Context) {
	session := sessions.Default(c)
	language := c.Param("language")

	session.Set("language", language)
	err := session.Save()

	if nil != err {
		blog.Errorf("user update language error:%s", err.Error())
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
			blog.Errorf("[UserInfo] json unmarshal error:%s, logID:%s", err.Error(), util.GetHTTPCCRequestID(c.Request.Header))
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
			blog.Errorf("[UserInfo] json unmarshal error:%s, logID:%s", err.Error(), util.GetHTTPCCRequestID(c.Request.Header))
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
	session.Save()
	ret := metadata.LoginChangeSupplierResult{}
	ret.Result = true
	ret.Data.ID = ownerID

	c.JSON(200, ret)
	return
}
