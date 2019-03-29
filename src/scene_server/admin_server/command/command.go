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
	"context"
	"fmt"
	"os"

	"configcenter/src/common"
	"configcenter/src/common/backbone/configcenter"
	"configcenter/src/storage/dal/mongo"

	"github.com/spf13/pflag"
)

const bkbizCmdName = "bkbiz"

const (
	scopeAll = "all"
)

// Parse run app command
func Parse(args []string) error {
	ctx := context.Background()
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
		bizName        string
		scope          string
	)

	// set flags
	bkbizfs := pflag.NewFlagSet(bkbizCmdName, pflag.ExitOnError)
	bkbizfs.BoolVar(&dryrunflag, "dryrun", false, "dryrun flag, if this flag seted, we will just print what we will do but not execute to db")
	bkbizfs.BoolVar(&exportflag, "export", false, "export flag")
	bkbizfs.BoolVar(&miniflag, "mini", false, "mini flag, only export required fields")
	bkbizfs.BoolVar(&importflag, "import", false, "import flag")
	bkbizfs.StringVar(&scope, "scope", scopeAll, "export scope, could be [biz] or [process], default all")
	bkbizfs.StringVar(&filepath, "file", "", "export/import filepath")
	bkbizfs.StringVar(&configposition, "config", "conf/api.conf", "The config path. e.g conf/api.conf")
	bkbizfs.StringVar(&bizName, "biz_name", "蓝鲸", "export/import the specified business topo")
	err := bkbizfs.Parse(args[1:])
	if err != nil {
		return err
	}

	// init config
	pconfig, err := configcenter.ParseConfigWithFile(configposition)
	if nil != err {
		return fmt.Errorf("parse config file error %s", err.Error())
	}
	config := mongo.ParseConfigFromKV("mongodb", pconfig.ConfigMap)
	// connect to mongo db
	db, err := mongo.NewMgo(config.BuildURI(), 0)
	if err != nil {
		return fmt.Errorf("connect mongo server failed %s", err.Error())
	}
	opt := &option{
		position: filepath,
		OwnerID:  common.BKDefaultOwnerID,
		dryrun:   dryrunflag,
		mini:     miniflag,
		scope:    scope,
		bizName:  bizName,
	}

	if exportflag {
		var mode string
		if miniflag {
			mode = "mini"
		} else {
			mode = "verbose"

		}
		fmt.Printf("exporting %s business to %s in \033[34m%s\033[0m mode\n", bizName, filepath, mode)
		if err := export(ctx, db, opt); err != nil {
			fmt.Printf("export error: %s", err.Error())
			os.Exit(2)
		}
		fmt.Printf("blueking %s has been export to %s\n", bizName, filepath)
	} else if importflag {
		if dryrunflag {
			fmt.Printf("dryrun import %s business from %s\n", bizName, filepath)
		} else {
			fmt.Printf("importing %s business from %s\n", bizName, filepath)
		}
		opt.mini = false
		opt.scope = scopeAll
		if err := importBKBiz(ctx, db, opt); err != nil {
			fmt.Printf("import error: %s", err.Error())
			os.Exit(2)
		}
		if !dryrunflag {
			fmt.Printf("%s business has been import from %s\n", bizName, filepath)
		}
	} else {
		fmt.Printf("invalide argument")
	}

	os.Exit(0)
	return nil
}
