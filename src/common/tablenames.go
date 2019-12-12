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

package common

// table names
const (
	// BKTableNameProcModule the table name of the process module
	BKTableNameProcModule = "cc_Proc2Module"

	// BKTableNameProcTempalte the table name of the process template config
	BKTableNameProcTemplate = "cc_Proc2Template"

	// BKTableNameProcInstanceModel the table name of the process instance
	BKTableNameProcInstanceModel = "cc_ProcInstanceModel"

	// BKTableNameProcInstaceDetail  the table name of the process instance detail
	BKTableNameProcInstaceDetail = "cc_ProcInstanceDetail"

	// BKTableNameProcOperateTask  the table name of the process instance operater task info
	BKTableNameProcOperateTask = "cc_ProcOpTask"

	// BKTableNamePropertyGroup the table name of the property group
	BKTableNamePropertyGroup = "cc_PropertyGroup"

	// BKTableNameObjDes the table name of the asst des
	BKTableNameAsstDes = "cc_AsstDes"

	// BKTableNameObjDes the table name of the object
	BKTableNameObjDes = "cc_ObjDes"

	// BKTableNameObjDes the table name of the object
	BKTableNameObjUnique = "cc_ObjectUnique"

	// BKTableNameObjAttDes the table name of the object attribute
	BKTableNameObjAttDes = "cc_ObjAttDes"

	// BKTableNameObjClassifiction the table name of the object classification
	BKTableNameObjClassifiction = "cc_ObjClassification"

	// BKTableNameInstAsst the table name of the inst association
	BKTableNameInstAsst = "cc_InstAsst"

	BKTableNameBaseApp     = "cc_ApplicationBase"
	BKTableNameBaseHost    = "cc_HostBase"
	BKTableNameBaseModule  = "cc_ModuleBase"
	BKTableNameBaseInst    = "cc_ObjectBase"
	BKTableNameBasePlat    = "cc_PlatBase"
	BKTableNameBaseSet     = "cc_SetBase"
	BKTableNameBaseProcess = "cc_Process"

	BKTableNameModuleHostConfig = "cc_ModuleHostConfig"
	BKTableNameSystem           = "cc_System"
	BKTableNameHistory          = "cc_History"
	BKTableNameHostFavorite     = "cc_HostFavourite"
	BKTableNameOperationLog     = "cc_OperationLog"
	BKTableNameSubscription     = "cc_Subscription"
	BKTableNameUserAPI          = "cc_UserAPI"
	BKTableNameUserCustom       = "cc_UserCustom"
	BKTableNameObjAsst          = "cc_ObjAsst"
	BKTableNameTopoGraphics     = "cc_TopoGraphics"
	BKTableNameTransaction      = "cc_Transaction"
	BKTableNameIDgenerator      = "cc_idgenerator"

	BKTableNameNetcollectDevice   = "cc_NetcollectDevice"
	BKTableNameNetcollectProperty = "cc_NetcollectProperty"

	BKTableNameNetcollectConfig  = "cc_NetcollectConfig"
	BKTableNameNetcollectReport  = "cc_NetcollectReport"
	BKTableNameNetcollectHistory = "cc_NetcollectHistory"

	BKTableNameHostLock = "cc_HostLock"

	// Cloud sync tables
	BKTableNameCloudTask              = "cc_CloudTask"
	BKTableNameCloudSyncHistory       = "cc_CloudSyncHistory"
	BKTableNameCloudResourceConfirm   = "cc_CloudResourceConfirm"
	BKTableNameResourceConfirmHistory = "cc_ResourceConfirmHistory"

	// Operation tables
	BKTableNameChartConfig   = "cc_ChartConfig"
	BKTableNameChartPosition = "cc_ChartPosition"
	BKTableNameChartData     = "cc_ChartData"

	// process tables
	BKTableNameServiceCategory         = "cc_ServiceCategory"
	BKTableNameServiceTemplate         = "cc_ServiceTemplate"
	BKTableNameServiceInstance         = "cc_ServiceInstance"
	BKTableNameProcessTemplate         = "cc_ProcessTemplate"
	BKTableNameProcessInstanceRelation = "cc_ProcessInstanceRelation"

	BKTableNameSetTemplate                = "cc_SetTemplate"
	BKTableNameSetServiceTemplateRelation = "cc_SetServiceTemplateRelation"
	BKTableNameAPITask                    = "cc_APITask"
	BKTableNameSetTemplateSyncStatus      = "cc_SetTemplateSyncStatus"
	BKTableNameSetTemplateSyncHistory     = "cc_SetTemplateSyncHistory"

	// rule for host property auto apply
	BKTableNameHostApplyRule = "cc_HostApplyRule"
)

