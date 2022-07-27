package topo_server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/test"
	testutil "configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var bizID, categoryId, serviceTemplateID, serviceTemplateID2, serviceTemplateID3 int64

var _ = Describe("create empty set template test", func() {
	bizID := int64(2)
	ctx := context.Background()
	It("create set template", func() {
		option := metadata.CreateSetTemplateOption{
			Name:               "setTpl1",
			ServiceTemplateIDs: nil,
		}
		rsp, err := topoServerClient.SetTemplate().CreateSetTemplate(ctx, header, bizID, option)
		Expect(err).To(HaveOccurred())
		Expect(err.GetCode()).Should(Equal(common.CCErrCommParamsNeedSet))
		Expect(rsp).To(BeNil())
	})
})

var _ = Describe("create normal set template test", func() {
	ctx := context.Background()

	BeforeEach(prepareSetTemplateData)

	It("normal set template test", func() {
		var setTemplateID int64
		By("create set template")
		func() {
			option := metadata.CreateSetTemplateOption{
				Name:               "setTpl2",
				ServiceTemplateIDs: []int64{serviceTemplateID},
			}
			rsp, err := topoServerClient.SetTemplate().CreateSetTemplate(ctx, header, bizID, option)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data.Name).To(Equal("setTpl2"))
			Expect(rsp.Data.ID).To(Not(Equal(int64(0))))
			setTemplateID = rsp.Data.ID
		}()

		By("update set template")
		func() {
			option := metadata.UpdateSetTemplateOption{
				Name:               "setTpl3",
				ServiceTemplateIDs: []int64{serviceTemplateID2, serviceTemplateID3},
			}
			rsp, err := topoServerClient.SetTemplate().UpdateSetTemplate(ctx, header, bizID, setTemplateID, option)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data.Name).To(Equal("setTpl3"))
			Expect(rsp.Data.ID).To(Equal(setTemplateID))
		}()

		By("list set-template")
		func() {
			option := metadata.ListSetTemplateOption{
				SetTemplateIDs: []int64{setTemplateID},
				Page: metadata.BasePage{
					Limit: common.BKNoLimit,
				},
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplate(ctx, header, bizID, option)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(BeNumerically(">", 0))
			s, e := json.Marshal(rsp.Info)
			Expect(e).Should(BeNil())
			Expect(string(s)).Should(ContainSubstring(fmt.Sprintf(`"id":%d`, setTemplateID)))
		}()

		By("list set-template related service templates")
		func() {
			rsp, err := topoServerClient.SetTemplate().ListSetTplRelatedSvcTpl(ctx, header, bizID, setTemplateID)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp).To(HaveLen(2))
			s, e := json.Marshal(rsp)
			Expect(e).Should(BeNil())
			Expect(string(s)).Should(ContainSubstring(fmt.Sprintf(`"id":%d`, serviceTemplateID2)))
			Expect(string(s)).Should(ContainSubstring(fmt.Sprintf(`"id":%d`, serviceTemplateID3)))
		}()

		var setID int64
		func() {
			// create set
			data := map[string]interface{}{
				"bk_set_name":       "set1",
				"bk_set_desc":       "",
				"bk_set_env":        "3",
				"bk_service_status": "1",
				"description":       "",
				"bk_capacity":       nil,
				"bk_biz_id":         bizID,
				"bk_parent_id":      bizID,
				"metadata": map[string]interface{}{
					"label": map[string]interface{}{
						"bk_biz_id": strconv.FormatInt(bizID, 10),
					},
				},
				"bk_supplier_account": "0",
				"set_template_id":     setTemplateID,
			}
			rsp, err := topoServerClient.Instance().CreateSet(ctx, bizID, header, data)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp[common.BKSetNameField]).To(Equal("set1"))
			Expect(rsp[common.BKSetIDField]).To(Not(Equal(int64(0))))
			var e error
			setID, e = util.GetInt64ByInterface(rsp[common.BKSetIDField])
			Expect(e).To(BeNil())

			s, e := json.Marshal(rsp)
			Expect(e).Should(BeNil())
			Expect(string(s)).Should(ContainSubstring(fmt.Sprintf(`"set_template_id":%d`, setTemplateID)))
		}()

		By("list set-template related set")
		func() {
			option := metadata.ListSetByTemplateOption{
				Page: metadata.BasePage{
					Limit: common.BKNoLimit,
				},
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTplRelatedSetsWeb(ctx, header, bizID, setTemplateID, option)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).Should(Equal(1))
			Expect(rsp.Info).Should(HaveLen(1))
			s, e := json.Marshal(rsp)
			Expect(e).Should(BeNil())
			Expect(string(s)).Should(ContainSubstring(fmt.Sprintf(`"bk_set_id":%d`, setID)))
		}()

		// update set-template and check diff result
		var serviceTemplateID4 int64
		By("create service template 4")
		func() {
			data := map[string]interface{}{
				"bk_biz_id":           bizID,
				"name":                "svcTpl4",
				"service_category_id": categoryId,
			}
			rsp, err := procServerClient.Service().CreateServiceTemplate(ctx, header, data)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data[common.BKFieldID]).To(Not(Equal(int64(0))))
			Expect(rsp.Data[common.BKFieldName]).To(Equal("svcTpl4"))
			serviceTemplateID4, err = util.GetInt64ByInterface(rsp.Data[common.BKFieldID])
			Expect(err).To(BeNil())
		}()

		By("update set template")
		func() {
			option := metadata.UpdateSetTemplateOption{
				Name:               "setTpl4",
				ServiceTemplateIDs: []int64{serviceTemplateID3, serviceTemplateID4},
			}
			rsp, err := topoServerClient.SetTemplate().UpdateSetTemplate(ctx, header, bizID, setTemplateID, option)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data.Name).To(Equal("setTpl4"))
			Expect(rsp.Data.ID).To(Equal(setTemplateID))
		}()

		By("diff set template with set")
		func() {
			option := metadata.DiffSetTplWithInstOption{
				SetID: setID,
			}
			setTplDiffResult, err := topoServerClient.SetTemplate().DiffSetTplWithInst(ctx, header, bizID, setTemplateID, option)
			Expect(err).NotTo(HaveOccurred())
			setDiff := setTplDiffResult.Difference
			Expect(setDiff.SetID).To(Equal(setID))
			Expect(setDiff.ModuleDiffs).To(HaveLen(3))
			m := MatchFields(IgnoreMissing|IgnoreExtras, Fields{
				"DiffType": Equal(metadata.ModuleDiffAdd),
			})
			Expect(setDiff.ModuleDiffs).Should(ContainElement(m))
			m = MatchFields(IgnoreMissing|IgnoreExtras, Fields{
				"DiffType": Equal(metadata.ModuleDiffRemove),
			})
			Expect(setDiff.ModuleDiffs).Should(ContainElement(m))
		}()

		// delete setTemplate be referenced
		By("delete set template")
		func() {
			option := metadata.DeleteSetTemplateOption{
				SetTemplateIDs: []int64{setTemplateID},
			}
			err := topoServerClient.SetTemplate().DeleteSetTemplate(ctx, header, bizID, option)
			Expect(err).NotTo(BeNil())
		}()

		By("list set-template")
		func() {
			option := metadata.ListSetTemplateOption{
				SetTemplateIDs: []int64{setTemplateID},
				Page: metadata.BasePage{
					Limit: common.BKNoLimit,
				},
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplate(ctx, header, bizID, option)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(BeNumerically(">", 0))
			s, e := json.Marshal(rsp.Info)
			Expect(e).Should(BeNil())
			Expect(string(s)).Should(ContainSubstring(fmt.Sprintf(`"id":%d`, setTemplateID)))
		}()

		// delete set, then delete set template
		By("delete set")
		func() {
			err := topoServerClient.Instance().DeleteSet(ctx, bizID, setID, header)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(BeNil())
		}()

		By("delete set template")
		func() {
			option := metadata.DeleteSetTemplateOption{
				SetTemplateIDs: []int64{setTemplateID},
			}
			err := topoServerClient.SetTemplate().DeleteSetTemplate(ctx, header, bizID, option)
			Expect(err).To(BeNil())
		}()

		By("list set-template")
		func() {
			option := metadata.ListSetTemplateOption{
				SetTemplateIDs: []int64{setTemplateID},
				Page: metadata.BasePage{
					Limit: common.BKNoLimit,
				},
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplate(ctx, header, bizID, option)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(BeZero())
			s, e := json.Marshal(rsp.Info)
			Expect(e).Should(BeNil())
			Expect(string(s)).ShouldNot(ContainSubstring(fmt.Sprintf(`"id":%d`, setTemplateID)))
		}()
	})
})

