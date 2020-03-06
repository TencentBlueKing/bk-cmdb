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
	"configcenter/src/auth/extensions"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/rdapi"
	"configcenter/src/scene_server/topo_server/app/options"
	"configcenter/src/scene_server/topo_server/core"
	"configcenter/src/storage/dal"
	"configcenter/src/thirdpartyclient/elasticsearch"
	"encoding/json"

	"github.com/emicklei/go-restful"
)

type Service struct {
	Engine      *backbone.Engine
	DB          dal.RDB
	Core        core.Core
	Config      options.Config
	AuthManager *extensions.AuthManager
	Es          *elasticsearch.EsSrv
	Error       errors.CCErrorIf
	Language    language.CCLanguageIf
	EnableTxn   bool
}

// WebService the web service
func (s *Service) WebService() *restful.Container {
	errors.SetGlobalCCError(s.Error)
	getErrFunc := func() errors.CCErrorIf {
		return s.Error
	}

	api := new(restful.WebService)
	api.Path("/topo/v3/").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)

	// init service actions
	s.initService(api)

	healthz := new(restful.WebService).Produces(restful.MIME_JSON)
	healthz.Route(healthz.GET("/healthz").To(s.Healthz))
	container := restful.NewContainer().Add(api)
	container.Add(healthz)

	return container
}

type MetaShell struct {
	Metadata *metadata.Metadata `json:"metadata"`
}

type MapStrWithMetadata struct {
	Metadata *metadata.Metadata
	Data     mapstr.MapStr
}

func (m *MapStrWithMetadata) UnmarshalJSON(data []byte) error {
	md := new(MetaShell)
	if err := json.Unmarshal(data, md); err != nil {
		return err
	}
	m.Metadata = md.Metadata
	if err := json.Unmarshal(data, &m.Data); err != nil {
		return err
	}

	return nil
}
