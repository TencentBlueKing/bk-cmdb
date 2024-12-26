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
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func generateObject(objectIDs []string) {
	for _, objID := range objectIDs {
		count, err := test.GetDB().Table(common.BKTableNameObjDes).Find(map[string]interface{}{
			"bk_obj_id": objID}).Count(context.Background())

		Expect(err).NotTo(HaveOccurred())
		if count == 0 {
			input := metadata.Object{
				ObjCls:     "bk_uncategorized",
				ObjIcon:    "icon-cc-business",
				ObjectID:   objID,
				ObjectName: objID,
				Creator:    "admin",
			}
			rsp, err := objectClient.CreateObject(context.Background(), header, input)
			util.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.ObjCls).To(Equal(input.ObjCls))
			Expect(rsp.Data.ObjIcon).To(Equal(input.ObjIcon))
			Expect(rsp.Data.ObjectID).To(Equal(input.ObjectID))
			Expect(rsp.Data.ObjectName).To(Equal(input.ObjectName))
			Expect(rsp.Data.Creator).To(Equal(input.Creator))
		}
	}
}

var _ = Describe("inst test", func() {
	var switchInstId1, switchInstId2, switchInstId3, switchInstId4, switchInstId5, routerInstId1, routerInstId2,
		routerInstId3, instAsst1, instAsst2, instAsst3, instAsst4, instAsst5 int64

	It("create inst bk_obj_id='switch'", func() {
		objectIDs := []string{"switch", "router"}
		generateObject(objectIDs)

		input := map[string]interface{}{
			"bk_inst_name": "switch_2",
		}
		rsp, err := instClient.CreateInst(context.Background(), "switch", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("switch_2"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("switch"))
		switchInstId2, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create inst bk_obj_id='switch'", func() {
		input := map[string]interface{}{
			"bk_inst_name": "switch_1",
		}
		rsp, err := instClient.CreateInst(context.Background(), "switch", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("switch_1"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("switch"))
		switchInstId1, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create inst bk_obj_id='switch'", func() {
		input := map[string]interface{}{
			"bk_inst_name": "switch_3",
		}
		rsp, err := instClient.CreateInst(context.Background(), "switch", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("switch_3"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("switch"))
		switchInstId3, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create inst bk_obj_id='switch'", func() {
		input := map[string]interface{}{
			"bk_inst_name": "switch_4",
		}
		rsp, err := instClient.CreateInst(context.Background(), "switch", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("switch_4"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("switch"))
		switchInstId4, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create inst bk_obj_id='switch'", func() {
		input := map[string]interface{}{
			"bk_inst_name": "switch_5",
		}
		rsp, err := instClient.CreateInst(context.Background(), "switch", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("switch_5"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("switch"))
		switchInstId5, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create inst bk_obj_id='router'", func() {
		input := map[string]interface{}{
			"bk_inst_name": "router_1",
		}
		rsp, err := instClient.CreateInst(context.Background(), "router", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("router_1"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("router"))
		routerInstId1, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create inst bk_obj_id='router'", func() {
		input := map[string]interface{}{
			"bk_inst_name": "router_2",
		}
		rsp, err := instClient.CreateInst(context.Background(), "router", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("router_2"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("router"))
		routerInstId2, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create inst bk_obj_id='router'", func() {
		input := map[string]interface{}{
			"bk_inst_name": "router_3",
		}
		rsp, err := instClient.CreateInst(context.Background(), "router", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data["bk_inst_name"].(string)).To(Equal("router_3"))
		Expect(rsp.Data["bk_obj_id"].(string)).To(Equal("router"))
		routerInstId3, err = commonutil.GetInt64ByInterface(rsp.Data["bk_inst_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("create association ='router_default_switch'", func() {
		input := &metadata.Association{
			AsstKindID:           "default",
			AsstObjID:            "switch",
			AssociationName:      "router_default_switch",
			AssociationAliasName: "",
			ObjectID:             "router",
			Mapping:              "n:n",
		}
		rsp, err := asstClient.CreateObject(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("create association ='router_belong_switch'", func() {
		input := &metadata.Association{
			AsstKindID:           "belong",
			AsstObjID:            "switch",
			AssociationName:      "router_belong_switch",
			AssociationAliasName: "",
			ObjectID:             "router",
			Mapping:              "1:1",
		}
		rsp, err := asstClient.CreateObject(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("create association ='router_connect_switch'", func() {
		input := &metadata.Association{
			AsstKindID:           "connect",
			AsstObjID:            "switch",
			AssociationName:      "router_connect_switch",
			AssociationAliasName: "",
			ObjectID:             "router",
			Mapping:              "1:n",
		}
		rsp, err := asstClient.CreateObject(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("create inst association ='router1_default_switch1'", func() {
		input := &metadata.CreateAssociationInstRequest{
			ObjectAsstID: "router_default_switch",
			InstID:       routerInstId1,
			AsstInstID:   switchInstId1,
		}
		rsp, err := asstClient.CreateInst(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		instAsst1 = rsp.Data.ID
	})

	It("create inst association ='router1_default_switch2'", func() {
		input := &metadata.CreateAssociationInstRequest{
			ObjectAsstID: "router_default_switch",
			InstID:       routerInstId1,
			AsstInstID:   switchInstId2,
		}
		rsp, err := asstClient.CreateInst(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		instAsst2 = rsp.Data.ID
	})

	It("create inst association ='router_belong_switch', belong mapping 1:1", func() {
		input := &metadata.CreateManyInstAsstRequest{
			ObjectAsstID: "router_belong_switch",
			ObjectID:     "router",
			AsstObjectID: "switch",
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
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(len(rsp.Data.SuccessCreated)).To(Equal(1))
		Expect(len(rsp.Data.Error)).To(Equal(2))
		instAsst3 = rsp.Data.SuccessCreated[0]
	})

	It("create inst association ='router_connect_switch', connect mapping 1:n", func() {
		input := &metadata.CreateManyInstAsstRequest{
			ObjectAsstID: "router_connect_switch",
			ObjectID:     "router",
			AsstObjectID: "switch",
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
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(len(rsp.Data.SuccessCreated)).To(Equal(2))
		Expect(len(rsp.Data.Error)).To(Equal(0))
		instAsst4 = rsp.Data.SuccessCreated[0]
		instAsst5 = rsp.Data.SuccessCreated[1]
	})

	It("create inst association ='router_belong_switch', belong mapping 1:1", func() {
		input := &metadata.CreateManyInstAsstRequest{
			ObjectAsstID: "router_belong_switch",
			ObjectID:     "router",
			AsstObjectID: "router",
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
		util.RegisterResponseWithRid(rsp, header)
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
				ObjectID: "router",
				InstID:   routerInstId1,
			},
		}
		rsp, err := asstClient.SearchAssociationRelatedInst(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(len(rsp.Data)).To(Equal(2))
		Expect(rsp.Data[0].ObjectAsstID).To(Equal("router_default_switch"))
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
				ObjectID: "router",
				InstID:   routerInstId1,
			},
		}
		rsp, err := asstClient.SearchAssociationRelatedInst(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
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
				ObjectID: "router",
				InstID:   routerInstId1,
			},
		}
		rsp, err := asstClient.SearchAssociationRelatedInst(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})

	It("search instance associations and instances detail", func() {
		input := &metadata.InstAndAssocRequest{}
		input.Condition.AsstFilter = &querybuilder.QueryFilter{Rule: querybuilder.CombinedRule{
			Condition: querybuilder.ConditionAnd,
			Rules: []querybuilder.Rule{
				&querybuilder.AtomRule{
					Field: common.BKObjIDField, Operator: querybuilder.OperatorEqual, Value: "router",
				},
				&querybuilder.AtomRule{
					Field: common.BKInstIDField, Operator: querybuilder.OperatorEqual, Value: routerInstId1,
				},
			},
		}}
		input.Condition.SrcDetail = true
		input.Page.Limit = 200

		rsp, err := asstClient.SearchInstAssocAndInstDetail(context.Background(), header, "router", input)
		util.RegisterResponseWithRid(rsp, header)
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
					Field: common.BKObjIDField, Operator: querybuilder.OperatorEqual, Value: "router",
				},
				&querybuilder.AtomRule{
					Field: common.BKInstIDField, Operator: querybuilder.OperatorEqual, Value: routerInstId1,
				},
			},
		}}
		input.Condition.SrcDetail = true
		input.Page.Limit = 201

		rsp, err := asstClient.SearchInstAssocAndInstDetail(context.Background(), header, "router", input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).To(HaveOccurred())
	})

	It("search instance associations", func() {
		input := &metadata.CommonSearchFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd,
					Rules: []querybuilder.Rule{
						&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual,
							Value: routerInstId1},
					},
				},
			},
			Fields: []string{"bk_asst_id", "bk_inst_id", "bk_obj_id", "bk_asst_inst_id", "bk_asst_obj_id",
				"bk_obj_asst_id"},
			Page: metadata.BasePage{Start: 0, Limit: 1},
		}

		rsp, err := asstClient.SearchInstanceAssociations(context.Background(), header, "router", input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))

		data, err := mapstr.NewFromInterface(rsp.Data)
		Expect(err).NotTo(HaveOccurred())

		info, ok := data["info"].([]interface{})
		Expect(ok).To(Equal(true))
		Expect(len(info)).To(Equal(1))

		association, ok := info[0].(map[string]interface{})
		Expect(ok).To(Equal(true))

		Expect(association["bk_obj_asst_id"].(string)).To(Equal("router_default_switch"))
	})

	It("search instance associations without conditions", func() {
		input := &metadata.CommonSearchFilter{
			Fields: []string{"bk_asst_id", "bk_inst_id", "bk_obj_id", "bk_asst_inst_id", "bk_asst_obj_id",
				"bk_obj_asst_id"},
			Page: metadata.BasePage{Start: 0, Limit: 1},
		}

		rsp, err := asstClient.SearchInstanceAssociations(context.Background(), header, "router", input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))

		data, err := mapstr.NewFromInterface(rsp.Data)
		Expect(err).NotTo(HaveOccurred())

		info, ok := data["info"].([]interface{})
		Expect(ok).To(Equal(true))
		Expect(len(info)).To(Equal(1))

		association, ok := info[0].(map[string]interface{})
		Expect(ok).To(Equal(true))

		Expect(association["bk_obj_asst_id"].(string)).To(Equal("router_default_switch"))
	})

	It("search instance associations without fields", func() {
		input := &metadata.CommonSearchFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd,
					Rules: []querybuilder.Rule{
						&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual,
							Value: routerInstId1},
					},
				},
			},
			Page: metadata.BasePage{Start: 0, Limit: 1},
		}

		rsp, err := asstClient.SearchInstanceAssociations(context.Background(), header, "router", input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))

		data, err := mapstr.NewFromInterface(rsp.Data)
		Expect(err).NotTo(HaveOccurred())

		info, ok := data["info"].([]interface{})
		Expect(ok).To(Equal(true))
		Expect(len(info)).To(Equal(1))

		association, ok := info[0].(map[string]interface{})
		Expect(ok).To(Equal(true))
		Expect(len(association)).To(Equal(7))
	})

	It("search instance associations without page", func() {
		input := &metadata.CommonSearchFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd,
					Rules: []querybuilder.Rule{
						&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual,
							Value: routerInstId1},
					},
				},
			},
			Fields: []string{"bk_asst_id", "bk_inst_id", "bk_obj_id", "bk_asst_inst_id", "bk_asst_obj_id",
				"bk_obj_asst_id"},
		}

		rsp, err := asstClient.SearchInstanceAssociations(context.Background(), header, "router", input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})

	It("search instance associations with limit more than 500", func() {
		input := &metadata.CommonSearchFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd,
					Rules: []querybuilder.Rule{
						&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual,
							Value: routerInstId1},
					},
				},
			},
			Fields: []string{"bk_asst_id", "bk_inst_id", "bk_obj_id", "bk_asst_inst_id", "bk_asst_obj_id",
				"bk_obj_asst_id"},
			Page: metadata.BasePage{Start: 0, Limit: common.BKMaxInstanceLimit + 1},
		}

		rsp, err := asstClient.SearchInstanceAssociations(context.Background(), header, "router", input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})

	It("count instance associations", func() {
		input := &metadata.CommonCountFilter{
			Conditions: &querybuilder.QueryFilter{
				Rule: querybuilder.CombinedRule{
					Condition: querybuilder.ConditionAnd,
					Rules: []querybuilder.Rule{
						&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual,
							Value: routerInstId1},
					},
				},
			},
		}

		rsp, err := asstClient.CountInstanceAssociations(context.Background(), header, "router", input)
		util.RegisterResponseWithRid(rsp, header)
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

		rsp, err := asstClient.CountInstanceAssociations(context.Background(), header, "router", input)
		util.RegisterResponseWithRid(rsp, header)
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
		ruleItem := &querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual,
			Value: routerInstId1}
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

		rsp, err := asstClient.CountInstanceAssociations(context.Background(), header, "router", input)
		util.RegisterResponseWithRid(rsp, header)
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
						&querybuilder.AtomRule{Field: "bk_inst_id", Operator: querybuilder.OperatorEqual,
							Value: values},
					},
				},
			},
		}

		rsp, err := asstClient.CountInstanceAssociations(context.Background(), header, "router", input)
		util.RegisterResponseWithRid(rsp, header)
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

		rsp, err := asstClient.CountInstanceAssociations(context.Background(), header, "router", input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})

	// check "DeleteInstBatch" "the number of IDs should be less than 500." function.
	It("delete inst association batch", func() {
		list := make([]int64, 501, 501)
		input := &metadata.DeleteAssociationInstBatchRequest{
			ID:       list,
			ObjectID: "router",
		}
		rsp, err := asstClient.DeleteInstBatch(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
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
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(false))
	})
	// check "DeleteInstBatch" features available.
	It("delete inst association batch", func() {
		input := &metadata.DeleteAssociationInstBatchRequest{
			ID:       []int64{instAsst1, instAsst2, instAsst3, instAsst4, instAsst5},
			ObjectID: "router",
		}
		rsp, err := asstClient.DeleteInstBatch(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(Equal(5))
		test.DeleteAllObjects()
	})
})
