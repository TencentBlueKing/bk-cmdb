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

package service

import (
	"net/http"

	"configcenter/src/common/http/rest"

	"github.com/emicklei/go-restful/v3"
)

func (s *Service) initKube(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	// namespace
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/kube/createmany/namespace/bk_biz_id/{bk_biz_id}",
		Handler: s.CreateNamespace})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/kube/updatemany/namespace/bk_biz_id/{bk_biz_id}",
		Handler: s.UpdateNamespace})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/kube/deletemany/namespace/bk_biz_id/{bk_biz_id}",
		Handler: s.DeleteNamespace})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/kube/findmany/namespace/bk_biz_id/{bk_biz_id}",
		Handler: s.ListNamespace})

	// workload
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/kube/createmany/workload/{kind}/{bk_biz_id}",
		Handler: s.CreateWorkload})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/kube/updatemany/workload/{kind}/{bk_biz_id}",
		Handler: s.UpdateWorkload})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/kube/deletemany/workload/{kind}/{bk_biz_id}",
		Handler: s.DeleteWorkload})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/kube/findmany/workload/{kind}/{bk_biz_id}",
		Handler: s.ListWorkload})

	// topo
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/kube/find/host_node_path",
		Handler: s.FindNodePathForHost})

	// pod
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/kube/find/pod_path/bk_biz_id/{bk_biz_id}",
		Handler: s.FindPodPath})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/kube/findmany/pod/bk_biz_id/{bk_biz_id}",
		Handler: s.ListPod})

	// container
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/kube/findmany/container/bk_biz_id/{bk_biz_id}",
		Handler: s.ListContainer})

	utility.AddToRestfulWebService(web)
}
