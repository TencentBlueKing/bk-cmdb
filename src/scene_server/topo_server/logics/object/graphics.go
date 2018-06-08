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
 
package object

import (
	"configcenter/src/common/errors"
	"configcenter/src/scene_server/topo_server/manager"
	api "configcenter/src/source_controller/api/object"
)

type topoGraphicsLogic struct {
	objcli *api.Client
	cfg    manager.Configer
	mgr    manager.Manager
}

var _ manager.TopoGraphicsLogic = (*topoGraphicsLogic)(nil) // check the interface

func init() {
	obj := &topoGraphicsLogic{}

	obj.objcli = api.NewClient("")
	manager.SetManager(obj)
	manager.RegisterLogic(manager.TopoGraphics, obj)
}

// Set implement SetConfiger interface
func (cli *topoGraphicsLogic) Set(cfg manager.Configer) {
	cli.cfg = cfg
}

// SetManager implement the manager's Hooker interface
func (cli *topoGraphicsLogic) SetManager(mgr manager.Manager) error {
	cli.mgr = mgr
	return nil
}

func (cli *topoGraphicsLogic) SearchGraphics(forward *api.ForwardParam, params *api.TopoGraphics, errProxy errors.DefaultCCErrorIf) ([]api.TopoGraphics, error) {
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.SearchTopoGraphics(forward, params)
}

func (cli *topoGraphicsLogic) UpdateGraphics(forward *api.ForwardParam, params []api.TopoGraphics, errProxy errors.DefaultCCErrorIf) error {
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.UpdateGraphics(forward, params)
}
