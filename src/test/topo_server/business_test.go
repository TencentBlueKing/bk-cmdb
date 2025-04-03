package topo_server_test

import (
	"context"
	"fmt"
	"strconv"

	"configcenter/src/common"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("business test", func() {
	var bizId, bizId2, bizID string
	var bizIdInt, bizId1, bizID2, hostID1 int64

	It("create business bk_biz_name = 'eereeede' for tenant system", func() {
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
		rsp, err := apiServerClient.CreateBiz(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(ContainElement("eereeede"))
		bizId = commonutil.GetStrByInterface(rsp.Data["bk_biz_id"])
	})

	It("biz test for multi-tenant", func() {
		// create biz for tenant test
		testHeader := test.GetTestTenantHeader()
		input := map[string]interface{}{
			"life_cycle":        "2",
			"language":          "1",
			"bk_biz_maintainer": "admin",
			"bk_biz_productor":  "",
			"bk_biz_tester":     "",
			"bk_biz_developer":  "",
			"operator":          "",
			"bk_biz_name":       "biz_test",
			"time_zone":         "Africa/Accra",
		}
		rsp, err := apiServerClient.CreateBiz(context.Background(), testHeader, input)
		util.RegisterResponseWithRid(rsp, test.GetTestTenantHeader())
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(ContainElement("biz_test"))
		testBizId, err := commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
		Expect(err).To(BeNil())
		// search biz for tenant test
		searchInput := &params.SearchParams{
			Condition: map[string]interface{}{
				"bk_biz_name": "eereeede",
			},
		}
		searchRsp, err := instClient.SearchApp(context.Background(), testHeader, searchInput)
		util.RegisterResponseWithRid(rsp, testHeader)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Data.Count).To(Equal(0))

		searchInput = &params.SearchParams{
			Condition: map[string]interface{}{
				"bk_biz_name": "biz_test",
			},
		}
		searchRsp, err = instClient.SearchApp(context.Background(), header, searchInput)
		util.RegisterResponseWithRid(rsp, testHeader)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Data.Count).To(Equal(0))

		searchInput = &params.SearchParams{
			Condition: map[string]interface{}{
				"bk_biz_name": "biz_test",
			},
		}
		searchRsp, err = instClient.SearchApp(context.Background(), testHeader, searchInput)
		util.RegisterResponseWithRid(rsp, testHeader)
		Expect(err).NotTo(HaveOccurred())
		Expect(searchRsp.Data.Count).To(Equal(1))
		id, err := commonutil.GetInt64ByInterface(searchRsp.Data.Info[0]["bk_biz_id"])
		Expect(err).To(BeNil())
		Expect(id).To(Equal(testBizId))
	})

	It("create business bk_biz_name = 'test_biz_status', set and module", func() {
		input := map[string]interface{}{
			"life_cycle":        "2",
			"language":          "1",
			"bk_biz_maintainer": "admin",
			"bk_biz_productor":  "",
			"bk_biz_tester":     "",
			"bk_biz_developer":  "",
			"operator":          "",
			"bk_biz_name":       "test_biz_status",
			"time_zone":         "Africa/Accra",
		}
		rsp, err := apiServerClient.CreateBiz(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(ContainElement("test_biz_status"))
		bizID = commonutil.GetStrByInterface(rsp.Data["bk_biz_id"])
		bizIdStr := commonutil.GetStrByInterface(bizID)
		bizId1, err = strconv.ParseInt(bizIdStr, 10, 64)
	})

	It("add host and get host id", func() {
		cloudID := test.GetCloudID()
		hostInput := map[string]interface{}{
			"bk_biz_id": bizId1,
			"host_info": map[string]interface{}{
				"1": map[string]interface{}{
					"bk_host_innerip": "127.0.0.8",
					"bk_cloud_id":     cloudID,
				},
			},
		}
		hostRsp, err := hostServerClient.AddHost(context.Background(), header, hostInput)
		util.RegisterResponseWithRid(hostRsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(hostRsp.Result).To(Equal(true), hostRsp.ToString())

		searchInput := &metadata.HostCommonSearch{
			AppID: bizId1,
		}
		resp, err := hostServerClient.SearchHostWithBiz(context.Background(), header, searchInput)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Result).To(Equal(true))
		Expect(resp.Data.Count).To(Equal(1))
		hostID1, err = commonutil.GetInt64ByInterface(resp.Data.Info[0]["host"].(map[string]interface{})["bk_host_id"])
		Expect(err).NotTo(HaveOccurred())
	})

	It("change biz status disable with host", func() {
		input := &metadata.UpdateBusinessStatusOption{
			BizName: "test_biz_status",
		}
		err := apiServerClient.UpdateBusinessStatus(context.Background(), string(common.DataStatusDisabled), bizId1,
			header, input)
		util.RegisterResponseWithRid(err, header)
		Expect(err).To(HaveOccurred())
	})

	It(fmt.Sprintf("transfer host to resource pool, biz: %d", bizId1), func() {
		transferData := &metadata.DefaultModuleHostConfigParams{
			ApplicationID: bizId1,
			HostIDs:       []int64{hostID1},
		}
		resp, transferErr := hostServerClient.MoveHostToResourcePool(context.Background(), header, transferData)
		util.RegisterResponseWithRid(resp, header)
		Expect(transferErr).NotTo(HaveOccurred())
		Expect(resp.CCError()).NotTo(HaveOccurred())
	})

	It("change biz status disable", func() {
		input := &metadata.UpdateBusinessStatusOption{
			BizName: "eereeede",
		}
		var err error
		bizID2, err = strconv.ParseInt(bizId, 10, 64)
		Expect(err).To(BeNil())
		err = apiServerClient.UpdateBusinessStatus(context.Background(), string(common.DataStatusDisabled), bizID2,
			header, input)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
	})

	It("check change biz status disable", func() {
		input := &metadata.QueryCondition{
			Condition: mapstr.MapStr{
				"bk_biz_name": "eereeede",
			},
		}
		resp, err := apiServerClient.ReadInstance(context.Background(), header, "biz", input)
		util.RegisterResponseWithRid(resp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Data.Count).To(Equal(0))
	})

	It("change biz status enable", func() {
		input := &metadata.UpdateBusinessStatusOption{
			BizName: "eereeede",
		}
		err := apiServerClient.UpdateBusinessStatus(context.Background(), string(common.DataStatusEnable), bizID2,
			header, input)
		Expect(err).NotTo(HaveOccurred())
	})

	It("check change biz status enable", func() {
		input := &metadata.QueryCondition{
			Condition: mapstr.MapStr{
				"bk_biz_name": "eereeede",
			},
		}
		resp, err := apiServerClient.ReadInstance(context.Background(), header, "biz", input)
		util.RegisterResponseWithRid(resp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Data.Count).To(Equal(1))
		Expect(resp.Data.Info[0]["bk_biz_name"]).To(Equal("eereeede"))
		id, err := commonutil.GetInt64ByInterface(resp.Data.Info[0]["bk_biz_id"])
		Expect(err).To(BeNil())
		Expect(id).To(Equal(bizID2))
	})

	It("change biz status disable", func() {
		input := &metadata.UpdateBusinessStatusOption{
			BizName: "test_biz_status",
		}
		err := apiServerClient.UpdateBusinessStatus(context.Background(), string(common.DataStatusDisabled), bizId1,
			header, input)
		Expect(err).NotTo(HaveOccurred())
	})

	It(fmt.Sprintf("delete business bk_biz_id = %d", bizId1), func() {
		input := metadata.DeleteBizParam{
			BizID: []int64{bizId1},
		}
		err := apiServerClient.DeleteBiz(context.Background(), header, input)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
	})

	It("search business basic info", func() {
		bizID, err := strconv.ParseInt(bizId, 10, 64)
		rsp, err := instClient.GetAppBasicInfo(context.Background(), header, bizID)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Data.BizName).To(Equal("eereeede"))
		Expect(rsp.Data.BizID).To(Equal(bizID))
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
		rsp, err := apiServerClient.CreateBiz(context.Background(), header, input)
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
		rsp, err := apiServerClient.CreateBiz(context.Background(), header, input)
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
		rsp, err := apiServerClient.CreateBiz(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(ContainElement("mmrmm"))
		bizIdInt, err = commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
		Expect(err).NotTo(HaveOccurred())
		bizId2 = strconv.FormatInt(bizIdInt, 10)
	})

	It("search resource business", func() {
		rsp, err := instClient.GetDefaultApp(context.Background(), header)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Data.Count).To(Equal(1))
		Expect(rsp.Data.Info[0]).To(ContainElement("资源池"))
	})

	It("search business bk_biz_name = 'mmrmm'", func() {
		input := &params.SearchParams{
			Condition: map[string]interface{}{
				"bk_biz_name": "mmrmm",
			},
		}
		rsp, err := instClient.SearchApp(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Data.Count).To(Equal(1))
		Expect(rsp.Data.Info[0]["bk_biz_name"]).To(Equal("mmrmm"))
		id, err := commonutil.GetInt64ByInterface(rsp.Data.Info[0]["bk_biz_id"])
		Expect(err).NotTo(HaveOccurred())
		Expect(id).To(Equal(bizIdInt))
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
		rsp, err := apiServerClient.SearchBiz(context.Background(), header, input)
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
		rsp, err := apiServerClient.SearchBiz(context.Background(), header, input)
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
		rsp, err := apiServerClient.SearchBiz(context.Background(), header, input)
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
		rsp, err := apiServerClient.SearchBiz(context.Background(), header, input)
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
		rsp, err := apiServerClient.SearchBiz(context.Background(), header, input)
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
		rsp, err := apiServerClient.SearchBiz(context.Background(), header, input)
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
		rsp, err := apiServerClient.UpdateBiz(context.Background(), bizId, header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("update nonexist business", func() {
		input := map[string]interface{}{
			"bk_biz_name": "cdewdercfee",
			"life_cycle":  "2",
		}
		rsp, err := apiServerClient.UpdateBiz(context.Background(), "1000", header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).Should(BeNil())
		Expect(rsp.Result).To(Equal(false))
	})

	It("update business using exist bk_biz_name", func() {
		input := map[string]interface{}{
			"bk_biz_name": "mmrmm",
			"life_cycle":  "2",
		}
		rsp, err := apiServerClient.UpdateBiz(context.Background(), bizId, header, input)
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
		err := apiServerClient.UpdateBizDataStatus(context.Background(), common.DataStatusDisabled, bizIdInt, header)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
	})

	It("update nonexist business enable status diable", func() {
		err := apiServerClient.UpdateBizDataStatus(context.Background(), common.DataStatusDisabled, 1000, header)
		util.RegisterResponseWithRid(err, header)
		Expect(err).ShouldNot(BeNil())
	})

	It("update nonexist business enable status enable", func() {
		err := apiServerClient.UpdateBizDataStatus(context.Background(), common.DataStatusEnable, 1000, header)
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
		rsp, err := apiServerClient.SearchBiz(context.Background(), header, input)
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
		rsp, err := apiServerClient.SearchBiz(context.Background(), header, input)
		util.RegisterResponseWithRid(rsp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data.Count).To(Equal(1))
		Expect(rsp.Data.Info).To(ContainElement(ContainElement("cdewdercfee")))
	})

	It(fmt.Sprintf("update business enable status bk_biz_id = %s enable", bizId2), func() {
		err := apiServerClient.UpdateBizDataStatus(context.Background(), common.DataStatusEnable, bizIdInt, header)
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
		rsp, err := apiServerClient.SearchBiz(context.Background(), header, input)
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
		err := apiServerClient.UpdateBizDataStatus(context.Background(), common.DataStatusDisabled, bizIdInt, header)
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
