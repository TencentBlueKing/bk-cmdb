/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package local

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"testing"
	"time"

	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
	redisdal "configcenter/src/storage/dal/redis"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"gopkg.in/redis.v5"
)

var (
	mongoURI        string
	redisAdress     string
	redisPort       string
	redisPasswd     string
	redisDatabase   string
	redisMasterName string
	redisClient     *redis.Client
)

func init() {
	flag.StringVar(&mongoURI, "mongo-addr", "mongodb://cc:cc@localhost:27011,localhost:27012,localhost:27013,localhost:27014/cmdb", "mongodb URI")
	flag.StringVar(&redisAdress, "redisAdress", "127.0.0.1", "redis host address")
	flag.StringVar(&redisPort, "redisPort", "6379", "redis host port")
	flag.StringVar(&redisPasswd, "redisPasswd", "cc", "redis password")
	flag.StringVar(&redisDatabase, "redisDatabase", "0", "redis database")
	flag.StringVar(&redisMasterName, "redisMasterName", "", "redis master name")
	flag.Parse()

	redisCfg := redisdal.Config{
		Address:    redisAdress,
		Port:       redisPort,
		Password:   redisPasswd,
		Database:   redisDatabase,
		MasterName: redisMasterName,
	}

	var err error
	redisClient, err = redisdal.NewFromConfig(redisCfg)
	if err != nil {
		panic("redisdal.NewFromConfig err:%" + err.Error())
	}
}

func GetClient() (*Mongo, error) {
	m, err := NewMgo(mongoURI, time.Minute)
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
func clearData(t *testing.T, tablename string) {
	client, err := GetClient()
	require.NoError(t, err)
	err = client.DropTable(context.Background(), tablename)
	require.NoError(t, err)
	err = client.CreateTable(context.Background(), tablename)
	require.NoError(t, err)
}

// 测试本地事务Insert的commit
func TestLocalInsertCommit(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table(tablename).Insert(context.Background(), map[string]string{"key001": "value001"})
	require.NoError(t, err)

	err = sess.CommitTransaction(context.Background())
	require.NoError(t, err)

	//校验结果
	cnt, err := client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)
}

// 测试本地事务Insert的abort
func TestLocalInsertAbort(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table(tablename).Insert(context.Background(), map[string]string{"key001": "value001"})
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.NoError(t, err)

	//校验结果
	cnt, err := client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)
}

// 测试本地事务Delete的commit
func TestLocalDeleteCommit(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	// 准备一些数据
	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"k1": "v1",
		},
		map[string]interface{}{
			"k2": "v2",
		},
	}
	err = client.Table(tablename).Insert(context.Background(), insertDataMany)
	require.NoError(t, err)

	cnt, err := client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table(tablename).Delete(context.Background(), map[string]string{"k1": "v1"})
	require.NoError(t, err)

	err = sess.CommitTransaction(context.Background())
	require.NoError(t, err)

	//校验结果
	cnt, err = client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)
}

// 测试本地事务Delete的abort
func TestLocalDeleteAbort(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	// 准备一些数据
	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"k1": "v1",
		},
		map[string]interface{}{
			"k2": "v2",
		},
	}
	err = client.Table(tablename).Insert(context.Background(), insertDataMany)
	require.NoError(t, err)

	cnt, err := client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table(tablename).Delete(context.Background(), map[string]string{"k1": "v1"})
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.NoError(t, err)

	//校验结果
	cnt, err = client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)
}

