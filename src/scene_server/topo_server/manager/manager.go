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
	"configcenter/src/common/errors"
	api "configcenter/src/source_controller/api/object"
	"fmt"
)

var mgr = &topoMgr{logics: make(map[string]interface{})}

type topoMgr struct {
	objctr func() string
	logics map[string]interface{}
}

// setLogic register logic object
func (cli *topoMgr) setLogic(logicName string, logic interface{}) error {
	fmt.Println("set logic object:", logicName)
	if _, ok := cli.logics[logicName]; ok {
		return fmt.Errorf("%s repeated", logicName)
	}

	switch t := logic.(type) {
	case SetConfiger:
		t.Set(cli)
	}

	cli.logics[logicName] = logic
	return nil
}

// Get implement Configer interface
func (cli *topoMgr) Get(target interface{}) string {
	/**
	TODO:need to enable , when all server implement the rdiscover interface
	switch target.(type) {
	case ObjectLogic, ObjectAsstLogic, ObjectAttGroupLogic, ObjectAttLogic, ObjectClassLogic:
		address, err := cli.rd.GetObjectCtrlServ()
		if nil != err {
			blog.Error("failed to get object controller address, error info is %s", err.Error())
			return address
		}
	}
	*/
	return cli.objctr() // TODO: need to delete
}

// object asst interface
func (cli *topoMgr) CreateObjectAsst(forward *api.ForwardParam, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) (int, error) {
	target := cli.logics[ObjectAsst].(ObjectAsstLogic)
	return target.CreateObjectAsst(forward, obj, errProxy)
}

func (cli *topoMgr) SelectObjectAsst(forward *api.ForwardParam, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) ([]api.ObjAsstDes, error) {
	target := cli.logics[ObjectAsst].(ObjectAsstLogic)
	return target.SelectObjectAsst(forward, obj, errProxy)
}

func (cli *topoMgr) UpdateObjectAsst(forward *api.ForwardParam, selector, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectAsst].(ObjectAsstLogic)
	return target.UpdateObjectAsst(forward, selector, obj, errProxy)
}

func (cli *topoMgr) DeleteObjectAsstByID(forward *api.ForwardParam, id int, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectAsst].(ObjectAsstLogic)
	return target.DeleteObjectAsstByID(forward, id, errProxy)
}

func (cli *topoMgr) DeleteObjectAsst(forward *api.ForwardParam, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectAsst].(ObjectAsstLogic)
	return target.DeleteObjectAsst(forward, obj, errProxy)
}

// object attribute interface
func (cli *topoMgr) CreateTopoModel(forward *api.ForwardParam, obj api.ObjAttDes, errProxy errors.DefaultCCErrorIf) (int, error) {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.CreateTopoModel(forward, obj, errProxy)
}
func (cli *topoMgr) SelectTopoModel(forward *api.ForwardParam, rstitems []TopoModelRsp, ownerid, objid, clsid, preid, prename string, errProxy errors.DefaultCCErrorIf) ([]TopoModelRsp, error) {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.SelectTopoModel(forward, rstitems, ownerid, objid, clsid, preid, prename, errProxy)
}
func (cli *topoMgr) DeleteTopoModel(forward *api.ForwardParam, ownerid, objid string, assotype int, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.DeleteTopoModel(forward, ownerid, objid, assotype, errProxy)
}
func (cli *topoMgr) CreateObjectAtt(forward *api.ForwardParam, params api.ObjAttDes, errProxy errors.DefaultCCErrorIf) (int, error) {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.CreateObjectAtt(forward, params, errProxy)
}

func (cli *topoMgr) SelectObjectAtt(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjAttDes, error) {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.SelectObjectAtt(forward, params, errProxy)
}

func (cli *topoMgr) UpdateObjectAtt(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.UpdateObjectAtt(forward, id, params, errProxy)
}

func (cli *topoMgr) DeleteObjectAtt(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.DeleteObjectAtt(forward, id, params, errProxy)
}

// object class interface
func (cli *topoMgr) CreateObjectClass(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) (int, error) {
	target := cli.logics[ObjectClass].(ObjectClassLogic)
	return target.CreateObjectClass(forward, params, errProxy)
}

func (cli *topoMgr) SelectObjectClass(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjClsDes, error) {
	target := cli.logics[ObjectClass].(ObjectClassLogic)
	return target.SelectObjectClass(forward, params, errProxy)
}

func (cli *topoMgr) SelectObjectClassWithObjects(forward *api.ForwardParam, ownerID string, params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjClsObjectDes, error) {
	target := cli.logics[ObjectClass].(ObjectClassLogic)
	return target.SelectObjectClassWithObjects(forward, ownerID, params, errProxy)
}
func (cli *topoMgr) UpdateObjectClass(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectClass].(ObjectClassLogic)
	return target.UpdateObjectClass(forward, id, params, errProxy)
}
func (cli *topoMgr) DeleteObjectClass(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectClass].(ObjectClassLogic)
	return target.DeleteObjectClass(forward, id, params, errProxy)
}

// object attribute group interface
func (cli *topoMgr) CreateObjectGroup(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) (int, error) {
	target := cli.logics[ObjectGroup].(ObjectAttGroupLogic)
	return target.CreateObjectGroup(forward, params, errProxy)
}
func (cli *topoMgr) UpdateObjectGroup(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectGroup].(ObjectAttGroupLogic)
	return target.UpdateObjectGroup(forward, params, errProxy)
}
func (cli *topoMgr) UpdateObjectGroupProperty(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectGroup].(ObjectAttGroupLogic)
	return target.UpdateObjectGroupProperty(forward, params, errProxy)
}
func (cli *topoMgr) DeleteObjectGroup(forward *api.ForwardParam, id int, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectGroup].(ObjectAttGroupLogic)
	return target.DeleteObjectGroup(forward, id, errProxy)
}
func (cli *topoMgr) DeleteObjectGroupProperty(forward *api.ForwardParam, ownerID, objectID, propertyID, groupID string, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectGroup].(ObjectAttGroupLogic)
	return target.DeleteObjectGroupProperty(forward, ownerID, objectID, propertyID, groupID, errProxy)
}

func (cli *topoMgr) SelectPropertyGroupByObjectID(forward *api.ForwardParam, ownerID, objectID string, data []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjAttGroupDes, error) {
	target := cli.logics[ObjectGroup].(ObjectAttGroupLogic)
	return target.SelectPropertyGroupByObjectID(forward, ownerID, objectID, data, errProxy)
}

// CreateObject create a new object
func (cli *topoMgr) CreateObject(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) (int, error) {
	target := cli.logics[Object].(ObjectLogic)
	return target.CreateObject(forward, params, errProxy)
}

// SelectObject select objects
func (cli *topoMgr) SelectObject(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjDes, error) {
	target := cli.logics[Object].(ObjectLogic)
	return target.SelectObject(forward, params, errProxy)
}

// UpdateObject update object info
func (cli *topoMgr) UpdateObject(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[Object].(ObjectLogic)
	return target.UpdateObject(forward, id, params, errProxy)
}

// DeleteObject delete object info
func (cli *topoMgr) DeleteObject(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[Object].(ObjectLogic)
	return target.DeleteObject(forward, id, params, errProxy)
}
