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

package metadata

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginUserInfo struct {
	UserName    string
	ChName      string
	Phone       string
	Email       string
	Role        int64
	BkToken     string
	OnwerUin    string
	OwnerUinArr []string               //user all owner uin
	IsOwner     bool                   // is master
	Extra       map[string]interface{} //custom information
	Language    string
}

type LoginPluginInfo struct {
	Name    string // plugin info
	Version string // In what version is used
	//CookieEnv string // Reserved word, not used now,  When the cookie has the current key, it is used preferentially.
	HandleFunc LoginUserPluginInerface
}

type LoginUserPluginParams struct {
	Url          string
	IsMultiOwner bool
	Cookie       []*http.Cookie // Reserved word, not used now
	Header       http.Header    // Reserved word, not used now
}

type LoginUserPluginInerface interface {
	LoginUser(c *gin.Context, config map[string]string, isMultiOwner bool) (user *LoginUserInfo, loginSucc bool)
	GetUserList(c *gin.Context, config map[string]string) ([]*LoginSystemUserInfo, error)
}

type LoginSystemUserInfo struct {
	CnName string `json:"chinese_name"`
	EnName string `json:"english_name"`
}

type LonginSystemUserListResult struct {
	BaseResp `json",inline"`
	Data     []*LoginSystemUserInfo `json:"data"`
}
