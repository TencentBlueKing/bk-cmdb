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

package user

import (
	"plugin"

	"configcenter/src/common/backbone"
	"configcenter/src/web_server/app/options"

	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v5"
)

type User interface {
	LoginUser(c *gin.Context) (isLogin bool)
	GetUserList(c *gin.Context) (int, interface{})
	GetLoginUrl(c *gin.Context) string
}

// NewUser return user instance by type
func NewUser(config options.Config, engine *backbone.Engine, cacheCli *redis.Client, loginPlg *plugin.Plugin) User {
	return &publicUser{config, engine, cacheCli, loginPlg}
}
