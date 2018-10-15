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
	"context"
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/spf13/pflag"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/datacollection/app"
	"configcenter/src/scene_server/datacollection/app/options"
)

func main() {
	common.SetIdentification(types.CC_MODULE_DATACOLLECTION)
	runtime.GOMAXPROCS(runtime.NumCPU())

	var mock bool
	pflag.CommandLine.BoolVar(&mock, "mock", false, "send mock message")

	op := options.NewServerOption()
	op.AddFlags(pflag.CommandLine)
	util.InitFlags()

	if mock {
		if err := sigmock(); err != nil {
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

	if err := app.Run(context.Background(), op); err != nil {
		blog.Fatal(err)
	}
}

func sigmock() error {
	pid, err := common.ReadPid()
	if err != nil {
		return err
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Signal(syscall.SIGUSR1)
}
