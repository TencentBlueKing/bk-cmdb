package topo_server_test

import (
	"context"
	"encoding/json"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("inst test", func() {
	var switchInstId1, switchInstId2, switchInstId3, switchInstId4, switchInstId5, routerInstId1, routerInstId2,
		routerInstId3, instAsst1, instAsst2, instAsst3, instAsst4, instAsst5 int64

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

	It("create inst bk_obj_id='bk_switch'", func() {
		input := map[string]interface{}{
			"bk_asset_id":  "103",
			"bk_inst_name": "switch_3",
			"bk_sn":        "203",
		}
		rsp, err := instClient.CreateInst(context.Background(), "bk_switch", header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("switch_3"))
		Expect(rsp.Data["bk_asset_id"].(string)).To(Equal("103"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("bk_switch"))
		Expect(rsp.Data["bk_sn"].(string)).To(Equal("203"))
		switchInstId3, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create inst bk_obj_id='bk_switch'", func() {
		input := map[string]interface{}{
			"bk_asset_id":  "104",
			"bk_inst_name": "switch_4",
			"bk_sn":        "204",
		}
		rsp, err := instClient.CreateInst(context.Background(), "bk_switch", header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("switch_4"))
		Expect(rsp.Data["bk_asset_id"].(string)).To(Equal("104"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("bk_switch"))
		Expect(rsp.Data["bk_sn"].(string)).To(Equal("204"))
		switchInstId4, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create inst bk_obj_id='bk_switch'", func() {
		input := map[string]interface{}{
			"bk_asset_id":  "105",
			"bk_inst_name": "switch_5",
			"bk_sn":        "205",
		}
		rsp, err := instClient.CreateInst(context.Background(), "bk_switch", header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("switch_5"))
		Expect(rsp.Data["bk_asset_id"].(string)).To(Equal("105"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("bk_switch"))
		Expect(rsp.Data["bk_sn"].(string)).To(Equal("205"))
		switchInstId5, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
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

	It("create inst bk_obj_id='bk_router'", func() {
		input := map[string]interface{}{
			"bk_asset_id":  "102",
			"bk_inst_name": "router_2",
			"bk_sn":        "202",
		}
		rsp, err := instClient.CreateInst(context.Background(), "bk_router", header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("router_2"))
		Expect(rsp.Data["bk_asset_id"].(string)).To(Equal("102"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("bk_router"))
		Expect(rsp.Data["bk_sn"].(string)).To(Equal("202"))
		routerInstId2, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create inst bk_obj_id='bk_router'", func() {
		input := map[string]interface{}{
			"bk_asset_id":  "103",
			"bk_inst_name": "router_3",
			"bk_sn":        "203",
		}
		rsp, err := instClient.CreateInst(context.Background(), "bk_router", header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("router_3"))
		Expect(rsp.Data["bk_asset_id"].(string)).To(Equal("103"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("bk_router"))
		Expect(rsp.Data["bk_sn"].(string)).To(Equal("203"))
		routerInstId3, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
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

	It("create association ='bk_router_belong_bk_switch'", func() {
		input := &metadata.Association{
			AsstKindID:           "belong",
			AsstObjID:            "bk_switch",
			AssociationName:      "bk_router_belong_bk_switch",
			AssociationAliasName: "",
			ObjectID:             "bk_router",
			Mapping:              "1:1",
		}
		rsp, err := asstClient.CreateObject(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("create association ='bk_router_connect_bk_switch'", func() {
		input := &metadata.Association{
			AsstKindID:           "connect",
			AsstObjID:            "bk_switch",
			AssociationName:      "bk_router_connect_bk_switch",
			AssociationAliasName: "",
			ObjectID:             "bk_router",
			Mapping:              "1:n",
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

	It("create inst association ='bk_router_belong_bk_switch', belong mapping 1:1", func() {
		input := &metadata.CreateManyInstAsstRequest{
			ObjectAsstID: "bk_router_belong_bk_switch",
			ObjectID:     "bk_router",
			AsstObjectID: "bk_switch",
			Details: []metadata.InstAsst{
				{
					InstID:     routerInstId2,
					AsstInstID: switchInstId3,
				},
				{
					InstID:     routerInstId2,
					AsstInstID: switchInstId4,
				},
				{
					InstID:     routerInstId2,
					AsstInstID: switchInstId5,
				},
			},
		}
		rsp, err := asstClient.CreateManyInstAssociation(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(len(rsp.Data.SuccessCreated)).To(Equal(1))
		Expect(len(rsp.Data.Error)).To(Equal(2))
		instAsst3 = rsp.Data.SuccessCreated[0]
	})

	It("create inst association ='bk_router_connect_bk_switch', connect mapping 1:n", func() {
		input := &metadata.CreateManyInstAsstRequest{
			ObjectAsstID: "bk_router_connect_bk_switch",
			ObjectID:     "bk_router",
			AsstObjectID: "bk_switch",
			Details: []metadata.InstAsst{
				{
					InstID:     routerInstId3,
					AsstInstID: switchInstId4,
				},
				{
					InstID:     routerInstId3,
					AsstInstID: switchInstId5,
				},
			},
		}
		rsp, err := asstClient.CreateManyInstAssociation(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(len(rsp.Data.SuccessCreated)).To(Equal(2))
		Expect(len(rsp.Data.Error)).To(Equal(0))
		instAsst4 = rsp.Data.SuccessCreated[0]
		instAsst5 = rsp.Data.SuccessCreated[1]
	})

	It("create inst association ='bk_router_belong_bk_switch', belong mapping 1:1", func() {
		input := &metadata.CreateManyInstAsstRequest{
			ObjectAsstID: "bk_router_belong_bk_switch",
			ObjectID:     "bk_router",
			AsstObjectID: "bk_router",
			Details: []metadata.InstAsst{
				{
					InstID:     routerInstId2,
					AsstInstID: switchInstId3,
				},
				{
					InstID:     routerInstId2,
					AsstInstID: switchInstId4,
				},
				{
					InstID:     routerInstId2,
					AsstInstID: switchInstId5,
				},
			},
		}
		rsp, err := asstClient.CreateManyInstAssociation(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})

	// check "SearchAssociationRelatedInst" features available.
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
	// check "SearchAssociationRelatedInst" "limit-check<=500" function.
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
	// check "SearchAssociationRelatedInst" "fields can not be empty." function.
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

	It("search instance associations and instances detail", func() {
		input := &metadata.InstAndAssocRequest{}
		input.Condition.AsstFilter = &querybuilder.QueryFilter{Rule: querybuilder.CombinedRule{
			Condition: querybuilder.ConditionAnd,
			Rules: []querybuilder.Rule{
				&querybuilder.AtomRule{
					Field: common.BKObjIDField, Operator: querybuilder.OperatorEqual, Value: "bk_router",
				},
				&querybuilder.AtomRule{
					Field: common.BKInstIDField, Operator: querybuilder.OperatorEqual, Value: routerInstId1,
				},
			},
		}}
		input.Condition.SrcDetail = true
		input.Page.Limit = 200

		rsp, err := asstClient.SearchInstAssocAndInstDetail(context.Background(), header, "bk_router", input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))

		Expect(len(rsp.Data.Asst)).To(Equal(2))
		Expect(len(rsp.Data.Src)).To(Equal(1))
	})

	It("search instance associations and instances detail too many limit", func() {
		input := &metadata.InstAndAssocRequest{}
		input.Condition.AsstFilter = &querybuilder.QueryFilter{Rule: querybuilder.CombinedRule{
			Condition: querybuilder.ConditionAnd,
			Rules: []querybuilder.Rule{
				&querybuilder.AtomRule{
					Field: common.BKObjIDField, Operator: querybuilder.OperatorEqual, Value: "bk_router",
				},
				&querybuilder.AtomRule{
					Field: common.BKInstIDField, Operator: querybuilder.OperatorEqual, Value: routerInstId1,
				},
			},
		}}
		input.Condition.SrcDetail = true
		input.Page.Limit = 201

		rsp, err := asstClient.SearchInstAssocAndInstDetail(context.Background(), header, "bk_router", input)
		util.RegisterResponse(rsp)
		Expect(err).To(HaveOccurred())
	})

	It("search instance associations", func() {
		input := &metadata.CommonSearchFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd,
					Rules: []querybuilder.Rule{
						&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual, Value: routerInstId1},
					},
				},
			},
			Fields: []string{"bk_asst_id", "bk_inst_id", "bk_obj_id", "bk_asst_inst_id", "bk_asst_obj_id", "bk_obj_asst_id"},
			Page:   metadata.BasePage{Start: 0, Limit: 1},
		}

		rsp, err := asstClient.SearchInstanceAssociations(context.Background(), header, "bk_router", input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))

		data, err := mapstr.NewFromInterface(rsp.Data)
		Expect(err).NotTo(HaveOccurred())

		info, ok := data["info"].([]interface{})
		Expect(ok).To(Equal(true))
		Expect(len(info)).To(Equal(1))

		association, ok := info[0].(map[string]interface{})
		Expect(ok).To(Equal(true))

		Expect(association["bk_obj_asst_id"].(string)).To(Equal("bk_router_default_bk_switch"))
	})

	It("search instance associations without conditions", func() {
		input := &metadata.CommonSearchFilter{
			Fields: []string{"bk_asst_id", "bk_inst_id", "bk_obj_id", "bk_asst_inst_id", "bk_asst_obj_id", "bk_obj_asst_id"},
			Page:   metadata.BasePage{Start: 0, Limit: 1},
		}

		rsp, err := asstClient.SearchInstanceAssociations(context.Background(), header, "bk_router", input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))

		data, err := mapstr.NewFromInterface(rsp.Data)
		Expect(err).NotTo(HaveOccurred())

		info, ok := data["info"].([]interface{})
		Expect(ok).To(Equal(true))
		Expect(len(info)).To(Equal(1))

		association, ok := info[0].(map[string]interface{})
		Expect(ok).To(Equal(true))

		Expect(association["bk_obj_asst_id"].(string)).To(Equal("bk_router_default_bk_switch"))
	})

	It("search instance associations without fields", func() {
		input := &metadata.CommonSearchFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd,
					Rules: []querybuilder.Rule{
						&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual, Value: routerInstId1},
					},
				},
			},
			Page: metadata.BasePage{Start: 0, Limit: 1},
		}

		rsp, err := asstClient.SearchInstanceAssociations(context.Background(), header, "bk_router", input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))

		data, err := mapstr.NewFromInterface(rsp.Data)
		Expect(err).NotTo(HaveOccurred())

		info, ok := data["info"].([]interface{})
		Expect(ok).To(Equal(true))
		Expect(len(info)).To(Equal(1))

		association, ok := info[0].(map[string]interface{})
		Expect(ok).To(Equal(true))
		Expect(len(association)).To(Equal(8))
	})

	It("search instance associations without page", func() {
		input := &metadata.CommonSearchFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd,
					Rules: []querybuilder.Rule{
						&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual, Value: routerInstId1},
					},
				},
			},
			Fields: []string{"bk_asst_id", "bk_inst_id", "bk_obj_id", "bk_asst_inst_id", "bk_asst_obj_id", "bk_obj_asst_id"},
		}

		rsp, err := asstClient.SearchInstanceAssociations(context.Background(), header, "bk_router", input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})

	It("search instance associations with limit more than 500", func() {
		input := &metadata.CommonSearchFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd,
					Rules: []querybuilder.Rule{
						&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual, Value: routerInstId1},
					},
				},
			},
			Fields: []string{"bk_asst_id", "bk_inst_id", "bk_obj_id", "bk_asst_inst_id", "bk_asst_obj_id", "bk_obj_asst_id"},
			Page:   metadata.BasePage{Start: 0, Limit: common.BKMaxInstanceLimit + 1},
		}

		rsp, err := asstClient.SearchInstanceAssociations(context.Background(), header, "bk_router", input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})

	It("count instance associations", func() {
		input := &metadata.CommonCountFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd,
					Rules: []querybuilder.Rule{
						&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual, Value: routerInstId1},
					},
				},
			},
		}

		rsp, err := asstClient.CountInstanceAssociations(context.Background(), header, "bk_router", input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))

		data, err := mapstr.NewFromInterface(rsp.Data)
		Expect(err).NotTo(HaveOccurred())

		count, err := data["count"].(json.Number).Int64()
		Expect(err).NotTo(HaveOccurred())
		Expect(int(count)).To(Equal(2))
	})

	It("count instance associations without conditions", func() {
		input := &metadata.CommonCountFilter{}

		rsp, err := asstClient.CountInstanceAssociations(context.Background(), header, "bk_router", input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))

		data, err := mapstr.NewFromInterface(rsp.Data)
		Expect(err).NotTo(HaveOccurred())

		count, err := data["count"].(json.Number).Int64()
		Expect(err).NotTo(HaveOccurred())
		Expect(int(count)).To(Equal(5))
	})

	It("count instance associations with OR conditions more than 20", func() {
		rules := []querybuilder.Rule{}
		ruleItem := &querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual, Value: routerInstId1}
		for i := 0; i < querybuilder.DefaultMaxConditionOrRulesCount+1; i++ {
			rules = append(rules, ruleItem)
		}

		input := &metadata.CommonCountFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionOr,
					Rules:     rules,
				},
			},
		}

		rsp, err := asstClient.CountInstanceAssociations(context.Background(), header, "bk_router", input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})

	It("count instance associations with operator value elements count more than 500", func() {
		var values []string
		for i := 0; i < querybuilder.DefaultMaxSliceElementsCount+1; i++ {
			values = append(values, fmt.Sprintf("%d", i))
		}

		input := &metadata.CommonCountFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionOr,
					Rules: []querybuilder.Rule{
						&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual, Value: values},
					},
				},
			},
		}

		rsp, err := asstClient.CountInstanceAssociations(context.Background(), header, "bk_router", input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})

	It("count instance associations with conditions deep more than 3", func() {
		deep4 := &querybuilder.CombinedRule{
			Condition: querybuilder.ConditionOr,
			Rules: []querybuilder.Rule{
				&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual, Value: routerInstId1},
			},
		}

		deep3 := &querybuilder.CombinedRule{
			Condition: querybuilder.ConditionOr,
			Rules: []querybuilder.Rule{
				&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual, Value: routerInstId1},
				deep4,
			},
		}

		deep2 := &querybuilder.CombinedRule{
			Condition: querybuilder.ConditionOr,
			Rules: []querybuilder.Rule{
				&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual, Value: routerInstId1},
				deep3,
			},
		}

		deep1 := &querybuilder.QueryFilter{
			Rule: deep2,
		}

		input := &metadata.CommonCountFilter{
			Conditions: deep1,
		}

		rsp, err := asstClient.CountInstanceAssociations(context.Background(), header, "bk_router", input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})

	// check "DeleteInstBatch" "the number of IDs should be less than 500." function.
	It("delete inst association batch", func() {
		list := make([]int64, 501, 501)
		input := &metadata.DeleteAssociationInstBatchRequest{
			ID:       list,
			ObjectID: "bk_router",
		}
		rsp, err := asstClient.DeleteInstBatch(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})
	// check "DeleteInstBatch" necessary input: bk_obj_id
	It("delete inst association batch", func() {
		input := &metadata.DeleteAssociationInstBatchRequest{
			ID:       []int64{instAsst1, instAsst2},
			ObjectID: "",
		}
		rsp, err := asstClient.DeleteInstBatch(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})
	// check "DeleteInstBatch" features available.
	It("delete inst association batch", func() {
		input := &metadata.DeleteAssociationInstBatchRequest{
			ID:       []int64{instAsst1, instAsst2, instAsst3, instAsst4, instAsst5},
			ObjectID: "bk_router",
		}
		rsp, err := asstClient.DeleteInstBatch(context.Background(), header, input)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(Equal(5))
	})
})
