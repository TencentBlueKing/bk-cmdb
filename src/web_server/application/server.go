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

package app

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	_ "configcenter/src/web_server/application/controllers"
	"configcenter/src/web_server/application/options"
	"configcenter/src/web_server/application/service"
)

//Run cc server
func Run(op *options.ServerOption) error {

	setConfig(op)

	serv, err := service.NewCCWebServer(op.ServConf)
	if err != nil {
		blog.Error("fail to create ccapi server. err:%s", err.Error())
		return err
	}

	//pid
	if err := common.SavePid(); err != nil {
		blog.Error("fail to save pid: err:%s", err.Error())
	}

	return serv.Start()
}

func setConfig(op *options.ServerOption) {
	//server cert directory

}
