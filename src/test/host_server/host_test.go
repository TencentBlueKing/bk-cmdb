package host_server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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

var _ = Describe("host test", func() {
	var bizId, setId, moduleId, idleModuleId, faultModuleId, dirID int64
	var hostId, hostId1, hostId2, hostId3 int64

	Describe("test preparation", func() {
		It("create business bk_biz_name = 'cc_biz'", func() {
			test.ClearDatabase()

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

		It("create resource pool directory", func() {
			dir := map[string]interface{}{
				"bk_module_name": "test",
			}
			rsp, err := dirClient.CreateResourceDirectory(context.Background(), header, dir)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search host created using api", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Ip: params.IPInfo{
					Data:  []string{"127.0.0.1"},
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search host created using excel", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Ip: params.IPInfo{
					Data:  []string{"127.0.0.2"},
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
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.2"))
			Expect(data["bk_asset_id"].(string)).To(Equal("addhost_excel_asset_1"))
			Expect(data["bk_host_name"].(string)).To(Equal("127.0.0.2"))
			hostId2, err = commonutil.GetInt64ByInterface(data["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("search host using multiple ips", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Ip: params.IPInfo{
					Data: []string{
						"127.0.0.1",
						"127.0.0.2",
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
		// 	util.RegisterResponse(rsp)
		// 	Expect(err).NotTo(HaveOccurred())
		// 	Expect(rsp.Result).To(Equal(true))
		// })

		// It("search host created using agent", func() {
		// 	input := &params.HostCommonSearch{
		// 		AppID: int(bizId),
		// 		Ip: params.IPInfo{
		// 			Data:  []string{"127.0.0.3"},
		// 			Exact: 1,
		// 			Flag:  "bk_host_innerip|bk_host_outerip",
		// 		},
		// 		Page: params.PageInfo{
		// 			Sort: "bk_host_id",
		// 		},
		// 	}
		// 	rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
		// 	util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search resource host", func() {
			input := &params.HostCommonSearch{
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.4"))
			hostId, err = commonutil.GetInt64ByInterface(data["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		It("get host base info", func() {
			rsp, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0", strconv.FormatInt(hostId, 10), header)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(len(rsp.Data)).To(Equal(2))
		})
	})

	Describe("transfer host test", func() {
		It("search biz host", func() {
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search biz host", func() {
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search biz host", func() {
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
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
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
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search cloned host", func() {
			input := &params.HostCommonSearch{
				AppID: -1,
				Ip: params.IPInfo{
					Data:  []string{"127.0.0.5"},
					Exact: 0,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
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
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.5"))
			Expect(data["bk_comment"].(string)).To(Equal("127.0.0.1 comment"))

		})

		It("get instance topo", func() {
			rsp, err := instClient.GetInternalModule(context.Background(), "0", strconv.FormatInt(bizId, 10), header)
			util.RegisterResponse(rsp)
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
			input := &params.HostCommonSearch{
				AppID: int(bizId),
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search fault host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search fault host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.2"))
		})

		It("search transfered module host", func() {
			input := &params.HostCommonSearch{
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
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
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			data1 := rsp.Data.Info[1]["host"].(map[string]interface{})
			Expect("127.0.0.2").To(SatisfyAny(Equal(data["bk_host_innerip"].(string)), Equal(data1["bk_host_innerip"].(string))))
		})

		It("search fault host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search resource host", func() {
			input := &params.HostCommonSearch{
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search resource host", func() {
			input := &params.HostCommonSearch{
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(3))
		})

		It("search resource host change start limit", func() {
			input := &params.HostCommonSearch{
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
					Sort:  "bk_host_id",
					Start: 2,
					Limit: 2,
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(3))
			Expect(len(rsp.Data.Info)).To(Equal(1))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			Expect(data["bk_host_innerip"].(string)).To(Equal("127.0.0.5"))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
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
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true), rsp.ToString())
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
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
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
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
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search resource host", func() {
			input := &params.HostCommonSearch{
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
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
		})
	})
})

var _ = Describe("list_hosts_topo test", func() {
	It("list_hosts_topo", func() {
		test.ClearDatabase()

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
		setRsp, err := instClient.CreateSet(context.Background(), strconv.FormatInt(bizId, 10), header, setInput)
		util.RegisterResponse(setRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(setRsp.Result).To(Equal(true))
		setId, err := commonutil.GetInt64ByInterface(setRsp.Data[common.BKSetIDField])
		Expect(err).NotTo(HaveOccurred())

		By("create module cc_module")
		moduleInput := map[string]interface{}{
			"bk_module_name":      "cc_module",
			"bk_parent_id":        setId,
			"service_category_id": 2,
			"service_template_id": 0,
		}
		moduleRsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, moduleInput)
		util.RegisterResponse(moduleRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(moduleRsp.Result).To(Equal(true))
		moduleId1, err := commonutil.GetInt64ByInterface(moduleRsp.Data[common.BKModuleIDField])
		Expect(err).NotTo(HaveOccurred())

		By("create module cc_module1")
		moduleInput1 := map[string]interface{}{
			"bk_module_name":      "cc_module1",
			"bk_parent_id":        setId,
			"service_category_id": 2,
			"service_template_id": 0,
		}
		moduleRsp1, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, moduleInput1)
		util.RegisterResponse(moduleRsp1)
		Expect(err).NotTo(HaveOccurred())
		Expect(moduleRsp1.Result).To(Equal(true))
		moduleId2, err := commonutil.GetInt64ByInterface(moduleRsp1.Data[common.BKModuleIDField])
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
		util.RegisterResponse(hostRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(true), hostRsp.ToString())

		By("search hosts")
		searchInput := &params.HostCommonSearch{
			AppID: int(bizId),
			Page: params.PageInfo{
				Sort: common.BKHostIDField,
			},
		}
		searchRsp, err := hostServerClient.SearchHost(context.Background(), header, searchInput)
		util.RegisterResponse(searchRsp)
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
		rsp, err := hostServerClient.ListBizHostsTopo(context.Background(), header, bizId, &metadata.ListHostsWithNoBizParameter{Page: metadata.BasePage{Sort: common.BKHostIDField, Limit: 10}, Fields: []string{"bk_host_id"}})
		util.RegisterResponse(rsp)
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
		test.ClearDatabase()

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
		util.RegisterResponse(hostRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(true), hostRsp.ToString())

		By("search hosts")
		searchInput := &params.HostCommonSearch{
			Page: params.PageInfo{
				Sort: common.BKHostIDField,
			},
		}
		searchRsp, err := hostServerClient.SearchHost(context.Background(), header, searchInput)
		util.RegisterResponse(searchRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		hostId1, err := commonutil.GetInt64ByInterface(searchRsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_id"])
		Expect(err).NotTo(HaveOccurred())
		hostId2, err := commonutil.GetInt64ByInterface(searchRsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_id"])
		Expect(err).NotTo(HaveOccurred())

		By("update host property batch, bk_host_name is not editable field")
		updateInput := map[string]interface{}{
			"update": []map[string]interface{}{
				{
					common.BKHostIDField: hostId1,
					"properties": map[string]interface{}{
						"bk_host_name": "batch_update1",
						"operator":     "admin",
						"bk_comment":   "test",
						"bk_isp_name":  "1",
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
		updateRsp, err := hostServerClient.UpdateHostPropertyBatch(context.Background(), header, updateInput)
		util.RegisterResponse(updateRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(updateRsp.Result).To(Equal(true))

		By("search updated host property")
		searchRsp, err = hostServerClient.SearchHost(context.Background(), header, searchInput)
		util.RegisterResponse(searchRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		Expect(searchRsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_name"].(string)).To(Equal(""))
		Expect(searchRsp.Data.Info[0]["host"].(map[string]interface{})["operator"].(string)).To(Equal("admin"))
		Expect(searchRsp.Data.Info[0]["host"].(map[string]interface{})["bk_comment"].(string)).To(Equal("test"))
		Expect(searchRsp.Data.Info[0]["host"].(map[string]interface{})["bk_isp_name"].(string)).To(Equal("1"))
		Expect(searchRsp.Data.Info[1]["host"].(map[string]interface{})["bk_bak_operator"].(string)).To(Equal("admin"))
		Expect(searchRsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_outerip"].(string)).To(Equal("127.2.3.4"))
	})
})

var _ = Describe("multiple ip host validation test", func() {
	It("multiple ip host validation", func() {
		test.ClearDatabase()

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
		util.RegisterResponse(hostRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(true), hostRsp.ToString())

		By("search hosts")
		searchInput := &params.HostCommonSearch{
			Page: params.PageInfo{
				Sort: common.BKHostIDField,
			},
		}
		searchRsp, err := hostServerClient.SearchHost(context.Background(), header, searchInput)
		util.RegisterResponse(searchRsp)
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
		addHostResult, err := test.GetClientSet().CoreService().Instance().CreateInstance(context.Background(), header, common.BKInnerObjIDHost, input)
		util.RegisterResponse(addHostResult)
		Expect(err).NotTo(HaveOccurred())
		Expect(addHostResult.Result).To(Equal(false), addHostResult.ToString())

		By("add hosts with one same ip using api")
		input = &metadata.CreateModelInstance{
			Data: map[string]interface{}{
				"bk_host_innerip": "1.0.0.1",
				"bk_cloud_id":     0,
			},
		}
		addHostResult, err = test.GetClientSet().CoreService().Instance().CreateInstance(context.Background(), header, common.BKInnerObjIDHost, input)
		util.RegisterResponse(addHostResult)
		Expect(err).NotTo(HaveOccurred())
		Expect(addHostResult.Result).To(Equal(false), addHostResult.ToString())
	})
})

var _ = Describe("add_host_to_resource_pool test", func() {
	It("add_host_to_resource_pool", func() {
		test.ClearDatabase()

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
		util.RegisterResponse(hostRsp)
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
		searchRsp, err := hostServerClient.SearchHost(context.Background(), header, &params.HostCommonSearch{})
		util.RegisterResponse(searchRsp)
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
		util.RegisterResponse(hostRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(false))

		By("search hosts")
		searchRsp, err = hostServerClient.SearchHost(context.Background(), header, &params.HostCommonSearch{})
		util.RegisterResponse(searchRsp)
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
		util.RegisterResponse(hostRsp)
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
		searchRsp, err = hostServerClient.SearchHost(context.Background(), header, &params.HostCommonSearch{
			Page: params.PageInfo{
				Sort: common.BKHostIDField,
			},
		})
		util.RegisterResponse(searchRsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Result).To(Equal(true))
		Expect(searchRsp.Data.Count).To(Equal(3))
		hostID, err := commonutil.GetInt64ByInterface(searchRsp.Data.Info[2]["host"].(map[string]interface{})[common.BKHostIDField])
		Expect(err).NotTo(HaveOccurred())
		Expect(hostID).To(Equal(result.Success[0].HostID))
	})
})
