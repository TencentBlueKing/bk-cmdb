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

package container_server_test

import (
	"context"

	"configcenter/src/common/metadata"
	"configcenter/src/framework/core/input"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("container test", func() {
	var bizId, setId, moduleId int64

	Describe("test preparation", func() {
		It("create business bk_biz_name = 'cc_biz'", func() {
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

		It("create set", func() {
			input := mapstr.MapStr{
				"bk_set_name":         "test",
				"bk_parent_id":        bizId,
				"bk_supplier_account": "0",
				"bk_biz_id":           bizId,
				"bk_service_status":   "1",
				"bk_set_env":          "3",
			}
			rsp, err := instClient.CreateSet(context.Background(), strconv.FormatInt(bizId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data["bk_set_name"].(string)).To(Equal("test"))
			parentIdRes, err := commonutil.GetInt64ByInterface(rsp.Data["bk_parent_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(parentIdRes).To(Equal(bizId))
			bizIdRes, err := commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(bizIdRes).To(Equal(bizId))
			setId, err = commonutil.GetInt64ByInterface(rsp.Data["bk_set_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("create module", func() {
			input := map[string]interface{}{
				"bk_module_name":      "cc_module",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data["bk_module_name"].(string)).To(Equal("cc_module"))

			setIdRes, err := commonutil.GetInt64ByInterface(rsp.Data["bk_set_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(setIdRes).To(Equal(setId))

			parentIdRes, err := commonutil.GetInt64ByInterface(rsp.Data["bk_parent_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(parentIdRes).To(Equal(setId))
			moduleId, err = commonutil.GetInt64ByInterface(rsp.Data["bk_module_id"])
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("add pod test", func() {
		It("add pod", func() {
			input := metadata.CreatePod{
				"pod": map[string]interface{}{
					"bk_pod_name":      "bcs_pod",
					"bk_pod_namespace": "bcs_namespace",
					"bk_pod_cluster":   "bcs_cluster",
					"bk_pod_uuid":      "bcs_uuid",
					"bk_cloud_id":      0,
					"bk_host_innerip":  "1.0.0.1",
				},
			}
			resp, err := containerServerClient.CreatePod(context.Background(), header, bizId, input)
			util.RegisterResponse(resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Result).To(Equal(true), resp.ToString())
		})

		It("list pod", func() {
			input := metadata.ListPods{
				BizID:     bizId,
				SetIDs:    []int64{setId},
				ModuleIDs: []int64{moduleId},
			}
			resp, err := containerServerClient.ListPods(context.Background(), header, bizId, input)
			util.RegisterResponse(resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Result).To(Equal(true), resp.ToString())
			Expect(resp.Count).To(Equal(1))
			Expect(resp.Info[0].(map[string]interface{})["bk_pod_name"].(string)).To(Equal("bcs_pod"))
			Expect(resp.Info[0].(map[string]interface{})["bk_pod_namespace"].(string)).To(Equal("bcs_namespace"))
			Expect(resp.Info[0].(map[string]interface{})["bk_pod_cluster"].(string)).To(Equal("bcs_cluster"))
			Expect(resp.Info[0].(map[string]interface{})["bk_pod_uuid"].(string)).To(Equal("bcs_uuid"))
			Expect(resp.Info[0].(map[string]interface{})["bk_cloud_id"].(int)).To(Equal(0))
			Expect(resp.Info[0].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("1.0.0.1"))
		})
	})
})
