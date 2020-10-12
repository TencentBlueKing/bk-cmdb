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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/dal/types"
	"configcenter/src/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("host ip array validation test", func() {
	var find types.Find
	It("test preparation", func() {
		tConf := test.GetTestConfig()
		mongoConfig := local.MongoConf{
			MaxOpenConns: mongo.DefaultMaxOpenConns,
			MaxIdleConns: mongo.MinimumMaxIdleOpenConns,
			URI:          tConf.MongoURI,
			RsName:       "rs0",
		}
		db, err := local.NewMgo(mongoConfig, time.Minute)
		Expect(err).To(BeNil())

		redisCfg := redis.Config{
			Address:  tConf.RedisCfg.RedisAdress,
			Password: tConf.RedisCfg.RedisPasswd,
			Database: "0",
		}
		redisClient, err := redis.NewFromConfig(redisCfg)
		Expect(err).To(BeNil())

		err = db.InitTxnManager(redisClient)
		Expect(err).To(BeNil())
		err = db.Table(common.BKTableNameBaseHost).Insert(context.Background(), map[string]interface{}{"bk_host_innerip": []string{"127.0.0.1"}})
		Expect(err).To(BeNil())
		find = db.Table(common.BKTableNameBaseHost).Find(nil).Fields(common.BKHostInnerIPField)
	})

	It("host ip array valid type test", func() {
		By("HostMapStr test", func() {
			err := find.One(context.Background(), &metadata.HostMapStr{})
			Expect(err).To(BeNil())
		})

		By("[]HostMapStr test", func() {
			err := find.All(context.Background(), &[]metadata.HostMapStr{})
			Expect(err).To(BeNil())
		})

		type validStructWithIP struct {
			InnerIP metadata.StringArrayToString `json:"bk_host_innerip" bson:"bk_host_innerip"`
		}

		By("validStructWithIP test", func() {
			err := find.One(context.Background(), &validStructWithIP{})
			Expect(err).To(BeNil())
		})

		By("[]validStructWithIP test", func() {
			err := find.All(context.Background(), &[]validStructWithIP{})
			Expect(err).To(BeNil())
		})

		type structWithoutIP struct {
			ID int64 `json:"bk_host_id" bson:"bk_host_id"`
		}

		By("structWithoutIP test", func() {
			err := find.One(context.Background(), &structWithoutIP{})
			Expect(err).To(BeNil())
		})

		By("[]structWithoutIP test", func() {
			err := find.All(context.Background(), &[]structWithoutIP{})
			Expect(err).To(BeNil())
		})
	})

	It("host ip array invalid type test", func() {
		By("map[string]interface{} test", func() {
			err := find.One(context.Background(), &map[string]interface{}{})
			Expect(err).NotTo(BeNil())
		})

		By("[]map[string]interface{} test", func() {
			err := find.All(context.Background(), &[]map[string]interface{}{})
			Expect(err).NotTo(BeNil())
		})

		type invalidStruct struct {
			InnerIP string `json:"bk_host_innerip" bson:"bk_host_innerip"`
		}

		By("invalidStruct test", func() {
			err := find.One(context.Background(), &invalidStruct{})
			Expect(err).NotTo(BeNil())
		})

		By("[]invalidStruct test", func() {
			err := find.All(context.Background(), &[]invalidStruct{})
			Expect(err).NotTo(BeNil())
		})
	})
})
