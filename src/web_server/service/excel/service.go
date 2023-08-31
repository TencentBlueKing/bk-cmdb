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

package excel

import (
	"configcenter/src/common/backbone"
	"configcenter/src/web_server/capability"

	"github.com/gin-gonic/gin"
)

type service struct {
	ws     *gin.Engine
	engine *backbone.Engine
}

// Init init excel service
func Init(c *capability.Capability) {
	s := &service{
		engine: c.Engine,
	}

	c.Ws.POST("/importtemplate/:bk_obj_id", s.BuildTemplate)

	c.Ws.POST("/insts/object/:bk_obj_id/export", s.ExportInst)

	c.Ws.POST("/hosts/export", s.ExportHost)

	c.Ws.POST("/insts/object/:bk_obj_id/import", s.AddInst)

	c.Ws.POST("/hosts/import", s.AddHost)

	c.Ws.POST("/hosts/update", s.UpdateHost)

	c.Ws.POST("/object/object/:bk_obj_id/export", s.ExportObject)

	c.Ws.POST("/object/object/:bk_obj_id/import", s.ImportObject)
}
