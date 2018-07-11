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
	httpcli "configcenter/src/common/http/httpclient"
	sencapi "configcenter/src/scene_server/api"
	"configcenter/src/scene_server/topo_server/topo_service/manager"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"

	simplejson "github.com/bitly/go-simplejson"
)

type objLogic struct {
	objcli *api.Client
	cfg    manager.Configer
	mgr    manager.Manager
}

var _ manager.ObjectLogic = (*objLogic)(nil)

func init() {
	obj := &objLogic{}
	obj.objcli = api.NewClient("")
	manager.SetManager(obj)
	manager.RegisterLogic(manager.Object, obj)
}

// Set implement SetConfiger interface
func (cli *objLogic) Set(cfg manager.Configer) {
	cli.cfg = cfg
}

// SetManager implement the manager's Hooker interface
func (cli *objLogic) SetManager(mgr manager.Manager) error {
	cli.mgr = mgr
	return nil
}

func (cli *objLogic) CreateObject(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) (int, error) {

	// unmarshal json
	var obj sencapi.ObjectDes
	if jserr := json.Unmarshal(params, &obj); nil != jserr {
		blog.Error("unmarshal json failed, data:%s error:%s", string(params), jserr.Error())
		return 0, jserr
	}

	// 基础校验
	if len(obj.OwnerID) == 0 {
		blog.Error("'OwnerID' is not set")
		return 0, fmt.Errorf("'OwnerID' is not set")
	}

	if len(obj.ObjID) == 0 {
		blog.Error("'ObjID' is not set")
		return 0, fmt.Errorf("'ObjID' is not set")
	}

	if len(obj.ObjName) == 0 {
		blog.Error("'ObjName' is not set")
		return 0, fmt.Errorf("'ObjName' is not set")
	}

	if len(obj.ClassificationID) == 0 {
		blog.Error("'ClassificationID' is not set")
		return 0, fmt.Errorf("'ClassificationID' is not set")
	}

	// disable the inner object id
	switch obj.ObjID {
	case common.BKInnerObjIDApp, common.BKInnerObjIDHost, common.BKInnerObjIDModule, common.BKInnerObjIDSet:
		return 0, fmt.Errorf("the built-in model-id[%s], please use a new name", obj.ObjID)
	}

	// check the classification id
	checkClsCond := make(map[string]interface{})
	checkClsCond["bk_classification_id"] = obj.ClassificationID
	checkClsCondVal, _ := json.Marshal(checkClsCond)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if items, err := cli.objcli.SearchMetaObjectCls(forward, checkClsCondVal); nil != err {

	} else if 0 == len(items) {
		blog.Error("bk_classification_id[%s] is invalid", obj.ClassificationID)
		return 0, fmt.Errorf("bk_classification_id[%s] is invalid", obj.ClassificationID)
	}

	// check the object id
	checkObjIDCond := make(map[string]interface{})
	checkObjIDCond[common.BKObjIDField] = obj.ObjID
	checkObjIDCond[common.BKOwnerIDField] = obj.OwnerID
	checkObjIDCondVal, _ := json.Marshal(checkObjIDCond)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if items, err := cli.objcli.SearchMetaObject(forward, checkObjIDCondVal); nil != err {
		blog.Error("select failed, error:%s", err.Error())
		return 0, err
	} else if 0 != len(items) {
		blog.Warn("repeat to create the object, objid: %s", obj.ObjID)
		return 0, fmt.Errorf("repeat to create the object, objid:%s", obj.ObjID)
	}

	// check the object name
	checkObjNameCond := make(map[string]interface{})
	checkObjNameCond["bk_obj_name"] = obj.ObjName
	checkObjNameCond[common.BKOwnerIDField] = obj.OwnerID
	checkObjNameCondVal, _ := json.Marshal(checkObjNameCond)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if items, err := cli.objcli.SearchMetaObject(forward, checkObjNameCondVal); nil != err {
		blog.Error("select failed, error:%s", err.Error())
		return 0, err
	} else if 0 != len(items) {
		blog.Warn("repeat to create the object, objname: %s", obj.ObjName)
		return 0, fmt.Errorf("repeat to create the object, objname:%s", obj.ObjName)
	}

	// creaet default group
	defaultGroup := map[string]interface{}{
		"bk_isdefault":        true,
		common.BKOwnerIDField: obj.OwnerID,
		common.BKObjIDField:   obj.ObjID,
		"bk_group_name":       "Default",
		"bk_group_index":      -1,
		"bk_group_id":         "default",
	}
	defaultGroupStr, _ := json.Marshal(defaultGroup)
	_, defErr := cli.mgr.CreateObjectGroup(forward, defaultGroupStr, errProxy)

	// create default the instname
	var objAtt api.ObjAttDes
	objAtt.OwnerID = obj.OwnerID
	objAtt.ObjectID = obj.ObjID
	objAtt.Editable = true
	objAtt.Creator = "user"
	objAtt.IsOnly = true
	objAtt.IsPre = true
	objAtt.IsRequired = true
	if nil != defErr {
		blog.Error("failed to create default group, error info is %s", defErr.Error())
		objAtt.PropertyGroup = "none"
	} else {
		// only one default for a group
		objAtt.PropertyGroup = "default"
	}

	objAtt.PropertyIndex = -1
	objAtt.PropertyType = common.FieldTypeSingleChar
	switch obj.ObjID {
	case common.BKInnerObjIDApp:
		objAtt.PropertyID = common.BKAppNameField
		objAtt.PropertyName = "业务名"
	case common.BKInnerObjIDModule:
		objAtt.PropertyID = common.BKModuleNameField
		objAtt.PropertyName = "模块名"
	case common.BKInnerObjIDSet:
		objAtt.PropertyID = common.BKSetNameField
		objAtt.PropertyName = "集群名"
	default:
		objAtt.PropertyID = common.BKInstNameField
		objAtt.PropertyName = "实例名"
	}

	// create the default object attribute
	objAttVal, _ := json.Marshal(objAtt)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if _, objAttErr := cli.objcli.CreateMetaObjectAtt(forward, objAttVal); nil != objAttErr {
		blog.Error(" failed to create the default properid , error info is %s", objAttErr.Error())
	}

	// 执行创建动作
	return cli.objcli.CreateMetaObject(forward, params)
}

