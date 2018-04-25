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
	"configcenter/src/scene_server/topo_server/topo_service/manager"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
)

type objAttLogic struct {
	objcli *api.Client
	cfg    manager.Configer
	mgr    manager.Manager
}

var _ manager.ObjectAttLogic = (*objAttLogic)(nil)

func init() {
	obj := &objAttLogic{}
	obj.objcli = api.NewClient("")
	manager.SetManager(obj)
	manager.RegisterLogic(manager.ObjectAttribute, obj)
}

// Set implement SetConfiger interface
func (cli *objAttLogic) Set(cfg manager.Configer) {
	cli.cfg = cfg
}

// SetManager implement the manager's Hooker interface
func (cli *objAttLogic) SetManager(mgr manager.Manager) error {
	cli.mgr = mgr
	return nil
}

// CreateModel create main line topo object
func (cli *objAttLogic) CreateTopoModel(forward *api.ForwardParam, obj api.ObjAttDes, errProxy errors.DefaultCCErrorIf) (int, error) {

	blog.Info("create %v", obj)

	// base check
	if 0 == len(obj.OwnerID) {
		blog.Error("'OwnerID' is not set")
		return 0, fmt.Errorf("'OwnerID'is not set")
	}

	if 0 == len(obj.ObjectID) {
		blog.Error("'ObjID' is not set")
		return 0, fmt.Errorf("'ObjID' is not set")
	}

	if obj.ObjectID == obj.AssociationID {
		return 0, fmt.Errorf("no self-correlation")
	}

	// check objid
	objtmp := map[string]interface{}{}
	objtmp[common.BKOwnerIDField] = obj.OwnerID
	objtmp[common.BKObjIDField] = obj.ObjectID
	objTmpJSON, _ := json.Marshal(objtmp)
	if items, err := cli.mgr.SelectObject(forward, objTmpJSON, errProxy); nil != err {
		blog.Error("the existence test failed, error:%s", err.Error())
		return 0, err
	} else if 0 == len(items) {
		blog.Error("the ObjID[%s] is invalid", obj.ObjectID)
		return 0, fmt.Errorf("the ObjID[%s] is invalid", obj.ObjectID)
	}

	// check the association map
	objtmp[common.BKOwnerIDField] = obj.OwnerID
	objtmp[common.BKObjIDField] = obj.AssociationID
	objTmpJSON, _ = json.Marshal(objtmp)
	if items, err := cli.mgr.SelectObject(forward, objTmpJSON, errProxy); nil != err {
		blog.Error("the existence test failed, error:%s", err.Error())
		return 0, err
	} else if 0 == len(items) {
		blog.Error("the AssociationID[%s] is invalid", obj.ObjectID)
		return 0, fmt.Errorf("the AssociationID[%s] is invalid", obj.ObjectID)
	}

	// create object association
	objasst := map[string]interface{}{}
	objasst[common.BKOwnerIDField] = obj.OwnerID
	objasst[common.BKObjIDField] = obj.ObjectID
	objasst["bk_asst_obj_id"] = obj.AssociationID
	obj.Editable = true
	obj.PropertyID = common.BKChildStr
	obj.PropertyType = common.FiledTypeSingleChar
	objasst["bk_object_att_id"] = obj.PropertyID

	// to create object association	, failed return
	id, operr := cli.mgr.CreateObjectAsst(forward, objasst, errProxy)
	// check objatt data
	if nil != operr {
		return 0, operr
	}

	// to create the inner  attribute
	objAtt := api.ObjAttDes{}
	objAtt.ObjectID = obj.ObjectID
	objAtt.OwnerID = obj.OwnerID
	objAtt.PropertyID = common.BKInstParentStr
	objAtt.PropertyType = common.FiledTypeInt
	objAtt.IsSystem = true
	objAtt.IsOnly = true
	objAtt.IsRequired = true

	val, jsErr := json.Marshal(objAtt)
	if nil != jsErr {
		blog.Error("marshal failed, error:%v", jsErr)
		return 0, jsErr
	}

	cli.objcli.SetAddress(cli.cfg.Get(cli))
	innerAttID, rstErr := cli.objcli.CreateMetaObjectAtt(forward, val)
	if nil != rstErr {
		blog.Error("failed to create the inner filed for the owner(%s) object(%s)", obj.OwnerID, obj.ObjectID)
		cli.mgr.DeleteObjectAsstByID(forward, id, errProxy)
		return innerAttID, rstErr
	}

	// to create the main line attribute, if it fails, you need to delete the association
	val, jsErr = json.Marshal(obj)
	if nil != jsErr {
		blog.Error("marshal failed, error:%v", jsErr)
		return 0, jsErr
	}

	cli.objcli.SetAddress(cli.cfg.Get(cli))
	rst, rstErr := cli.objcli.CreateMetaObjectAtt(forward, val)
	if nil != rstErr {
		cli.mgr.DeleteObjectAsstByID(forward, id, errProxy)
		cli.objcli.DeleteMetaObjectAtt(forward, innerAttID, []byte("{}"))
		return rst, rstErr
	}

	return rst, rstErr
}

