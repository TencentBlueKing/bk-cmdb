package cloud_server_test

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var accountData1 = map[string]interface{}{
	"bk_account_name": "awsAccount1",
	"bk_cloud_vendor": "1",
	"bk_secret_id":    "aaaaa",
	"bk_secret_key":   "bbbbb",
	"bk_description":  "aws账户1",
	"bk_creator":      "admin",
}

var accountData2 = map[string]interface{}{
	"bk_account_name": "tcAccount1",
	"bk_cloud_vendor": "2",
	"bk_secret_id":    "ccccc",
	"bk_secret_key":   "ddddd",
	"bk_description":  "腾讯云账户1",
	"bk_creator":      "admin",
}

func NewTmpAccount() map[string]interface{} {
	return map[string]interface{}{
		"bk_account_name": "tmpAccount",
		"bk_cloud_vendor": "2",
		"bk_secret_id":    "eeeee",
		"bk_secret_key":   "fffff",
		"bk_description":  "腾讯云临时账户",
		"bk_creator":      "admin",
	}
}

// 清除表数据，保证测试用例之间互不干扰
func clearAccountData() {
	//清空云账户表
	err := test.GetDB().Table(common.BKTableNameCloudAccount).Delete(context.Background(), map[string]interface{}{})
	Expect(err).NotTo(HaveOccurred())

	//删除云账户id计数
	err = test.GetDB().Table(common.BKTableNameIDgenerator).Delete(context.Background(), map[string]interface{}{"_id": common.BKTableNameCloudAccount})
	Expect(err).NotTo(HaveOccurred())
}

// 准备测试用例需要的数据
func prepareAccountData() {
	accountData := []map[string]interface{}{accountData1, accountData2}
	for i := range accountData {
		rsp, err := cloudServerClient.CreateAccount(context.Background(), header, accountData[i])
		util.RegisterResponse(rsp)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	}
}

var _ = Describe("cloud account test", func() {

	BeforeEach(func() {
		//清空数据
		clearAccountData()
		clearSyncTaskData()
		//准备数据
		prepareAccountData()
	})

	var _ = Describe("create cloud account test", func() {

		It("create with normal data", func() {
			tmpAccount := NewTmpAccount()
			rsp, err := cloudServerClient.CreateAccount(context.Background(), header, tmpAccount)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("create with cloud name which is already exist", func() {
			tmpAccount := NewTmpAccount()
			tmpAccount["bk_account_name"] = accountData1["bk_account_name"]
			rsp, err := cloudServerClient.CreateAccount(context.Background(), header, tmpAccount)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(common.CCErrCloudAccountNameAlreadyExist))
		})

		It("create with cloud vendor which is not valid", func() {
			tmpAccount := NewTmpAccount()
			tmpAccount["bk_cloud_vendor"] = "aaa"
			rsp, err := cloudServerClient.CreateAccount(context.Background(), header, tmpAccount)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(common.CCErrCloudVendorNotSupport))

		})

	})

	var _ = Describe("update cloud account test", func() {

		It("update with normal data", func() {
			accountID := int64(1)
			data := map[string]interface{}{"bk_account_name": "Jack"}
			rsp, err := cloudServerClient.UpdateAccount(context.Background(), header, accountID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update with cloud name which is already exist", func() {
			accountID := int64(1)
			data := map[string]interface{}{"bk_account_name": accountData2["bk_account_name"]}
			rsp, err := cloudServerClient.UpdateAccount(context.Background(), header, accountID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(common.CCErrCloudAccountNameAlreadyExist))
		})

		It("update with cloud vendor which is not valid", func() {
			accountID := int64(1)
			data := map[string]interface{}{"bk_account_name": accountData2["bk_account_name"]}
			rsp, err := cloudServerClient.UpdateAccount(context.Background(), header, accountID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(common.CCErrCloudAccountNameAlreadyExist))
		})

		It("update with cloud accountID which is not exist", func() {
			accountID := int64(99999)
			data := map[string]interface{}{"bk_account_name": "Jack"}
			rsp, err := cloudServerClient.UpdateAccount(context.Background(), header, accountID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
			Expect(rsp.Code).To(Equal(common.CCErrCloudAccountIDNoExistFail))
		})

	})

	var _ = Describe("delete cloud account test", func() {

		It("delete with normal data", func() {
			accountID := int64(1)
			rsp, err := cloudServerClient.DeleteAccount(context.Background(), header, accountID)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("delete with cloud accountID which is not exist", func() {
			accountID := int64(99999)
			rsp, err := cloudServerClient.DeleteAccount(context.Background(), header, accountID)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

	})

	var _ = Describe("search cloud account test", func() {

		It("search with default query condition", func() {
			rsp, err := cloudServerClient.SearchAccount(context.Background(), header, nil)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(2)))
		})

		It("search with configured conditon", func() {
			queryData := map[string]interface{}{"condition": map[string]interface{}{"bk_account_name": accountData1["bk_account_name"]}}
			rsp, err := cloudServerClient.SearchAccount(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(1)))
			Expect(rsp.Data.Info[0].String("bk_account_name")).To(Equal(accountData1["bk_account_name"]))
		})

		It("search with configured sort", func() {
			queryData := map[string]interface{}{"page": map[string]interface{}{"sort": "bk_account_name"}}
			rsp, err := cloudServerClient.SearchAccount(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(2)))
			Expect(rsp.Data.Info[0].String("bk_account_name")).To(Equal(accountData1["bk_account_name"]))
			Expect(rsp.Data.Info[1].String("bk_account_name")).To(Equal(accountData2["bk_account_name"]))
		})

		It("search with configured limit", func() {
			queryData := map[string]interface{}{"page": map[string]interface{}{"limit": 1}}
			rsp, err := cloudServerClient.SearchAccount(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(len(rsp.Data.Info)).To(Equal(1))
		})

		It("search with configured is_fuzzy is false", func() {
			queryData := map[string]interface{}{"is_fuzzy": false, "condition": map[string]interface{}{"bk_account_name": "aws"}}
			rsp, err := cloudServerClient.SearchAccount(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(0)))

			queryData = map[string]interface{}{"is_fuzzy": false, "condition": map[string]interface{}{"bk_account_name": "awsAccount1"}}
			rsp, err = cloudServerClient.SearchAccount(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(1)))
			Expect(rsp.Data.Info[0].String("bk_account_name")).To(Equal("awsAccount1"))
		})

		It("search with configured is_fuzzy is true", func() {
			queryData := map[string]interface{}{"is_fuzzy": true, "condition": map[string]interface{}{"bk_account_name": "aws"}}
			rsp, err := cloudServerClient.SearchAccount(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(1)))
			Expect(rsp.Data.Info[0].String("bk_account_name")).To(Equal("awsAccount1"))
		})

	})
})