func (cli *objLogic) SelectObject(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjDes, error) {

	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.SearchMetaObject(forward, params)
}

func (cli *objLogic) UpdateObject(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	if id <= 0 {
		blog.Error("id is invalid, %d", id)
		return fmt.Errorf("id[%d] is invalid", id)
	}

	js, err := simplejson.NewJson(params)
	if nil != err {
		blog.Error("update object failed, error:%s", err.Error())
		return err
	}

	objMap, _ := js.Map()
	if _, ok := objMap[common.BKObjIDField]; ok {
		blog.Error("'%s' is forbidden to be updated", common.BKObjIDField)
		return fmt.Errorf("'%s' is forbidden to be updated", common.BKObjIDField)
	}

	checkObjNameCond := make(map[string]interface{})

	// 检查数据是否存在
	checkIDCond := make(map[string]interface{})
	checkIDCond["id"] = id
	checkIDCondVal, _ := json.Marshal(checkIDCond)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if items, err := cli.objcli.SearchMetaObject(forward, checkIDCondVal); nil != err {
		blog.Error("select failed, error:%s", err.Error())
		return err
	} else if 0 == len(items) {
		blog.Warn("nothing to updated, id is: %d", id)
		return fmt.Errorf("nothing can be updated, please check the condition")
	} else {
		checkObjNameCond[common.BKOwnerIDField] = items[0].OwnerID
	}

	// 如果设置了ObjName 就要判断名字是否冲突
	if name, ok := objMap["bk_obj_name"]; ok && 0 == len(name.(string)) {
		blog.Error("'bk_obj_name' is not set")
		return fmt.Errorf("'bk_obj_name' is not set")
	} else if ok {

		checkObjNameCond["bk_obj_name"] = name
		checkObjNameCondVal, _ := json.Marshal(checkObjNameCond)
		cli.objcli.SetAddress(cli.cfg.Get(cli))
		if items, err := cli.objcli.SearchMetaObject(forward, checkObjNameCondVal); nil != err {
			blog.Error("select failed, error:%s", err.Error())
			return err
		} else if 0 != len(items) {
			for _, tmpitem := range items {
				if tmpitem.ID != id { // 排除自身
					blog.Warn("repeat to create the object, objname: %v", checkObjNameCond["bk_obj_name"])
					return fmt.Errorf("repeat to create the object, objname:%v", checkObjNameCond["bk_obj_name"])
				}
			}
		}
	}
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.UpdateMetaObject(forward, id, params)
}

