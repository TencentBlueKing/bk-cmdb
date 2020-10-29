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

// 作为mongo故障检测用例，用于mongo切主、增加从节点等DB变动时，观察mongo操作的情况和db故障情况
// 执行命令例子：
// go run avaliabitity_check.go -mongo-addr mongodb://localhost:27011,localhost:27012,localhost:27013/test
// -concurrent 10 -sustain-seconds 10

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/test"
	"configcenter/src/test/db/mongo/operator"
)

const (
	chanSize = 100
)

var (
	// db测试用例的wg
	wg sync.WaitGroup
	// 统计用的wg
	swg sync.WaitGroup

	// 停止程序标识
	stopCh = make(chan struct{})
)

type metrics struct {
	// 写主总数
	writeTotal int
	// 写主成功数
	writeSuccess int
	// 写主失败数
	writeFail int
	// 读主总数
	readTotal int
	// 读主成功数
	readSuccess int
	// 读主失败数
	readFail int
	// 读从总数
	readSecTotal int
	// 读从成功数
	readSecSuccess int
	// 读从失败数
	readSecFail int

	// QPS
	writeQPS   int
	readQPS    int
	readSecQPS int

	writeTotalChan     chan int
	writeSuccessChan   chan int
	writeFailChan      chan int
	readTotalChan      chan int
	readSuccessChan    chan int
	readFailChan       chan int
	readSecTotalChan   chan int
	readSecSuccessChan chan int
	readSecFailChan    chan int
}

func NewMetrics() *metrics {
	return &metrics{
		writeTotalChan:     make(chan int, chanSize),
		writeSuccessChan:   make(chan int, chanSize),
		writeFailChan:      make(chan int, chanSize),
		readTotalChan:      make(chan int, chanSize),
		readSuccessChan:    make(chan int, chanSize),
		readFailChan:       make(chan int, chanSize),
		readSecTotalChan:   make(chan int, chanSize),
		readSecSuccessChan: make(chan int, chanSize),
		readSecFailChan:    make(chan int, chanSize),
	}
}

// Output 输出统计数据
func (m *metrics) Output() string {
	return fmt.Sprintf("writeTotal:%d, writeSuccess:%d, writeFail:%d, readTotal:%d, readSuccess:%d"+
		", readFail:%d, readSecTotal:%d, readSecSuccess:%d, readSecFail:%d, writeQPS:%d, readQPS:%d, readSecQPS:%d",
		m.writeTotal, m.writeSuccess, m.writeFail, m.readTotal, m.readSuccess,
		m.readFail, m.readSecTotal, m.readSecSuccess, m.readSecFail, m.writeQPS, m.readQPS, m.readSecQPS)
}

// dbCheck 跑db测试用例
func dbCheck(operator *operator.MongoOperator, metrics *metrics) {

	// 写操作
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				return
			default:
				err := operator.WriteWithTxn()
				metrics.writeTotalChan <- 1
				if err != nil {
					metrics.writeFailChan <- 1
					blog.Errorf("write failed, err:%s", err)
					continue
				}
				metrics.writeSuccessChan <- 1
			}

		}
	}()

	// 读主操作
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				return
			default:
			}
			err := operator.ReadNoTxn()
			metrics.readTotalChan <- 1
			if err != nil {
				metrics.readFailChan <- 1
				blog.Errorf("read primary failed, err:%s", err)
				continue
			}
			metrics.readSuccessChan <- 1
		}
	}()

	// 优先读从操作
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				return
			default:
			}
			err := operator.ReadSecondaryPrefer()
			metrics.readSecTotalChan <- 1
			if err != nil {
				metrics.readSecFailChan <- 1
				blog.Errorf("read secondary prefer failed, err:%s", err)
				continue
			}
			metrics.readSecSuccessChan <- 1
		}
	}()
}

