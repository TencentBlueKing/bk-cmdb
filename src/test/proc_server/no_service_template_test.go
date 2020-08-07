package proc_server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("no service template test", func() {
	var categoryId1, categoryId2, categoryId3, moduleId, serviceId, serviceId1, serviceId2, serviceId3, processId int64
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
				input := map[string]interface{}{
					common.BKAppIDField: bizId,
				}
				rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				j, err := json.Marshal(rsp)
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
				input := map[string]interface{}{
					common.BKAppIDField: bizId,
				}
				rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				j, err := json.Marshal(rsp)
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
				input := map[string]interface{}{
					common.BKAppIDField: bizId,
				}
				rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				j, err := json.Marshal(rsp)
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
				input := map[string]interface{}{
					common.BKAppIDField: bizId,
				}
				rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				j, err := json.Marshal(rsp)
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
				rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				Expect(rsp.Data["bk_module_name"].(string)).To(Equal("test"))
				setIdRes, err := commonutil.GetInt64ByInterface(rsp.Data["bk_set_id"])
				Expect(err).NotTo(HaveOccurred())
				Expect(setIdRes).To(Equal(setId))
				parentIdRes, err := commonutil.GetInt64ByInterface(rsp.Data["bk_parent_id"])
				Expect(err).NotTo(HaveOccurred())
				Expect(parentIdRes).To(Equal(setId))
				moduleId, err = commonutil.GetInt64ByInterface(rsp.Data["bk_module_id"])
				Expect(err).NotTo(HaveOccurred())
			})

			It("search module", func() {
				input := &params.SearchParams{
					Condition: map[string]interface{}{},
					Page: map[string]interface{}{
						"sort": "id",
					},
				}
				rsp, err := instClient.SearchModule(context.Background(), "0", strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				j, err := json.Marshal(rsp)
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
				rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
			})

			It("search module", func() {
				input := &params.SearchParams{
					Condition: map[string]interface{}{},
					Page: map[string]interface{}{
						"sort": "id",
					},
				}
				rsp, err := instClient.SearchModule(context.Background(), "0", strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				j, err := json.Marshal(rsp)
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
				input := map[string]interface{}{
					common.BKAppIDField: bizId,
				}
				rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
				j, err := json.Marshal(rsp)
				Expect(j).To(Equal(resMap["service_category"]))
			})
		})
	})

	Describe("create service instance test", func() {
		It("create service instance without template with processes", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"bk_module_id":      moduleId,
				"instances": []map[string]interface{}{
					{
						"bk_host_id": hostId1,
						"processes": []map[string]interface{}{
							{
								"process_info": map[string]interface{}{
									"bk_func_name":         "p1",
									"bk_process_name":      "p1",
									"bk_start_param_regex": "",
								},
							},
						},
					},
				},
			}
			rsp, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			serviceId, err = commonutil.GetInt64ByInterface(rsp.Data.([]interface{})[0])
			Expect(err).NotTo(HaveOccurred())
		})

		It("search service instance", func() {
			input := map[string]interface{}{
				"bk_module_id":      moduleId,
				common.BKAppIDField: bizId,
				"with_name":         true,
			}
			rsp, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"id\":%d", serviceId)))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"bk_host_id\":%d", hostId1)))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"bk_module_id\":%d", moduleId)))
			resMap["service_instance"] = j
		})

		It("create service instance with invalid module", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"bk_module_id":      10000,
				"instances": []map[string]interface{}{
					{
						"bk_host_id": hostId1,
					},
				},
			}
			rsp, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
		})

		It("create service instance with invalid host", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"bk_module_id":      moduleId,
				"instances": []map[string]interface{}{
					{
						"bk_host_id": 10000,
					},
				},
			}
			rsp, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
		})

		// TODO: ADD TRANSACTION TO FIX THIS
		It("create service instance with invalid process", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"bk_module_id":      moduleId,
				"instances": []map[string]interface{}{
					{
						"bk_host_id": hostId1,
						"processes": []map[string]interface{}{
							{
								"process_info": map[string]interface{}{
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.BaseResp.ToString())
		})

		PIt("search service instance", func() {
			input := map[string]interface{}{
				"bk_module_id":      moduleId,
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(Equal(resMap["service_instance"]))
		})

		It("clone service instance to source host", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"bk_module_id":      moduleId,
				"instances": []map[string]interface{}{
					{
						"bk_host_id": hostId1,
						"processes": []map[string]interface{}{
							{
								"process_info": map[string]interface{}{
									"bk_func_name":         "p1",
									"bk_process_name":      "p1",
									"bk_start_param_regex": "",
								},
							},
						},
					},
				},
			}
			rsp, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			serviceId2, err = commonutil.GetInt64ByInterface(rsp.Data.([]interface{})[0])
			Expect(err).NotTo(HaveOccurred())
		})

		It("clone service instance to other host", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"bk_module_id":      moduleId,
				"instances": []map[string]interface{}{
					{
						"bk_host_id": hostId2,
						"processes": []map[string]interface{}{
							{
								"process_info": map[string]interface{}{
									"bk_func_name":         "p1",
									"bk_process_name":      "p1",
									"bk_start_param_regex": "",
								},
							},
						},
					},
				},
			}
			rsp, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			serviceId3, err = commonutil.GetInt64ByInterface(rsp.Data.([]interface{})[0])
			Expect(err).NotTo(HaveOccurred())
		})

		It("create service instance without template with no process", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"bk_module_id":      moduleId,
				"instances": []map[string]interface{}{
					{
						"bk_host_id": hostId1,
					},
				},
			}
			rsp, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			serviceId1, err = commonutil.GetInt64ByInterface(rsp.Data.([]interface{})[0])
			Expect(err).NotTo(HaveOccurred())
		})

		It("search service instance", func() {
			input := map[string]interface{}{
				"bk_module_id":      moduleId,
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"id\":%d", serviceId1)))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"id\":%d", serviceId2)))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"id\":%d", serviceId3)))
		})

		It("delete service instance with no process", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"service_instance_ids": []int64{
					serviceId1,
				},
			}
			rsp, err := serviceClient.DeleteServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
		})

		It("search service instance", func() {
			input := map[string]interface{}{
				"bk_module_id":      moduleId,
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).NotTo(ContainSubstring(fmt.Sprintf("\"id\":%d", serviceId1)))
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
			input := map[string]interface{}{
				"service_instance_id": serviceId,
				common.BKAppIDField:   bizId,
			}
			rsp, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			j, err := json.Marshal(rsp.Data)
			data := []metadata.ProcessInstance{}
			json.Unmarshal(j, &data)
			Expect(len(data)).To(Equal(2))
			Expect(data[0].Property["bk_process_name"]).To(Equal("p1"))
			Expect(data[0].Property["bk_func_name"]).To(Equal("p1"))
			Expect(data[0].Relation.HostID).To(Equal(hostId1))
			Expect(data[1].Property["bk_process_name"]).To(Equal("p2"))
			Expect(data[1].Property["bk_func_name"]).To(Equal("p2"))
			Expect(data[1].Property["bk_start_param_regex"]).To(Equal("123"))
			Expect(data[1].Relation.HostID).To(Equal(hostId1))
			resMap["process_instance"] = j
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
			input := map[string]interface{}{
				"service_instance_id": serviceId,
				common.BKAppIDField:   bizId,
			}
			rsp, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			j, err := json.Marshal(rsp.Data)
			Expect(resMap["process_instance"]).To(Equal(j))
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
					},
				},
			}
			rsp, err := processClient.UpdateProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
		})

		It("search process instance", func() {
			input := map[string]interface{}{
				"service_instance_id": serviceId,
				common.BKAppIDField:   bizId,
			}
			rsp, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			j, err := json.Marshal(rsp.Data)
			data := []metadata.ProcessInstance{}
			json.Unmarshal(j, &data)
			Expect(len(data)).To(Equal(2))
			Expect(data[1].Property["bk_process_name"]).To(Equal("p3"))
			Expect(data[1].Property["bk_func_name"]).To(Equal("p3"))
			Expect(data[1].Property["bk_start_param_regex"]).To(Equal("1234"))
			Expect(data[1].Relation.HostID).To(Equal(hostId1))
			resMap["process_instance"] = j
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
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_instance_id": serviceId,
			}
			rsp, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			j, err := json.Marshal(rsp.Data)
			Expect(resMap["process_instance"]).To(Equal(j))
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

		It("update process instances by their ids", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"process_ids":       []int64{processId},
				"update_data": map[string]interface{}{
					common.BKProcPortEnable:   true,
					common.BKDescriptionField: "aaa",
					common.BKProtocol:         "1",
				},
			}
			rsp, err := processClient.UpdateProcessInstancesByIDs(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
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
			Expect(data.Info[0].Property[common.BKProcPortEnable]).To(Equal(true))
			Expect(data.Info[0].Property[common.BKDescriptionField]).To(Equal("aaa"))
			Expect(data.Info[0].Property[common.BKProtocol]).To(Equal("1"))
		})

		It("delete process instance", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"process_instance_ids": []int64{
					processId,
				},
			}
			rsp, err := processClient.DeleteProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
		})

		It("search process instance", func() {
			input := map[string]interface{}{
				"service_instance_id": serviceId,
				common.BKAppIDField:   bizId,
			}
			rsp, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			j, err := json.Marshal(rsp.Data)
			data := []metadata.ProcessInstance{}
			json.Unmarshal(j, &data)
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
			input := map[string]interface{}{
				"bk_module_id":      moduleId,
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).NotTo(ContainSubstring(fmt.Sprintf("\"id\":%d", serviceId)))
		})
	})
})

