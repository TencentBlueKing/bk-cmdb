package command

import (
	"fmt"
	"os"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/conf"
	"configcenter/src/common/core/cc/api"

	"github.com/spf13/pflag"
)

const bkbizCmdName = "bkbiz"

// Parse run app command
func Parse(args []string) error {
	if len(args) <= 1 || args[1] != bkbizCmdName {
		return nil
	}

	var (
		exportflag     bool
		importflag     bool
		dryrun         bool
		filepath       string
		configposition string
	)

	// set flags
	bkbizfs := pflag.NewFlagSet(bkbizCmdName, pflag.ExitOnError)
	bkbizfs.BoolVar(&dryrun, "dryrun", false, "dryrun flag, if this flag seted, we will just print what we will do but not execute to db")
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
		fmt.Printf("blueking business has been export to %s\n", filepath)
	} else if importflag {
		fmt.Printf("importing blueking business from %s\n", filepath)
		if err := importBKBiz(a.InstCli, &option{position: filepath, OwnerID: common.BKDefaultOwnerID, dryrun: dryrun}); err != nil {
			blog.Errorf("import error: %s", err.Error())
			os.Exit(2)
		}
		if !dryrun {
			fmt.Printf("blueking business has been import from %s\n", filepath)
		}
	} else {
		blog.Errorf("invalide argument")
	}

	os.Exit(0)
	return nil
}
