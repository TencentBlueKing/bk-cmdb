package host_server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("host test", func() {
	var bizId, setId, moduleId, idleModuleId, faultModuleId, dirID int64
	var hostId, hostId1, hostId2, hostId3 int64

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
			rsp, e := instClient.CreateSet(context.Background(), bizId, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(e).NotTo(HaveOccurred())
			Expect(rsp["bk_set_name"].(string)).To(Equal("test"))
			parentIdRes, err := commonutil.GetInt64ByInterface(rsp["bk_parent_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(parentIdRes).To(Equal(bizId))
			bizIdRes, err := commonutil.GetInt64ByInterface(rsp["bk_biz_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(bizIdRes).To(Equal(bizId))
			setId, err = commonutil.GetInt64ByInterface(rsp["bk_set_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("create module", func() {
			input := map[string]interface{}{
				"bk_module_name":      "cc_module",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, e := instClient.CreateModule(context.Background(), bizId, setId, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(e).NotTo(HaveOccurred())
			Expect(rsp["bk_module_name"].(string)).To(Equal("cc_module"))

			setIdRes, err := commonutil.GetInt64ByInterface(rsp["bk_set_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(setIdRes).To(Equal(setId))

			parentIdRes, err := commonutil.GetInt64ByInterface(rsp["bk_parent_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(parentIdRes).To(Equal(setId))
			moduleId, err = commonutil.GetInt64ByInterface(rsp["bk_module_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("create resource pool directory", func() {
			dir := map[string]interface{}{
				"bk_module_name": "test",
			}
			rsp, err := dirClient.CreateResourceDirectory(context.Background(), header, dir)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			dirID = int64(rsp.Data.Created.ID)
		})
	})

	Describe("add host test", func() {
		It("add host using api", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"host_info": map[string]interface{}{
					"4": map[string]interface{}{
						"bk_host_innerip": "127.0.0.1",
						"bk_asset_id":     "addhost_api_asset_1",
						"bk_cloud_id":     0,
						"bk_comment":      "127.0.0.1 comment",
					},
				},
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search host created using api", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Ipv4Ip: metadata.IPInfo{
					Data:  []string{"127.0.0.1"},
					Exact: 1,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.1"))
			Expect(data["bk_asset_id"].(string)).To(Equal("addhost_api_asset_1"))
			hostId1, err = commonutil.GetInt64ByInterface(data["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("add host using excel", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"host_info": map[string]interface{}{
					"5": map[string]interface{}{
						"bk_host_innerip": "127.0.0.2",
						"bk_asset_id":     "addhost_excel_asset_1",
						"bk_host_name":    "127.0.0.2",
					},
				},
				"input_type": "excel",
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search host created using excel", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Ipv4Ip: metadata.IPInfo{
					Data:  []string{"127.0.0.2"},
					Exact: 1,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.2"))
			Expect(data["bk_asset_id"].(string)).To(Equal("addhost_excel_asset_1"))
			Expect(data["bk_host_name"].(string)).To(Equal("127.0.0.2"))
			hostId2, err = commonutil.GetInt64ByInterface(data["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("search host using multiple ips", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Ipv4Ip: metadata.IPInfo{
					Data: []string{
						"127.0.0.1",
						"127.0.0.2",
					},
					Exact: 1,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		// This api is marked as deprecated
		// It("add host using agent", func() {
		// 	input := map[string]interface{}{
		// 		"host_info": map[string]interface{}{
		// 			"bk_host_innerip": "127.0.0.3",
		// 			"bk_asset_id":     "addhost_agent_asset_1",
		// 			"bk_cloud_id":     0,
		// 		},
		// 	}
		// 	rsp, err := hostServerClient.AddHostFromAgent(context.Background(), header, input)
		// 	util.RegisterResponseWithRid(rsp, header)
		// 	Expect(err).NotTo(HaveOccurred())
		// 	Expect(rsp.Result).To(Equal(true))
		// })

		// It("search host created using agent", func() {
		// 	input := &metadata.HostCommonSearch{
		// 		AppID: bizId,
		// 		Ip: metadata.IPInfo{
		// 			Data:  []string{"127.0.0.3"},
		// 			Exact: 1,
		// 			Flag:  "bk_host_innerip|bk_host_outerip",
		// 		},
		// 		Page: metadata.BasePage{
		// 			Sort: "bk_host_id",
		// 		},
		// 	}
		// 	rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
		// 	util.RegisterResponseWithRid(rsp, header)
		// 	Expect(err).NotTo(HaveOccurred())
		// 	Expect(rsp.Result).To(Equal(true))
		// 	Expect(rsp.Data.Count).To(Equal(1))
		// 	data := rsp.Data.Info[0]["host"].(map[string]interface{})
		// 	Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.3"))
		// 	Expect(data["bk_asset_id"].(string)).To(Equal("addhost_agent_asset_1"))
		// })

		It("add host to resource", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"4": map[string]interface{}{
						"bk_host_innerip": "127.0.0.4",
						"bk_cloud_id":     0,
					},
				},
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search resource host", func() {
			input := &metadata.HostCommonSearch{
				AppID: -1,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "biz",
						Condition: []metadata.ConditionItem{
							{
								Field:    "default",
								Operator: "$eq",
								Value:    1,
							},
						},
						Fields: []string{},
					},
				},
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.4"))
			hostId, err = commonutil.GetInt64ByInterface(data["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("get host base info", func() {
			rsp, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0",
				strconv.FormatInt(hostId, 10), header)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			for _, data := range rsp.Data {
				if data.PropertyID == "bk_host_innerip" {
					Expect(data.PropertyValue).To(Equal("127.0.0.4"))
					break
				}
			}
		})

		It("get host count in multi cloud area", func() {
			opt := metadata.CloudAreaHostCount{CloudIDs: []int64{0, 100}}
			rsp, err := hostServerClient.FindCloudAreaHostCount(context.Background(), header, opt)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(len(rsp.Data)).To(Equal(2))
		})
	})

	Describe("transfer host test", func() {
		It("search biz host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			Expect(rsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("127.0.0.1"))
			Expect(rsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("127.0.0.2"))
		})

		It("transfer resourcehost to idlemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostIDs: []int64{
					hostId,
				},
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search biz host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(3))
			Expect(rsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("127.0.0.1"))
			Expect(rsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("127.0.0.2"))
			Expect(rsp.Data.Info[2]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("127.0.0.4"))
		})

		It("transfer host to resourcemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostIDs: []int64{
					hostId1,
				},
				ModuleID: dirID,
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search biz host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			Expect(rsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("127.0.0.2"))
			Expect(rsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_innerip"].(string)).To(Equal("127.0.0.4"))
		})

		It("transfer host module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId2,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": true,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("transfer host same module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId2,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": true,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search module host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    moduleId,
							},
						},
						Fields: []string{},
					},
				},
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			host := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(host["bk_host_innerip"].(string)).To(Equal("127.0.0.2"))

			hostIdRes, err := commonutil.GetInt64ByInterface(host["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(hostIdRes).To(Equal(hostId2))
			module := rsp.Data.Info[0]["module"].([]interface{})[0].(map[string]interface{})
			Expect(module["bk_module_name"].(string)).To(Equal("cc_module"))

			moduleIdRes, err := commonutil.GetInt64ByInterface(module["bk_module_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(moduleIdRes).To(Equal(moduleId))
		})

		It("add clone destion host", func() {
			input := map[string]interface{}{
				"bk_biz_id": 1,
				"host_info": map[string]interface{}{
					"4": map[string]interface{}{
						"bk_host_innerip": "127.0.0.5",
						"bk_asset_id":     "add_clone_destion_host",
						"bk_cloud_id":     0,
					},
				},
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})
		It("clone host", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   1,
				OrgIP:   "127.0.0.1",
				DstIP:   "127.0.0.5",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search cloned host", func() {
			input := &metadata.HostCommonSearch{
				AppID: -1,
				Ipv4Ip: metadata.IPInfo{
					Data:  []string{"127.0.0.5"},
					Exact: 0,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "biz",
						Condition: []metadata.ConditionItem{
							{
								Field:    "default",
								Operator: "$eq",
								Value:    1,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.5"))
			Expect(data["bk_comment"].(string)).To(Equal("127.0.0.1 comment"))

		})

		It("get instance topo", func() {
			rsp, err := instClient.GetInternalModule(context.Background(), "0", strconv.FormatInt(bizId, 10), header)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.SetName).To(Equal("空闲机池"))
			Expect(len(rsp.Data.Module)).To(Equal(3))
			for _, module := range rsp.Data.Module {
				switch module.ModuleName {
				case "空闲机":
					idleModuleId = module.ModuleID
				case "故障机":
					faultModuleId = module.ModuleID
				}
			}
		})

		It("search fault host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    faultModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("transfer host to fault module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId2,
				},
				"bk_module_id": []int64{
					faultModuleId,
				},
				"is_increment": true,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search fault host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    faultModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("transfer host to fault module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostIDs: []int64{
					hostId2,
				},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search fault host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    faultModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.2"))
		})

		It("search transfered module host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    moduleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("transfer fault host to idle module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId2,
				},
				"bk_module_id": []int64{
					idleModuleId,
				},
				"is_increment": true,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search idle host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.4"))
		})

		It("transfer host to idle module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostIDs: []int64{
					hostId2,
				},
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search idle host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			data1 := rsp.Data.Info[1]["host"].(map[string]interface{})
			Expect("127.0.0.2").To(SatisfyAny(Equal(data["bk_host_innerip"].(string)),
				Equal(data1["bk_host_innerip"].(string))))
		})

		It("search fault host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    faultModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search resource host", func() {
			input := &metadata.HostCommonSearch{
				AppID: -1,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "biz",
						Condition: []metadata.ConditionItem{
							{
								Field:    "default",
								Operator: "$eq",
								Value:    1,
							},
						},
						Fields: []string{},
					},
				},
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		It("transfer host to resource pool default directory", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostIDs: []int64{
					hostId2,
				},
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search resource host", func() {
			input := &metadata.HostCommonSearch{
				AppID: -1,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "biz",
						Condition: []metadata.ConditionItem{
							{
								Field:    "default",
								Operator: "$eq",
								Value:    1,
							},
						},
						Fields: []string{},
					},
				},
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(3))
		})

		It("search resource host change start limit", func() {
			input := &metadata.HostCommonSearch{
				AppID: -1,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "biz",
						Condition: []metadata.ConditionItem{
							{
								Field:    "default",
								Operator: "$eq",
								Value:    1,
							},
						},
						Fields: []string{},
					},
				},
				Page: metadata.BasePage{
					Sort:  "bk_host_id",
					Start: 2,
					Limit: 2,
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(3))
			Expect(len(rsp.Data.Info)).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.5"))
		})

		It("search idle host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
		})

		It("add host", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "127.0.0.6",
						"bk_asset_id":     "host_sync_asset_1",
						"bk_cloud_id":     0,
					},
				},
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search idle host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			data := rsp.Data.Info[1]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.6"))
			hostId3, err = commonutil.GetInt64ByInterface(data["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("transfer host module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": true,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("transfer host module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId3,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": true,
			}
			rsp, err := hostServerClient.TransferHostModule(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search idle host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search module host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    moduleId,
							},
						},
						Fields: []string{},
					},
				},
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.4"))
			data1 := rsp.Data.Info[1]["host"].(map[string]interface{})
			Expect(data1["bk_host_innerip"].(string)).To(Equal("127.0.0.6"))
		})

		It("move all module hosts to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: bizId,
				SetID:         setId,
				ModuleID:      moduleId,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search module host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    moduleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search idle host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		It("update idle host", func() {
			input := map[string]interface{}{
				"bk_host_id": fmt.Sprintf("%v,%v", hostId, hostId3),
				"bk_sn":      "update_bk_sn",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search idle host", func() {
			input := &metadata.HostCommonSearch{
				AppID: bizId,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "module",
						Condition: []metadata.ConditionItem{
							{
								Field:    "bk_module_id",
								Operator: "$eq",
								Value:    idleModuleId,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_sn"].(string)).To(Equal("update_bk_sn"))
			data1 := rsp.Data.Info[1]["host"].(map[string]interface{})
			Expect(data1["bk_sn"].(string)).To(Equal("update_bk_sn"))
		})

		It("delete resource host", func() {
			input := map[string]interface{}{
				"bk_host_id":          fmt.Sprintf("%v,%v", hostId1, hostId2),
				"bk_supplier_account": "0",
			}
			rsp, err := hostServerClient.DeleteHostBatch(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search resource host", func() {
			input := &metadata.HostCommonSearch{
				AppID: -1,
				Condition: []metadata.SearchCondition{
					{
						ObjectID: "biz",
						Condition: []metadata.ConditionItem{
							{
								Field:    "default",
								Operator: "$eq",
								Value:    1,
							},
						},
						Fields: []string{},
					},
				},
				Page: metadata.BasePage{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
		})
	})
})

var _ = Describe("list_hosts_topo test", func() {
	It("list_hosts_topo", func() {
		test.DeleteAllBizs()

		By("create biz cc_biz_test")
		bizInput := map[string]interface{}{
			"life_cycle":        "2",
			"language":          "1",
			"bk_biz_maintainer": "admin",
			"bk_biz_name":       "cc_biz_test",
			"time_zone":         "Africa/Accra",
		}
		bizRsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, bizInput)
		util.RegisterResponse(bizRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(bizRsp.Result).To(Equal(true))
		bizId, err := commonutil.GetInt64ByInterface(bizRsp.Data[common.BKAppIDField])
		Expect(err).NotTo(HaveOccurred())

		By("create set cc_set")
		setInput := mapstr.MapStr{
			"bk_set_name":         "cc_set",
			"bk_parent_id":        bizId,
			"bk_supplier_account": "0",
			"bk_biz_id":           bizId,
			"bk_service_status":   "1",
			"bk_set_env":          "3",
		}
		setRsp, err := instClient.CreateSet(context.Background(), bizId, header, setInput)
		util.RegisterResponse(setRsp)
		Expect(err).NotTo(HaveOccurred())
		setId, err := commonutil.GetInt64ByInterface(setRsp[common.BKSetIDField])
		Expect(err).NotTo(HaveOccurred())

		By("create module cc_module")
		moduleInput := map[string]interface{}{
			"bk_module_name":      "cc_module",
			"bk_parent_id":        setId,
			"service_category_id": 2,
			"service_template_id": 0,
		}
		moduleRsp, err := instClient.CreateModule(context.Background(), bizId, setId, header, moduleInput)
		util.RegisterResponseWithRid(moduleRsp, header)
		Expect(err).NotTo(HaveOccurred())
		moduleId1, err := commonutil.GetInt64ByInterface(moduleRsp[common.BKModuleIDField])
		Expect(err).NotTo(HaveOccurred())

		By("create module cc_module1")
		moduleInput1 := map[string]interface{}{
			"bk_module_name":      "cc_module1",
			"bk_parent_id":        setId,
			"service_category_id": 2,
			"service_template_id": 0,
		}
		moduleRsp1, err := instClient.CreateModule(context.Background(), bizId, setId, header, moduleInput1)
		util.RegisterResponseWithRid(moduleRsp1, header)
		Expect(err).NotTo(HaveOccurred())
		moduleId2, err := commonutil.GetInt64ByInterface(moduleRsp1[common.BKModuleIDField])
		Expect(err).NotTo(HaveOccurred())

		By("add host using api")
		hostInput := map[string]interface{}{
			"bk_biz_id": bizId,
			"host_info": map[string]interface{}{
				"4": map[string]interface{}{
					"bk_host_innerip": "127.0.0.1",
				},
				"5": map[string]interface{}{
					"bk_host_innerip": "127.0.0.2",
				},
			},
		}
		hostRsp, err := hostServerClient.AddHost(context.Background(), header, hostInput)
		util.RegisterResponseWithRid(hostRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(true), hostRsp.ToString())

		By("search hosts")
		searchInput := &metadata.HostCommonSearch{
			AppID: bizId,
			Page: metadata.BasePage{
				Sort: common.BKHostIDField,
			},
		}
		searchRsp, err := hostServerClient.SearchHost(context.Background(), header, searchInput)
		util.RegisterResponseWithRid(searchRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		hostId1, err := commonutil.GetInt64ByInterface(searchRsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_id"])
		Expect(err).NotTo(HaveOccurred())
		hostId2, err := commonutil.GetInt64ByInterface(searchRsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_id"])
		Expect(err).NotTo(HaveOccurred())

		By("transfer host module")
		transferInput := map[string]interface{}{
			"bk_biz_id": bizId,
			"bk_host_id": []int64{
				hostId1,
			},
			"bk_module_id": []int64{
				moduleId1,
				moduleId2,
			},
			"is_increment": true,
		}
		transferRsp, err := hostServerClient.TransferHostModule(context.Background(), header, transferInput)
		util.RegisterResponse(transferRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(transferRsp.Result).To(Equal(true))

		By("list hosts topo")
		rsp, err := hostServerClient.ListBizHostsTopo(context.Background(), header, bizId,
			&metadata.ListHostsWithNoBizParameter{Page: metadata.BasePage{Sort: common.BKHostIDField, Limit: 10},
				Fields: []string{"bk_host_id"}})
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		j, err := json.Marshal(rsp.Data)
		Expect(j).To(MatchRegexp(`.*"count":2.*`))
		Expect(j).To(MatchRegexp(fmt.Sprintf(`.*"bk_host_id":%d.*`, hostId1)))
		Expect(j).To(MatchRegexp(fmt.Sprintf(`.*"bk_host_id":%d.*`, hostId2)))
		Expect(j).To(MatchRegexp(fmt.Sprintf(`.*"bk_set_id":%d.*`, setId)))
		Expect(j).To(MatchRegexp(fmt.Sprintf(`.*"bk_module_id":%d.*`, moduleId1)))
		Expect(j).To(MatchRegexp(fmt.Sprintf(`.*"bk_module_id":%d.*`, moduleId2)))
	})
})

var _ = Describe("batch_update_host test", func() {
	It("batch_update_host", func() {
		test.DeleteAllHosts()

		By("add host using api")
		hostInput := map[string]interface{}{
			"host_info": map[string]interface{}{
				"4": map[string]interface{}{
					"bk_host_innerip": "127.0.0.1",
				},
				"5": map[string]interface{}{
					"bk_host_innerip": "127.0.0.2",
				},
			},
		}
		hostRsp, err := hostServerClient.AddHost(context.Background(), header, hostInput)
		util.RegisterResponseWithRid(hostRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(true), hostRsp.ToString())

		By("search hosts")
		searchInput := &metadata.HostCommonSearch{
			Page: metadata.BasePage{
				Sort: common.BKHostIDField,
			},
		}
		searchRsp, err := hostServerClient.SearchHost(context.Background(), header, searchInput)
		util.RegisterResponseWithRid(searchRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		hostId1, err := commonutil.GetInt64ByInterface(searchRsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_id"])
		Expect(err).NotTo(HaveOccurred())
		hostId2, err := commonutil.GetInt64ByInterface(searchRsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_id"])
		Expect(err).NotTo(HaveOccurred())

		By("update host property batch, bk_asset_id is not editable field")
		updateInput := map[string]interface{}{
			"update": []map[string]interface{}{
				{
					common.BKHostIDField: hostId1,
					"properties": map[string]interface{}{
						"bk_asset_id": "batch_update1",
						"operator":    "admin",
						"bk_comment":  "test",
						"bk_isp_name": "1",
					},
				},
				{
					common.BKHostIDField: hostId2,
					"properties": map[string]interface{}{
						"bk_bak_operator": "admin",
						"bk_host_outerip": "127.2.3.4",
					},
				},
			},
		}
		err = hostServerClient.UpdateHostPropertyBatch(context.Background(), header, updateInput)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())

		By("search updated host property")
		searchRsp, err = hostServerClient.SearchHost(context.Background(), header, searchInput)
		util.RegisterResponseWithRid(searchRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		Expect(searchRsp.Data.Info[0]["host"].(map[string]interface{})["bk_asset_id"].(string)).To(Equal(""))
		Expect(searchRsp.Data.Info[0]["host"].(map[string]interface{})["operator"].(string)).To(Equal("admin"))
		Expect(searchRsp.Data.Info[0]["host"].(map[string]interface{})["bk_comment"].(string)).To(Equal("test"))
		Expect(searchRsp.Data.Info[0]["host"].(map[string]interface{})["bk_isp_name"].(string)).To(Equal("1"))
		Expect(searchRsp.Data.Info[1]["host"].(map[string]interface{})["bk_bak_operator"].(string)).To(Equal("admin"))
		Expect(searchRsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_outerip"].(string)).To(Equal("127.2.3.4"))
	})
})

var _ = Describe("multiple ip host validation test", func() {
	It("multiple ip host validation", func() {
		test.DeleteAllHosts()

		By("add hosts with different ip using api")
		hostInput := map[string]interface{}{
			"host_info": map[string]interface{}{
				"1": map[string]interface{}{
					"bk_host_innerip": "1.0.0.1,1.0.0.2",
				},
				"2": map[string]interface{}{
					"bk_host_innerip": "1.0.0.3",
				},
			},
		}
		hostRsp, err := hostServerClient.AddHost(context.Background(), header, hostInput)
		util.RegisterResponseWithRid(hostRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(true), hostRsp.ToString())

		By("search hosts")
		searchInput := &metadata.HostCommonSearch{
			Page: metadata.BasePage{
				Sort: common.BKHostIDField,
			},
		}
		searchRsp, err := hostServerClient.SearchHost(context.Background(), header, searchInput)
		util.RegisterResponseWithRid(searchRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		Expect(searchRsp.Data.Count).To(Equal(2))

		By("add same multiple ip host using api")
		input := &metadata.CreateModelInstance{
			Data: map[string]interface{}{
				"bk_host_innerip": "1.0.0.1,1.0.0.2",
				"bk_cloud_id":     0,
			},
		}
		addHostResult, err := test.GetClientSet().CoreService().Instance().CreateInstance(context.Background(), header,
			common.BKInnerObjIDHost, input)
		util.RegisterResponse(addHostResult)
		Expect(err).To(HaveOccurred())

		By("add hosts with one same ip using api")
		input = &metadata.CreateModelInstance{
			Data: map[string]interface{}{
				"bk_host_innerip": "1.0.0.1",
				"bk_cloud_id":     0,
			},
		}
		addHostResult, err = test.GetClientSet().CoreService().Instance().CreateInstance(context.Background(), header,
			common.BKInnerObjIDHost, input)
		util.RegisterResponse(addHostResult)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("add_host_to_resource_pool test", func() {
	It("add_host_to_resource_pool", func() {
		test.DeleteAllHosts()

		By("add hosts to resource pool default module")
		hostInput := metadata.AddHostToResourcePoolHostList{
			HostInfo: []map[string]interface{}{
				{
					common.BKHostInnerIPField: "1.0.0.1",
					common.BKCloudIDField:     0,
				},
				{
					common.BKHostInnerIPField: "1.0.0.2",
					common.BKCloudIDField:     0,
				},
			},
		}
		hostRsp, err := hostServerClient.AddHostToResourcePool(context.Background(), header, hostInput)
		util.RegisterResponseWithRid(hostRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(true), hostRsp.ToString())
		js, err := json.Marshal(hostRsp.Data)
		Expect(err).NotTo(HaveOccurred())
		result := metadata.AddHostToResourcePoolResult{}
		err = json.Unmarshal(js, &result)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result.Success)).To(Equal(2))
		Expect(len(result.Error)).To(Equal(0))
		var hostID1, hostID2 int64
		if result.Success[0].Index == 0 {
			hostID1 = result.Success[0].HostID
			hostID2 = result.Success[1].HostID
		} else {
			hostID1 = result.Success[1].HostID
			hostID2 = result.Success[0].HostID
		}

		By("search hosts")
		searchRsp, err := hostServerClient.SearchHost(context.Background(), header, &metadata.HostCommonSearch{})
		util.RegisterResponseWithRid(searchRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		Expect(searchRsp.Data.Count).To(Equal(2))
		host1 := searchRsp.Data.Info[0]["host"].(map[string]interface{})
		host2 := searchRsp.Data.Info[1]["host"].(map[string]interface{})
		host1ID, err := commonutil.GetInt64ByInterface(host1[common.BKHostIDField])
		Expect(err).NotTo(HaveOccurred())
		host2ID, err := commonutil.GetInt64ByInterface(host2[common.BKHostIDField])
		Expect(err).NotTo(HaveOccurred())
		if host1[common.BKHostInnerIPField] == "1.0.0.1" {
			Expect(host1ID).To(Equal(hostID1))
			Expect(host2ID).To(Equal(hostID2))
		} else {
			Expect(host1ID).To(Equal(hostID2))
			Expect(host2ID).To(Equal(hostID1))
		}

		By("add hosts to resource pool invalid module")
		hostInput = metadata.AddHostToResourcePoolHostList{
			HostInfo: []map[string]interface{}{
				{
					common.BKHostInnerIPField: "1.0.0.3",
					common.BKCloudIDField:     0,
				},
				{
					common.BKHostInnerIPField: "1.0.0.4",
					common.BKCloudIDField:     0,
				},
			},
			Directory: 1000,
		}
		hostRsp, err = hostServerClient.AddHostToResourcePool(context.Background(), header, hostInput)
		util.RegisterResponseWithRid(hostRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(false))

		By("search hosts")
		searchRsp, err = hostServerClient.SearchHost(context.Background(), header, &metadata.HostCommonSearch{})
		util.RegisterResponseWithRid(searchRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		Expect(searchRsp.Data.Count).To(Equal(2))

		By("add hosts to resource pool one invalid host")
		hostInput = metadata.AddHostToResourcePoolHostList{
			HostInfo: []map[string]interface{}{
				{
					common.BKHostInnerIPField: "1.0.0.5",
					common.BKCloudIDField:     0,
				},
				{
					"bk_host_innerip":     "",
					common.BKCloudIDField: 0,
				},
			},
		}
		hostRsp, err = hostServerClient.AddHostToResourcePool(context.Background(), header, hostInput)
		util.RegisterResponseWithRid(hostRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(false))
		js, err = json.Marshal(hostRsp.Data)
		Expect(err).NotTo(HaveOccurred())
		result = metadata.AddHostToResourcePoolResult{}
		err = json.Unmarshal(js, &result)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result.Success)).To(Equal(1))
		Expect(result.Success[0].Index).To(Equal(0))
		Expect(len(result.Error)).To(Equal(1))
		Expect(result.Error[0].Index).To(Equal(1))

		By("search hosts")
		searchRsp, err = hostServerClient.SearchHost(context.Background(), header, &metadata.HostCommonSearch{
			Page: metadata.BasePage{
				Sort: common.BKHostIDField,
			},
		})
		util.RegisterResponseWithRid(searchRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		Expect(searchRsp.Data.Count).To(Equal(3))
		hostID, err := commonutil.GetInt64ByInterface(searchRsp.Data.Info[2]["host"].(map[string]interface{})[common.BKHostIDField])
		Expect(err).NotTo(HaveOccurred())
		Expect(hostID).To(Equal(result.Success[0].HostID))
	})
})

var _ = Describe("bind & unbind host agent test", func() {
	var hostID, hostID1 int64
	agentID, agentID1 := "11111", "22222"

	hostParam := &metadata.HostCommonSearch{
		Page: metadata.BasePage{Sort: common.BKHostIDField},
	}

	It("bind host agent", func() {
		By("add hosts using api")
		hostInput := map[string]interface{}{
			"host_info": map[string]interface{}{
				"0": map[string]interface{}{
					"bk_host_innerip": "127.0.0.10",
				},
				"1": map[string]interface{}{
					"bk_host_innerip": "127.0.0.11",
				},
			},
		}
		hostRsp, err := hostServerClient.AddHost(context.Background(), header, hostInput)
		util.RegisterResponseWithRid(hostRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(true), hostRsp.ToString())

		By("get host ids")
		searchRsp, err := hostServerClient.SearchHost(context.Background(), header, &metadata.HostCommonSearch{
			Page: metadata.BasePage{Sort: common.BKHostIDField},
		})
		util.RegisterResponseWithRid(searchRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		Expect(searchRsp.Data.Count >= 2).To(Equal(true))
		host := searchRsp.Data.Info[0]["host"].(map[string]interface{})
		hostID, err = commonutil.GetInt64ByInterface(host[common.BKHostIDField])
		Expect(err).NotTo(HaveOccurred())
		host1 := searchRsp.Data.Info[1]["host"].(map[string]interface{})
		hostID1, err = commonutil.GetInt64ByInterface(host1[common.BKHostIDField])
		Expect(err).NotTo(HaveOccurred())

		By("bind agent to host")
		err = hostServerClient.BindAgent(context.Background(), header, &metadata.BindAgentParam{
			List: []metadata.HostAgentRelation{{
				HostID:  hostID,
				AgentID: agentID,
			}, {
				HostID:  hostID1,
				AgentID: agentID1,
			}},
		})
		util.RegisterResponseWithRid(nil, header)
		Expect(err).NotTo(HaveOccurred())

		By("check if agent id is bound to host")
		searchRsp, err = hostServerClient.SearchHost(context.Background(), header, hostParam)
		util.RegisterResponseWithRid(searchRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		Expect(searchRsp.Data.Count >= 2).To(Equal(true))
		host = searchRsp.Data.Info[0]["host"].(map[string]interface{})
		Expect(commonutil.GetStrByInterface(host[common.BKAgentIDField])).To(Equal(agentID))
		host = searchRsp.Data.Info[1]["host"].(map[string]interface{})
		Expect(commonutil.GetStrByInterface(host[common.BKAgentIDField])).To(Equal(agentID1))

		By("bind agent to host again")
		err = hostServerClient.BindAgent(context.Background(), header, &metadata.BindAgentParam{
			List: []metadata.HostAgentRelation{{
				HostID:  hostID,
				AgentID: agentID,
			}},
		})
		util.RegisterResponseWithRid(nil, header)
		Expect(err).NotTo(HaveOccurred())

		By("update host agents")
		agentID, agentID1 = "33333", "444444"
		err = hostServerClient.BindAgent(context.Background(), header, &metadata.BindAgentParam{
			List: []metadata.HostAgentRelation{{
				HostID:  hostID,
				AgentID: agentID,
			}, {
				HostID:  hostID1,
				AgentID: agentID1,
			}},
		})
		util.RegisterResponseWithRid(nil, header)
		Expect(err).NotTo(HaveOccurred())

		By("check if agent id is updated to host")
		searchRsp, err = hostServerClient.SearchHost(context.Background(), header, hostParam)
		util.RegisterResponseWithRid(searchRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		Expect(searchRsp.Data.Count >= 2).To(Equal(true))
		host = searchRsp.Data.Info[0]["host"].(map[string]interface{})
		Expect(commonutil.GetStrByInterface(host[common.BKAgentIDField])).To(Equal(agentID))
		host = searchRsp.Data.Info[1]["host"].(map[string]interface{})
		Expect(commonutil.GetStrByInterface(host[common.BKAgentIDField])).To(Equal(agentID1))

		By("bind empty agent id to host")
		err = hostServerClient.BindAgent(context.Background(), header, &metadata.BindAgentParam{
			List: []metadata.HostAgentRelation{{
				HostID:  hostID,
				AgentID: "",
			}},
		})
		util.RegisterResponseWithRid(nil, header)
		Expect(err).To(HaveOccurred())

		By("bind agent id to not exist host")
		err = hostServerClient.BindAgent(context.Background(), header, &metadata.BindAgentParam{
			List: []metadata.HostAgentRelation{{
				HostID:  111111,
				AgentID: agentID,
			}},
		})
		util.RegisterResponseWithRid(nil, header)
		Expect(err).To(HaveOccurred())

		By("bind duplicate agent id to host")
		err = hostServerClient.BindAgent(context.Background(), header, &metadata.BindAgentParam{
			List: []metadata.HostAgentRelation{{
				HostID:  hostID,
				AgentID: agentID1,
			}},
		})
		util.RegisterResponseWithRid(nil, header)
		Expect(err).To(HaveOccurred())
	})

	It("unbind host agent", func() {
		By("unbind host agent with mismatch agent")
		err := hostServerClient.UnbindAgent(context.Background(), header, &metadata.UnbindAgentParam{
			List: []metadata.HostAgentRelation{{
				HostID:  hostID,
				AgentID: "333333",
			}},
		})
		util.RegisterResponseWithRid(nil, header)
		Expect(err).To(HaveOccurred())

		By("unbind agent id with not exist host")
		err = hostServerClient.UnbindAgent(context.Background(), header, &metadata.UnbindAgentParam{
			List: []metadata.HostAgentRelation{{
				HostID:  111111,
				AgentID: agentID,
			}},
		})
		util.RegisterResponseWithRid(nil, header)
		Expect(err).To(HaveOccurred())

		By("check if agent id is still bound to host")
		searchRsp, rawErr := hostServerClient.SearchHost(context.Background(), header, hostParam)
		util.RegisterResponseWithRid(searchRsp, header)
		Expect(rawErr).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		Expect(searchRsp.Data.Count >= 1).To(Equal(true))
		host := searchRsp.Data.Info[0]["host"].(map[string]interface{})
		Expect(commonutil.GetStrByInterface(host[common.BKAgentIDField])).To(Equal(agentID))

		By("unbind host agent")
		param := &metadata.UnbindAgentParam{
			List: []metadata.HostAgentRelation{{
				HostID:  hostID,
				AgentID: agentID,
			}, {
				HostID:  hostID1,
				AgentID: agentID1,
			}},
		}
		err = hostServerClient.UnbindAgent(context.Background(), header, param)
		util.RegisterResponseWithRid(nil, header)
		Expect(err).NotTo(HaveOccurred())

		By("check if agent id is not bound to host")
		searchRsp, rawErr = hostServerClient.SearchHost(context.Background(), header, hostParam)
		util.RegisterResponseWithRid(searchRsp, header)
		Expect(rawErr).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		Expect(searchRsp.Data.Count >= 2).To(Equal(true))
		host = searchRsp.Data.Info[0]["host"].(map[string]interface{})
		Expect(commonutil.GetStrByInterface(host[common.BKAgentIDField])).To(Equal(""))
		host = searchRsp.Data.Info[1]["host"].(map[string]interface{})
		Expect(commonutil.GetStrByInterface(host[common.BKAgentIDField])).To(Equal(""))

		By("unbind host agent again")
		err = hostServerClient.UnbindAgent(context.Background(), header, param)
		util.RegisterResponseWithRid(nil, header)
		Expect(err).NotTo(HaveOccurred())
	})
})

var _ = Describe("cloud host test", func() {
	ctx := context.Background()
	var bizID, bizID1, setID, moduleID, hostID1, hostID2, cloudID int64

	It("test preparation", func() {
		By("create cloud area name = 'cloud_host_plat'", func() {
			resp, err := hostServerClient.CreateCloudArea(context.Background(), header, map[string]interface{}{
				"bk_cloud_name":   "cloud_host_plat",
				"bk_status":       "1",
				"bk_cloud_vendor": "1",
				"bk_account_id":   2,
				"creator":         "admin",
			})
			util.RegisterResponseWithRid(resp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Result).To(Equal(true))
			cloudID = int64(resp.Data.Created.ID)
		})

		By("create business bk_biz_name = 'cloud_host_biz'", func() {
			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_name":       "cloud_host_biz",
				"time_zone":         "Africa/Accra",
			}
			rsp, err := apiServerClient.CreateBiz(ctx, "0", header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			bizID, err = commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		By("create set", func() {
			input := mapstr.MapStr{
				"bk_set_name":       "cloud_host_set",
				"bk_parent_id":      bizID,
				"bk_biz_id":         bizID,
				"bk_service_status": "1",
				"bk_set_env":        "3",
			}
			rsp, err := instClient.CreateSet(ctx, bizID, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			var e error
			setID, e = commonutil.GetInt64ByInterface(rsp["bk_set_id"])
			Expect(e).NotTo(HaveOccurred())
		})

		By("create module", func() {
			input := map[string]interface{}{
				"bk_module_name":      "cloud_host_module",
				"bk_parent_id":        setID,
				"service_category_id": 2,
			}
			rsp, err := instClient.CreateModule(ctx, bizID, setID, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			var e error
			moduleID, e = commonutil.GetInt64ByInterface(rsp["bk_module_id"])
			Expect(e).NotTo(HaveOccurred())
		})

		By("create business bk_biz_name = 'trans_cloud_host_biz'", func() {
			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_name":       "trans_cloud_host_biz",
				"time_zone":         "Africa/Accra",
			}
			rsp, err := apiServerClient.CreateBiz(ctx, "0", header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			bizID1, err = commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("add cloud host to biz idle module test", func() {
		By("add cloud hosts to biz idle module", func() {
			input := &metadata.AddCloudHostToBizParam{
				BizID: bizID,
				HostInfo: []mapstr.MapStr{
					{
						common.BKHostInnerIPField: "127.0.0.111",
						common.BKCloudIDField:     cloudID,
						common.BKCloudVendor:      "1",
						common.BKCloudInstIDField: "000000001",
					},
					{
						common.BKHostInnerIPField: "127.0.0.112",
						common.BKCloudIDField:     cloudID,
						common.BKCloudVendor:      "2",
						common.BKCloudInstIDField: "000000002",
					},
				},
			}

			rsp, err := hostServerClient.AddCloudHostToBiz(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(rsp.IDs)).To(Equal(2))
			hostID1, hostID2 = rsp.IDs[0], rsp.IDs[1]
		})

		By("check created cloud hosts", func() {
			rsp, err := hostServerClient.SearchHost(ctx, header, &metadata.HostCommonSearch{
				AppID: bizID,
				Page:  metadata.BasePage{Sort: common.BKHostInnerIPField},
			})
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))

			for _, data := range rsp.Data.Info {
				host, ok := data["host"].(map[string]interface{})
				Expect(ok).To(Equal(true))
				hostID, err := commonutil.GetInt64ByInterface(host[common.BKHostIDField])
				Expect(err).NotTo(HaveOccurred())
				innerIP := commonutil.GetStrByInterface(host[common.BKHostInnerIPField])
				Expect(hostID == hostID1 && innerIP == "127.0.0.111" || hostID == hostID2 && innerIP == "127.0.0.112").
					To(Equal(true))
			}
		})

		By("import cloud host to same biz again", func() {
			input := &metadata.AddCloudHostToBizParam{
				BizID: bizID,
				HostInfo: []mapstr.MapStr{{
					common.BKHostInnerIPField: "127.0.0.111",
					common.BKCloudIDField:     cloudID,
					common.BKHostNameField:    "127.0.0.123",
					common.BKCloudVendor:      "1",
					common.BKCloudInstIDField: "000000001",
				}},
			}

			rsp, err := hostServerClient.AddCloudHostToBiz(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())

			searchRsp, e := hostServerClient.SearchHost(ctx, header, &metadata.HostCommonSearch{AppID: bizID})
			util.RegisterResponseWithRid(rsp, header)
			Expect(e).NotTo(HaveOccurred())
			Expect(searchRsp.Result).To(Equal(true))

			for _, data := range searchRsp.Data.Info {
				host, ok := data["host"].(map[string]interface{})
				Expect(ok).To(Equal(true))
				hostID, err := commonutil.GetInt64ByInterface(host[common.BKHostIDField])
				Expect(err).NotTo(HaveOccurred())
				if hostID == hostID1 {
					Expect(commonutil.GetStrByInterface(host[common.BKHostNameField])).To(Equal("127.0.0.123"))
				}
			}
		})

		By("transfer cloud host to resource pool", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizID1,
				HostIDs:       []int64{hostID1},
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		By("transfer cloud host to another biz", func() {
			input := &metadata.TransferHostAcrossBusinessParameter{
				SrcAppID: bizID,
				DstAppID: bizID1,
				HostID:   []int64{hostID1, hostID2},
			}
			err := hostServerClient.TransferHostAcrossBusiness(ctx, header, input)
			util.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("add cloud hosts to invalid biz", func() {
			input := &metadata.AddCloudHostToBizParam{
				BizID: 1000,
				HostInfo: []mapstr.MapStr{
					{
						common.BKHostInnerIPField: "127.0.0.113",
						common.BKCloudIDField:     cloudID,
					},
				},
			}

			rsp, err := hostServerClient.AddCloudHostToBiz(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		})

		By("add cloud hosts with invalid host", func() {
			input := &metadata.AddCloudHostToBizParam{
				BizID: bizID,
				HostInfo: []mapstr.MapStr{
					{
						common.BKHostInnerIPField: "127.0.0.114",
						common.BKCloudIDField:     cloudID,
						common.BKCloudVendor:      "1",
						common.BKCloudInstIDField: "000000001",
					},
					{
						common.BKHostInnerIPField: "127.0.0.1144",
						common.BKCloudIDField:     cloudID,
					},
				},
			}

			rsp, err := hostServerClient.AddCloudHostToBiz(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		})

		By("add cloud hosts with duplicate host in different biz", func() {
			input := &metadata.AddCloudHostToBizParam{
				BizID: bizID1,
				HostInfo: []mapstr.MapStr{{
					common.BKHostInnerIPField: "127.0.0.111",
					common.BKCloudIDField:     cloudID,
					common.BKCloudVendor:      "1",
					common.BKCloudInstIDField: "000000001",
				}},
			}

			rsp, err := hostServerClient.AddCloudHostToBiz(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		})

		By("add cloud hosts with no biz", func() {
			input := &metadata.AddCloudHostToBizParam{
				HostInfo: []mapstr.MapStr{{
					common.BKHostInnerIPField: "127.0.0.115",
					common.BKCloudIDField:     cloudID,
					common.BKCloudVendor:      "1",
					common.BKCloudInstIDField: "000000001",
				}},
			}

			rsp, err := hostServerClient.AddCloudHostToBiz(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		})

		By("add cloud hosts with no host", func() {
			input := &metadata.AddCloudHostToBizParam{
				BizID: bizID,
			}

			rsp, err := hostServerClient.AddCloudHostToBiz(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		})
	})

	It("delete cloud host from biz idle set test", func() {
		By("delete cloud hosts from invalid biz", func() {
			input := &metadata.DeleteCloudHostFromBizParam{
				BizID:   bizID1,
				HostIDs: []int64{hostID1},
			}

			err := hostServerClient.DeleteCloudHostFromBiz(ctx, header, input)
			util.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("delete cloud hosts with invalid host", func() {
			input := &metadata.DeleteCloudHostFromBizParam{
				BizID:   bizID,
				HostIDs: []int64{hostId},
			}

			err := hostServerClient.DeleteCloudHostFromBiz(ctx, header, input)
			util.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("transfer cloud host to biz module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizID,
				"bk_host_id": []int64{
					hostID1,
				},
				"bk_module_id": []int64{
					moduleID,
				},
			}
			rsp, err := hostServerClient.TransferHostModule(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		By("delete cloud hosts from biz module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizID1,
				HostIDs:       []int64{hostID1, hostID2},
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		By("transfer cloud host to idle set", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizID,
				HostIDs:       []int64{hostID1},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(ctx, header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		By("delete cloud hosts from idle set", func() {
			input := &metadata.DeleteCloudHostFromBizParam{
				BizID:   bizID,
				HostIDs: []int64{hostID1, hostID2},
			}
			err := hostServerClient.DeleteCloudHostFromBiz(ctx, header, input)
			util.RegisterResponseWithRid(err, header)
			Expect(err).NotTo(HaveOccurred())
		})

		By("check deleted cloud hosts", func() {
			rsp, err := hostServerClient.SearchHost(ctx, header, &metadata.HostCommonSearch{AppID: bizID})
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})
	})
})
