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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	commonutil "configcenter/src/common/util"
	"configcenter/src/kube/types"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kube cluster test", func() {
	ctx := context.Background()
	var bizId, clusterID, hostId1, nodeID, nodeID2 int64
	Describe("test preparation", func() {
		It("create business bk_biz_name = 'cc_biz'", func() {
			test.ClearDatabase()

			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_name":       "cc_biz",
				"time_zone":         "Africa/Accra",
			}
			rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data).To(ContainElement("cc_biz"))
			bizId, err = commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("add host using api", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"host_info": map[string]interface{}{
					"4": map[string]interface{}{
						"bk_host_innerip": "127.0.0.1",
						"bk_asset_id":     "addhost_api_asset_1",
						"bk_cloud_id":     0,
						"bk_comment":      "127.0.0.1 comment",
					},
				},
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search host created using api", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Ip: params.IPInfo{
					Data:  []string{"127.0.0.1"},
					Exact: 1,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.1"))
			Expect(data["bk_asset_id"].(string)).To(Equal("addhost_api_asset_1"))
			hostId1, err = commonutil.GetInt64ByInterface(data["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("create kube cluster", func() {
			clusterName := "cluster"
			schedulingEngine := "k8s"
			uid := "BCS-xxx-xxx004"
			xid := "xid-008"
			version := "0.1"
			networkType := "underlay"
			region := "shenzhen"
			vpc := "vpc-002"
			network := []string{"1.1.1.0/21"}
			clusterType := types.ClusterShareTypeField
			createCLuster := &types.Cluster{
				Name:             &clusterName,
				SchedulingEngine: &schedulingEngine,
				Uid:              &uid,
				Xid:              &xid,
				Version:          &version,
				NetworkType:      &networkType,
				Region:           &region,
				Vpc:              &vpc,
				NetWork:          &network,
				Type:             &clusterType,
			}

			id, err := kubeClient.CreateCluster(ctx, header, bizId, createCLuster)
			util.RegisterResponse(id)
			Expect(err).NotTo(HaveOccurred())
			clusterID = id
		})
	})

	It("create kube node", func() {
		By("create node")
		func() {
			node := "node"
			unschedulable := false
			hostName := "hostname"
			internalIP := []string{"1.1.1.1", "2.2.2.2"}
			externalIP := []string{"3.3.3.3", "4.4.4.4"}
			createNode := &types.CreateNodesOption{
				Nodes: []types.OneNodeCreateOption{
					{
						HostID:    hostId1,
						ClusterID: clusterID,
						Node: types.Node{
							Name:          &node,
							Unschedulable: &unschedulable,
							InternalIP:    &internalIP,
							ExternalIP:    &externalIP,
							HostName:      &hostName,
						},
					},
				},
			}
			result, err := kubeClient.BatchCreateNode(ctx, header, bizId, createNode)
			util.RegisterResponse(result)
			Expect(err).NotTo(HaveOccurred())
			nodeID = result[0]
		}()

		By("create node")
		func() {
			node := "node2"
			unschedulable := false
			hostName := "hostname"
			internalIP := []string{"1.1.1.1", "2.2.2.2"}
			externalIP := []string{"3.3.3.3", "4.4.4.4"}
			createNode := &types.CreateNodesOption{
				Nodes: []types.OneNodeCreateOption{
					{
						HostID:    hostId1,
						ClusterID: clusterID,
						Node: types.Node{
							Name:          &node,
							Unschedulable: &unschedulable,
							InternalIP:    &internalIP,
							ExternalIP:    &externalIP,
							HostName:      &hostName,
						},
					},
				},
			}
			result, err := kubeClient.BatchCreateNode(ctx, header, bizId, createNode)
			util.RegisterResponse(result)
			Expect(err).NotTo(HaveOccurred())
			nodeID2 = result[0]
		}()

		By("create node without name")

		func() {
			unschedulable := false
			hostName := "hostname"
			internalIP := []string{"1.1.1.1", "2.2.2.2"}
			externalIP := []string{"3.3.3.3", "4.4.4.4"}
			createNode := &types.CreateNodesOption{
				Nodes: []types.OneNodeCreateOption{
					{
						HostID:    hostId1,
						ClusterID: clusterID,
						Node: types.Node{
							Unschedulable: &unschedulable,
							InternalIP:    &internalIP,
							ExternalIP:    &externalIP,
							HostName:      &hostName,
						},
					},
				},
			}
			result, err := kubeClient.BatchCreateNode(ctx, header, bizId, createNode)
			util.RegisterResponse(result)
			Expect(err.Error()).Should(ContainSubstring("name"))
		}()

	})

	It("update kube node", func() {
		By("update node fields")
		func() {
			internalIP := []string{"5.5.5.5", "6.6.6.6"}
			createNode := &types.UpdateNodeOption{
				IDs: []int64{nodeID},
				Data: types.Node{
					InternalIP: &internalIP,
				},
			}
			result, err := kubeClient.UpdateNodeFields(ctx, header, bizId, createNode)
			util.RegisterResponse(result)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Result).To(Equal(true))
		}()

		By("update node non-editable field")
		func() {
			name := "nodetest"
			option := &types.UpdateNodeOption{
				IDs: []int64{nodeID},
				Data: types.Node{
					Name: &name,
				},
			}
			result, err := kubeClient.UpdateNodeFields(ctx, header, bizId, option)
			util.RegisterResponse(result)
			Expect(err.Error()).Should(ContainSubstring("name"))
		}()

		By("search node by node name")

		func() {
			input := &types.QueryNodeOption{
				Filter: &filter.Expression{
					RuleFactory: &filter.AtomRule{
						Field:    types.KubeNameField,
						Operator: filter.OpFactory(filter.Equal),
						Value:    "node",
					},
				},
				Page: metadata.BasePage{
					Start: 0,
					Limit: 10,
				},
			}
			result, err := kubeClient.SearchNode(ctx, header, bizId, input)
			util.RegisterResponseWithRid(result, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Result).To(Equal(true))
		}()

		By("search node count by node name")

		func() {
			input := &types.QueryNodeOption{
				Filter: &filter.Expression{
					RuleFactory: &filter.AtomRule{
						Field:    types.KubeNameField,
						Operator: filter.OpFactory(filter.In),
						Value:    []string{"node", "node2"},
					},
				},
				Page: metadata.BasePage{
					EnableCount: true,
				},
			}
			result, err := kubeClient.SearchNode(ctx, header, bizId, input)
			util.RegisterResponseWithRid(result, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Result).To(Equal(true))

			var info rest.CountInfo
			j, _ := json.Marshal(&result.Data)
			json.Unmarshal(j, &info)
			Expect(info.Count).Should(Equal(int64(2)))

		}()

		By("delete kube node")

		func() {
			option := &types.BatchDeleteNodeOption{
				IDs: []int64{nodeID2},
			}
			rsp, err := kubeClient.BatchDeleteNode(ctx, header, bizId, option)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		}()
	})
})
