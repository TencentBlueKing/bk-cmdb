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
	"fmt"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"configcenter/src/common"
	"configcenter/src/common/blog"
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
	user := user.NewUser(s.Config)
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