// 测试本地事务Update的commit
func TestLocalUpdateCommit(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	// 准备一些数据
	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"k1": "v1",
		},
		map[string]interface{}{
			"k2": "v2",
		},
	}
	err = client.Table(tablename).Insert(context.Background(), insertDataMany)
	require.NoError(t, err)

	cnt, err := client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table(tablename).Update(context.Background(), map[string]interface{}{"k1": "v1"}, map[string]interface{}{"k1": "update1"})
	require.NoError(t, err)

	err = sess.CommitTransaction(context.Background())
	require.NoError(t, err)

	//校验结果
	cnt, err = client.Table(tablename).Find(map[string]interface{}{"k1": "v1"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)

	cnt, err = client.Table(tablename).Find(map[string]interface{}{"k1": "update1"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)
}

// 测试本地事务Update的commit
func TestLocalUpdateAbort(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	// 准备一些数据
	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"k1": "v1",
		},
		map[string]interface{}{
			"k2": "v2",
		},
	}
	err = client.Table(tablename).Insert(context.Background(), insertDataMany)
	require.NoError(t, err)

	cnt, err := client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table(tablename).Update(context.Background(), map[string]interface{}{"k1": "v1"}, map[string]interface{}{"k1": "update1"})
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.NoError(t, err)

	//校验结果
	cnt, err = client.Table(tablename).Find(map[string]interface{}{"k1": "v1"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	cnt, err = client.Table(tablename).Find(map[string]interface{}{"k1": "update1"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)
}

// 测试本地事务的隔离性
func TestLocalIsolation(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	// 开启事务
	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)
	// 激活服务端的txnNumber，从这时起，该事务外的操作对该事务不可见
	_, err = sess.Table("-_-aaa-_-").Find(nil).Count(context.Background())
	require.NoError(t, err)

	// 事务外插入一条数据，自动进行了提交
	err = client.Table(tablename).Insert(context.Background(), map[string]string{"key002": "value002"})
	require.NoError(t, err)

	// 事务外能看到事务外插入的数据
	cnt, err := client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	// 正在操作中的事务看不到事务外插入的数据
	cnt, err = sess.Table(tablename).Find(nil).Count(context.Background())
	//resultMany := make([]map[string]interface{}, 0)
	//err = sess.Table(tablename).Find(nil).All(context.Background(), &resultMany)
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)
	//require.Equal(t, 0, len(resultMany))

	// 事务内插入一条数据
	err = sess.Table(tablename).Insert(context.Background(), map[string]string{"key001": "value001"})
	require.NoError(t, err)

	// 事务内可以看到自身事务插入的数据
	cnt, err = sess.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	//事务外不能看到未提交事务插入的数据
	cnt, err = client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	// 提交事务
	err = sess.CommitTransaction(context.Background())
	require.NoError(t, err)

	// 事务完成提交后，事务外可以看到插入的数据
	cnt, err = client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)
}

// 测试本地事务的重复开启事务
func TestLocalStartStart(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.StartTransaction(context.Background())
	require.EqualError(t, err, "transaction already in progress")
}

// 测试本地事务提交后又取消
func TestLocalCommitAbort(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.CommitTransaction(context.Background())
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.EqualError(t, err, "cannot call abortTransaction after calling commitTransaction")
}

// 测试本地事务取消后又提交
func TestLocalAbortCommit(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.NoError(t, err)

	err = sess.CommitTransaction(context.Background())
	require.EqualError(t, err, "cannot call commitTransaction after calling abortTransaction")
}

// 测试本地事务的重复取消
func TestLocalAbortAbort(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.EqualError(t, err, "cannot call abortTransaction twice")
}

// 测试本地事务的重复提交,不会报错，mongoDB不把重复提交作为错误
func TestLocalCommitCommit(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.CommitTransaction(context.Background())
	require.NoError(t, err)

	err = sess.CommitTransaction(context.Background())
	require.NoError(t, err)
}

// 测试本地事务提交后,再继续进行DB操作，最后再提交
func TestLocalCommitOpCommit(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.CommitTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table("cc_tranTest").Insert(context.Background(), map[string]string{"key001": "value001"})
	require.NoError(t, err)

	err = sess.CommitTransaction(context.Background())
	require.EqualError(t, err, "no transaction started")
}

// 测试本地事务取消后,再继续进行操作，最后再取消
func TestLocalAbortOpAbort(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table("cc_tranTest").Insert(context.Background(), map[string]string{"key001": "value001"})
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.EqualError(t, err, "no transaction started")
}

