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
	// SynchronizeOperateTypeAdd synchroneize data add
	SynchronizeOperateTypeAdd SynchronizeOperateType = iota + 1
	// SynchronizeOperateTypeUpdate synchroneize data update
	SynchronizeOperateTypeUpdate
	// SynchronizeOperateTypeRepalce synchroneize data add or update
	SynchronizeOperateTypeRepalce
	// SynchronizeOperateTypeDelete synchroneize data delete
	SynchronizeOperateTypeDelete
)

// SynchronizeInstanceParameter synchronize instance data http request parameter
type SynchronizeInstanceParameter struct {
	OperateType       SynchronizeOperateType `json:"op_type"`
	ObjectID          string                 `json:"obj_id"`
	InstacneInfoArray []*SynchronizeItem     `json:"instance_info_array"`
	// source data sign
	SynchronizeSign string `json:"sync_sign"`
}

// SynchronizeItem synchronize data information
type SynchronizeItem struct {
	Info mapstr.MapStr `json:"info"`
	ID   int64         `json:"id"`
}

// SynchronizeModelDescType synchronize model data type
type SynchronizeModelDescType int

const (
	// SynchronizeModelDescTypeGroupInfo synchroneize model ggroup
	SynchronizeModelDescTypeGroupInfo SynchronizeModelDescType = iota + 1
	// SynchronizeModelDescTypeModuleAttribute synchroneize model attribute
	SynchronizeModelDescTypeModuleAttribute
	// SynchronizeModelDescTypeModuleAttributeGroup synchroneize model attribute group
	SynchronizeModelDescTypeModuleAttributeGroup
)

// SynchronizeModelDataParameter  synchronize model data http request parameter
type SynchronizeModelDescParameter struct {
	OperateType        SynchronizeOperateType   `json:"op_type"`
	ModelDescription   SynchronizeModelDescType `json:"model_desc"`
	ModelDescInfoArray []*SynchronizeItem       `json:"model_attr_array"`
	// source data sign
	SynchronizeSign string `json:"sync_sign"`
}
