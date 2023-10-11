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
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/source_controller/coreservice/service/capability"
)

type service struct {
	core core.Core
}

// InitKube init kube related service
func InitKube(c *capability.Capability) {
	s := &service{
		core: c.Core,
	}

	// cluster
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/kube/cluster", Handler: s.CreateCluster})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/kube/cluster",
		Handler: s.BatchUpdateCluster})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/kube/cluster",
		Handler: s.BatchDeleteCluster})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/cluster", Handler: s.SearchClusters})

	// shared cluster
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/shared/cluster/ns/relation",
		Handler: s.ListNsSharedClusterRel})

	// node
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/kube/node", Handler: s.BatchUpdateNode})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/kube/node",
		Handler: s.BatchDeleteNode})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/kube/node", Handler: s.BatchCreateNode})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/node", Handler: s.SearchNodes})

	// namespace
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/namespace", Handler: s.CreateNamespace})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/namespace", Handler: s.UpdateNamespace})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/namespace",
		Handler: s.DeleteNamespace})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/namespace", Handler: s.ListNamespace})

	// workload
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/workload/{kind}",
		Handler: s.CreateWorkload})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/workload/{kind}",
		Handler: s.UpdateWorkload})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/workload/{kind}",
		Handler: s.DeleteWorkload})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/workload/{kind}", Handler: s.ListWorkload})

	// pod
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/kube/pod", Handler: s.BatchCreatePod})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/pod", Handler: s.ListPod})
	c.Utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/pod", Handler: s.DeletePods})

	// container
	c.Utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/container", Handler: s.ListContainer})
}
