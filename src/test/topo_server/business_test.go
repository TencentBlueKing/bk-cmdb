package topo_server_test

import (
	"configcenter/src/common"
	params "configcenter/src/common/paraparse"
	"context"
	"fmt"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("business test", func() {

	var bizId, bizId2 string

	It("create business bk_biz_name = 'eereeede'", func() {
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
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(ContainElement("eereeede"))
		bizId = strconv.FormatInt(int64(rsp.Data["bk_biz_id"].(float64)), 10)
	})

	It("create business bk_biz_name = 'mmrmm'", func() {
		input := map[string]interface{}{
			"life_cycle":        "2",
			"language":          "1",
			"bk_biz_maintainer": "admin",
			"bk_biz_productor":  "",
			"bk_biz_tester":     "",
			"bk_biz_developer":  "",
			"operator":          "",
			"bk_biz_name":       "mmrmm",
			"time_zone":         "Africa/Accra",
		}
		rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(ContainElement("mmrmm"))
		bizId2 = strconv.FormatInt(int64(rsp.Data["bk_biz_id"].(float64)), 10)
	})

	It(fmt.Sprintf("update business bk_biz_id = %s", bizId), func() {
		input := map[string]interface{}{
			"bk_biz_name": "cdewdercfee",
			"life_cycle":  "2",
		}
		rsp, err := apiServerClient.UpdateBiz(context.Background(), "0", bizId, header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It(fmt.Sprintf("update business enable status bk_biz_id = %s", bizId2), func() {
		rsp, err := apiServerClient.UpdateBizDataStatus(context.Background(), "0", common.DataStatusDisabled, bizId2, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("search business", func() {
		input := &params.SearchParams{
			Page: map[string]interface{}{
				"start": 0,
				"limit": 10,
				"sort":  "",
			},
			Fields:    []string{},
			Condition: map[string]interface{}{},
		}
		rsp, err := apiServerClient.SearchBiz(context.Background(), "0", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data.Info).To(ContainElement(ContainElement("cdewdercfee")))
		Expect(rsp.Data.Info).NotTo(ContainElement(ContainElement("eereeede")))
		Expect(rsp.Data.Info).NotTo(ContainElement(ContainElement("mmrmm")))
	})
})
