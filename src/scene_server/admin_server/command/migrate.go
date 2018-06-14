package command

import (
	"configcenter/src/common/blog"
	"configcenter/src/common/conf"
	"configcenter/src/common/core/cc/api"
	"github.com/spf13/pflag"
)

type Command struct {
}

const bkbizCmdName = "bkbiz"

func Parse(args []string) error {
	if len(args) < 4 && args[1] != bkbizCmdName {
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
	err := bkbizfs.Parse(args[2:])
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
		return export(filepath, a.InstCli)
	}

	return nil
}