// 测试本地事务提交后,再继续进行操作，最后再取消
func TestLocalCommitOpAbort(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.CommitTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table("cc_tranTest").Insert(context.Background(), map[string]string{"key001": "value001"})
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.EqualError(t, err, "no transaction started")
}

// 测试本地事务取消后,再继续进行操作，最后再提交
func TestLocalAbortOpCommit(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table("cc_tranTest").Insert(context.Background(), map[string]string{"key001": "value001"})
	require.NoError(t, err)

	err = sess.CommitTransaction(context.Background())
	require.EqualError(t, err, "no transaction started")
}

// 测试本地事务的嵌套, 需要用到两个不同的session，内部的事务不受外部事务的影响
func TestLocalEmbed(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	sess2, err := client.StartSession()
	require.NoError(t, err)
	defer sess2.EndSession(context.Background())

	err = sess2.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess2.Table(tablename).Insert(context.Background(), map[string]string{"key": "inner"})
	require.NoError(t, err)

	err = sess2.CommitTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table(tablename).Insert(context.Background(), map[string]string{"key": "outter"})
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.NoError(t, err)

	//校验结果
	cnt, err := client.Table(tablename).Find(map[string]string{"key": "inner"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	cnt, err = client.Table(tablename).Find(map[string]string{"key": "outter"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)
}

// 测试本地事务,使用已经结束的session
func TestLocalEndSession(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess, err := client.StartSession()
	require.NoError(t, err)
	sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table(tablename).Insert(context.Background(), map[string]string{"key": "value"})
	require.EqualError(t, err, "ended session was used")
}

// 测试本地事务的客户端超时情况
// 可以看出客户端超时了，但服务端会继续执行操作，客户端一段时间后进行提交，数据落地成功
// 为防止数据出问题，要在客户端超时报错进行错误判断，有错误就进行回滚
func PTestLocalClientTimeout(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	// 设置客户端超时时间为微妙级别，不超过1毫秒
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Microsecond)
	defer cancel()
	err = sess.Table(tablename).Insert(ctx, map[string]string{"key001": "value001"})
	require.Error(t, err)

	time.Sleep(200 * time.Millisecond)
	err = sess.CommitTransaction(context.Background())
	require.NoError(t, err)

	//校验结果
	cnt, err := client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)
}

// 测试本地事务的服务端超时情况
func PTestLocalServerTimeout(t *testing.T) {
	client, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess, err := client.StartSession()
	require.NoError(t, err)
	defer sess.EndSession(context.Background())

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table("cc_tranTest").Insert(context.Background(), map[string]string{"key001": "value001"})
	require.NoError(t, err)

	// 休眠足够长时间，导致事务超时
	time.Sleep(100 * time.Second)

	err = sess.CommitTransaction(context.Background())
	require.EqualError(t, err, "(NoSuchTransaction) Transaction 1 has been aborted.")

	err = sess.StartTransaction(context.Background())
	require.NoError(t, err)

	err = sess.Table("cc_tranTest").Insert(context.Background(), map[string]string{"key002": "value002"})
	require.NoError(t, err)

	err = sess.AbortTransaction(context.Background())
	require.NoError(t, err)

	//校验结果
	cnt, err := client.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)
}

// 测试分布式事务的重复开启
func TestDistributedStartStart(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = client2.StartTransaction(ctx)
	require.EqualError(t, err, "transaction already in progress")
}

// 测试分布式事务提交后又取消
func TestDistributedCommitAbort(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = sess1.CommitTransaction(ctx)
	require.NoError(t, err)

	err = client2.AbortTransaction(ctx)
	require.EqualError(t, err, "cannot call abortTransaction after calling commitTransaction")
}

// 测试分布式事务取消后又提交
func TestDistributedAbortCommit(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = sess1.AbortTransaction(ctx)
	require.NoError(t, err)

	err = client2.CommitTransaction(ctx)
	require.EqualError(t, err, "cannot call commitTransaction after calling abortTransaction")
}

