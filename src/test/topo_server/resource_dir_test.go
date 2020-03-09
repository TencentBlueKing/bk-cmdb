package topo_server_test

import (
	"configcenter/src/common/metadata"
	"configcenter/src/test"
	"configcenter/src/test/util"
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("resource pool directory test", func() {

	It("create resource pool directory", func() {
		test.ClearDatabase()
		input := map[string]interface{}{
			"life_cycle":        "2",
			"language":          "1",
			"bk_biz_maintainer": "admin",
			"bk_biz_productor":  "",
			"bk_biz_tester":     "",
			"bk_biz_developer":  "",
			"operator":          "",
			"bk_biz_name":       "eereeede",
			"time_zone":         "Africa/Accra",
		}
		rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(ContainElement("eereeede"))
	})
})

func prepareHost() {
	prepareHost := map[string]interface{}{
		"bk_biz_id":  1,
		"input_type": "excel",
		"host_info": map[string]interface{}{
			"1": map[string]interface{}{
				"bk_host_innerip": "1.0.0.1",
				"bk_asset_id":     "addhost_excel_asset_1",
				"bk_host_name":    "1.value1.0.1",
			},
			"2": map[string]interface{}{
				"bk_host_innerip": "1.0.0.2",
				"bk_asset_id":     "addhost_excel_asset_1",
				"bk_host_name":    "1.value1.0.2",
			},
			"3": map[string]interface{}{
				"bk_host_innerip": "1.0.0.3",
				"bk_asset_id":     "addhost_excel_asset_1",
				"bk_host_name":    "1.value1.0.3",
			},
		},
	}

	rsp := metadata.Response{}
	err := apiServerClient.Client().Post().
		WithContext(context.Background()).
		Body(prepareHost).
		SubResourcef("/hosts/add").
		WithHeaders(header).
		Do().Into(&rsp)
	util.RegisterResponse(rsp)
	Expect(err).NotTo(HaveOccurred())
	Expect(rsp.Result).To(Equal(true), rsp.ToString())
}
