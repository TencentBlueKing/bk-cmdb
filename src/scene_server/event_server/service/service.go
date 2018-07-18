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
	"github.com/emicklei/go-restful"
	redis "gopkg.in/redis.v5"

	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/rdapi"
	"configcenter/src/storage"
)

type Service struct {
	*backbone.Engine
	db    storage.DI
	cache *redis.Client
}

func (s *Service) SetDB(db storage.DI) {
	s.db = db
}

func (s *Service) SetCache(db *redis.Client) {
	s.cache = db
}

func (s *Service) WebService() *restful.WebService {
	ws := new(restful.WebService)
	getErrFun := func() errors.CCErrorIf {
		return s.CCErr
	}
	ws.Path("/event/v3").Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)

	ws.Route(ws.POST("/subscribe/search/{ownerID}/{appID}").To(s.Query))
	ws.Route(ws.POST("/subscribe/ping").To(s.Ping))
	ws.Route(ws.POST("/subscribe/telnet").To(s.Telnet))
	ws.Route(ws.POST("/subscribe/{ownerID}/{appID}").To(s.Subscribe))
	ws.Route(ws.DELETE("/subscribe/{ownerID}/{appID}/{subscribeID}").To(s.UnSubscribe))
	ws.Route(ws.PUT("/subscribe/{ownerID}/{appID}/{subscribeID}").To(s.Rebook))

	return ws
}
