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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	webcom "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/user"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
)

type userDataResult struct {
	Message string      `json:"bk_error_msg"`
	Data    interface{} `json:"data"`
	Code    string      `json:"bk_error_code"`
	Result  bool        `json:"result"`
}

// GetUserList get user list
func (s *Service) GetUserList(c *gin.Context) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	rspBody := metadata.LonginSystemUserListResult{}

	userManger := user.NewUser(*s.Config, s.Engine, s.CacheCli)
	userList, rawErr := userManger.GetUserList(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(c.Request.Header))
	if rawErr != nil && rawErr.ErrCode != 0 {
		blog.Error("GetUserList failed, err: %s, rid: %s", rawErr.ToCCError(defErr).Error(), rid)
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
		blog.ErrorJSON("session not owner info, rid:%s", rid)
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
		blog.ErrorJSON("session not owner info. owner:%s, ownerlist:%s, rid:%s", ownerID, ownerUinList, rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommNotFound,
			ErrMsg: "not found",
		})
		return
	}
	session.Set(common.WEBSessionOwnerUinKey, supplier.OwnerID)
	session.Set(common.WEBSessionRoleKey, strconv.FormatInt(supplier.Role, 10))
	if err := session.Save(); err != nil {
		blog.Errorf("save session failed, err: %+v, rid: %s", err, rid)
	}

	// not user, notice not privilege
	uin, _ := session.Get(common.WEBSessionUinKey).(string)
	language := webcom.GetLanguageByHTTPRequest(c)

	ownerM := user.NewOwnerManager(uin, supplier.OwnerID, language)
	ownerM.CacheCli = s.CacheCli
	ownerM.Engine = s.Engine
	permissions, err := ownerM.InitOwner()
	if nil != err {
		blog.Errorf("InitOwner error: %v, rid:%s", err, rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result:      false,
			Code:        err.GetCode(),
			ErrMsg:      err.Error(),
			Permissions: permissions,
		})
		return
	}

	ret := metadata.LoginChangeSupplierResult{}
	ret.Result = true
	ret.Data.ID = ownerID

	c.JSON(200, ret)
	return
}
