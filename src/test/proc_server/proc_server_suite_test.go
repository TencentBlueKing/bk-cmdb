package proc_server_test

import (
	"context"
	"strconv"
	"testing"

	"configcenter/src/common/mapstr"
	params "configcenter/src/common/paraparse"
	commonutil "configcenter/src/common/util"
	"configcenter/src/test"
	"configcenter/src/test/reporter"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var clientSet = test.GetClientSet()
var serviceClient = clientSet.ProcServer().Service()
var processClient = clientSet.ProcServer().Process()
var instClient = test.GetClientSet().TopoServer().Instance()
var hostServerClient = test.GetClientSet().HostServer()
var apiServerClient = test.GetClientSet().ApiServer()

var bizId, hostId1, hostId2, setId int64

func TestProcServer(t *testing.T) {
	RegisterFailHandler(util.Fail)
	reporters := []Reporter{
		reporter.NewHtmlReporter(test.GetReportDir()+"procserver.html", test.GetReportUrl(), true),
	}
	RunSpecsWithDefaultAndCustomReporters(t, "ProcServer Suite", reporters)
}

var _ = BeforeSuite(func() {
	test.ClearDatabase()

	Describe("test preparation", func() {
		Describe("create biz", func() {
			input := map[string]interface{}{
				"life_cycle":        "2",
				"language":          "1",
				"bk_biz_maintainer": "admin",
				"bk_biz_name":       "cc_biz",
				"time_zone":         "Africa/Accra",
			}
			rsp, err := apiServerClient.CreateBiz(context.Background(), "0", header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			bizId, err = commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("add host", func() {
			input := map[string]interface{}{
				"bk_biz_id": bizId,
				"host_info": map[string]interface{}{
					"1": map[string]interface{}{
						"bk_host_innerip": "127.0.0.1",
						"bk_asset_id":     "addhost_api_asset_1",
						"bk_cloud_id":     0,
					},
					"2": map[string]interface{}{
						"bk_host_innerip": "127.0.0.2",
						"bk_asset_id":     "addhost_api_asset_2",
						"bk_cloud_id":     0,
					},
				},
			}
			rsp, err := hostServerClient.AddHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		Describe("search host", func() {
			input := &params.HostCommonSearch{
				AppID: int(bizId),
			}
			rsp, err := hostServerClient.SearchHost(context.Background(), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(2))
			hostId1, err = commonutil.GetInt64ByInterface(rsp.Data.Info[0]["host"].(map[string]interface{})["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
			hostId2, err = commonutil.GetInt64ByInterface(rsp.Data.Info[1]["host"].(map[string]interface{})["bk_host_id"])
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("create set", func() {
			input := mapstr.MapStr{
				"bk_set_name":         "test",
				"bk_parent_id":        bizId,
				"bk_supplier_account": "0",
				"bk_biz_id":           bizId,
				"bk_service_status":   "1",
				"bk_set_env":          "3",
			}
			rsp, err := instClient.CreateSet(context.Background(), strconv.FormatInt(bizId, 10), header, input)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data["bk_set_name"].(string)).To(Equal("test"))
			parentIdRes, err := commonutil.GetInt64ByInterface(rsp.Data["bk_parent_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(parentIdRes).To(Equal(bizId))
			bizIdRes, err := commonutil.GetInt64ByInterface(rsp.Data["bk_biz_id"])
			Expect(err).NotTo(HaveOccurred())
			Expect(bizIdRes).To(Equal(bizId))
			setId, err = commonutil.GetInt64ByInterface(rsp.Data["bk_set_id"])
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
