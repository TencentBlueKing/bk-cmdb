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

package load_test

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common/util"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
	redisdal "configcenter/src/storage/dal/redis"
	"configcenter/src/test"
	"configcenter/src/test/run"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	mongoURI        string
	redisAdress     string
	redisPort       string
	redisPasswd     string
	redisDatabase   string
	redisMasterName string
	redisClient     redis.Client
	tConfig         test.TestConfig
)

var str_1KB string = `
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
`
var DBWriteStr string

func initData() {
	tConfig = test.GetTestConfig()
	redisCfg := redisdal.Config{
		Address:  tConfig.RedisCfg.RedisAdress + ":" + tConfig.RedisCfg.RedisPort,
		Password: tConfig.RedisCfg.RedisPasswd,
		Database: "0",
		//MasterName: redisMasterName,
	}

	var err error
	redisClient, err = redisdal.NewFromConfig(redisCfg)
	if err != nil {
		panic("redisdal.NewFromConfig err:%" + err.Error())
	}

	DBWriteStr = ""
	for i := 0; i < tConfig.DBWriteKBSize; i++ {
		DBWriteStr += str_1KB
	}
}

func GetClient() (*local.Mongo, error) {
	//fmt.Println("********mongouri:", test.GetTestConfig().MongoURI)
	mongoConfig := local.MongoConf{
		MaxOpenConns: mongo.DefaultMaxOpenConns,
		MaxIdleConns: mongo.MinimumMaxIdleOpenConns,
		URI:          tConfig.MongoURI,
		RsName:       "rs0",
	}
	m, err := local.NewMgo(mongoConfig, time.Minute)
	if err != nil {
		return nil, err
	}
	err = m.InitTxnManager(redisClient)
	if err != nil {
		return nil, err
	}
	return m, nil
}

//清空表数据
func clearData(tablename string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}
	err = client.DropTable(context.Background(), tablename)
	if err != nil {
		return err
	}
	err = client.CreateTable(context.Background(), tablename)
	if err != nil {
		return err
	}
	return nil
}

// 分布式事务的Insert
func DistributedInsertCommit(client1, client2 *local.Mongo, tablename string) error {
	sess1, err := client1.StartSession()
	if err != nil {
		return err
	}

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	if err != nil {
		return err
	}
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	if err != nil {
		return err
	}

	ts := fmt.Sprintf("%d", time.Now().UnixNano())

	err = sess1.Table(tablename).Insert(ctx, map[string]string{"distrkey001_" + ts: DBWriteStr})
	if err != nil {
		return err
	}

	//err = client2.Table(tablename).Insert(ctx, map[string]string{"distrkey002_" + ts: DBWriteStr})
	//if err != nil {
	//	return err
	//}

	err = client2.CommitTransaction(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 本地事务Insert
func NotxnInsertCommit(client *local.Mongo, tablename string) error {
	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	return client.Table(tablename).Insert(context.Background(), map[string]string{"notxnkey001" + ts: DBWriteStr})
}

// 无事务的Insert
func LocalInsertCommit(client *local.Mongo, tablename string) error {

	sess, err := client.StartSession()
	if err != nil {
		return err
	}
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	if err != nil {
		return err
	}

	ts := fmt.Sprintf("%d", time.Now().UnixNano())

	err = sess.Table(tablename).Insert(context.Background(), map[string]string{"localkey001" + ts: DBWriteStr})
	if err != nil {
		return err
	}

	err = sess.CommitTransaction(context.Background())
	if err != nil {
		return err
	}
	return nil
}

var _ = Describe("DB Operation Load Test", func() {
	initData()
	client1, err := GetClient()
	//fmt.Println("********getclient err:", err)
	Expect(err).Should(BeNil())
	client2, err := GetClient()
	Expect(err).Should(BeNil())
	tablename := "cc_tranTest"
	//事务操作前，清空数据
	err = clearData(tablename)
	Expect(err).Should(BeNil())

	Describe("distributed insert load test", func() {
		//var header = test.GetHeader()
		Measure("emit the request", func(b Benchmarker) {
			runtime := b.Time("runtime", func() {
				err := DistributedInsertCommit(client1, client2, tablename)
				Expect(err).Should(BeNil())
			})
			Expect(runtime.Seconds()).Should(BeNumerically("<", 0.9))
		}, 10)

		It("running load test", func() {
			m := run.FireLoadTest(func() error {
				err := DistributedInsertCommit(client1, client2, tablename)
				if err != nil {
					return err
				}
				return nil
			})
			fmt.Printf("distributed insert load perform: \n" + m.Format())
		})

	})

	Describe("local insert load test", func() {
		//var header = test.GetHeader()
		Measure("emit the request", func(b Benchmarker) {
			runtime := b.Time("runtime", func() {
				err := LocalInsertCommit(client1, tablename)
				Expect(err).Should(BeNil())
			})
			Expect(runtime.Seconds()).Should(BeNumerically("<", 0.9))
		}, 10)

		It("running load test", func() {
			m := run.FireLoadTest(func() error {
				err := LocalInsertCommit(client1, tablename)
				if err != nil {
					return err
				}
				return nil
			})
			fmt.Printf("local insert load perform: \n" + m.Format())
		})

	})

	Describe("notxn insert load test", func() {
		//var header = test.GetHeader()
		Measure("emit the request", func(b Benchmarker) {
			runtime := b.Time("runtime", func() {
				err := NotxnInsertCommit(client1, tablename)
				Expect(err).Should(BeNil())
			})
			Expect(runtime.Seconds()).Should(BeNumerically("<", 0.9))
		}, 10)

		It("running load test", func() {
			m := run.FireLoadTest(func() error {
				err := NotxnInsertCommit(client1, tablename)
				if err != nil {
					return err
				}
				return nil
			})
			fmt.Printf("local insert load perform: \n" + m.Format())
		})

	})
})
