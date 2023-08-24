/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package importer

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/web_server/service/excel/core"
)

// ImportParamI import excel instance parameter interface
type ImportParamI interface {
	// BuildParam get import instances parameter
	BuildParam(insts map[int]map[string]interface{}) (mapstr.MapStr, error)

	// GetOpType get operation type
	GetOpType() int64

	// GetHandleType get handle type
	GetHandleType() core.HandleType

	// GetBizID get business id
	GetBizID() int64
}

// BaseParam base add data parameter
type BaseParam struct {
	OpType int64 `json:"op"`
	// 用来限定导出关联关系，map[bk_obj_id]object_unique_id 2021年05月17日

	AssociationCond map[string]int64 `json:"association_condition"`

	// 用来限定当前操作对象导出数据的时候，需要使用的唯一校验关系，
	// 自关联的时候，规定左边对象使用到的唯一索引
	ObjectUniqueID int64 `json:"object_unique_id"`
}

// GetOpType get operation type
func (b *BaseParam) GetOpType() int64 {
	return b.OpType
}

// InstParam import instance parameter
type InstParam struct {
	BaseParam `json:",inline"`
	BizID     int64 `json:"bk_biz_id"`
}

// BuildParam get import instances parameter
func (i *InstParam) BuildParam(insts map[int]map[string]interface{}) (mapstr.MapStr, error) {
	param := mapstr.MapStr{
		"input_type":        common.InputTypeExcel,
		"BatchInfo":         insts,
		common.BKAppIDField: i.BizID,
	}

	return param, nil
}

// GetHandleType get handle type
func (i *InstParam) GetHandleType() core.HandleType {
	return core.AddInst
}

// GetBizID get business id
func (i *InstParam) GetBizID() int64 {
	return i.BizID
}

// AddHostParam import add host parameter
type AddHostParam struct {
	BaseParam `json:",inline"`
	ModuleID  int64 `json:"bk_module_id"`
}

// BuildParam get import instances parameter
func (a *AddHostParam) BuildParam(insts map[int]map[string]interface{}) (mapstr.MapStr, error) {
	param := map[string]interface{}{
		"host_info":            insts,
		"input_type":           common.InputTypeExcel,
		common.BKModuleIDField: a.ModuleID,
	}

	return param, nil
}

// GetHandleType get handle type
func (a *AddHostParam) GetHandleType() core.HandleType {
	return core.AddHost
}

// GetBizID get business id
func (a *AddHostParam) GetBizID() int64 {
	return 0
}

// UpdateHostParam excel import update host parameter
type UpdateHostParam struct {
	BaseParam `json:",inline"`
	BizID     int64 `json:"bk_biz_id"`
}

// BuildParam get import instances parameter
func (u *UpdateHostParam) BuildParam(insts map[int]map[string]interface{}) (mapstr.MapStr, error) {
	param := map[string]interface{}{
		"host_info":  insts,
		"input_type": common.InputTypeExcel,
	}

	return param, nil
}

// GetHandleType get handle type
func (u *UpdateHostParam) GetHandleType() core.HandleType {
	return core.UpdateHost
}

// GetBizID get business id
func (u *UpdateHostParam) GetBizID() int64 {
	return u.BizID
}
