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
	"configcenter/src/common/backbone"

	"github.com/emicklei/go-restful"
)

type Service struct {
	*backbone.Engine
}

func (s *Service) WebService(filter restful.FilterFunction) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/api/v2").Filter(filter).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/host/batch").To(s.DeleteHostBatch))

	return ws
}
