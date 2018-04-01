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
	"fmt"
	// "log"
	// "net/http"
	"os"
	"runtime"

	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/app"
	"configcenter/src/scene_server/event_server/app/options"

	_ "configcenter/src/scene_server/event_server/event_service/actions/subscription"
	_ "configcenter/src/scene_server/event_server/event_service/distribution"
	"github.com/spf13/pflag"
	// _ "net/http/pprof"

	"configcenter/src/common"
	"configcenter/src/common/types"
)

func main() {
	common.SetIdentification(types.CC_MODULE_EVENTSERVER)

	runtime.GOMAXPROCS(runtime.NumCPU())

	blog.InitLogs()
	defer blog.CloseLogs()

	op := options.NewServerOption()
	op.AddFlags(pflag.CommandLine)

	util.InitFlags()

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	if err := app.Run(op); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
