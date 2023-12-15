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

package common

import (
	"configcenter/src/common"
	httpheader "configcenter/src/common/http/header"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SetProxyHeader TODO
func SetProxyHeader(c *gin.Context) {
	// http request header add user
	session := sessions.Default(c)
	userName, _ := session.Get(common.WEBSessionUinKey).(string)
	ownerID, _ := session.Get(common.WEBSessionOwnerUinKey).(string)

	// 删除 Accept-Encoding 避免返回值被压缩
	c.Request.Header.Del("Accept-Encoding")
	httpheader.AddUser(c.Request.Header, userName)
	httpheader.AddLanguage(c.Request.Header, GetLanguageByHTTPRequest(c))
	httpheader.AddSupplierAccount(c.Request.Header, ownerID)
}

// GetLanguageByHTTPRequest get language by http request cookie
func GetLanguageByHTTPRequest(c *gin.Context) string {
	cookieLanguage, err := c.Cookie(common.HTTPCookieLanguage)
	if err == nil && cookieLanguage != "" {
		return cookieLanguage
	}

	return string(common.Chinese)
}
