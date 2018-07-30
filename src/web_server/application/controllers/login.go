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
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/wactions"
)

//LogOutUser log out user
func LogOutUser(c *gin.Context) {
	a := api.NewAPIResource()
	config, _ := a.ParseConfig()
	site := config["site.domain_url"]
	loginURL := config["site.bk_login_url"]
	appCode := config["site.app_code"]
	loginPage := fmt.Sprintf(loginURL, appCode, site)
	session := sessions.Default(c)
	session.Clear()
	c.Redirect(301, loginPage)
}

func init() {
	wactions.RegisterNewAction(wactions.Action{common.HTTPSelectGet, "/logout", nil, LogOutUser})
}
