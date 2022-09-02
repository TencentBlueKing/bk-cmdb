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
package proc_server_test

import (
	"context"
	"encoding/json"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("no service template test", func() {
	var categoryId1, categoryId2, categoryId3, moduleId, serviceId, serviceId2, serviceId3, processId int64
	resMap := make(map[string]interface{}, 0)

	Describe("service category test", func() {
		Describe("create service category test", func() {
			It("create service category", func() {
				input := map[string]interface{}{
					"bk_parent_id":      0,
					common.BKAppIDField: bizId,
					"name":              "test",
				}
				rsp, err := serviceClient.CreateServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				j, err := json.Marshal(rsp.Data)
				data := metadata.ServiceCategory{}
				json.Unmarshal(j, &data)
				Expect(data.Name).To(Equal("test"))
				Expect(data.ParentID).To(Equal(int64(0)))
				Expect(data.RootID).To(Equal(data.ID))
				categoryId1 = data.ID
			})

			It("create service category with parent", func() {
				input := map[string]interface{}{
					"bk_parent_id":      categoryId1,
					common.BKAppIDField: bizId,
					"name":              "test1",
				}
				rsp, err := serviceClient.CreateServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				j, err := json.Marshal(rsp.Data)
				data := metadata.ServiceCategory{}
				json.Unmarshal(j, &data)
				Expect(data.Name).To(Equal("test1"))
				Expect(data.ParentID).To(Equal(int64(categoryId1)))
				Expect(data.RootID).To(Equal(int64(categoryId1)))
				categoryId2 = data.ID
			})

			It("create service category with grandparent", func() {
				input := map[string]interface{}{
					"bk_parent_id":      categoryId2,
					common.BKAppIDField: bizId,
					"name":              "test2",
				}
				rsp, err := serviceClient.CreateServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				j, err := json.Marshal(rsp.Data)
				data := metadata.ServiceCategory{}
				json.Unmarshal(j, &data)
				Expect(data.Name).To(Equal("test2"))
				Expect(data.ParentID).To(Equal(int64(categoryId2)))
				Expect(data.RootID).To(Equal(int64(categoryId1)))
				categoryId3 = data.ID
			})

			It("search service category", func() {
				input := &metadata.ListServiceCategoryOption{
					BizID: bizId,
				}
				rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
				util.RegisterResponseWithRid(rsp, header)
				Expect(err).NotTo(HaveOccurred())
				j, _ := json.Marshal(rsp)
				Expect(j).To(ContainSubstring("\"name\":\"test\""))
				Expect(j).To(ContainSubstring("\"name\":\"test1\""))
				Expect(j).To(ContainSubstring("\"name\":\"test2\""))
				resMap["service_category"] = j
			})

			It("create service category with invalid parent", func() {
				input := map[string]interface{}{
					"bk_parent_id":      10000,
					common.BKAppIDField: bizId,
					"name":              "test3",
				}
				rsp, err := serviceClient.CreateServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
			})

			It("create service category with empty name", func() {
				input := map[string]interface{}{
					"bk_parent_id":      0,
					common.BKAppIDField: bizId,
					"name":              "",
				}
				rsp, err := serviceClient.CreateServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
			})

			It("create service category with duplicate name", func() {
				input := map[string]interface{}{
					"bk_parent_id":      0,
					common.BKAppIDField: bizId,
					"name":              "test",
				}
				rsp, err := serviceClient.CreateServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
			})

			It("search service category", func() {
				input := &metadata.ListServiceCategoryOption{
					BizID: bizId,
				}
				rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
				util.RegisterResponseWithRid(rsp, header)
				Expect(err).NotTo(HaveOccurred())
				j, _ := json.Marshal(rsp)
				Expect(j).To(Equal(resMap["service_category"]))
			})
		})

		Describe("modify service category test", func() {
			It("update service category", func() {
				input := map[string]interface{}{
					"name": "test3",
					"id":   categoryId3,
				}
				rsp, err := serviceClient.UpdateServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				j, err := json.Marshal(rsp.Data)
				data := metadata.ServiceCategory{}
				json.Unmarshal(j, &data)
				Expect(data.Name).To(Equal("test3"))
				Expect(data.ParentID).To(Equal(int64(categoryId2)))
				Expect(data.RootID).To(Equal(int64(categoryId1)))
			})

			It("search service category", func() {
				input := &metadata.ListServiceCategoryOption{
					BizID: bizId,
				}
				rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
				util.RegisterResponseWithRid(rsp, header)
				Expect(err).NotTo(HaveOccurred())
				j, _ := json.Marshal(rsp)
				Expect(j).NotTo(ContainSubstring("\"name\":\"test2\""))
				Expect(j).To(ContainSubstring("\"name\":\"test3\""))
				resMap["service_category"] = j
			})

			It("update service category with empty name", func() {
				input := map[string]interface{}{
					"name": "",
					"id":   categoryId3,
				}
				rsp, err := serviceClient.UpdateServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
			})

			It("create service category with same parent", func() {
				input := map[string]interface{}{
					"bk_parent_id":      categoryId2,
					common.BKAppIDField: bizId,
					"name":              "test4",
				}
				rsp, err := serviceClient.CreateServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				j, err := json.Marshal(rsp.Data)
				data := metadata.ServiceCategory{}
				json.Unmarshal(j, &data)
				Expect(data.Name).To(Equal("test4"))
				Expect(data.ParentID).To(Equal(int64(categoryId2)))
				Expect(data.RootID).To(Equal(int64(categoryId1)))
			})

			It("search service category", func() {
				input := &metadata.ListServiceCategoryOption{
					BizID: bizId,
				}
				rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
				util.RegisterResponseWithRid(rsp, header)
				Expect(err).NotTo(HaveOccurred())
				j, _ := json.Marshal(rsp)
				Expect(j).To(ContainSubstring("\"name\":\"test4\""))
				resMap["service_category"] = j
			})

			It("update service category with duplicate name", func() {
				input := map[string]interface{}{
					"name": "test4",
					"id":   categoryId3,
				}
				rsp, err := serviceClient.UpdateServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
			})

			It("delete service category with children", func() {
				input := map[string]interface{}{
					common.BKAppIDField: bizId,
					"id":                categoryId1,
				}
				rsp, err := serviceClient.DeleteServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
			})

			It("create module without template using service category", func() {
				input := map[string]interface{}{
					"bk_module_name":      "test",
					"bk_parent_id":        setId,
					"service_category_id": categoryId3,
					"service_template_id": 0,
				}
				rsp, e := instClient.CreateModule(context.Background(), bizId, setId, header, input)
				util.RegisterResponseWithRid(rsp, header)
				Expect(e).NotTo(HaveOccurred())
				var err error
				Expect(rsp["bk_module_name"].(string)).To(Equal("test"))
				setIdRes, err := commonutil.GetInt64ByInterface(rsp["bk_set_id"])
				Expect(err).NotTo(HaveOccurred())
				Expect(setIdRes).To(Equal(setId))
				parentIdRes, err := commonutil.GetInt64ByInterface(rsp["bk_parent_id"])
				Expect(err).NotTo(HaveOccurred())
				Expect(parentIdRes).To(Equal(setId))
				moduleId, err = commonutil.GetInt64ByInterface(rsp["bk_module_id"])
				Expect(err).NotTo(HaveOccurred())
			})

			It("search module", func() {
				input := &params.SearchParams{
					Condition: map[string]interface{}{},
					Page: map[string]interface{}{
						"sort": "id",
					},
				}
				rsp, err := instClient.SearchModule(context.Background(), "0", bizId, setId, header, input)
				util.RegisterResponseWithRid(rsp, header)
				Expect(err).NotTo(HaveOccurred())
				j, _ := json.Marshal(rsp)
				Expect(j).To(ContainSubstring("\"bk_module_name\":\"test\""))
				Expect(j).To(ContainSubstring(fmt.Sprintf("\"service_category_id\":%d", categoryId3)))
			})

			It("create module with invalid service_category_id", func() {
				input := map[string]interface{}{
					"bk_module_name":      "module1",
					"bk_parent_id":        setId,
					"service_category_id": 12345,
					"service_template_id": 0,
				}
				rsp, err := instClient.CreateModule(context.Background(), bizId, setId, header, input)
				util.RegisterResponse(rsp)
				Expect(err).To(HaveOccurred())
			})

			It("search module", func() {
				input := &params.SearchParams{
					Condition: map[string]interface{}{},
					Page: map[string]interface{}{
						"sort": "id",
					},
				}
				rsp, err := instClient.SearchModule(context.Background(), "0", bizId, setId, header, input)
				util.RegisterResponseWithRid(rsp, header)
				Expect(err).NotTo(HaveOccurred())
				j, _ := json.Marshal(rsp)
				Expect(j).NotTo(ContainSubstring("module1"))
			})

			It("delete service category with module", func() {
				input := map[string]interface{}{
					common.BKAppIDField: bizId,
					"id":                categoryId3,
				}
				rsp, err := serviceClient.DeleteServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
			})

			It("search service category", func() {
				input := &metadata.ListServiceCategoryOption{
					BizID: bizId,
				}
				rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
				util.RegisterResponseWithRid(rsp, header)
				Expect(err).NotTo(HaveOccurred())
				j, _ := json.Marshal(rsp)
				Expect(j).To(Equal(resMap["service_category"]))
			})
		})
	})

	Describe("create service instance test", func() {
		It("create service instance with host not in the module", func() {
			svcInput := &metadata.CreateServiceInstanceInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Instances: []metadata.CreateServiceInstanceDetail{
					{
						HostID: hostId1,
					},
				},
			}
			serviceIds, err := serviceClient.CreateServiceInstance(context.Background(), header, svcInput)
			util.RegisterResponse(serviceIds)
			Expect(err).To(HaveOccurred())

			By(fmt.Sprintf("transfer host %d & %d to the module %d", hostId1, hostId2, moduleId))
			transInput := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId1,
					hostId2,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": true,
			}
			rsp, rawErr := hostServerClient.TransferHostModule(context.Background(), header, transInput)
			util.RegisterResponse(rsp)
			Expect(rawErr).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("create service instance without template with processes", func() {
			input := &metadata.CreateServiceInstanceInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Instances: []metadata.CreateServiceInstanceDetail{
					{
						HostID: hostId1,
						Processes: []metadata.ProcessInstanceDetail{
							{
								ProcessData: map[string]interface{}{
									"bk_func_name":         "p1",
									"bk_process_name":      "p1",
									"bk_start_param_regex": "",
								},
							},
						},
					},
				},
			}
			serviceIds, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(serviceIds)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(serviceIds)).NotTo(Equal(0))
			serviceId = serviceIds[0]
			Expect(err).NotTo(HaveOccurred())
		})

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())

			for _, svcInst := range data.Info {
				Expect(svcInst.ID).To(Equal(serviceId))
				Expect(svcInst.HostID).To(Equal(hostId1))
				Expect(svcInst.ModuleID).To(Equal(moduleId))
			}
			resMap["service_instance"] = data
		})

		It("create service instance with invalid module", func() {
			input := &metadata.CreateServiceInstanceInput{
				BizID:    bizId,
				ModuleID: 10000,
				Instances: []metadata.CreateServiceInstanceDetail{
					{
						HostID: hostId1,
					},
				},
			}
			_, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			Expect(err).To(HaveOccurred())
		})

		It("create service instance with invalid host", func() {
			input := &metadata.CreateServiceInstanceInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Instances: []metadata.CreateServiceInstanceDetail{
					{
						HostID: 10000,
					},
				},
			}
			rsp, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).To(HaveOccurred())
		})

		// TODO: ADD TRANSACTION TO FIX THIS
		It("create service instance with invalid process", func() {
			input := &metadata.CreateServiceInstanceInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Instances: []metadata.CreateServiceInstanceDetail{
					{
						HostID: hostId1,
						Processes: []metadata.ProcessInstanceDetail{
							{
								ProcessData: map[string]interface{}{
									"bk_func_name":         "",
									"bk_process_name":      "",
									"bk_start_param_regex": "",
								},
							},
						},
					},
				},
			}
			rsp, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).To(HaveOccurred())
		})

		PIt("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data).To(Equal(resMap["service_instance"]))
		})

		It("clone service instance to source host", func() {
			input := &metadata.CreateServiceInstanceInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Instances: []metadata.CreateServiceInstanceDetail{
					{
						HostID: hostId1,
						Processes: []metadata.ProcessInstanceDetail{
							{
								ProcessData: map[string]interface{}{
									"bk_func_name":         "p1",
									"bk_process_name":      "p1",
									"bk_start_param_regex": "",
								},
							},
						},
					},
				},
			}
			serviceIds, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(serviceIds)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(serviceIds)).NotTo(Equal(0))
			serviceId2 = serviceIds[0]
		})

		It("clone service instance to other host", func() {
			input := &metadata.CreateServiceInstanceInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Instances: []metadata.CreateServiceInstanceDetail{
					{
						HostID: hostId2,
						Processes: []metadata.ProcessInstanceDetail{
							{
								ProcessData: map[string]interface{}{
									"bk_func_name":         "p1",
									"bk_process_name":      "p1",
									"bk_start_param_regex": "",
								},
							},
						},
					},
				},
			}
			serviceIds, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(serviceIds)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(serviceIds)).NotTo(Equal(0))
			serviceId3 = serviceIds[0]
		})

		It("create service instance without template with no process", func() {
			input := &metadata.CreateServiceInstanceInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Instances: []metadata.CreateServiceInstanceDetail{
					{
						HostID: hostId1,
					},
				},
			}
			serviceIds, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(serviceIds)
			Expect(err).To(HaveOccurred())
		})

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(data.Info)).To(Equal(3))
		})

		It("update service instance", func() {
			input := map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"service_instance_id": serviceId3,
						"update": map[string]interface{}{
							"name": "inst_update_test",
						},
					},
				},
			}
			rsp, err := serviceClient.UpdateServiceInstances(context.Background(), header, bizId, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
		})

		It("delete service instance", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"service_instance_ids": []int64{
					serviceId3,
				},
			}
			rsp, err := serviceClient.DeleteServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
		})

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			for _, svcInst := range data.Info {
				Expect(svcInst.ID).NotTo(Equal(serviceId3))
			}
		})
	})

	Describe("process instance test", func() {
		It("create process instance", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_instance_id": serviceId,
				"processes": []map[string]interface{}{
					{
						"process_info": map[string]interface{}{
							"bk_func_name":         "p2",
							"bk_process_name":      "p2",
							"bk_start_param_regex": "123",
						},
					},
				},
			}
			rsp, err := processClient.CreateProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			processId, err = commonutil.GetInt64ByInterface(rsp.Data.([]interface{})[0])
			Expect(err).NotTo(HaveOccurred())
		})

		It("search process instance", func() {
			input := &metadata.ListProcessInstancesOption{
				BizID:             bizId,
				ServiceInstanceID: serviceId,
			}
			data, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(data)).To(Equal(2))
			Expect(data[0].Property["bk_process_name"]).To(Equal("p1"))
			Expect(data[0].Property["bk_func_name"]).To(Equal("p1"))
			Expect(data[0].Relation.HostID).To(Equal(hostId1))
			Expect(data[1].Property["bk_process_name"]).To(Equal("p2"))
			Expect(data[1].Property["bk_func_name"]).To(Equal("p2"))
			Expect(data[1].Property["bk_start_param_regex"]).To(Equal("123"))
			Expect(data[1].Relation.HostID).To(Equal(hostId1))
			resMap["process_instance"] = data
		})

		It("create process instance with same name", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_instance_id": serviceId,
				"processes": []map[string]interface{}{
					{
						"process_info": map[string]interface{}{
							"bk_func_name":         "p",
							"bk_process_name":      "p2",
							"bk_start_param_regex": "123",
						},
					},
				},
			}
			rsp, err := processClient.CreateProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
		})

		It("create process instance with same bk_func_name and bk_start_param_regex", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_instance_id": serviceId,
				"processes": []map[string]interface{}{
					{
						"process_info": map[string]interface{}{
							"bk_func_name":         "p2",
							"bk_process_name":      "p1234",
							"bk_start_param_regex": "123",
						},
					},
				},
			}
			rsp, err := processClient.CreateProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
		})

		It("create process instance with empty name", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_instance_id": serviceId,
				"processes": []map[string]interface{}{
					{
						"process_info": map[string]interface{}{
							"bk_process_name": "",
						},
					},
				},
			}
			rsp, err := processClient.CreateProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
		})

		It("search process instance", func() {
			input := &metadata.ListProcessInstancesOption{
				BizID:             bizId,
				ServiceInstanceID: serviceId,
			}
			data, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(resMap["process_instance"]).To(Equal(data))
		})

		It("udpate process instance", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"processes": []map[string]interface{}{
					{
						"bk_func_name":         "p3",
						"bk_process_name":      "p3",
						"bk_start_param_regex": "1234",
						"bk_process_id":        processId,
						"bind_info": []map[string]interface{}{
							{
								"ip":       "127.0.0.1",
								"port":     "1024",
								"protocol": "1",
								"enable":   true,
							},
						},
					},
				},
			}
			rsp, err := processClient.UpdateProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
		})

		It("search process instance", func() {
			input := &metadata.ListProcessInstancesOption{
				BizID:             bizId,
				ServiceInstanceID: serviceId,
			}
			data, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())

			Expect(len(data)).To(Equal(2))
			Expect(data[1].Property["bk_process_name"]).To(Equal("p3"))
			Expect(data[1].Property["bk_func_name"]).To(Equal("p3"))
			Expect(data[1].Property["bk_start_param_regex"]).To(Equal("1234"))
			Expect(data[1].Relation.HostID).To(Equal(hostId1))
			resMap["process_instance"] = data
		})

		It("update process instance with same name", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"processes": []map[string]interface{}{
					{
						"bk_process_name": "p1",
						"bk_process_id":   processId,
					},
				},
			}
			rsp, err := processClient.UpdateProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
		})

		It("update process instance with same bk_func_name and bk_start_param_regex", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"processes": []map[string]interface{}{
					{
						"bk_func_name":         "p1",
						"bk_process_name":      "p1234",
						"bk_start_param_regex": "",
						"bk_process_id":        processId,
					},
				},
			}
			rsp, err := processClient.UpdateProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
		})

		It("udpate process instance with empty name", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"processes": []map[string]interface{}{
					{
						"bk_process_name": "",
						"bk_process_id":   processId,
					},
				},
			}
			rsp, err := processClient.UpdateProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
		})

		It("search process instance", func() {
			input := &metadata.ListProcessInstancesOption{
				BizID:             bizId,
				ServiceInstanceID: serviceId,
			}
			data, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(resMap["process_instance"]).To(Equal(data))
		})

		It("list process related info", func() {
			input := metadata.ListProcessRelatedInfoOption{
				Set:    metadata.SetCondOfP{},
				Module: metadata.ModuleCondOfP{},
				ServiceInstance: metadata.ServiceInstanceCondOfP{
					IDs: []int64{serviceId},
				},
				Fields: []string{},
				Page: metadata.BasePage{
					Start: 0,
					Limit: 100,
				},
			}
			rsp, err := processClient.ListProcessRelatedInfo(context.Background(), header, bizId, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			Expect(rsp.Data.Count).To(Not(Equal(0)))
		})

		It("list process instance names with their ids in one module", func() {
			input := map[string]interface{}{
				"bk_module_id":      moduleId,
				common.BKAppIDField: bizId,
				"page": map[string]interface{}{
					"start": 0,
					"limit": 10,
					"sort":  "bk_process_name",
				},
			}
			rsp, err := processClient.ListProcessInstancesNameIDsInModule(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			data := struct {
				Count int64                             `json:"count"`
				Info  []metadata.ProcessInstanceNameIDs `json:"info"`
			}{}
			j, err := json.Marshal(rsp.Data)
			json.Unmarshal(j, &data)
			Expect(data.Count).To(Equal(int64(2)))
			Expect(data.Info[0].ProcessName).To(Equal("p1"))
		})

		It("list process instance details by their ids", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"process_ids":       []int64{processId},
				"page": map[string]interface{}{
					"start": 0,
					"limit": 10,
					"sort":  "bk_process_id",
				},
			}
			rsp, err := processClient.ListProcessInstancesDetailsByIDs(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			data := struct {
				Count int64                                `json:"count"`
				Info  []metadata.ProcessInstanceDetailByID `json:"info"`
			}{}
			j, err := json.Marshal(rsp.Data)
			json.Unmarshal(j, &data)
			Expect(data.Count).To(Equal(int64(1)))
			Expect(data.Info[0].Property[common.BKProcessNameField]).To(Equal("p3"))
		})

		It("list process instance details", func() {
			input := metadata.ListProcessInstancesDetailsOption{
				ProcessIDs: []int64{processId},
				Fields:     []string{common.BKProcessIDField, common.BKProcessNameField},
			}
			rsp, err := processClient.ListProcessInstancesDetails(context.Background(), header, bizId, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			Expect(len(rsp.Data)).To(Equal(1))
			pName, _ := rsp.Data[0].String(common.BKProcessNameField)
			Expect(pName).To(Equal("p3"))
		})

		It("update process instances by their ids", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"process_ids":       []int64{processId},
				"update_data": map[string]interface{}{
					common.BKDescriptionField: "aaa",
				},
			}
			rsp, err := processClient.UpdateProcessInstancesByIDs(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
		})

		It("list process instance details with bind info", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"process_ids":       []int64{processId},
				"page": map[string]interface{}{
					"start": 0,
					"limit": 10,
					"sort":  "bk_process_id",
				},
			}
			rsp, err := processClient.ListProcessInstancesDetailsByIDs(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			data := struct {
				Count int64                                `json:"count"`
				Info  []metadata.ProcessInstanceDetailByID `json:"info"`
			}{}
			j, err := json.Marshal(rsp.Data)
			json.Unmarshal(j, &data)
			Expect(data.Count).To(Equal(int64(1)))
			Expect(data.Info[0].Property[common.BKProcessNameField]).To(Equal("p3"))
			Expect(data.Info[0].Property[common.BKDescriptionField]).To(Equal("aaa"))
			bindInfo := map[string]interface{}{
				"ip":       "127.0.0.1",
				"port":     "1024",
				"protocol": "1",
				"enable":   true,
			}
			ExpectBindInfoArr, err := commonutil.GetMapInterfaceByInterface(
				data.Info[0].Property[common.BKProcBindInfo])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(ExpectBindInfoArr)).To(Equal(int(1)))
			expectBindInfo, ok := ExpectBindInfoArr[0].(map[string]interface{})
			delete(expectBindInfo, "template_row_id")
			Expect(ok).To(Equal(true))
			Expect(expectBindInfo).To(Equal(bindInfo))

		})

		It("delete process instance", func() {
			input := &metadata.DeleteProcessInstanceInServiceInstanceInput{
				BizID:              bizId,
				ProcessInstanceIDs: []int64{processId},
			}
			err := processClient.DeleteProcessInstance(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
		})

		It("search process instance", func() {
			input := &metadata.ListProcessInstancesOption{
				BizID:             bizId,
				ServiceInstanceID: serviceId,
			}
			data, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(data)).To(Equal(1))
		})

		It("delete service instance with process", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"service_instance_ids": []int64{
					serviceId,
				},
			}
			rsp, err := serviceClient.DeleteServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
		})

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())

			for _, svcInst := range data.Info {
				Expect(svcInst.ID).NotTo(Equal(serviceId))
			}
		})

		It("delete all processes in service instances", func() {
			By("search all processes in service instances")
			searchInput := &metadata.ListProcessInstancesOption{
				BizID:             bizId,
				ServiceInstanceID: serviceId2,
			}
			procData, err := processClient.SearchProcessInstance(context.Background(), header, searchInput)
			util.RegisterResponse(procData)
			Expect(err).NotTo(HaveOccurred())

			processIDs := make([]int64, len(procData))
			for index, proc := range procData {
				procID, err := commonutil.GetInt64ByInterface(proc.Property[common.BKProcessIDField])
				Expect(err).NotTo(HaveOccurred())
				processIDs[index] = procID
			}

			By("delete all processes in service instances")
			deleteInput := &metadata.DeleteProcessInstanceInServiceInstanceInput{
				BizID:              bizId,
				ProcessInstanceIDs: processIDs,
			}
			err = processClient.DeleteProcessInstance(context.Background(), header, deleteInput)
			Expect(err).NotTo(HaveOccurred())

			By("check if service instances has been deleted")
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			svcInstData, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(svcInstData)
			Expect(err).NotTo(HaveOccurred())
			Expect(int(svcInstData.Count)).To(Equal(0))
		})
	})
})
