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

package mongo_test

import (
	"configcenter/src/common/blog"
	"configcenter/src/test/db/mongo/operator"
	"configcenter/src/test/run"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DB Operation Load Test", func() {
	tableName := "cc_tranTest"
	// 实例化mongo操作者对象
	operator := operator.NewMongoOperator(tableName)
	// 清空数据
	err := operator.ClearData()
	if err := operator.ClearData(); err != nil {
		blog.Errorf("ClearData err:%s", err)
	}
	Expect(err).Should(BeNil())

	Describe("with txn load test", func() {
		It("Write", func() {
			m := run.FireLoadTest(func() error {
				err := operator.WriteWithTxn()
				if err != nil {
					return err
				}
				return nil
			})
			fmt.Printf("write with txn load perform: \n%s", m.Format())
		})

		It("Read", func() {
			m := run.FireLoadTest(func() error {
				err := operator.ReadWithTxn()
				if err != nil {
					return err
				}
				return nil
			})
			fmt.Printf("read with txn load perform: \n%s", m.Format())
		})

	})

	Describe("no txn load test", func() {
		It("Write", func() {
			m := run.FireLoadTest(func() error {
				err := operator.WriteNoTxn()
				if err != nil {
					return err
				}
				return nil
			})
			fmt.Printf("write no txn load perform: \n%s", m.Format())
		})

		It("Read", func() {
			m := run.FireLoadTest(func() error {
				err := operator.ReadNoTxn()
				if err != nil {
					return err
				}
				return nil
			})
			fmt.Printf("read no txn load perform: \n%s", m.Format())
		})

	})
})
