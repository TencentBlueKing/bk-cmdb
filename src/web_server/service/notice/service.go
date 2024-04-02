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

package notice

import (
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/thirdparty/apigw/notice"
	"configcenter/src/web_server/app/options"
	"configcenter/src/web_server/capability"
	webCommon "configcenter/src/web_server/common"

	"github.com/gin-gonic/gin"
)

type service struct {
	ws        *gin.Engine
	engine    *backbone.Engine
	config    *options.Config
	noticeCli notice.ClientI
}

// Init notice service
func Init(c *capability.Capability) {
	s := &service{
		engine:    c.Engine,
		config:    c.Config,
		noticeCli: c.NoticeCli,
	}

	c.Ws.GET("/notice/get_current_announcements", s.GetCurAnn)
}

// GetCurAnn get current announcements
func (s *service) GetCurAnn(c *gin.Context) {
	language := webCommon.GetLanguageByHTTPRequest(c)
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(language)

	if !s.config.EnableNotification {
		c.JSON(http.StatusBadRequest, metadata.BaseResp{Code: common.CCErrWebDisableNotification,
			ErrMsg: defErr.Error(common.CCErrWebDisableNotification).Error()})
		return
	}

	params := make(map[string]string)
	for key, val := range c.Request.URL.Query() {
		params[key] = val[0]
	}

	rid := httpheader.GetRid(c.Request.Header)
	ctx := util.NewContextFromGinContext(c)

	ann, err := s.noticeCli.GetCurAnn(ctx, c.Request.Header, params)
	if err != nil {
		blog.Errorf("get current announcements failed, req: %+v, err: %v, rid: %s", params, err, rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{Code: common.CCErrWebGetAnnFail,
			ErrMsg: defErr.Errorf(common.CCErrWebGetAnnFail, err.Error()).Error()})
		return
	}

	c.JSON(http.StatusOK, metadata.NewSuccessResp(ann))
}
