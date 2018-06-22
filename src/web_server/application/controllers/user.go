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

package controllers

import (
	"fmt"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/wactions"
	"configcenter/src/web_server/application/middleware/user"
)

const BkAccountUrl = "site.bk_account_url"

func init() {
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectGet, "/user/list", nil, GetUserList})
	wactions.RegisterNewAction(wactions.Action{common.HTTPUpdate, "/user/language/:language", nil, UpdateUserLanguage})

}

type userResult struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Code    string      `json:"code"`
	Result  bool        `json:"result"`
}

var getUserFailData = map[string]interface{}{
	"result":        false,
	"bk_error_msg":  "get user list false",
	"bk_error_code": "",
	"data":          nil,
}

// GetUserList get user list
func GetUserList(c *gin.Context) {
	a := api.NewAPIResource()
	config, _ := a.ParseConfig()
	accountURL := config[BkAccountUrl]
	user := user.NewUser()
	code, data := user.GetUserList(c, accountURL)
	c.JSON(code, data)
	return
}

func UpdateUserLanguage(c *gin.Context) {
	session := sessions.Default(c)
	language := c.Param("language")

	session.Set("language", language)
	err := session.Save()

	if nil != err {
		blog.Errorf("user update language error:%s", err.Error())
		c.JSON(200, gin.H{
			"result":        false,
			"bk_error_msg":  "user update language error",
			"bk_error_code": fmt.Sprintf("%d", common.CCErrCommHTTPDoRequestFailed),
			"data":          nil,
		})
		return
	}

	c.SetCookie("blueking_language", language, 0, "/", "", false, true)

	c.JSON(200, gin.H{
		"result":        true,
		"bk_error_msg":  "",
		"bk_error_code": "00",
		"data":          nil,
	})
	return
}
