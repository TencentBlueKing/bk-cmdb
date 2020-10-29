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

package main

// 使用压力测试，在mongo切主、增加从节点等DB变动时，观察mongo操作的情况和db故障情况

import (
	"fmt"
	"sync"
	"time"

	"configcenter/src/test/run"
	"configcenter/src/common/blog"
	"configcenter/src/test/db/mongo/operator"
)

var (
	// db测试用例的lwg
	lwg sync.WaitGroup
)


// dbLoadCheck 跑db测试用例
func dbLoadCheck(operator *operator.MongoOperator) {

	// 写操作
	lwg.Add(1)
	go func() {
		defer lwg.Done()
		m := run.FireLoadTest(func() error {
			err := operator.WriteWithTxn()
			if err != nil {
				return err
			}
			return nil
		})
		fmt.Printf("write primary load perform: \n%s", m.Format())
	}()

	// 读主操作
	lwg.Add(1)
	go func() {
		defer lwg.Done()
		m := run.FireLoadTest(func() error {
			err := operator.ReadNoTxn()
			if err != nil {
				return err
			}
			return nil
		})
		fmt.Printf("read primary load perform: \n%s", m.Format())
	}()

	// 优先读从操作
	lwg.Add(1)
	go func() {
		defer lwg.Done()
		m := run.FireLoadTest(func() error {
			err := operator.ReadSecondaryPrefer()
			if err != nil {
				return err
			}
			return nil
		})
		fmt.Printf("read secondary load perform: \n%s", m.Format())
	}()
}


func main() {

	operator := operator.NewMongoOperator("cc_mongo_check")
	if err := operator.ClearData(); err != nil {
		blog.Errorf("ClearData failed, err:%s", err)
		return
	}

	start := time.Now()
	dbLoadCheck(operator)
	lwg.Wait()
	blog.Infof("running time %dms", time.Since(start)/time.Millisecond)

}