// 测试分布式事务的重复取消
func TestDistributedAbortAbort(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = sess1.AbortTransaction(ctx)
	require.NoError(t, err)

	err = client2.AbortTransaction(ctx)
	require.EqualError(t, err, "cannot call abortTransaction twice")
}

// 测试分布式事务的重复提交,不会报错，mongoDB不把重复提交作为错误
func TestDistributedCommitCommit(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	// local测试不用DB操作，是因为通过判断单个session对象的s.didCommitAfterStart = true而不给服务端发送commit请求
	// distributed测试因为没有共享dsession对象的idCommitAfterStart，所以会给服务端发送commit请求，发送前需要进行DB操作来激活txnNumber
	err = sess1.Table("cc_tranTest").Insert(ctx, map[string]string{"key001": "value001"})
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = sess1.CommitTransaction(ctx)
	require.NoError(t, err)

	err = client2.CommitTransaction(ctx)
	require.NoError(t, err)
}

// 测试分布式事务的提交后，继续进行DB操作，然后再次提交
func TestDistributedCommitOpCommit(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table("cc_tranTest").Insert(ctx, map[string]string{"key001": "value001"})
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = sess1.CommitTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table("cc_tranTest").Insert(ctx, map[string]string{"key001": "value001"})
	require.NoError(t, err)

	err = client2.CommitTransaction(ctx)
	require.EqualError(t, err, "no transaction started")
}

// 测试分布式事务的Insert提交
func TestDistributedInsertCommit(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table(tablename).Insert(ctx, map[string]string{"key001": "value001"})
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = client2.Table(tablename).Insert(ctx, map[string]string{"key002": "value002"})
	require.NoError(t, err)

	err = client2.CommitTransaction(ctx)
	require.NoError(t, err)

	//校验结果
	cnt, err := client2.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)
}

// 测试分布式事务的Insert取消
func TestDistributedInsertAbort(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table(tablename).Insert(ctx, map[string]string{"key001": "value001"})
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = client2.Table(tablename).Insert(ctx, map[string]string{"key002": "value002"})
	require.NoError(t, err)

	err = client2.AbortTransaction(ctx)
	require.NoError(t, err)

	//校验结果
	cnt, err := client2.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)
}

// 测试分布式事务的Delete提交
func TestDistributedDeleteCommit(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	// 准备一些数据
	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"k1": "v1",
		},
		map[string]interface{}{
			"k2": "v2",
		},
	}
	err = client1.Table(tablename).Insert(context.Background(), insertDataMany)
	require.NoError(t, err)

	cnt, err := client1.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table(tablename).Delete(ctx, map[string]string{"k1": "v1"})
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = client2.Table(tablename).Delete(ctx, map[string]string{"k2": "v2"})
	require.NoError(t, err)

	err = client2.CommitTransaction(ctx)
	require.NoError(t, err)

	//校验结果
	cnt, err = client1.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)
}

// 测试分布式事务的Delete取消
func TestDistributedDeleteAbort(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	// 准备一些数据
	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"k1": "v1",
		},
		map[string]interface{}{
			"k2": "v2",
		},
	}
	err = client1.Table(tablename).Insert(context.Background(), insertDataMany)
	require.NoError(t, err)

	cnt, err := client1.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table(tablename).Delete(ctx, map[string]string{"k1": "v1"})
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = client2.Table(tablename).Delete(ctx, map[string]string{"k2": "v2"})
	require.NoError(t, err)

	err = client2.AbortTransaction(ctx)
	require.NoError(t, err)

	//校验结果
	cnt, err = client1.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)
}

