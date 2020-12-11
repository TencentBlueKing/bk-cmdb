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

package topology

type BizBriefTopology struct {
	// basic business info
	Biz *BizBase `json:"biz"`
	// the idle set nodes info
	Idle []*Node `json:"idle"`
	// the other common nodes
	Nodes []*Node `json:"nds"`
}

type Node struct {
	// the object of this node, like set or module
	Object string `json:"obj"`
	// the node's instance id, like set id or module id
	ID int64 `json:"id"`
	// the node's name, like set name or module name
	Name string `json:"nm"`
	// only set, module has this field.
	// describe what kind of set or module this node is.
	// 0: normal module or set.
	// >1: special set or module
	Default *int `json:"dft,omitempty"`
	// the sub-nodes of current node
	SubNodes []*Node `json:"nds"`
}

var bizBaseFields = []string{"bk_biz_id", "bk_biz_name", "default"}

type BizBase struct {
	// business id
	ID int64 `json:"id" bson:"bk_biz_id"`
	// business name
	Name string `json:"nm" bson:"bk_biz_name"`
	// describe it's a resource pool business or normal business.
	// 0: normal business
	// >0: special business, like resource pool business.
	Default int `json:"dft" bson:"default"`
}

var customBaseFields = []string{"bk_biz_id", "bk_parent_id", "bk_inst_id", "bk_inst_name", "bk_obj_id"}

type customBase struct {
	Business int64  `bson:"bk_biz_id"`
	ParentID int64  `bson:"bk_parent_id"`
	ID       int64  `bson:"bk_inst_id"`
	Name     string `bson:"bk_inst_name"`
	Object   string `bson:"bk_obj_id"`
}

type customArchive struct {
	Oid    string      `bson:"oid"`
	Detail *customBase `bson:"detail"`
}

var setBaseFields = []string{"bk_biz_id", "bk_parent_id", "bk_set_id", "bk_set_name", "default"}

type setBase struct {
	Business int64  `bson:"bk_biz_id"`
	ParentID int64  `bson:"bk_parent_id"`
	ID       int64  `bson:"bk_set_id"`
	Name     string `bson:"bk_set_name"`
	Default  int    `bson:"default"`
}

type setArchive struct {
	Oid    string   `bson:"oid"`
	Detail *setBase `bson:"detail"`
}

var moduleBaseFields = []string{"bk_biz_id", "bk_set_id", "bk_module_id", "bk_module_name", "default"}

type moduleBase struct {
	Business int64  `bson:"bk_biz_id"`
	SetID    int64  `bson:"bk_set_id"`
	ID       int64  `bson:"bk_module_id"`
	Name     string `bson:"bk_module_name"`
	Default  int    `bson:"default"`
}

type moduleArchive struct {
	Oid    string      `bson:"oid"`
	Detail *moduleBase `bson:"detail"`
}

var mainlineAsstFields = []string{"bk_asst_obj_id", "bk_obj_id"}

type mainlineAssociation struct {
	AssociateTo string `bson:"bk_asst_obj_id"`
	ObjectID    string `bson:"bk_obj_id"`
}

// page step
const (
	step                          = 100
	defaultRefreshIntervalMinutes = 15
)
