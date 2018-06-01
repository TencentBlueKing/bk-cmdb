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
	"configcenter/src/source_controller/api/metadata"
	"errors"
	"net/http"
)

var (
	Err_Not_Found_Anything = errors.New("found nothing")
	Err_Not_Set_Input      = errors.New("input nothing")
	Err_Not_Set_ObjID      = errors.New("not set objid")
	Err_Request_Object     = errors.New("http request failed")
	Err_Decode_Json        = errors.New("decode json failed")
	Err_Creaate_Object     = errors.New("create object failed")
)

// ObjDes
type ObjDes struct {
	metadata.ObjectDes `json:",inline"`
}

// ObjAsstDes association
type ObjAsstDes struct {
	metadata.ObjectAsst `json:",inline"`
}

// ObjAttDes 对象模型属性
type ObjAttDes struct {
	metadata.ObjectAttDes `json:",inline"`
	AssoType              int    `json:"bk_asst_type"`
	AsstForward           string `json:"bk_asst_forward"`
	AssociationID         string `json:"bk_asst_obj_id"`
	PropertyGroupName     string `json:"bk_property_group_name"`
}

// ForwardParam define logic layer common param
type ForwardParam struct {
	Header http.Header
}

// ObjAttGroupDes define property group
type ObjAttGroupDes struct {
	metadata.PropertyGroup `json:",inline"`
}

// ObjClsDes 对象分类（分栏/分组)
type ObjClsDes struct {
	metadata.ObjClassification `json:",inline"`
}

// ObjClsObjectDes 分类下的对象模型
type ObjClsObjectDes struct {
	metadata.ObjClassificationObject `json:",inline"`
}

// TopoGraphics Topo Graphics
type TopoGraphics struct {
	metadata.TopoGraphics `json:",inline"`
	Assts                 []GraphAsst `json:"assts,omitempty"`
}

// Asst the node association node
type GraphAsst struct {
	AsstType string            `json:"bk_asst_type"`
	NodeType string            `json:"node_type"`
	ObjID    string            `json:"bk_obj_id"`
	InstID   int               `json:"bk_inst_id"`
	ObjAtt   string            `json:"bk_object_att_id"`
	Lable    map[string]string `json:"lable"`
}

// ObjDesRsp 用于提取congtroller 返回的数据结构
type ObjDesRsp struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    []ObjDes    `json:"data"`
}

// ObjAsstRsp 用于提取congtroller 返回的数据结构
type ObjAsstRsp struct {
	Result  bool         `json:"result"`
	Code    int          `json:"code"`
	Message interface{}  `json:"message"`
	Data    []ObjAsstDes `json:"data"`
}

// ObjAttRsp  用于提取controller 返回的数据结构
type ObjAttRsp struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    []ObjAttDes `json:"data"`
}

// ObjAttGroupRsp  用于提取controller 返回的数据结构
type ObjAttGroupRsp struct {
	Result  bool             `json:"result"`
	Code    int              `json:"code"`
	Message interface{}      `json:"message"`
	Data    []ObjAttGroupDes `json:"data"`
}

// ObjClsRsp 用于提起controller返回的数据结构
type ObjClsRsp struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    []ObjClsDes `json:"data"`
}

// ObjClsObjectRsp 用于提起controller返回的数据结构
type ObjClsObjectRsp struct {
	Result  bool              `json:"result"`
	Code    int               `json:"code"`
	Message interface{}       `json:"message"`
	Data    []ObjClsObjectDes `json:"data"`
}