func (cli *objAttLogic) CreateObjectAtt(forward *api.ForwardParam, obj api.ObjAttDes, errProxy errors.DefaultCCErrorIf) (int, error) {

	// base check
	if 0 == len(obj.OwnerID) {
		blog.Error("'" + common.BKOwnerIDField + "' is not set")
		return 0, fmt.Errorf("'" + common.BKOwnerIDField + "' is not set")
	}

	if 0 == len(obj.ObjectID) {
		blog.Error("'" + common.BKObjIDField + "' is not set")
		return 0, fmt.Errorf("'" + common.BKObjIDField + "' is not set")
	}

	if 0 == len(obj.PropertyID) {
		blog.Error("'bk_property_id' is not set")
		return 0, fmt.Errorf("'bk_property_id' is not set")
	}

	if 0 == len(obj.PropertyName) {
		if obj.PropertyID != common.BKChildStr && obj.PropertyID != common.BKParentStr {
			blog.Error("'bk_property_name' is not set")
			return 0, fmt.Errorf("'bk_property_name' is not set")
		}
	}

	if 0 == len(obj.PropertyType) {
		if obj.PropertyID != common.BKChildStr && obj.PropertyID != common.BKParentStr {
			blog.Error("'bk_property_type' is not set")
			return 0, fmt.Errorf("'bk_property_type' is not set")
		}
	}

	if obj.PropertyID == common.BKChildStr || obj.PropertyID == common.BKParentStr {
		blog.Error("'%s' is the built-in property", obj.PropertyID)
		return 0, fmt.Errorf("'%s' is the built-in property", obj.PropertyID)
	}

	// check objid
	checkObjCond := make(map[string]interface{})
	checkObjCond[common.BKOwnerIDField] = obj.OwnerID
	checkObjCond[common.BKObjIDField] = obj.ObjectID
	checkObjCondVal, _ := json.Marshal(checkObjCond)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if items, err := cli.objcli.SearchMetaObject(forward, checkObjCondVal); nil != err {
		blog.Error("the existence test failed, error:%s", err.Error())
		return 0, err
	} else if 0 == len(items) {
		blog.Error("bk_obj_id[%s] is invalid", obj.ObjectID)
		return 0, fmt.Errorf("bk_obj_id[%s] is invalid", obj.ObjectID)
	}

	// check property name
	checkAttNameCond := make(map[string]interface{})
	checkAttNameCond[common.BKOwnerIDField] = obj.OwnerID
	checkAttNameCond[common.BKObjIDField] = obj.ObjectID
	checkAttNameCond["bk_property_name"] = obj.PropertyName
	checkAttNameCondVal, _ := json.Marshal(checkAttNameCond)

	if items, itemErr := cli.objcli.SearchMetaObjectAtt(forward, checkAttNameCondVal); nil != itemErr {
		blog.Error("create objectt failed, error:%s", itemErr.Error())
		return 0, itemErr
	} else if 0 != len(items) {
		blog.Warn("duplicate property name, PropertyName: %s", obj.PropertyName)
		return 0, fmt.Errorf("duplicate property name, PropertyName: %s", obj.PropertyName)
	}

	// check property id
	checkAttIDCond := make(map[string]interface{})
	checkAttIDCond[common.BKOwnerIDField] = obj.OwnerID
	checkAttIDCond[common.BKObjIDField] = obj.ObjectID
	checkAttIDCond["bk_property_id"] = obj.PropertyID
	checkAttIDCondVal, _ := json.Marshal(checkAttIDCond)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if items, itemErr := cli.objcli.SearchMetaObjectAtt(forward, checkAttIDCondVal); nil != itemErr {
		blog.Error("create objectt failed, error:%s", itemErr.Error())
		return 0, itemErr
	} else if 0 != len(items) {
		blog.Warn("duplicate propertyid, PropertyName: %s", obj.PropertyID)
		return 0, fmt.Errorf("duplicate propertyid, PropertyID: %s", obj.PropertyID)
	}

	// to create an object attribute, if it fails, you need to delete the association
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if 0 == len(obj.PropertyGroup) {
		// search property group
		groupCondition := map[string]bool{
			"bk_isdefault": true,
		}
		groupConditionStr, _ := json.Marshal(groupCondition)
		groupDes, groupDesErr := cli.objcli.SelectPropertyGroupByObjectID(forward, obj.OwnerID, obj.ObjectID, groupConditionStr)
		if nil != groupDesErr {
			blog.Error("failed to found the group config, error info is %s", groupDesErr.Error())
		}
		if 0 != len(groupDes) {
			obj.PropertyGroup = groupDes[0].GroupID // default is only one

		}
	}

	// create object association
	objasst := map[string]interface{}{}
	objasst[common.BKOwnerIDField] = obj.OwnerID
	objasst[common.BKObjIDField] = obj.ObjectID
	objasst["bk_asst_obj_id"] = obj.AssociationID
	objasst["bk_object_att_id"] = obj.PropertyID
	objasst["bk_asst_forward"] = obj.AsstForward

	// to create object association	, failed return
	var asstID int
	if 0 != len(obj.AssociationID) {

		// 检测AssociationID是否为已存在的模型
		checkObjCond[common.BKOwnerIDField] = obj.OwnerID
		checkObjCond[common.BKObjIDField] = obj.AssociationID
		checkObjCondVal, _ := json.Marshal(checkObjCond)
		cli.objcli.SetAddress(cli.cfg.Get(cli))
		if items, err := cli.objcli.SearchMetaObject(forward, checkObjCondVal); nil != err {
			blog.Error("the existence test failed, error:%s", err.Error())
			return 0, err
		} else if 0 == len(items) {
			blog.Error("AssociationID[%s] is invalid", obj.AssociationID)
			return 0, fmt.Errorf("AssociationID[%s] is invalid", obj.AssociationID)
		}

		id, opErr := cli.mgr.CreateObjectAsst(forward, objasst, errProxy)
		// check objatt data
		if nil != opErr {
			return 0, opErr
		}

		asstID = id
	}

	val, _ := json.Marshal(obj)
	rst, rstErr := cli.objcli.CreateMetaObjectAtt(forward, val)
	if nil != rstErr && 0 != asstID {
		cli.mgr.DeleteObjectAsstByID(forward, asstID, errProxy)
		return rst, rstErr
	}

	return rst, rstErr
}

