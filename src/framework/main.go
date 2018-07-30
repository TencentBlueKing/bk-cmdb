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
	"configcenter/src/framework/core/discovery"
	"configcenter/src/framework/core/httpserver"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/monitor/metric"
	"configcenter/src/framework/core/option"
	"configcenter/src/framework/core/output/module/client"
	"time"

	"configcenter/src/common/blog"

	"fmt"
	"github.com/spf13/pflag"

	"os"
	"os/signal"
	"runtime"
	"syscall"

	_ "configcenter/src/framework/plugins" // load all plugins
)

// APPNAME the name of this application, will be use as identification mark for monitoring
const APPNAME = "DemoApp"

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	opt := &option.Options{AppName: APPNAME}
	opt.AddFlags(pflag.CommandLine)
	util.InitFlags()

	blog.InitLogs()

	log.SetLoger(&log.Logger{
		Info: func(args ...interface{}) {
			blog.Info("%v", args)
		},
		Infof:  blog.Infof,
		Fatal:  blog.Fatal,
		Fatalf: blog.Fatalf,
		Error: func(args ...interface{}) {
			blog.Error("%v", args)
		},
		Errorf: blog.Errorf,
	})

	if err := config.Init(opt); err != nil {
		log.Errorf("init config error: %v", err)
		return
	}

	server, err := httpserver.NewServer(opt)
	if err != nil {
		log.Errorf("NewServer error: %v", err)
		return
	}

	if "" != opt.Regdiscv {
		rd := discovery.NewRegDiscover(APPNAME, opt.Regdiscv, server.GetAddr(), server.GetPort(), false)
		go func() {
			rd.Start()
		}()
		for {
			_, err := rd.GetApiServ()
			if err == nil {
				break
			}
			log.Errorf("there is no api server, will reget it after 2s")
			time.Sleep(time.Second * 2)
		}
		client.NewForConfig(config.Get(), rd)
	} else {
		client.NewForConfig(config.Get(), nil)
	}

	api.Init()
	defer func() {
		blog.CloseLogs()
		api.UnInit()
	}()

	// init the framework
	if err := common.SavePid(); nil != err {
		fmt.Printf("\n can not save the pidfile, error info is %s\n", err.Error())
		return
	}

	metricManager := metric.NewManager(opt)

	server.RegisterActions(api.Actions()...)
	server.RegisterActions(metricManager.Actions()...)

	httpChan := make(chan error)
	go func() { httpChan <- server.ListenAndServe() }()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	select {
	case err := <-httpChan:
		log.Errorf("http exit, error: %v", err)
		return
	case s := <-sigs:
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("the signal:", s.String())
		case syscall.SIGURG:
			// the reserved
		case syscall.SIGUSR1:
			// the reserved
		case syscall.SIGUSR2:
		default:
			fmt.Printf("\nunknown the signal (%s) \n", s.String())
		}
	}

}
