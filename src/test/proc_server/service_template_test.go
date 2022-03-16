package proc_server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	"configcenter/src/common/selector"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("service template test", func() {
	var categoryId, serviceTemplateId, moduleId, serviceId, serviceId1, processTemplateId, processId int64
	resMap := make(map[string]interface{}, 0)

	Describe("service template test", func() {
		It("create service category", func() {
			input := map[string]interface{}{
				"bk_parent_id":      0,
				common.BKAppIDField: bizId,
				"name":              "test10",
			}
			rsp, err := serviceClient.CreateServiceCategory(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp.Data)
			data := metadata.ServiceCategory{}
			json.Unmarshal(j, &data)
			categoryId = data.ID
		})

		It("search service category", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			resMap["service_category"] = j
		})

		It("create service template", func() {
			input := map[string]interface{}{
				"service_category_id": categoryId,
				common.BKAppIDField:   bizId,
				"name":                "st",
			}
			rsp, err := serviceClient.CreateServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp.Data)
			data := metadata.ServiceTemplate{}
			json.Unmarshal(j, &data)
			Expect(data.Name).To(Equal("st"))
			Expect(data.ServiceCategoryID).To(Equal(categoryId))
			serviceTemplateId = data.ID
		})

		It("search service template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_category_id": categoryId,
			}
			rsp, err := serviceClient.SearchServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"count\":1"))
			Expect(j).To(ContainSubstring("\"name\":\"st\""))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"service_category_id\":%d", categoryId)))
			resMap["service_template"] = j
		})

		It("find service template count info", func() {
			input := map[string]interface{}{
				"service_template_ids": []int64{serviceTemplateId},
			}
			rsp, err := serviceClient.FindServiceTemplateCountInfo(context.Background(), header, bizId, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(len(rsp.Data)).To(Equal(1))
			Expect(rsp.Data[0].(map[string]interface{})["service_template_id"].(json.Number).Int64()).To(Equal(serviceTemplateId))
			Expect(rsp.Data[0].(map[string]interface{})["process_template_count"].(json.Number).Int64()).To(Equal(int64(0)))
			Expect(rsp.Data[0].(map[string]interface{})["service_instance_count"].(json.Number).Int64()).To(Equal(int64(0)))
			Expect(rsp.Data[0].(map[string]interface{})["module_count"].(json.Number).Int64()).To(Equal(int64(0)))
		})

		It("create service template with empty name", func() {
			input := map[string]interface{}{
				"service_category_id": categoryId,
				common.BKAppIDField:   bizId,
				"name":                "",
			}
			rsp, err := serviceClient.CreateServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("create service template with same name", func() {
			input := map[string]interface{}{
				"service_category_id": categoryId,
				common.BKAppIDField:   bizId,
				"name":                "st",
			}
			rsp, err := serviceClient.CreateServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("create service template with invalid service category", func() {
			input := map[string]interface{}{
				"service_category_id": 12345,
				common.BKAppIDField:   bizId,
				"name":                "st1",
			}
			rsp, err := serviceClient.CreateServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("search service template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_category_id": categoryId,
			}
			rsp, err := serviceClient.SearchServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(Equal(resMap["service_template"]))
		})

		It("delete service category with template", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"id":                categoryId,
			}
			rsp, err := serviceClient.DeleteServiceCategory(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("search service category", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(Equal(resMap["service_category"]))
		})

		It("create module with template", func() {
			input := map[string]interface{}{
				"bk_module_name":      "test12345",
				"bk_parent_id":        setId,
				"service_category_id": categoryId,
				"service_template_id": serviceTemplateId,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
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
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"bk_module_name\":\"st\""))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"service_template_id\":%d", serviceTemplateId)))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"service_category_id\":%d", categoryId)))
			resMap["module"] = j
		})

		It("create module with invalid template", func() {
			input := map[string]interface{}{
				"bk_module_name":      "12345",
				"bk_parent_id":        setId,
				"service_category_id": categoryId,
				"service_template_id": 10000,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("create module with unmatch category and template", func() {
			input := map[string]interface{}{
				"bk_module_name":      "123",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": serviceTemplateId,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("create module with same template", func() {
			input := map[string]interface{}{
				"bk_module_name":      "test1234567",
				"bk_parent_id":        setId,
				"service_category_id": categoryId,
				"service_template_id": serviceTemplateId,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
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
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(Equal(resMap["module"]))
		})

		It("update module with template", func() {
			input := map[string]interface{}{
				"bk_module_name":      "TEST",
				"service_category_id": 2,
				"service_template_id": 1000,
			}
			rsp, err := instClient.UpdateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), strconv.FormatInt(moduleId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
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
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"bk_module_name\":\"st\""))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"service_template_id\":%d", serviceTemplateId)))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"service_category_id\":%d", categoryId)))
		})

		It("delete service template with module", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_template_id": serviceTemplateId,
			}
			rsp, err := serviceClient.DeleteServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("search service template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_category_id": categoryId,
			}
			rsp, err := serviceClient.SearchServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(Equal(resMap["service_template"]))
		})
	})

	Describe("service instance test", func() {
		It("create service instance with service template that has no process template", func() {
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

			input := &metadata.CreateServiceInstanceInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Instances: []metadata.CreateServiceInstanceDetail{
					{
						HostID: hostId1,
						Processes: []metadata.ProcessInstanceDetail{
							{
								ProcessData: map[string]interface{}{
									"bk_func_name":    "p1",
									"bk_process_name": "p1",
								},
							},
						},
					},
				},
			}
			serviceIds, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			util.RegisterResponse(serviceIds)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("process template test", func() {
		It("create process template", func() {
			input := map[string]interface{}{
				"service_template_id": serviceTemplateId,
				common.BKAppIDField:   bizId,
				"processes": []map[string]interface{}{
					{
						"spec": map[string]interface{}{
							"bk_func_name": map[string]interface{}{
								"value":            "p1",
								"as_default_value": true,
							},
							"bk_process_name": map[string]interface{}{
								"value":            "p1",
								"as_default_value": true,
							},
							"bk_start_param_regex": map[string]interface{}{
								"value":            "123",
								"as_default_value": false,
							},
							common.BKProcPortEnable: map[string]interface{}{
								"value":            false,
								"as_default_value": false,
							},
						},
					},
				},
			}
			rsp, err := processClient.CreateProcessTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			processTemplateId, err = commonutil.GetInt64ByInterface(rsp.Data.([]interface{})[0])
			Expect(err).NotTo(HaveOccurred())
		})

		It("search process template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_template_id": serviceTemplateId,
			}
			rsp, err := processClient.SearchProcessTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"bk_func_name\":{\"as_default_value\":true,\"value\":\"p1\"}"))
			Expect(j).To(ContainSubstring("\"bk_process_name\":{\"as_default_value\":true,\"value\":\"p1\"}"))
			Expect(j).To(ContainSubstring("\"bk_start_param_regex\":{\"as_default_value\":false,\"value\":\"123\"}"))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"id\":%d", processTemplateId)))
			resMap["process_template"] = j
		})

		It("create process template with same name", func() {
			input := map[string]interface{}{
				"service_template_id": serviceTemplateId,
				common.BKAppIDField:   bizId,
				"processes": []map[string]interface{}{
					{
						"spec": map[string]interface{}{
							"bk_func_name": map[string]interface{}{
								"value":            "p123",
								"as_default_value": true,
							},
							"bk_process_name": map[string]interface{}{
								"value":            "p1",
								"as_default_value": true,
							},
						},
					},
				},
			}
			rsp, err := processClient.CreateProcessTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("create process template with same bk_func_name and bk_start_param_regex", func() {
			input := map[string]interface{}{
				"service_template_id": serviceTemplateId,
				common.BKAppIDField:   bizId,
				"processes": []map[string]interface{}{
					{
						"spec": map[string]interface{}{
							"bk_func_name": map[string]interface{}{
								"value":            "p1",
								"as_default_value": true,
							},
							"bk_process_name": map[string]interface{}{
								"value":            "p123",
								"as_default_value": true,
							},
							"bk_start_param_regex": map[string]interface{}{
								"value":            "123",
								"as_default_value": false,
							},
						},
					},
				},
			}
			rsp, err := processClient.CreateProcessTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("create process template with empty name", func() {
			input := map[string]interface{}{
				"service_template_id": serviceTemplateId,
				common.BKAppIDField:   bizId,
				"processes": []map[string]interface{}{
					{
						"spec": map[string]interface{}{
							"bk_func_name": map[string]interface{}{
								"value":            "",
								"as_default_value": true,
							},
							"bk_process_name": map[string]interface{}{
								"value":            "",
								"as_default_value": true,
							},
						},
					},
				},
			}
			rsp, err := processClient.CreateProcessTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("create process template with invalid service template", func() {
			input := map[string]interface{}{
				"service_template_id": 10000,
				common.BKAppIDField:   bizId,
				"processes": []map[string]interface{}{
					{
						"spec": map[string]interface{}{
							"bk_func_name": map[string]interface{}{
								"value":            "123",
								"as_default_value": true,
							},
							"bk_process_name": map[string]interface{}{
								"value":            "123",
								"as_default_value": true,
							},
						},
					},
				},
			}
			rsp, err := processClient.CreateProcessTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("search process template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_template_id": serviceTemplateId,
			}
			rsp, err := processClient.SearchProcessTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(Equal(resMap["process_template"]))
		})

		It("create service instance with template", func() {
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
				"is_increment":        true,
				"disable_auto_create": true,
			}
			rsp, rawErr := hostServerClient.TransferHostModule(context.Background(), header, transInput)
			util.RegisterResponse(rsp)
			Expect(rawErr).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			input := &metadata.CreateServiceInstanceInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Instances: []metadata.CreateServiceInstanceDetail{
					{
						HostID: hostId1,
						Processes: []metadata.ProcessInstanceDetail{
							{
								ProcessData: map[string]interface{}{
									"bk_func_name":    "p1",
									"bk_process_name": "p1",
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
		})

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			exists := false
			for _, svcInst := range data.Info {
				if svcInst.ID == serviceId {
					exists = true
					break
				}
			}
			Expect(exists).To(Equal(true))
			resMap["service_instance"] = data
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
									"bk_func_name":    "p1",
									"bk_process_name": "p1",
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

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data).To(Equal(resMap["service_instance"]))
		})

		It("search service instance by set template id", func() {
			input := map[string]interface{}{
				"set_template_id": 1,
			}
			rsp, err := serviceClient.SearchServiceInstanceBySetTemplate(context.Background(), strconv.FormatInt(bizId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("clone service instance to another host", func() {
			input := &metadata.CreateServiceInstanceInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Instances: []metadata.CreateServiceInstanceDetail{
					{
						HostID: hostId2,
						Processes: []metadata.ProcessInstanceDetail{
							{
								ProcessData: map[string]interface{}{
									"bk_func_name":    "p2",
									"bk_process_name": "p3",
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
			serviceId1 = serviceIds[0]
		})

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			svcExists, svc2Exists := false, false
			for _, svcInst := range data.Info {
				if svcInst.ID == serviceId {
					svcExists = true
				}
				if svcInst.ID == serviceId1 {
					svc2Exists = true
				}
			}
			Expect(svcExists && svc2Exists).To(Equal(true))
			resMap["service_instance"] = data
		})

		It("search process instance", func() {
			input := &metadata.ListProcessInstancesOption{
				BizID:             bizId,
				ServiceInstanceID: serviceId1,
			}
			data, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(data)).To(Equal(1))
			Expect(data[0].Property["bk_process_name"]).To(Equal("p1"))
			Expect(data[0].Property["bk_func_name"]).To(Equal("p1"))
			Expect(data[0].Property["bk_start_param_regex"]).To(Equal("123"))
			Expect(data[0].Relation.HostID).To(Equal(hostId2))
			processId, err = commonutil.GetInt64ByInterface(data[0].Property["bk_process_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("create process instance", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_instance_id": serviceId,
				"processes": []map[string]interface{}{
					{
						"process_info": map[string]interface{}{
							"bk_process_name": "p2",
							"bk_func_name":    "p2",
						},
					},
				},
			}
			rsp, err := processClient.CreateProcessInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("update process instance", func() {
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
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search process instance", func() {
			input := &metadata.ListProcessInstancesOption{
				BizID:             bizId,
				ServiceInstanceID: serviceId1,
			}
			data, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(data)).To(Equal(1))
			Expect(data[0].Property["bk_process_name"]).To(Equal("p1"))
			Expect(data[0].Property["bk_func_name"]).To(Equal("p1"))
			Expect(data[0].Property["bk_start_param_regex"]).To(Equal("1234"))
			Expect(data[0].Relation.HostID).To(Equal(hostId2))
		})

		It("update process template", func() {
			input := map[string]interface{}{
				"process_template_id": processTemplateId,
				common.BKAppIDField:   bizId,
				"process_property": map[string]interface{}{
					"bk_func_name": map[string]interface{}{
						"value":            "123",
						"as_default_value": false,
					},
					"bk_process_name": map[string]interface{}{
						"value":            "123",
						"as_default_value": false,
					},
					"bk_start_param_regex": map[string]interface{}{
						"value":            "123456",
						"as_default_value": true,
					},
				},
			}
			rsp, err := processClient.UpdateProcessTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"bk_func_name\":{\"as_default_value\":true,\"value\":\"p1\"}"))
			Expect(j).To(ContainSubstring("\"bk_process_name\":{\"as_default_value\":true,\"value\":\"123\"}"))
			Expect(j).To(ContainSubstring("\"bk_start_param_regex\":{\"as_default_value\":true,\"value\":\"123456\"}"))
		})

		It("search process template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_template_id": serviceTemplateId,
			}
			rsp, err := processClient.SearchProcessTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"bk_func_name\":{\"as_default_value\":true,\"value\":\"p1\"}"))
			Expect(j).To(ContainSubstring("\"bk_process_name\":{\"as_default_value\":true,\"value\":\"123\"}"))
			Expect(j).To(ContainSubstring("\"bk_start_param_regex\":{\"as_default_value\":true,\"value\":\"123456\"}"))
			resMap["process_template"] = j
		})
	})

	Describe("update template test", func() {
		It("update service template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"id":                  serviceTemplateId,
				"service_category_id": 2,
				"name":                "abcdefg",
			}
			rsp, err := serviceClient.UpdateServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp.Data)
			data := metadata.ServiceTemplate{}
			json.Unmarshal(j, &data)
			Expect(data.Name).To(Equal("abcdefg"))
			Expect(data.ServiceCategoryID).To(Equal(int64(2)))
		})

		It("search service template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_category_id": 2,
			}
			rsp, err := serviceClient.SearchServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"count\":1"))
			Expect(j).To(ContainSubstring("\"name\":\"abcdefg\""))
			Expect(j).To(ContainSubstring("\"service_category_id\":2"))
			resMap["service_template"] = j
		})

		It("create service template with name 'service_template'", func() {
			input := map[string]interface{}{
				"service_category_id": categoryId,
				common.BKAppIDField:   bizId,
				"name":                "service_template",
			}
			rsp, err := serviceClient.CreateServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("update service template with same name as another service template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"id":                  serviceTemplateId,
				"service_category_id": categoryId,
				"name":                "service_template",
			}
			rsp, err := serviceClient.UpdateServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("create module with name 'service_template_module' in the same set", func() {
			input := map[string]interface{}{
				"bk_module_name":      "service_template_module",
				"bk_parent_id":        setId,
				"service_category_id": categoryId,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("update service template with same name as another module in the same set with a module using this template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"id":                  serviceTemplateId,
				"service_category_id": categoryId,
				"name":                "service_template_module",
			}
			rsp, err := serviceClient.UpdateServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("update service template with invalid service category", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"id":                  serviceTemplateId,
				"service_category_id": 100000,
				"name":                "abcdefg",
			}
			rsp, err := serviceClient.UpdateServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("search service template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_category_id": 2,
			}
			rsp, err := serviceClient.SearchServiceTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(Equal(resMap["service_template"]))
		})

		It("remove service instance with template with process", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"service_instance_ids": []int64{
					serviceId,
				},
			}
			rsp, err := serviceClient.DeleteServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("sync service instance and template after add and change process template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"bk_module_ids":       []int64{moduleId},
				"service_template_id": serviceTemplateId,
			}
			rsp, err := serviceClient.SyncServiceInstanceByTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())

			// wait for some time to get the async task done
			time.Sleep(time.Second * 20)
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
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"bk_module_name\":\"abcdefg\""))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"service_template_id\":%d", serviceTemplateId)))
			Expect(j).To(ContainSubstring("\"service_category_id\":2"))
		})

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				HostIDs:  []int64{hostId1},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			serviceId = data.Info[0].ID
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
			Expect(data[0].Property["bk_process_name"]).To(Equal("123"))
			Expect(data[0].Property["bk_func_name"]).To(Equal("p1"))
			Expect(data[0].Property["bk_start_param_regex"]).To(Equal("123456"))
			Expect(data[0].Relation.HostID).To(Equal(hostId1))
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
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("delete process instance", func() {
			input := &metadata.DeleteProcessInstanceInServiceInstanceInput{
				BizID:              bizId,
				ProcessInstanceIDs: []int64{processId},
			}
			err := processClient.DeleteProcessInstance(context.Background(), header, input)
			Expect(err).To(HaveOccurred())
		})

		It("search process instance", func() {
			input := &metadata.ListProcessInstancesOption{
				BizID:             bizId,
				ServiceInstanceID: serviceId1,
			}
			data, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(data)).To(Equal(1))
			Expect(data[0].Property["bk_process_name"]).To(Equal("123"))
			Expect(data[0].Property["bk_func_name"]).To(Equal("p1"))
			Expect(data[0].Property["bk_start_param_regex"]).To(Equal("123456"))
		})
	})

	Describe("service instance label test", func() {
		It("service instance add labels", func() {
			input := map[string]interface{}{
				"labels": map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
				"instance_ids": []int64{
					serviceId,
					serviceId1,
				},
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.ServiceInstanceAddLabels(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("service instance add and edit labels", func() {
			input := map[string]interface{}{
				"labels": map[string]interface{}{
					"key2": "value",
					"key3": "value3",
				},
				"instance_ids": []int64{
					serviceId1,
				},
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.ServiceInstanceAddLabels(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			resMap["service_instance"] = data
			Expect(data.Count).To(Equal(uint64(2)))
			Expect(len(data.Info[0].Labels)).To(Equal(3))
			Expect(data.Info[0].Labels["key1"]).To(Equal("value1"))
			Expect(data.Info[0].Labels["key2"]).To(Equal("value"))
			Expect(data.Info[0].Labels["key3"]).To(Equal("value3"))
			Expect(len(data.Info[1].Labels)).To(Equal(2))
			Expect(data.Info[1].Labels["key1"]).To(Equal("value1"))
			Expect(data.Info[1].Labels["key2"]).To(Equal("value2"))
		})

		It("service instance add labels with empty key values", func() {
			input := map[string]interface{}{
				"labels": map[string]interface{}{
					"":     "value1",
					"key1": "",
					"key4": "value4",
				},
				"instance_ids": []int64{
					serviceId,
					serviceId1,
				},
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.ServiceInstanceAddLabels(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("service instance add labels with invalid instance id", func() {
			input := map[string]interface{}{
				"labels": map[string]interface{}{
					"key5": "value5",
				},
				"instance_ids": []int64{
					serviceId,
					10000,
					serviceId1,
				},
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.ServiceInstanceAddLabels(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("search module service instances labels", func() {
			input := map[string]interface{}{
				"bk_module_id":      moduleId,
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.ServiceInstanceFindLabels(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp.Data)
			data := make(map[string][]string)
			json.Unmarshal(j, &data)
			Expect(len(data)).To(Equal(3))
			Expect(data["key1"]).To(ConsistOf("value1"))
			Expect(data["key2"]).To(ConsistOf("value2", "value"))
			Expect(data["key3"]).To(ConsistOf("value3"))
		})

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(resMap["service_instance"]).To(Equal(data))
		})

		It("search service instance without key", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					selector.Selector{
						Key:      "key3",
						Operator: "!",
					}},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data.Count).To(Equal(uint64(1)))
			Expect(data.Info[0].ID).To(Equal(serviceId))
		})

		It("search service instance without key with values", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					selector.Selector{
						Key:      "key3",
						Operator: "!",
						Values:   []string{"123"},
					}},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).To(HaveOccurred())
		})

		It("search service instance exists key", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					selector.Selector{
						Key:      "key3",
						Operator: "exists",
					}},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data.Count).To(Equal(uint64(1)))
			Expect(data.Info[0].ID).To(Equal(serviceId1))
		})

		It("search service instance exists key with values", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					selector.Selector{
						Key:      "key3",
						Operator: "exists",
						Values:   []string{"123"},
					}},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).To(HaveOccurred())
		})

		It("search service instance with equal key value", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					selector.Selector{
						Key:      "key1",
						Operator: "=",
						Values:   []string{"value1"},
					},
					selector.Selector{
						Key:      "key2",
						Operator: "=",
						Values:   []string{"value2"},
					},
				},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data.Count).To(Equal(uint64(1)))
			Expect(data.Info[0].ID).To(Equal(serviceId))
		})

		It("search service instance with equal key zero value", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					selector.Selector{
						Key:      "key1",
						Operator: "=",
					}},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).To(HaveOccurred())
		})

		It("search service instance with equal key many values", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					selector.Selector{
						Key:      "key1",
						Operator: "=",
						Values:   []string{"value1", "value2"},
					}},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).To(HaveOccurred())
		})

		It("search service instance with not equal key value", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					{
						Key:      "key1",
						Operator: "!=",
						Values: []string{
							"value2",
						},
					},
					{
						Key:      "key2",
						Operator: "!=",
						Values: []string{
							"value",
						},
					},
				},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data.Count).To(Equal(uint64(1)))
			Expect(data.Info[0].ID).To(Equal(serviceId))
		})

		It("search service instance with not equal key zero value", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					{
						Key:      "key1",
						Operator: "!=",
					},
				},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).To(HaveOccurred())
		})

		It("search service instance with not equal key many values", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					{
						Key:      "key1",
						Operator: "!=",
						Values: []string{
							"value1",
							"value2",
						},
					},
				},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).To(HaveOccurred())
		})

		It("search service instance with value in values", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					{
						Key:      "key1",
						Operator: "in",
						Values: []string{
							"value1",
						},
					},
					{
						Key:      "key2",
						Operator: "in",
						Values: []string{
							"value",
							"value2",
						},
					},
				},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data.Count).To(Equal(uint64(2)))
			Expect(data.Info[0].ID).To(Equal(serviceId1))
			Expect(data.Info[1].ID).To(Equal(serviceId))
		})

		It("search service instance with value in zero values", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					{
						Key:      "key1",
						Operator: "in",
					},
				},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).To(HaveOccurred())
		})

		It("search service instance with value not in values", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					{
						Key:      "key3",
						Operator: "notin",
						Values: []string{
							"value",
						},
					},
					{
						Key:      "key1",
						Operator: "notin",
						Values: []string{
							"value",
							"value2",
						},
					},
				},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data.Count).To(Equal(uint64(2)))
			Expect(data.Info[0].ID).To(Equal(serviceId1))
			Expect(data.Info[1].ID).To(Equal(serviceId))
		})

		It("search service instance with value not in zero values", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					{
						Key:      "key1",
						Operator: "notin",
					},
				},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).To(HaveOccurred())
		})

		It("search service instance with invalid operator", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					selector.Selector{
						Key:      "key1",
						Operator: "123",
					}},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).To(HaveOccurred())
		})

		It("search service instance with empty key", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					{
						Key:      "",
						Operator: "exists",
					},
				},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).To(HaveOccurred())
		})

		It("search service instance with no matching data", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
				Selectors: selector.Selectors{
					{
						Key:      "key1",
						Operator: "!",
					},
					{
						Key:      "key3",
						Operator: "exists",
					},
					{
						Key:      "key3",
						Operator: "notin",
						Values: []string{
							"value",
						},
					},
					{
						Key:      "key2",
						Operator: "!=",
						Values: []string{
							"value3",
						},
					},
					{
						Key:      "key2",
						Operator: "=",
						Values: []string{
							"value2",
						},
					},
				},
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(int(data.Count)).To(Equal(0))
		})

		It("service instance remove labels", func() {
			input := map[string]interface{}{
				"keys": []string{
					"key1",
					"",
					"key3",
				},
				"instance_ids": []int64{
					serviceId,
					serviceId1,
				},
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.ServiceInstanceRemoveLabels(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("service instance remove labels with invalid service instance id", func() {
			input := map[string]interface{}{
				"keys": []string{
					"key2",
				},
				"instance_ids": []int64{
					serviceId,
					100000,
					serviceId1,
				},
				common.BKAppIDField: bizId,
			}
			rsp, err := serviceClient.ServiceInstanceRemoveLabels(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false), rsp.ToString())
		})

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data.Count).To(Equal(uint64(2)))
			Expect(len(data.Info[0].Labels)).To(Equal(1))
			Expect(data.Info[0].Labels["key2"]).To(Equal("value"))
			Expect(len(data.Info[1].Labels)).To(Equal(1))
			Expect(data.Info[1].Labels["key2"]).To(Equal("value2"))
		})
	})

	Describe("removal test", func() {
		It("remove process template", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"process_templates": []int64{
					processTemplateId,
				},
			}
			rsp, err := processClient.DeleteProcessTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search process template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"service_template_id": serviceTemplateId,
			}
			rsp, err := processClient.SearchProcessTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"count\":0"))
		})

		It("remove service instance with template with process", func() {
			input := map[string]interface{}{
				common.BKAppIDField: bizId,
				"service_instance_ids": []int64{
					serviceId1,
				},
			}
			rsp, err := serviceClient.DeleteServiceInstance(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data.Count).To(Equal(uint64(1)))
			Expect(data.Info[0].ID).To(Equal(serviceId))
		})

		It("sync service instance and template after remove process template", func() {
			input := map[string]interface{}{
				common.BKAppIDField:   bizId,
				"bk_module_ids":       []int64{moduleId},
				"service_template_id": serviceTemplateId,
			}
			rsp, err := serviceClient.SyncServiceInstanceByTemplate(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())

			// wait for some time to get the async task done
			time.Sleep(time.Second * 20)
		})

		It("search process instance", func() {
			input := &metadata.ListProcessInstancesOption{
				BizID:             bizId,
				ServiceInstanceID: serviceId,
			}
			data, err := processClient.SearchProcessInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(data)).To(Equal(0))
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
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"service_category_id\":2"))
		})

		It("search service instance", func() {
			input := &metadata.GetServiceInstanceInModuleInput{
				BizID:    bizId,
				ModuleID: moduleId,
			}
			data, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			util.RegisterResponse(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data.Count).To(Equal(uint64(0)))
		})

		// unbind service template on module is prohibited
		PDescribe("unbind service template on module", func() {
			It("unbind service template on module", func() {
				input := map[string]interface{}{
					"metadata": map[string]interface{}{
						"label": map[string]interface{}{
							"bk_biz_id": strconv.FormatInt(bizId, 10),
						},
					},
					"bk_module_id": moduleId,
				}
				rsp, err := serviceClient.RemoveTemplateBindingOnModule(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.ToString())
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
				Expect(rsp.Result).To(Equal(true), rsp.ToString())
				j, err := json.Marshal(rsp)
				Expect(j).To(ContainSubstring("\"service_template_id\":0"))
			})

			It("delete service template", func() {
				input := map[string]interface{}{
					"metadata": map[string]interface{}{
						"label": map[string]interface{}{
							"bk_biz_id": strconv.FormatInt(bizId, 10),
						},
					},
					"service_template_id": serviceTemplateId,
				}
				rsp, err := serviceClient.DeleteServiceTemplate(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.ToString())
			})

			It("search service template", func() {
				input := map[string]interface{}{
					"metadata": map[string]interface{}{
						"label": map[string]interface{}{
							"bk_biz_id": strconv.FormatInt(bizId, 10),
						},
					},
					"service_category_id": 2,
				}
				rsp, err := serviceClient.SearchServiceTemplate(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.ToString())
				j, err := json.Marshal(rsp)
				Expect(j).To(ContainSubstring("\"count\":0"))
			})

			It("delete service category", func() {
				input := map[string]interface{}{
					"metadata": map[string]interface{}{
						"label": map[string]interface{}{
							"bk_biz_id": strconv.FormatInt(bizId, 10),
						},
					},
					"id": categoryId,
				}
				rsp, err := serviceClient.DeleteServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.ToString())
			})

			It("search service category", func() {
				input := map[string]interface{}{
					"metadata": map[string]interface{}{
						"label": map[string]interface{}{
							"bk_biz_id": strconv.FormatInt(bizId, 10),
						},
					},
				}
				rsp, err := serviceClient.SearchServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true), rsp.ToString())
				j, err := json.Marshal(rsp)
				Expect(j).NotTo(ContainSubstring(fmt.Sprintf("\"id\":%d", categoryId)))
			})

			It("delete service category twice", func() {
				input := map[string]interface{}{
					"metadata": map[string]interface{}{
						"label": map[string]interface{}{
							"bk_biz_id": strconv.FormatInt(bizId, 10),
						},
					},
					"id": categoryId,
				}
				rsp, err := serviceClient.DeleteServiceCategory(context.Background(), header, input)
				util.RegisterResponse(rsp)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false), rsp.ToString())
			})
		})
	})
})