func (cli *objAttLogic) UpdateObjectAtt(forward *api.ForwardParam, attrID int, val []byte, errProxy errors.DefaultCCErrorIf) error {

	if attrID <= 0 {
		blog.Error("ID is invalid, %d", attrID)
		return fmt.Errorf("ID is invalid, %d", attrID)
	}

	// check whether it is exists
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	itemObj, itemErr := cli.objcli.SearchMetaObjectAttByID(forward, attrID)
	if nil != itemErr {
		blog.Error("update objectt failed, error:%s", itemErr.Error())
		return itemErr
	} else if nil == itemObj {
		blog.Warn("nothing can be updated, ID: %d", attrID)
		return fmt.Errorf("nothing can be updated, please check the condition")
	}

	// base check
	// var obj api.ObjAttDes
	var obj map[string]interface{}

	if jsErr := json.Unmarshal(val, &obj); nil != jsErr {
		blog.Error("unmarshal json failed, error information is %v", jsErr)
		return jsErr
	}
	blog.Debug("object att data:%s", string(val))

	// objAtt will be used to saving data whitout ignore fields,
	// bk_property_id field is not allow to edit so we will not copy id value
	objAtt := make(map[string]interface{})

	// check whether bk_property_name duplicated
	if propertyName, ok := obj["bk_property_name"]; ok && "" != propertyName {

		checkAttNameCond := make(map[string]interface{})
		checkAttNameCond[common.BKOwnerIDField] = itemObj.OwnerID
		checkAttNameCond[common.BKObjIDField] = itemObj.ObjectID
		checkAttNameCond["bk_property_name"] = propertyName
		checkAttNameCondVal, _ := json.Marshal(checkAttNameCond)
		cli.objcli.SetAddress(cli.cfg.Get(cli))
		if items, itemErr := cli.objcli.SearchMetaObjectAtt(forward, checkAttNameCondVal); nil != itemErr {
			blog.Error("update objectt failed, error:%s", itemErr.Error())
			return itemErr
		} else if 0 != len(items) {
			for _, tmpItem := range items {
				if tmpItem.ID != attrID { // except self
					blog.Warn("duplicate property name, PropertyName: %s", propertyName)
					return fmt.Errorf("duplicate property name, PropertyName: %s", propertyName)
				}
			}
		}

		objAtt["bk_property_name"] = propertyName
	}

	// update object association
	if AssociationID, ok := obj["bk_asst_obj_id"]; ok && "" != AssociationID {

		// base check
		PropertyID, ok := obj["bk_property_id"]
		if !ok && "" != PropertyID {
			blog.Error("'bk_property_id' is not set")
			return fmt.Errorf("'bk_property_id' is not set")
		}

		ObjectID, ok := obj["bk_obj_id"]
		if !ok && "" != ObjectID {
			blog.Error("'bk_obj_id' is not set")
			return fmt.Errorf("'bk_obj_id' is not set")
		}

		OwnerID, ok := obj["bk_supplier_account"]
		if !ok && "" != OwnerID {
			blog.Error("'bk_supplier_account' is not set")
			return fmt.Errorf("'bk_supplier_account' is not set")
		}

		if PropertyID == common.BKChildStr || PropertyID == common.BKParentStr {
			blog.Error("'%s' is the built-in property", PropertyID)
			return fmt.Errorf("'%s' is the built-in property, it can't be modified", PropertyID)
		}

		// delete the association map
		delAsst := make(map[string]interface{})

		delAsst[common.BKObjIDField] = ObjectID
		delAsst[common.BKOwnerIDField] = OwnerID
		delAsst["bk_object_att_id"] = PropertyID
		//delasst["AsstObjID"] = obj.AssociationID
		delAsstVal, _ := json.Marshal(delAsst)
		cli.objcli.SetAddress(cli.cfg.Get(cli))
		if delErr := cli.objcli.DeleteMetaObjectAsst(forward, 0, delAsstVal); nil != delErr {
			blog.Error("delete association failed, error:%s", delErr.Error())
			return delErr
		}

		// recreate association map
		newAsst := api.ObjAsstDes{}
		newAsst.ObjectID = fmt.Sprint(ObjectID)
		newAsst.OwnerID = fmt.Sprint(OwnerID)
		newAsst.AsstObjID = fmt.Sprint(AssociationID)
		newAsst.ObjectAttID = fmt.Sprint(PropertyID)
		newAsst.AsstForward = fmt.Sprint(obj["bk_asst_forward"])

		newAsstVal, _ := json.Marshal(newAsst)
		cli.objcli.SetAddress(cli.cfg.Get(cli))
		if _, crtErr := cli.objcli.CreateMetaObjectAsst(forward, newAsstVal); nil != crtErr {

			blog.Error("create association failed, error:%s", crtErr.Error())
			if _, crtErr = cli.objcli.CreateMetaObjectAsst(forward, delAsstVal); nil != crtErr {

				blog.Error("create new association failed, and reset the old association failed, error:%s", crtErr.Error())
				return fmt.Errorf("create new association failed, and reset the old association failed, error:%s", crtErr.Error())
			}
			return crtErr
		}
	}

	// update object attribute
	// bk_supplier_account 和 bk_obj_id are not allowed to edit
	if fieldValue, ok := obj["bk_property_group"]; ok {
		objAtt["bk_property_group"] = fieldValue
	}
	if fieldValue, ok := obj["option"]; ok {
		objAtt["option"] = fieldValue
	}
	if fieldValue, ok := obj["creator"]; ok {
		objAtt["creator"] = fieldValue
	}
	if fieldValue, ok := obj["description"]; ok {
		objAtt["description"] = fieldValue
	}
	if fieldValue, ok := obj["unit"]; ok {
		objAtt["unit"] = fieldValue
	}
	if fieldValue, ok := obj["placeholder"]; ok {
		objAtt["placeholder"] = fieldValue
	}
	if fieldValue, ok := obj["editable"]; ok {
		objAtt["editable"] = fieldValue
	}
	if fieldValue, ok := obj["isrequired"]; ok {
		objAtt["isrequired"] = fieldValue
	}
	if fieldValue, ok := obj["isreadonly"]; ok {
		objAtt["isreadonly"] = fieldValue
	}
	if fieldValue, ok := obj["isonly"]; ok {
		objAtt["isonly"] = fieldValue
	}
	if fieldValue, ok := obj["bk_property_type"]; ok {
		objAtt["bk_property_type"] = fieldValue
	}

	objAttVal, _ := json.Marshal(objAtt)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.UpdateMetaObjectAtt(forward, attrID, objAttVal)
}

