/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package host_server_test

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("host abnormal test", func() {
	ctx := context.Background()
	supplierAccount := "0"
	responses := make(map[string]interface{})
	Describe("test host apply", func() {

		It("1. CreateBusiness", func() {

			input := map[string]interface{}{
				"bk_biz_name":       util.RandSeq(16),
				"life_cycle":        "2",
				"bk_biz_maintainer": "admin",
				"bk_biz_productor":  "",
				"bk_biz_tester":     "",
				"bk_biz_developer":  "",
				"operator":          "",
				"time_zone":         "Asia/Shanghai",
				"language":          "1",
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef("/biz/%s", supplierAccount).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			mapData, err := mapstruct.Struct2Map(rsp)
			Expect(err).NotTo(HaveOccurred())
			responses["req_cedb268c4487418baedab1d08843505d"] = mapData
		})

		It("2. CreateAttributeGroup", func() {

			input := map[string]interface{}{
				"bk_group_id":         util.RandSeq(16),
				"bk_group_index":      rand.Int(),
				"bk_group_name":       util.RandSeq(16),
				"bk_obj_id":           "host",
				"bk_supplier_account": "0",
				"is_collapse":         false,
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef("/create/objectattgroup").
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_b41cdfeea4134d06b1eb558ec7cb71ac"] = rsp
		})

		It("3. CreateAttribute", func() {

			input := map[string]interface{}{
				"bk_property_name":    util.RandSeq(16),
				"bk_property_id":      util.RandSeq(16),
				"unit":                "",
				"placeholder":         "",
				"bk_property_type":    "singlechar",
				"editable":            true,
				"isrequired":          false,
				"option":              "",
				"creator":             "admin",
				"bk_property_group":   "value1",
				"bk_property_index":   0,
				"bk_obj_id":           "host",
				"bk_supplier_account": "0",
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef("/create/objectattr").
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_269f039e70864ca29c0ca7bfce344ed5"] = rsp
		})

		It("4. CreateSet", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			input := map[string]interface{}{
				"bk_set_name":         util.RandSeq(16),
				"bk_set_desc":         "",
				"bk_set_env":          "3",
				"bk_service_status":   "1",
				"description":         "",
				"bk_capacity":         nil,
				"bk_biz_id":           value1,
				"bk_parent_id":        value1,
				"bk_supplier_account": "0",
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef("/set/%d", value1).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_97c8a2eacfe243b7b910bbbb15299641"] = rsp
		})

		It("5. CreateModule", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_97c8a2eacfe243b7b910bbbb15299641", "name:$.data.bk_set_id", "{.data.bk_set_id}")

			input := map[string]interface{}{
				"bk_module_name":      util.RandSeq(16),
				"bk_biz_id":           value1,
				"bk_parent_id":        value2,
				"bk_supplier_account": "0",
				"service_template_id": 0,
				"service_category_id": 2,
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef("/module/%d/%d", value1, value2).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_ba23073a86944d0cb6386a1b2a724ff1"] = rsp
		})

		It("5. UpdateModule", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_97c8a2eacfe243b7b910bbbb15299641", "name:$.data.bk_set_id", "{.data.bk_set_id}")
			value3 := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a724ff1", "name:$.data.bk_module_id", "{.data.bk_module_id}")

			input := map[string]interface{}{
				"host_apply_enabled": true,
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Put().
				WithContext(ctx).
				Body(input).
				SubResourcef("/module/%d/%d/%d", value1, value2, value3).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_6e100c803ff54da9ab54d4888ec52f09"] = rsp
		})

		It("5. SetModuleHostApplyStatus", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a724ff1", "name:$.data.bk_module_id", "{.data.bk_module_id}")
			input := map[string]interface{}{
				"enabled":     true,
				"clear_rules": false,
				"ids":         []int64{value2},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Put().
				WithContext(ctx).
				Body(input).
				SubResourcef("/module/host_apply_enable_status/bk_biz_id/%d", value1).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_62d2f03f11e7413aada33b34389fd2f2"] = rsp
		})

		It("5.1 ImportHosts", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")

			input := map[string]interface{}{
				"bk_biz_id":  value1,
				"input_type": "excel",
				"host_info": map[string]interface{}{
					"1": map[string]interface{}{
						"bk_host_innerip": fmt.Sprintf("127.%d.0.1", value1),
						"bk_asset_id":     "addhost_excel_asset_1",
						"bk_host_name":    "127.value1.0.1",
					},
					"2": map[string]interface{}{
						"bk_host_innerip": fmt.Sprintf("127.%d.0.2", value1),
						"bk_asset_id":     "addhost_excel_asset_1",
						"bk_host_name":    "127.value1.0.2",
					},
					"3": map[string]interface{}{
						"bk_host_innerip": fmt.Sprintf("127.%d.0.3", value1),
						"bk_asset_id":     "addhost_excel_asset_1",
						"bk_host_name":    "127.value1.0.3",
					},
				},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef("/hosts/add").
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_bf73b19cf2cc4587b141476da95321ba"] = rsp
		})

		It("5.2 ListBizHosts", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")

			input := map[string]interface{}{
				"host_property_filter": map[string]interface{}{
					"condition": "AND",
					"rules": []map[string]interface{}{
						{
							"operator": "in",
							"field":    "bk_host_innerip",
							"value": []string{
								fmt.Sprintf("127.%d.0.1", value1),
								fmt.Sprintf("127.%d.0.2", value1),
								fmt.Sprintf("127.%d.0.3", value1),
							},
						},
					},
				},
				"page": map[string]interface{}{
					"limit": 10,
				},
				"fields": []string{"bk_host_id"},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef("/hosts/app/%d/list_hosts", value1).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_ad1c416acb594b8a924163dad615df7a"] = rsp
		})

		It("5.3 TransferHost", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_ad1c416acb594b8a924163dad615df7a", "name:$.data.info[0].bk_host_id", "{.data.info[0].bk_host_id}")
			value3 := util.JsonPathExtractInt(responses, "req_ad1c416acb594b8a924163dad615df7a", "name:$.data.info[1].bk_host_id", "{.data.info[1].bk_host_id}")
			value4 := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a724ff1", "name:$.data.bk_module_id", "{.data.bk_module_id}")
			urlTemplate := "/host/transfer_with_auto_clear_service_instance/bk_biz_id/%d"

			input := map[string]interface{}{
				"bk_host_ids": []interface{}{
					value2,
					value3,
				},
				"remove_from_node": map[string]interface{}{
					"bk_inst_id": value1,
					"bk_obj_id":  "biz",
				},
				"add_to_modules": []interface{}{
					value4,
				},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate, value1).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_080c967a7d59488f96f289c1f06851c2"] = rsp
		})

		It("6. CreateHostApplyRule", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a724ff1", "name:$.data.bk_module_id", "{.data.bk_module_id}")
			value3 := util.JsonPathExtractInt(responses, "req_269f039e70864ca29c0ca7bfce344ed5", "name:$.data.id", "{.data.id}")
			urlTemplate := "/create/host_apply_rule/bk_biz_id/%d"

			input := map[string]interface{}{
				"bk_module_id":      value2,
				"bk_attribute_id":   value3,
				"bk_property_value": "value",
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate, value1).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_d0479d956b6c40b5aedb6c2fbb5451d5"] = rsp
		})

		It("7. UpdateHostApplyRule", func() {
			value1 := util.JsonPathExtractInt(responses, "req_d0479d956b6c40b5aedb6c2fbb5451d5", "name:$.data.id", "{.data.id}")
			value2 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")

			input := map[string]interface{}{
				"bk_property_value": "value",
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Put().
				WithContext(ctx).
				Body(input).
				SubResourcef("/update/host_apply_rule/%d/bk_biz_id/%d", value1, value2).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_d2630b0367d74aa9ab42857995535fe6"] = rsp
		})

		It("8. UpdateHostApplyRule", func() {
			value1 := util.JsonPathExtractInt(responses, "req_d0479d956b6c40b5aedb6c2fbb5451d5", "name:$.data.id", "{.data.id}")
			value2 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			urlTemplate := "/update/host_apply_rule/%d/bk_biz_id/%d"

			input := map[string]interface{}{
				"bk_property_value": "value",
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Put().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate, value1, value2).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_d2630b0367d74aa9ab42857995535fe6"] = rsp
		})

		It("9. ListHostApplyRule", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a724ff1", "name:$.data.bk_module_id", "{.data.bk_module_id}")
			urlTemplate := "/findmany/host_apply_rule/bk_biz_id/%d"

			input := map[string]interface{}{
				"bk_module_ids": []interface{}{value2},
				"page": map[string]interface{}{
					"start": 0,
					"limit": 100,
					"sort":  "",
				},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate, value1).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_5feb4bec3e904242a36bfc259d3521ba"] = rsp
		})

		It("9.1. ListHostRelatedApplyRule", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_ad1c416acb594b8a924163dad615df7a", "name:$.data.info[0].bk_host_id", "{.data.info[0].bk_host_id}")
			value3 := util.JsonPathExtractInt(responses, "req_ad1c416acb594b8a924163dad615df7a", "name:$.data.info[1].bk_host_id", "{.data.info[1].bk_host_id}")
			urlTemplate := "/findmany/host_apply_rule/bk_biz_id/%d/host_related_rules"

			input := map[string]interface{}{
				"bk_host_ids": []interface{}{
					value2,
					value3,
				},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate, value1).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_9c42cab1e85744459c67b60b6cda9844"] = rsp
		})

		It("10. CreateAttribute", func() {
			value1 := util.JsonPathExtract(responses, "req_b41cdfeea4134d06b1eb558ec7cb71ac", "name:$.data.bk_group_id", "{.data.bk_group_id}")
			urlTemplate := "/create/objectattr"

			input := map[string]interface{}{
				"bk_property_name":    util.RandSeq(16),
				"bk_property_id":      util.RandSeq(16),
				"unit":                "",
				"placeholder":         "",
				"bk_property_type":    "singlechar",
				"editable":            true,
				"isrequired":          false,
				"option":              "",
				"creator":             "admin",
				"bk_property_group":   value1,
				"bk_property_index":   0,
				"bk_obj_id":           "host",
				"bk_supplier_account": "0",
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_9a5a5cc52e624951a6954ee2e1c6d512"] = rsp
		})

		It("11. BatchCreateOrUpdateHostApplyRule ", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a724ff1", "name:$.data.bk_module_id", "{.data.bk_module_id}")
			value3 := util.JsonPathExtractInt(responses, "req_269f039e70864ca29c0ca7bfce344ed5", "name:$.data.id", "{.data.id}")
			value4 := util.JsonPathExtractInt(responses, "req_d0479d956b6c40b5aedb6c2fbb5451d5", "name:$.data.id", "{.data.id}")
			value5 := util.JsonPathExtractInt(responses, "req_9a5a5cc52e624951a6954ee2e1c6d512", "name:$.data.id", "{.data.id}")
			urlTemplate := "/createmany/host_apply_rule/bk_biz_id/%d/batch_create_or_update"

			input := map[string]interface{}{
				"host_apply_rules": []map[string]interface{}{
					{
						"bk_module_id":       value2,
						"bk_attribute_id":    value3,
						"host_apply_rule_id": value4,
					},
					{
						"bk_attribute_id":    value5,
						"bk_module_id":       value2,
						"bk_property_value":  util.RandSeq(16),
						"host_apply_rule_id": 0,
					},
				},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate, value1).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_19a5ec9ae86f47dcbe1c9f17d12e3fc2"] = rsp
		})

		It("12. RunHostApplyPlan", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a724ff1", "name:$.data.bk_module_id", "{.data.bk_module_id}")
			urlTemplate := "/host/updatemany/module/host_apply_plan/run"

			input := map[string]interface{}{
				"bk_module_ids": []interface{}{value2},
				"bk_biz_id":     value1,
				"additional_rules": []map[string]interface{}{
					{
						"bk_attribute_id":   125,
						"bk_module_id":      value2,
						"bk_property_value": "运营中[需告警]",
					},
					{
						"bk_attribute_id":   34,
						"bk_module_id":      value2,
						"bk_property_value": "3",
					},
				},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_676540002c8547f29b7004b6620e6c53"] = rsp
		})

		It("13. GenerateModuleApplyPlan", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a724ff1", "name:$.data.bk_module_id", "{.data.bk_module_id}")
			urlTemplate := "/host/createmany/module/host_apply_plan/preview"

			input := map[string]interface{}{
				"bk_biz_id":     value1,
				"bk_module_ids": []interface{}{value2},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())

			responses["req_8585626c27b845bbae201d6f8c80be59"] = rsp
		})

		It("14. CreateModule", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_97c8a2eacfe243b7b910bbbb15299641", "name:$.data.bk_set_id", "{.data.bk_set_id}")
			urlTemplate := "/module/%d/%d"

			input := map[string]interface{}{
				"bk_module_name":      util.RandSeq(16),
				"bk_biz_id":           value1,
				"bk_parent_id":        value2,
				"bk_supplier_account": "0",
				"service_template_id": 0,
				"service_category_id": 2,
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate, value1, value2).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_4a6d78e1316a44fe90dab8c55f31a941"] = rsp
		})

		It("15 TransferHost", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_ad1c416acb594b8a924163dad615df7a", "name:$.data.info[0].bk_host_id", "{.data.info[0].bk_host_id}")
			value3 := util.JsonPathExtractInt(responses, "req_ad1c416acb594b8a924163dad615df7a", "name:$.data.info[1].bk_host_id", "{.data.info[1].bk_host_id}")
			value4 := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a724ff1", "name:$.data.bk_module_id", "{.data.bk_module_id}")
			value5 := util.JsonPathExtractInt(responses, "req_4a6d78e1316a44fe90dab8c55f31a941", "name:$.data.bk_module_id", "{.data.bk_module_id}")
			urlTemplate := "/host/transfer_with_auto_clear_service_instance/bk_biz_id/%d"

			input := map[string]interface{}{
				"bk_host_ids": []interface{}{
					value2,
					value3,
				},
				"remove_from_node": map[string]interface{}{
					"bk_inst_id": value1,
					"bk_obj_id":  "biz",
				},
				"add_to_modules": []interface{}{
					value4,
					value5,
				},
				"host_apply_trans_rule": map[string]interface{}{
					"changed": true,
					"final_rules": map[string]interface{}{
						"bk_attribute_id":   125,
						"bk_property_value": "运营中[需告警]",
					},
				},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate, value1).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_aa54110e3d354701bd0f367fb5a8f852"] = rsp
		})

		It("16. BatchCreateOrUpdateHostApplyRule ", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_9a5a5cc52e624951a6954ee2e1c6d512", "name:$.data.id", "{.data.id}")
			value3 := util.JsonPathExtractInt(responses, "req_4a6d78e1316a44fe90dab8c55f31a941", "name:$.data.bk_module_id", "{.data.bk_module_id}")
			urlTemplate := "/createmany/host_apply_rule/bk_biz_id/%d/batch_create_or_update"

			input := map[string]interface{}{
				"host_apply_rules": []map[string]interface{}{
					{
						"bk_attribute_id":    value2,
						"bk_module_id":       value3,
						"bk_property_value":  util.RandSeq(16),
						"host_apply_rule_id": 0,
					},
				},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate, value1).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_d3a0e99e1865429d883029dc1df44b6d"] = rsp
		})

		It("17. SearchHostApplyRuleRelatedModules", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_9a5a5cc52e624951a6954ee2e1c6d512", "name:$.data.id", "{.data.id}")
			value3 := util.JsonPathExtract(responses, "req_9a5a5cc52e624951a6954ee2e1c6d512", "name:$.data.bk_property_name", "{.data.bk_property_name}")
			urlTemplate := "/find/topoinst/bk_biz_id/%d/host_apply_rule_related"

			input := map[string]interface{}{
				"query_filter": map[string]interface{}{
					"condition": "OR",
					"rules": []map[string]interface{}{
						{
							"operator": "exist",
							"field":    strconv.FormatInt(value2, 10),
							"value":    value3,
						},
					},
				},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate, value1).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_6d1f79ae6ad4429b8f40689713b44e2c"] = rsp
		})

		It("18. DeleteHostApplyRule", func() {
			value1 := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			value2 := util.JsonPathExtractInt(responses, "req_d0479d956b6c40b5aedb6c2fbb5451d5", "name:$.data.id", "{.data.id}")
			value3 := util.JsonPathExtractInt(responses, "req_19a5ec9ae86f47dcbe1c9f17d12e3fc2", "name:$.data.items[1].rule.id", "{.data.items[1].rule.id}")
			value4 := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a724ff1", "name:$.data.bk_module_id", "{.data.bk_module_id}")

			urlTemplate := "/host/deletemany/module/host_apply_rule/bk_biz_id/%d"
			input := map[string]interface{}{
				"host_apply_rule_ids": []interface{}{value2, value3},
				"bk_module_ids":       []interface{}{value4},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Delete().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate, value1).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			responses["req_315b9b59bef44dcb8af27ef9f0c2cae8"] = rsp
		})

		It("prepare environment", func() {
			// create service template
			bizID := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d",
				"name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			option := map[string]interface{}{
				"service_category_id": 2,
				"bk_biz_id":           bizID,
				"name":                util.RandSeq(16),
				"host_apply_enabled":  true,
			}
			resp, err := procServerClient.CreateServiceTemplate(context.Background(), header, option)
			util.RegisterResponse(resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Result).To(Equal(true))
			responses["req_ba23073a86944d0cb6386a1b2a331wq2"] = resp

			// create module
			templateID := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a331wq2", "name:$.data.id",
				"{.data.id}")
			setID := util.JsonPathExtractInt(responses, "req_97c8a2eacfe243b7b910bbbb15299641", "name:$.data.bk_set_id",
				"{.data.bk_set_id}")

			input := map[string]interface{}{
				"bk_module_name":      util.RandSeq(16),
				"bk_parent_id":        setID,
				"service_category_id": 2,
				"service_template_id": templateID,
			}
			rsp, err := instClient.CreateModule(context.Background(), bizID, setID, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			responses["req_ba23073a86944d0cb6386a1b2a7341ed1"] = rsp

			// create module rule
			moduleID := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a7341ed1",
				"name:$.bk_module_id", "{.bk_module_id}")
			attributeID := util.JsonPathExtractInt(responses, "req_269f039e70864ca29c0ca7bfce344ed5", "name:$.data.id",
				"{.data.id}")
			moduleRuleInput := map[string]interface{}{
				"bk_module_id":      moduleID,
				"bk_attribute_id":   attributeID,
				"bk_property_value": "testModule",
			}

			moduleRuleRsp := metadata.Response{}
			err = apiServerClient.Client().Post().
				WithContext(ctx).
				Body(moduleRuleInput).
				SubResourcef("/create/host_apply_rule/bk_biz_id/%d", bizID).
				WithHeaders(header).
				Do().Into(&moduleRuleRsp)
			util.RegisterResponse(moduleRuleRsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(moduleRuleRsp.Result).To(Equal(true))

			// create service template rule
			ruleInput := map[string]interface{}{
				"service_template_id": templateID,
				"bk_attribute_id":     attributeID,
				"bk_property_value":   "testServiceTemplate",
			}

			templateRuleRsp := metadata.Response{}
			err = apiServerClient.Client().Post().
				WithContext(ctx).
				Body(ruleInput).
				SubResourcef("/create/host_apply_rule/bk_biz_id/%d", bizID).
				WithHeaders(header).
				Do().Into(&templateRuleRsp)
			util.RegisterResponse(templateRuleRsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(templateRuleRsp.Result).To(Equal(true))

			// transfer host
			host1 := util.JsonPathExtractInt(responses, "req_ad1c416acb594b8a924163dad615df7a",
				"name:$.data.info[0].bk_host_id", "{.data.info[0].bk_host_id}")
			host2 := util.JsonPathExtractInt(responses, "req_ad1c416acb594b8a924163dad615df7a",
				"name:$.data.info[1].bk_host_id", "{.data.info[1].bk_host_id}")
			transInput := map[string]interface{}{
				"bk_host_ids": []interface{}{
					host1,
					host2,
				},
				"is_remove_from_all": false,
				"add_to_modules": []interface{}{
					moduleID,
				},
			}

			transRsp := metadata.Response{}
			err = apiServerClient.Client().Post().
				WithContext(ctx).
				Body(transInput).
				SubResourcef("/host/transfer_with_auto_clear_service_instance/bk_biz_id/%d", bizID).
				WithHeaders(header).
				Do().Into(&transRsp)
			util.RegisterResponse(transRsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(transRsp.Result).To(Equal(true))
		})

		It("19. GetTemplateHostApplyStatus", func() {
			bizID := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d",
				"name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			moduleID := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a7341ed1",
				"name:$.bk_module_id", "{.bk_module_id}")
			option := metadata.GetHostApplyStatusParam{
				ApplicationID: bizID,
				ModuleIDs:     []int64{moduleID},
			}
			resp := metadata.MapArrayResponse{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(option).
				SubResourcef("/host/find/service_template/host_apply_status").
				WithHeaders(header).
				Do().Into(&resp)

			util.RegisterResponse(resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Result).To(Equal(true))
			Expect(commonutil.GetInt64ByInterface(resp.Data[0]["bk_module_id"])).To(Equal(moduleID))
			Expect(commonutil.GetStrByInterface(resp.Data[0]["host_apply_enabled"])).To(Equal("true"))
		})

		It("20. GenerateTemplateApplyPlan", func() {
			bizID := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id",
				"{.data.bk_biz_id}")
			templateID := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a331wq2", "name:$.data.id",
				"{.data.id}")
			urlTemplate := "/host/createmany/service_template/host_apply_plan/preview"

			input := map[string]interface{}{
				"bk_biz_id":            bizID,
				"service_template_ids": []interface{}{templateID},
			}

			rsp := metadata.Response{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef(urlTemplate).
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("21. GetServiceTemplateHostApplyRule", func() {
			bizID := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id",
				"{.data.bk_biz_id}")
			templateID := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a331wq2", "name:$.data.id",
				"{.data.id}")

			input := map[string]interface{}{
				"bk_biz_id":            bizID,
				"service_template_ids": []interface{}{templateID},
				"page": metadata.BasePage{
					Limit: 1,
				},
			}

			rsp := metadata.ResponseInstData{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef("/host/findmany/service_template/host_apply_rule").
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Info[0].Int64("service_template_id")).To(Equal(templateID))
		})

		It("22. GetServiceTemplateInvalidHostCount", func() {
			bizID := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id",
				"{.data.bk_biz_id}")
			templateID := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a331wq2", "name:$.data.id",
				"{.data.id}")

			input := map[string]interface{}{
				"bk_biz_id": bizID,
				"id":        templateID,
			}

			rsp := metadata.ResponseDataMapStr{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef("/host/findmany/service_template/host_apply_plan/invalid_host_count").
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Int64("count")).To(Equal(int64(2)))
		})

		It("23. GetModuleInvalidHostCount", func() {
			bizID := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d", "name:$.data.bk_biz_id",
				"{.data.bk_biz_id}")
			moduleID := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a7341ed1",
				"name:$.bk_module_id", "{.bk_module_id}")

			input := map[string]interface{}{
				"bk_biz_id": bizID,
				"id":        moduleID,
			}

			rsp := metadata.ResponseDataMapStr{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(input).
				SubResourcef("/host/findmany/module/host_apply_plan/invalid_host_count").
				WithHeaders(header).
				Do().Into(&rsp)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Int64("count")).To(Equal(int64(2)))
		})

		It("24. GetServiceTemplateHostApplyRuleCount", func() {
			bizID := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d",
				"name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			templateID := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a331wq2", "name:$.data.id",
				"{.data.id}")
			option := metadata.HostApplyRuleCountOption{
				ApplicationID:      bizID,
				ServiceTemplateIDs: []int64{templateID},
			}
			resp := metadata.MapArrayResponse{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(option).
				SubResourcef("/host/findmany/service_template/host_apply_rule_count").
				WithHeaders(header).
				Do().Into(&resp)

			util.RegisterResponse(resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Result).To(Equal(true))
			Expect(commonutil.GetInt64ByInterface(resp.Data[0]["service_template_id"])).To(Equal(templateID))
			Expect(commonutil.GetInt64ByInterface(resp.Data[0]["count"])).To(Equal(int64(1)))
		})

		It("25. GetModuleFinalRules", func() {
			bizID := util.JsonPathExtractInt(responses, "req_cedb268c4487418baedab1d08843505d",
				"name:$.data.bk_biz_id", "{.data.bk_biz_id}")
			moduleID := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a7341ed1",
				"name:$.bk_module_id", "{.bk_module_id}")
			templateID := util.JsonPathExtractInt(responses, "req_ba23073a86944d0cb6386a1b2a331wq2", "name:$.data.id",
				"{.data.id}")
			option := metadata.ModuleFinalRulesParam{
				ApplicationID: bizID,
				ModuleIDs:     []int64{moduleID},
			}
			resp := metadata.MapArrayResponse{}
			err := apiServerClient.Client().Post().
				WithContext(ctx).
				Body(option).
				SubResourcef("/host/findmany/module/get_module_final_rules").
				WithHeaders(header).
				Do().Into(&resp)

			util.RegisterResponse(resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Result).To(Equal(true))
			Expect(commonutil.GetInt64ByInterface(resp.Data[0]["service_template_id"])).To(Equal(templateID))
		})
	})
})
