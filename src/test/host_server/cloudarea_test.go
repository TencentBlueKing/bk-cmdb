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
	"bk_vpc_id":       1,
	"bk_vpc_name":     "Default-VPC",
	"bk_account_id":   2,
	"creator":         "admin",
}

// wrong testData
var tmpTestData = map[string]interface{}{
	"bk_cloud_name":   "LPL39区",
	"bk_status":       "3",
	"bk_cloud_vendor": "2",
	"bk_region":       "1111",
	"bk_vpc_id":       1,
	"bk_vpc_name":     "Default-VPC",
	"bk_account_id":   2,
	"creator":         "admin",
}

var _ = Describe("cloud area test", func() {

	var _ = Describe("create cloud area test", func() {

		It("create with normal data1", func() {
			rsp, err := hostServerClient.CreateCloudArea(context.Background(), header, testData1)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("create with normal data2", func() {
			rsp, err := hostServerClient.CreateCloudArea(context.Background(), header, testData2)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("create with cloud area which is already exist", func() {
			tmpTestData["bk_cloud_name"] = testData1["bk_cloud_name"]
			rsp, err := hostServerClient.CreateCloudArea(context.Background(), header, tmpTestData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("create with cloud vendor which is not valid", func() {
			testData1["bk_cloud_vendor"] = "hello"
			rsp, err := hostServerClient.CreateCloudArea(context.Background(), header, testData1)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

	})

	var _ = Describe("update cloud area test", func() {

		It("update with normal data", func() {
			cloudID := int64(2)
			data := map[string]interface{}{"bk_cloud_name": "LPL200区"}
			rsp, err := hostServerClient.UpdateCloudArea(context.Background(), header, cloudID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})

		It("update with cloud name which is already exist", func() {
			cloudID := int64(3)
			data := map[string]interface{}{"bk_cloud_name": "LPL200区"}
			rsp, err := hostServerClient.UpdateCloudArea(context.Background(), header, cloudID, data)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(false))
		})

		It("update with cloud vendor which is not valid", func() {
			cloudID := int64(2)
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

	var _ = Describe("search cloud area test", func() {

		It("search with default query condition", func() {
			cond := make(map[string]interface{})
			rsp, err := hostServerClient.SearchCloudArea(context.Background(), header, cond)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(3)))
		})

		It("search with configured conditon", func() {
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
			Expect(rsp.Data.Count).To(Equal(int64(3)))
			Expect(rsp.Data.Info[1].String("bk_cloud_name")).To(Equal("LPL200区"))
			Expect(rsp.Data.Info[2].String("bk_cloud_name")).To(Equal(testData2["bk_cloud_name"]))
		})

		It("search with configured limit", func() {
			queryData := map[string]interface{}{"page": map[string]interface{}{"limit": 1}}
			rsp, err := hostServerClient.SearchCloudArea(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(len(rsp.Data.Info)).To(Equal(int(1)))
		})

		It("search with configured exact is true", func() {
			queryData := map[string]interface{}{"exact": true, "condition": map[string]interface{}{"bk_cloud_name": "LPL"}}
			rsp, err := hostServerClient.SearchCloudArea(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(0)))

			queryData = map[string]interface{}{"exact": true, "condition": map[string]interface{}{"bk_cloud_name": testData2["bk_cloud_name"]}}
			rsp, err = hostServerClient.SearchCloudArea(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(1)))
			Expect(rsp.Data.Info[0].String("bk_cloud_name")).To(Equal(testData2["bk_cloud_name"]))
		})

		It("search with configured exact is false", func() {
			queryData := map[string]interface{}{"exact": false, "condition": map[string]interface{}{"bk_cloud_name": "LPL"}}
			rsp, err := hostServerClient.SearchCloudArea(context.Background(), header, queryData)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
			Expect(rsp.Data.Count).To(Equal(int64(2)))
		})

	})

	var _ = Describe("delete cloud area test", func() {

		It("delete with normal data", func() {
			cloudID := int64(2)
			rsp, err := hostServerClient.DeleteCloudArea(context.Background(), header, cloudID)
			util.RegisterResponse(rsp)
			Expect(err).NotTo(HaveOccurred())
			Expect(rsp.Result).To(Equal(true))
		})
	})
})