// deleteTopoModel delete main line topo object
func (cli *objAttLogic) DeleteTopoModel(forward *api.ForwardParam, ownerID, objID string, assoType int, errProxy errors.DefaultCCErrorIf) error {

	objAsst := map[string]interface{}{}
	obj := map[string]interface{}{}

	switch assoType {
	case common.BKParent:
		obj["bk_property_id"] = common.BKParentStr
	case common.BKChild:
		obj["bk_property_id"] = common.BKChildStr
	default:
		return fmt.Errorf("unknown asso type %d", assoType)
	}

	obj[common.BKOwnerIDField] = ownerID
	obj[common.BKObjIDField] = objID
	//obj.AssoType = assotype
	//obj.Editable = true

	objAsst[common.BKOwnerIDField] = ownerID
	objAsst[common.BKObjIDField] = objID
	objAsst["bk_object_att_id"] = obj["bk_property_id"]

	blog.Debug("delete object att %v", obj)

	// delete object attribute
	val, jsErr := json.Marshal(obj)
	if nil != jsErr {
		return fmt.Errorf("marshal failed, error:%v", jsErr.Error())
	}

	blog.Debug("delete objectatt, %s", string(val))
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	rstErr := cli.objcli.DeleteMetaObjectAtt(forward, 0, val)
	if nil != rstErr {
		blog.Error("delete object att failed, error informationis %v", rstErr)
		return rstErr
	}

	// delete asst
	blog.Debug("delete association %v", objAsst)
	if delErr := cli.mgr.DeleteObjectAsst(forward, objAsst, errProxy); nil != delErr {
		blog.Error("delete association failed, error information is %v", delErr)
		return delErr
	}

	return rstErr
}

