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

	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/create/kube/cluster/bk_biz_id/{bk_biz_id}",
		Handler: s.CreateCluster})

	utility.AddHandler(rest.Action{Verb: http.MethodDelete,
		Path:    "/delete/kube/cluster/bk_biz_id/{bk_biz_id}",
		Handler: s.DeleteCluster})

	utility.AddHandler(rest.Action{Verb: http.MethodPut,
		Path:    "/updatemany/kube/cluster/bk_biz_id/{bk_biz_id}",
		Handler: s.UpdateClusterFields})

	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/findmany/kube/cluster/bk_biz_id/{bk_biz_id}",
		Handler: s.SearchClusters})

	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/createmany/kube/node/bk_biz_id/{bk_biz_id}",
		Handler: s.BatchCreateNode})

	utility.AddHandler(rest.Action{Verb: http.MethodDelete,
		Path:    "/deletemany/kube/node/bk_biz_id/{bk_biz_id}",
		Handler: s.BatchDeleteNode})

	utility.AddHandler(rest.Action{Verb: http.MethodPut,
		Path:    "/updatemany/kube/node/bk_biz_id/{bk_biz_id}",
		Handler: s.UpdateNodeFields})

	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/findmany/kube/node/bk_biz_id/{bk_biz_id}",
		Handler: s.SearchNodes})

	// namespace
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/kube/namespace/bk_biz_id/{bk_biz_id}",
		Handler: s.CreateNamespace})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/kube/namespace/bk_biz_id/{bk_biz_id}",
		Handler: s.UpdateNamespace})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/kube/namespace/bk_biz_id/{bk_biz_id}",
		Handler: s.DeleteNamespace})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/namespace/bk_biz_id/{bk_biz_id}",
		Handler: s.ListNamespace})

	// workload
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/kube/workload/{kind}/{bk_biz_id}",
		Handler: s.CreateWorkload})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/kube/workload/{kind}/{bk_biz_id}",
		Handler: s.UpdateWorkload})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/kube/workload/{kind}/{bk_biz_id}",
		Handler: s.DeleteWorkload})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/workload/{kind}/{bk_biz_id}",
		Handler: s.ListWorkload})

	// topo
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/kube/host_node_path",
		Handler: s.FindNodePathForHost})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/find/kube/topo_path/bk_biz_id/{bk_biz_id}",
		Handler: s.SearchKubeTopoPath})

	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/find/kube/{bk_biz_id}/topo_node/{type}/count",
		Handler: s.CountKubeTopoHostsOrPods})

	utility.AddHandler(rest.Action{Verb: http.MethodGet,
		Path:    "/find/kube/{object}/attributes",
		Handler: s.FindResourceAttrs})

	// pod
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/createmany/kube/pod",
		Handler: s.BatchCreatePod})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/kube/pod_path/bk_biz_id/{bk_biz_id}",
		Handler: s.FindPodPath})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/pod/bk_biz_id/{bk_biz_id}",
		Handler: s.ListPod})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/kube/pod", Handler: s.DeletePods})

	// container
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/container/bk_biz_id/{bk_biz_id}",
		Handler: s.ListContainer})

	utility.AddToRestfulWebService(web)
}