// 操作统计
func statistics(metrics *metrics) {

	swg.Add(1)
	go func() {
		defer swg.Done()
		for range metrics.writeTotalChan {
			metrics.writeTotal++
		}
	}()

	swg.Add(1)
	go func() {
		defer swg.Done()
		for range metrics.writeSuccessChan {
			metrics.writeSuccess++
		}
	}()

	swg.Add(1)
	go func() {
		defer swg.Done()
		for range metrics.writeFailChan {
			metrics.writeFail++
		}
	}()

	swg.Add(1)
	go func() {
		defer swg.Done()
		for range metrics.readTotalChan {
			metrics.readTotal++
		}
	}()

	swg.Add(1)
	go func() {
		defer swg.Done()
		for range metrics.readSuccessChan {
			metrics.readSuccess++
		}
	}()

	swg.Add(1)
	go func() {
		defer swg.Done()
		for range metrics.readFailChan {
			metrics.readFail++
		}
	}()

	swg.Add(1)
	go func() {
		defer swg.Done()
		for range metrics.readSecTotalChan {
			metrics.readSecTotal++
		}
	}()

	swg.Add(1)
	go func() {
		defer swg.Done()
		for range metrics.readSecSuccessChan {
			metrics.readSecSuccess++
		}
	}()

	swg.Add(1)
	go func() {
		defer swg.Done()
		for range metrics.readSecFailChan {
			metrics.readSecFail++
		}
	}()
}

// 打印统计数据
func printStatistcs(metrics *metrics) {
	go func() {
		for {
			select {
			case <-stopCh:
				return
			default:
				// 不打印metrics数据，防止资源竞争，避免加锁，留到最后打印
				//blog.Infof("metrics:%#v", *metrics)
				blog.Info("running")
				// 控制打印频率
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
}

// 优雅退出通知
func exitNotify() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		select {
		case sig := <-exit:
			blog.Infof("receive exit signal %v, begin to exit", sig)
			close(stopCh)
			return
		}
	}()
}

// 到达设置的持续时间后退出程序
func exitAfterSustainSeconds(sustainSeconds float64) {
	go func() {
		select {
		case <-time.After(time.Duration(sustainSeconds) * time.Second):
			close(stopCh)
		}
	}()
}

// closeMetricsChan 关闭channel，让统计退出循环
func closeMetricsChan(metrics *metrics) {
	close(metrics.writeTotalChan)
	close(metrics.writeSuccessChan)
	close(metrics.writeFailChan)
	close(metrics.readTotalChan)
	close(metrics.readSuccessChan)
	close(metrics.readFailChan)
	close(metrics.readSecTotalChan)
	close(metrics.readSecSuccessChan)
	close(metrics.readSecFailChan)
}

// calcQPS 计算QPS
func calcQPS(metrics *metrics, sustainSeconds float64) {
	metrics.writeQPS = int(float64(metrics.writeTotal) / sustainSeconds)
	metrics.readQPS = int(float64(metrics.readTotal) / sustainSeconds)
	metrics.readSecQPS = int(float64(metrics.readSecTotal) / sustainSeconds)
}

func main() {

	operator := operator.NewMongoOperator("cc_mongo_check")
	if err := operator.ClearData(); err != nil {
		blog.Errorf("ClearData failed, err:%s", err)
		return
	}
	metrics := NewMetrics()

	start := time.Now()
	tConf := test.GetTestConfig()
	exitNotify()
	exitAfterSustainSeconds(tConf.SustainSeconds)

	for i := 1; i <= tConf.Concurrent; i++ {
		dbCheck(operator, metrics)
	}

	statistics(metrics)
	printStatistcs(metrics)

	// 等待db测试用例终止，所有写channel操作结束
	wg.Wait()

	closeMetricsChan(metrics)
	// 等待统计数据完成
	swg.Wait()

	calcQPS(metrics, tConf.SustainSeconds)
	blog.Infof("metrics:%s", metrics.Output())
	blog.Infof("running time %dms", time.Since(start)/time.Millisecond)
}
