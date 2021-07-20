/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package host_server_test

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// testData1
var testData1 = map[string]interface{}{
	"bk_cloud_name":   "LPL17区",
	"bk_status":       "1",
	"bk_cloud_vendor": "1",
	"bk_account_id":   2,
	"creator":         "admin",
}

// testData2
var testData2 = map[string]interface{}{
	"bk_cloud_name":   "LPL29区",
	"bk_status":       "2",
	"bk_cloud_vendor": "2",
	"bk_region":       "1111",
	"bk_vpc_id":       "1",
	"bk_vpc_name":     "Default-VPC",
	"bk_account_id":   2,
	"creator":         "admin",
}

// NewTmpCloudArea new a tmp cloudarea
func NewTmpCloudArea() map[string]interface{} {
	return map[string]interface{}{
		"bk_cloud_name":   "LPL39区",
		"bk_status":       "1",
		"bk_cloud_vendor": "2",
		"bk_region":       "1111",
		"bk_vpc_id":       "tmp_vpc",
		"bk_vpc_name":     "Default-VPC",
		"bk_account_id":   2,
		"creator":         "admin",
	}
}

var _ = Describe("cloud area test", func() {

	BeforeEach(func() {
		//准备数据
		prepareCloudData()
	})

	var _ = Describe("cloud area test create", func() {

		It("create with normal data", func() {
			rsp, err := hostServerClient.CreateCloudArea(context.Background(), header, NewTmpCloudArea())
			cloudIDTmp = int64(rsp.Data.Created.ID)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("create with cloud area which is already exist", func() {
			tmpTestData := NewTmpCloudArea()
			tmpTestData["bk_cloud_name"] = testData1["bk_cloud_name"]
			rsp, err := hostServerClient.CreateCloudArea(context.Background(), header, tmpTestData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create with cloud vendor which is not valid", func() {
			tmpTestData := NewTmpCloudArea()
			tmpTestData["bk_cloud_name"] = "best mind"
			tmpTestData["bk_cloud_vendor"] = "hello"
			rsp, err := hostServerClient.CreateCloudArea(context.Background(), header, tmpTestData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

	})

	var _ = Describe("cloud area test batch create", func() {

		It("batch create with normal data", func() {
			rsp, err := hostServerClient.CreateManyCloudArea(context.Background(), header, map[string]interface{}{"data": []interface{}{NewTmpCloudArea()}})
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("batch create with cloud area which is already exist", func() {
			tmpTestData := NewTmpCloudArea()
			tmpTestData["bk_cloud_name"] = testData1["bk_cloud_name"]
			rsp, err := hostServerClient.CreateManyCloudArea(context.Background(), header, map[string]interface{}{"data": []interface{}{tmpTestData}})
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("batch create with cloud vendor which is not valid", func() {
			tmpTestData := NewTmpCloudArea()
			tmpTestData["bk_cloud_name"] = "best mind"
			tmpTestData["bk_cloud_vendor"] = "hello"
			rsp, err := hostServerClient.CreateManyCloudArea(context.Background(), header, map[string]interface{}{"data": []interface{}{tmpTestData}})
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

	})

	var _ = Describe("cloud area test update", func() {

		It("update with normal data", func() {
			cloudID := cloudID1
			data := map[string]interface{}{"bk_cloud_name": "LPL200区"}
			rsp, err := hostServerClient.UpdateCloudArea(context.Background(), header, cloudID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update with cloud name which is already exist", func() {
			cloudID := cloudID2
			data := map[string]interface{}{"bk_cloud_name": testData1["bk_cloud_name"]}
			rsp, err := hostServerClient.UpdateCloudArea(context.Background(), header, cloudID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update with cloud vendor which is not valid", func() {
			cloudID := cloudID1
			data := map[string]interface{}{"bk_cloud_vendor": "hello"}
			rsp, err := hostServerClient.UpdateCloudArea(context.Background(), header, cloudID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update with bk_cloud_id which is not exist", func() {
			cloudID := int64(99999)
			data := map[string]interface{}{"bk_cloud_name": "cloudIDNotExist"}
			rsp, err := hostServerClient.UpdateCloudArea(context.Background(), header, cloudID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

	})

	var _ = Describe("cloud area test search", func() {

		It("search with default query condition", func() {
			cond := make(map[string]interface{})
			rsp, err := hostServerClient.SearchCloudArea(context.Background(), header, cond)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(2)))
		})

		It("search with configured condition", func() {
			queryData := map[string]interface{}{"condition": map[string]interface{}{"bk_cloud_name": testData2["bk_cloud_name"]}}
			rsp, err := hostServerClient.SearchCloudArea(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(1)))
			Expect(rsp.Data.Info[0].String("bk_cloud_name")).To(Equal(testData2["bk_cloud_name"]))
		})

		It("search with configured sort", func() {
			queryData := map[string]interface{}{"page": map[string]interface{}{"sort": "create_time"}}
			rsp, err := hostServerClient.SearchCloudArea(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(2)))
			Expect(rsp.Data.Info[0].String("bk_cloud_name")).To(Equal(testData1["bk_cloud_name"]))
			Expect(rsp.Data.Info[1].String("bk_cloud_name")).To(Equal(testData2["bk_cloud_name"]))
		})

		It("search with configured limit", func() {
			queryData := map[string]interface{}{"page": map[string]interface{}{"limit": 1}}
			rsp, err := hostServerClient.SearchCloudArea(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(len(rsp.Data.Info)).To(Equal(int(1)))
		})

		It("search with configured is_fuzzy is false", func() {
			queryData := map[string]interface{}{"is_fuzzy": false, "condition": map[string]interface{}{"bk_cloud_name": "LPL"}}
			rsp, err := hostServerClient.SearchCloudArea(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(0)))

			queryData = map[string]interface{}{"is_fuzzy": false, "condition": map[string]interface{}{"bk_cloud_name": testData2["bk_cloud_name"]}}
			rsp, err = hostServerClient.SearchCloudArea(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(1)))
			Expect(rsp.Data.Info[0].String("bk_cloud_name")).To(Equal(testData2["bk_cloud_name"]))
		})

		It("search with configured is_fuzzy is true", func() {
			queryData := map[string]interface{}{"is_fuzzy": true, "condition": map[string]interface{}{"bk_cloud_name": "LPL"}}
			rsp, err := hostServerClient.SearchCloudArea(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(2)))
		})

	})

	var _ = Describe("cloud area test delete", func() {

		It("delete with normal data", func() {
			cloudID := cloudID1
			rsp, err := hostServerClient.DeleteCloudArea(context.Background(), header, cloudID)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})
	})
})

var cloudID1, cloudID2, cloudIDTmp int64

func prepareCloudData() {
	//清空云区域表
	err := test.GetDB().Table(common.BKTableNameBasePlat).Delete(context.Background(), map[string]interface{}{})
	Expect(err).NotTo(HaveOccurred())

	// 准备数据
	resp, err := hostServerClient.CreateCloudArea(context.Background(), header, testData1)
	cloudID1 = int64(resp.Data.Created.ID)
	util.RegisterResponse(resp)
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.Result).To(Equal(true))

	resp, err = hostServerClient.CreateCloudArea(context.Background(), header, testData2)
	cloudID2 = int64(resp.Data.Created.ID)
	util.RegisterResponse(resp)
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.Result).To(Equal(true))
}