func (cli *objAttLogic) DeleteObjectAtt(forward *api.ForwardParam, attrID int, val []byte, errProxy errors.DefaultCCErrorIf) error {

	if 0 > attrID {
		blog.Error("ID is invalid, ID is %d", attrID)
		return fmt.Errorf("ID is invalid")
	}

	if 0 == attrID && 0 == len(val) {
		blog.Error("there are no delete conditions available")
		return fmt.Errorf("there are no delete conditions available")
	}

	objAsst := map[string]interface{}{}
	obj := api.ObjAttDes{}

	if 0 == attrID {

		if nil != val && 0 != len(val) {
			if objErr := json.Unmarshal(val, &obj); nil != objErr {
				blog.Error("failed to unmarshal the json, error info is %s", objErr.Error())
				return objErr
			}

			objAsst[common.BKOwnerIDField] = obj.OwnerID
			objAsst[common.BKObjIDField] = obj.ObjectID
			//objasst.ObjectAttID = obj.PropertyID

			// base check
			if 0 == len(obj.OwnerID) {
				blog.Error("'%s' is not set", common.BKOwnerIDField)
				return fmt.Errorf("'%s' is not set", common.BKOwnerIDField)
			}

			if 0 == len(obj.ObjectID) {
				blog.Error("'%s' is not set", common.BKObjIDField)
				return fmt.Errorf("'%s' is not set", common.BKObjIDField)
			}

			// check the objectid
			checkObjAttCond := make(map[string]interface{})
			checkObjAttCond[common.BKOwnerIDField] = obj.OwnerID
			checkObjAttCond[common.BKObjIDField] = obj.ObjectID

			checkObjAttCondVal, _ := json.Marshal(checkObjAttCond)
			cli.objcli.SetAddress(cli.cfg.Get(cli))
			if items, itemsErr := cli.objcli.SearchMetaObjectAtt(forward, checkObjAttCondVal); nil != itemsErr {
				blog.Error("failed to search meta object attribute, error info is %s", itemsErr.Error())
				return itemsErr
			} else if 0 == len(items) {
				blog.Error("nothing to be delete, condition:%s", string(checkObjAttCondVal))
				// objatt not found
				return fmt.Errorf("nothing to be deleted, please the condition")
			}
		}

	} else {
		// read object attribute by id
		cli.objcli.SetAddress(cli.cfg.Get(cli))
		objAtt, rstErr := cli.objcli.SearchMetaObjectAttByID(forward, attrID)
		if nil != rstErr {
			blog.Error("call subsearch failed for object attribute, objatt id %v", attrID)
			return fmt.Errorf("nothing to be deleted, please check the condition")
		}

		if objAtt.IsPre {
			return fmt.Errorf("could not delete preset attribute")
		}

		objAsst[common.BKObjIDField] = objAtt.ObjectID
		objAsst[common.BKOwnerIDField] = objAtt.OwnerID
		objAsst["bk_object_att_id"] = objAtt.PropertyID

	}

	// save the old association map
	oldAsstItems, tmpErr := cli.mgr.SelectObjectAsst(forward, objAsst, errProxy)
	if nil != tmpErr {
		blog.Error("cache the old data failed, error:%s", tmpErr.Error())
		return tmpErr
	}
	// delete association map
	blog.Debug("delete association %v", objAsst)
	if delErr := cli.mgr.DeleteObjectAsst(forward, objAsst, errProxy); nil != delErr {
		blog.Error("delete association failed, error information is %v", delErr)
		//return delerr
	}

	// delete object attribute
	blog.Debug("delete objectatt %d", attrID)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	tmpErr = cli.objcli.DeleteMetaObjectAtt(forward, attrID, val)
	if nil != tmpErr {
		for _, tmp := range oldAsstItems {
			tmpVal, _ := json.Marshal(tmp)
			cli.objcli.SetAddress(cli.cfg.Get(cli))
			if _, subTmpErr := cli.objcli.CreateMetaObjectAsst(forward, tmpVal); nil != subTmpErr {
				blog.Error("delete objectatt failed, reset the associtation failed, error:%s", subTmpErr.Error())
			}
		}
		return fmt.Errorf("delete the data failed, error:%s", tmpErr.Error())
	}
	return nil
}

