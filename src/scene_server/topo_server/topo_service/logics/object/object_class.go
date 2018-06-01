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
	sencapi "configcenter/src/scene_server/api"
	"configcenter/src/scene_server/topo_server/topo_service/manager"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"

	simplejson "github.com/bitly/go-simplejson"
)

type objClsLogic struct {
	objcli *api.Client
	cfg    manager.Configer
	mgr    manager.Manager
}

var _ manager.ObjectClassLogic = (*objClsLogic)(nil)

func init() {
	obj := &objClsLogic{}
	obj.objcli = api.NewClient("")
	manager.SetManager(obj)
	manager.RegisterLogic(manager.ObjectClass, obj)
}

// Set implement SetConfiger interface
func (cli *objClsLogic) Set(cfg manager.Configer) {
	cli.cfg = cfg
}

// SetManager implement the manager's Hooker interface
func (cli *objClsLogic) SetManager(mgr manager.Manager) error {
	cli.mgr = mgr
	return nil
}

func (cli *objClsLogic) CreateObjectClass(forward *api.ForwardParam, val []byte, errProxy errors.DefaultCCErrorIf) (int, error) {

	var obj sencapi.ObjectClsDes
	if jsErr := json.Unmarshal(val, &obj); nil != jsErr {
		blog.Error("unmarshal json failed, error:%s", jsErr.Error())
		return 0, jsErr
	}

	// base check
	if 0 == len(obj.ClsID) {
		blog.Error("'bk_classification_id' is not set")
		return 0, fmt.Errorf("'bk_classification_id' is not set ")
	}

	if 0 == len(obj.ClsName) {
		blog.Error("'bk_classification_name' is not set")
		return 0, fmt.Errorf("'bk_classification_name' is not set ")
	}

	// uniqueness check
	checkIDCond := make(map[string]interface{})
	checkIDCond["bk_classification_id"] = obj.ClsID
	checkIDCondVal, _ := json.Marshal(checkIDCond)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if items, err := cli.objcli.SearchMetaObjectCls(forward, checkIDCondVal); nil != err {
		blog.Error("select failed, error:%s", err.Error())
		return 0, err
	} else if len(items) > 0 {
		blog.Warn("repeat to create the objectcls, ClassificationID: %s", obj.ClsID)
		return 0, fmt.Errorf("repeat to create the objectcls, ClassificationID: %s", obj.ClsID)
	}
	// uniqueness check
	checkNameCond := make(map[string]interface{})
	checkNameCond["bk_classification_name"] = obj.ClsName
	checkNameCondVal, _ := json.Marshal(checkNameCond)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if items, err := cli.objcli.SearchMetaObjectCls(forward, checkNameCondVal); nil != err {
		blog.Error("select failed, error:%s", err.Error())
		return 0, err
	} else if len(items) > 0 {
		blog.Warn("repeat to create the objectcls, ClassificationName: %s", obj.ClsName)
		return 0, fmt.Errorf("repeat to create the objectcls, ClassificationName: %s", obj.ClsName)
	}

	blog.Debug("create ", string(val))
	return cli.objcli.CreateMetaObjectCls(forward, val)
}

func (cli *objClsLogic) UpdateObjectClass(forward *api.ForwardParam, clsID int, val []byte, errProxy errors.DefaultCCErrorIf) error {

	if clsID <= 0 {
		blog.Error("attrid is invalid, %d", clsID)
		return fmt.Errorf("ID is invalid, %d", clsID)
	}

	js, err := simplejson.NewJson(val)
	if nil != err {
		blog.Error("update objectcls failed, error:%s", err.Error())
		return err
	}

	objMap, _ := js.Map()
	if _, ok := objMap["bk_classification_id"]; ok {
		blog.Error("'bk_classification_id' is forbidden to be updated")
		return fmt.Errorf("'bk_classification_id' is forbidden to be updated")
	}

	// uniqueness check
	checkIDCond := make(map[string]interface{})
	checkIDCond["id"] = clsID
	checkIDCondVal, _ := json.Marshal(checkIDCond)
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if items, err := cli.objcli.SearchMetaObjectCls(forward, checkIDCondVal); nil != err {
		blog.Error("update failed, error:%s", err.Error())
		return err
	} else if 0 == len(items) {
		blog.Warn("nothing can be updated, id: %d", clsID)
		return fmt.Errorf("nothing can be updated, please check the condition")
	}

	// uniqueness check
	if clsName, clsNameOk := objMap["bk_classification_name"]; clsNameOk && 0 == len(clsName.(string)) {
		blog.Error("'bk_classification_name' is not set")
		return fmt.Errorf("'bk_classification_name' is not set")
	} else if clsNameOk {
		// uniqueness check
		checkNameCond := make(map[string]interface{})
		checkNameCond["bk_classification_name"] = clsName
		checkNameCondVal, _ := json.Marshal(checkNameCond)
		cli.objcli.SetAddress(cli.cfg.Get(cli))
		if items, err := cli.objcli.SearchMetaObjectCls(forward, checkNameCondVal); nil != err {
			blog.Error("select failed, error:%s", err.Error())
			return err
		} else if 0 != len(items) {
			for _, tmpitem := range items {
				if tmpitem.ID != clsID { // 排除自身
					blog.Warn("repeat to create the objectcls, ClassificationName: %s", objMap["bk_classification_name"])
					return fmt.Errorf("repeat to create the objectcls, ClassificationName: %s", objMap["bk_classification_name"])
				}
			}
		}
	}

	// execute
	return cli.objcli.UpdateMetaObjectCls(forward, clsID, val)
}

func (cli *objClsLogic) DeleteObjectClass(forward *api.ForwardParam, clsID int, val []byte, errProxy errors.DefaultCCErrorIf) error {

	if 0 > clsID {
		blog.Error("id is invalid, %d", clsID)
		return fmt.Errorf("id is invalid, %d", clsID)
	}

	if 0 == clsID && 0 == len(val) {
		blog.Error("there are no delete conditions available")
		return fmt.Errorf("there are no delete conditions available")
	}

	if 0 != clsID {
		checkCond := make(map[string]interface{})
		checkCond["id"] = clsID
		val, _ = json.Marshal(checkCond)
	}
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	if items, err := cli.objcli.SearchMetaObjectCls(forward, val); nil != err {
		blog.Error("check the target failed, error:%s ", err.Error())
		return err
	} else if 0 == len(items) {
		blog.Warn("nothing to be deleted")
		return fmt.Errorf("nothing can be deleted, please check the condition")
	}

	return cli.objcli.DeleteMetaObjectCls(forward, clsID, val)
}

// SelectObjectClass select object classification
func (cli *objClsLogic) SelectObjectClass(forward *api.ForwardParam, val []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjClsDes, error) {

	blog.Debug("search classification %v", string(val))
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.SearchMetaObjectCls(forward, val)
}

// SelectObjectClassWithObjects select object classification with associatin objects
func (cli *objClsLogic) SelectObjectClassWithObjects(forward *api.ForwardParam, ownerID string, val []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjClsObjectDes, error) {

	blog.Debug("search classification objects %v", string(val))
	cli.objcli.SetAddress(cli.cfg.Get(cli))
	return cli.objcli.SearchMetaObjectClsObjects(forward, ownerID, val)
}