var _ = Describe("set template attribute test", func() {
	ctx := context.Background()

	setAttrMap := make(map[string]metadata.Attribute)

	BeforeEach(func() {
		prepareSetTemplateData()

		By("create set attributes and then get all set attributes for later use")
		func() {
			input := &metadata.CreateModelAttributes{
				Attributes: []metadata.Attribute{{
					ObjectID:     common.BKInnerObjIDSet,
					PropertyID:   "int_attr",
					PropertyName: "int_attr",
					IsEditable:   true,
					PropertyType: common.FieldTypeInt,
				}, {
					ObjectID:     common.BKInnerObjIDSet,
					PropertyID:   "str_attr",
					PropertyName: "str_attr",
					IsEditable:   true,
					PropertyType: common.FieldTypeSingleChar,
					BizID:        bizID,
				}, {
					ObjectID:     common.BKInnerObjIDSet,
					PropertyID:   "enum_attr",
					PropertyName: "enum_attr",
					IsEditable:   true,
					PropertyType: common.FieldTypeEnum,
					Option: []metadata.EnumVal{{
						ID:        "key1",
						Name:      "value1",
						Type:      "text",
						IsDefault: true,
					}, {
						ID:        "key2",
						Name:      "value2",
						Type:      "text",
						IsDefault: false,
					}},
					BizID: bizID,
				}},
			}
			res, err := clientSet.CoreService().Model().CreateModelAttrs(ctx, header, common.BKInnerObjIDSet, input)
			testutil.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())

			readInput := &metadata.QueryCondition{
				Page:           metadata.BasePage{Limit: common.BKNoLimit},
				DisableCounter: true,
			}
			rsp, err := clientSet.CoreService().Model().ReadModelAttr(ctx, header, common.BKInnerObjIDSet, readInput)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())

			for _, attr := range rsp.Info {
				setAttrMap[attr.PropertyID] = attr
			}
		}()
	})

	It("normal set template attribute test", func() {
		var setTemplateID int64
		svcTempAttrs := []metadata.SetTempAttr{{
			AttributeID:   setAttrMap["int_attr"].ID,
			PropertyValue: 1,
		}, {
			AttributeID:   setAttrMap["str_attr"].ID,
			PropertyValue: "str",
		}}

		By("create set template all info", func() {
			option := &metadata.CreateSetTempAllInfoOption{
				BizID:              bizID,
				Name:               "set_template",
				ServiceTemplateIDs: []int64{serviceTemplateID, serviceTemplateID2},
				Attributes:         svcTempAttrs,
			}

			var err error
			setTemplateID, err = topoServerClient.SetTemplate().CreateSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(setTemplateID, header)
			Expect(err).NotTo(HaveOccurred())
		})

		By("get set template all info", func() {
			option := &metadata.GetSetTempAllInfoOption{
				ID:    setTemplateID,
				BizID: bizID,
			}
			rsp, err := topoServerClient.SetTemplate().GetSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.ID).To(Equal(setTemplateID))
			Expect(rsp.BizID).To(Equal(bizID))
			Expect(rsp.Name).To(Equal("set_template"))
			Expect(rsp.ServiceTemplateIDs).To(ConsistOf(serviceTemplateID, serviceTemplateID2))
			Expect(len(rsp.Attributes)).To(Equal(2))
			Expect(rsp.Attributes[0].AttributeID).To(Equal(svcTempAttrs[0].AttributeID))
			intVal, e := util.GetIntByInterface(rsp.Attributes[0].PropertyValue)
			Expect(e).NotTo(HaveOccurred())
			Expect(intVal).To(Equal(1))
			Expect(rsp.Attributes[1].AttributeID).To(Equal(svcTempAttrs[1].AttributeID))
			Expect(rsp.Attributes[1].PropertyValue).To(Equal(svcTempAttrs[1].PropertyValue))
		})

		var setID int64
		By("create set using template", func() {
			data := map[string]interface{}{
				"bk_set_name":     "set1",
				"bk_biz_id":       bizID,
				"bk_parent_id":    bizID,
				"set_template_id": setTemplateID,
			}
			rsp, e := topoServerClient.Instance().CreateSet(ctx, bizID, header, data)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(e).NotTo(HaveOccurred())
			var err error
			setID, err = util.GetInt64ByInterface(rsp[common.BKSetIDField])
			Expect(err).To(BeNil())
		})

		By("check set template related set has the attributes", func() {
			opt := metadata.ListSetByTemplateOption{
				Page: metadata.BasePage{Limit: common.BKNoLimit},
			}
			rsp, e := topoServerClient.SetTemplate().ListSetTplRelatedSetsWeb(ctx, header, bizID, setTemplateID, opt)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(e).NotTo(HaveOccurred())
			Expect(rsp.Count).Should(Equal(1))
			Expect(rsp.Info).Should(HaveLen(1))
			createdSetID, err := util.GetInt64ByInterface(rsp.Info[0][common.BKSetIDField])
			Expect(err).To(BeNil())
			Expect(createdSetID).Should(Equal(setID))
			intVal, err := util.GetInt64ByInterface(rsp.Info[0]["int_attr"])
			Expect(err).To(BeNil())
			Expect(intVal).Should(Equal(int64(1)))
			Expect(util.GetStrByInterface(rsp.Info[0]["str_attr"])).Should(Equal("str"))
		})

		By("check set modules", func() {
			input := &params.SearchParams{
				Condition: mapstr.MapStr{common.BKSetIDField: setID},
				Page:      mapstr.MapStr{"sort": common.BKModuleNameField},
			}
			rsp, err := instClient.SearchModule(context.Background(), "0", bizID, setID, header, input)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(Equal(2))
			Expect(util.GetStrByInterface(rsp.Info[0][common.BKModuleNameField])).To(Equal("svcTpl1"))
			Expect(util.GetStrByInterface(rsp.Info[1][common.BKModuleNameField])).To(Equal("svcTpl2"))
		})

		By("update set without set template attributes", func() {
			input := map[string]interface{}{
				"bk_set_name": "set2",
				"int_attr":    1,
				"str_attr":    "str",
				"enum_attr":   "key2",
			}
			err := instClient.UpdateSet(ctx, bizID, setID, header, input)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).NotTo(HaveOccurred())
		})

		By("update set template all info", func() {
			svcTempAttrs = []metadata.SetTempAttr{{
				AttributeID:   setAttrMap["int_attr"].ID,
				PropertyValue: 2,
			}, {
				AttributeID:   setAttrMap["enum_attr"].ID,
				PropertyValue: "key1",
			}}

			option := &metadata.UpdateSetTempAllInfoOption{
				ID:                 setTemplateID,
				BizID:              bizID,
				Name:               "set_template1",
				ServiceTemplateIDs: []int64{serviceTemplateID2, serviceTemplateID3},
				Attributes:         svcTempAttrs,
			}

			err := topoServerClient.SetTemplate().UpdateSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(nil, header)
			Expect(err).NotTo(HaveOccurred())
		})

		By("check updated set template all info", func() {
			option := &metadata.GetSetTempAllInfoOption{
				ID:    setTemplateID,
				BizID: bizID,
			}
			rsp, err := topoServerClient.SetTemplate().GetSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.ID).To(Equal(setTemplateID))
			Expect(rsp.BizID).To(Equal(bizID))
			Expect(rsp.Name).To(Equal("set_template1"))
			Expect(rsp.ServiceTemplateIDs).To(ConsistOf(serviceTemplateID2, serviceTemplateID3))
			Expect(len(rsp.Attributes)).To(Equal(2))
			Expect(rsp.Attributes[0].AttributeID).To(Equal(svcTempAttrs[0].AttributeID))
			intVal, e := util.GetIntByInterface(rsp.Attributes[0].PropertyValue)
			Expect(e).NotTo(HaveOccurred())
			Expect(intVal).To(Equal(2))
			Expect(rsp.Attributes[1].AttributeID).To(Equal(svcTempAttrs[1].AttributeID))
			Expect(rsp.Attributes[1].PropertyValue).To(Equal(svcTempAttrs[1].PropertyValue))
		})

		By("diff set template with set", func() {
			option := metadata.DiffSetTplWithInstOption{
				SetID: setID,
			}
			res, err := topoServerClient.SetTemplate().DiffSetTplWithInst(ctx, header, bizID, setTemplateID, option)
			testutil.RegisterResponseWithRid(res, header)
			Expect(err).NotTo(HaveOccurred())
			setDiff := res.Difference
			Expect(setDiff.SetID).To(Equal(setID))
			Expect(setDiff.ModuleDiffs).To(HaveLen(3))
			for _, diff := range setDiff.ModuleDiffs {
				if diff.DiffType == metadata.ModuleDiffAdd {
					Expect(diff.ServiceTemplateID).To(Equal(serviceTemplateID3))
				} else if diff.DiffType == metadata.ModuleDiffUnchanged {
					Expect(diff.ServiceTemplateID).To(Equal(serviceTemplateID2))
				} else {
					Expect(diff.DiffType).To(Equal(metadata.ModuleDiffRemove))
					Expect(diff.ServiceTemplateID).To(Equal(serviceTemplateID))
				}
			}
			Expect(len(setDiff.Attributes)).To(Equal(2))
			for _, attr := range setDiff.Attributes {
				if attr.ID == setAttrMap["int_attr"].ID {
					templatePropertyValue, err := util.GetIntByInterface(attr.TemplatePropertyValue)
					Expect(err).To(BeNil())
					Expect(templatePropertyValue).To(Equal(2))
					instancePropertyValue, err := util.GetIntByInterface(attr.InstancePropertyValue)
					Expect(err).To(BeNil())
					Expect(instancePropertyValue).To(Equal(1))
				} else {
					Expect(attr.ID).To(Equal(setAttrMap["enum_attr"].ID))
					Expect(util.GetStrByInterface(attr.TemplatePropertyValue)).To(Equal("key1"))
					Expect(util.GetStrByInterface(attr.InstancePropertyValue)).To(Equal("key2"))
				}
			}
		})

		By("sync set", func() {
			option := &metadata.SyncSetTplToInstOption{
				SetIDs: []int64{setID},
			}
			err := topoServerClient.SetTemplate().SyncSetTplToInst(ctx, header, bizID, setTemplateID, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(time.Second * 10)
		})

		By("check set attributes has changed", func() {
			opt := metadata.ListSetByTemplateOption{
				Page: metadata.BasePage{Limit: common.BKNoLimit},
			}
			rsp, e := topoServerClient.SetTemplate().ListSetTplRelatedSetsWeb(ctx, header, bizID, setTemplateID, opt)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(e).NotTo(HaveOccurred())
			Expect(rsp.Count).Should(Equal(1))
			Expect(rsp.Info).Should(HaveLen(1))
			createdSetID, err := util.GetInt64ByInterface(rsp.Info[0][common.BKSetIDField])
			Expect(err).To(BeNil())
			Expect(createdSetID).Should(Equal(setID))
			intVal, err := util.GetInt64ByInterface(rsp.Info[0]["int_attr"])
			Expect(err).To(BeNil())
			Expect(intVal).Should(Equal(int64(2)))
			Expect(util.GetStrByInterface(rsp.Info[0]["enum_attr"])).Should(Equal("key1"))
		})

		By("check set modules has changed", func() {
			input := &params.SearchParams{
				Condition: mapstr.MapStr{common.BKSetIDField: setID},
				Page:      mapstr.MapStr{"sort": common.BKModuleNameField},
			}
			rsp, err := instClient.SearchModule(context.Background(), "0", bizID, setID, header, input)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(Equal(2))
			Expect(util.GetStrByInterface(rsp.Info[0][common.BKModuleNameField])).To(Equal("svcTpl2"))
			Expect(util.GetStrByInterface(rsp.Info[1][common.BKModuleNameField])).To(Equal("svcTpl3"))
		})

		By("update set template attributes", func() {
			svcTempAttrs = []metadata.SetTempAttr{{
				AttributeID:   setAttrMap["int_attr"].ID,
				PropertyValue: 4,
			}, {
				AttributeID:   setAttrMap["enum_attr"].ID,
				PropertyValue: "key2",
			}}

			option := &metadata.UpdateSetTempAttrOption{
				BizID:      bizID,
				ID:         setTemplateID,
				Attributes: svcTempAttrs,
			}
			err := topoServerClient.SetTemplate().UpdateSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(BeNil())
		})

		var setTempAttrIDs []int64
		By("list set template attributes", func() {
			option := &metadata.ListSetTempAttrOption{
				BizID: bizID,
				ID:    setTemplateID,
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).To(BeNil())
			Expect(len(rsp.Attributes)).To(Equal(2))
			Expect(rsp.Attributes[0].AttributeID).To(Equal(svcTempAttrs[0].AttributeID))
			intVal, e := util.GetIntByInterface(rsp.Attributes[0].PropertyValue)
			Expect(e).NotTo(HaveOccurred())
			Expect(intVal).To(Equal(4))
			setTempAttrIDs = append(setTempAttrIDs, rsp.Attributes[0].AttributeID)
			Expect(rsp.Attributes[1].AttributeID).To(Equal(svcTempAttrs[1].AttributeID))
			Expect(rsp.Attributes[1].PropertyValue).To(Equal(svcTempAttrs[1].PropertyValue))
			setTempAttrIDs = append(setTempAttrIDs, rsp.Attributes[1].AttributeID)
		})

		By("delete set template attributes", func() {
			option := &metadata.DeleteSetTempAttrOption{
				BizID:        bizID,
				ID:           setTemplateID,
				AttributeIDs: []int64{setTempAttrIDs[0]},
			}
			err := topoServerClient.SetTemplate().DeleteSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(BeNil())
		})

		By("check set template attribute is deleted", func() {
			option := &metadata.ListSetTempAttrOption{
				BizID: bizID,
				ID:    setTemplateID,
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).To(BeNil())
			Expect(len(rsp.Attributes)).To(Equal(1))
			Expect(rsp.Attributes[0].AttributeID).To(Equal(setTempAttrIDs[1]))
		})

		By("delete set", func() {
			err := topoServerClient.Instance().DeleteSet(ctx, bizID, setID, header)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(BeNil())
		})

		By("delete set template", func() {
			option := metadata.DeleteSetTemplateOption{
				SetTemplateIDs: []int64{setTemplateID},
			}
			err := topoServerClient.SetTemplate().DeleteSetTemplate(ctx, header, bizID, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(BeNil())
		})

		By("check if set template is deleted", func() {
			option := metadata.ListSetTemplateOption{
				SetTemplateIDs: []int64{setTemplateID},
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplate(ctx, header, bizID, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(BeZero())
		})

		By("check if set template attributes are deleted", func() {
			option := &metadata.ListSetTempAttrOption{
				BizID: bizID,
				ID:    setTemplateID,
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		})
	})

	It("abnormal set template attribute test", func() {
		svcTempAttrs := []metadata.SetTempAttr{{
			AttributeID:   setAttrMap["int_attr"].ID,
			PropertyValue: 1,
		}, {
			AttributeID:   setAttrMap["str_attr"].ID,
			PropertyValue: "str",
		}}

		By("create set template all info with no name", func() {
			option := &metadata.CreateSetTempAllInfoOption{
				BizID:              bizID,
				ServiceTemplateIDs: []int64{serviceTemplateID, serviceTemplateID2},
				Attributes:         svcTempAttrs,
			}

			_, err := topoServerClient.SetTemplate().CreateSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("create set template all info with no service templates", func() {
			option := &metadata.CreateSetTempAllInfoOption{
				BizID:      bizID,
				Name:       "test1",
				Attributes: svcTempAttrs,
			}

			_, err := topoServerClient.SetTemplate().CreateSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("create set template all info with invalid attributes", func() {
			option := &metadata.CreateSetTempAllInfoOption{
				BizID:              bizID,
				Name:               "test2",
				ServiceTemplateIDs: []int64{serviceTemplateID, serviceTemplateID2},
				Attributes: []metadata.SetTempAttr{{
					AttributeID:   setAttrMap["str_attr"].ID,
					PropertyValue: 222,
				}},
			}

			_, err := topoServerClient.SetTemplate().CreateSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())

			option.Attributes = []metadata.SetTempAttr{{
				AttributeID:   setAttrMap[common.BKSetNameField].ID,
				PropertyValue: "test3",
			}}
			_, err = topoServerClient.SetTemplate().CreateSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		var setTemplateID int64
		By("create set template all info", func() {
			option := &metadata.CreateSetTempAllInfoOption{
				BizID:              bizID,
				Name:               "set_template",
				ServiceTemplateIDs: []int64{serviceTemplateID, serviceTemplateID2},
				Attributes:         svcTempAttrs,
			}

			var err error
			setTemplateID, err = topoServerClient.SetTemplate().CreateSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(setTemplateID, header)
			Expect(err).NotTo(HaveOccurred())
		})

		By("create set template all info with duplicate name", func() {
			option := &metadata.CreateSetTempAllInfoOption{
				BizID:              bizID,
				Name:               "set_template",
				ServiceTemplateIDs: []int64{serviceTemplateID, serviceTemplateID2},
			}

			_, err := topoServerClient.SetTemplate().CreateSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("check only one set template exists", func() {
			option := metadata.ListSetTemplateOption{
				Page: metadata.BasePage{Limit: common.BKNoLimit},
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplate(ctx, header, bizID, option)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Count).To(Equal(int64(1)))
			Expect(rsp.Info[0].ID).To(Equal(setTemplateID))
		})

		By("get set template all info with invalid id", func() {
			option := &metadata.GetSetTempAllInfoOption{
				ID:    10000,
				BizID: bizID,
			}
			_, err := topoServerClient.SetTemplate().GetSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		var setID int64
		By("create set using template", func() {
			data := map[string]interface{}{
				"bk_set_name":     "set1",
				"bk_biz_id":       bizID,
				"bk_parent_id":    bizID,
				"set_template_id": setTemplateID,
			}
			rsp, e := topoServerClient.Instance().CreateSet(ctx, bizID, header, data)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(e).NotTo(HaveOccurred())
			var err error
			setID, err = util.GetInt64ByInterface(rsp[common.BKSetIDField])
			Expect(err).To(BeNil())
		})

		By("update set using set template attributes", func() {
			input := map[string]interface{}{
				"int_attr": 5,
			}
			err := instClient.UpdateSet(ctx, bizID, setID, header, input)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("update set template all info with no id", func() {
			option := &metadata.UpdateSetTempAllInfoOption{
				BizID:              bizID,
				Name:               "test4",
				ServiceTemplateIDs: []int64{serviceTemplateID2, serviceTemplateID3},
			}

			err := topoServerClient.SetTemplate().UpdateSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("update set template all info with invalid id", func() {
			option := &metadata.UpdateSetTempAllInfoOption{
				ID:                 1000,
				BizID:              bizID,
				Name:               "test4",
				ServiceTemplateIDs: []int64{serviceTemplateID2, serviceTemplateID3},
			}

			err := topoServerClient.SetTemplate().UpdateSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("update set template all info with invalid attributes", func() {
			option := &metadata.UpdateSetTempAllInfoOption{
				ID:                 setTemplateID,
				BizID:              bizID,
				Name:               "test6",
				ServiceTemplateIDs: []int64{serviceTemplateID2, serviceTemplateID3},
				Attributes: []metadata.SetTempAttr{{
					AttributeID:   setAttrMap["int_attr"].ID,
					PropertyValue: "test",
				}},
			}

			err := topoServerClient.SetTemplate().UpdateSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())

			option.Attributes = []metadata.SetTempAttr{{
				AttributeID:   setAttrMap[common.BKSetNameField].ID,
				PropertyValue: "test7",
			}}
			err = topoServerClient.SetTemplate().UpdateSetTemplateAllInfo(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("diff set template with no set id", func() {
			option := metadata.DiffSetTplWithInstOption{}
			_, err := topoServerClient.SetTemplate().DiffSetTplWithInst(ctx, header, bizID, setTemplateID, option)
			Expect(err).To(HaveOccurred())
		})

		By("diff set template with no set template id", func() {
			option := metadata.DiffSetTplWithInstOption{
				SetID: setID,
			}
			_, err := topoServerClient.SetTemplate().DiffSetTplWithInst(ctx, header, bizID, 0, option)
			Expect(err).To(HaveOccurred())
		})

		By("diff set template with invalid set template id", func() {
			option := metadata.DiffSetTplWithInstOption{
				SetID: setID,
			}
			_, err := topoServerClient.SetTemplate().DiffSetTplWithInst(ctx, header, bizID, 1000, option)
			Expect(err).To(HaveOccurred())
		})

		By("sync set template with no set ids", func() {
			option := new(metadata.SyncSetTplToInstOption)
			err := topoServerClient.SetTemplate().SyncSetTplToInst(ctx, header, bizID, setTemplateID, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("sync set template with no set template id", func() {
			option := &metadata.SyncSetTplToInstOption{
				SetIDs: []int64{setID},
			}
			err := topoServerClient.SetTemplate().SyncSetTplToInst(ctx, header, bizID, 0, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("sync set template with invalid set template id", func() {
			option := &metadata.SyncSetTplToInstOption{
				SetIDs: []int64{setID},
			}
			err := topoServerClient.SetTemplate().SyncSetTplToInst(ctx, header, bizID, 1000, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("update set template attributes with not exist attribute", func() {
			option := &metadata.UpdateSetTempAttrOption{
				BizID: bizID,
				ID:    setTemplateID,
				Attributes: []metadata.SetTempAttr{{
					AttributeID:   setAttrMap["enum_attr"].ID,
					PropertyValue: "key2",
				}},
			}
			err := topoServerClient.SetTemplate().UpdateSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("update set template attributes with invalid attribute", func() {
			option := &metadata.UpdateSetTempAttrOption{
				BizID: bizID,
				ID:    setTemplateID,
				Attributes: []metadata.SetTempAttr{{
					AttributeID:   setAttrMap["int_attr"].ID,
					PropertyValue: "111",
				}},
			}
			err := topoServerClient.SetTemplate().UpdateSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("list set template attributes with no biz id", func() {
			option := &metadata.ListSetTempAttrOption{
				ID: setTemplateID,
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		})

		By("list set template attributes with no set template id", func() {
			option := &metadata.ListSetTempAttrOption{
				ID: setTemplateID,
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		})

		By("list set template attributes with invalid biz id", func() {
			option := &metadata.ListSetTempAttrOption{
				BizID: 1000,
				ID:    setTemplateID,
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		})

		By("list set template attributes with invalid set template id", func() {
			option := &metadata.ListSetTempAttrOption{
				BizID: bizID,
				ID:    1000,
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).To(HaveOccurred())
		})

		By("delete set template attributes with no ids", func() {
			option := &metadata.DeleteSetTempAttrOption{
				BizID: bizID,
				ID:    setTemplateID,
			}
			err := topoServerClient.SetTemplate().DeleteSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		var setTempAttrIDs []int64
		By("list set template attributes", func() {
			option := &metadata.ListSetTempAttrOption{
				BizID: bizID,
				ID:    setTemplateID,
			}
			rsp, err := topoServerClient.SetTemplate().ListSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(rsp, header)
			Expect(err).To(BeNil())
			Expect(len(rsp.Attributes)).To(Equal(2))
			for _, attribute := range rsp.Attributes {
				setTempAttrIDs = append(setTempAttrIDs, attribute.AttributeID)
			}
		})

		By("delete set template attributes with no biz id", func() {
			option := &metadata.DeleteSetTempAttrOption{
				ID:           setTemplateID,
				AttributeIDs: setTempAttrIDs,
			}
			err := topoServerClient.SetTemplate().DeleteSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("delete set template attributes with no template id", func() {
			option := &metadata.DeleteSetTempAttrOption{
				BizID:        bizID,
				AttributeIDs: setTempAttrIDs,
			}
			err := topoServerClient.SetTemplate().DeleteSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("delete set template attributes with invalid biz id", func() {
			option := &metadata.DeleteSetTempAttrOption{
				BizID:        1000,
				ID:           setTemplateID,
				AttributeIDs: setTempAttrIDs,
			}
			err := topoServerClient.SetTemplate().DeleteSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("delete set template attributes with invalid template id", func() {
			option := &metadata.DeleteSetTempAttrOption{
				BizID:        bizID,
				ID:           1000,
				AttributeIDs: setTempAttrIDs,
			}
			err := topoServerClient.SetTemplate().DeleteSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})

		By("delete set template attributes with invalid ids", func() {
			option := &metadata.DeleteSetTempAttrOption{
				BizID:        bizID,
				ID:           setTemplateID,
				AttributeIDs: []int64{1000},
			}
			err := topoServerClient.SetTemplate().DeleteSetTemplateAttr(ctx, header, option)
			testutil.RegisterResponseWithRid(err, header)
			Expect(err).To(HaveOccurred())
		})
	})
})

func prepareSetTemplateData() {
	ctx := context.Background()

	By("clear redundant data")
	func() {
		assts := make([]metadata.Association, 0)
		asstCond := mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline}
		err := test.GetDB().Table(common.BKTableNameObjAsst).Find(asstCond).All(ctx, &assts)
		Expect(err).NotTo(HaveOccurred())

		for _, asst := range assts {
			if !common.IsInnerModel(asst.ObjectID) {
				rsp, err := objectClient.DeleteModel(context.Background(), asst.ObjectID, header)
				Expect(err).NotTo(HaveOccurred())
				Expect(rsp.Result).To(Equal(true))
			}
		}

		biz := new(metadata.BizInst)
		err = test.GetDB().Table(common.BKTableNameBaseApp).Find(mapstr.MapStr{"bk_biz_name": "set_Template_biz"}).
			One(ctx, biz)
		if test.GetDB().IsNotFoundError(err) {
			return
		}
		Expect(err).NotTo(HaveOccurred())

		delCond := mapstr.MapStr{common.BKAppIDField: biz.BizID}
		err = test.GetDB().Table(common.BKTableNameObjAttDes).Delete(ctx, delCond)
		Expect(err).NotTo(HaveOccurred())

		err = apiServerClient.UpdateBizDataStatus(ctx, "0", common.DataStatusDisabled, biz.BizID, header)
		testutil.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())

		rsp, err := apiServerClient.DeleteBiz(ctx, header, metadata.DeleteBizParam{BizID: []int64{biz.BizID}})
		testutil.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.BaseResp.Result).To(Equal(true))
	}()

	By("create business")
	func() {
		data := map[string]interface{}{
			"bk_biz_name":       "set_Template_biz",
			"life_cycle":        "2",
			"bk_biz_maintainer": "admin",
			"bk_biz_productor":  "",
			"bk_biz_tester":     "",
			"bk_biz_developer":  "",
			"operator":          "",
			"time_zone":         "Asia/Shanghai",
			"language":          "1",
		}
		rsp, err := topoServerClient.Instance().CreateApp(ctx, common.BKDefaultOwnerID, header, data)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.BaseResp.Result).To(Equal(true))
		Expect(rsp.Data[common.BKAppIDField]).To(Not(Equal(int64(0))))
		bizID, err = util.GetInt64ByInterface(rsp.Data[common.BKAppIDField])
		Expect(err).To(BeNil())
	}()

	By("create parent service category")
	func() {
		input := map[string]interface{}{
			"bk_parent_id":      0,
			common.BKAppIDField: bizID,
			"name":              "root0",
		}
		rsp, err := serviceClient.CreateServiceCategory(context.Background(), header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
			"Result": Equal(true),
		}))
		j, err := json.Marshal(rsp.Data)
		data := metadata.ServiceCategory{}
		err = json.Unmarshal(j, &data)
		Expect(err).To(BeNil())
		categoryId = data.ID
	}()

	By("create service sub category")
	func() {
		input := map[string]interface{}{
			"bk_parent_id":      categoryId,
			common.BKAppIDField: bizID,
			"name":              "child0",
		}
		rsp, err := serviceClient.CreateServiceCategory(context.Background(), header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
			"Result": Equal(true),
		}))
		j, err := json.Marshal(rsp.Data)
		data := metadata.ServiceCategory{}
		err = json.Unmarshal(j, &data)
		Expect(err).To(BeNil())
		categoryId = data.ID
	}()

	By("create service template")
	func() {
		data := map[string]interface{}{
			"bk_biz_id":           bizID,
			"name":                "svcTpl1",
			"service_category_id": categoryId,
		}
		rsp, err := procServerClient.Service().CreateServiceTemplate(ctx, header, data)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
			"Result": Equal(true),
		}))
		Expect(rsp.Data[common.BKFieldID]).To(Not(Equal(int64(0))))
		Expect(rsp.Data[common.BKFieldName]).To(Equal("svcTpl1"))
		serviceTemplateID, err = util.GetInt64ByInterface(rsp.Data[common.BKFieldID])
		Expect(err).To(BeNil())
	}()

	By("create service template 2")
	func() {
		data := map[string]interface{}{
			"bk_biz_id":           bizID,
			"name":                "svcTpl2",
			"service_category_id": categoryId,
		}
		rsp, err := procServerClient.Service().CreateServiceTemplate(ctx, header, data)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
			"Result": Equal(true),
		}))
		Expect(rsp.Data[common.BKFieldID]).To(Not(Equal(int64(0))))
		Expect(rsp.Data[common.BKFieldName]).To(Equal("svcTpl2"))
		serviceTemplateID2, err = util.GetInt64ByInterface(rsp.Data[common.BKFieldID])
		Expect(err).To(BeNil())
	}()

	By("create service template 3")
	func() {
		data := map[string]interface{}{
			"bk_biz_id":           bizID,
			"name":                "svcTpl3",
			"service_category_id": categoryId,
		}
		rsp, err := procServerClient.Service().CreateServiceTemplate(ctx, header, data)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
			"Result": Equal(true),
		}))
		Expect(rsp.Data[common.BKFieldID]).To(Not(Equal(int64(0))))
		Expect(rsp.Data[common.BKFieldName]).To(Equal("svcTpl3"))
		serviceTemplateID3, err = util.GetInt64ByInterface(rsp.Data[common.BKFieldID])
		Expect(err).To(BeNil())
	}()
}
