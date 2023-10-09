/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package login

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/esbserver/esbutil"
)

// ClientInterface login client interface
type ClientInterface interface {
	GetUser(ctx context.Context, h http.Header) (resp *UserResponse, err error)
}

// NewClientInterface new login client interface
func NewClientInterface(client rest.ClientInterface, config *esbutil.EsbConfigSrv) ClientInterface {
	return &login{
		client: client,
		config: config,
	}
}

type login struct {
	config *esbutil.EsbConfigSrv
	client rest.ClientInterface
}

// UserResponse user response
type UserResponse struct {
	metadata.EsbBaseResponse `json:",inline"`
	Data                     UserInfo `json:"data"`
}

// UserInfo user info
type UserInfo struct {
	BkUsername string `json:"bk_username"`
	QQ         string `json:"qq"`
	Language   string `json:"language"`
	Phone      string `json:"phone"`
	WxUserid   string `json:"wx_userid"`
	Email      string `json:"email"`
	CHName     string `json:"chname"`
	TimeZone   string `json:"time_zone"`
}
