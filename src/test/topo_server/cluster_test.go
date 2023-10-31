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
	commonutil "configcenter/src/common/util"
	"configcenter/src/kube/types"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kube cluster test", func() {
	ctx := context.Background()
	var clusterID, clusterID2 int64
	Describe("test preparation", func() {
		It("create business bk_biz_name = 'cc_biz'", func() {
			test.DeleteAllBizs()

			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_name":       "cc_biz",
				"time_zone":         "Africa/Accra",
			}
			rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data).To(ContainElement("cc_biz"))
			bizID, err = commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("create kube cluster", func() {

		By("create cluster")
		func() {
			clusterName := "cluster"
			schedulingEngine := "k8s"
			uid := "BCS-K8S-25001"
			xid := "cls-hox2lkf2"
			version := "0.1"
			networkType := "underlay"
			region := "shenzhen"
			vpc := "vpc-q6awe02n"
			environment := "prod"
			network := []string{"1.1.1.0/21"}
			clusterType := types.IndependentClusterType
			createCluster := &types.Cluster{
				BizID:            bizID,
				Name:             &clusterName,
				SchedulingEngine: &schedulingEngine,
				Uid:              &uid,
				Xid:              &xid,
				Environment:      &environment,
				Version:          &version,
				NetworkType:      &networkType,
				Region:           &region,
				Vpc:              &vpc,
				NetWork:          &network,
				Type:             &clusterType,
			}

			id, err := kubeClient.CreateCluster(ctx, header, createCluster)
			util.RegisterResponseWithRid(id, header)
			Expect(err).NotTo(HaveOccurred())
			clusterID = id
		}()

		By("create cluster again")

		func() {
			clusterName := "cluster1"
			schedulingEngine := "k8s"
			uid := "BCS-K8S-250011"
			xid := "cls-hox2lkf21"
			version := "0.11"
			networkType := "underlay"
			environment := "prod"
			region := "shenzhen"
			vpc := "vpc-q6awe02n"
			network := []string{"1.1.1.0/21"}
			clusterType := types.IndependentClusterType
			createCluster := &types.Cluster{
				BizID:            bizID,
				Name:             &clusterName,
				SchedulingEngine: &schedulingEngine,
				Uid:              &uid,
				Xid:              &xid,
				Version:          &version,
				Environment:      &environment,
				NetworkType:      &networkType,
				Region:           &region,
				Vpc:              &vpc,
				NetWork:          &network,
				Type:             &clusterType,
			}

			id, err := kubeClient.CreateCluster(ctx, header, createCluster)
			Expect(err).NotTo(HaveOccurred())
			clusterID2 = id
		}()

		By("create kube cluster without cluster name")

		func() {
			schedulingEngine := "k8s"
			uid := "BCS-K8S-250012"
			xid := "cls-hox2lkf21"
			version := "0.11"
			networkType := "underlay"
			region := "shenzhen"
			vpc := "vpc-q6awe02n"
			environment := "prod"
			network := []string{"1.1.1.0/21"}
			clusterType := types.SharedClusterType
			createCluster := &types.Cluster{
				BizID:            bizID,
				SchedulingEngine: &schedulingEngine,
				Uid:              &uid,
				Xid:              &xid,
				Version:          &version,
				Environment:      &environment,
				NetworkType:      &networkType,
				Region:           &region,
				Vpc:              &vpc,
				NetWork:          &network,
				Type:             &clusterType,
			}

			id, err := kubeClient.CreateCluster(ctx, header, createCluster)
			util.RegisterResponseWithRid(id, header)
			Expect(err.Error()).Should(ContainSubstring("name"))
		}()
	})

	It("update kube cluster field", func() {

		By("update cluster version")
		func() {
			version := "0.2"
			data := &types.UpdateClusterOption{
				BizID: bizID,
				UpdateClusterByIDsOption: types.UpdateClusterByIDsOption{
					IDs: []int64{clusterID},
					Data: types.Cluster{
						Version: &version,
					},
				},
			}
			err := kubeClient.UpdateClusterFields(ctx, header, data)
			Expect(err).NotTo(HaveOccurred())
		}()

		By("update kube cluster non-editable field")

		func() {
			uid := "uid"
			data := &types.UpdateClusterOption{
				BizID: bizID,
				UpdateClusterByIDsOption: types.UpdateClusterByIDsOption{
					IDs: []int64{clusterID},
					Data: types.Cluster{
						Uid: &uid,
					},
				},
			}
			err := kubeClient.UpdateClusterFields(ctx, header, data)
			util.RegisterResponseWithRid(err, header)
			Expect(err.Error()).Should(ContainSubstring("uid"))
		}()

		By("delete kube cluster")

		func() {
			option := &types.DeleteClusterOption{
				BizID: bizID,
				DeleteClusterByIDsOption: types.DeleteClusterByIDsOption{
					IDs: []int64{clusterID2},
				},
			}
			err := kubeClient.DeleteCluster(ctx, header, option)
			Expect(err).NotTo(HaveOccurred())
		}()

		By("search kube cluster by name")
		func() {
			input := &types.QueryClusterOption{
				BizID: bizID,
				Filter: &filter.Expression{
					RuleFactory: &filter.CombinedRule{
						Condition: filter.And,
						Rules: []filter.RuleFactory{
							&filter.AtomRule{
								Field:    types.KubeNameField,
								Operator: filter.OpFactory(filter.Equal),
								Value:    "cluster",
							},
						},
					},
				},
				Page: metadata.BasePage{
					Start: 0,
					Limit: 10,
				},
			}
			result, err := kubeClient.SearchCluster(ctx, header, input)
			util.RegisterResponseWithRid(result, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Result).To(Equal(true))

			var info rest.CountInfo
			j, _ := json.Marshal(&result.Data)
			json.Unmarshal(j, &info)
			Expect(info.Count).Should(Equal(int64(0)))

		}()

		By("search kube count cluster by name")

		func() {
			input := &types.QueryClusterOption{
				BizID: bizID,
				Filter: &filter.Expression{
					RuleFactory: &filter.AtomRule{
						Field:    types.KubeNameField,
						Operator: filter.OpFactory(filter.Equal),
						Value:    "cluster",
					},
				},
				Page: metadata.BasePage{
					EnableCount: true,
				},
			}
			result, err := kubeClient.SearchCluster(ctx, header, input)
			util.RegisterResponseWithRid(result, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Result).To(Equal(true))

			var info rest.CountInfo
			j, _ := json.Marshal(&result.Data)
			json.Unmarshal(j, &info)
			Expect(info.Count).Should(Equal(int64(1)))
		}()
	})
})
