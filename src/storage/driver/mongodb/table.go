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

package mongodb

import (
	"configcenter/src/common"
	"configcenter/src/storage/dal/types"
)

/*
FILE:  构建数据库中的表对象
Timed: 2020年11月08日
Description:

TB 前缀: 为Table 缩写，表示为collection 对象类型的数据

注意：差异环境table 对象，可以新开一个文件来存取

*/

func TBHostBase() types.Table {
	return Table(common.BKTableNameBaseHost)
}

func TBApplicationBase() types.Table {
	return Table(common.BKTableNameBaseApp)
}

func TBSetBase() types.Table {
	return Table(common.BKTableNameBaseSet)
}

func TBModuleBase() types.Table {
	return Table(common.BKTableNameBaseModule)
}

// 模型实例数据表，TBInstance别名函数
func TBModelInstance() types.Table {
	return Table(common.BKTableNameBaseInst)
}

// 模型实例数据表, TBModelInstance 别名函数
func TBInstance() types.Table {
	return Table(common.BKTableNameBaseInst)
}

func TBCloudArea() types.Table {
	return Table(common.BKTableNameBasePlat)
}

func TBProcess() types.Table {
	return Table(common.BKTableNameBaseProcess)
}

func TBPropertyGroup() types.Table {
	return Table(common.BKTableNamePropertyGroup)
}

// 模型关联关系
func TBAsstDes() types.Table {
	return Table(common.BKTableNameAsstDes)
}

func TBObjDesc() types.Table {
	return Table(common.BKTableNameObjDes)
}

func TBObjUnique() types.Table {
	return Table(common.BKTableNameObjUnique)
}

func TBObjAttrDesc() types.Table {
	return Table(common.BKTableNameObjAttDes)
}

func TBObjClassification() types.Table {
	return Table(common.BKTableNameObjClassification)
}

// 实例关系数据， TBInstAsst别名函数
func TBInstnceAsst() types.Table {
	return Table(common.BKTableNameInstAsst)
}

// 实例关系数据， TBInstnceAsst别名函数
func TBInstAsst() types.Table {
	return Table(common.BKTableNameInstAsst)
}

func TBModuleHostConfig() types.Table {
	return Table(common.BKTableNameModuleHostConfig)
}

func TBAuditLog() types.Table {
	return Table(common.BKTableNameAuditLog)
}

func TBObjectAsst() types.Table {
	return Table(common.BKTableNameObjAsst)
}

func TB() types.Table {
	return Table(common.BKTableNameObjAsst)
}

func TBServiceCategory() types.Table {
	return Table(common.BKTableNameServiceCategory)
}

func TBServiceTemplate() types.Table {
	return Table(common.BKTableNameServiceTemplate)
}

func TBServiceInstance() types.Table {
	return Table(common.BKTableNameServiceInstance)
}

func TBProcessTemplate() types.Table {
	return Table(common.BKTableNameProcessTemplate)
}

func TBProcessInstanceRelation() types.Table {
	return Table(common.BKTableNameProcessInstanceRelation)
}

func TBSetTemplate() types.Table {
	return Table(common.BKTableNameSetTemplate)
}

func TBSetServiceTemplateRelation() types.Table {
	return Table(common.BKTableNameSetServiceTemplateRelation)
}

func TBHostApplyRule() types.Table {
	return Table(common.BKTableNameHostApplyRule)
}
