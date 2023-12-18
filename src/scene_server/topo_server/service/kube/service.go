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

package kube

import (
	"net/http"

	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/topo_server/service/capability"
)

type service struct {
	*capability.Capability
}

// InitKube init kube service
func InitKube(utility *rest.RestUtility, c *capability.Capability) {
	s := &service{
		Capability: c,
	}

	// cluster
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/kube/cluster", Handler: s.CreateCluster})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/kube/cluster", Handler: s.DeleteCluster})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/kube/cluster",
		Handler: s.UpdateClusterFields})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/kube/cluster/type",
		Handler: s.UpdateClusterType})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/cluster", Handler: s.SearchClusters})

	// node
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/kube/node", Handler: s.BatchCreateNode})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/kube/node", Handler: s.BatchDeleteNode})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/kube/node", Handler: s.UpdateNodeFields})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/node", Handler: s.SearchNodes})

	// namespace
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/kube/namespace",
		Handler: s.CreateNamespace})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/kube/namespace",
		Handler: s.UpdateNamespace})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/kube/namespace",
		Handler: s.DeleteNamespace})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/namespace", Handler: s.ListNamespace})

	// workload
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/kube/workload/{kind}",
		Handler: s.CreateWorkload})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/kube/workload/{kind}",
		Handler: s.UpdateWorkload})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/kube/workload/{kind}",
		Handler: s.DeleteWorkload})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/workload/{kind}",
		Handler: s.ListWorkload})

	// topo
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/kube/host_node_path",
		Handler: s.FindNodePathForHost})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/kube/topo_path", Handler: s.SearchKubeTopoPath})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/kube/topo_node/{type}/count",
		Handler: s.CountKubeTopoHostsOrPods})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/kube/{object}/attributes",
		Handler: s.FindResourceAttrs})

	// pod
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/kube/pod", Handler: s.BatchCreatePod})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/kube/pod_path", Handler: s.FindPodPath})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/pod", Handler: s.ListPod})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/kube/pod", Handler: s.DeletePods})

	// container
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/container", Handler: s.ListContainer})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/container/by_topo",
		Handler: s.ListContainerByTopo})
}
