package main

import (
	"fmt"
	"os"
	"runtime"

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

	blog.InitLogs()
	defer blog.CloseLogs()

	op := options.NewServerOption()
	op.AddFlags(pflag.CommandLine)

	util.InitFlags()

	if err := app.Run(op); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
