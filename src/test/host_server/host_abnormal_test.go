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
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	noExistID    int64  = 99999
	noExistID1   int64  = 99998
	noExistIDStr string = "99999"
)
var bizId, bizId1, setId, setId1, moduleId, moduleId1, moduleId2, idleModuleId, faultModuleId, defaultDirID int64
var hostId, hostId1, hostId2, hostId3, hostId4 int64

var _ = Describe("host abnormal test", func() {

	Describe("test preparation", func() {
		It("create business bk_biz_name = 'Christina'", func() {
			test.ClearDatabase()

			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_name":       "Christina",
				"time_zone":         "Africa/Accra",
			}
			rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			bizId, err = commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("create business bk_biz_name = 'Angela'", func() {
			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_name":       "Angela",
				"time_zone":         "Africa/Accra",
			}
			rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			bizId1, err = commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
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
			setId, err = commonutil.GetInt64ByInterface(rsp.Data["bk_set_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("create set", func() {
			input := mapstr.MapStr{
				"bk_set_name":         "test",
				"bk_parent_id":        bizId1,
				"bk_supplier_account": "0",
				"bk_biz_id":           bizId1,
				"bk_service_status":   "1",
				"bk_set_env":          "3",
			}
			rsp, err := instClient.CreateSet(context.Background(), strconv.FormatInt(bizId1, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			setId1, err = commonutil.GetInt64ByInterface(rsp.Data["bk_set_id"])
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
			moduleId, err = commonutil.GetInt64ByInterface(rsp.Data["bk_module_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("create module", func() {
			input := map[string]interface{}{
				"bk_module_name":      "cc_module1",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			moduleId1, err = commonutil.GetInt64ByInterface(rsp.Data["bk_module_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("create module", func() {
			input := map[string]interface{}{
				"bk_module_name":      "cc_module1",
				"bk_parent_id":        setId1,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId1, 10), strconv.FormatInt(setId1, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			moduleId2, err = commonutil.GetInt64ByInterface(rsp.Data["bk_module_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("get instance topo", func() {
			rsp, err := instClient.GetInternalModule(context.Background(), "0", strconv.FormatInt(bizId1, 10), header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.SetName).To(Equal("空闲机池"))
			Expect(len(rsp.Data.Module)).To(Equal(3))
			for _, module := range rsp.Data.Module {
				switch int(module.Default) {
				case common.DefaultResModuleFlag:
					idleModuleId = module.ModuleID
				case common.DefaultFaultModuleFlag:
					faultModuleId = module.ModuleID
				}
			}
		})

		// 云区域ID不存在新加主机报错
		It("add host using api with noexist cloud_id", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"host_info": map[string]interface{}{
					"4": map[string]interface{}{
						"bk_host_innerip": "127.0.1.1",
						"bk_cloud_id":     noExistID,
					},
				},
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("get resource pool default module id", func() {
			cond := map[string]interface{}{
				common.BKDefaultField: common.DefaultResModuleFlag,
			}
			dirRsp, err := dirClient.SearchResourceDirectory(context.Background(), header, cond)
			util.RegisterResponse(dirRsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(dirRsp.Result).To(Equal(true))
			Expect(len(dirRsp.Data.Info)).To(Equal(1))

			defaultDirID, err = commonutil.GetInt64ByInterface(dirRsp.Data.Info[0][common.BKModuleIDField])
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("add host test", func() {

		//清空数据
		BeforeEach(func() {
			clearData()
		})

		Describe("add host using api test", func() {
			//测试用例运行后，主机数量应为0
			AfterEach(func() {
				// 查询业务下的主机
				input := &params.HostCommonSearch{
					AppID: int(bizId),
					Page: params.PageInfo{
						Sort: "bk_host_id",
					},
				}
				rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				Expect(rsp.Data.Count).To(Equal(0))
			})

			It("add host using api with noexist biz_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": noExistID,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.1",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with invalid biz_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": "test",
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.1",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with invalid cloud_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.1",
							"bk_cloud_id":     -1,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with invalid bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "333.0.0.1",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with invalid bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.e",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with no host_info", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with no bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_cloud_id": 0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

		})

		Describe("add host using api test2", func() {
			// 测试用例运行后，主机数量应为1
			AfterEach(func() {
				// 查询业务下的主机
				input := &params.HostCommonSearch{
					AppID: int(bizId),
					Page: params.PageInfo{
						Sort: "bk_host_id",
					},
				}
				rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				Expect(rsp.Data.Count).To(Equal(1))
			})

			// 如果云区域ID没有给出，默认是0
			It("add host using api with no bk_cloud_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.1",
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("add host using api to biz", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.1",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("add host using api to biz twice", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.1",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))

				rsp, err = hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})
		})

		Describe("add host using api test3", func() {
			It("add host using api to biz multiple ip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.13",
							"bk_cloud_id":     0,
						},
						"5": map[string]interface{}{
							"bk_host_innerip": "127.0.1.14",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})
		})

		Describe("add host using excel test", func() {
			//测试用例运行后，主机数量应为0
			AfterEach(func() {
				// 查询业务下的主机
				input := &params.HostCommonSearch{
					AppID: int(bizId),
					Page: params.PageInfo{
						Sort: "bk_host_id",
					},
				}
				rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				Expect(rsp.Data.Count).To(Equal(0))
			})

			It("add host using excel with noexist biz_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": noExistID,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.1",
							"bk_cloud_id":     0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with invalid biz_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": "test",
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.1",
							"bk_cloud_id":     0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with invalid cloud_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.1",
							"bk_cloud_id":     -1,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with invalid bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "333.0.0.1",
							"bk_cloud_id":     0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with invalid bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.e",
							"bk_cloud_id":     0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with no host_info", func() {
				input := map[string]interface{}{
					"bk_biz_id":  bizId,
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with no bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_cloud_id": 0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

		})

		Describe("add host using excel test2", func() {
			//测试用例运行后，主机数量应为1
			AfterEach(func() {
				// 查询业务下的主机
				input := &params.HostCommonSearch{
					AppID: int(bizId),
					Page: params.PageInfo{
						Sort: "bk_host_id",
					},
				}
				rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				Expect(rsp.Data.Count).To(Equal(1))
			})

			It("add host using excel with no bk_cloud_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.1",
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("add host using excel to biz", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.1",
							"bk_cloud_id":     0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("add host using excel to biz twice", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "127.0.1.1",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))

				rsp, err = hostServerClient.AddHost(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})
		})
	})

	Describe("search host test", func() {

		// 准备数据
		JustBeforeEach(func() {
			prepareData()
		})

		It("search host using invalid bk_host_id", func() {
			rsp, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0", "eceer", header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search host using noexist bk_host_id", func() {
			rsp, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0", noExistIDStr, header)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search biz host using noexist biz id", func() {
			input := &params.HostCommonSearch{
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
				Condition: []params.SearchCondition{
					{
						ObjectID: "biz",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_biz_id",
								"operator": "$eq",
								"value":    noExistID,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search host using invalid ip", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Ip: params.IPInfo{
					Data:  []string{"127.0.0"},
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
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search host using multiple ips with an invalid value", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Ip: params.IPInfo{
					Data: []string{
						"127.0.1.1",
						"127.0.1",
						"127.0.1.2",
					},
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
			Expect(rsp.Data.Count).To(Equal(2))
		})
	})

	Describe("transfer host test", func() {

		// 清空数据
		BeforeEach(func() {
			clearData()
		})

		// 准备数据
		JustBeforeEach(func() {
			prepareData()
		})

		It("transfer resourcehost to nonexist biz's idlemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: noExistID,
				HostIDs: []int64{
					hostId4,
				},
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer biz host to other biz's idlemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostIDs: []int64{
					hostId,
				},
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer resourcehost to idlemodule less biz id", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				HostIDs: []int64{
					hostId,
				},
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer resourcehost to idlemodule less host ids", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer resourcehost to idlemodule noexist host ids", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostIDs: []int64{
					hostId4,
					noExistID,
				},
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to other biz's module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_host_id": []int64{
					hostId2,
				},
				"bk_module_id": []int64{
					moduleId1,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to nonexist module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId2,
				},
				"bk_module_id": []int64{
					noExistID,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer nonexist host to module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_host_id": []int64{
					noExistID,
				},
				"bk_module_id": []int64{
					moduleId2,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer multiple nonexist host to module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_host_id": []int64{
					noExistID,
					noExistID,
				},
				"bk_module_id": []int64{
					moduleId2,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to module less biz_id", func() {
			input := map[string]interface{}{
				"bk_host_id": []int64{
					hostId3,
				},
				"bk_module_id": []int64{
					moduleId2,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to module less host_id", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_module_id": []int64{
					moduleId2,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to module less module_id", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_host_id": []int64{
					hostId3,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		// 可以不传is_increment，默认为false
		It("transfer host to module less is_increment", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_module_id": []int64{
					moduleId2,
				},
				"bk_host_id": []int64{
					hostId3,
				},
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("transfer multiple hosts with a nonexist host to module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_host_id": []int64{
					hostId1,
					noExistID,
				},
				"bk_module_id": []int64{
					moduleId2,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))

			// 查看转移后的模块主机数量, 因之前转移时有不存在的hostid，导致报错，使用事务则会回滚，故无任何主机能被转移成功
			input1 := &params.HostCommonSearch{
				AppID: int(bizId1),
				Condition: []params.SearchCondition{
					{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    moduleId2,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp1, err := hostServerClient.SearchHost(context.Background(), header, input1)
			util.RegisterResponse(rsp1)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp1.Result).To(Equal(true))
			Expect(rsp1.Data.Count).To(Equal(0))

		})

		It("transfer multiple host to module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_host_id": []int64{
					hostId3,
					hostId1,
				},
				"bk_module_id": []int64{
					moduleId2,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			// 查看转移后的模块主机数量
			input1 := &params.HostCommonSearch{
				AppID: int(bizId1),
				Condition: []params.SearchCondition{
					{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    moduleId2,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp1, err := hostServerClient.SearchHost(context.Background(), header, input1)
			util.RegisterResponse(rsp1)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp1.Result).To(Equal(true))
			Expect(rsp1.Data.Count).To(Equal(2))
		})

		It("transfer multiple host to noexist module in a biz", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_host_id": []int64{
					hostId3,
					hostId1,
				},
				"bk_module_id": []int64{
					moduleId1,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))

			input1 := &params.HostCommonSearch{
				AppID: int(bizId1),
				Condition: []params.SearchCondition{
					{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    moduleId1,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp1, err := hostServerClient.SearchHost(context.Background(), header, input1)
			util.RegisterResponse(rsp1)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp1.Result).To(Equal(true))
			Expect(rsp1.Data.Count).To(Equal(0))
		})

		It("transfer host to idle module nonexist biz", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: noExistID,
				HostIDs: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to idle module one nonexist host", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostIDs: []int64{
					hostId1,
					noExistID,
				},
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		// 查询条件导致找不到主机，返回true，不会造成脏数据
		It("move nonexist biz's module hosts to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: noExistID,
				SetID:         setId,
				ModuleID:      moduleId,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		// 查询条件导致找不到主机，返回true，不会造成脏数据
		It("move nonexist module's hosts to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: bizId1,
				SetID:         setId,
				ModuleID:      noExistID,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		// 查询条件导致找不到主机，返回true，不会造成脏数据
		It("move nonexist set's hosts to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: bizId1,
				SetID:         noExistID,
				ModuleID:      moduleId,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		// 查询条件导致找不到主机，返回true，不会造成脏数据
		It("move unmatching set module relationship hosts to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: bizId1,
				SetID:         setId1,
				ModuleID:      moduleId,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("move host whose set and module is 0 to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: bizId1,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to idle module less hostid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to idle module empty hostid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostIDs:       []int64{},
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to idle module less bizid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				HostIDs: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to nonexist biz's fault module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: noExistID,
				HostIDs: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer nonexist host to fault module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostIDs: []int64{
					noExistID,
				},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		// 主机ID不存在导致返回结果为false，因使用事务导致回滚，故无任何主机转移成功
		It("transfer a nonexist host to fault module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostIDs: []int64{
					hostId1,
					noExistID,
				},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))

			input1 := &params.HostCommonSearch{
				AppID: int(bizId1),
				Condition: []params.SearchCondition{
					{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    faultModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp1, err := hostServerClient.SearchHost(context.Background(), header, input1)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp1.Result).To(Equal(true))
			Expect(rsp1.Data.Count).To(Equal(0))
		})

		It("transfer host to fault module less hostid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to fault module less bizid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				HostIDs: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer multiple hosts to fault module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostIDs: []int64{
					hostId1,
					hostId3,
				},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			input1 := &params.HostCommonSearch{
				AppID: int(bizId1),
				Condition: []params.SearchCondition{
					{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    faultModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp1, err := hostServerClient.SearchHost(context.Background(), header, input1)
			util.RegisterResponse(rsp1)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp1.Result).To(Equal(true))
			Expect(rsp1.Data.Count).To(Equal(2))

		})

		It("transfer unmatching biz host to resourcemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostIDs: []int64{
					hostId1,
				},
				ModuleID: defaultDirID,
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer nonexist biz host to resourcemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: noExistID,
				HostIDs: []int64{
					hostId,
				},
				ModuleID: defaultDirID,
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer a nonexist host to resourcemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostIDs: []int64{
					hostId,
					noExistID,
				},
				ModuleID: defaultDirID,
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer nonidle host to resourcemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostIDs: []int64{
					hostId1,
				},
				ModuleID: defaultDirID,
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("transfer host to resourcemodule less hostid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				ModuleID:      defaultDirID,
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("transfer host to resourcemodule less bizid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				HostIDs: []int64{
					hostId1,
				},
				ModuleID: defaultDirID,
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer multiple hosts to idle module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostIDs: []int64{
					hostId1,
					hostId3,
				},
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			input1 := &params.HostCommonSearch{
				AppID: int(bizId1),
				Condition: []params.SearchCondition{
					{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp1, err := hostServerClient.SearchHost(context.Background(), header, input1)
			util.RegisterResponse(rsp1)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp1.Result).To(Equal(true))
			Expect(rsp1.Data.Count).To(Equal(2))
		})

		It("transfer multiple hosts to resource module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostIDs: []int64{
					hostId1,
					hostId3,
				},
				ModuleID: defaultDirID,
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			input1 := &params.HostCommonSearch{
				AppID: -1,
				Condition: []params.SearchCondition{
					{
						ObjectID: "biz",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "default",
								"operator": "$eq",
								"value":    1,
							},
						},
						Fields: []string{},
					},
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp1, err := hostServerClient.SearchHost(context.Background(), header, input1)
			util.RegisterResponse(rsp1)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp1.Result).To(Equal(true))
			Expect(rsp1.Data.Count).To(Equal(3))
		})
	})

	PDescribe("sync host test", func() {

		// 清空数据
		BeforeEach(func() {
			clearData()
		})

		// 准备数据
		JustBeforeEach(func() {
			prepareData()
		})

		It("sync host less biz", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.1.16",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
				},
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host less host", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_module_id": []int64{
					moduleId,
				},
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host less module", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.1.17",
						"bk_cloud_id":     0,
					},
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host empty module", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.1.17",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{},
				"bk_biz_id":    bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host invalid host", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.1",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host invalid biz", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.1.18",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"bk_biz_id": noExistID,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host invalid module", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.1.19",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					noExistID,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host one invalid module", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.1.20",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
					noExistID,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host one invalid host", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.1.21",
						"bk_cloud_id":     0,
					},
					"1": map[string]interface{}{
						"bk_host_innerip": "127.0.1",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host one invalid host", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.1.3",
						"bk_cloud_id":     0,
					},
					"1": map[string]interface{}{
						"bk_host_innerip": "127.0.1",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host duplicate module", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.1.5",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
					moduleId,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId1),
				Condition: []params.SearchCondition{
					{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    moduleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("sync multiple host multiple module", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.1.8",
						"bk_cloud_id":     0,
					},
					"1": map[string]interface{}{
						"bk_host_innerip": "127.0.1.9",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
					moduleId1,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			input1 := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					{
						ObjectID: "module",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_module_id",
								"operator": "$eq",
								"value":    moduleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp1, err := hostServerClient.SearchHost(context.Background(), header, input1)
			util.RegisterResponse(rsp1)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp1.Result).To(Equal(true))
			Expect(rsp1.Data.Count).To(Equal(2))
		})

		PIt("sync host multiple module in different biz", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.1.7",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
					moduleId2,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})
	})

	Describe("clone host host", func() {
		It("clone host less biz", func() {
			input := &metadata.CloneHostPropertyParams{
				OrgIP:   "127.0.1.1",
				DstIP:   "127.0.1.2",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("clone host less srcip", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   bizId,
				DstIP:   "127.0.1.2",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("clone host less dstip", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   bizId,
				OrgIP:   "127.0.1.1",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("clone host invalid biz", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   noExistID,
				OrgIP:   "127.0.1.1",
				DstIP:   "127.0.1.2",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("clone host invalid srcip", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   bizId,
				OrgIP:   "127.0.1",
				DstIP:   "127.0.1.2",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("clone host invalid dstip", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   bizId,
				OrgIP:   "127.0.1.1",
				DstIP:   "127.0.1",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("clone host exist dstip", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   bizId,
				OrgIP:   "127.0.1.1",
				DstIP:   "127.0.1.2",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})
	})

	Describe("batch operate host", func() {

		// 清空数据
		BeforeEach(func() {
			clearData()
		})

		// 准备数据
		JustBeforeEach(func() {
			prepareData()
		})

		It("update host less hostid", func() {
			input := map[string]interface{}{
				"bk_sn": "update_bk_sn",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update host invalid hostid", func() {
			input := map[string]interface{}{
				"bk_host_id": "2ew213,fe",
				"bk_sn":      "update_bk_sn",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update host empty hostid", func() {
			input := map[string]interface{}{
				"bk_host_id": "",
				"bk_sn":      "update_bk_sn",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		// the hostId1's bk_sn will be updated successfully, noExistID is ignored
		It("update host one nonexist hostid", func() {
			input := map[string]interface{}{
				"bk_host_id": fmt.Sprintf("%v,%v", hostId1, noExistID),
				"bk_sn":      "update_bk_sn",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update host one nonexist attr", func() {
			input := map[string]interface{}{
				"bk_host_id":      fmt.Sprintf("%v,%v", hostId1, hostId3),
				"bk_sn":           "update_bk_sn",
				"fecfecefrrwdxww": "test",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			rsp1, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0", strconv.FormatInt(hostId1, 10), header)
			util.RegisterResponse(rsp1)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp1.Result).To(Equal(true))
			for _, data := range rsp1.Data {
				if data.PropertyID == "bk_sn" {
					Expect(data.PropertyValue).To(Equal("update_bk_sn"))
					break
				}
			}
		})

		It("update host one invalid attr value", func() {
			input := map[string]interface{}{
				"bk_host_id": fmt.Sprintf("%v,%v", hostId1, hostId3),
				"bk_sn":      1,
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update host using created attr", func() {
			input := map[string]interface{}{
				"bk_host_id": fmt.Sprintf("%v", hostId),
				"a":          "2",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update host using delete attr", func() {
			input := map[string]interface{}{
				"bk_host_id": fmt.Sprintf("%v", hostId),
				"a":          "",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		// one host format check is unpass ,lead to all the host fail to delete
		It("delete host one invalid bk_host_id", func() {
			input := map[string]interface{}{
				"bk_host_id": fmt.Sprintf("%v,abc", hostId4),
			}
			rsp, err := hostServerClient.DeleteHostBatch(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))

			rsp1, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0", strconv.FormatInt(hostId4, 10), header)
			util.RegisterResponse(rsp1)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp1.Result).To(Equal(true))
		})

		// delete host batch does not judge if host exists since it has no side effect
		It("delete host one nonexist bk_host_id", func() {
			input := map[string]interface{}{
				"bk_host_id": fmt.Sprintf("%v,%v", hostId4, noExistID),
			}
			rsp, err := hostServerClient.DeleteHostBatch(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			rsp1, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0", strconv.FormatInt(hostId4, 10), header)
			util.RegisterResponse(rsp1)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp1.Result).To(Equal(false))
		})
	})
})

// 初始化数据，为后续的操作做准备
func prepareData() {
	// 在业务bizId中加入主机
	input := map[string]interface{}{
		"bk_biz_id": bizId,
		"host_info": map[string]interface{}{
			"4": map[string]interface{}{
				"bk_host_innerip": "127.0.1.1",
				"bk_cloud_id":     0,
			},
			"5": map[string]interface{}{
				"bk_host_innerip": "127.0.1.2",
				"bk_cloud_id":     0,
			},
		},
	}
	rsp, err := hostServerClient.AddHost(context.Background(), header, input)
	util.RegisterResponse(rsp)
	Expect(err).NotTo(HaveOccurred())
	Expect(rsp.Result).To(Equal(true))

	input1 := &params.HostCommonSearch{
		AppID: int(bizId),
		Page: params.PageInfo{
			Sort: "bk_host_id",
		},
	}
	rsp1, err := hostServerClient.SearchHost(context.Background(), header, input1)
	util.RegisterResponse(rsp1)
	Expect(err).NotTo(HaveOccurred())
	Expect(rsp1.Result).To(Equal(true))
	Expect(rsp1.Data.Count).To(Equal(2))
	hostId, err = commonutil.GetInt64ByInterface(rsp1.Data.Info[0]["host"].(map[string]interface{})["bk_host_id"])
	Expect(err).NotTo(HaveOccurred())
	hostId2, err = commonutil.GetInt64ByInterface(rsp1.Data.Info[1]["host"].(map[string]interface{})["bk_host_id"])
	Expect(err).NotTo(HaveOccurred())

	// 在业务bizId1中加入主机
	input2 := map[string]interface{}{
		"bk_biz_id": bizId1,
		"host_info": map[string]interface{}{
			"4": map[string]interface{}{
				"bk_host_innerip": "127.0.1.3",
				"bk_cloud_id":     0,
			},
			"5": map[string]interface{}{
				"bk_host_innerip": "127.0.1.4",
				"bk_cloud_id":     0,
			},
		},
	}
	rsp2, err := hostServerClient.AddHost(context.Background(), header, input2)
	util.RegisterResponse(rsp2)
	Expect(err).NotTo(HaveOccurred())
	Expect(rsp2.Result).To(Equal(true))

	input3 := &params.HostCommonSearch{
		AppID: int(bizId1),
		Page: params.PageInfo{
			Sort: "bk_host_id",
		},
	}
	rsp3, err := hostServerClient.SearchHost(context.Background(), header, input3)
	util.RegisterResponse(rsp3)
	Expect(err).NotTo(HaveOccurred())
	Expect(rsp3.Result).To(Equal(true))
	Expect(rsp3.Data.Count).To(Equal(2))
	hostId1, err = commonutil.GetInt64ByInterface(rsp3.Data.Info[0]["host"].(map[string]interface{})["bk_host_id"])
	Expect(err).NotTo(HaveOccurred())
	hostId3, err = commonutil.GetInt64ByInterface(rsp3.Data.Info[1]["host"].(map[string]interface{})["bk_host_id"])
	Expect(err).NotTo(HaveOccurred())

	// 在资源池中加入主机
	input4 := map[string]interface{}{
		"host_info": map[string]interface{}{
			"4": map[string]interface{}{
				"bk_host_innerip": "127.0.1.5",
				"bk_cloud_id":     0,
			},
		},
	}
	rsp4, err := hostServerClient.AddHost(context.Background(), header, input4)
	util.RegisterResponse(rsp4)
	Expect(err).NotTo(HaveOccurred())
	Expect(rsp4.Result).To(Equal(true))

	// 查看资源池中的主机数量
	input5 := &params.HostCommonSearch{
		AppID: -1,
		Condition: []params.SearchCondition{
			{
				ObjectID: "biz",
				Condition: []interface{}{
					map[string]interface{}{
						"field":    "default",
						"operator": "$eq",
						"value":    1,
					},
				},
				Fields: []string{},
			},
		},
		Page: params.PageInfo{
			Sort: "bk_host_id",
		},
	}
	rsp5, err := hostServerClient.SearchHost(context.Background(), header, input5)
	util.RegisterResponse(rsp5)
	Expect(err).NotTo(HaveOccurred())
	Expect(rsp5.Result).To(Equal(true))
	data := rsp5.Data.Info[0]["host"].(map[string]interface{})
	hostId4, err = commonutil.GetInt64ByInterface(data["bk_host_id"])
	Expect(err).NotTo(HaveOccurred())
	Expect(rsp5.Data.Count).To(Equal(1))
}

// 清除所有数据，保证测试用例之间互不干扰
func clearData() {
	// 业务包括两个自建的业务和资源池
	bizIds := []int64{bizId, bizId1, -1}
	for _, bizId := range bizIds {
		// 获取业务下的所有主机
		input := &params.HostCommonSearch{
			AppID: int(bizId),
			Page: params.PageInfo{
				Sort: "bk_host_id",
			},
		}
		rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		hostIds := []int64{}
		hostIds2 := []string{}
		for _, hostInfo := range rsp.Data.Info {
			hostIdInt, err := commonutil.GetInt64ByInterface(hostInfo["host"].(map[string]interface{})["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
			hostIds = append(hostIds, hostIdInt)
			hostIds2 = append(hostIds2, commonutil.GetStrByInterface(hostInfo["host"].(map[string]interface{})["bk_host_id"]))
		}

		if len(hostIds) > 0 {
			if bizId != -1 {
				// 将业务下的主机全部转到该业务下的空闲模块
				input1 := &metadata.DefaultModuleHostConfigParams{
					ApplicationID: bizId,
					HostIDs:       hostIds,
				}
				rsp1, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input1)
				util.RegisterResponse(rsp1)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp1.Result).To(Equal(true))

				// 将业务下的空闲模块主机全部转到资源池
				input2 := &metadata.DefaultModuleHostConfigParams{
					ApplicationID: bizId,
					HostIDs:       hostIds,
					ModuleID:      defaultDirID,
				}
				rsp2, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input2)
				util.RegisterResponse(rsp2)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp2.Result).To(Equal(true))
			}
			// 删除资源池里的主机
			input4 := map[string]interface{}{
				"bk_host_id": strings.Join(hostIds2, ","),
			}
			// By(fmt.Sprintf("*********DeleteHostBatch bid:%v, input4:%+v*******", bizId, input4))
			rsp4, err := hostServerClient.DeleteHostBatch(context.Background(), header, input4)
			util.RegisterResponse(rsp4)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp4.Result).To(Equal(true))
		}

		// 查询业务下的主机
		input3 := &params.HostCommonSearch{
			AppID: int(bizId),
			Page: params.PageInfo{
				Sort: "bk_host_id",
			},
		}
		rsp3, err := hostServerClient.SearchHost(context.Background(), header, input3)
		util.RegisterResponse(rsp3)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp3.Result).To(Equal(true))
		// By(fmt.Sprintf("*********bid:%v, data:%+v*******", bizId, rsp3.Data))
		Expect(rsp3.Data.Count).To(Equal(0))
	}
}
