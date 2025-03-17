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

	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
)

// SetCookie set tenant cookie
func SetCookie(c *gin.Context, session sessions.Session) error {
	rid := httpheader.GetRid(c.Request.Header)
	enableTenantMode, err := cc.Bool("tenant.enableMultiTenantMode")
	if err != nil {
		blog.Errorf("get enable tenant mode failed, err: %v, rid: %s", err, rid)
		return err
	}

	cookieTenantID, err := c.Cookie(common.HTTPCookieTenant)
	if cookieTenantID == "" || err != nil {
		if enableTenantMode {
			blog.Errorf("tenant mode is enabled but tenant cookie is not set, rid: %s", rid)
			return fmt.Errorf("tenant mode is enabled but tenant cookie is not set")
		} else {
			c.SetCookie(common.HTTPCookieTenant, common.BKUnconfiguredTenantID, 0, "/", "", false, false)
			session.Set(common.WEBSessionTenantUinKey, common.BKUnconfiguredTenantID)
		}
	} else {
		if enableTenantMode {
			session.Set(common.WEBSessionTenantUinKey, cookieTenantID)
		} else {
			session.Set(common.WEBSessionTenantUinKey, common.BKUnconfiguredTenantID)
		}
	}

	if err = session.Save(); err != nil {
		blog.Warnf("save session failed, err: %s, rid: %s", err.Error(), rid)
	}
	return nil
}
