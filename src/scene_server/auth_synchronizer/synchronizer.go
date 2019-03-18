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

    "github.com/spf13/pflag"
	
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_synchronizer/app"
	"configcenter/src/scene_server/auth_synchronizer/app/options"
)

func main() {
	common.SetIdentification(types.CC_MODULE_AUTH_SYNCHROIZER)
	runtime.GOMAXPROCS(runtime.NumCPU())

	blog.InitLogs()
	defer blog.CloseLogs()

	serverOptions := options.NewServerOption()
	serverOptions.AddFlags(pflag.CommandLine)

	util.InitFlags()

	if err := common.SavePid(); err != nil {
		blog.Errorf("fail to save pid: err:%s", err.Error())
	}

	if err := app.Run(context.Background(), serverOptions); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		blog.Errorf("process stopped by %v", err)
		blog.CloseLogs()
		os.Exit(1)
	}
}