// SelectTopoModel 根据模型ID查询拓扑，仅向下查询
func (cli *objAttLogic) SelectTopoModel(forward *api.ForwardParam, rstItems []manager.TopoModelRsp, ownerID, objID, clsID, preID, preName string, errProxy errors.DefaultCCErrorIf) ([]manager.TopoModelRsp, error) {

	blog.Info("ownerid %s objid %s", ownerID, objID)
	// read parent object
	parentObj := map[string]interface{}{}

	parentObj[common.BKOwnerIDField] = ownerID
	parentObj[common.BKObjIDField] = objID

	if 0 != len(clsID) {
		parentObj["bk_classification_id"] = clsID
	}

	parentObjJSON, _ := json.Marshal(parentObj)
	blog.Debug("json:%v", string(parentObjJSON))
	parentObjMsg, parentObjErr := cli.mgr.SelectObject(forward, parentObjJSON, errProxy)
	if nil != parentObjErr {
		blog.Error("search parent object failed, error:%v", parentObjErr)
		return nil, fmt.Errorf("%v", parentObjErr)
	}

	if 0 >= len(parentObjMsg) {

		blog.Warn("can not found the object[%s:%s]", ownerID, objID)
		if nil != rstItems {
			return rstItems, nil
		}

		return nil, fmt.Errorf("there is not any object for  ownerid %s objid %s", ownerID, objID)
	}

	// construct result
	if nil == rstItems {
		rstItems = make([]manager.TopoModelRsp, 0)
	}

	topoObj := manager.TopoModelRsp{}

	// should only one parent
	parent := parentObjMsg[0]

	topoObj.OwnerID = ownerID
	topoObj.ObjID = parent.ObjectID
	topoObj.ObjName = parent.ObjectName
	topoObj.PreObjID = preID
	topoObj.PreObjName = preName

	// search all results of the current object
	selector := map[string]interface{}{}

	selector[common.BKOwnerIDField] = ownerID
	selector["bk_asst_obj_id"] = objID
	selector["bk_object_att_id"] = common.BKChildStr

	// 读取关联关系
	asstMsg, opErr := cli.mgr.SelectObjectAsst(forward, selector, errProxy)
	if nil != opErr {
		blog.Error("search association failed, error:%v", opErr)
		return nil, opErr
	}

	if 0 >= len(asstMsg) {
		rstItems = append(rstItems, topoObj)
		return rstItems, nil
	}

	itemAsso := asstMsg[0]

	// 读取子模型信息

	childObj := map[string]interface{}{}
	childObj[common.BKOwnerIDField] = itemAsso.OwnerID
	childObj[common.BKObjIDField] = itemAsso.ObjectID

	if 0 != len(clsID) {
		childObj["bk_classification_id"] = clsID
	}

	childObjJSON, _ := json.Marshal(childObj)
	childObjMsg, childObjErr := cli.mgr.SelectObject(forward, childObjJSON, errProxy)
	if nil != childObjErr {
		blog.Error("search child object failed, error:%v", childObjErr)
		return nil, childObjErr
	}

	if 0 >= len(childObjMsg) {

		blog.Warn("can not found the object[%s:%s]", itemAsso.OwnerID, itemAsso.ObjectID)
		rstItems = append(rstItems, topoObj)
		return rstItems, nil
	}

	// should be only one child
	child := childObjMsg[0]

	topoObj.NextObj = child.ObjectID
	topoObj.NextName = child.ObjectName

	rstItems = append(rstItems, topoObj)

	// recursion search
	return cli.SelectTopoModel(forward, rstItems, itemAsso.OwnerID, itemAsso.ObjectID, clsID, topoObj.ObjID, topoObj.ObjName, errProxy)
}

