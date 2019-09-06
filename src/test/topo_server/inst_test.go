package topo_server_test

import (
	"context"
	"encoding/json"
	"strconv"

	"configcenter/src/common/metadata"
	"configcenter/src/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("inst test", func() {
	var instId, instId1 int64
	instClient := topoServerClient.Instance()

	It("create inst bk_obj_id='bk_switch'", func() {
		test.ClearDatabase()
		input := map[string]interface{}{
			"bk_asset_id":  "123",
			"bk_inst_name": "wejeidjew",
			"bk_sn":        "1234",
		}
		rsp, err := instClient.CreateInst(context.Background(), "0", "bk_switch", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("wejeidjew"))
		Expect(rsp.Data["bk_asset_id"].(string)).To(Equal("123"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("bk_switch"))
		Expect(rsp.Data["bk_sn"].(string)).To(Equal("1234"))
		instId = int64(rsp.Data["bk_inst_id"].(float64))
	})

	It("create inst bk_obj_id='bk_switch'", func() {
		input := map[string]interface{}{
			"bk_asset_id":  "234",
			"bk_inst_name": "wejeidjew",
			"bk_sn":        "1234",
		}
		rsp, err := instClient.CreateInst(context.Background(), "0", "bk_switch", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("wejeidjew"))
		Expect(rsp.Data["bk_asset_id"].(string)).To(Equal("234"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("bk_switch"))
		Expect(rsp.Data["bk_sn"].(string)).To(Equal("1234"))
		instId1 = int64(rsp.Data["bk_inst_id"].(float64))
	})

	It("create inst invalid bk_inst_name", func() {
		input := map[string]interface{}{
			"bk_asset_id":  "234",
			"bk_inst_name": "~!@#$%^&*()_+-=",
			"bk_sn":        "1234",
		}
		rsp, err := instClient.CreateInst(context.Background(), "0", "bk_switch", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})

	It("update inst", func() {
		input := map[string]interface{}{
			"bk_inst_name": "aaa",
		}
		rsp, err := instClient.UpdateInst(context.Background(), "0", "bk_switch", instId, header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("update inst invalid bk_inst_name", func() {
		input := map[string]interface{}{
			"bk_inst_name": "~!@#$%^&*()_+-=",
		}
		rsp, err := instClient.UpdateInst(context.Background(), "0", "bk_switch", instId, header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})

	It("delete inst", func() {
		rsp, err := instClient.DeleteInst(context.Background(), "0", "bk_switch", instId1, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("search inst", func() {
		input := &metadata.SearchParams{
			Condition: map[string]interface{}{},
			Page: map[string]interface{}{
				"sort": "id",
			},
		}
		rsp, err := instClient.SelectInsts(context.Background(), "0", "bk_switch", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data.Count).To(Equal(1))
		Expect(rsp.Data.Info[0]["bk_inst_name"].(string)).To(Equal("aaa"))
		Expect(rsp.Data.Info[0]["bk_asset_id"].(string)).To(Equal("123"))
		Expect(rsp.Data.Info[0]["bk_obj_id"].(string)).To(Equal("bk_switch"))
		Expect(rsp.Data.Info[0]["bk_sn"].(string)).To(Equal("1234"))
	})

	It("search inst by object", func() {
		input := &metadata.SearchParams{
			Condition: map[string]interface{}{},
			Page: map[string]interface{}{
				"sort": "id",
			},
		}
		rsp, err := instClient.InstSearch(context.Background(), "0", "bk_switch", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data.Count).To(Equal(1))
		Expect(rsp.Data.Info[0]["bk_inst_name"].(string)).To(Equal("aaa"))
		Expect(rsp.Data.Info[0]["bk_asset_id"].(string)).To(Equal("123"))
		Expect(rsp.Data.Info[0]["bk_obj_id"].(string)).To(Equal("bk_switch"))
		Expect(rsp.Data.Info[0]["bk_sn"].(string)).To(Equal("1234"))
	})

	It("search inst association detail", func() {
		input := &metadata.SearchParams{
			Condition: map[string]interface{}{},
			Page: map[string]interface{}{
				"sort": "id",
			},
		}
		rsp, err := instClient.SelectInst(context.Background(), "0", "bk_switch", strconv.FormatInt(instId, 10), header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data.Count).To(Equal(1))
		Expect(rsp.Data.Info[0]["bk_inst_name"].(string)).To(Equal("aaa"))
		Expect(rsp.Data.Info[0]["bk_asset_id"].(string)).To(Equal("123"))
		Expect(rsp.Data.Info[0]["bk_obj_id"].(string)).To(Equal("bk_switch"))
		Expect(rsp.Data.Info[0]["bk_sn"].(string)).To(Equal("1234"))
	})

	It("search inst by association", func() {
		input := &metadata.AssociationParams{
			Page: metadata.BasePage{
				Sort: "id",
			},
		}
		rsp, err := instClient.SelectInstsByAssociation(context.Background(), "0", "bk_switch", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data.Count).To(Equal(1))
		Expect(rsp.Data.Info[0]["bk_inst_name"].(string)).To(Equal("aaa"))
		Expect(rsp.Data.Info[0]["bk_asset_id"].(string)).To(Equal("123"))
		Expect(rsp.Data.Info[0]["bk_obj_id"].(string)).To(Equal("bk_switch"))
		Expect(rsp.Data.Info[0]["bk_sn"].(string)).To(Equal("1234"))
	})

	It("search inst topo", func() {
		input := &metadata.SearchParams{
			Condition: map[string]interface{}{},
			Page: map[string]interface{}{
				"sort": "id",
			},
		}
		rsp, err := instClient.SelectTopo(context.Background(), "0", "bk_switch", strconv.FormatInt(instId, 10), header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("search inst association topo", func() {
		input := &metadata.SearchParams{
			Condition: map[string]interface{}{},
			Page: map[string]interface{}{
				"sort": "id",
			},
		}
		rsp, err := instClient.SelectAssociationTopo(context.Background(), "0", "bk_switch", strconv.FormatInt(instId, 10), header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		j, err := json.Marshal(rsp.Data[0].Curr)
		data := map[string]interface{}{}
		json.Unmarshal(j, &data)
		Expect(data["bk_obj_id"].(string)).To(Equal("bk_switch"))
		Expect(data["bk_inst_name"].(string)).To(Equal("aaa"))
		Expect(int64(data["bk_inst_id"].(float64))).To(Equal(instId))
	})
})
