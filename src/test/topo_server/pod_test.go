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

package topo_server_test

import (
	"context"
	"encoding/json"

	"configcenter/pkg/filter"
	"configcenter/src/common"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	commonutil "configcenter/src/common/util"
	"configcenter/src/kube/types"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pod test", func() {
	ctx := context.Background()

	var bizID int64
	bizName := "biz_for_kube"
	var clusterID int64
	clusterUID := "clusterUID"
	clusterName := "clusterName"
	var namespaceID int64
	nsName := "nsName"
	var wlID int64
	wlName := "wlName"
	wlKind := "deployment"
	var nodeID int64
	var hostID1 int64
	var hostID2 int64
	podName := "podName"
	var podID int64
	containerName := "containerName"
	containerUID := "containerUID"
	It("prepare environment, create business, cluster, namespace", func() {
		test.ClearDatabase()

		// create business
		biz := map[string]interface{}{
			common.BKMaintainersField: "kube",
			common.BKAppNameField:     bizName,
			"time_zone":               "Africa/Accra",
			"language":                "1",
		}
		bizResp, err := apiServerClient.CreateBiz(ctx, "0", header, biz)
		util.RegisterResponseWithRid(bizResp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(bizResp.Result).To(Equal(true))
		bizID, err = commonutil.GetInt64ByInterface(bizResp.Data[common.BKAppIDField])
		Expect(err).NotTo(HaveOccurred())

		// create host
		input := map[string]interface{}{
			common.BKAppIDField: bizID,
			"host_info": map[string]interface{}{
				"0": map[string]interface{}{
					"bk_host_innerip": "127.0.0.1",
					"bk_cloud_id":     0,
				},
				"1": map[string]interface{}{
					"bk_host_innerip": "127.0.0.2",
					"bk_cloud_id":     0,
				},
			},
		}
		rsp, err := hostServerClient.AddHost(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true), rsp.ToString())
		searchOpt := &params.HostCommonSearch{
			AppID: int(bizID),
			Ip: params.IPInfo{
				Data:  []string{"127.0.0.1", "127.0.0.2"},
				Exact: 1,
				Flag:  "bk_host_innerip|bk_host_outerip",
			},
		}
		hostRep, err := hostServerClient.SearchHost(context.Background(), header, searchOpt)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRep.Result).To(Equal(true))
		Expect(hostRep.Data.Count).To(Equal(2))
		data1 := hostRep.Data.Info[0]["host"].(map[string]interface{})
		hostID1, err = commonutil.GetInt64ByInterface(data1[common.BKHostIDField])
		data2 := hostRep.Data.Info[1]["host"].(map[string]interface{})
		hostID2, err = commonutil.GetInt64ByInterface(data2[common.BKHostIDField])

		schedulingEngine := "k8s"
		xid := "cls-hox2lkf2"
		version := "0.1"
		networkType := "underlay"
		region := "shenzhen"
		vpc := "vpc-q6awe02n"
		network := []string{"1.1.1.0/21"}
		clusterType := "public"
		createCLuster := &types.Cluster{
			Name:             &clusterName,
			SchedulingEngine: &schedulingEngine,
			Uid:              &clusterUID,
			Xid:              &xid,
			Version:          &version,
			NetworkType:      &networkType,
			Region:           &region,
			Vpc:              &vpc,
			NetWork:          &network,
			Type:             &clusterType,
		}

		id, err := kubeClient.CreateCluster(ctx, header, bizID, createCLuster)

		util.RegisterResponse(id)
		Expect(err).NotTo(HaveOccurred())
		clusterID = id

		// create namespace
		ns := types.Namespace{
			ClusterSpec: types.ClusterSpec{
				ClusterID: clusterID,
			},
			Name: nsName,
		}
		createNsOpt := types.NsCreateOption{
			Data: []types.Namespace{ns},
		}

		nsResult, err := kubeClient.CreateNamespace(ctx, header, bizID, &createNsOpt)
		util.RegisterResponseWithRid(nsResult, header)
		Expect(err).NotTo(HaveOccurred())
		namespaceID = nsResult.IDs[0]

		// create workload
		wl := types.Deployment{
			WorkloadBase: types.WorkloadBase{
				NamespaceSpec: types.NamespaceSpec{
					ClusterSpec: types.ClusterSpec{
						BizID: bizID,
					},
					NamespaceID: namespaceID,
				},
				Name: wlName,
			},
		}
		createWOpt := types.WlCreateOption{
			Data: []types.WorkloadInterface{
				&wl,
			},
		}
		wlResult, err := kubeClient.CreateWorkload(ctx, header, bizID, types.KubeDeployment, &createWOpt)
		util.RegisterResponseWithRid(wlResult, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(wlResult.IDs)).To(Equal(1))
		wlID = wlResult.IDs[0]

		// create node
		nodeName1 := "node1"
		nodeName2 := "node2"
		unschedulable := false
		hostName1 := "hostname1"
		hostName2 := "hostname2"
		internalIP1 := []string{"127.0.0.1"}
		externalIP1 := []string{"127.0.0.1"}
		internalIP2 := []string{"127.0.0.2"}
		externalIP2 := []string{"127.0.0.2"}
		createNode := &types.CreateNodesOption{
			Nodes: []types.OneNodeCreateOption{
				{
					HostID:    hostID1,
					ClusterID: clusterID,
					Node: types.Node{
						Name:          &nodeName1,
						Unschedulable: &unschedulable,
						InternalIP:    &internalIP1,
						ExternalIP:    &externalIP1,
						HostName:      &hostName1,
					},
				}, {
					HostID:    hostID2,
					ClusterID: clusterID,
					Node: types.Node{
						Name:          &nodeName2,
						Unschedulable: &unschedulable,
						InternalIP:    &internalIP2,
						ExternalIP:    &externalIP2,
						HostName:      &hostName2,
					},
				},
			},
		}
		nodeResult, err := kubeClient.BatchCreateNode(ctx, header, bizID, createNode)
		util.RegisterResponse(nodeResult)
		Expect(err).NotTo(HaveOccurred())
		nodeID = nodeResult[0]
	})

	It("create pod and container", func() {
		containerImage := "image"
		createOpt := types.CreatePodsOption{
			Data: []types.PodsInfoArray{
				{
					BizID: bizID,
					Pods: []types.PodsInfo{
						{
							Spec: types.SpecSimpleInfo{
								ClusterID:   clusterID,
								NamespaceID: namespaceID,
								NodeID:      nodeID,
								Ref: types.Reference{
									Kind: types.WorkloadType(wlKind),
									ID:   wlID,
								},
							},
							Pod: types.Pod{
								Name: &podName,
							},
							Containers: []types.Container{
								{
									Name:        &containerName,
									ContainerID: &containerUID,
									Image:       &containerImage,
								},
							},
							HostID: hostID1,
						},
					},
				},
			},
		}

		result, err := kubeClient.BatchCreatePod(ctx, header, &createOpt)
		util.RegisterResponseWithRid(result, header)
		Expect(err).NotTo(HaveOccurred())
		podID = result[0]
	})

	It("find pod", func() {
		filter := &filter.Expression{
			RuleFactory: &filter.CombinedRule{
				Condition: filter.And,
				Rules: []filter.RuleFactory{
					&filter.AtomRule{
						Field:    types.ClusterUIDField,
						Operator: filter.Equal.Factory(),
						Value:    clusterUID,
					},
					&filter.AtomRule{
						Field:    types.NamespaceField,
						Operator: filter.Equal.Factory(),
						Value:    nsName,
					},
					&filter.AtomRule{
						Field:    types.RefField,
						Operator: filter.Object.Factory(),
						Value: &filter.CombinedRule{
							Condition: filter.And,
							Rules: []filter.RuleFactory{
								&filter.AtomRule{
									Field:    types.KindField,
									Operator: filter.Equal.Factory(),
									Value:    wlKind,
								},
								&filter.AtomRule{
									Field:    types.KubeNameField,
									Operator: filter.Equal.Factory(),
									Value:    wlName,
								},
							},
						},
					},
					&filter.AtomRule{
						Field:    common.BKFieldName,
						Operator: filter.Equal.Factory(),
						Value:    podName,
					},
				},
			},
		}

		// get pod data
		page := metadata.BasePage{
			Start: 0,
			Limit: 10,
		}
		fields := []string{common.BKFieldID}
		queryOpt := types.PodQueryOption{
			Filter: filter,
			Page:   page,
			Fields: fields,
		}
		result, err := kubeClient.ListPod(ctx, header, bizID, &queryOpt)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result.Info)).To(Equal(1))
		Expect(result.Info[0][common.BKFieldID].(json.Number).Int64()).To(Equal(podID))

		// get pod count
		page = metadata.BasePage{
			EnableCount: true,
		}
		queryOpt = types.PodQueryOption{
			Filter: filter,
			Page:   page,
		}
		result, err = kubeClient.ListPod(ctx, header, bizID, &queryOpt)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Count).To(Equal(1))
	})

	It("find container", func() {
		filter := &filter.Expression{
			RuleFactory: &filter.CombinedRule{
				Condition: filter.And,
				Rules: []filter.RuleFactory{
					&filter.AtomRule{
						Field:    common.BKFieldName,
						Operator: filter.Equal.Factory(),
						Value:    containerName,
					},
					&filter.AtomRule{
						Field:    types.BKPodIDField,
						Operator: filter.Equal.Factory(),
						Value:    podID,
					},
				},
			},
		}

		// get container data
		page := metadata.BasePage{
			Start: 0,
			Limit: 10,
		}
		fields := []string{types.ContainerUIDField}
		queryOpt := types.ContainerQueryOption{
			Filter: filter,
			Page:   page,
			Fields: fields,
		}
		result, err := kubeClient.ListContainer(ctx, header, bizID, &queryOpt)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result.Info)).To(Equal(1))
		Expect(result.Info[0][types.ContainerUIDField].(string)).To(Equal(containerUID))

		// get container count
		page = metadata.BasePage{
			EnableCount: true,
		}
		queryOpt = types.ContainerQueryOption{
			Filter: filter,
			Page:   page,
		}
		result, err = kubeClient.ListContainer(ctx, header, bizID, &queryOpt)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Count).To(Equal(1))
	})

	It("find node path for host", func() {
		req := types.HostPathOption{
			HostIDs: []int64{
				hostID1,
			},
		}

		result, err := kubeClient.FindNodePathForHost(ctx, header, &req)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result.Info)).To(Equal(1))
		Expect(result.Info[0].HostID).To(Equal(hostID1))
		Expect(result.Info[0].Paths[0].BizID).To(Equal(bizID))
		Expect(result.Info[0].Paths[0].BizName).To(Equal(bizName))
		Expect(result.Info[0].Paths[0].ClusterID).To(Equal(clusterID))
		Expect(result.Info[0].Paths[0].ClusterName).To(Equal(clusterName))
	})

	It("find pod path", func() {

		req := types.PodPathOption{
			PodIDs: []int64{podID},
		}

		result, err := kubeClient.FindPodPath(ctx, header, bizID, &req)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result.Info)).To(Equal(1))
		Expect(result.Info[0].BizName).To(Equal(bizName))
		Expect(result.Info[0].ClusterID).To(Equal(clusterID))
		Expect(result.Info[0].ClusterName).To(Equal(clusterName))
		Expect(result.Info[0].NamespaceID).To(Equal(namespaceID))
		Expect(result.Info[0].Namespace).To(Equal(nsName))
		Expect(result.Info[0].Kind).To(Equal(types.WorkloadType(wlKind)))
		Expect(result.Info[0].WorkloadID).To(Equal(wlID))
		Expect(result.Info[0].WorkloadName).To(Equal(wlName))
		Expect(result.Info[0].PodID).To(Equal(podID))
	})

	It("find host with k8s condition", func() {
		req := types.SearchHostOption{
			BizID:     bizID,
			ClusterID: clusterID,
		}

		result, err := hostServerClient.SearchKubeHost(ctx, header, req)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result.Info)).To(Equal(2))

		req = types.SearchHostOption{
			BizID:       bizID,
			ClusterID:   clusterID,
			NamespaceID: namespaceID,
			WorkloadID:  wlID,
			WlKind:      types.WorkloadType(wlKind),
		}
		result, err = hostServerClient.SearchKubeHost(ctx, header, req)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result.Info)).To(Equal(1))

		req = types.SearchHostOption{
			BizID:     bizID,
			ClusterID: clusterID,
			Folder:    true,
		}
		result, err = hostServerClient.SearchKubeHost(ctx, header, req)

		Expect(err).NotTo(HaveOccurred())
		Expect(len(result.Info)).To(Equal(1))
	})
})