// AllTables alltables
var AllTables = []string{
	BKTableNameProcModule,
	BKTableNameProcTemplate,
	BKTableNameProcInstanceModel,
	BKTableNameProcInstaceDetail,
	BKTableNameProcOperateTask,
	BKTableNamePropertyGroup,
	BKTableNameObjDes,
	BKTableNameObjAttDes,
	BKTableNameObjClassifiction,
	BKTableNameInstAsst,
	BKTableNameBaseApp,
	BKTableNameBaseHost,
	BKTableNameBaseModule,
	BKTableNameBaseInst,
	BKTableNameBasePlat,
	BKTableNameBaseSet,
	BKTableNameBaseProcess,
	BKTableNameModuleHostConfig,
	BKTableNameSystem,
	BKTableNameHistory,
	BKTableNameHostFavorite,
	BKTableNameOperationLog,
	BKTableNameSubscription,
	BKTableNameUserAPI,
	BKTableNameUserCustom,
	BKTableNameObjAsst,
	BKTableNameTopoGraphics,
	BKTableNameNetcollectConfig,
	BKTableNameNetcollectDevice,
	BKTableNameNetcollectProperty,
	BKTableNameNetcollectReport,
	BKTableNameNetcollectHistory,
	BKTableNameTransaction,
	BKTableNameIDgenerator,
	BKTableNameHostLock,
	BKTableNameCloudTask,
	BKTableNameCloudSyncHistory,
	BKTableNameCloudResourceConfirm,
	BKTableNameResourceConfirmHistory,
	BKTableNameObjUnique,
	BKTableNameAsstDes,
	BKTableNameServiceCategory,
	BKTableNameServiceTemplate,
	BKTableNameServiceInstance,
	BKTableNameProcessTemplate,
	BKTableNameProcessInstanceRelation,
	BKTableNameSetTemplate,
	BKTableNameSetServiceTemplateRelation,
	BKTableNameChartConfig,
	BKTableNameChartPosition,
	BKTableNameChartData,
	BKTableNameHostApplyRule,
}

// GetInstTableName returns inst data table name
func GetInstTableName(objID string) string {
	switch objID {
	case BKInnerObjIDApp:
		return BKTableNameBaseApp
	case BKInnerObjIDSet:
		return BKTableNameBaseSet
	case BKInnerObjIDModule:
		return BKTableNameBaseModule
	// case BKInnerObjIDObject:
	// 	return BKTableNameBaseInst
	case BKInnerObjIDHost:
		return BKTableNameBaseHost
	case BKInnerObjIDProc:
		return BKTableNameBaseProcess
	case BKInnerObjIDPlat:
		return BKTableNameBasePlat
	case BKTableNameInstAsst:
		return BKTableNameInstAsst
	case BKTableNameModuleHostConfig:
		return BKTableNameModuleHostConfig
	case BKTableNameObjClassifiction:
		return BKTableNameObjClassifiction
	case BKTableNameObjAttDes:
		return BKTableNameObjAttDes
	case BKTableNamePropertyGroup:
		return BKTableNamePropertyGroup
	case BKTableNameObjUnique:
		return BKTableNameObjUnique
	case BKTableNameAsstDes:
		return BKTableNameAsstDes
	case BKTableNameOperationLog:
		return BKTableNameOperationLog
	case BKTableNameUserAPI:
		return BKTableNameUserAPI
	default:
		return BKTableNameBaseInst
	}
}
