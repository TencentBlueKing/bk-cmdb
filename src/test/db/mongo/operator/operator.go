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

package operator

import (
	"context"
	"fmt"
	"strconv"
	//"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/test"
)

var (
	str_1KB = `
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

	str_1Byte = "a"
)

// MongoOperator mongo operator which contains mongo test cases
type MongoOperator struct {
	tableName string
	data      string
	db        *local.Mongo
	tConfig   test.TestConfig
}

// NewMongoOperator new a mongo operator instance
func NewMongoOperator(tableName string) *MongoOperator {
	tConfig := test.GetTestConfig()
	data := ""
	for i := 0; i < tConfig.DBWriteByteSize; i++ {
		data += str_1Byte
	}
	return &MongoOperator{
		tableName: tableName,
		data:      data,
		db:        test.GetDB(),
		tConfig:   tConfig,
	}
}

//清空表数据
func (m *MongoOperator) ClearData() error {
	err := m.db.DropTable(context.Background(), m.tableName)
	if err != nil {
		return err
	}

	err = m.db.CreateTable(context.Background(), m.tableName)
	if err != nil {
		return err
	}

	// 准备一条数据，用于读
	m.db.Table(m.tableName).Insert(context.Background(), map[string]string{"testread": m.data})

	return nil
}

// GenerateTxnKit 生成事务需要的组件
func (m *MongoOperator) GenerateTxnKit() (context.Context, *metadata.TxnCapable, error) {
	sessionID, err := local.GenSessionID()
	if err != nil {
		return nil, nil, fmt.Errorf("generate session id failed, err: %v", err)
	}

	timeout := common.TransactionDefaultTimeout

	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionIdHeader, sessionID)
	ctx = context.WithValue(ctx, common.TransactionTimeoutHeader, strconv.FormatInt(int64(timeout), 10))

	cap := &metadata.TxnCapable{
		Timeout:   timeout,
		SessionID: sessionID,
	}

	return ctx, cap, nil
}

// 有分布式事务的Write
func (m *MongoOperator) WriteWithTxn() error {
	ctx, cap, err := m.GenerateTxnKit()
	if err != nil {
		blog.Info("genTxnKit err: %s", err)
		return err
	}

	txnErr := m.db.Table(m.tableName).Insert(ctx, map[string]string{"txnfield": m.data})
	if txnErr == nil {
		err := m.db.CommitTransaction(ctx, cap)
		if err != nil {
			blog.Info("WriteWithTxn CommitTransaction err: %s", err)
			return err
		}
	} else {
		blog.Info("WriteWithTxn Insert err: %s", txnErr)
		err := m.db.AbortTransaction(ctx, cap)
		if err != nil {
			blog.Info("WriteWithTxn AbortTransaction err: %s", err)
			return txnErr
		}
		return txnErr
	}

	return nil
}

// 有分布式事务的Read
func (m *MongoOperator) ReadWithTxn() error {
	ctx, cap, err := m.GenerateTxnKit()
	if err != nil {
		blog.Info("genTxnKit err: %s", err)
		return err
	}

	result := make([]interface{}, 0)
	//filter := map[string]interface{}{"_id": fmt.Sprintf("%d", time.Now().UnixNano())}
	filter := map[string]interface{}{"testread": m.data}
	txnErr := m.db.Table(m.tableName).Find(filter).All(ctx, &result)
	if txnErr == nil {
		err := m.db.CommitTransaction(ctx, cap)
		if err != nil {
			blog.Info("ReadWithTxn CommitTransaction err: %s", err)
			return err
		}
	} else {
		blog.Info("ReadWithTxn Find err: %s", txnErr)
		err := m.db.AbortTransaction(ctx, cap)
		if err != nil {
			blog.Info("ReadWithTxn AbortTransaction err: %s", err)
			return txnErr
		}
		return txnErr
	}

	return nil
}

// 无事务的Write
func (m *MongoOperator) WriteNoTxn() error {
	err := m.db.Table(m.tableName).Insert(context.Background(), map[string]string{"notxnfield": m.data})
	if err != nil {
		blog.Info("WriteNoTxn Insert err: %s", err)
		return err
	}
	return nil
}

// 无事务的Read
func (m *MongoOperator) ReadNoTxn() error {
	result := make([]interface{}, 0)
	//filter := map[string]interface{}{"_id": fmt.Sprintf("%d", time.Now().UnixNano())}
	filter := map[string]interface{}{"testread": m.data}
	err := m.db.Table(m.tableName).Find(filter).All(context.Background(), &result)
	if err != nil {
		blog.Info("ReadNoTxn Find err: %s", err)
		return err
	}
	return nil
}

// 优先读从节点的Read
func (m *MongoOperator) ReadSecondaryPrefer() error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.BKHTTPReadReference, common.SecondaryPreferredMode)
	result := make([]interface{}, 0)
	//filter := map[string]interface{}{"_id": fmt.Sprintf("%d", time.Now().UnixNano())}
	filter := map[string]interface{}{"testread": m.data}
	err := m.db.Table(m.tableName).Find(filter).All(context.Background(), &result)
	if err != nil {
		blog.Info("ReadSecondaryPrefer Find err: %s", err)
		return err
	}
	return nil
}
