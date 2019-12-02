package topo_server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

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
	var bizID int64
	var categoryId int64

	It("normal set template test", func() {
		By("clear database")
		test.ClearDatabase()

		By("create business")
		func() {
			data := map[string]interface{}{
				"bk_biz_name":       "biz3",
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
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
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

		var serviceTemplateID int64
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

		var serviceTemplateID2 int64
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

		var serviceTemplateID3 int64
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
			bizIDStr := strconv.FormatInt(bizID, 10)
			rsp, err := topoServerClient.Instance().CreateSet(ctx, bizIDStr, header, data)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
			Expect(rsp.Data[common.BKSetNameField]).To(Equal("set1"))
			Expect(rsp.Data[common.BKSetIDField]).To(Not(Equal(int64(0))))
			setID, err = util.GetInt64ByInterface(rsp.Data[common.BKSetIDField])
			Expect(err).To(BeNil())

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
				SetIDs: []int64{setID},
			}
			setTplDiffResult, err := topoServerClient.SetTemplate().DiffSetTplWithInst(ctx, header, bizID, setTemplateID, option)
			Expect(err).NotTo(HaveOccurred())
			Expect(setTplDiffResult.Difference).To(HaveLen(1))
			setDiff := setTplDiffResult.Difference[0]
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
			bizIDstr := strconv.FormatInt(bizID, 10)
			setIDstr := strconv.FormatInt(setID, 10)
			resp, err := topoServerClient.Instance().DeleteSet(ctx, bizIDstr, setIDstr, header)
			Expect(err).To(BeNil())
			Expect(resp.BaseResp).To(MatchFields(IgnoreExtras, Fields{
				"Result": Equal(true),
			}))
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
