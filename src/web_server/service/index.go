/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package service

import (
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/version"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Index html file
func (s *Service) Index(c *gin.Context) {
	session := sessions.Default(c)
	userName, _ := session.Get(common.WEBSessionUinKey).(string)
	tenantID, _ := session.Get(common.WEBSessionTenantUinKey).(string)
	timeZone, _ := session.Get(common.WEBSessionTimeZoneKey).(string)

	pageConf := gin.H{
		"site":                  s.Config.Site.DomainUrl,
		"version":               s.Config.Version,
		"ccversion":             version.CCVersion,
		"authscheme":            s.Config.Site.AuthScheme,
		"fullTextSearch":        s.Config.Site.FullTextSearch,
		"curl":                  s.Config.LoginUrl,
		"userName":              userName,
		"tenantId":              tenantID,
		"agentAppUrl":           s.Config.AgentAppUrl,
		"helpDocUrl":            s.Config.Site.HelpDocUrl,
		"cookieDomain":          s.Config.Site.BkDomain,
		"componentApiUrl":       s.Config.Site.BkComponentApiUrl,
		"publicPath":            getPublicPath(s.Config.Site.DomainUrl),
		"enableNotification":    s.Config.EnableNotification,
		"bkSharedResUrl":        s.Config.Site.BkSharedResUrl,
		"enableMultiTenantMode": s.Config.EnableMultiTenantMode,
		"timeZone":              timeZone,
	}

	if s.Config.Site.PaasDomainUrl != "" {
		pageConf["userManage"] = strings.TrimSuffix(s.Config.Site.PaasDomainUrl,
			"/") + "/api/c/compapi/v2/usermanage/fs_list_users/"
	}

	if s.Config.Site.BkApiUrlTmpl != "" {
		pageConf["userManageUrl"] = strings.TrimSuffix(strings.ReplaceAll(s.Config.Site.BkApiUrlTmpl, "{api_name}",
			"bk-user-web"), "/") + "/prod"
	}

	c.HTML(200, "index.html", pageConf)
}

// getPublicPath 获取前端需要的资源目录
// 如：http://127.0.0.1/test  -> /test/
func getPublicPath(site string) string {
	site = strings.TrimPrefix(site, "http://")
	site = strings.TrimPrefix(site, "https://")

	segments := strings.Split(site, "/")
	publicPath := strings.Join(segments[1:], "/")
	if publicPath == "" {
		return "/"
	}

	if !strings.HasPrefix(publicPath, "/") {
		publicPath = "/" + publicPath
	}

	if !strings.HasSuffix(publicPath, "/") {
		publicPath = publicPath + "/"
	}

	return publicPath
}
