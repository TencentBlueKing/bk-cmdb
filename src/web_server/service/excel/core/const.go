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

package core

// 此文件存放公共的，需要暴露给其他文件的常量定义

// excel file const define
const (
	// InstHeaderLen excel实例表头所占行数
	InstHeaderLen = 6

	// HeaderTableLen excel表头里，表格相关所占用行数
	HeaderTableLen = 3

	// NameRowIdx excel表头「字段名」所在行位置
	NameRowIdx = 0

	// TypeRowIdx excel表头「字段类型」所在行位置
	TypeRowIdx = 1

	// IDRowIdx excel表头「字段标识」所在行位置
	IDRowIdx = 2

	// TableNameRowIdx excel表头「表格字段名」所在行位置
	TableNameRowIdx = 3

	// TableTypeRowIdx excel表头「表格字段类型」所在行位置
	TableTypeRowIdx = 4

	// TableIDRowIdx excel表头「表格字段标识」所在行位置
	TableIDRowIdx = 5

	// InstRowIdx excel 实例数据开始的位置
	InstRowIdx = 6

	// AsstSheet excel关联关系sheet名称
	AsstSheet = "association"

	// AsstStartRowIdx excel关联关系sheet开始行位置
	AsstStartRowIdx = 0

	// AsstExampleRowIdx excel关联关系sheet例子所在行位置
	AsstExampleRowIdx = 1

	// AsstIDColIdx excel关联关系sheet「关联标识」所在列位置
	AsstIDColIdx = 1

	// AsstOPColIdx excel关联关系sheet「操作」所在列位置
	AsstOPColIdx = 2

	// AsstSrcInstColIdx excel关联关系sheet「源实例」所在列位置
	AsstSrcInstColIdx = 3

	// AsstDstInstColIdx excel关联关系sheet「目标实例」所在列位置
	AsstDstInstColIdx = 4

	// AsstDataRowIdx excel关联关系sheet数据开始的位置
	AsstDataRowIdx = 2
)

// export instance const define
const (
	// TopoObjID 导出主机实例时，「业务拓扑」这一字段的objID
	TopoObjID = "field_topo"

	// IDPrefix 导出主机实例时，字段标识的前缀
	IDPrefix = "bk_ext_"
)

const (
	// PropDefaultColIdx 属性所在列的默认值
	PropDefaultColIdx = 0
)

// HandleType 导入excel操作类型
type HandleType string

const (
	// AddHost 添加主机
	AddHost HandleType = "addHost"

	// UpdateHost 更新主机
	UpdateHost HandleType = "updateHost"

	// AddInst 添加实例
	AddInst HandleType = "addInst"
)

// AsstOp 关联关系操作
type AsstOp string

const (
	// AsstOpAdd 添加关联关系操作
	AsstOpAdd AsstOp = "add"

	// AsstOpDel 删除关联关系操作
	AsstOpDel AsstOp = "delete"
)

// AsstOps 关联关系操作数组
var AsstOps = []string{string(AsstOpAdd), string(AsstOpDel)}
