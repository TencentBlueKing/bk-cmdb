package host_server_test

import (
	"context"
	"fmt"
	"strconv"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	"configcenter/src/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("host abnormal test", func() {
	var bizId, bizId1, setId, setId1, moduleId, moduleId1, moduleId2, idleModuleId, faultModuleId int64
	var hostId, hostId1, hostId2, hostId3 int64

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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			bizId = int64(rsp.Data["bk_biz_id"].(float64))
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			bizId1 = int64(rsp.Data["bk_biz_id"].(float64))
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
			rsp, err := instClient.CreateSet(context.Background(), strconv.FormatInt(bizId, 10), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			setId = int64(rsp.Data["bk_set_id"].(float64))
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			setId1 = int64(rsp.Data["bk_set_id"].(float64))
		})

		It("create module", func() {
			input := map[string]interface{}{
				"bk_module_name":      "cc_module",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			moduleId = int64(rsp.Data["bk_module_id"].(float64))
		})

		It("create module", func() {
			input := map[string]interface{}{
				"bk_module_name":      "cc_module1",
				"bk_parent_id":        setId,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			moduleId1 = int64(rsp.Data["bk_module_id"].(float64))
		})

		It("create module", func() {
			input := map[string]interface{}{
				"bk_module_name":      "cc_module1",
				"bk_parent_id":        setId1,
				"service_category_id": 2,
				"service_template_id": 0,
			}
			rsp, err := instClient.CreateModule(context.Background(), strconv.FormatInt(bizId, 10), strconv.FormatInt(setId, 10), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			moduleId2 = int64(rsp.Data["bk_module_id"].(float64))
		})

		It("get instance topo", func() {
			rsp, err := instClient.GetInternalModule(context.Background(), "0", strconv.FormatInt(bizId, 10), header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.SetName).To(Equal("空闲机池"))
			Expect(len(rsp.Data.Module)).To(Equal(2))
			idleModuleId = rsp.Data.Module[0].ModuleID
			faultModuleId = rsp.Data.Module[1].ModuleID
		})
	})

	Describe("add host test", func() {
		Describe("add host using api test", func() {
			It("add host using api with large biz_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": 1000,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.1",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with invalid biz_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": "test",
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.2",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with large cloud_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.3",
							"bk_cloud_id":     100,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with invalid cloud_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.4",
							"bk_cloud_id":     -1,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with invalid bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "0.0.0.1",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
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
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with invalid bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.e",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with no host_info", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
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
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api with no bk_cloud_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.5",
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using api to biz", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.6",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("add host using api to biz twice", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.6",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("add host using api to biz multiple ip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.13",
							"bk_cloud_id":     0,
						},
						"5": map[string]interface{}{
							"bk_host_innerip": "2.0.0.14",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
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
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				data := rsp.Data.Info[0]["host"].(map[string]interface{})
				hostId = int64(data["bk_host_id"].(float64))
				data = rsp.Data.Info[1]["host"].(map[string]interface{})
				hostId2 = int64(data["bk_host_id"].(float64))
				Expect(rsp.Data.Count).To(Equal(3))
			})

			It("create object attribute for host", func() {
				input := &metadata.ObjAttDes{
					Attribute: metadata.Attribute{
						OwnerID:       "0",
						ObjectID:      "host",
						PropertyID:    "a",
						PropertyName:  "a",
						PropertyGroup: "default",
						IsEditable:    true,
						PropertyType:  "singleasst",
					},
				}
				rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("update host using created attr", func() {
				input := map[string]interface{}{
					"bk_host_id": fmt.Sprintf("%v", hostId),
					"a":          "2",
				}
				rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("update host using delete attr", func() {
				input := map[string]interface{}{
					"bk_host_id": fmt.Sprintf("%v", hostId),
					"a":          "",
				}
				rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("search biz host", func() {
				input := &params.HostCommonSearch{
					AppID: int(bizId),
					Page: params.PageInfo{
						Sort: "bk_host_id",
					},
				}
				rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				Expect(rsp.Data.Count).To(Equal(3))
			})
		})

		Describe("add host using excel test", func() {
			It("add host using excel with large biz_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": 1000,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.7",
							"bk_cloud_id":     0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with invalid biz_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": "test",
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.8",
							"bk_cloud_id":     0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with large cloud_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId1,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.9",
							"bk_cloud_id":     100,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with invalid cloud_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId1,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.10",
							"bk_cloud_id":     -1,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with invalid bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId1,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "0.0.0.1",
							"bk_cloud_id":     0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with invalid bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId1,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "333.0.0.1",
							"bk_cloud_id":     0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with invalid bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId1,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.e",
							"bk_cloud_id":     0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with no host_info", func() {
				input := map[string]interface{}{
					"bk_biz_id":  bizId1,
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with no bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId1,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_cloud_id": 0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with no bk_cloud_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId1,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.11",
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel to biz", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId1,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.12",
							"bk_cloud_id":     0,
						},
					},
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("add host using excel to biz twice", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId1,
					"host_info": map[string]interface{}{
						"4": map[string]interface{}{
							"bk_host_innerip": "2.0.0.12",
							"bk_cloud_id":     0,
						},
					},
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			})

			It("search biz host", func() {
				input := &params.HostCommonSearch{
					AppID: int(bizId1),
					Page: params.PageInfo{
						Sort: "bk_host_id",
					},
				}
				rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				data := rsp.Data.Info[0]["host"].(map[string]interface{})
				hostId3 = int64(data["bk_host_id"].(float64))
				Expect(rsp.Data.Count).To(Equal(1))
			})
		})
	})

	Describe("search host test", func() {
		It("search host using invalid bk_host_id", func() {
			rsp, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0", "eceer", header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search host using large bk_host_id", func() {
			rsp, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0", "1000", header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search biz host using large biz id", func() {
			input := &params.HostCommonSearch{
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
				Condition: []params.SearchCondition{
					params.SearchCondition{
						ObjectID: "biz",
						Condition: []interface{}{
							map[string]interface{}{
								"field":    "bk_biz_id",
								"operator": "$eq",
								"value":    100,
							},
						},
						Fields: []string{},
					},
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search host using invalid ip", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Ip: params.IPInfo{
					Data:  []string{"1.0.0"},
					Exact: 1,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("search host using multiple ips with an invalid value", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Ip: params.IPInfo{
					Data: []string{
						"2.0.0.6",
						"2.0.0",
						"2.0.0.13",
					},
					Exact: 1,
					Flag:  "bk_host_innerip|bk_host_outerip",
				},
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})
	})

	Describe("transfer host test", func() {
		It("add host to resource", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"4": map[string]interface{}{
						"bk_host_innerip": "2.0.0.15",
						"bk_cloud_id":     0,
					},
				},
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search resource host", func() {
			input := &params.HostCommonSearch{
				AppID: -1,
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			data := rsp.Data.Info[0]["host"].(map[string]interface{})
			hostId1 = int64(data["bk_host_id"].(float64))
			Expect(rsp.Data.Count).To(Equal(1))
		})

		It("transfer resourcehost to nonexist biz's idlemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: 1000,
				HostID: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer biz host to other biz's idlemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostID: []int64{
					hostId,
				},
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer resourcehost to idlemodule less biz id", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				HostID: []int64{
					hostId,
				},
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer resourcehost to idlemodule less host ids", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer resourcehost to idlemodule invalid host ids", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostID: []int64{
					hostId1,
					1000,
				},
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to other biz's module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"bk_host_id": []int64{
					hostId2,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
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
					1000,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer nonexist host to module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_host_id": []int64{
					1000,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer multiple nonexist host to module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_host_id": []int64{
					1000,
					2000,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to module less biz_id", func() {
			input := map[string]interface{}{
				"bk_host_id": []int64{
					hostId3,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to module less host_id", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
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
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to module less is_increment", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_module_id": []int64{
					moduleId,
				},
				"bk_host_id": []int64{
					hostId3,
				},
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer multiple hosts with a nonexist host to module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_host_id": []int64{
					hostId3,
					2000,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search biz host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId1),
				Page: params.PageInfo{
					Sort: "bk_host_id",
				},
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(1))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("transfer resourcehost to idlemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostID: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.AssignHostToApp(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("transfer multiple host to module", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId1,
				"bk_host_id": []int64{
					hostId3,
					hostId1,
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"is_increment": false,
			}
			rsp, err := hostServerClient.HostModuleRelation(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		It("transfer host to idle module nonexist biz", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: 1000,
				HostID: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to idle module one nonexist host", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: 1000,
				HostID: []int64{
					hostId1,
					1000,
				},
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("move nonexist biz's module hosts to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: 1000,
				SetID:         setId,
				ModuleID:      moduleId,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("move nonexist module's hosts to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: bizId1,
				SetID:         setId,
				ModuleID:      1000,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("move nonexist set's module hosts to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: bizId1,
				SetID:         1000,
				ModuleID:      moduleId,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("move nonexist biz's module hosts to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: 1000,
				SetID:         setId,
				ModuleID:      moduleId,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("move unmatching set module hosts to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: bizId1,
				SetID:         setId1,
				ModuleID:      moduleId,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("move unmatching set module hosts to idle", func() {
			input := &metadata.SetHostConfigParams{
				ApplicationID: bizId1,
			}
			rsp, err := hostServerClient.MoveSetHost2IdleModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		It("transfer host to idle module less hostid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to idle module empty hostid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostID:        []int64{},
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to idle module less bizid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				HostID: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to nonexist biz's fault module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: 1000,
				HostID: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer nonexist host to fault module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostID: []int64{
					1000,
				},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer a nonexist host to fault module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostID: []int64{
					hostId1,
					1000,
				},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to fault module less hostid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to fault module less bizid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				HostID: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search fault host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId1),
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("transfer multiple hosts to fault module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostID: []int64{
					hostId1,
					hostId3,
				},
			}
			rsp, err := hostServerClient.MoveHost2FaultModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search fault host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId1),
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId1),
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("transfer unmatching biz host to resourcemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostID: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer nonexist biz host to resourcemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: 100,
				HostID: []int64{
					hostId,
				},
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer a nonexist host to resourcemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId,
				HostID: []int64{
					hostId,
					1000,
				},
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer nonidle host to resourcemodule", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostID: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer host to resourcemodule less hostid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("transfer host to resourcemodule less bizid", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				HostID: []int64{
					hostId1,
				},
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("transfer multiple hosts to idle module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostID: []int64{
					hostId1,
					hostId3,
				},
			}
			rsp, err := hostServerClient.MoveHost2EmptyModule(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search idle host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		It("transfer multiple hosts to resource module", func() {
			input := &metadata.DefaultModuleHostConfigParams{
				ApplicationID: bizId1,
				HostID: []int64{
					hostId1,
					hostId3,
				},
			}
			rsp, err := hostServerClient.MoveHostToResourcePool(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("search resource host", func() {
			input := &params.HostCommonSearch{
				AppID: -1,
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		It("sync host less biz", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "2.0.0.16",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
				},
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host less module", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "2.0.0.17",
						"bk_cloud_id":     0,
					},
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host empty module", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "2.0.0.17",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{},
				"bk_biz_id":    bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host invalid host", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "2.0.0",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host invalid biz", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "2.0.0.18",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"bk_biz_id": 1000,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host invalid module", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "2.0.0.19",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					1000,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host one invalid module", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "2.0.0.20",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
					1000,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host one invalid host", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "2.0.0.21",
						"bk_cloud_id":     0,
					},
					"1": map[string]interface{}{
						"bk_host_innerip": "2.0.0",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host one invalid host", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "2.0.0.21",
						"bk_cloud_id":     0,
					},
					"1": map[string]interface{}{
						"bk_host_innerip": "2.0.0",
						"bk_cloud_id":     0,
					},
				},
				"bk_module_id": []int64{
					moduleId,
				},
				"bk_biz_id": bizId,
			}
			rsp, err := hostServerClient.SyncHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("sync host duplicate module", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "2.0.0.21",
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId1),
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(0))
		})

		It("sync multiple host multiple module", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "2.0.0.22",
						"bk_cloud_id":     0,
					},
					"1": map[string]interface{}{
						"bk_host_innerip": "2.0.0.23",
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("sync host multiple module in different biz", func() {
			input := map[string]interface{}{
				"host_info": map[string]interface{}{
					"0": map[string]interface{}{
						"bk_host_innerip": "2.0.0.30",
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId1),
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})

		It("clone host less biz", func() {
			input := &metadata.CloneHostPropertyParams{
				OrgIP:   "2.0.0.22",
				DstIP:   "2.0.0.24",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("clone host less srcip", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   bizId,
				DstIP:   "2.0.0.25",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("clone host less dstip", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   bizId,
				OrgIP:   "2.0.0.22",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("clone host invalid biz", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   1000,
				OrgIP:   "2.0.0.22",
				DstIP:   "2.0.0.26",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("clone host invalid srcip", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   1000,
				OrgIP:   "2.0.0",
				DstIP:   "2.0.0.27",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("clone host invalid dstip", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   1000,
				OrgIP:   "2.0.0.22",
				DstIP:   "2.0.0",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("clone host exist dstip", func() {
			input := &metadata.CloneHostPropertyParams{
				AppID:   bizId,
				OrgIP:   "2.0.0.22",
				DstIP:   "2.0.0.23",
				CloudID: 0,
			}
			rsp, err := hostServerClient.CloneHostProperty(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("search module host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
				Condition: []params.SearchCondition{
					params.SearchCondition{
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
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
		})
	})

	Describe("batch operate host", func() {
		It("update host less hostid", func() {
			input := map[string]interface{}{
				"bk_host_name": "update_host_name",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update host invalid hostid", func() {
			input := map[string]interface{}{
				"bk_host_id":   "2ew213,fe",
				"bk_host_name": "update_host_name",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update host empty hostid", func() {
			input := map[string]interface{}{
				"bk_host_id":   "",
				"bk_host_name": "update_host_name",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update host one nonexist hostid", func() {
			input := map[string]interface{}{
				"bk_host_id":   fmt.Sprintf("%v,%v", hostId1, 100),
				"bk_host_name": "update_host_name",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update host one nonexist attr", func() {
			input := map[string]interface{}{
				"bk_host_id":      fmt.Sprintf("%v,%v", hostId1, hostId3),
				"bk_host_name":    "update_host_name",
				"fecfecefrrwdxww": "test",
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update host one invalid attr value", func() {
			input := map[string]interface{}{
				"bk_host_id":   fmt.Sprintf("%v,%v", hostId1, hostId3),
				"bk_host_name": 1,
			}
			rsp, err := hostServerClient.UpdateHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("get host base info", func() {
			rsp, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0", strconv.FormatInt(hostId1, 10), header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			for _, data := range rsp.Data {
				if data.PropertyID == "bk_host_name" {
					Expect(data.PropertyValue).To(Equal(""))
					break
				}
			}
		})

		It("delete host less bk_host_id", func() {
			input := map[string]interface{}{
				"bk_supplier_account": "0",
			}
			rsp, err := hostServerClient.DeleteHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("delete host less bk_supplier_account", func() {
			input := map[string]interface{}{
				"bk_host_id": fmt.Sprintf("%v,%v", hostId1, hostId3),
			}
			rsp, err := hostServerClient.DeleteHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("delete host empty bk_host_id", func() {
			input := map[string]interface{}{
				"bk_host_id":          "",
				"bk_supplier_account": "0",
			}
			rsp, err := hostServerClient.DeleteHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("delete host one invalid bk_host_id", func() {
			input := map[string]interface{}{
				"bk_host_id":          fmt.Sprintf("%v,abc", hostId1),
				"bk_supplier_account": "0",
			}
			rsp, err := hostServerClient.DeleteHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("delete host one nonexist bk_host_id", func() {
			input := map[string]interface{}{
				"bk_host_id":          fmt.Sprintf("%v,%v", hostId3, 100),
				"bk_supplier_account": "0",
			}
			rsp, err := hostServerClient.DeleteHostBatch(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("get host base info", func() {
			rsp, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0", strconv.FormatInt(hostId1, 10), header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("get host base info", func() {
			rsp, err := hostServerClient.GetHostInstanceProperties(context.Background(), "0", strconv.FormatInt(hostId3, 10), header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})
	})
})
