package host_server_test

import (
	"context"

	params "configcenter/src/common/paraparse"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("host abnormal test", func() {
	var bizId int64

	Describe("test preparation", func() {
		It("create business bk_biz_name = 'Christina'", func() {
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
					"bk_biz_id": bizId,
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
					"bk_biz_id": bizId,
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
					"bk_biz_id": bizId,
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
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with invalid bk_host_innerip", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
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
					"bk_biz_id":  bizId,
					"input_type": "excel",
				}
				rsp, err := hostServerClient.AddHost(context.Background(), header, input)
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
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(false))
			})

			It("add host using excel with no bk_cloud_id", func() {
				input := map[string]interface{}{
					"bk_biz_id": bizId,
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
					"bk_biz_id": bizId,
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
					"bk_biz_id": bizId,
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
					AppID: int(bizId),
					Page: params.PageInfo{
						Sort: "bk_host_id",
					},
				}
				rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
				Expect(rsp.Data.Count).To(Equal(4))
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
	})
})
