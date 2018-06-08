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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/scene_server/topo_server/manager"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
)

type objAssoLogic struct {
	objcli *api.Client
	cfg    manager.Configer
	mgr    manager.Manager
}

var _ manager.ObjectAsstLogic = (*objAssoLogic)(nil) // check the interface

func init() {
	obj := &objAssoLogic{}

	obj.objcli = api.NewClient("")
	manager.SetManager(obj)
	manager.RegisterLogic(manager.ObjectAsst, obj)
}

// Set implement SetConfiger interface
func (cli *objAssoLogic) Set(cfg manager.Configer) {
	cli.cfg = cfg
}

// SetManager implement the manager's Hooker interface
func (cli *objAssoLogic) SetManager(mgr manager.Manager) error {
	cli.mgr = mgr
	return nil
}

func (cli *objAssoLogic) CreateObjectAsst(forward *api.ForwardParam, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) (int, error) {

	objasstval, jserr := json.Marshal(obj)
	if nil != jserr {
		blog.Error("mashar json failed, error information is %v", jserr)
		return 0, jserr
	}

	// 基础校验
	if _, ok := obj[common.BKOwnerIDField]; !ok {
		blog.Error("'%s' is not set", common.BKOwnerIDField)
		return 0, fmt.Errorf("'%s' is not set", common.BKOwnerIDField)
	}

	if _, ok := obj[common.BKObjIDField]; !ok {
		blog.Error("'%s' is not set", common.BKObjIDField)
		return 0, fmt.Errorf("'%s' is not set", common.BKObjIDField)
	}

	if _, ok := obj["bk_asst_obj_id"]; !ok {
		blog.Error("'bk_asst_obj_id' is not set")
		return 0, fmt.Errorf("'bk_asst_obj_id' is not set")
	}

	if _, ok := obj["bk_object_att_id"]; !ok {
		blog.Error("'bk_object_att_id' is not set")
		return 0, fmt.Errorf("'bk_object_att_id' is not set")
	}

	// 关联关系存在性校验
	checkCond := make(map[string]interface{})

	checkCond[common.BKObjIDField] = obj[common.BKObjIDField]
	checkCond[common.BKOwnerIDField] = obj[common.BKOwnerIDField]
	checkCond["bk_asst_obj_id"] = obj["bk_asst_obj_id"]

	checkCondVal, _ := json.Marshal(checkCond)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if items, err := cli.objcli.SearchMetaObjectAsst(forward, checkCondVal); nil != err {
		blog.Error("select failed, error:%s", err.Error())
		return 0, err
	} else if len(items) > 0 {
		blog.Warn("repeat to create the object association, objid: %s", string(checkCondVal))
		return 0, fmt.Errorf("repeat associated, objid:%v with asstobjid:%v", obj[common.BKObjIDField], obj["bk_asst_obj_id"])
	}

	return cli.objcli.CreateMetaObjectAsst(forward, objasstval)
}

func (cli *objAssoLogic) SelectObjectAsst(forward *api.ForwardParam, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) ([]api.ObjAsstDes, error) {

	objasstval, jserr := json.Marshal(obj)
	if nil != jserr {
		blog.Error("mashar json failed, error information is %v", jserr)
		return nil, jserr
	}

	blog.Info("search objassociation %v", string(objasstval))
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.SearchMetaObjectAsst(forward, objasstval)
}

func (cli *objAssoLogic) UpdateObjectAsst(forward *api.ForwardParam, selector, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) error {

	objasstval, jserr := json.Marshal(obj)
	if nil != jserr {
		blog.Error("mashar json failed, error information is %v", jserr)
		return jserr
	}

	rstmsg, operr := cli.SelectObjectAsst(forward, selector, errProxy)
	if nil != operr {
		blog.Error("search object association failed")
		return operr
	}

	if 0 == len(rstmsg) {
		if _, err := cli.CreateObjectAsst(forward, obj, errProxy); nil != err {
			blog.Error("update object association failed, error information is %v", err)
			return err
		}
		return nil
	}

	blog.Info("update objassociation %v", string(objasstval))
	for _, tmp := range rstmsg {
		cli.objcli.SetAddress(cli.cfg.Get(cli))
		if rsterr := cli.objcli.UpdateMetaObjectAsst(forward, tmp.ID, objasstval); nil != rsterr {
			blog.Error("http put request failed")
			return rsterr
		}

	}

	return nil
}

func (cli *objAssoLogic) DeleteObjectAsstByID(forward *api.ForwardParam, id int, errProxy errors.DefaultCCErrorIf) error {

	// 关联关系存在性校验
	checkCond := make(map[string]interface{})

	checkCond["id"] = id
	checkCondVal, _ := json.Marshal(checkCond)

	if items, err := cli.objcli.SearchMetaObjectAsst(forward, checkCondVal); nil != err {
		blog.Error("select failed, error:%s", err.Error())
		return err
	} else if 0 == len(items) {
		blog.Warn("nothing can be deleted, id is %d", id)
		return fmt.Errorf("nothing can be deleted, please check the condition")
	}
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.DeleteMetaObjectAsst(forward, id, nil)
}

func (cli *objAssoLogic) DeleteObjectAsst(forward *api.ForwardParam, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) error {

	objasstval, jserr := json.Marshal(obj)
	if nil != jserr {
		blog.Error("mashar json failed, error information is %v", jserr)
		return jserr
	}
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if items, err := cli.objcli.SearchMetaObjectAsst(forward, objasstval); nil != err {
		blog.Error("select failed, error:%s", err.Error())
		return err
	} else if 0 == len(items) {
		blog.Warn("nothing can be deleted, condition is %s", objasstval)
		return fmt.Errorf("nothing can be deleted, please check the condition")
	}
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.DeleteMetaObjectAsst(forward, 0, objasstval)
}
