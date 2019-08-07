package proc_server_test

import (
	"context"
	"testing"

	"configcenter/src/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var clientSet = test.GetClientSet()
var serviceClient = clientSet.ProcServer().Service()
var instClient = test.GetClientSet().TopoServer().Instance()
var hostServerClient = test.GetClientSet().HostServer()
var apiServerClient = test.GetClientSet().ApiServer()

var bizId, hostId1, hostId2 int64

func TestProcServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProcServer Suite")
}

var _ = BeforeSuite(func() {
	test.ClearDatabase()

	Describe("test preparation", func() {
		Describe("create biz", func() {
			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_name":       "cc_biz",
				"time_zone":         "Africa/Accra",
			}
			rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			bizId = int64(rsp.Data["bk_biz_id"].(float64))
		})

		Describe("add host", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"host_info": map[string]interface{}{
					"1": map[string]interface{}{
						"bk_host_innerip": "1.0.0.1",
						"bk_asset_id":     "addhost_api_asset_1",
						"bk_cloud_id":     0,
					},
					"2": map[string]interface{}{
						"bk_host_innerip": "1.0.0.2",
						"bk_asset_id":     "addhost_api_asset_2",
						"bk_cloud_id":     0,
					},
				},
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		Describe("search host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			hostId1 = int64(rsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_id"].(float64))
			hostId2 = int64(rsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_id"].(float64))
		})
	})
})
