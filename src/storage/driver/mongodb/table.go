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

// TBHostBase TODO
func TBHostBase() types.Table {
	return Table(common.BKTableNameBaseHost)
}

// TBApplicationBase TODO
func TBApplicationBase() types.Table {
	return Table(common.BKTableNameBaseApp)
}

// TBSetBase TODO
func TBSetBase() types.Table {
	return Table(common.BKTableNameBaseSet)
}

// TBModuleBase TODO
func TBModuleBase() types.Table {
	return Table(common.BKTableNameBaseModule)
}

// TBModelInstance 模型实例数据表，TBInstance别名函数
func TBModelInstance(objID, supplierAccount string) types.Table {
	return Table(common.GetObjectInstTableName(objID, supplierAccount))
}

// TBInstance 模型实例数据表, TBModelInstance 别名函数
func TBInstance(objID, supplierAccount string) types.Table {
	return TBModelInstance(objID, supplierAccount)
}

// TBCloudArea TODO
func TBCloudArea() types.Table {
	return Table(common.BKTableNameBasePlat)
}

// TBProcess TODO
func TBProcess() types.Table {
	return Table(common.BKTableNameBaseProcess)
}

// TBPropertyGroup TODO
func TBPropertyGroup() types.Table {
	return Table(common.BKTableNamePropertyGroup)
}

// TBAsstDes 模型关联关系
func TBAsstDes() types.Table {
	return Table(common.BKTableNameAsstDes)
}

// TBObjDesc TODO
func TBObjDesc() types.Table {
	return Table(common.BKTableNameObjDes)
}

// TBObjUnique TODO
func TBObjUnique() types.Table {
	return Table(common.BKTableNameObjUnique)
}

// TBObjAttrDesc TODO
func TBObjAttrDesc() types.Table {
	return Table(common.BKTableNameObjAttDes)
}

// TBObjClassification TODO
func TBObjClassification() types.Table {
	return Table(common.BKTableNameObjClassification)
}

// TBInstanceAsst 实例关系数据， TBInstAsst别名函数
func TBInstanceAsst(objID, supplierAccount string) types.Table {
	return Table(common.GetObjectInstAsstTableName(objID, supplierAccount))
}

// TBInstAsst 实例关系数据， TBInstanceAsst别名函数
func TBInstAsst(objID, supplierAccount string) types.Table {
	return TBInstanceAsst(objID, supplierAccount)
}

// TBModuleHostConfig TODO
func TBModuleHostConfig() types.Table {
	return Table(common.BKTableNameModuleHostConfig)
}

// TBAuditLog TODO
func TBAuditLog() types.Table {
	return Table(common.BKTableNameAuditLog)
}

// TBObjectAsst TODO
func TBObjectAsst() types.Table {
	return Table(common.BKTableNameObjAsst)
}

// TBServiceCategory TODO
func TBServiceCategory() types.Table {
	return Table(common.BKTableNameServiceCategory)
}

// TBServiceTemplate TODO
func TBServiceTemplate() types.Table {
	return Table(common.BKTableNameServiceTemplate)
}

// TBServiceInstance TODO
func TBServiceInstance() types.Table {
	return Table(common.BKTableNameServiceInstance)
}

// TBProcessTemplate TODO
func TBProcessTemplate() types.Table {
	return Table(common.BKTableNameProcessTemplate)
}

// TBProcessInstanceRelation TODO
func TBProcessInstanceRelation() types.Table {
	return Table(common.BKTableNameProcessInstanceRelation)
}

// TBSetTemplate TODO
func TBSetTemplate() types.Table {
	return Table(common.BKTableNameSetTemplate)
}

// TBSetServiceTemplateRelation TODO
func TBSetServiceTemplateRelation() types.Table {
	return Table(common.BKTableNameSetServiceTemplateRelation)
}

// TBHostApplyRule TODO
func TBHostApplyRule() types.Table {
	return Table(common.BKTableNameHostApplyRule)
}
