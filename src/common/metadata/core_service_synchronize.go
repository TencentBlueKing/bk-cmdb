/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	"configcenter/src/common/mapstr"
)

// SynchronizeOperateType synchronize data operate type
type SynchronizeOperateType int64

const (
	// SynchronizeOperateTypeAdd synchronize data add
	SynchronizeOperateTypeAdd SynchronizeOperateType = iota + 1
	// SynchronizeOperateTypeUpdate synchronize data update
	SynchronizeOperateTypeUpdate
	// SynchronizeOperateTypeRepalce synchronize data add or update
	SynchronizeOperateTypeRepalce
	// SynchronizeOperateTypeDelete synchronize data delete
	SynchronizeOperateTypeDelete
)

// SynchronizeOperateType synchronize data operate type
type SynchronizeDataType int64

const (
	// SynchronizeDataTypeInstance synchronize data is instance
	SynchronizeDataTypeInstance SynchronizeDataType = iota + 1
	// SynchronizeDataTypeModel synchronize data is model
	SynchronizeDataTypeModel
	//SynchronizeDataTypeAssociation synchronize data is association
	SynchronizeDataTypeAssociation
)

// SynchronizeInstanceParameter synchronize instance data http request parameter
type SynchronizeParameter struct {
	OperateType SynchronizeOperateType `json:"op_type"`
	// synchronize data type
	DataType SynchronizeDataType `json:"data_type"`
	// DataType = SynchronizeDataTypeInstance,
	// DataSign = object_id,  eg:host,plat,module,proc etc.
	// DataType = SynchronizeDataTypeModel,
	// DataSign = common.SynchronizeModelDescTypeGroupInfo,common.SynchronizeModelDescTypeModuleAttribute etc
	// DataType = SynchronizeDataTypeAssociation
	// DataSign = common.SynchronizeAssociationTypeModelHost etc.
	DataSign  string             `json:"data_sign"`
	InfoArray []*SynchronizeItem `json:"instance_info_array"`
	// source data sign
	SynchronizeSign string `json:"sync_sign"`
}

// SynchronizeItem synchronize data information
type SynchronizeItem struct {
	Info mapstr.MapStr `json:"info"`
	ID   int64         `json:"id"`
}
