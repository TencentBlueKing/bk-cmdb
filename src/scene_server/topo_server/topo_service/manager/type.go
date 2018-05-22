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
)

// ObjectAsst const definition
const ObjectAsst = "object_asst"

// ObjectAttribute const definition
const ObjectAttribute = "object_attribute"

// ObjectClass const definition
const ObjectClass = "object_class"

// ObjectGroup const definition
const ObjectGroup = "object_group"

// Object const definition
const Object = "object"

// TopoModelRsp 拓扑模型结构
type TopoModelRsp struct {
	ObjID      string `json:"bk_obj_id"`
	ObjName    string `json:"bk_obj_name"`
	OwnerID    string `json:"bk_supplier_account"`
	NextObj    string `json:"bk_next_obj"`
	NextName   string `json:"bk_next_name"`
	PreObjID   string `json:"bk_pre_obj_id"`
	PreObjName string `json:"bk_pre_obj_name"`
}

// TopoInst 实例拓扑结构
type TopoInst struct {
	InstID   int    `json:"bk_inst_id"`
	InstName string `json:"bk_inst_name"`
	ObjID    string `json:"bk_obj_id"`
	ObjName  string `json:"bk_obj_name"`
	Default  int    `json:"default"`
}

// TopoInstRst 拓扑实例
type TopoInstRst struct {
	TopoInst
	Child []TopoInstRst `json:"child"`
}

// ObjectAsstLogic define the logic interface
type ObjectAsstLogic interface {
	CreateObjectAsst(forward *api.ForwardParam, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) (int, error)
	SelectObjectAsst(forward *api.ForwardParam, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) ([]api.ObjAsstDes, error)
	UpdateObjectAsst(forward *api.ForwardParam, selector, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) error
	DeleteObjectAsstByID(forward *api.ForwardParam, id int, errProxy errors.DefaultCCErrorIf) error
	DeleteObjectAsst(forward *api.ForwardParam, obj map[string]interface{}, errProxy errors.DefaultCCErrorIf) error
}

// ObjectClassLogic define the logic interface
type ObjectClassLogic interface {
	CreateObjectClass(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) (int, error)
	SelectObjectClass(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjClsDes, error)
	SelectObjectClassWithObjects(forward *api.ForwardParam, ownerID string, params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjClsObjectDes, error)
	UpdateObjectClass(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error
	DeleteObjectClass(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error
}

// ObjectLogic define the logic interface
type ObjectLogic interface {
	CreateObject(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) (int, error)
	SelectObject(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjDes, error)
	UpdateObject(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error
	DeleteObject(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error
}

// TopoGraphics define the TopoGraphics interface
type TopoGraphics interface {
	SearchGraphics(forward *api.ForwardParam, params map[string]interface{}, errProxy errors.DefaultCCErrorIf) ([]api.TopoGraphics, error)
}

// ObjectAttLogic define the logic interface
type ObjectAttLogic interface {
	CreateTopoModel(forward *api.ForwardParam, obj api.ObjAttDes, errProxy errors.DefaultCCErrorIf) (int, error)
	SelectTopoModel(forward *api.ForwardParam, rstitems []TopoModelRsp, ownerid, objid, clsid, preid, prename string, errProxy errors.DefaultCCErrorIf) ([]TopoModelRsp, error)
	DeleteTopoModel(forward *api.ForwardParam, ownerid, objid string, assotype int, errProxy errors.DefaultCCErrorIf) error
	CreateObjectAtt(forward *api.ForwardParam, params api.ObjAttDes, errProxy errors.DefaultCCErrorIf) (int, error)
	SelectObjectAtt(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjAttDes, error)
	UpdateObjectAtt(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error
	DeleteObjectAtt(forward *api.ForwardParam, id int, params []byte, errProxy errors.DefaultCCErrorIf) error
}

// ObjectAttGroupLogic define the logic interface
type ObjectAttGroupLogic interface {
	CreateObjectGroup(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) (int, error)
	UpdateObjectGroup(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) error
	UpdateObjectGroupProperty(forward *api.ForwardParam, params []byte, errProxy errors.DefaultCCErrorIf) error
	DeleteObjectGroup(forward *api.ForwardParam, id int, errProxy errors.DefaultCCErrorIf) error
	DeleteObjectGroupProperty(forward *api.ForwardParam, ownerID, objectID, propertyID, groupID string, errProxy errors.DefaultCCErrorIf) error
	SelectPropertyGroupByObjectID(forward *api.ForwardParam, ownerID, objectID string, data []byte, errProxy errors.DefaultCCErrorIf) ([]api.ObjAttGroupDes, error)
}

// Manager define manager interface
type Manager interface {

	// object asst interface
	ObjectAsstLogic

	// object class interface
	ObjectClassLogic

	// object interface
	ObjectLogic

	// object attribute interface
	ObjectAttLogic

	// object attribute group interface
	ObjectAttGroupLogic

	// TopoGraphics interface
	TopoGraphics
}

// Hooker define callback hook
type Hooker interface {
	SetManager(mgr Manager) error
}

// Configer define the configer interface
type Configer interface {
	Get(target interface{}) string
}

// SetConfiger define the set configer interface
type SetConfiger interface {
	Set(cfg Configer)
}
