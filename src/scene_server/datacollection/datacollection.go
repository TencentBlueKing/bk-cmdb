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

func main() {
	common.SetIdentification(types.CC_MODULE_DATACOLLECTION)
	runtime.GOMAXPROCS(runtime.NumCPU())

	var mock bool
	var collector string
	pflag.CommandLine.BoolVar(&mock, "mock", false, "send mock message")
	pflag.CommandLine.StringVar(&collector, "collector", "", "collector name that send mock message send to")

	op := options.NewServerOption()
	op.AddFlags(pflag.CommandLine)
	util.InitFlags()

	if mock {
		if err := sigMock(collector); err != nil {
			fmt.Printf("sigmock failed %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("sigmock success\n")
		return
	}

	blog.InitLogs()
	defer blog.CloseLogs()
	if err := common.SavePid(); err != nil {
		blog.Error("fail to save pid. err: %s", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	if err := app.Run(ctx, cancel, op); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		blog.Errorf("process stopped by %v", err)
		blog.CloseLogs()
		os.Exit(1)
	}
}

func sigMock(collector string) error {
	body := bytes.NewBufferString(`{"name":"` + collector + `"}`)
	resp, err := http.Post("http://127.0.0.1:12140", "application/json", body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s", respBody)
	}
	return nil
}
