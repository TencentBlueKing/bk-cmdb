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
	"crypto/md5"
	"encoding/base64"
	"fmt"

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

// SynchronizeOperateDataType synchronize data operate type
type SynchronizeOperateDataType int64

const (
	// SynchronizeOperateDataTypeInstance synchronize data is instance
	SynchronizeOperateDataTypeInstance SynchronizeOperateDataType = iota + 1
	// SynchronizeOperateDataTypeModel synchronize data is model
	SynchronizeOperateDataTypeModel
	//SynchronizeOperateDataTypeAssociation synchronize data is association
	SynchronizeOperateDataTypeAssociation
)

// SynchronizeDataInfo synchronize instance data http request parameter
type SynchronizeDataInfo struct {
	OperateDataType SynchronizeOperateDataType `json:"operate_data_type"`
	// OperateDataType = SynchronizeOperateDataTypeInstance,
	// DataClassify = object_id,  eg:host,plat,module,proc etc.
	// OperateDataType = SynchronizeOperateDataTypeModel,
	// DataClassify = common.SynchronizeModelDescTypeGroupInfo,common.SynchronizeModelDescTypeModuleAttribute etc
	// OperateDataType = SynchronizeOperateDataTypeAssociation
	// DataClassify = common.SynchronizeAssociationTypeModelHost etc.
	DataClassify string          `json:"data_classify"`
	InfoArray    []mapstr.MapStr `json:"instance_info_array"`
	// OffSet current data offset  start location
	Offset int64 `json:"offset"`
	// Count total data count
	Count           int64  `json:"count"`
	Version         int64  `json:"version"`
	SynchronizeFlag string `json:"synchronize_flag"`
}

// SynchronizeParameter synchronize instance data http request parameter
type SynchronizeParameter struct {
	OperateType SynchronizeOperateType `json:"op_type"`
	// synchronize data type
	OperateDataType SynchronizeOperateDataType `json:"operate_data_type"`
	// OperateDataType = SynchronizeOperateDataTypeInstance,
	// DataClassify = object_id,  eg:host,plat,module,proc etc.
	// OperateDataType = SynchronizeOperateDataTypeModel,
	// DataClassify = common.SynchronizeModelDescTypeGroupInfo,common.SynchronizeModelDescTypeModuleAttribute etc
	// OperateDataType = SynchronizeOperateDataTypeAssociation
	// DataClassify = common.SynchronizeAssociationTypeModelHost etc.
	DataClassify    string             `json:"data_classify"`
	InfoArray       []*SynchronizeItem `json:"instance_info_array"`
	Version         int64              `json:"version"`
	SynchronizeFlag string             `json:"synchronize_flag"`
}

// SynchronizeItem synchronize data information
type SynchronizeItem struct {
	Info mapstr.MapStr `json:"info"`
	ID   int64         `json:"id"`
}

// SynchronizeFindInfoParameter synchronize  data fetch data http request parameter
type SynchronizeFindInfoParameter struct {
	DataType     SynchronizeOperateDataType `json:"data_type"`
	DataClassify string                     `json:"data_classify"`
	Condition    mapstr.MapStr              `json:"condition"`
	Start        uint64                     `json:"start"`
	Limit        uint64                     `json:"limit"`
}

// SynchronizeResult synchronize result
type SynchronizeResult struct {
	BaseResp `json:",inline"`
	Data     SetDataResult `json:"data"`
}

// SynchronizeDataResult common Synchronize result definition
type SynchronizeDataResult struct {
	//Created    []CreatedDataResult `json:"created"`
	//Updated    []UpdatedDataResult `json:"updated"`
	Exceptions []ExceptionResult `json:"exception"`
}

// SynchronizeClearDataParameter synchronize  data clear data http request parameter
type SynchronizeClearDataParameter struct {
	Tamestamp       int64  `json:"tamestamp"`
	Sign            string `json:"sign"`
	Version         int64  `json:"version"`
	SynchronizeFlag string `json:"synchronizeFlag"`
}

// GenerateSign generate sign
func (s *SynchronizeClearDataParameter) GenerateSign(key string) {
	m := md5.New()
	m.Write([]byte(s.signContext(key)))
	s.Sign = base64.StdEncoding.EncodeToString((m.Sum(nil)))
}

//  Legality sign is legal
func (s *SynchronizeClearDataParameter) Legality(key string) bool {
	m := md5.New()
	m.Write([]byte(s.signContext(key)))
	// tamestamp need legality
	return s.Sign == base64.StdEncoding.EncodeToString((m.Sum(nil)))
}

// GenerateSign generate sign
func (s *SynchronizeClearDataParameter) signContext(key string) string {
	return fmt.Sprintf("key-%s-%s-%d-%d", key, s.SynchronizeFlag, s.Tamestamp, s.Version)
}

// SetIdenifierFlag update idenifier flag data
type SetIdenifierFlag struct {
	DataType SynchronizeOperateDataType `json:"data_type"`
	// DataType = SynchronizeOperateDataTypeInstance,
	// DataType = object_id,  eg:host,plat,module,proc etc.
	// DataType = SynchronizeOperateDataTypeModel,
	// DataClassify = common.SynchronizeModelDescTypeGroupInfo,common.SynchronizeModelDescTypeModuleAttribute etc
	DataClassify string `json:"data_classify"`

	IdentifierID []int64 `json:"identifier_id"`

	// 需要同步到第三方系统的身份标志
	Flag string `json:"flag"`
	// 1:新增,在原有的基础新加同步标志。
	// 2:覆盖，删除 原有的同步标志
	// 3:删除, 删除同步标志
	OperateType SynchronizeOperateType `json:"op_type"`
}
