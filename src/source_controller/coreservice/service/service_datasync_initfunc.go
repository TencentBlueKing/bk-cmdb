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
	"net/http"

	"configcenter/src/common/http/rest"

	"github.com/emicklei/go-restful/v3"
)

func (s *coreService) initDataSynchronize(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/set/synchronize/instance", Handler: s.SynchronizeInstance})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/set/synchronize/model", Handler: s.SynchronizeModel})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/set/synchronize/association", Handler: s.SynchronizeAssociation})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/read/synchronize", Handler: s.SynchronizeFind})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/clear/synchronize/data", Handler: s.SynchronizeClearData})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/set/synchronize/identifier/flag", Handler: s.SetIdentifierFlag})

	utility.AddToRestfulWebService(web)
}
