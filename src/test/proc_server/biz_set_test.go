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

package proc_server_test

import (
	"context"
	"encoding/json"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("biz set test", func() {

	var bizID1, bizID2, setID1, setID2, moduleID1, moduleID2, moduleID3, sampleBizSetID1, hostID1, hostID2, hostID3,
		serviceId, processTemplateId, serviceTemplateId, processTemplateID int64

	It("prepare environment, create a biz set and biz in it with topo for searching biz and topo in biz set", func() {
		ctx := context.Background()
		biz := map[string]interface{}{
			common.BKMaintainersField: "biz_set",
			common.BKAppNameField:     "biz_for_biz_set",
			"time_zone":               "Africa/Accra",
			"language":                "1",
		}
		bizResp, err := apiServerClient.CreateBiz(ctx, header, biz)
		util.RegisterResponseWithRid(bizResp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(bizResp.Result).To(Equal(true))
		bizID1, err = commonutil.GetInt64ByInterface(bizResp.Data[common.BKAppIDField])
		Expect(err).NotTo(HaveOccurred())

		biz[common.BKAppNameField] = "biz_for_biz_set_1"
		bizResp, err = apiServerClient.CreateBiz(ctx, header, biz)
		util.RegisterResponseWithRid(bizResp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(bizResp.Result).To(Equal(true))
		bizID2, err = commonutil.GetInt64ByInterface(bizResp.Data[common.BKAppIDField])
		Expect(err).NotTo(HaveOccurred())

		biz[common.BKAppNameField] = "biz_not_for_biz_set_1"
		bizResp, err = apiServerClient.CreateBiz(ctx, header, biz)
		util.RegisterResponseWithRid(bizResp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(bizResp.Result).To(Equal(true))
		_, err = commonutil.GetInt64ByInterface(bizResp.Data[common.BKAppIDField])
		Expect(err).NotTo(HaveOccurred())

		set := map[string]interface{}{
			common.BKSetNameField:  "set_for_biz_set_1",
			common.BKAppIDField:    bizID1,
			common.BKParentIDField: bizID1,
		}
		setResp, err := instClient.CreateSet(ctx, bizID1, header, set)
		util.RegisterResponseWithRid(setResp, header)
		Expect(err).NotTo(HaveOccurred())
		setID1, err = commonutil.GetInt64ByInterface(setResp[common.BKSetIDField])
		Expect(err).NotTo(HaveOccurred())

		set1 := map[string]interface{}{
			common.BKSetNameField:  "set_for_biz_set_2",
			common.BKAppIDField:    bizID2,
			common.BKParentIDField: bizID2,
		}
		setResp1, err := instClient.CreateSet(ctx, bizID2, header, set1)
		util.RegisterResponseWithRid(setResp1, header)
		Expect(err).NotTo(HaveOccurred())
		setID2, err = commonutil.GetInt64ByInterface(setResp1[common.BKSetIDField])
		Expect(err).NotTo(HaveOccurred())

		module := map[string]interface{}{
			common.BKModuleNameField: "module_for_biz_set_1",
			common.BKAppIDField:      bizID1,
			common.BKParentIDField:   setID1,
		}
		moduleResp, err := instClient.CreateModule(ctx, bizID1, setID1, header, module)
		util.RegisterResponseWithRid(moduleResp, header)
		moduleID1, err = commonutil.GetInt64ByInterface(moduleResp["bk_module_id"])
		Expect(err).NotTo(HaveOccurred())

		module3 := map[string]interface{}{
			common.BKModuleNameField: "module_for_biz_set_3",
			common.BKAppIDField:      bizID2,
			common.BKParentIDField:   setID2,
		}
		moduleResp3, err := instClient.CreateModule(context.Background(), bizID2, setID2, header, module3)
		util.RegisterResponseWithRid(moduleResp3, header)
		moduleID3, err = commonutil.GetInt64ByInterface(moduleResp3["bk_module_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create service and process template", func() {
		// get default service category id
		svrCond := map[string]interface{}{
			"name":         common.DefaultServiceCategoryName,
			"bk_parent_id": mapstr.MapStr{common.BKDBNE: 0},
		}
		serviceCategory := new(metadata.ServiceCategory)
		svrErr := test.GetDB().Table(common.BKTableNameServiceCategory).Find(svrCond).One(context.Background(),
			serviceCategory)
		Expect(svrErr).NotTo(HaveOccurred())
		categoryId := serviceCategory.ID

		// create service template
		serviceTemplateInput := map[string]interface{}{
			"service_category_id": categoryId,
			common.BKAppIDField:   bizID1,
			"name":                "service_template_1",
		}
		serviceTemplateRsp, err := serviceClient.CreateServiceTemplate(context.Background(), header,
			serviceTemplateInput)
		util.RegisterResponseWithRid(serviceTemplateRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(serviceTemplateRsp.Result).To(Equal(true), serviceTemplateRsp.ToString())
		j, err := json.Marshal(serviceTemplateRsp.Data)
		data := metadata.ServiceTemplate{}
		json.Unmarshal(j, &data)
		Expect(data.Name).To(Equal("service_template_1"))
		Expect(data.ServiceCategoryID).To(Equal(categoryId))
		serviceTemplateId = data.ID

		// create process template
		processTemplateInput := map[string]interface{}{
			"service_template_id": serviceTemplateId,
			common.BKAppIDField:   bizID1,
			"processes": []map[string]interface{}{
				{
					"spec": map[string]interface{}{
						"bk_func_name": map[string]interface{}{
							"value":            "p1_1",
							"as_default_value": true,
						},
						"bk_process_name": map[string]interface{}{
							"value":            "p1_1",
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
		processRsp, err := processClient.CreateProcessTemplate(context.Background(), header, processTemplateInput)
		util.RegisterResponseWithRid(processRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(processRsp.Result).To(Equal(true), processRsp.ToString())
		processTemplateId, err = commonutil.GetInt64ByInterface(processRsp.Data.([]interface{})[0])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create module by service template", func() {
		input := map[string]interface{}{
			"bk_module_name":      "module_for_biz_set_2",
			"bk_parent_id":        setID1,
			"service_template_id": serviceTemplateId,
		}
		rsp, e := instClient.CreateModule(context.Background(), bizID1, setID1, header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(e).NotTo(HaveOccurred())
		var err error
		moduleID2, err = commonutil.GetInt64ByInterface(rsp["bk_module_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create biz set", func() {
		createBizSetOpt := metadata.CreateBizSetRequest{
			BizSetAttr: map[string]interface{}{
				common.BKBizSetNameField: "sample_biz_set_1",
			},
			BizSetScope: &metadata.BizSetScope{
				MatchAll: false,
				Filter: &querybuilder.QueryFilter{
					Rule: querybuilder.CombinedRule{
						Condition: querybuilder.ConditionAnd,
						Rules: []querybuilder.Rule{
							querybuilder.AtomRule{
								Field:    common.BKAppIDField,
								Operator: querybuilder.OperatorIn,
								Value:    []int64{bizID1, bizID2},
							},
						},
					}},
			},
		}

		var err error
		sampleBizSetID1, err = instClient.CreateBizSet(context.Background(), header, createBizSetOpt)
		util.RegisterResponseWithRid(nil, header)
		Expect(err).NotTo(HaveOccurred())
	})

	It("prepare environment, add host to resource pool and transfer to biz set", func() {
		// add host
		input := map[string]interface{}{
			"bk_biz_id": bizID1,
			"host_info": map[string]interface{}{
				"1": map[string]interface{}{
					"bk_host_innerip": "127.0.0.3",
					"bk_asset_id":     "addhost_api_asset_1",
					"bk_cloud_id":     0,
				},
				"2": map[string]interface{}{
					"bk_host_innerip": "127.0.0.4",
					"bk_asset_id":     "addhost_api_asset_2",
					"bk_cloud_id":     0,
				},
			},
		}
		rsp, err := hostServerClient.AddHost(context.Background(), header, input)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))

		input1 := map[string]interface{}{
			"bk_biz_id": bizID2,
			"host_info": map[string]interface{}{
				"1": map[string]interface{}{
					"bk_host_innerip": "127.0.0.5",
					"bk_asset_id":     "addhost_api_asset_3",
					"bk_cloud_id":     0,
				},
			},
		}
		rsp1, err := hostServerClient.AddHost(context.Background(), header, input1)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp1.Result).To(Equal(true))

		searchInput := &metadata.HostCommonSearch{
			AppID: bizID1,
		}
		resp, err := hostServerClient.SearchHostWithBiz(context.Background(), header, searchInput)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Result).To(Equal(true))
		Expect(resp.Data.Count).To(Equal(2))
		hostID1, err = commonutil.GetInt64ByInterface(resp.Data.Info[0]["host"].(map[string]interface{})["bk_host_id"])
		Expect(err).NotTo(HaveOccurred())
		hostID2, err = commonutil.GetInt64ByInterface(resp.Data.Info[1]["host"].(map[string]interface{})["bk_host_id"])
		Expect(err).NotTo(HaveOccurred())

		searchInput1 := &metadata.HostCommonSearch{
			AppID: bizID2,
		}
		resp1, err := hostServerClient.SearchHostWithBiz(context.Background(), header, searchInput1)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp1.Result).To(Equal(true))
		Expect(resp1.Data.Count).To(Equal(1))
		hostID3, err = commonutil.GetInt64ByInterface(resp1.Data.Info[0]["host"].(map[string]interface{})["bk_host_id"])
		Expect(err).NotTo(HaveOccurred())

		transInput := map[string]interface{}{
			"bk_biz_id": bizID1,
			"bk_host_id": []int64{
				hostID1,
			},
			"bk_module_id": []int64{
				moduleID1,
			},
			"is_increment": true,
		}
		rsp, rawErr := hostServerClient.TransferHostModule(context.Background(), header, transInput)
		util.RegisterResponseWithRid(rsp, header)
		Expect(rawErr).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))

		transInput2 := map[string]interface{}{
			"bk_biz_id": bizID1,
			"bk_host_id": []int64{
				hostID2,
			},
			"bk_module_id": []int64{
				moduleID2,
			},
			"is_increment": true,
		}
		rsp, rawErr = hostServerClient.TransferHostModule(context.Background(), header, transInput2)
		util.RegisterResponseWithRid(rsp, header)
		Expect(rawErr).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))

		transInput3 := map[string]interface{}{
			"bk_biz_id": bizID2,
			"bk_host_id": []int64{
				hostID3,
			},
			"bk_module_id": []int64{
				moduleID3,
			},
			"is_increment": true,
		}
		rsp, rawErr = hostServerClient.TransferHostModule(context.Background(), header, transInput3)
		util.RegisterResponseWithRid(rsp, header)
		Expect(rawErr).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It(fmt.Sprintf("create service instance in biz %d", bizID1), func() {
		instInput := &metadata.CreateServiceInstanceInput{
			BizID:    bizID1,
			ModuleID: moduleID1,
			Instances: []metadata.CreateServiceInstanceDetail{
				{
					HostID: hostID1,
					Processes: []metadata.ProcessInstanceDetail{
						{
							ProcessTemplateID: processTemplateId,
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
		serviceIds, err := serviceClient.CreateServiceInstance(context.Background(), header, instInput)
		util.RegisterResponseWithRid(serviceIds, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(serviceIds)).NotTo(Equal(0))
		serviceId = serviceIds[0]
	})

	It(fmt.Sprintf("create service instance in biz %d", bizID2), func() {
		instInput := &metadata.CreateServiceInstanceInput{
			BizID:    bizID2,
			ModuleID: moduleID3,
			Instances: []metadata.CreateServiceInstanceDetail{
				{
					HostID: hostID3,
					Processes: []metadata.ProcessInstanceDetail{
						{
							ProcessTemplateID: processTemplateId,
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
		serviceIds, err := serviceClient.CreateServiceInstance(context.Background(), header, instInput)
		util.RegisterResponseWithRid(serviceIds, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(serviceIds)).NotTo(Equal(0))
	})

	It("search service instances in module", func() {
		svrInstInput := &metadata.GetServiceInstanceInModuleInput{
			BizID:    bizID1,
			ModuleID: moduleID1,
			HostIDs:  []int64{hostID1},
			Page:     metadata.BasePage{Start: 0, Limit: 10},
		}
		svrInstRsp, ccErr := processClient.SearchBizSetSrvInstInModule(context.Background(), sampleBizSetID1, header,
			svrInstInput)
		util.RegisterResponseWithRid(svrInstRsp, header)
		Expect(ccErr).NotTo(HaveOccurred())
		Expect(int(svrInstRsp.Count)).To(Equal(1))
		Expect(svrInstRsp.Info[0].Name).To(ContainSubstring("p1"))
		Expect(svrInstRsp.Info[0].HostID).To(Equal(hostID1))
		Expect(svrInstRsp.Info[0].ModuleID).To(Equal(moduleID1))
		Expect(svrInstRsp.Info[0].BizID).To(Equal(bizID1))
		Expect(svrInstRsp.Info[0].ID).To(Equal(serviceId))
	})

	It("list process instances", func() {
		processInstInput := &metadata.ListProcessInstancesOption{
			BizID:             bizID1,
			ServiceInstanceID: serviceId,
		}
		processInstData, ccErr := processClient.ListBizSetProcessInstances(context.Background(), sampleBizSetID1,
			header, processInstInput)
		util.RegisterResponseWithRid(processInstData, header)
		Expect(ccErr).NotTo(HaveOccurred())
		Expect(len(processInstData)).To(Equal(1))
		Expect(processInstData[0].Relation.BizID).To(Equal(bizID1))
		Expect(processInstData[0].Relation.HostID).To(Equal(hostID1))
		Expect(processInstData[0].Relation.ServiceInstanceID).To(Equal(serviceId))
		Expect(processInstData[0].Property["bk_process_name"]).To(Equal("p1"))
		Expect(processInstData[0].Property["bk_func_name"]).To(Equal("p1"))
	})

	It("test search process template and process info", func() {
		processTemplateData := &metadata.GetBizSetProcTemplateOption{
			BizID: bizID1,
		}
		processTemplateRsp, err := processClient.GetBizSetProcessTemplate(context.Background(), sampleBizSetID1,
			processTemplateId, header, processTemplateData)
		Expect(err).NotTo(HaveOccurred())
		Expect(processTemplateRsp.ProcessName).To(Equal("p1_1"))
		Expect(processTemplateRsp.BizID).To(Equal(bizID1))
		Expect(processTemplateRsp.ServiceTemplateID).To(Equal(serviceTemplateId))
		processTemplateID = processTemplateRsp.ID

		// list service instance
		svrInstanceData := &metadata.ListServiceInstancesWithHostInput{
			BizID:  bizID1,
			HostID: hostID2,
		}
		svrInstanceRsp, err := processClient.ListBizSetSrvInstWithHost(context.Background(), sampleBizSetID1,
			header, svrInstanceData)
		Expect(err).NotTo(HaveOccurred())
		Expect(int(svrInstanceRsp.Count)).To(Equal(1))
		Expect(svrInstanceRsp.Info[0].Name).To(Equal("127.0.0.4_p1_1"))
		Expect(svrInstanceRsp.Info[0].BizID).To(Equal(bizID1))
		Expect(svrInstanceRsp.Info[0].HostID).To(Equal(hostID2))
		Expect(svrInstanceRsp.Info[0].ModuleID).To(Equal(moduleID2))
		svrID := svrInstanceRsp.Info[0].ID

		// list process info
		input := &metadata.ListProcessInstancesOption{
			BizID:             bizID1,
			ServiceInstanceID: svrID,
		}
		data, err := processClient.SearchProcessInstance(context.Background(), header, input)
		util.RegisterResponseWithRid(data, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(data)).To(Equal(1))
		Expect(data[0].Property["bk_process_name"]).To(Equal("p1_1"))
		Expect(data[0].Property["bk_func_name"]).To(Equal("p1_1"))
		serviceInstanceID, err := commonutil.GetInt64ByInterface(data[0].Property["service_instance_id"])
		Expect(err).NotTo(HaveOccurred())
		Expect(serviceInstanceID).To(Equal(svrID))
		Expect(data[0].Relation.ProcessTemplateID).To(Equal(processTemplateID))
		Expect(data[0].Relation.HostID).To(Equal(hostID2))
	})
})
