package topo_server_test

import (
	"context"
	"fmt"
	"strconv"

	"configcenter/src/common"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("business test", func() {
	var bizId, bizId2 string
	var bizIdInt int64

	It("create business bk_biz_name = 'eereeede'", func() {
		test.DeleteAllBizs()
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
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(ContainElement("eereeede"))
		bizId = commonutil.GetStrByInterface(rsp.Data["bk_biz_id"])
	})

	It("create business bk_biz_name = 'eereeede' again", func() {
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
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).Should(BeNil())
		Expect(rsp.Result).To(Equal(false))
	})

	It("create business less bk_biz_name", func() {
		input := map[string]interface{}{
			"life_cycle":        "2",
			"language":          "1",
			"bk_biz_maintainer": "admin",
			"bk_biz_productor":  "",
			"bk_biz_tester":     "",
			"bk_biz_developer":  "",
			"operator":          "",
			"time_zone":         "Africa/Accra",
		}
		rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).Should(BeNil())
		Expect(rsp.Result).To(Equal(false))
	})

	It("create business bk_biz_name = 'mmrmm'", func() {
		input := map[string]interface{}{
			"life_cycle":        "1",
			"language":          "1",
			"bk_biz_maintainer": "admin",
			"bk_biz_productor":  "",
			"bk_biz_tester":     "admin",
			"bk_biz_developer":  "",
			"operator":          "",
			"bk_biz_name":       "mmrmm",
			"time_zone":         "Africa/Accra",
		}
		rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(ContainElement("mmrmm"))
		bizIdInt, err = commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
		Expect(err).NotTo(HaveOccurred())
		bizId2 = strconv.FormatInt(bizIdInt, 10)
	})

	It("search business change start limit", func() {
		input := &metadata.QueryBusinessRequest{
			Page: metadata.BasePage{
				Start: 1,
				Limit: 1,
				Sort:  "bk_biz_id",
			},
			Fields:    []string{},
			Condition: map[string]interface{}{},
		}
		rsp, err := apiServerClient.SearchBiz(context.Background(), "0", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		switch httpheader.GetTenantID(header) {
		case common.BKDefaultTenantID:
			Expect(rsp.Data.Count).To(Equal(3))
		default:
			Expect(rsp.Data.Count).To(Equal(2))
		}
		Expect(len(rsp.Data.Info)).To(Equal(1))
	})

	It("search business using bk_biz_maintainer", func() {
		input := &metadata.QueryBusinessRequest{
			Page: metadata.BasePage{
				Start: 0,
				Limit: 10,
				Sort:  "bk_biz_id",
			},
			Fields: []string{},
			Condition: map[string]interface{}{
				"bk_biz_maintainer": "admin",
			},
		}
		rsp, err := apiServerClient.SearchBiz(context.Background(), "0", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		switch httpheader.GetTenantID(header) {
		case common.BKDefaultTenantID:
			Expect(rsp.Data.Count).To(Equal(3))
		default:
			Expect(rsp.Data.Count).To(Equal(2))
		}
	})

	It("search business using bk_biz_tester", func() {
		input := &metadata.QueryBusinessRequest{
			Page: metadata.BasePage{
				Start: 0,
				Limit: 10,
				Sort:  "bk_biz_id",
			},
			Fields: []string{},
			Condition: map[string]interface{}{
				"bk_biz_tester": "admin",
			},
		}
		rsp, err := apiServerClient.SearchBiz(context.Background(), "0", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data.Count).To(Equal(1))
		Expect(rsp.Data.Info).To(ContainElement(ContainElement("mmrmm")))
	})

	It("search business using time_zone", func() {
		input := &metadata.QueryBusinessRequest{
			Page: metadata.BasePage{
				Start: 0,
				Limit: 10,
				Sort:  "bk_biz_id",
			},
			Fields: []string{},
			Condition: map[string]interface{}{
				"time_zone": "Africa/Accra",
			},
		}
		rsp, err := apiServerClient.SearchBiz(context.Background(), "0", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data.Count).To(Equal(2))
	})

	It("search business using language", func() {
		input := &metadata.QueryBusinessRequest{
			Page: metadata.BasePage{
				Start: 0,
				Limit: 10,
				Sort:  "bk_biz_id",
			},
			Fields: []string{},
			Condition: map[string]interface{}{
				"language": "1",
			},
		}
		rsp, err := apiServerClient.SearchBiz(context.Background(), "0", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		switch httpheader.GetTenantID(header) {
		case common.BKDefaultTenantID:
			Expect(rsp.Data.Count).To(Equal(3))
		default:
			Expect(rsp.Data.Count).To(Equal(2))
		}
	})

	It("search business using life_cycle", func() {
		input := &metadata.QueryBusinessRequest{
			Page: metadata.BasePage{
				Start: 0,
				Limit: 10,
				Sort:  "bk_biz_id",
			},
			Fields: []string{},
			Condition: map[string]interface{}{
				"life_cycle": "1",
			},
		}
		rsp, err := apiServerClient.SearchBiz(context.Background(), "0", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data.Count).To(Equal(1))
		Expect(rsp.Data.Info).To(ContainElement(ContainElement("mmrmm")))
	})

	It(fmt.Sprintf("update business bk_biz_id = %s", bizId), func() {
		input := map[string]interface{}{
			"bk_biz_name": "cdewdercfee",
			"life_cycle":  "2",
		}
		rsp, err := apiServerClient.UpdateBiz(context.Background(), "0", bizId, header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("update nonexist business", func() {
		input := map[string]interface{}{
			"bk_biz_name": "cdewdercfee",
			"life_cycle":  "2",
		}
		rsp, err := apiServerClient.UpdateBiz(context.Background(), "0", "1000", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).Should(BeNil())
		Expect(rsp.Result).To(Equal(false))
	})

	It("update business using exist bk_biz_name", func() {
		input := map[string]interface{}{
			"bk_biz_name": "mmrmm",
			"life_cycle":  "2",
		}
		rsp, err := apiServerClient.UpdateBiz(context.Background(), "0", bizId, header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).Should(BeNil())
		Expect(rsp.Result).To(Equal(false))
	})

	It(fmt.Sprintf("batch update business properties by condition bk_biz_id in [%s]", bizId2), func() {
		bizID, err := strconv.ParseInt(bizId2, 10, 64)
		Expect(err).Should(BeNil())
		input := metadata.UpdateBizPropertyBatchParameter{
			Properties: map[string]interface{}{
				"operator": "test",
			},
			Condition: map[string]interface{}{
				"bk_biz_id": map[string]interface{}{
					"$in": []int64{bizID},
				},
			},
		}

		rsp, err := apiServerClient.UpdateBizPropertyBatch(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It(fmt.Sprintf("batch update business properties by condition bk_biz_id in []"), func() {
		input := metadata.UpdateBizPropertyBatchParameter{
			Properties: map[string]interface{}{
				"operator": "test",
			},
			Condition: map[string]interface{}{
				"bk_biz_id": map[string]interface{}{
					"$in": make([]int64, 0),
				},
			},
		}

		rsp, err := apiServerClient.UpdateBizPropertyBatch(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).Should(BeNil())
		Expect(rsp.Result).To(Equal(false))
	})

	It(fmt.Sprintf("update business enable status bk_biz_id = %s", bizId2), func() {
		err := apiServerClient.UpdateBizDataStatus(context.Background(), "0", common.DataStatusDisabled, bizIdInt,
			header)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
	})

	It("update nonexist business enable status diable", func() {
		err := apiServerClient.UpdateBizDataStatus(context.Background(), "0", common.DataStatusDisabled, 1000, header)
		util.RegisterResponseWithRid(err, header)
		Expect(err).ShouldNot(BeNil())
	})

	It("update nonexist business enable status enable", func() {
		err := apiServerClient.UpdateBizDataStatus(context.Background(), "0", common.DataStatusEnable, 1000, header)
		util.RegisterResponseWithRid(err, header)
		Expect(err).ShouldNot(BeNil())
	})

	It("search business", func() {
		input := &metadata.QueryBusinessRequest{
			Page: metadata.BasePage{
				Start: 0,
				Limit: 10,
				Sort:  "",
			},
			Fields:    []string{},
			Condition: map[string]interface{}{},
		}
		rsp, err := apiServerClient.SearchBiz(context.Background(), "0", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		switch httpheader.GetTenantID(header) {
		case common.BKDefaultTenantID:
			Expect(rsp.Data.Count).To(Equal(2))
		default:
			Expect(rsp.Data.Count).To(Equal(1))
		}
		Expect(rsp.Data.Info).To(ContainElement(ContainElement("cdewdercfee")))
		Expect(rsp.Data.Info).NotTo(ContainElement(ContainElement("eereeede")))
		Expect(rsp.Data.Info).NotTo(ContainElement(ContainElement("mmrmm")))
	})

	It("search business", func() {
		input := &metadata.QueryBusinessRequest{
			Page: metadata.BasePage{
				Start: 0,
				Limit: 10,
				Sort:  "",
			},
			Fields: []string{},
			Condition: map[string]interface{}{
				"bk_biz_name": "cdewdercfee",
			},
		}
		rsp, err := apiServerClient.SearchBiz(context.Background(), "0", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data.Count).To(Equal(1))
		Expect(rsp.Data.Info).To(ContainElement(ContainElement("cdewdercfee")))
	})

	It(fmt.Sprintf("update business enable status bk_biz_id = %s enable", bizId2), func() {
		err := apiServerClient.UpdateBizDataStatus(context.Background(), "0", common.DataStatusEnable, bizIdInt, header)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
	})

	It("search business using bk_biz_id", func() {
		input := &metadata.QueryBusinessRequest{
			Page: metadata.BasePage{
				Start: 0,
				Limit: 10,
				Sort:  "",
			},
			Fields: []string{},
			Condition: map[string]interface{}{
				"bk_biz_id": bizIdInt,
			},
		}
		rsp, err := apiServerClient.SearchBiz(context.Background(), "0", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data.Count).To(Equal(1))
		Expect(rsp.Data.Info).To(ContainElement(ContainElement(ContainSubstring("mmrmm"))))
	})

	It("get brief biz topo", func() {
		input := map[string]interface{}{
			"set_fields": []string{
				"bk_set_id",
				"bk_set_name",
				"bk_set_env",
			},
			"module_fields": []string{
				"bk_module_id",
				"bk_module_name",
			},
			"host_fields": []string{
				"bk_host_id",
				"bk_host_innerip",
				"bk_host_name",
			},
		}
		rsp, err := instClient.SearchBriefBizTopo(context.Background(), header, bizIdInt, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(len(rsp.Data)).To(Equal(1))
		Expect(rsp.Data[0].Set["bk_set_name"]).To(Equal("空闲机池"))
		Expect(len(rsp.Data[0].ModuleTopos)).To(Equal(3))
		modulesMap := map[string]bool{
			"空闲机": true,
			"故障机": true,
			"待回收": true,
		}
		Expect(modulesMap[rsp.Data[0].ModuleTopos[0].Module["bk_module_name"].(string)]).To(Equal(true))
		Expect(modulesMap[rsp.Data[0].ModuleTopos[1].Module["bk_module_name"].(string)]).To(Equal(true))
		Expect(modulesMap[rsp.Data[0].ModuleTopos[2].Module["bk_module_name"].(string)]).To(Equal(true))
	})

	It(fmt.Sprintf("delete unarchived business bk_biz_id = %s", bizId2), func() {
		bizID, err := strconv.ParseInt(bizId2, 10, 64)
		Expect(err).Should(BeNil())
		input := metadata.DeleteBizParam{
			BizID: []int64{bizID},
		}

		err = apiServerClient.DeleteBiz(context.Background(), header, input)
		util.RegisterResponseWithRid(err, header)
		Expect(err).Should(HaveOccurred())
	})

	It(fmt.Sprintf("update business disabled status bk_biz_id = %s", bizId2), func() {
		err := apiServerClient.UpdateBizDataStatus(context.Background(), "0", common.DataStatusDisabled,
			bizIdInt, header)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
	})

	It(fmt.Sprintf("delete archived business bk_biz_id = %s", bizId2), func() {
		bizID, err := strconv.ParseInt(bizId2, 10, 64)
		Expect(err).Should(BeNil())
		input := metadata.DeleteBizParam{
			BizID: []int64{bizID},
		}

		err = apiServerClient.DeleteBiz(context.Background(), header, input)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
	})

	It(fmt.Sprintf("delete default business bk_biz_id = 1"), func() {
		bizResult := new(metadata.BizInst)
		dbErr := test.GetDB().Table(common.BKTableNameBaseApp).Find(mapstr.MapStr{
			common.BKAppNameField: common.DefaultAppName}).Fields(common.BKAppIDField).One(context.Background(),
			bizResult)
		Expect(dbErr).NotTo(HaveOccurred())
		input := metadata.DeleteBizParam{
			BizID: []int64{bizResult.BizID},
		}

		err := apiServerClient.DeleteBiz(context.Background(), header, input)
		util.RegisterResponseWithRid(err, header)
		Expect(err).Should(HaveOccurred())
	})

	It(fmt.Sprintf("delete business in []"), func() {
		input := metadata.DeleteBizParam{
			BizID: make([]int64, 0),
		}

		err := apiServerClient.DeleteBiz(context.Background(), header, input)
		util.RegisterResponseWithRid(err, header)
		Expect(err).Should(HaveOccurred())
	})
})
