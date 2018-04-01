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
func (cli *topoMgr) CreateObjectAsst(obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) (int, error) {
	target := cli.logics[ObjectAsst].(ObjectAsstLogic)
	return target.CreateObjectAsst(obj, errProxy)
}

func (cli *topoMgr) SelectObjectAsst(obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) ([]api.ObjAsstDes, error) {
	target := cli.logics[ObjectAsst].(ObjectAsstLogic)
	return target.SelectObjectAsst(obj, errProxy)
}

func (cli *topoMgr) UpdateObjectAsst(selector, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectAsst].(ObjectAsstLogic)
	return target.UpdateObjectAsst(selector, obj, errProxy)
}

func (cli *topoMgr) DeleteObjectAsstByID(id int, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectAsst].(ObjectAsstLogic)
	return target.DeleteObjectAsstByID(id, errProxy)
}

func (cli *topoMgr) DeleteObjectAsst(obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectAsst].(ObjectAsstLogic)
	return target.DeleteObjectAsst(obj, errProxy)
}

// object attribute interface
func (cli *topoMgr) CreateTopoModel(obj api.ObjAttDes, errProxy errors.DefaultCCErrorIf) (int, error) {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.CreateTopoModel(obj, errProxy)
}
func (cli *topoMgr) SelectTopoModel(rstitems []TopoModelRsp, ownerid, objid, clsid, preid, prename string, errProxy errors.DefaultCCErrorIf) ([]TopoModelRsp, error) {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.SelectTopoModel(rstitems, ownerid, objid, clsid, preid, prename, errProxy)
}
func (cli *topoMgr) DeleteTopoModel(ownerid, objid string, assotype int, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.DeleteTopoModel(ownerid, objid, assotype, errProxy)
}
func (cli *topoMgr) CreateObjectAtt(params api.ObjAttDes, errProxy errors.DefaultCCErrorIf) (int, error) {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.CreateObjectAtt(params, errProxy)
}

func (cli *topoMgr) SelectObjectAtt(params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjAttDes, error) {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.SelectObjectAtt(params, errProxy)
}

func (cli *topoMgr) UpdateObjectAtt(id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.UpdateObjectAtt(id, params, errProxy)
}

func (cli *topoMgr) DeleteObjectAtt(id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectAttribute].(ObjectAttLogic)
	return target.DeleteObjectAtt(id, params, errProxy)
}

// object class interface
func (cli *topoMgr) CreateObjectClass(params []byte, errProxy errors.DefaultCCErrorIf) (int, error) {
	target := cli.logics[ObjectClass].(ObjectClassLogic)
	return target.CreateObjectClass(params, errProxy)
}

func (cli *topoMgr) SelectObjectClass(params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjClsDes, error) {
	target := cli.logics[ObjectClass].(ObjectClassLogic)
	return target.SelectObjectClass(params, errProxy)
}

func (cli *topoMgr) SelectObjectClassWithObjects(ownerID string, params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjClsObjectDes, error) {
	target := cli.logics[ObjectClass].(ObjectClassLogic)
	return target.SelectObjectClassWithObjects(ownerID, params, errProxy)
}
func (cli *topoMgr) UpdateObjectClass(id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectClass].(ObjectClassLogic)
	return target.UpdateObjectClass(id, params, errProxy)
}
func (cli *topoMgr) DeleteObjectClass(id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectClass].(ObjectClassLogic)
	return target.DeleteObjectClass(id, params, errProxy)
}

// object attribute group interface
func (cli *topoMgr) CreateObjectGroup(params []byte, errProxy errors.DefaultCCErrorIf) (int, error) {
	target := cli.logics[ObjectGroup].(ObjectAttGroupLogic)
	return target.CreateObjectGroup(params, errProxy)
}
func (cli *topoMgr) UpdateObjectGroup(params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectGroup].(ObjectAttGroupLogic)
	return target.UpdateObjectGroup(params, errProxy)
}
func (cli *topoMgr) UpdateObjectGroupProperty(params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectGroup].(ObjectAttGroupLogic)
	return target.UpdateObjectGroupProperty(params, errProxy)
}
func (cli *topoMgr) DeleteObjectGroup(id int, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectGroup].(ObjectAttGroupLogic)
	return target.DeleteObjectGroup(id, errProxy)
}
func (cli *topoMgr) DeleteObjectGroupProperty(ownerID, objectID, propertyID, groupID string, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[ObjectGroup].(ObjectAttGroupLogic)
	return target.DeleteObjectGroupProperty(ownerID, objectID, propertyID, groupID, errProxy)
}

func (cli *topoMgr) SelectPropertyGroupByObjectID(ownerID, objectID string, data []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjAttGroupDes, error) {
	target := cli.logics[ObjectGroup].(ObjectAttGroupLogic)
	return target.SelectPropertyGroupByObjectID(ownerID, objectID, data, errProxy)
}

// CreateObject create a new object
func (cli *topoMgr) CreateObject(params []byte, errProxy errors.DefaultCCErrorIf) (int, error) {
	target := cli.logics[Object].(ObjectLogic)
	return target.CreateObject(params, errProxy)
}

// SelectObject select objects
func (cli *topoMgr) SelectObject(params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjDes, error) {
	target := cli.logics[Object].(ObjectLogic)
	return target.SelectObject(params, errProxy)
}

// UpdateObject update object info
func (cli *topoMgr) UpdateObject(id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[Object].(ObjectLogic)
	return target.UpdateObject(id, params, errProxy)
}

// DeleteObject delete object info
func (cli *topoMgr) DeleteObject(id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	target := cli.logics[Object].(ObjectLogic)
	return target.DeleteObject(id, params, errProxy)
}