func (cli *objAttLogic) SelectObjectAtt(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjAttDes, error) {
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.searchWithParams(forward, params, errProxy)
}

func (cli *objAttLogic) subSearchWithParams(forward *api.ForwardParam, val []byte) ([]api.ObjAttDes, error) {

	cli.objcli.SetAddress(cli.cfg.Get(cli))
	objs, err := cli.objcli.SearchMetaObjectAttExceptInnerFiled(forward, val)

	// TODO need to delete

	delArrayFunc := func(s []api.ObjAttDes, i int) []api.ObjAttDes {
		return append(s[:i], s[i+1:]...)
	}

retry:
	for idx, tmp := range objs {
		if tmp.PropertyID == common.BKChildStr || tmp.PropertyID == common.BKParentStr || tmp.IsSystem {
			// 清理当前的值
			objs = delArrayFunc(objs, idx)
			goto retry
		}
	}

	return objs, err
}

func (cli *objAttLogic) searchWithParams(forward *api.ForwardParam, val []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjAttDes, error) {

	// read object attributes
	objAtts, rstErr := cli.subSearchWithParams(forward, val)
	if nil != rstErr {
		blog.Error("call subsearch failed for object attribute")
		return nil, rstErr
	}

	for idx, tmp := range objAtts {

		// read association map attribute
		objAsst := map[string]interface{}{}
		objAsst[common.BKObjIDField] = tmp.ObjectID
		objAsst[common.BKOwnerIDField] = tmp.OwnerID
		objAsst["bk_object_att_id"] = tmp.PropertyID

		asstMsg, asstErr := cli.mgr.SelectObjectAsst(forward, objAsst, errProxy)
		// read property group
		condition := map[string]interface{}{
			"bk_group_id":         tmp.PropertyGroup,
			common.BKOwnerIDField: tmp.OwnerID,
			"bk_obj_id":           tmp.ObjectID,
		}
		conditionStr, _ := json.Marshal(condition)
		groups, groupErr := cli.mgr.SelectPropertyGroupByObjectID(forward, tmp.OwnerID, tmp.ObjectID, conditionStr, errProxy)
		if nil != groupErr {
			blog.Error("failed to search the property group, error info is %s", groupErr.Error())
		} else if 0 != len(groups) {
			objAtts[idx].PropertyGroupName = groups[0].GroupName
		}

		if nil != asstErr {
			blog.Error("search obj association failed, error information is %v", asstErr)
			return nil, fmt.Errorf("search obj association failed, error information is %v", asstErr)
		}

		if 0 < len(asstMsg) {
			objAtts[idx].AssociationID = asstMsg[0].AsstObjID // by the rules, only one id
			objAtts[idx].AsstForward = asstMsg[0].AsstForward // by the rules, only one id
		}

		if 0 == len(objAtts[idx].PropertyGroup) {
			objAtts[idx].PropertyGroup = "none"
		}

	}

	// return
	return objAtts, nil
}
