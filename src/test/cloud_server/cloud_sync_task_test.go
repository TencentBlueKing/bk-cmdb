package cloud_server_test

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testData1 = map[string]interface{}{
	"bk_task_name":     "王者荣耀1",
	"bk_account_id":    1,
	"bk_resource_type": "host",
	"bk_sync_all":      false,
	"bk_sync_vpcs": []map[string]interface{}{
		{
			"bk_vpc_id":     "vpc-001",
			"bk_vpc_name":   "vpc-default",
			"bk_region": "广东一区",
			"bk_host_count": 56,
			"bk_sync_dir":   33,
		},
		{
			"bk_vpc_id":     "vpc-002",
			"bk_vpc_name":   "vpc-default2",
			"bk_region": "广东二区",
			"bk_host_count": 26,
			"bk_sync_dir":   55,
		},
	},
}

var testData3 = map[string]interface{}{
	"bk_task_name":     "王者荣耀132",
	"bk_account_id":    2,
	"bk_resource_type": "host",
	"bk_sync_all":      false,
	"bk_sync_vpcs": []map[string]interface{}{
		{
			"bk_vpc_name":   "vpc-default",
			"bk_region": "广东一区",
			"bk_host_count": 56,
			"bk_sync_dir":   33,
		},
	},
}

var testData2 = map[string]interface{}{
	"bk_task_name":     "王者荣耀2",
	"bk_account_id":    2,
	"bk_resource_type": "host",
	"bk_sync_all":      true,
	"bk_sync_all_dir":  55,
	"bk_sync_vpcs":     []string{},
}

var tmpData = map[string]interface{}{
	"bk_task_name":     "王者荣耀23",
	"bk_account_id":    2,
	"bk_resource_type": "host",
	"bk_sync_all":      true,
	"bk_sync_all_dir":  55,
	"bk_sync_vpcs":     []string{},
}

// 清除表数据，保证测试用例之间互不干扰
func clearSyncTaskData() {
	err := test.GetDB().DropTable(context.Background(), common.BKTableNameCloudSyncTask)
	Expect(err).NotTo(HaveOccurred())

	err = test.GetDB().Table(common.BKTableNameIDgenerator).Delete(context.Background(), map[string]interface{}{"_id": common.BKTableNameCloudSyncTask})
	Expect(err).NotTo(HaveOccurred())
}

// 准备测试用例需要的数据
func prepareSyncTaskData() {
	accountData := []map[string]interface{}{testData1, testData2}
	for _, data := range accountData {
		rsp, err := cloudServerClient.CreateSyncTask(context.Background(), header, data)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	}
}

var _ = Describe("cloud sync task test", func() {

	BeforeEach(func() {
		//清空数据
		clearSyncTaskData()
		//准备数据
		prepareSyncTaskData()

		// 准备需要云账户数据
		clearData()
		prepareData()
	})

	var _ = Describe("create cloud sync task test", func() {

		It("create task with normal data", func() {
			clearSyncTaskData()
			rsp, err := cloudServerClient.CreateSyncTask(context.Background(), header, tmpData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("create with task name which is already exist", func() {
			tmpTask := tmpData
			tmpTask["bk_task_name"] = testData1["bk_task_name"]
			rsp, err := cloudServerClient.CreateSyncTask(context.Background(), header, tmpTask)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(common.CCErrCloudSyncTaskNameAlreadyExist))
		})

		It("create task with vpc but without vpcID", func() {
			rsp, err := cloudServerClient.CreateSyncTask(context.Background(), header, testData3)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(common.CCErrCloudVpcIDIsRequired))
		})

		It("create task with invalid accountID", func() {
			data := tmpData
			data["bk_task_name"] = "hello world"
			data["bk_account_id"] = int64(999)
			rsp, err := cloudServerClient.CreateSyncTask(context.Background(), header, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(common.CCErrCloudValidSyncTaskParamFail))
		})

	})

	var _ = Describe("update cloud sync task test", func() {

		It("update with normal data", func() {
			taskID := int64(1)
			data := map[string]interface{}{"bk_task_name": "你好啊，雷猴啊"}
			rsp, err := cloudServerClient.UpdateSyncTask(context.Background(), header, taskID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update with task name which is already exist", func() {
			taskID := int64(1)
			data := map[string]interface{}{"bk_task_name": testData2["bk_task_name"]}
			rsp, err := cloudServerClient.UpdateSyncTask(context.Background(), header, taskID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(common.CCErrCloudSyncTaskNameAlreadyExist))
		})

		It("update with invalid accountID", func() {
			taskID := int64(1)
			data := map[string]interface{}{"bk_account_id": int64(999)}
			rsp, err := cloudServerClient.UpdateSyncTask(context.Background(), header, taskID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(common.CCErrCloudValidSyncTaskParamFail))
		})

	})

	var _ = Describe("delete cloud sync task test", func() {

		It("delete with normal data", func() {
			accountID := int64(1)
			rsp, err := cloudServerClient.DeleteSyncTask(context.Background(), header, accountID)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

	})

	var _ = Describe("search cloud sync task test", func() {

		It("search with default query condition", func() {
			cond := make(map[string]interface{})
			rsp, err := cloudServerClient.SearchSyncTask(context.Background(), header, cond)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(2)))
		})

		It("search with configured condition", func() {
			queryData := map[string]interface{}{"condition": map[string]interface{}{"bk_task_name": testData1["bk_task_name"]}}
			rsp, err := cloudServerClient.SearchSyncTask(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(1)))
			Expect(rsp.Data.Info[0].String("bk_task_name")).To(Equal(testData1["bk_task_name"]))
		})

		It("search with configured sort", func() {
			queryData := map[string]interface{}{"page": map[string]interface{}{"sort": "bk_task_name"}}
			rsp, err := cloudServerClient.SearchSyncTask(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(2)))
			Expect(rsp.Data.Info[0].String("bk_task_name")).To(Equal(testData1["bk_task_name"]))
			Expect(rsp.Data.Info[1].String("bk_task_name")).To(Equal(testData2["bk_task_name"]))
		})

		It("search with configured limit", func() {
			queryData := map[string]interface{}{"page": map[string]interface{}{"limit": 1}}
			rsp, err := cloudServerClient.SearchSyncTask(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(len(rsp.Data.Info)).To(Equal(1))
		})

		It("search with configured exact is true", func() {
			queryData := map[string]interface{}{"exact": true, "condition": map[string]interface{}{"bk_task_name": "王者荣耀"}}
			rsp, err := cloudServerClient.SearchSyncTask(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(0)))

			queryData = map[string]interface{}{"exact": true, "condition": map[string]interface{}{"bk_task_name": testData1["bk_task_name"]}}
			rsp, err = cloudServerClient.SearchSyncTask(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(1)))
			Expect(rsp.Data.Info[0].String("bk_task_name")).To(Equal(testData1["bk_task_name"]))
		})

		It("search with configured exact is false", func() {
			queryData := map[string]interface{}{"exact": false, "condition": map[string]interface{}{"bk_task_name": "王者荣耀"}}
			rsp, err := cloudServerClient.SearchSyncTask(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(2)))
		})

	})
})
