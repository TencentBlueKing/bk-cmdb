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

// Position the node position in graph
type Position struct {
	X *int64 `json:"x" bson:"x"`
	Y *int64 `json:"y" bson:"y"`
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
	Assts           []GraphAsst            `json:"assts,omitempty"`
}

func (t *TopoGraphics) FillBlank() *TopoGraphics {
	if t.BizID == nil {
		t.SetBizID(0)
	}
	if t.SupplierAccount == nil {
		t.SetSupplierAccount("0")
	}
	if t.Ext == nil {
		t.SetExt(map[string]interface{}{})
	}
	if t.Position == nil {
		t.SetPosition(&Position{})
	}
	return t
}

func (t *TopoGraphics) SetNodeType(val string) { t.NodeType = &val }
func (t *TopoGraphics) SetObjID(val string)    { t.ObjID = &val }
func (t *TopoGraphics) SetInstID(val int)      { t.InstID = &val }
func (t *TopoGraphics) SetNodeName(val string) { t.NodeName = &val }
func (t *TopoGraphics) SetIsPre(val bool)      { t.IsPre = &val }
func (t *TopoGraphics) SetPosition(val *Position) {
	if val == nil {
		t.Position = &Position{}
	} else {
		t.Position = val
	}
}
func (t *TopoGraphics) SetExt(val map[string]interface{}) { t.Ext = val }
func (t *TopoGraphics) SetIcon(val string)                { t.Icon = &val }
func (t *TopoGraphics) SetScopeType(val string)           { t.ScopeType = &val }
func (t *TopoGraphics) SetScopeID(val string)             { t.ScopeID = &val }
func (t *TopoGraphics) SetBizID(val int)                  { t.BizID = &val }
func (t *TopoGraphics) SetSupplierAccount(val string)     { t.SupplierAccount = &val }

type SearchTopoGraphicsResult struct {
	BaseResp `json:",inline"`
	Data     []TopoGraphics `json:"data"`
}
