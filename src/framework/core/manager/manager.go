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

package manager

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"configcenter/src/framework/core/httpserver"
	"configcenter/src/framework/core/input"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/output"
	"configcenter/src/framework/core/types"
	"github.com/emicklei/go-restful"
)

// Manager contains the basic framework data.
type Manager struct {
	cancel      context.CancelFunc
	ms          []Action
	OutputerMgr output.Manager
	InputerMgr  input.Manager
}

// Actions returns metricActions
func (m *Manager) Actions() []httpserver.Action {
	var httpactions []httpserver.Action
	for _, a := range m.ms {
		httpactions = append(httpactions, httpserver.Action{Method: a.Method, Path: a.Path, Handler: func(req *restful.Request, resp *restful.Response) {

			value, err := ioutil.ReadAll(req.Request.Body)
			if err != nil {
				log.Errorf("read http request body failed, error:%s", err.Error())
				return
			}

			mData := types.MapStr{}
			if err := json.Unmarshal(value, &mData); nil != err {
				log.Errorf("failed to unmarshal the data, error %s", err.Error())
				return
			}

			data, dataErr := a.HandlerFunc(mData)
			if nil != dataErr {
				log.Errorf("%s", dataErr.Error())
			}

			// TODO:需要处理返回值的情况
			if nil != data {
				_ = data
			}

		}})
	}
	return httpactions
}

// CreateFrameworkContext create a new framework context instance
func (cli *Manager) CreateFrameworkContext() FrameworkContext {
	return cli
}

// stop used to stop the business cycles.
func (cli *Manager) stop() error {

	if nil != cli.cancel {
		cli.cancel()
	}

	return cli.InputerMgr.Stop()
}

// Run start the business cycle until the stop method is called.
func (cli *Manager) Run(ctx context.Context, cancel context.CancelFunc) {

	cli.cancel = cancel

	cli.InputerMgr.Run(ctx, cli)
}
