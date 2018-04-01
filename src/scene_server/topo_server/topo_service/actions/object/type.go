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
	"configcenter/src/scene_server/topo_server/topo_service/manager"

	restful "github.com/emicklei/go-restful"
)

// Page 用于分页查询
type Page struct {
	Sort  string `json:"sort"`
	Limit int    `json:"limit"`
	Skip  int    `json:"start"`
}

// InstItem item 结构
type InstItem struct {
	Count int                      `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

// ObjInstRsp 用于提取controller 返回的实例数据结构
type ObjInstRsp struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    InstItem    `json:"data"`
}

// TopoItem define topo item
type TopoItem struct {
	ClassificationID string `json:"bk_classification_id"`
	Position         string `json:"position"`
	ObjID            string `json:"bk_obj_id"`
	OwnerID          string `json:"bk_supplier_account"`
	ObjName          string `json:"bk_obj_name"`
}

// ObjectTopo define the common object topo
type ObjectTopo struct {
	LabelType string   `json:"label_type"`
	LabelName string   `json:"label_name"`
	Label     string   `json:"label"`
	From      TopoItem `json:"from"`
	To        TopoItem `json:"to"`
	Arrows    string   `json:"arrows"`
}

// HelperFunction helper function definition
type HelperFunction struct {
	// SelectInstTopo return the inst topo
	SelectInstTopo func(ownerid, objid string, appid, instid, level int, req *restful.Request) ([]manager.TopoInstRst, error)
}

// CreateHelperFunction return the helper function items
func CreateHelperFunction() *HelperFunction {
	return &HelperFunction{
		SelectInstTopo: topo.SelectInstTopo,
	}
}
