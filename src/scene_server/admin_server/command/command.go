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

package command

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"configcenter/src/common"
	"configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/storage/mgoclient"
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
		miniflag       bool
		dryrunflag     bool
		filepath       string
		configposition string
		scope          string
	)

	// set flags
	bkbizfs := pflag.NewFlagSet(bkbizCmdName, pflag.ExitOnError)
	bkbizfs.BoolVar(&dryrunflag, "dryrun", false, "dryrun flag, if this flag seted, we will just print what we will do but not execute to db")
	bkbizfs.BoolVar(&exportflag, "export", false, "export flag")
	bkbizfs.BoolVar(&miniflag, "mini", false, "mini flag, only export required fields")
	bkbizfs.BoolVar(&importflag, "import", false, "import flag")
	bkbizfs.StringVar(&scope, "scope", "all", "export scope, could be [biz] or [process], default all")
	bkbizfs.StringVar(&filepath, "file", "", "export or import filepath")
	bkbizfs.StringVar(&configposition, "config", "conf/api.conf", "The config path. e.g conf/api.conf")
	err := bkbizfs.Parse(args[1:])
	if err != nil {
		return err
	}

	// init config
	pconfig, err := configcenter.ParseConfigWithFile(configposition)
	if nil != err {
		return fmt.Errorf("parse config file error %s", err.Error())
	}
	config := mgoclient.NewMongoConfig(pconfig.ConfigMap)
	// connect to mongo db
	db, err := mgoclient.NewFromConfig(*config)
	if err != nil {
		return fmt.Errorf("connect mongo server failed %s", err.Error())
	}
	err = db.Open()
	if err != nil {
		return fmt.Errorf("connect mongo server failed %s", err.Error())
	}

	opt := &option{
		position: filepath,
		OwnerID:  common.BKDefaultOwnerID,
		dryrun:   dryrunflag,
		mini:     miniflag,
		scope:    scope,
	}

	if exportflag {
		mode := ""
		if miniflag {
			mode = "mini"
		} else {
			mode = "verbose"

		}
		fmt.Printf("exporting blueking business to %s in \033[34m%s\033[0m mode\n", filepath, mode)
		if err := export(db, opt); err != nil {
			blog.Errorf("export error: %s", err.Error())
			os.Exit(2)
		}
		fmt.Printf("blueking business has been export to %s\n", filepath)
	} else if importflag {
		fmt.Printf("importing blueking business from %s\n", filepath)
		opt.mini = false
		opt.scope = "all"
		if err := importBKBiz(db, opt); err != nil {
			blog.Errorf("import error: %s", err.Error())
			os.Exit(2)
		}
		if !dryrunflag {
			fmt.Printf("blueking business has been import from %s\n", filepath)
		}
	} else {
		blog.Errorf("invalide argument")
	}

	os.Exit(0)
	return nil
}
