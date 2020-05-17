/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
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

func (s *ContainerService) initPod(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/biz/{bk_biz_id}/pod", Handler: s.CreatePod})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/biz/{bk_biz_id}/pod", Handler: s.CreateManyPod})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/biz/{bk_biz_id}/pod", Handler: s.UpdatePod})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/biz/{bk_biz_id}/pod", Handler: s.DeletePod})
	//utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/container/biz/{bk_biz_id}/pod", Handler: s.DeleteManyPod})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/list/pod", Handler: s.ListPods})

	utility.AddToRestfulWebService(web)
}

func (s *ContainerService) initService(web *restful.WebService) {
	s.initPod(web)
}
