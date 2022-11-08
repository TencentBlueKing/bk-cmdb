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
	commonutil "configcenter/src/common/util"
	"configcenter/src/kube/types"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("workload test", func() {
	ctx := context.Background()

	var bizID int64
	var clusterID int64
	clusterUID := "clusterUID"
	var namespaceID int64
	nsName := "nsName"
	var wlID int64
	wlName := "wlName"
	It("prepare environment, create business, cluster, namespace", func() {
		test.ClearDatabase()

		biz := map[string]interface{}{
			common.BKMaintainersField: "kube",
			common.BKAppNameField:     "biz_for_kube",
			"time_zone":               "Africa/Accra",
			"language":                "1",
		}
		bizResp, err := apiServerClient.CreateBiz(ctx, "0", header, biz)
		util.RegisterResponseWithRid(bizResp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(bizResp.Result).To(Equal(true))
		bizID, err = commonutil.GetInt64ByInterface(bizResp.Data[common.BKAppIDField])
		Expect(err).NotTo(HaveOccurred())

		clusterName := "cluster"
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

		ns := types.Namespace{
			ClusterSpec: types.ClusterSpec{
				ClusterID: clusterID,
			},
			Name: nsName,
		}
		createOpt := types.NsCreateOption{
			Data: []types.Namespace{ns},
		}

		wlResult, err := kubeClient.CreateNamespace(ctx, header, bizID, &createOpt)
		util.RegisterResponseWithRid(wlResult, header)
		Expect(err).NotTo(HaveOccurred())
		namespaceID = wlResult.IDs[0]
	})

	It("create workload", func() {
		label := map[string]string{
			"test":  "test",
			"test2": "test2",
		}

		selector := types.LabelSelector{
			MatchLabels: map[string]string{
				"test":  "test",
				"test2": "test2",
			},
			MatchExpressions: []types.LabelSelectorRequirement{
				{
					Key:      "tier",
					Operator: "In",
					Values:   []string{"cache"},
				},
			},
		}

		var replicas int64 = 1
		strategyType := types.RollingUpdateDeploymentStrategyType
		var minReadySeconds int64 = 1
		rollingUpdateStrategy := types.RollingUpdateDeployment{
			MaxUnavailable: &types.IntOrString{
				Type:   types.IntType,
				IntVal: 2,
			},
			MaxSurge: &types.IntOrString{
				Type:   types.StringType,
				StrVal: "12",
			},
		}
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
			Labels:                &label,
			Selector:              &selector,
			Replicas:              &replicas,
			MinReadySeconds:       &minReadySeconds,
			StrategyType:          &strategyType,
			RollingUpdateStrategy: &rollingUpdateStrategy,
		}
		createOpt := types.WlCreateOption{
			Data: []types.WorkloadInterface{
				&wl,
			},
		}

		var err error
		result, err := kubeClient.CreateWorkload(ctx, header, bizID, types.KubeDeployment, &createOpt)
		util.RegisterResponseWithRid(result, header)
		Expect(err).NotTo(HaveOccurred())
		wlID = result.IDs[0]
	})

	It("update workload", func() {
		label := map[string]string{
			"test": "test",
		}

		selector := types.LabelSelector{
			MatchLabels: map[string]string{
				"test": "test",
			},
			MatchExpressions: []types.LabelSelectorRequirement{
				{
					Key:      "tier",
					Operator: "In",
					Values:   []string{"cache"},
				},
			},
		}

		var replicas int64 = 2
		strategyType := types.RollingUpdateDeploymentStrategyType
		var minReadySeconds int64 = 3
		rollingUpdateStrategy := types.RollingUpdateDeployment{
			MaxUnavailable: &types.IntOrString{
				Type:   types.StringType,
				StrVal: "2",
			},
			MaxSurge: &types.IntOrString{
				Type:   types.StringType,
				StrVal: "12",
			},
		}
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
			Labels:                &label,
			Selector:              &selector,
			Replicas:              &replicas,
			MinReadySeconds:       &minReadySeconds,
			StrategyType:          &strategyType,
			RollingUpdateStrategy: &rollingUpdateStrategy,
		}
		updateOpt := types.WlUpdateOption{
			IDs:  []int64{wlID},
			Data: &wl,
		}

		err := kubeClient.UpdateWorkload(ctx, header, bizID, types.KubeDeployment, &updateOpt)
		Expect(err).NotTo(HaveOccurred())
	})

	It("find workload", func() {
		filter := &filter.Expression{
			RuleFactory: &filter.CombinedRule{
				Condition: filter.And,
				Rules: []filter.RuleFactory{
					&filter.AtomRule{
						Field:    types.BKClusterIDFiled,
						Operator: filter.Equal.Factory(),
						Value:    clusterID,
					},
					&filter.AtomRule{
						Field:    types.BKNamespaceIDField,
						Operator: filter.Equal.Factory(),
						Value:    namespaceID,
					},
					&filter.AtomRule{
						Field:    common.BKFieldName,
						Operator: filter.Equal.Factory(),
						Value:    wlName,
					},
					&filter.AtomRule{
						Field:    "labels",
						Operator: filter.Object.Factory(),
						Value: &filter.AtomRule{
							Field:    "test",
							Operator: filter.Equal.Factory(),
							Value:    "test",
						},
					},
				},
			},
		}

		// get workload data
		page := metadata.BasePage{
			Start: 0,
			Limit: 10,
		}
		fields := []string{common.BKFieldID}
		queryOpt := types.WlQueryOption{
			Filter: filter,
			Page:   page,
			Fields: fields,
		}
		result, err := kubeClient.ListWorkload(ctx, header, bizID, "deployment", &queryOpt)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result.Info)).To(Equal(1))
		Expect(result.Info[0][common.BKFieldID].(json.Number).Int64()).To(Equal(wlID))

		// get workload count
		page = metadata.BasePage{
			EnableCount: true,
		}
		queryOpt = types.WlQueryOption{
			Filter: filter,
			Page:   page,
		}
		result, err = kubeClient.ListWorkload(ctx, header, bizID, "deployment", &queryOpt)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Count).To(Equal(1))
	})

	It("delete workload", func() {
		deleteOpt := types.WlDeleteOption{
			IDs: []int64{wlID},
		}

		err := kubeClient.DeleteWorkload(ctx, header, bizID, types.KubeDeployment, &deleteOpt)
		Expect(err).NotTo(HaveOccurred())
	})
})
