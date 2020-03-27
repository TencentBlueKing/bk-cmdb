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

package operation

import (
	"time"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/model"
)

type opcondition struct {
	InstID []int64 `json:"inst_ids"`
}

type deleteCondition struct {
	opcondition `json:",inline"`
}

type updateCondition struct {
	InstID   int64                  `json:"inst_id"`
	InstInfo map[string]interface{} `json:"datas"`
}

// OpCondition the condition operation
type OpCondition struct {
	Delete deleteCondition   `json:"delete"`
	Update []updateCondition `json:"update"`
}

type InstBatchInfo struct {
	// BatchInfo batch info
	// map[rownumber]map[property_id][date]
	BatchInfo map[int64]map[string]interface{} `json:"BatchInfo"`
	InputType string                           `json:"input_type"`
}

// ConditionItem subcondition
type ConditionItem struct {
	Field    string      `json:"field,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}

// AssociationParams  association params
type AssociationParams struct {
	Page      metadata.BasePage          `json:"page,omitempty"`
	Fields    map[string][]string        `json:"fields,omitempty"`
	Condition map[string][]ConditionItem `json:"condition,omitempty"`
}

// commonInstTopo common inst topo
type CommonInstTopo struct {
	metadata.InstNameAsst
	Count    int                     `json:"count"`
	Children []metadata.InstNameAsst `json:"children"`
}

type CommonInstTopoV2 struct {
	Prev []*CommonInstTopo `json:"prev"`
	Next []*CommonInstTopo `json:"next"`
	Curr interface{}       `json:"curr"`
}

type deletedInst struct {
	instID int64
	obj    model.Object
}

// OperationLog opeartion log item definition
type OperationLog struct {
	OwnerID       string      `bson:"bk_supplier_account"    json:"bk_supplier_account"`
	ApplicationID int         `bson:"bk_biz_id"              json:"bk_biz_id"`
	ExtKey        string      `bson:"ext_key"             json:"ext_key"`
	OpDesc        string      `bson:"op_desc"             json:"op_desc"`
	OpType        int         `bson:"op_type"             json:"op_type"`
	OpTarget      string      `bson:"op_target"           json:"op_target"`
	Content       interface{} `bson:"content"             json:"content"`
	User          string      `bson:"operator"                json:"operator"`
	OpFrom        string      `bson:"op_from"             json:"op_from"`
	ExtInfo       string      `bson:"ext_info"            json:"ext_info"`
	CreateTime    time.Time   `bson:"op_time"         json:"op_time"`
	InstID        int         `bson:"inst_id"             json:"inst_id"`
}

type Content struct {
	PreData interface{} `json:"pre_data"`
	CurData interface{} `json:"cur_data"`
	Headers []Header    `json:"header"`
}

type Header struct {
	PropertyID   string `json:"bk_property_id"`
	PropertyName string `json:"bk_property_name"`
}

type Ref struct {
	RefID   int    `json:"ref_id"`
	RefName string `json:"ref_name"`
}

type ExportObjectCondition struct {
	ObjIDS []string `json:"condition"`
}

type ImportObjectData struct {
	Meta mapstr.MapStr           `json:"meta"`
	Attr map[int64]mapstr.MapStr `json:"attr"`
}
