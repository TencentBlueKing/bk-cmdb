package topo_server_test

import (
	"context"
	"encoding/json"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("inst test", func() {
	var instId, instId1 int64

	It("create object bk_classification_id = 'bk_middleware' and bk_obj_id='cc_test'", func() {
		test.ClearDatabase()
		input := metadata.Object{
			ObjCls:     "bk_middleware",
			ObjIcon:    "icon-cc-business",
			ObjectID:   "cc_test",
			ObjectName: "cc_test",
			OwnerID:    "0",
			Creator:    "admin",
		}
		rsp, err := objectClient.CreateObject(context.Background(), header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("create object attribute bk_obj_id='cc_test' and bk_property_id='test_sglchar' and bk_property_name='test_sglchar'", func() {
		input := &metadata.ObjAttDes{
			Attribute: metadata.Attribute{
				ObjectID:     "cc_test",
				PropertyID:   "test_sglchar",
				PropertyName: "test_sglchar",
				IsEditable:   false,
				PropertyType: "singlechar",
				Option:       "a+b*",
				IsRequired:   true,
			},
		}
		rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("create object attribute bk_obj_id='cc_test' and bk_property_id='test_unique' and bk_property_name='test_unique'", func() {
		input := &metadata.ObjAttDes{
			Attribute: metadata.Attribute{
				ObjectID:     "cc_test",
				PropertyID:   "test_unique",
				PropertyName: "test_unique",
				IsEditable:   true,
				PropertyType: "singlechar",
				IsRequired:   false,
			},
		}
		rsp, err := apiServerClient.CreateObjectAtt(context.Background(), header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("create inst bk_obj_id='cc_test'", func() {
		input := map[string]interface{}{
			"test_sglchar": "ab",
			"bk_inst_name": "wejeidjew",
			"test_unique":  "1234",
		}
		rsp, err := instClient.CreateInst(context.Background(), "0", "cc_test", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("wejeidjew"))
		Expect(rsp.Data["test_sglchar"].(string)).To(Equal("ab"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("cc_test"))
		Expect(rsp.Data["test_unique"].(string)).To(Equal("1234"))
	})

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
			"bk_asset_id":  "345",
			"bk_inst_name": "~!@#$%^&*()_+-=",
			"bk_sn":        "1234",
		}
		rsp, err := instClient.CreateInst(context.Background(), "0", "bk_switch", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
		Expect(rsp.Code).To(Equal(common.CCErrCommParamsIsInvalid))
	})

	It("create inst bk_obj_id='bk_switch' duplicate bk_asset_id", func() {
		input := map[string]interface{}{
			"bk_asset_id":  "234",
			"bk_inst_name": "abcdefg",
		}
		rsp, err := instClient.CreateInst(context.Background(), "0", "bk_switch", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
		Expect(rsp.Code).To(Equal(common.CCErrCommDuplicateItem))
	})

	It("create inst bk_obj_id='bk_switch' missing bk_inst_name", func() {
		input := map[string]interface{}{
			"bk_asset_id": "456",
		}
		rsp, err := instClient.CreateInst(context.Background(), "0", "bk_switch", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
		Expect(rsp.Code).To(Equal(common.CCErrCommParamsNeedSet))
	})

	It("create inst bk_obj_id='bk_switch' missing bk_asset_id", func() {
		input := map[string]interface{}{
			"bk_inst_name": "456",
		}
		rsp, err := instClient.CreateInst(context.Background(), "0", "bk_switch", header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
		Expect(rsp.Code).To(Equal(common.CCErrCommParamsNeedSet))
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
		Expect(rsp.Code).To(Equal(common.CCErrCommParamsIsInvalid))
	})

	It("update inst invalid instId", func() {
		input := map[string]interface{}{
			"bk_inst_name": "aaa",
		}
		rsp, err := instClient.UpdateInst(context.Background(), "0", "bk_switch", int64(1000), header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
		Expect(rsp.Code).To(Equal(common.CCErrCommNotFound))
	})

	It("update inst with unmatch object", func() {
		input := map[string]interface{}{
			"bk_inst_name": "123",
		}
		rsp, err := instClient.UpdateInst(context.Background(), "0", "cc_test", instId, header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
		Expect(rsp.Code).To(Equal(common.CCErrTopoObjectSelectFailed))
	})

	It("delete inst", func() {
		rsp, err := instClient.DeleteInst(context.Background(), "0", "bk_switch", instId1, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("delete inst with unmatch object", func() {
		rsp, err := instClient.DeleteInst(context.Background(), "0", "cc_test", instId, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
		Expect(rsp.Code).To(Equal(common.CCErrTopoObjectSelectFailed))
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
