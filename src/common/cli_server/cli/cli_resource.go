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
 
package cli

import (
	"configcenter/src/common/cli_server/actions"
	"configcenter/src/common/cli_server/config"
	"configcenter/src/common/conf"
	storage "configcenter/src/storage"
	dbcli "configcenter/src/storage/dbclient"
)

type cliResource struct {
	Config string
	Action []*actions.Action
	Mysql  storage.DI
	Monogo storage.DI
}

type CliResource interface {
	InitCli()
	ParseConfig() (map[string]string, error)
	GetDataCli(config map[string]string, dType string) error
	SetConfig(conf *config.CCCliConfig)
	GetAction() []*actions.Action
	GetActionDb(dbType string) storage.DI
}

func NewCliResource() CliResource {
	return &cliResource{}
}

func (c *cliResource) InitCli() {
	cliFuncs := actions.GetCliAction()
	c.Action = append(c.Action, cliFuncs...)
}

func (cc *cliResource) ParseConfig() (map[string]string, error) {
	ccConfig := new(conf.Config)
	ccConfig.InitConfig(cc.Config)

	return ccConfig.Configmap, nil
}

//get data cli
func (cc *cliResource) GetDataCli(config map[string]string, dType string) error {
	host := config[dType+".host"]
	port := config[dType+".port"]
	user := config[dType+".usr"]
	pwd := config[dType+".pwd"]
	dbName := config[dType+".database"]
	mechanism := config[dType+".mechanism"]

	dataCli, err := dbcli.NewDB(host, port, user, pwd, mechanism, dbName, dType)
	if err != nil {
		return err
	}
	err = dataCli.Open()
	if err != nil {
		return err
	}
	if dType == "mysql" {
		cc.Mysql = dataCli
	} else {
		cc.Monogo = dataCli
	}

	return nil
}

func (cc *cliResource) SetConfig(conf *config.CCCliConfig) {
	cc.Config = conf.ExConfig
}

func (cc *cliResource) GetAction() []*actions.Action {
	return cc.Action
}

func (cc *cliResource) GetActionDb(dbType string) storage.DI {
	if "mysql" == dbType || "Mysql" == dbType {
		return cc.Mysql
	} else {
		return cc.Monogo
	}
}