var _ = Describe("list_biz_host_process test", func() {
	It("list_biz_host_process", func() {
		By("add host")
		input := map[string]interface{}{
			"bk_biz_id": bizId,
			"host_info": map[string]interface{}{
				"1": map[string]interface{}{
					"bk_host_innerip": "127.0.0.3",
					"bk_cloud_id":     0,
				},
				"2": map[string]interface{}{
					"bk_host_innerip": "127.0.0.4",
					"bk_cloud_id":     0,
				},
			},
		}
		addHostRsp, err := hostServerClient.AddHost(context.Background(), header, input)
		util.RegisterResponse(addHostRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(addHostRsp.Result).To(Equal(true))

		By("search host")
		input1 := &params.HostCommonSearch{
			AppID: int(bizId),
			Ip: params.IPInfo{
				Data:  []string{"127.0.0.3", "127.0.0.4"},
				Exact: 1,
				Flag:  "bk_host_innerip|bk_host_outerip",
			},
		}
		hostRsp, err := hostServerClient.SearchHost(context.Background(), header, input1)
		util.RegisterResponse(hostRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(true))
		hostId3, err := commonutil.GetInt64ByInterface(hostRsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_id"])
		Expect(err).NotTo(HaveOccurred())
		hostId4, err := commonutil.GetInt64ByInterface(hostRsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_id"])
		Expect(err).NotTo(HaveOccurred())

		By("create module")
		input = map[string]interface{}{
			"bk_module_name":      "list_biz_host_process_test",
			"bk_parent_id":        setId,
			"service_category_id": 2,
			"service_template_id": 0,
		}
		rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true), rsp.BaseResp.ToString())
		moduleId, err := commonutil.GetInt64ByInterface(rsp.Data["bk_module_id"])
		Expect(err).NotTo(HaveOccurred())

		By("create service instance without template with processes")
		input = map[string]interface{}{
			common.BKAppIDField: bizId,
			"bk_module_id":      moduleId,
			"instances": []map[string]interface{}{
				{
					"bk_host_id": hostId3,
					"processes": []map[string]interface{}{
						{
							"process_info": map[string]interface{}{
								common.BKPort:             "123",
								common.BKFuncName:         "p11111",
								common.BKProcessNameField: "p11111",
							},
						},
						{
							"process_info": map[string]interface{}{
								common.BKProtocol:         "2",
								common.BKFuncName:         "p22222",
								common.BKProcessNameField: "p22222",
							},
						},
					},
				},
				{
					"bk_host_id": hostId4,
					"processes": []map[string]interface{}{
						{
							"process_info": map[string]interface{}{
								common.BKBindIP:           "0.0.0.0",
								common.BKFuncName:         "p33333",
								common.BKProcessNameField: "p33333",
							},
						},
					},
				},
			},
		}
		rsp1, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
		util.RegisterResponse(rsp1)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp1.Result).To(Equal(true), rsp1.BaseResp.ToString())

		By("list all biz host process")
		input = map[string]interface{}{
			"bk_host_ids":       []int64{hostId3, hostId4},
			common.BKAppIDField: bizId,
			"page": metadata.BasePage{
				Sort:  "-bk_host_id",
				Limit: 10,
				Start: 0,
			},
		}
		rsp2, err := processClient.ListProcessInstancesWithHost(context.Background(), header, input)
		util.RegisterResponse(rsp2)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp2.Result).To(Equal(true), rsp2.BaseResp.ToString())
		j, err := json.Marshal(rsp2.Data)
		data := struct {
			Count int64                          `json:"count"`
			Info  []metadata.HostProcessInstance `json:"info"`
		}{}
		json.Unmarshal(j, &data)
		Expect(data.Count).To(Equal(int64(3)))
		Expect(data.Info[0].HostID).To(Or(Equal(hostId3), Equal(hostId4)))
		Expect(data.Info[1].HostID).To(Or(Equal(hostId3), Equal(hostId4)))
		Expect(data.Info[2].HostID).To(Or(Equal(hostId3), Equal(hostId4)))
		Expect("0.0.0.0").To(Or(Equal(data.Info[0].BindIP), Equal(data.Info[1].BindIP), Equal(data.Info[2].BindIP)))
		Expect("123").To(Or(Equal(data.Info[0].Port), Equal(data.Info[1].Port), Equal(data.Info[2].Port)))
		Expect("2").To(Or(Equal(string(data.Info[0].Protocol)), Equal(string(data.Info[1].Protocol)), Equal(string(data.Info[2].Protocol))))
	})
})
