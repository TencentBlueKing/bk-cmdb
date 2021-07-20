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
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/storage/dal/mongo/local"

	"github.com/spf13/pflag"
)

const bkbizCmdName = "bkbiz"

const (
	scopeAll = "all"
)

// Parse run app command
func Parse(args []string) error {
	ctx := context.Background()
	var (
		exportFlag     bool
		importFlag     bool
		miniFlag       bool
		dryRunFlag     bool
		filePath       string
		configPosition string
		bizName        string
		scope          string
	)

	if len(args) <= 1 || args[1] != bkbizCmdName {
		return nil
	}

	// set flags
	cmdFlags := pflag.NewFlagSet(bkbizCmdName, pflag.ExitOnError)
	cmdFlags.BoolVar(&dryRunFlag, "dryrun", false, "dryrun flag, if this flag seted, we will just print what we will do but not execute to db")
	cmdFlags.BoolVar(&exportFlag, "export", false, "export flag")
	cmdFlags.BoolVar(&miniFlag, "mini", false, "mini flag, only export required fields")
	cmdFlags.BoolVar(&importFlag, "import", false, "import flag")
	cmdFlags.StringVar(&scope, "scope", "all", "export scope, could be [biz] or [process], default all")
	cmdFlags.StringVar(&filePath, "file", "", "export/import filepath")
	cmdFlags.StringVar(&configPosition, "config", "conf/api.conf", "The config path. e.g conf/api.conf")
	cmdFlags.StringVar(&bizName, "biz_name", "蓝鲸", "export/import the specified business topo")
	err := cmdFlags.Parse(args[1:])

	if err != nil {
		return err
	}

	// read config
    if err := cc.SetMigrateFromFile(configPosition); err != nil {
		return fmt.Errorf("parse config file error %s", err.Error())
	}
	mongoConfig, err := cc.Mongo("mongodb")
	if err != nil {
		return err
	}

	// connect to mongo db
	db, err := local.NewMgo(mongoConfig.GetMongoConf(), 0)
	if err != nil {
		return fmt.Errorf("connect mongo server failed %s", err.Error())
	}
	opt := &option{
		position: filePath,
		OwnerID:  common.BKDefaultOwnerID,
		dryrun:   dryRunFlag,
		mini:     miniFlag,
		scope:    scope,
		bizName:  bizName,
	}

	if exportFlag {
		var mode string
		if miniFlag {
			mode = "mini"
		} else {
			mode = "verbose"

		}
		fmt.Printf("exporting %s business to %s in \033[34m%s\033[0m mode\n", bizName, filePath, mode)
		if err := export(ctx, db, opt); err != nil {
			fmt.Printf("export error: %s\n", err.Error())
			os.Exit(2)
		}
		fmt.Printf("blueking %s has been export to %s\n", bizName, filePath)
	} else if importFlag {
		if dryRunFlag {
			fmt.Printf("dryrun import %s business from %s\n", bizName, filePath)
		} else {
			fmt.Printf("importing %s business from %s\n", bizName, filePath)
		}
		opt.mini = false
		opt.scope = scopeAll
		if err := importBKBiz(ctx, db, opt); err != nil {
			fmt.Printf("import error: %s\n", err.Error())
			os.Exit(2)
		}
		if !dryRunFlag {
			fmt.Printf("%s business has been import from %s\n", bizName, filePath)
		}
	} else {
		fmt.Printf("invalide argument")
	}

	os.Exit(0)
	return nil
}
