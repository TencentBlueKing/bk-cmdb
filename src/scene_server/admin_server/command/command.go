package command

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/conf"
	"configcenter/src/common/core/cc/api"
	"github.com/spf13/pflag"
	"os"
)

type Command struct {
}

const bkbizCmdName = "bkbiz"

func Parse(args []string) error {
	if len(args) < 4 || args[1] != bkbizCmdName {
		blog.Error("args %d,arg = %s", len(args), args[1])
		return nil
	}

	var (
		exportflag     bool
		importflag     bool
		filepath       string
		configposition string
	)

	// set flags
	bkbizfs := pflag.NewFlagSet(bkbizCmdName, pflag.ExitOnError)
	bkbizfs.BoolVar(&exportflag, "export", false, "export flag")
	bkbizfs.BoolVar(&importflag, "import", false, "import flag")
	bkbizfs.StringVar(&filepath, "file", "", "export or import filepath")
	bkbizfs.StringVar(&configposition, "config", "conf/api.conf", "The config path. e.g conf/api.conf")
	err := bkbizfs.Parse(args[1:])
	if err != nil {
		return err
	}

	// init config
	config := new(conf.Config)
	config.InitConfig(configposition)

	// connect to mongo db
	a := api.NewAPIResource()
	err = a.GetDataCli(config.Configmap, "mongodb")
	if err != nil {
		blog.Error("connect mongodb error exit! err:%s", err.Error())
		return err
	}

	if exportflag {
		if err := export(a.InstCli, &option{position: filepath, OwnerID: common.BKDefaultOwnerID}); err != nil {
			blog.Errorf("export error: %s", err.Error())
			os.Exit(2)
		}
		blog.Infof("blueking business has been export to %s", filepath)
	}

	os.Exit(0)
	return nil
}
