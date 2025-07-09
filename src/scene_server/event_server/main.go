/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package main

import (
	"context"
	"os"
	"runtime"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/app"
	"configcenter/src/scene_server/event_server/app/options"

	"github.com/spf13/pflag"
)

func main() {
	// setup system runtimes.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// set module datacollection id.
	common.SetIdentification(types.CC_MODULE_EVENTSERVER)

	// init logger.
	blog.InitLogs()
	defer blog.CloseLogs()

	// parse flags.
	op := options.NewServerOption()
	op.AddFlags(pflag.CommandLine)
	util.InitFlags()

	// main context for app.
	ctx, cancel := context.WithCancel(context.Background())

	// run the app now.
	if err := app.Run(ctx, cancel, op); err != nil {
		blog.Errorf("process stopped by %v", err)
		blog.CloseLogs()
		os.Exit(1)
	}
}
