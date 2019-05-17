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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/framework/api"
	cccommon "configcenter/src/framework/common"
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/discovery"
	"configcenter/src/framework/core/httpserver"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/option"
	"configcenter/src/framework/core/output/module/client"
	_ "configcenter/src/framework/plugins"
	"fmt"
	"runtime"

	// load all plugins

	"github.com/spf13/pflag"
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
			blog.Infof("%v", args)
		},
		Infof:  blog.Infof,
		Fatal:  blog.Fatal,
		Fatalf: blog.Fatalf,
		Error: func(args ...interface{}) {
			blog.Errorf("%v", args)
		},
		Errorf:   blog.Errorf,
		Warningf: blog.Warnf,
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
		disClient := zk.NewZkClient(opt.Regdiscv, 5*time.Second)
		if err := disClient.Start(); err != nil {
			log.Errorf("connect regdiscv [%s] failed: %v", opt.Regdiscv, err)
			return
		}
		if err := disClient.Ping(); err != nil {
			log.Errorf("connect regdiscv [%s] failed: %v", opt.Regdiscv, err)
			return
		}
		rd := discovery.NewRegDiscover(APPNAME, disClient, server.GetAddr(), server.GetPort(), false)
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

	// initial the background framework manager.
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

	// statics host by biz
	fmt.Println("==== begin to static host by biz ====")
	ccClient := client.GetClient()
	params := client.Params{SupplierAccount: "0", UserName: "admin"}
	ccv3 := ccClient.CCV3(params)
	// list biz
	businesses, err := ccv3.Business().SearchBusiness(cccommon.CreateCondition())
	if nil != err {
		fmt.Println("invoke cc v3 SearchBusiness api err")
	} else {
		// fmt.Println(businesses)
		for _, v := range businesses {
			bkBizId, _ := v.Int("bk_biz_id")
			hostCond := cccommon.CreateCondition()
			// hostCond.Field("bk_biz_id").Eq(2) //.
			//	Field("bk_supplier_account").Eq(0)
			//fmt.Println(hostCond.ToMapStr())
			hosts, err := ccv3.Host().SearchHost(hostCond)
			if nil != err {
				fmt.Println("invoke cc v3 SearchHost api err")
			} else {
				hostsCount := len(hosts)
				fmt.Printf("biz id: %d, biz name: %s, host count: %d\n",
					bkBizId, v.String("bk_biz_name"), hostsCount)
			}
		}
	}

	// test invoke search_host by ip is ok, by biz is wrong
	// hostCond := cccommon.CreateCondition()
	// hostCond.Field("bk_host_innerip").Eq("172.21.99.8")
	// hostCond.Field("bk_biz_id").Eq("3")
	// fmt.Println(hostCond.ToMapStr())
	// hosts, _ := ccv3.Host().SearchHost(hostCond)
	// fmt.Println(hosts)
	// fmt.Println(len(hosts))
}
