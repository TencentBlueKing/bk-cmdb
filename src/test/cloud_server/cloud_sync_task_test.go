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
			"bk_region":     "广东一区",
			"bk_host_count": 56,
			"bk_sync_dir":   1,
			"bk_cloud_id":   1,
		},
		{
			"bk_vpc_id":     "vpc-002",
			"bk_vpc_name":   "vpc-default2",
			"bk_region":     "广东二区",
			"bk_host_count": 26,
			"bk_sync_dir":   1,
			"bk_cloud_id":   1,
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
			"bk_region":     "广东一区",
			"bk_host_count": 56,
			"bk_sync_dir":   1,
		},
	},
}

var testData2 = map[string]interface{}{
	"bk_task_name":     "王者荣耀2",
	"bk_account_id":    2,
	"bk_resource_type": "host",
	"bk_sync_all":      true,
	"bk_sync_all_dir":  1,
	"bk_sync_vpcs":     []string{},
}

var tmpData = map[string]interface{}{
	"bk_task_name":     "王者荣耀23",
	"bk_account_id":    2,
	"bk_resource_type": "host",
	"bk_sync_all":      true,
	"bk_sync_all_dir":  1,
	"bk_sync_vpcs":     []string{},
}

// 清除表数据，保证测试用例之间互不干扰
func clearSyncTaskData() {
	// 清空云同步任务表
	err := test.GetDB().Table(common.BKTableNameCloudSyncTask).Delete(context.Background(), map[string]interface{}{})
	Expect(err).NotTo(HaveOccurred())

	//删除云同步任务id计数
	err = test.GetDB().Table(common.BKTableNameIDgenerator).Delete(context.Background(), map[string]interface{}{"_id": common.BKTableNameCloudSyncTask})
	Expect(err).NotTo(HaveOccurred())
}

// 准备测试用例需要的数据
func prepareSyncTaskData() {
	taskData := []map[string]interface{}{testData1, testData2}
	for _, data := range taskData {
		rsp, err := cloudServerClient.CreateSyncTask(context.Background(), header, data)
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	}
}

var cloudID1 int64

func prepareCloudData() {
	//清空云区域表
	err := test.GetDB().Table(common.BKTableNameBasePlat).Delete(context.Background(), map[string]interface{}{})
	Expect(err).NotTo(HaveOccurred())

	//删除云区域id计数
	err = test.GetDB().Table(common.BKTableNameIDgenerator).Delete(context.Background(), map[string]interface{}{"_id": common.BKTableNameBasePlat})
	Expect(err).NotTo(HaveOccurred())

	// 准备数据
	resp, err := hostServerClient.CreateCloudArea(context.Background(), header, map[string]interface{}{
		"bk_cloud_name":   "LPL17区",
		"bk_status":       "1",
		"bk_cloud_vendor": "1",
		"bk_account_id":   2,
		"creator":         "admin",
	})
	cloudID1 = int64(resp.Data.Created.ID)
	util.RegisterResponse(resp)
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.Result).To(Equal(true))
}

var _ = Describe("cloud sync task test", func() {

	BeforeEach(func() {
		// 准备云账户数据
		clearAccountData()
		prepareAccountData()

		// 准备云区域数据
		prepareCloudData()

		//清空同步任务数据
		clearSyncTaskData()
		//准备同步任务数据
		prepareSyncTaskData()
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

		It("search with configured is_fuzzy is false", func() {
			queryData := map[string]interface{}{"is_fuzzy": false, "condition": map[string]interface{}{"bk_task_name": "王者荣耀"}}
			rsp, err := cloudServerClient.SearchSyncTask(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(0)))

			queryData = map[string]interface{}{"is_fuzzy": false, "condition": map[string]interface{}{"bk_task_name": testData1["bk_task_name"]}}
			rsp, err = cloudServerClient.SearchSyncTask(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(1)))
			Expect(rsp.Data.Info[0].String("bk_task_name")).To(Equal(testData1["bk_task_name"]))
		})

		It("search with configured is_fuzzy is true", func() {
			queryData := map[string]interface{}{"is_fuzzy": true, "condition": map[string]interface{}{"bk_task_name": "王者荣耀"}}
			rsp, err := cloudServerClient.SearchSyncTask(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(2)))
		})

	})
})
