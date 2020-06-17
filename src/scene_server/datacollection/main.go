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

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/datacollection/app"
	"configcenter/src/scene_server/datacollection/app/options"

	"github.com/spf13/pflag"
)

// variables here for your common logics.
var (
	// options
	op *options.ServerOption

	// needSendMock flags if need to send mock message or not.
	needSendMock bool

	// mockCollectorName is name of mock collector.
	mockCollectorName string

	// mockDataCollectioAddr is address of target mock datacollection.
	mockDataCollectioAddr string
)

// parseFlags parses all flags for the module. DO NOT add flags anywhere,
// make it a centralized flags handler here.
func parseFlags() {
	// defines flags with specified name, default value, and usage string here.
	pflag.CommandLine.BoolVar(&needSendMock, "mock", false, "flag if need to send mock message or not.")
	pflag.CommandLine.StringVar(&mockCollectorName, "collector", "", "collector name that mock message send to.")
	pflag.CommandLine.StringVar(&mockDataCollectioAddr, "target", "http://127.0.0.1:12140", "address of target mock datacollection.")

	// module options.
	op = options.NewServerOption()

	// add ex-flags to main module options.
	op.AddFlags(pflag.CommandLine)

	// normalizes and parses the command line flags.
	util.InitFlags()
}

// mock cmd.
func mock() {
	if !needSendMock {
		return
	}

	// mock now.
	if err := sendMock(mockCollectorName); err != nil {
		fmt.Printf("mock failed, %+v\n", err)
		os.Exit(1)
	}

	fmt.Printf("mock success!\n")
	os.Exit(0)
}

// send mock message to target datacollection.
func sendMock(collectorName string) error {
	resp, err := http.Post(mockDataCollectioAddr, "application/json", bytes.NewBufferString(`{"name":"`+collectorName+`"}`))
	if err != nil {
		return fmt.Errorf("send mock message, %+v", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		// mock failed.
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("send mock message, status code[%d], %+v", resp.StatusCode, err)
		}
		return fmt.Errorf("send mock message, status code[%d], %+v", resp.StatusCode, respBody)
	}
	return nil
}

func main() {
	// setup system runtimes.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// set module datacollection id.
	common.SetIdentification(types.CC_MODULE_DATACOLLECTION)

	// parse flags.
	parseFlags()

	// mock.
	mock()

	// init logger.
	blog.InitLogs()
	defer blog.CloseLogs()

	// save local module pid.
	if err := common.SavePid(); err != nil {
		blog.Warnf("save module local PID, %+v", err)
	}

	// main context for app.
	ctx, cancel := context.WithCancel(context.Background())

	// run the app now.
	if err := app.Run(ctx, cancel, op); err != nil {
		blog.Errorf("app stopped by %+v", err)
		os.Exit(1)
	}
}
