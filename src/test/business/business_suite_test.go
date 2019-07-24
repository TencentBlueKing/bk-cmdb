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

package business_test

import (
    "flag"
    "fmt"
    "testing"
    "time"
    "encoding/json"

    "configcenter/src/test/run"
    "configcenter/src/apimachinery"
    "configcenter/src/apimachinery/discovery"
    "configcenter/src/apimachinery/util"
    "configcenter/src/common/backbone/service_mange/zk"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var clientSet apimachinery.ClientSetInterface
var TConfig TestConfig

type TestConfig struct {
	ZkAddr         string
	Concurrent     int
	SustainSeconds int
}

func init() {
	flag.StringVar(&TConfig.ZkAddr, "zk-addr", "127.0.0.1:2181", "zk discovery addresses, comma separated.")
	flag.IntVar(&TConfig.Concurrent, "concurrent", 100, "concurrent request during the load test.")
	flag.IntVar(&TConfig.SustainSeconds, "sustain-seconds", 10, "the load test sustain time in seconds ")
	flag.Parse()
	
	run.Concurrent = TConfig.Concurrent
	run.SustainSeconds = TConfig.SustainSeconds
	
	
}

func TestBusiness(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Business Suite")
}

var _ = BeforeSuite(func() {
	fmt.Println("before suit")
    js, _ := json.MarshalIndent(TConfig, "", "    ")
    fmt.Printf("test config: %s\n", run.SetRed(string(js)))
	client := zk.NewZkClient(TConfig.ZkAddr, 5*time.Second)
	Expect(client.Start()).Should(BeNil())
	Expect(client.Ping()).Should(BeNil())
	disc, err := discovery.NewServiceDiscovery(client)
	Expect(err).Should(BeNil())
	c := &util.APIMachineryConfig{
		QPS:       20000,
		Burst:     10000,
		TLSConfig: nil,
	}
	clientSet, err = apimachinery.NewApiMachinery(c, disc)
	Expect(err).Should(BeNil())
	// wait for get the apiserver address.
	time.Sleep(1 * time.Second)
	
	fmt.Println("**** initialize clientSet success ***")
})
