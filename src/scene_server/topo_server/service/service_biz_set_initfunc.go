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
	"net/http"

	"configcenter/src/common/http/rest"

	"github.com/emicklei/go-restful"
)

// business set operation
func (s *Service) initBusinessSet(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/create/biz_set",
		Handler: s.CreateBusinessSet})

	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/findmany/biz_set/with_reduced",
		Handler: s.SearchReducedBusinessSetList})

	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/findmany/biz_set",
		Handler: s.SearchBusinessSet})

	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/preview/biz_set",
		Handler: s.PreviewBusinessSet})

	utility.AddToRestfulWebService(web)
}
