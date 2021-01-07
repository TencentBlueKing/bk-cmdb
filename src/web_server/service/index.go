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

package service

import (
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/version"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
)

// Index html file
func (s *Service) Index(c *gin.Context) {
	session := sessions.Default(c)
	role := session.Get("role")
	userName, _ := session.Get(common.WEBSessionUinKey).(string)

	pageConf := gin.H{
		"site":           s.Config.Site.DomainUrl,
		"version":        s.Config.Version,
		"ccversion":      version.CCVersion,
		"authscheme":     s.Config.Site.AuthScheme,
		"fullTextSearch": s.Config.Site.FullTextSearch,
		"role":           role,
		"curl":           s.Config.LoginUrl,
		"userName":       userName,
		"agentAppUrl":    s.Config.AgentAppUrl,
		"authCenter":     s.Config.AuthCenter,
		"helpDocUrl":     s.Config.Site.HelpDocUrl,
	}

	if s.Config.Site.PaasDomainUrl != "" {
		pageConf["userManage"] = strings.TrimSuffix(s.Config.Site.PaasDomainUrl, "/") + "/api/c/compapi/v2/usermanage/fs_list_users/"
	}

	c.HTML(200, "index.html", pageConf)
}
