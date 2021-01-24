package topo_server_test

import (
	"context"

	"configcenter/src/common/metadata"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("inst test", func() {
	var switchInstId1, switchInstId2, routerInstId1, instAsst1, instAsst2 int64

	It("create inst bk_obj_id='bk_switch'", func() {
		input := map[string]interface{}{
			"bk_asset_id":  "101",
			"bk_inst_name": "switch_1",
			"bk_sn":        "201",
		}
		rsp, err := instClient.CreateInst(context.Background(), "bk_switch", header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("switch_1"))
		Expect(rsp.Data["bk_asset_id"].(string)).To(Equal("101"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("bk_switch"))
		Expect(rsp.Data["bk_sn"].(string)).To(Equal("201"))
		switchInstId1, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create inst bk_obj_id='bk_switch'", func() {
		input := map[string]interface{}{
			"bk_asset_id":  "102",
			"bk_inst_name": "switch_2",
			"bk_sn":        "202",
		}
		rsp, err := instClient.CreateInst(context.Background(), "bk_switch", header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("switch_2"))
		Expect(rsp.Data["bk_asset_id"].(string)).To(Equal("102"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("bk_switch"))
		Expect(rsp.Data["bk_sn"].(string)).To(Equal("202"))
		switchInstId2, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create inst bk_obj_id='bk_router'", func() {
		input := map[string]interface{}{
			"bk_asset_id":  "101",
			"bk_inst_name": "router_1",
			"bk_sn":        "201",
		}
		rsp, err := instClient.CreateInst(context.Background(), "bk_router", header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("router_1"))
		Expect(rsp.Data["bk_asset_id"].(string)).To(Equal("101"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("bk_router"))
		Expect(rsp.Data["bk_sn"].(string)).To(Equal("201"))
		routerInstId1, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create association ='bk_router_default_bk_switch'", func() {
		input := &metadata.Association{
			AsstKindID:           "default",
			AsstObjID:            "bk_switch",
			AssociationName:      "bk_router_default_bk_switch",
			AssociationAliasName: "",
			ObjectID:             "bk_router",
			Mapping:              "n:n",
		}
		rsp, err := asstClient.CreateObject(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("create inst association ='bk_router1_default_bk_switch1'", func() {
		input := &metadata.CreateAssociationInstRequest{
			ObjectAsstID: "bk_router_default_bk_switch",
			InstID:       routerInstId1,
			AsstInstID:   switchInstId1,
		}
		rsp, err := asstClient.CreateInst(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		instAsst1 = rsp.Data.ID
	})

	It("create inst association ='bk_router1_default_bk_switch2'", func() {
		input := &metadata.CreateAssociationInstRequest{
			ObjectAsstID: "bk_router_default_bk_switch",
			InstID:       routerInstId1,
			AsstInstID:   switchInstId2,
		}
		rsp, err := asstClient.CreateInst(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		instAsst2 = rsp.Data.ID
	})
	//check "SearchAssociationRelatedInst" features available.
	It("search inst association related", func() {
		input := &metadata.SearchAssociationRelatedInstRequest{
			Fields: []string{
				"bk_asst_id",
				"bk_inst_id",
				"bk_obj_id",
				"bk_asst_inst_id",
				"bk_asst_obj_id",
				"bk_obj_asst_id",
			},
			Page: metadata.BasePage{
				Start: 0,
				Limit: 10,
			},
			Condition: metadata.SearchAssociationRelatedInstRequestCond{
				ObjectID: "bk_router",
				InstID:   routerInstId1,
			},
		}
		rsp, err := asstClient.SearchAssociationRelatedInst(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(len(rsp.Data)).To(Equal(2))
		Expect(rsp.Data[0].ObjectAsstID).To(Equal("bk_router_default_bk_switch"))
	})
	//check "SearchAssociationRelatedInst" "limit-check<=500" function.
	It("search inst association related", func() {
		input := &metadata.SearchAssociationRelatedInstRequest{
			Fields: []string{
				"bk_asst_id",
				"bk_inst_id",
				"bk_obj_id",
				"bk_asst_inst_id",
				"bk_asst_obj_id",
				"bk_obj_asst_id",
			},
			Page: metadata.BasePage{
				Start: 0,
				Limit: 501,
			},
			Condition: metadata.SearchAssociationRelatedInstRequestCond{
				ObjectID: "bk_router",
				InstID:   routerInstId1,
			},
		}
		rsp, err := asstClient.SearchAssociationRelatedInst(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})
	//check "SearchAssociationRelatedInst" "fields can not be empty." function.
	It("search inst association related", func() {
		input := &metadata.SearchAssociationRelatedInstRequest{
			Fields: []string{},
			Page: metadata.BasePage{
				Start: 0,
				Limit: 501,
			},
			Condition: metadata.SearchAssociationRelatedInstRequestCond{
				ObjectID: "bk_router",
				InstID:   routerInstId1,
			},
		}
		rsp, err := asstClient.SearchAssociationRelatedInst(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})
	//check "DeleteInstBatch" "the number of IDs should be less than 500." function.
	It("delete inst association batch", func() {
		list := make([]int64, 501, 501)
		input := &metadata.DeleteAssociationInstBatchRequest{
			ID: list,
		}
		rsp, err := asstClient.DeleteInstBatch(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})
	//check "DeleteInstBatch" features available.
	It("delete inst association batch", func() {
		input := &metadata.DeleteAssociationInstBatchRequest{
			ID: []int64{instAsst1, instAsst2},
		}
		rsp, err := asstClient.DeleteInstBatch(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(Equal(2))
	})

})
