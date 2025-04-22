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

package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/thirdparty/apigw/apigwutil"
	"configcenter/src/thirdparty/apigw/user/types"
)

// VirtualUserClientI is the user api gateway client for virtual user
type VirtualUserClientI interface {
	BatchSearchVirtualUser(ctx context.Context, h http.Header, loginNames []string) ([]types.VirtualUserItem, error)
}

var (
	virtualUserAuth = make(map[string]string)
	lock            sync.RWMutex
)

// getAuthConfigByTenant get virtual user by tenantID
func getAuthConfigByTenant(tenantID string) (string, bool) {
	lock.RLock()
	defer lock.RUnlock()
	authConfig, exist := virtualUserAuth[tenantID]
	if !exist {
		return "", false
	}
	return authConfig, true
}

// setAuthConfig set virtual user by tenantID
func setAuthConfig(tenantID string, authConfig string) {
	lock.Lock()
	defer lock.Unlock()
	virtualUserAuth[tenantID] = authConfig
	return
}

// SetBKAuthHeader set api gateway authorization header
func SetBKAuthHeader(ctx context.Context, conf *apigwutil.ApiGWConfig, header http.Header,
	userCli VirtualUserClientI) (http.Header, error) {

	tenantID := httpheader.GetTenantID(header)
	if tenantID == "" {
		fmt.Errorf("tenant id is empty")
		return nil, fmt.Errorf("tenant id is empty")
	}

	if authInfo, exist := getAuthConfigByTenant(tenantID); exist {
		header = httpheader.SetBkAuth(header, authInfo)
		return header, nil
	}

	resp, err := userCli.BatchSearchVirtualUser(ctx, header, []string{"bk_admin"})
	if err != nil {
		blog.Errorf("search virtual user failed, err: %v", err)
		return nil, err
	}

	if len(resp) != 1 {
		blog.Errorf("search virtual user failed, resp: %v", resp)
		return nil, fmt.Errorf("search virtual user failed, resp: %v", resp)
	}

	authConf := apigwutil.AuthConfig{
		AppAuthConfig: apigwutil.AppAuthConfig{
			AppCode:   conf.AppCode,
			AppSecret: conf.AppSecret,
		},
		UserName: resp[0].VirtualUserName,
	}

	authInfo, err := json.Marshal(authConf)
	if err != nil {
		blog.Errorf("marshal default api auth config %+v failed, err: %v", authConf, err)
		return nil, err
	}

	header = httpheader.SetBkAuth(header, string(authInfo))
	setAuthConfig(tenantID, string(authInfo))
	return header, nil
}
