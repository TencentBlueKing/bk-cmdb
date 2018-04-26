package main

import (
	"configcenter/src/common"
	"configcenter/src/common/util"
	"configcenter/src/framework/api"
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

	// init framework
	api.Init()

	(&option.Options{}).AddFlags(pflag.CommandLine)

	util.InitFlags()

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
