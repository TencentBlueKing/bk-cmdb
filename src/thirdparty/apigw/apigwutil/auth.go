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

package apigwutil

import (
	"encoding/json"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
)

// AuthConfig defines the api gateway authorization config
type AuthConfig struct {
	AppAuthConfig `json:",inline"`
	BkToken       string `json:"bk_token,omitempty"`
	BkTicket      string `json:"bk_ticket,omitempty"`
	UserName      string `json:"bk_username,omitempty"`
	AccessToken   string `json:"access_token,omitempty"`
}

// AppAuthConfig defines the api gateway authorization config for blueking app
type AppAuthConfig struct {
	AppCode   string `json:"bk_app_code,omitempty"`
	AppSecret string `json:"bk_app_secret,omitempty"`
}

// SetAuthHeader set api gateway authorization header
func SetAuthHeader(appConf AppAuthConfig, header http.Header) http.Header {
	conf := AuthConfig{
		AppAuthConfig: appConf,
		BkToken:       util.GetBkToken(header),
		BkTicket:      util.GetBkTicket(header),
		UserName:      util.GetUser(header),
	}

	authInfo, err := json.Marshal(conf)
	if err != nil {
		blog.Errorf("marshal api auth config %+v failed, err: %v, rid: %s", conf, err, util.GetHTTPCCRequestID(header))
		return header
	}

	header.Set(common.BkHTTPHeaderAuth, string(authInfo))
	return header
}
