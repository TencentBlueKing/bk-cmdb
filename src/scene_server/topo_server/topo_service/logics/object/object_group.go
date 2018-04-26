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
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/scene_server/topo_server/topo_service/manager"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
)

type objAttGroupLogic struct {
	objcli *api.Client
	cfg    manager.Configer
	mgr    manager.Manager
}

var _ manager.ObjectAttGroupLogic = (*objAttGroupLogic)(nil)

func init() {
	obj := &objAttGroupLogic{}
	obj.objcli = api.NewClient("")
	manager.SetManager(obj)
	manager.RegisterLogic(manager.ObjectGroup, obj)
}

// Set implement SetConfiger interface
func (cli *objAttGroupLogic) Set(cfg manager.Configer) {
	cli.cfg = cfg
}

// SetManager implement the manager's Hooker interface
func (cli *objAttGroupLogic) SetManager(mgr manager.Manager) error {
	cli.mgr = mgr
	return nil
}

func (cli *objAttGroupLogic) CreateObjectGroup(forward *api.ForwardParam, data []byte, errProxy errors.DefaultCCErrorIf) (int, error) {
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	var selector api.ObjAttGroupDes
	jsErr := json.Unmarshal(data, &selector)
	if nil != jsErr {
		blog.Error("can not unmarshal the data (%s), error info is %s", string(data), jsErr.Error())
		return 0, jsErr
	}
	val := api.ObjAttGroupDes{}
	val.GroupIndex = selector.GroupIndex
	val.GroupName = selector.GroupName
	val.GroupID = selector.GroupID
	valstr, _ := json.Marshal(val)
	items, err := cli.objcli.SelectPropertyGroupByObjectID(forward, selector.OwnerID, selector.ObjectID, valstr)
	if nil != err {
		blog.Error("can not found the data(%s), error info is %s", string(valstr), err.Error())
		return 0, err
	}
	if len(items) > 0 {
		blog.Error("repeat the group info %+v", selector)
		return 0, fmt.Errorf("repeat the group info")
	}
	return cli.objcli.CreateMetaObjectAttGroup(forward, data)
}

func (cli *objAttGroupLogic) UpdateObjectGroup(forward *api.ForwardParam, data []byte, errProxy errors.DefaultCCErrorIf) error {
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.UpdateMetaObjectAttGroup(forward, data)
}

func (cli *objAttGroupLogic) UpdateObjectGroupProperty(forward *api.ForwardParam, data []byte, errProxy errors.DefaultCCErrorIf) error {
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.UpdateMetaObjectAttGroupProperty(forward, data)
}

func (cli *objAttGroupLogic) DeleteObjectGroup(forward *api.ForwardParam, id int, errProxy errors.DefaultCCErrorIf) error {
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.DeleteMetaObjectAttGroup(forward, id, nil)
}

func (cli *objAttGroupLogic) DeleteObjectGroupProperty(forward *api.ForwardParam, ownerID, objectID, propertyID, groupID string, errProxy errors.DefaultCCErrorIf) error {
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.DeleteMetaObjectAttGroupProperty(forward, ownerID, objectID, propertyID, groupID)
}

func (cli *objAttGroupLogic) SelectPropertyGroupByObjectID(forward *api.ForwardParam, ownerID, objectID string, data []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjAttGroupDes, error) {
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.SelectPropertyGroupByObjectID(forward, ownerID, objectID, data)
}
