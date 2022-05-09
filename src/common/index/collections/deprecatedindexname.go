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

package collections

import "configcenter/src/common"

// 2021年03月12日 这里的数据不能删除，在deprecatedIndexName存在的索引名字，没有数据库表中索引存在，就表示这个字段需要被删除
// 例如 index/collections/cc_DynamicGroup中的deprecatedDynamicGroupIndexes中的index存在id_1和id_1_bk_biz_id_1 而下面map中存在
// 的是id_1，id_1_bk_biz_id_1和 bk_biz_id_1_name_1，所以就会将索引 bk_biz_id_1_name_1删除。
var deprecatedIndexName = map[string][]string{
	common.BKTableNameDynamicGroup: {
		"id_1",
		"id_1_bk_biz_id_1",
		"bk_biz_id_1_name_1",
	},
	common.BKTableNameObjDes: {
		"bk_obj_id_1",
		"bk_classification_id_1",
		"bk_obj_name_1",
		"bk_supplier_account_1",
		"idx_unique_id",
	},
	common.BKTableNamePropertyGroup: {
		"bk_obj_id_1",
		"bk_supplier_account_1",
		"bk_group_id_1",
		"idx_unique_id",
		"idx_unique_groupName",
		"idx_unique_groupId",
	},
	common.BKTableNameSetTemplate: {
		"idx_id",
		"idx_unique_bizID_name",
	},
	common.BKTableNameAsstDes: {
		"idx_unique_id",
		"idx_supplierId",
		"idx_asstId",
	},
	common.BKTableNameHostApplyRule: {
		"bk_biz_id",
		"id",
		"bk_module_id",
		"host_property_under_module",
		"idx_unique_bizID_moduleID_attrID",
	},
	common.BKTableNameHostLock: {
		"bk_host_id_1",
	},
	common.BKTableNameBaseSet: {
		"bk_parent_id_1",
		"bk_biz_id_1",
		"bk_supplier_account_1",
		"bk_set_name_1",
		"bk_set_id_1_bk_biz_id_1",
		"idx_unique_setID",
		"bk_set_id_1",
	},
	common.BKTableNameObjAttDes: {
		"bk_obj_id_1",
		"bk_supplier_account_1",
		"idx_unique_objID_propertyID_bizID",
		"idx_unique_Id",
		"id_1",
	},
	common.BKTableNameSetServiceTemplateRelation: {
		"idx_unque_setTemplateID_serviceTemplateID",
	},
	common.BKTableNameBaseHost: {
		"bk_host_name_1",
		"bk_host_innerip_1",
		"bk_host_outerip_1",
		"bk_host_id_1_bk_supplier_account_1",
		"innerIP_platID",
		"bk_supplier_account_1",
		"bk_cloud_id_1",
		"idx_unique_hostID",
		"cloudInstID",
		"bk_idx_bk_asset_id",
		"bk_os_type_1",
	},
	common.BKTableNameBaseModule: {
		"bk_module_name_1",
		"default_1",
		"bk_biz_id_1",
		"bk_supplier_account_1",
		"bk_set_id_1",
		"bk_parent_id_1",
		"bk_module_id_1_bk_biz_id_1",
		"idx_unique_moduleID",
		"bk_idx_set_template_id_service_template_id",
		"bk_idx_set_template_id",
		"bk_idx_service_template_id",
		"bk_module_id_1",
	},
	common.BKTableNameChartConfig: {
		"config_id",
		"bk_obj_id",
	},
	common.BKTableNameCloudSyncTask: {
		"bk_task_id",
	},
	common.BKTableNameBasePlat: {
		"bk_supplier_account_1",
		"vpcID",
		"idx_unique_cloudName",
	},
	common.BKTableNameTopoGraphics: {
		"scope_type_1_scope_id_1_node_type_1_bk_obj_id_1_bk_inst_id_1",
	},
	common.BKTableNameDelArchive: {
		"oid_1",
		"idx_oid_coll",
		"idx_coll",
	},
	common.BKTableNameNetcollectDevice: {
		"device_id_1",
		"device_name_1",
		"bk_supplier_account_1",
	},
	common.BKTableNameServiceInstance: {
		"idx_bkBizID",
		"idx_serviceTemplateID",
		"moduleID",
		"bk_module_id_1_bk_biz_id_1",
		"bk_biz_id_1_bk_host_id_1",
		"idx_unique_id",
		"bk_idx_host_id",
		"idx_id",
	},
	common.BKTableNameProcessInstanceRelation: {
		"idx_bkServiceInstanceID",
		"idx_bkProcessTemplateID",
		"idx_bkBizID",
		"idx_bkProcessID",
		"idx_bkHostID",
		"idx_unique_serviceInstID_ProcID",
		"idx_unique_procID_hostID",
	},
	common.BKTableNameProcessTemplate: {
		"idx_serviceTemplateID",
		"idx_bkBizID",
		"idx_unique_id",
		"bk_idx_service_template_idd_bk_process_name",
		"bk_idx_service_template_id_bk_process_name",
		"idx_id",
	},
	common.BKTableNameCloudAccount: {
		"bk_account_id",
	},
	common.BKTableNameObjAsst: {
		"bk_obj_id_1",
		"bk_asst_obj_id_1",
		"bk_supplier_account_1",
		"bk_obj_id_1_bk_asst_obj_id_1_bk_asst_id_1",
		"idx_unique_id",
	},
	common.BKTableNameChartPosition: {
		"bk_biz_id",
	},
	common.BKTableNameCloudSyncHistory: {
		"bk_history_id",
	},
	common.BKTableNameServiceTemplate: {
		"id_1_bk_biz_id_1",
		"idx_bkBizID",
		"idx_unique_id",
		"bk_idx_bk_biz_id_name",
		"idx_id",
	},
	common.BKTableNameAPITask: {
		"idx_taskID",
		"idx_name_status_createTime",
		"idx_status_lastTime",
		"idx_name_flag_createTime",
	},
	common.BKTableNameNetcollectProperty: {
		"netcollect_property_id_1",
		"bk_supplier_account_1",
	},
	common.BKTableNameBaseProcess: {
		"bk_biz_id_1",
		"bk_supplier_account_1",
		"idx_unique_procID",
	},
	common.BKTableNameObjClassification: {
		"bk_classification_id_1",
		"bk_classification_name_1",
		"idx_unique_id",
	},
	common.BKTableNameObjUnique: {
		"bk_obj_id",
		"idx_unique_id",
	},
	common.BKTableNameModuleHostConfig: {
		"bk_biz_id_1",
		"bk_host_id_1",
		"bk_module_id_1",
		"bk_set_id_1",
		"bk_module_id_1_bk_biz_id_1",
		"bk_set_id_1_bk_biz_id_1",
		"idx_unique_moduleID_hostID",
	},
	common.BKTableNameBaseApp: {
		"bk_biz_name_1",
		"default_1",
		"bk_biz_id_1_bk_supplier_account_1",
		"default_1_bk_supplier_account_1",
		"idx_unique_bizID",
		"bk_biz_id_1",
	},
	common.BKTableNameServiceCategory: {
		"idx_unique_id",
		"idx_unique_Name_parentID_bizID",
	},
	common.BKTableNameAuditLog: {
		"index_id",
		"index_operationTime",
		"index_user",
		"index_resourceName",
		"index_operationTime_auditType_resourceType_action",
	},
}
