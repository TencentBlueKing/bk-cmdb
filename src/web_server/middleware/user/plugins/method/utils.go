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

package method

import (
	"fmt"

	"configcenter/pkg/tenant/logics"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/web_server/app/options"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SetTenantFromCookie get tenant from cookie, set tenant to session
func SetTenantFromCookie(c *gin.Context, config options.Config, session sessions.Session) (string, error) {
	rid := httpheader.GetRid(c.Request.Header)

	cookieTenantID, cookieErr := c.Cookie(common.HTTPCookieTenant)
	tenantID, err := logics.GetTenantWithMode(cookieTenantID, config.EnableMultiTenantMode)
	if err != nil {
		return "", err
	}

	if cookieTenantID == "" || cookieErr != nil {
		if config.EnableMultiTenantMode {
			blog.Errorf("tenant mode is enabled but tenant cookie is not set, rid: %s", rid)
			return tenantID, fmt.Errorf("tenant mode is enabled but tenant cookie is not set")
		}
		c.SetCookie(common.HTTPCookieTenant, tenantID, 0, "/", "", false, false)
		session.Set(common.WEBSessionTenantUinKey, tenantID)
	} else {
		session.Set(common.WEBSessionTenantUinKey, tenantID)
	}

	if err = session.Save(); err != nil {
		blog.Warnf("save session failed, err: %v, rid: %s", err, rid)
	}
	return tenantID, nil
}
