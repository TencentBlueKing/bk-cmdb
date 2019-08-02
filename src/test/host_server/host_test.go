package host_server_test

import (
	params "configcenter/src/common/paraparse"
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("host test", func() {
	var bizId int64

	It("create business bk_biz_name = 'cc_biz'", func() {
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
		Expect(rsp.Data).To(ContainElement("cc_biz"))
		bizId = int64(rsp.Data["bk_biz_id"].(float64))
	})

	It("add host using api", func() {
		input := map[string]interface{}{
			"bk_biz_id": bizId,
			"host_info": map[string]interface{}{
				"4": map[string]interface{}{
					"bk_host_innerip": "1.2.3.4",
					"bk_asset_id":     "addhost_api_asset_1",
					"bk_cloud_id":     0,
				},
			},
		}
		rsp, err := hostServerClient.AddHost(context.Background(), header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("search host created using api", func() {
		input := &params.HostCommonSearch{
			AppID: int(bizId),
			Ip: params.IPInfo{
				Data:  []string{"1.2.3.4"},
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
		Expect(rsp.Data.Count).To(Equal(1))
		data := rsp.Data.Info[0]["host"].(map[string]interface{})
		Expect(data["bk_host_innerip"].(string)).To(Equal("1.2.3.4"))
		Expect(data["bk_asset_id"].(string)).To(Equal("addhost_api_asset_1"))
	})

})
