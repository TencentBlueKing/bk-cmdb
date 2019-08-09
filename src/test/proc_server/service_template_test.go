package proc_server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("service template test", func() {
	var categoryId, serviceTemplateId, moduleId, serviceId, serviceId1, processTemplateId int64
	resMap := make(map[string]interface{}, 0)

	Describe("service template test", func() {
		It("create service category", func() {
			input := map[string]interface{}{
				"bk_parent_id": 0,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"name": "test10",
			}
			rsp, err := serviceClient.CreateServiceCategory(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := metadata.ServiceCategory{}
			json.Unmarshal(j, &data)
			categoryId = data.ID
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			resMap["service_category"] = j
		})

		It("create service template", func() {
			input := map[string]interface{}{
				"service_category_id": categoryId,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"name": "st",
			}
			rsp, err := serviceClient.CreateServiceTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := metadata.ServiceTemplate{}
			json.Unmarshal(j, &data)
			Expect(data.Name).To(Equal("st"))
			Expect(data.ServiceCategoryID).To(Equal(categoryId))
			serviceTemplateId = data.ID
		})

		It("search service template", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"service_category_id": categoryId,
			}
			rsp, err := serviceClient.SearchServiceTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"count\":\"1\""))
			Expect(j).To(ContainSubstring("\"name\":\"st\""))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"service_category_id\":%d", categoryId)))
			resMap["service_template"] = j
		})

		It("create service template with empty name", func() {
			input := map[string]interface{}{
				"service_category_id": categoryId,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"name": "",
			}
			rsp, err := serviceClient.CreateServiceTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1199006))
		})

		It("create service template with invalid service category", func() {
			input := map[string]interface{}{
				"service_category_id": 12345,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"name": "st1",
			}
			rsp, err := serviceClient.CreateServiceTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1199006))
		})

		It("search service template", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"service_category_id": categoryId,
			}
			rsp, err := serviceClient.SearchServiceTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			Expect(j).To(Equal(resMap["service_template"]))
		})

		It("delete service category with template", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"id": categoryId,
			}
			rsp, err := serviceClient.DeleteServiceCategory(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1199056))
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			moduleId = int64(rsp.Data["bk_module_id"].(float64))
		})

		It("search module", func() {
			input := &params.SearchParams{
				Condition: map[string]interface{}{},
				Page: map[string]interface{}{
					"sort": "id",
				},
			}
			rsp, err := instClient.SearchModule(context.Background(), "0", strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"bk_module_name\":\"test12345\""))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"service_template_id\":%d", serviceTemplateId)))
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1199006))
		})

		It("create module with unmatch category and template", func() {
			input := map[string]interface{}{
				"bk_module_name":      "123",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": serviceTemplateId,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1108036))
		})

		It("create module with same template", func() {
			input := map[string]interface{}{
				"bk_module_name":      "test1234567",
				"bk_parent_id":        setId,
				"service_category_id": categoryId,
				"service_template_id": serviceTemplateId,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1199014))
		})

		It("search module", func() {
			input := &params.SearchParams{
				Condition: map[string]interface{}{},
				Page: map[string]interface{}{
					"sort": "id",
				},
			}
			rsp, err := instClient.SearchModule(context.Background(), "0", strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			Expect(j).To(Equal(resMap["module"]))
		})

		It("delete service template with module", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"service_template_id": serviceTemplateId,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1199056))
		})

		It("search service template", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"service_category_id": categoryId,
			}
			rsp, err := serviceClient.SearchServiceTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			Expect(j).To(Equal(resMap["service_template"]))
		})
	})

	Describe("service instance test", func() {
		It("create service instance with template", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"bk_module_id": moduleId,
				"instances": []map[string]interface{}{
					{
						"bk_host_id": hostId1,
						"processes": []map[string]interface{}{
							{
								"process_info": map[string]interface{}{
									"bk_func_name":    "p1",
									"bk_process_name": "p1",
								},
							},
						},
					},
				},
			}
			rsp, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			serviceId = int64(rsp.Data.([]interface{})[0].(float64))
		})

		It("search service instance", func() {
			input := map[string]interface{}{
				"bk_module_id": moduleId,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
			}
			rsp, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"id\":%d", serviceId)))
			resMap["service_instance"] = j
		})

		It("clone service instance to source host", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"bk_module_id": moduleId,
				"instances": []map[string]interface{}{
					{
						"bk_host_id": hostId1,
						"processes": []map[string]interface{}{
							{
								"process_info": map[string]interface{}{
									"bk_func_name":    "p1",
									"bk_process_name": "p1",
								},
							},
						},
					},
				},
			}
			rsp, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1113016))
		})

		It("search service instance", func() {
			input := map[string]interface{}{
				"bk_module_id": moduleId,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
			}
			rsp, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			Expect(j).To(Equal(resMap["service_instance"]))
		})
	})

	Describe("process template test", func() {
		It("create process template", func() {
			input := map[string]interface{}{
				"service_template_id": serviceTemplateId,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
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
						},
					},
				},
			}
			rsp, err := processClient.CreateProcessTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			processTemplateId = int64(rsp.Data.([]interface{})[0].(float64))
		})

		It("search process template", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"service_template_id": serviceTemplateId,
			}
			rsp, err := processClient.SearchProcessTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"bk_func_name\":{\"value\":\"p1\",\"as_default_value\":\"true\"}"))
			Expect(j).To(ContainSubstring("\"bk_process_name\":{\"value\":\"p1\",\"as_default_value\":\"true\"}"))
			Expect(j).To(ContainSubstring("\"bk_start_param_regex\":{\"value\":\"123\",\"as_default_value\":\"false\"}"))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"id\":%d", processTemplateId)))
			resMap["process_template"] = j
		})

		It("create process template with same name", func() {
			input := map[string]interface{}{
				"service_template_id": serviceTemplateId,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1113019))
		})

		It("create process template with same bk_func_name and bk_start_param_regex", func() {
			input := map[string]interface{}{
				"service_template_id": serviceTemplateId,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1113020))
		})

		It("create process template with empty name", func() {
			input := map[string]interface{}{
				"service_template_id": serviceTemplateId,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1199006))
		})

		It("create process template with invalid service template", func() {
			input := map[string]interface{}{
				"service_template_id": 10000,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1199006))
		})

		It("search process template", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"service_template_id": serviceTemplateId,
			}
			rsp, err := processClient.SearchProcessTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"bk_func_name\":{\"value\":\"p1\",\"as_default_value\":\"true\"}"))
			Expect(j).To(ContainSubstring("\"bk_process_name\":{\"value\":\"p1\",\"as_default_value\":\"true\"}"))
			Expect(j).To(ContainSubstring("\"bk_start_param_regex\":{\"value\":\"123\",\"as_default_value\":\"false\"}"))
			Expect(j).To(Equal(resMap["process_template"]))
		})

		It("clone service instance to other host", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"bk_module_id": moduleId,
				"instances": []map[string]interface{}{
					{
						"bk_host_id": hostId2,
						"processes": []map[string]interface{}{
							{
								"process_info": map[string]interface{}{
									"bk_func_name":    "p1",
									"bk_process_name": "p1",
								},
							},
						},
					},
				},
			}
			rsp, err := serviceClient.CreateServiceInstance(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			serviceId1 = int64(rsp.Data.([]interface{})[0].(float64))
		})

		It("search service instance", func() {
			input := map[string]interface{}{
				"bk_module_id": moduleId,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
			}
			rsp, err := serviceClient.SearchServiceInstance(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"id\":%d", serviceId)))
			Expect(j).To(ContainSubstring(fmt.Sprintf("\"id\":%d", serviceId1)))
			resMap["service_instance"] = j
		})

		It("search process instance", func() {
			input := map[string]interface{}{
				"service_instance_id": serviceId1,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
			}
			rsp, err := processClient.SearchProcessInstance(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := []metadata.ProcessInstance{}
			json.Unmarshal(j, &data)
			Expect(len(data)).To(Equal(2))
			Expect(data[1].Property["bk_process_name"]).To(Equal("p1"))
			Expect(data[1].Property["bk_func_name"]).To(Equal("p1"))
			Expect(data[1].Property["bk_start_param_regex"]).To(Equal("123"))
			Expect(data[1].Relation.HostID).To(Equal(hostId2))
			resMap["process_instance"] = j
		})

		It("create process instance", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(1108035))
		})

		It("udpate process instance", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search process instance", func() {
			input := map[string]interface{}{
				"service_instance_id": serviceId1,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
			}
			rsp, err := processClient.SearchProcessInstance(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			Expect(j).To(Equal(resMap["process_instance"]))
		})

		It("update process template", func() {
			input := map[string]interface{}{
				"process_template_id": processTemplateId,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"processes": []map[string]interface{}{
					{
						"spec": map[string]interface{}{
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
					},
				},
			}
			rsp, err := processClient.UpdateProcessTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(j).To(ContainSubstring("\"bk_func_name\":{\"value\":\"123\",\"as_default_value\":\"true\"}"))
			Expect(j).To(ContainSubstring("\"bk_process_name\":{\"value\":\"123\",\"as_default_value\":\"true\"}"))
			Expect(j).To(ContainSubstring("\"bk_start_param_regex\":{\"value\":\"123456\",\"as_default_value\":\"true\"}"))
		})

		It("search process template", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"service_template_id": serviceTemplateId,
			}
			rsp, err := processClient.SearchProcessTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"bk_func_name\":{\"value\":\"123\",\"as_default_value\":\"true\"}"))
			Expect(j).To(ContainSubstring("\"bk_process_name\":{\"value\":\"123\",\"as_default_value\":\"true\"}"))
			Expect(j).To(ContainSubstring("\"bk_start_param_regex\":{\"value\":\"123456\",\"as_default_value\":\"true\"}"))
			resMap["process_template"] = j
		})
	})

	Describe("update template test", func() {
		It("update service template", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"id":                  serviceTemplateId,
				"service_category_id": 2,
				"name":                "abcdefg",
			}
			rsp, err := serviceClient.UpdateServiceTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp.Data)
			data := metadata.ServiceTemplate{}
			json.Unmarshal(j, &data)
			Expect(data.Name).To(Equal("st"))
			Expect(data.ServiceCategoryID).To(Equal(2))
		})

		It("search service template", func() {
			input := map[string]interface{}{
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizId, 10),
					},
				},
				"service_category_id": categoryId,
			}
			rsp, err := serviceClient.SearchServiceTemplate(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			j, err := json.Marshal(rsp)
			Expect(j).To(ContainSubstring("\"count\":\"1\""))
			Expect(j).To(ContainSubstring("\"name\":\"st\""))
			Expect(j).To(ContainSubstring("\"service_category_id\":2"))
		})
	})
})