func (cli *objLogic) DeleteObject(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error {
	if id < 0 {
		blog.Error("attrid is invalid, %d", id)
		return fmt.Errorf("attrid is invalid, %d", id)
	}
	if 0 != id {
		checkCond := make(map[string]interface{})
		checkCond["id"] = id
		params, _ = json.Marshal(checkCond)
	}

	if 0 == id && 0 == len(params) {
		blog.Error("there are no delete conditions available")
		return fmt.Errorf("there are no delete conditions available")
	}
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	items, err := cli.objcli.SearchMetaObject(forward, params)
	if nil != err {
		blog.Error("check the target failed, error:%s ", err.Error())
		return err
	} else if 0 == len(items) {
		blog.Warn("nothing to be deleted")
		return fmt.Errorf("nothing can be deleted, please check the condition")
	} else {
		// check the inner object
		for _, tmpitem := range items {
			if tmpitem.IsPre {
				blog.Warn("the built-in model[%s:%s] is forbidden to delete", tmpitem.ObjectID, tmpitem.ObjectName)
				return fmt.Errorf("the built-in model[%s:%s] is forbidden to delete", tmpitem.ObjectID, tmpitem.ObjectName)
			}
			/*
				// 检测ObjID是否存在实例，如果存在也不允许删除
				if ok, okerr := cli.existInst(tmpitem.ObjID, tmpitem.OwnerID); nil != okerr {
					blog.Error("the model instance detection failed, error:%s", okerr.Error())
					return fmt.Errorf("the model instance detection failed, error:%s", okerr.Error())
				} else if ok {
					return fmt.Errorf("it is forbidden to delete a model[%s] that has been instantiated", tmpitem.ObjID)
				}
			*/
		}
	}
	blog.Debug("the target object items:%+v", items)
	// clean the object attribute and association
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	for _, tmpItem := range items {

		// delete association
		att := map[string]interface{}{}
		att[common.BKOwnerIDField] = tmpItem.OwnerID
		att[common.BKObjIDField] = tmpItem.ObjectID
		attData, _ := json.Marshal(att)
		if err := cli.objcli.DeleteMetaObjectAtt(forward, 0, attData); nil != err {
			blog.Error("failed to delete the object att, error is %s", err.Error())
			return err
		}

		// delete the association
		asst := map[string]interface{}{}
		asst[common.BKOwnerIDField] = tmpItem.OwnerID
		asst[common.BKObjIDField] = tmpItem.ObjectID
		asstData, _ := json.Marshal(asst)
		if err := cli.objcli.DeleteMetaObjectAsst(forward, 0, asstData); nil != err {
			blog.Error("failed to delete the object asst")
			return err
		}

		// delete group
		grp, grpErr := cli.mgr.SelectPropertyGroupByObjectID(forward, tmpItem.OwnerID, tmpItem.ObjectID, []byte("{}"), errProxy)
		if nil != grpErr {
			blog.Error("failed to search group with object, error info is %s", grpErr.Error())
			return grpErr
		}
		for _, grpItem := range grp {
			if delErr := cli.mgr.DeleteObjectGroup(forward, grpItem.ID, errProxy); nil != delErr {
				blog.Error("failed to delete the group[%d], error info is %s", grpItem.ID, delErr.Error())
				return delErr
			}
		}
	}

	return cli.objcli.DeleteMetaObject(forward, id, params)
}

func (cli *objLogic) setAddress(forward *api.ForwardParam, address string) {
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	cli.objcli.SetAddress(address)
}

func (cli *objLogic) existInst(forward *api.ForwardParam, objID, ownerID string, errProxy errors.DefaultCCErrorIf) (bool, error) {

	cli.objcli.SetAddress(cli.cfg.Get(cli))
	searchCond := make(map[string]interface{})

	searchCond["sort"] = common.BKAppIDField
	searchCond["start"] = 0
	searchCond["limit"] = 100
	searchCond["fields"] = common.BKObjIDField

	cond := make(map[string]interface{})
	cond[common.BKObjIDField] = objID
	cond[common.BKOwnerIDField] = ownerID
	searchCond["condition"] = cond

	val, _ := json.Marshal(searchCond)

	cURL := cli.objcli.GetAddress() + "/object/v1/insts/object/search"

	httpClient := httpcli.NewHttpClient()

	httpClient.SetHeader("Content-Type", "application/json")
	httpClient.SetHeader("Accept", "application/json")

	dres, err := httpClient.POST(cURL, forward.Header, val)
	if nil != err {
		return false, fmt.Errorf("it fails to check the inst, error:%s", err.Error())
	}

	js, err := simplejson.NewJson([]byte(dres))
	sData, _ := js.Map()
	data, _ := sData["data"].(map[string]interface{})
	cnt, _ := data["count"].(json.Number).Int64()
	if cnt > 0 {
		return true, nil
	}

	return false, nil
}
