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

package metadata

// TopoGraphics define
type TopoGraphics struct {
	NodeType string                 `json:"node_type" bson:"node_type"` // obj inst
	ObjID    string                 `json:"bk_obj_id" bson:"bk_obj_id"`
	InstID   int                    `json:"bk_inst_id" bson:"bk_inst_id"`
	NodeName string                 `json:"node_name" bson:"node_name"`
	Position Position               `json:"position" bson:"position"`
	Ext      map[string]interface{} `json:"ext" bson:"fext"`
	Icon     string                 `json:"bk_obj_icon" bson:"bk_obj_icon"`

	ScopeType string `json:"scope_type" bson:"scope_type"` // biz,user,global,classification
	ScopeID   string `json:"scope_id" bson:"scope_id"`     // ID for ScopeType

	BizID           int    `json:"bk_biz_id" bson:"bk_biz_id"`
	SupplierAccount string `json:"bk_supplier_account" bson:"bk_supplier_account"` // bk_supplier_account
}

// Position the node position in graph
type Position struct {
	X int64 `json:"x" bson:"x"`
	Y int64 `json:"y" bson:"y"`
}
