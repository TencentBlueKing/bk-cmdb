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
	"configcenter/src/common/metadata"
	"configcenter/src/web_server/middleware/user"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
)

// LogOutUser log out user
func (s *Service) LogOutUser(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	c.Request.URL.Path = ""
	userManger := user.NewUser(*s.Config, s.Engine, s.CacheCli, s.VersionPlg)
	loginURL := userManger.GetLoginUrl(c)
	ret := metadata.LogoutResult{}
	ret.BaseResp.Result = true
	ret.Data.LogoutURL = loginURL
	c.JSON(200, ret)
	return
}
