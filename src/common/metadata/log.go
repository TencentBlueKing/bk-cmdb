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

import "time"

// OperationLog opeartion log item definition
type OperationLog struct {
	OwnerID       string      `bson:"bk_supplier_account"    json:"bk_supplier_account"`
	ApplicationID int64       `bson:"bk_biz_id"              json:"bk_biz_id"`
	ExtKey        string      `bson:"ext_key"             json:"ext_key"`
	OpDesc        string      `bson:"op_desc"             json:"op_desc"`
	OpType        int         `bson:"op_type"             json:"op_type"`
	OpTarget      string      `bson:"op_target"           json:"op_target"`
	Content       interface{} `bson:"content"             json:"content"`
	User          string      `bson:"operator"                json:"operator"`
	OpFrom        string      `bson:"op_from"             json:"op_from"`
	ExtInfo       string      `bson:"ext_info"            json:"ext_info"`
	CreateTime    time.Time   `bson:"op_time"         json:"op_time"`
	InstID        int64       `bson:"inst_id"             json:"inst_id"`
}

// TableName return the table name
func (OperationLog) TableName() string {
	return "cc_OperationLog"
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
	RefID   int64  `json:"ref_id"`
	RefName string `json:"ref_name"`
}