// 测试分布式事务的Update提交
func TestDistributedUpdateCommit(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	// 准备一些数据
	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"k1": "v1",
		},
		map[string]interface{}{
			"k2": "v2",
		},
	}
	err = client1.Table(tablename).Insert(context.Background(), insertDataMany)
	require.NoError(t, err)

	cnt, err := client1.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table(tablename).Update(ctx, map[string]string{"k1": "v1"}, map[string]string{"k1": "update1"})
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = client2.Table(tablename).Update(ctx, map[string]string{"k2": "v2"}, map[string]string{"k2": "update2"})
	require.NoError(t, err)

	err = client2.CommitTransaction(ctx)
	require.NoError(t, err)

	//校验结果
	cnt, err = client1.Table(tablename).Find(map[string]string{"k1": "v1"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)

	cnt, err = client1.Table(tablename).Find(map[string]string{"k2": "v2"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)

	cnt, err = client1.Table(tablename).Find(map[string]string{"k1": "update1"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	cnt, err = client1.Table(tablename).Find(map[string]string{"k2": "update2"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)
}

// 测试分布式事务的Update取消
func TestDistributedUpdateAbort(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	// 准备一些数据
	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"k1": "v1",
		},
		map[string]interface{}{
			"k2": "v2",
		},
	}
	err = client1.Table(tablename).Insert(context.Background(), insertDataMany)
	require.NoError(t, err)

	cnt, err := client1.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table(tablename).Update(ctx, map[string]string{"k1": "v1"}, map[string]string{"k1": "update1"})
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = client2.Table(tablename).Update(ctx, map[string]string{"k2": "v2"}, map[string]string{"k2": "update2"})
	require.NoError(t, err)

	err = client2.AbortTransaction(ctx)
	require.NoError(t, err)

	//校验结果
	cnt, err = client1.Table(tablename).Find(map[string]string{"k1": "v1"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	cnt, err = client1.Table(tablename).Find(map[string]string{"k2": "v2"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	cnt, err = client1.Table(tablename).Find(map[string]string{"k1": "update1"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)

	cnt, err = client1.Table(tablename).Find(map[string]string{"k2": "update2"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)
}

// 测试分布式事务的Upsert提交
func TestDistributedUpsertCommit(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	// 准备一些数据
	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"k1": "v1",
		},
		map[string]interface{}{
			"k2": "v2",
		},
	}
	err = client1.Table(tablename).Insert(context.Background(), insertDataMany)
	require.NoError(t, err)

	cnt, err := client1.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table(tablename).Upsert(ctx, map[string]string{"k1": "v1"}, map[string]string{"k1": "update1"})
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = client2.Table(tablename).Upsert(ctx, map[string]string{"k2": "v5"}, map[string]string{"k9": "update2"})
	require.NoError(t, err)

	err = client2.CommitTransaction(ctx)
	require.NoError(t, err)

	//校验结果
	cnt, err = client1.Table(tablename).Find(map[string]string{"k1": "v1"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)

	cnt, err = client1.Table(tablename).Find(map[string]string{"k2": "v2"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	cnt, err = client1.Table(tablename).Find(map[string]string{"k1": "update1"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	resultOne := make(map[string]string, 0)
	err = client1.Table(tablename).Find(map[string]string{"k2": "v5"}).One(context.Background(), &resultOne)
	require.NoError(t, err)
	require.Equal(t, map[string]string{"k2": "v5", "k9": "update2"}, resultOne)
}

// 测试分布式事务的Upsert取消
func TestDistributedUpsertAbort(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	// 准备一些数据
	insertDataMany := []map[string]interface{}{
		map[string]interface{}{
			"k1": "v1",
		},
		map[string]interface{}{
			"k2": "v2",
		},
	}
	err = client1.Table(tablename).Insert(context.Background(), insertDataMany)
	require.NoError(t, err)

	cnt, err := client1.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table(tablename).Upsert(ctx, map[string]string{"k1": "v1"}, map[string]string{"k1": "update1"})
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = client2.Table(tablename).Upsert(ctx, map[string]string{"k2": "v5"}, map[string]string{"k9": "update2"})
	require.NoError(t, err)

	err = client2.AbortTransaction(ctx)
	require.NoError(t, err)

	//校验结果
	cnt, err = client1.Table(tablename).Find(map[string]string{"k1": "v1"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	cnt, err = client1.Table(tablename).Find(map[string]string{"k2": "v2"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	cnt, err = client1.Table(tablename).Find(map[string]string{"k1": "update1"}).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)

	resultOne := make(map[string]string, 0)
	err = client1.Table(tablename).Find(map[string]string{"k2": "v5"}).One(context.Background(), &resultOne)
	require.EqualError(t, err, "document not found")
}

// 测试分布式事务的UpdateMultiModel提交
func TestDistributedUpdateMultiModelCommit(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	// 准备一些数据
	type RowStruct struct {
		A     string  `bson:"a"`
		B     string  `bson:"b"`
		Ext   string  `bson:"ext"`
		Sort  string  `bson:"sort"`
		Inc   int64   `bson:"inc"`
		Unset *string `bson:"unset"`
	}
	unsetVal := "test_val"
	insertData := RowStruct{
		A:     "a",
		B:     "b",
		Ext:   "ext",
		Sort:  "2",
		Inc:   1,
		Unset: &unsetVal,
	}
	err = client1.Table(tablename).Insert(context.Background(), insertData)
	require.NoError(t, err)

	resultData := RowStruct{
		A:     "a_update_multi_model",
		B:     "b_update_multi_model",
		Ext:   "ext",
		Sort:  "2",
		Inc:   5,
		Unset: nil,
	}

	filter := map[string]string{"ext": "ext"}
	update1 := []dal.ModeUpdate{
		dal.ModeUpdate{Op: "set", Doc: map[string]string{"a": "a_update_multi_model"}},
		dal.ModeUpdate{Op: "unset", Doc: map[string]string{"unset": ""}},
		dal.ModeUpdate{Op: "inc", Doc: map[string]interface{}{"inc": 1}},
	}
	update2 := []dal.ModeUpdate{
		dal.ModeUpdate{Op: "set", Doc: map[string]string{"b": "b_update_multi_model"}},
		dal.ModeUpdate{Op: "inc", Doc: map[string]interface{}{"inc": 3}},
	}

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table(tablename).UpdateMultiModel(ctx, filter, update1...)
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = client2.Table(tablename).UpdateMultiModel(ctx, filter, update2...)
	require.NoError(t, err)

	err = client2.CommitTransaction(ctx)
	require.NoError(t, err)

	//校验结果
	resultOne := RowStruct{}
	err = client1.Table(tablename).Find(nil).One(context.Background(), &resultOne)
	require.NoError(t, err)
	require.Equal(t, resultData, resultOne)
}

// 测试分布式事务的UpdateMultiModel取消
func TestDistributedUpdateMultiModelAbort(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	// 准备一些数据
	type RowStruct struct {
		A     string  `bson:"a"`
		B     string  `bson:"b"`
		Ext   string  `bson:"ext"`
		Sort  string  `bson:"sort"`
		Inc   int64   `bson:"inc"`
		Unset *string `bson:"unset"`
	}
	unsetVal := "test_val"
	insertData := RowStruct{
		A:     "a",
		B:     "b",
		Ext:   "ext",
		Sort:  "2",
		Inc:   1,
		Unset: &unsetVal,
	}
	err = client1.Table(tablename).Insert(context.Background(), insertData)
	require.NoError(t, err)

	filter := map[string]string{"ext": "ext"}
	update1 := []dal.ModeUpdate{
		dal.ModeUpdate{Op: "set", Doc: map[string]string{"a": "a_update_multi_model"}},
		dal.ModeUpdate{Op: "unset", Doc: map[string]string{"unset": ""}},
		dal.ModeUpdate{Op: "inc", Doc: map[string]interface{}{"inc": 1}},
	}
	update2 := []dal.ModeUpdate{
		dal.ModeUpdate{Op: "set", Doc: map[string]string{"b": "b_update_multi_model"}},
		dal.ModeUpdate{Op: "inc", Doc: map[string]interface{}{"inc": 3}},
	}

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table(tablename).UpdateMultiModel(ctx, filter, update1...)
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	err = client2.Table(tablename).UpdateMultiModel(ctx, filter, update2...)
	require.NoError(t, err)

	err = client2.AbortTransaction(ctx)
	require.NoError(t, err)

	//校验结果
	resultOne := RowStruct{}
	err = client1.Table(tablename).Find(nil).One(context.Background(), &resultOne)
	require.NoError(t, err)
	require.Equal(t, insertData, resultOne)
}

// 测试分布式事务的Aggregate提交
func TestDistributedAggregateCommit(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table(tablename).Insert(ctx, map[string]string{"aa": "aa"})
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	aggregateCond := []interface{}{
		map[string]interface{}{
			"$group": map[string]interface{}{
				"_id": "$aa",
				"num": map[string]interface{}{"$sum": 1},
			},
		},
	}

	type aggregateRowStruct struct {
		ID  string `bson:"_id"`
		Num int64  `bson:"num"`
	}
	resultOne := &aggregateRowStruct{}
	err = sess1.Table(tablename).AggregateOne(ctx, aggregateCond, resultOne)
	require.NoError(t, err)
	require.Equal(t, aggregateRowStruct{
		ID:  "aa",
		Num: 1,
	}, *resultOne)

	resultAll := make([]aggregateRowStruct, 0)
	err = client2.Table(tablename).AggregateAll(ctx, aggregateCond, &resultAll)
	require.NoError(t, err)
	if len(resultAll) == 0 {
		t.Errorf("AggregateOne error")
		return
	}
	require.Equal(t, aggregateRowStruct{
		ID:  "aa",
		Num: 1,
	}, resultAll[0])

	err = client2.CommitTransaction(ctx)
	require.NoError(t, err)

	//校验结果
	err = client1.Table(tablename).AggregateOne(context.Background(), aggregateCond, resultOne)
	require.NoError(t, err)
	require.Equal(t, aggregateRowStruct{
		ID:  "aa",
		Num: 1,
	}, *resultOne)
}

// 测试分布式事务的Aggregate取消
func TestDistributedAggregateAbort(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	err = sess1.Table(tablename).Insert(ctx, map[string]string{"aa": "aa"})
	require.NoError(t, err)

	client2, err := GetClient()
	require.NoError(t, err)

	aggregateCond := []interface{}{
		map[string]interface{}{
			"$group": map[string]interface{}{
				"_id": "$aa",
				"num": map[string]interface{}{"$sum": 1},
			},
		},
	}

	type aggregateRowStruct struct {
		ID  string `bson:"_id"`
		Num int64  `bson:"num"`
	}
	resultOne := &aggregateRowStruct{}
	err = sess1.Table(tablename).AggregateOne(ctx, aggregateCond, resultOne)
	require.NoError(t, err)
	require.Equal(t, aggregateRowStruct{
		ID:  "aa",
		Num: 1,
	}, *resultOne)

	resultAll := make([]aggregateRowStruct, 0)
	err = client2.Table(tablename).AggregateAll(ctx, aggregateCond, &resultAll)
	require.NoError(t, err)
	if len(resultAll) == 0 {
		t.Errorf("AggregateOne error")
		return
	}
	require.Equal(t, aggregateRowStruct{
		ID:  "aa",
		Num: 1,
	}, resultAll[0])

	err = client2.AbortTransaction(ctx)
	require.NoError(t, err)

	//校验结果
	err = client1.Table(tablename).AggregateOne(context.Background(), aggregateCond, resultOne)
	require.EqualError(t, err, "document not found")
}

// 测试分布式事务的隔离性
// context里包含事务信息则是事务内操作，否则为事务外操作
func TestDistributedIsolation(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	defer sess1.EndSession(ctx)

	// 开启事务
	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)

	// 激活服务端的txnNumber，从这时起，该事务外的操作对该事务不可见
	_, err = sess1.Table("-_-aaa-_-").Find(nil).Count(ctx)
	require.NoError(t, err)

	// 事务外插入一条数据，自动进行了提交
	err = client1.Table(tablename).Insert(context.Background(), map[string]string{"key001": "value001"})
	require.NoError(t, err)

	// 事务外能看到事务外插入的数据
	cnt, err := client1.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	// 正在操作中的事务看不到事务外插入的数据
	cnt, err = sess1.Table(tablename).Find(nil).Count(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(0), cnt)

	// 事务内用client1插入一条数据
	err = sess1.Table(tablename).Insert(ctx, map[string]string{"key002": "value002"})
	require.NoError(t, err)

	// 事务内用client2插入一条数据
	client2, err := GetClient()
	require.NoError(t, err)
	err = client2.Table(tablename).Insert(ctx, map[string]string{"key003": "value003"})
	require.NoError(t, err)

	// 事务内可以看到自身事务插入的数据
	cnt, err = client2.Table(tablename).Find(nil).Count(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(2), cnt)

	//事务外不能看到未提交事务插入的数据
	cnt, err = client1.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(1), cnt)

	// 提交事务
	err = sess1.CommitTransaction(ctx)
	require.NoError(t, err)

	// 事务完成提交后，事务外可以看到插入的数据
	cnt, err = client1.Table(tablename).Find(nil).Count(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(3), cnt)
}

// 测试分布式事务, 使用已经结束的session
// 只要ctx中含有事务信息，每次操作都会用生成新session操作
func TestDistributedInsertEndSession(t *testing.T) {
	client1, err := GetClient()
	require.NoError(t, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	clearData(t, tablename)

	sess1, err := client1.StartSession()
	require.NoError(t, err)

	// 获取事务信息，将其存入context中
	tnxInfo, err := sess1.TxnInfo()
	require.NoError(t, err)
	ctx := util.TnxIntoContext(context.Background(), tnxInfo)

	sess1.EndSession(ctx)

	err = sess1.StartTransaction(ctx)
	require.NoError(t, err)
}

func BenchmarkDistributedCUD(b *testing.B) {
	client1, err := GetClient()
	require.NoError(b, err)
	tablename := "cc_tranTest"

	//事务操作前，清空数据
	require.NoError(b, err)
	err = client1.DropTable(context.Background(), tablename)
	require.NoError(b, err)
	err = client1.CreateTable(context.Background(), tablename)
	require.NoError(b, err)

	client2, err := GetClient()
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		sess1, err := client1.StartSession()
		require.NoError(b, err)

		// 获取事务信息，将其存入context中
		tnxInfo, err := sess1.TxnInfo()
		require.NoError(b, err)
		ctx := util.TnxIntoContext(context.Background(), tnxInfo)

		defer sess1.EndSession(ctx)

		err = sess1.StartTransaction(ctx)
		require.NoError(b, err)

		err = sess1.Table(tablename).Insert(ctx, map[string]interface{}{"name": "a"})
		require.NoError(b, err)

		err = client1.Table(tablename).Update(ctx, map[string]interface{}{"name": "a"}, map[string]interface{}{"name": "b"})
		require.NoError(b, err)

		err = client2.Table(tablename).Delete(ctx, map[string]interface{}{"name": "b"})
		require.NoError(b, err)

		err = client2.CommitTransaction(ctx)
		require.NoError(b, err)
	}
}

// 测试验证sessionID生成的唯一性，保证并发操作时，事务上下文信息不会相同冲突
func aTestSessionID(t *testing.T) {
	m1 := map[string]bool{}
	m2 := map[string]bool{}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			id, _ := uuid.New()
			m1[fmt.Sprintf("%v", [16]byte(id))] = true
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			id, _ := uuid.New()
			m2[fmt.Sprintf("%v", [16]byte(id))] = true
		}
	}()
	wg.Wait()
	i := 0
	for k := range m1 {
		i++
		if _, ok := m2[k]; ok {
			t.Errorf("has same sessionID:%s\n", k)
		}
	}
	fmt.Println("compare count:", i)
}
