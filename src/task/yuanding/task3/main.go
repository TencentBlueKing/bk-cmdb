package main

import (
	"time"

	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"configcenter/src/common"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/framework/api"
	fcommon "configcenter/src/framework/common"
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/discovery"
	"configcenter/src/framework/core/httpserver"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/monitor/metric"
	"configcenter/src/framework/core/option"
	"configcenter/src/framework/core/output/module/client"
	_ "configcenter/src/framework/plugins"

	// load all plugins

	"github.com/spf13/pflag"
)

// APPNAME the name of this application, will be use as identification mark for monitoring
const APPNAME = "YdApp"


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

	metricManager := metric.NewManager(opt)

	server.RegisterActions(api.Actions()...)
	server.RegisterActions(metricManager.Actions()...)

	// 查询主机数量
	cond := fcommon.CreateCondition()
	items, err := client.GetClient().CCV3(client.Params{SupplierAccount: config.Get().Get("core.supplierAccount")}).Host().SearchHost(cond)
	if nil != err {
		log.Errorf("error info is %s\n", err.Error())
	}
	log.Infof("host num is %d\n", len(items))

	httpChan := make(chan error, 1)
	go func() { httpChan <- server.ListenAndServe() }()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-httpChan:
		log.Errorf("http exit, error: %v", err)
		return
	case s := <-sigs:
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("the signal:", s.String())
		//case syscall.SIGURG:
		//	// the reserved
		//case syscall.SIGUSR1:
		//	// the reserved
		//case syscall.SIGUSR2:
		default:
			fmt.Printf("\nunknown the signal (%s) \n", s.String())
		}
	}

}
