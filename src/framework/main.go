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
	"configcenter/src/common"
	"configcenter/src/common/util"
	"configcenter/src/framework/api"
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/options"
	"fmt"
	"github.com/spf13/pflag"

	"os"
	"os/signal"
	"runtime"
	"syscall"

	_ "configcenter/src/framework/plugins" // load all plugins
)

func setParams() {

}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	opt := &option.Options{}
	opt.AddFlags(pflag.CommandLine)
	util.InitFlags()

	if err := config.Init(opt); err != nil {
		panic(err)
	}

	// config.Get()

	// init the framework
	if err := common.SavePid(); nil != err {
		fmt.Printf("\n can not save the pidfile, error info is %s\n", err.Error())
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	for s := range sigs {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("the signal:", s.String())
			goto end
		case syscall.SIGURG:
			// the reserved
		case syscall.SIGUSR1:
			// the reserved
		case syscall.SIGUSR2:
			// the reserved
		default:
			fmt.Printf("\nunknown the signal (%s) \n", s.String())
		}

	}

end:
	// unint the framework
	api.UnInit()
}
