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

package logics

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"configcenter/src/common"
	"configcenter/src/common/blog"
)

func SetProxyHeader(c *gin.Context) {
	//http request header add user
	session := sessions.Default(c)
	userName, _ := session.Get("userName").(string)
	language, _ := session.Get("language").(string)
	ownerID, _ := session.Get("owner_uin").(string)
	c.Request.Header.Add(common.BKHTTPHeaderUser, userName)
	c.Request.Header.Add(common.BKHTTPLanguage, language)
	c.Request.Header.Add(common.BKHTTPOwnerID, ownerID)
}

func GetLanguageByHTTPRequest(c *gin.Context) string {

	cookieLanuage, err := c.Cookie(common.BKHTTPCookieLanugageKey)
	if "" != cookieLanuage && nil == err {
		return cookieLanuage
	}

	session := sessions.Default(c)
	language := session.Get(common.BKSessionLanugageKey)
	if nil == language {
		return ""
	}
	strLang, ok := language.(string)

	if false == ok {
		blog.Errorf("get language from session error, %v", language)
		return strLang
	}

	return strLang
}
