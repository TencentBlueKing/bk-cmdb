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

type TopoGraphics struct {
    ScopeType       *string                `json:"scope_type,omitempty" bson:"scope_type,omitempty"` // biz,user,global,classification
    ScopeID         *string                `json:"scope_id,omitempty" bson:"scope_id,omitempty"`     // ID for ScopeType
    NodeType        *string                `json:"node_type,omitempty" bson:"node_type,omitempty"`   // obj inst
    ObjID           *string                `json:"bk_obj_id,omitempty" bson:"bk_obj_id,omitempty"`
    IsPre           *bool                  `json:"ispre,omitempty"             bson:"ispre,omitempty"`
    InstID          *int                   `json:"bk_inst_id,omitempty" bson:"bk_inst_id,omitempty"`
    NodeName        *string                `json:"node_name,omitempty" bson:"node_name,omitempty"`
    Position        *Position              `json:"position,omitempty" bson:"position,omitempty"`
    Ext             map[string]interface{} `json:"ext,omitempty" bson:"ext,omitempty"`
    Icon            *string                `json:"bk_obj_icon,omitempty" bson:"bk_obj_icon,omitempty"`
    BizID           *int                   `json:"bk_biz_id,omitempty" bson:"bk_biz_id,omitempty"`
    SupplierAccount *string                `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account,omitempty"` // bk_supplier_account
}

// Position the node position in graph
type Position struct {
    X *int64 `json:"x" bson:"x"`
    Y *int64 `json:"y" bson:"y"`
}
