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

// GraphAsst TODO
// Asst the node association node
type GraphAsst struct {
	AsstType              string            `json:"bk_asst_type"`
	NodeType              string            `json:"node_type"`
	ObjID                 string            `json:"bk_obj_id"`
	InstID                int64             `json:"bk_inst_id"`
	AssociationKindInstID int64             `json:"bk_asst_inst_id"`
	AsstName              string            `json:"bk_asst_name"`
	Label                 map[string]string `json:"label"`
}

// TopoGraphics TODO
type TopoGraphics struct {
	ScopeType       string                 `json:"scope_type,omitempty" bson:"scope_type,omitempty"` // biz,user,global,classification
	ScopeID         string                 `json:"scope_id,omitempty" bson:"scope_id,omitempty"`     // ID for ScopeType
	NodeType        string                 `json:"node_type" bson:"node_type"`                       // obj inst
	ObjID           string                 `json:"bk_obj_id" bson:"bk_obj_id"`
	IsPre           bool                   `json:"ispre"             bson:"ispre"`
	InstID          int                    `json:"bk_inst_id" bson:"bk_inst_id"`
	NodeName        string                 `json:"node_name,omitempty" bson:"node_name,omitempty"`
	Position        Position               `json:"position" bson:"position"`
	Ext             map[string]interface{} `json:"ext,omitempty" bson:"ext,omitempty"`
	Icon            string                 `json:"bk_obj_icon,omitempty" bson:"bk_obj_icon,omitempty"`
	SupplierAccount string                 `json:"bk_supplier_account,omitempty" bson:"bk_supplier_account,omitempty"` // bk_supplier_account
	Assts           []GraphAsst            `json:"assts,omitempty"`
}

// UpdateTopoGraphicsInput TODO
type UpdateTopoGraphicsInput struct {
	Origin []TopoGraphics `field:"origin" json:"origin" bson:"origin"`
}

// FillBlank TODO
func (t *TopoGraphics) FillBlank() *TopoGraphics {
	t.SetSupplierAccount("0")
	t.SetExt(map[string]interface{}{})
	t.SetPosition(Position{})
	return t
}

// SetNodeType TODO
func (t *TopoGraphics) SetNodeType(val string) { t.NodeType = val }

// SetObjID TODO
func (t *TopoGraphics) SetObjID(val string) { t.ObjID = val }

// SetInstID TODO
func (t *TopoGraphics) SetInstID(val int) { t.InstID = val }

// SetNodeName TODO
func (t *TopoGraphics) SetNodeName(val string) { t.NodeName = val }

// SetIsPre TODO
func (t *TopoGraphics) SetIsPre(val bool) { t.IsPre = val }

// SetPosition TODO
func (t *TopoGraphics) SetPosition(val Position) {
	t.Position = val
}

// SetExt TODO
func (t *TopoGraphics) SetExt(val map[string]interface{}) { t.Ext = val }

// SetIcon TODO
func (t *TopoGraphics) SetIcon(val string) { t.Icon = val }

// SetScopeType TODO
func (t *TopoGraphics) SetScopeType(val string) { t.ScopeType = val }

// SetScopeID TODO
func (t *TopoGraphics) SetScopeID(val string) { t.ScopeID = val }

// SetSupplierAccount TODO
func (t *TopoGraphics) SetSupplierAccount(val string) { t.SupplierAccount = val }

// SearchTopoGraphicsResult TODO
type SearchTopoGraphicsResult struct {
	BaseResp `json:",inline"`
	Data     []TopoGraphics `json:"data"`
}
