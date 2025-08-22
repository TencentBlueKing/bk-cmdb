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

package y3_15_202411071530

import (
	tenanttmp "configcenter/pkg/types/tenant-template"
	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	aPITaskSyncHistoryIndexes = []types.Index{
		{
			Name:       common.CCLogicIndexNamePrefix + "lastTime",
			Keys:       bson.D{{common.LastTimeField, -1}},
			Background: true,
			// delete redundant tasks from 6 months ago
			ExpireAfterSeconds: 6 * 30 * 24 * 60 * 60,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "tenantID_taskID_taskType",
			Keys: bson.D{
				{common.TenantID, 1},
				{common.BKTaskIDField, 1},
				{common.BKTaskTypeField, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "tenantID_instID_taskType_createTime",
			Keys: bson.D{
				{common.TenantID, 1},
				{common.BKInstIDField, 1},
				{common.BKTaskTypeField, 1},
				{common.CreateTimeField, -1},
			},
			Background: true,
		},
	}
	objAttDesIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkObjID",
			Keys:                    bson.D{{"bk_obj_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkObjID_bkPropertyID_bkBizID",
			Keys:                    bson.D{{"bk_obj_id", 1}, {"bk_property_id", 1}, {"bk_biz_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkTemplateID",
			Keys:                    bson.D{{"bk_template_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	apiTaskIndexes = []types.Index{
		{
			Name:       common.CCLogicIndexNamePrefix + "lastTime",
			Keys:       bson.D{{common.LastTimeField, -1}},
			Background: true,
			// delete redundant tasks from 6 months ago
			ExpireAfterSeconds: 6 * 30 * 24 * 60 * 60,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "taskType_status_createTime",
			Keys: bson.D{
				{common.BKTaskTypeField, 1},
				{common.BKStatusField, 1},
				{common.CreateTimeField, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "tenantID_taskType_instID_extra",
			Keys: bson.D{
				{common.TenantID, 1},
				{common.BKTaskTypeField, 1},
				{common.BKInstIDField, 1},
				{metadata.APITaskExtraField, 1},
			},
			Background: true,
			Unique:     true,
			PartialFilterExpression: map[string]interface{}{
				common.BKStatusField: map[string]interface{}{
					common.BKDBIN: []metadata.APITaskStatus{metadata.APITaskStatusNew,
						metadata.APITaskStatusWaitExecute,
						metadata.APITaskStatusExecute},
				},
			},
		},
		{
			Name: common.CCLogicIndexNamePrefix + "tenantID_instID_taskType_createTime",
			Keys: bson.D{
				{common.TenantID, 1},
				{common.BKInstIDField, 1},
				{common.BKTaskTypeField, 1},
				{common.CreateTimeField, -1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "taskID",
			Keys: bson.D{{
				"task_id", 1},
			},
			Unique:     true,
			Background: true,
		},
	}
	applicationBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "default",
			Keys:                    bson.D{{"default", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizName",
			Keys:                    bson.D{{"bk_biz_name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: map[string]interface{}{"bk_biz_name": map[string]string{common.BKDBType: "string"}},
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	asstDesIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkAsstID",
			Keys:                    bson.D{{"bk_asst_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	auditLogIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "operationTime",
			Keys:                    bson.D{{"operation_time", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "user",
			Keys:                    bson.D{{"user", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "resourceName",
			Keys:                    bson.D{{"resource_name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name: common.CCLogicIndexNamePrefix + "auditType_resourceName_action_operationTime",
			Keys: bson.D{{"audit_type", 1}, {"resource_type", 1}, {"action", 1},
				{"operation_time", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	bizSetBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizSetID",
			Keys:                    bson.D{{"bk_biz_set_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizSetName",
			Keys:                    bson.D{{"bk_biz_set_name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizSetID_bkBizSetName",
			Keys:                    bson.D{{"bk_biz_set_id", 1}, {"bk_biz_set_name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	clusterBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "UID",
			Keys:                    bson.D{{"uid", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID_name",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "xid",
			Keys:                    bson.D{{"xid", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	containerBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkPodID_containerUID",
			Keys:                    bson.D{{"bk_pod_id", 1}, {"container_uid", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkPodID",
			Keys:                    bson.D{{"bk_pod_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkClusterID",
			Keys:                    bson.D{{"bk_cluster_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkNamespaceID",
			Keys:                    bson.D{{"bk_namespace_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "refID_refKind",
			Keys:                    bson.D{{"ref.id", 1}, {"ref.kind", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	cronJobBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkNamespaceID_name",
			Keys:                    bson.D{{"bk_namespace_id", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "clusterUID",
			Keys:                    bson.D{{"cluster_uid", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkClusterID",
			Keys:                    bson.D{{"bk_cluster_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "name",
			Keys:                    bson.D{{"name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	daemonSetBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkNamespaceID_name",
			Keys:                    bson.D{{"bk_namespace_id", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "clusterUID",
			Keys:                    bson.D{{"cluster_uid", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkClusterID",
			Keys:                    bson.D{{"bk_cluster_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "name",
			Keys:                    bson.D{{"name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	deploymentBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkNamespaceID_name",
			Keys:                    bson.D{{"bk_namespace_id", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "clusterUID",
			Keys:                    bson.D{{"cluster_uid", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkClusterID",
			Keys:                    bson.D{{"bk_cluster_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "name",
			Keys:                    bson.D{{"name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	dynamicGroupIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID_ID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizID_name",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	fieldTemplateIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "name",
			Keys:                    bson.D{{"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	fullSyncCondIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "resource_subResource",
			Keys:                    bson.D{{"resource", 1}, {"sub_resource", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: map[string]interface{}{"is_all": true},
		},
	}
	gameDeploymentBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkNamespaceID_name",
			Keys:                    bson.D{{"bk_namespace_id", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "clusterUID",
			Keys:                    bson.D{{"cluster_uid", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkClusterID",
			Keys:                    bson.D{{"bk_cluster_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "name",
			Keys:                    bson.D{{"name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	gameStatefulSetBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkNamespaceID_name",
			Keys:                    bson.D{{"bk_namespace_id", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "clusterUID",
			Keys:                    bson.D{{"cluster_uid", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkClusterID",
			Keys:                    bson.D{{"bk_cluster_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "name",
			Keys:                    bson.D{{"name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	hostApplyRuleIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Background:              false,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              false,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkModuleID",
			Keys:                    bson.D{{"bk_module_id", 1}},
			Background:              false,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkBizID_bkModuleID_serviceTemplateID_bkAttributeID",
			Keys: bson.D{{"bk_biz_id", 1}, {"bk_module_id", 1}, {"service_template_id", 1},
				{"bk_attribute_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "serviceTemplateID_bkAttributeID",
			Keys:                    bson.D{{"service_template_id", 1}, {"bk_attribute_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID_serviceTemplateID_bkAttributeID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"service_template_id", 1}, {"bk_attribute_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID_bkModuleID_bkAttributeID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"bk_module_id", 1}, {"bk_attribute_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkModuleID_bkAttributeID",
			Keys:                    bson.D{{"bk_module_id", 1}, {"bk_attribute_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	hostBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkHostName",
			Keys:                    bson.D{{"bk_host_name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkHostInnerIP",
			Keys:                    bson.D{{"bk_host_innerip", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkCloudInstID",
			Keys:                    bson.D{{"bk_cloud_inst_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkCloudID",
			Keys:                    bson.D{{"bk_cloud_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkOsType",
			Keys:                    bson.D{{"bk_os_type", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkHostOuterIP",
			Keys:                    bson.D{{"bk_host_outerip", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: map[string]interface{}{"bk_host_outerip": map[string]string{common.BKDBType: "string"}},
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + "bkCloudInstID_bkCloudVendor",
			Keys:       bson.D{{"bk_cloud_inst_id", 1}, {"bk_cloud_vendor", 1}},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				"bk_cloud_inst_id": map[string]string{common.BKDBType: "string"},
				"bk_cloud_vendor":  map[string]string{common.BKDBType: "string"},
			},
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkHostID",
			Keys:                    bson.D{{"bk_host_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkAssetID",
			Keys:                    bson.D{{"bk_asset_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + "bkAgentID",
			Keys:       bson.D{{"bk_agent_id", 1}},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{"bk_agent_id": map[string]string{common.BKDBType: "string",
				common.BKDBGT: ""}},
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + "bkHostInnerIP_bkCloudID",
			Keys:       bson.D{{"bk_host_innerip", 1}, {"bk_cloud_id", 1}},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				"bk_host_innerip": map[string]string{common.BKDBType: "string"},
				"bk_cloud_id":     map[string]string{common.BKDBType: "number"},
				"bk_addressing":   "static",
			},
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + "bkHostInnerIP_v6_bkCloudID",
			Keys:       bson.D{{"bk_host_innerip_v6", 1}, {"bk_cloud_id", 1}},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				"bk_cloud_id":        map[string]string{common.BKDBType: "number"},
				"bk_addressing":      "static",
				"bk_host_innerip_v6": map[string]string{common.BKDBType: "string"},
			},
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + "bkCloudID_bkHostInnerIP",
			Keys:       bson.D{{"bk_cloud_id", 1}, {"bk_host_innerip", 1}},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				"bk_cloud_id":     map[string]string{common.BKDBType: "number"},
				"bk_host_innerip": map[string]string{common.BKDBType: "string"},
				"bk_addressing":   "static",
			},
		},
	}
	hostLockIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkHostID",
			Keys:                    bson.D{{"bk_host_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	instAsstIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkObjId_bkInstId",
			Keys:                    bson.D{{"bk_obj_id", 1}, {"bk_inst_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", -1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkObjID_bkAsstObjID_bkAsstID",
			Keys:                    bson.D{{"bk_obj_id", -1}, {"bk_asst_obj_id", -1}, {"bk_asst_id", -1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkObjAsstID_ID",
			Keys:                    bson.D{{"bk_obj_asst_id", -1}, {"id", -1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID_",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkAsstObjID_bkAsstInstID",
			Keys:                    bson.D{{"bk_asst_obj_id", 1}, {"bk_asst_inst_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	jobBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkNamespaceID_name",
			Keys:                    bson.D{{"bk_namespace_id", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "clusterUID",
			Keys:                    bson.D{{"cluster_uid", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkClusterID",
			Keys:                    bson.D{{"bk_cluster_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "name",
			Keys:                    bson.D{{"name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	objAsstIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkObjID",
			Keys:                    bson.D{{"bk_obj_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkAsstObjID",
			Keys:                    bson.D{{"bk_asst_obj_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkObjID_bkAsstObjID_bkAsstID",
			Keys:                    bson.D{{"bk_obj_id", 1}, {"bk_asst_obj_id", 1}, {"bk_asst_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	projectBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkProjectID",
			Keys:                    bson.D{{"bk_project_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkProjectName",
			Keys:                    bson.D{{"bk_project_name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkProjectCode",
			Keys:                    bson.D{{"bk_project_code", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkStatus",
			Keys:                    bson.D{{"bk_status", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	platBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkVpcID",
			Keys:                    bson.D{{"bk_vpc_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkCloudID",
			Keys:                    bson.D{{"bk_cloud_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkCloudName",
			Keys:                    bson.D{{"bk_cloud_name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: map[string]interface{}{"bk_cloud_name": map[string]string{common.BKDBType: "string"}},
		},
	}
	setTemplateIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizID_name",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	processInstanceRelationIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "serviceInstanceID",
			Keys:                    bson.D{{"service_instance_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "processTemplateID",
			Keys:                    bson.D{{"process_template_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkProcessID",
			Keys:                    bson.D{{"bk_process_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "serviceInstanceID_bkProcessID",
			Keys:                    bson.D{{"service_instance_id", 1}, {"bk_process_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkProcessID_bkHostID",
			Keys:                    bson.D{{"bk_process_id", 1}, {"bk_host_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkHostID",
			Keys:                    bson.D{{"bk_host_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	objectBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkObjID",
			Keys:                    bson.D{{"bk_obj_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkInstID",
			Keys:                    bson.D{{"bk_inst_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkInstName",
			Keys:                    bson.D{{"bk_inst_name", 1}},
			Background:              false,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	propertyGroupIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkObjID",
			Keys:                    bson.D{{"bk_obj_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkGroupID",
			Keys:                    bson.D{{"bk_group_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizID_bkGroupName_bkObjID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"bk_group_name", 1}, {"bk_obj_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkObjID_bkBizID_bkGroupIndex",
			Keys:                    bson.D{{"bk_obj_id", 1}, {"bk_biz_id", 1}, {"bk_group_index", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkObjID_bkBizID_bkGroupID",
			Keys:                    bson.D{{"bk_obj_id", 1}, {"bk_biz_id", 1}, {"bk_group_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	objAttDesTemplateIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkTemplateID_bkPropertyID",
			Keys:                    bson.D{{"bk_template_id", 1}, {"bk_property_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkTemplateID_bkPropertyName",
			Keys:                    bson.D{{"bk_template_id", 1}, {"bk_property_name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	serviceCategoryIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "name_bkParentID_bkBizID",
			Keys:                    bson.D{{"name", 1}, {"bk_parent_id", 1}, {"bk_biz_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	statefulSetBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkNamespaceID_name",
			Keys:                    bson.D{{"bk_namespace_id", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "clusterUID",
			Keys:                    bson.D{{"cluster_uid", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkClusterID",
			Keys:                    bson.D{{"bk_cluster_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "name",
			Keys:                    bson.D{{"name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	moduleBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkModuleName",
			Keys:                    bson.D{{"bk_module_name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "default",
			Keys:                    bson.D{{"default", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkSetID",
			Keys:                    bson.D{{"bk_set_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkParentID",
			Keys:                    bson.D{{"bk_parent_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + "bkBizID_bkSetID_bkModuleName",
			Keys:       bson.D{{"bk_biz_id", 1}, {"bk_set_id", 1}, {"bk_module_name", 1}},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{"bk_biz_id": map[string]string{common.BKDBType: "number"},
				"bk_set_id":      map[string]string{common.BKDBType: "number"},
				"bk_module_name": map[string]string{common.BKDBType: "string"}},
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkModuleID_bkBizID",
			Keys:                    bson.D{{"bk_module_id", 1}, {"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkModuleID",
			Keys:                    bson.D{{"bk_module_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "setTemplateID_serviceTemplateID",
			Keys:                    bson.D{{"set_template_id", 1}, {"service_template_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "setTemplateID",
			Keys:                    bson.D{{"set_template_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "serviceTemplateID",
			Keys:                    bson.D{{"service_template_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	serviceTemplateAttrIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizID_serviceTemplateID_bkAttributeID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"service_template_id", 1}, {"bk_attribute_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	setServiceTemplateRelationIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "setTemplateID_serviceTemplateID",
			Keys:                    bson.D{{"set_template_id", 1}, {"service_template_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizID_setTemplateID_serviceTemplateID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"set_template_id", 1}, {"service_template_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	processTemplateIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "serviceTemplateID",
			Keys:                    bson.D{{"service_template_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "serviceTemplateID_bkProcessName",
			Keys:                    bson.D{{"service_template_id", 1}, {"bk_process_name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	setTemplateAttrIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizID_setTemplateID_bkAttributeID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"set_template_id", 1}, {"bk_attribute_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	objectUniqueTemplateIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkTemplateID",
			Keys:                    bson.D{{"bk_template_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	nsSharedClusterRelationIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkNamespaceID",
			Keys:                    bson.D{{"bk_namespace_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkAsstBizID",
			Keys:                    bson.D{{"bk_asst_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	modelQuoteRelationIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "destModel",
			Keys:                    bson.D{{"dest_model", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "srcModel_bkPropertyID",
			Keys:                    bson.D{{"src_model", 1}, {"bk_property_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "srcModel",
			Keys:                    bson.D{{"src_model", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	objectBaseMappingIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkInstID",
			Keys:                    bson.D{{"bk_inst_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	moduleHostConfigIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkHostID",
			Keys:                    bson.D{{"bk_host_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkModuleID",
			Keys:                    bson.D{{"bk_module_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkSetID",
			Keys:                    bson.D{{"bk_set_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkModuleID_bkHostID",
			Keys:                    bson.D{{"bk_module_id", 1}, {"bk_host_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID_bkHostID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"bk_host_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkModuleID_bkBizID",
			Keys:                    bson.D{{"bk_module_id", 1}, {"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkSetID_bkBizID",
			Keys:                    bson.D{{"bk_set_id", 1}, {"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	topoGraphicsIndexes = []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "scopeType_scopeID_nodeType_bkObjID_bkInstID",
			Keys: bson.D{{"scope_type", 1}, {"scope_id", 1}, {"node_type", 1}, {"bk_obj_id", 1},
				{"bk_inst_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	serviceInstanceIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "serviceTemplateID",
			Keys:                    bson.D{{"service_template_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkModuleID",
			Keys:                    bson.D{{"bk_module_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID_ID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkModuleID_bkBizID",
			Keys:                    bson.D{{"bk_module_id", 1}, {"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID_bkHostID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"bk_host_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkHostID",
			Keys:                    bson.D{{"bk_host_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	podWorkloadBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkNamespaceID_name",
			Keys:                    bson.D{{"bk_namespace_id", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "clusterUID",
			Keys:                    bson.D{{"cluster_uid", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkClusterID",
			Keys:                    bson.D{{"bk_cluster_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "name",
			Keys:                    bson.D{{"name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	objectUniqueIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkObjID",
			Keys:                    bson.D{{"bk_obj_id", 1}},
			Background:              false,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkTemplateID",
			Keys:                    bson.D{{"bk_template_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	nodeBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkClusterID_name",
			Keys:                    bson.D{{"bk_cluster_id", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID_clusterUID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"cluster_uid", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID_bkClusterID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"bk_cluster_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID_bkHostID",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"bk_host_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID_name",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	podBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "refID_refKind_name",
			Keys:                    bson.D{{"ref.id", 1}, {"ref.kind", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkHostID",
			Keys:                    bson.D{{"bk_host_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "refName_refID",
			Keys:                    bson.D{{"ref.name", 1}, {"ref.id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkClusterID",
			Keys:                    bson.D{{"bk_cluster_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "clusterUID",
			Keys:                    bson.D{{"cluster_uid", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "name",
			Keys:                    bson.D{{"name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkNodeID",
			Keys:                    bson.D{{"bk_node_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkNamespaceID",
			Keys:                    bson.D{{"bk_namespace_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "refID_refKind",
			Keys:                    bson.D{{"ref.id", 1}, {"ref.kind", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	setBaseIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkParentID",
			Keys:                    bson.D{{"bk_parent_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkSetName",
			Keys:                    bson.D{{"bk_set_name", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkParentID_bkSetName",
			Keys:                    bson.D{{"bk_parent_id", 1}, {"bk_set_name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + "bkBizID_bkSetName_bkParentID",
			Keys:       bson.D{{"bk_biz_id", 1}, {"bk_set_name", 1}, {"bk_parent_id", 1}},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				"bk_set_name":  map[string]string{common.BKDBType: "string"},
				"bk_parent_id": map[string]string{common.BKDBType: "number"},
				"bk_biz_id":    map[string]string{common.BKDBType: "number"}},
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkSetID_bkBizID",
			Keys:                    bson.D{{"bk_set_id", 1}, {"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkSetID",
			Keys:                    bson.D{{"bk_set_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	processIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkProcessID",
			Keys:                    bson.D{{"bk_process_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + "serviceInstanceID_bkProcessName",
			Keys:       bson.D{{"service_instance_id", 1}, {"bk_process_name", 1}},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{"service_instance_id": map[string]string{common.BKDBType: "number"},
				"bk_process_name": map[string]string{common.BKDBType: "string"}},
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "serviceInstanceID_bkFuncName_bkStartParamRegex",
			Keys: bson.D{{"service_instance_id", 1}, {"bk_func_name", 1},
				{"bk_start_param_regex", 1}},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{"service_instance_id": map[string]string{common.BKDBType: "number"},
				"bk_func_name":         map[string]string{common.BKDBType: "string"},
				"bk_start_param_regex": map[string]string{common.BKDBType: "string"}},
		},
	}
	objDesIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkClassificationID",
			Keys:                    bson.D{{"bk_classification_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkObjID",
			Keys:                    bson.D{{"bk_obj_id", 1}},
			Unique:                  true,
			Background:              false,
			PartialFilterExpression: map[string]interface{}{"bk_obj_id": map[string]string{common.BKDBType: "string"}},
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkObjName",
			Keys:                    bson.D{{"bk_obj_name", 1}},
			Unique:                  true,
			Background:              false,
			PartialFilterExpression: map[string]interface{}{"bk_obj_name": map[string]string{common.BKDBType: "string"}},
		},
	}
	serviceTemplateIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
			Keys:                    bson.D{{"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizID_name",
			Keys:                    bson.D{{"bk_biz_id", 1}, {"name", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "ID_bkBizID",
			Keys:                    bson.D{{"id", 1}, {"bk_biz_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	objClassificationIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkClassificationID",
			Keys:                    bson.D{{"bk_classification_id", 1}},
			Unique:                  true,
			Background:              false,
			PartialFilterExpression: map[string]interface{}{"bk_classification_id": map[string]string{common.BKDBType: "string"}},
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkClassificationName",
			Keys:                    bson.D{{"bk_classification_name", 1}},
			Unique:                  true,
			Background:              false,
			PartialFilterExpression: map[string]interface{}{"bk_classification_name": map[string]string{common.BKDBType: "string"}},
		},
	}
	objFieldTemplateRelationIndexes = []types.Index{
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkTemplateID_objectID",
			Keys:                    bson.D{{"bk_template_id", 1}, {"object_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "objectID",
			Keys:                    bson.D{{"object_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	namespaceBaseIndexes = []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys: bson.D{
				{"id", 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkClusterID_name",
			Keys: bson.D{
				{"bk_cluster_id", 1},
				{"name", 1},
			},
			Background: true,
			Unique:     true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "clusterUID",
			Keys: bson.D{
				{"cluster_uid", 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "clusterID",
			Keys: bson.D{
				{"bk_cluster_id", 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "name",
			Keys: bson.D{
				{"name", 1},
			},
			Background: true,
		},
	}
	instAsstCommonIndexes = []types.Index{
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkObjID_bkInstID",
			Keys:                    bson.D{{"bk_obj_id", 1}, {"bk_inst_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
			Keys:                    bson.D{{"id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkAsstObjID_bkAsstInstID",
			Keys:                    bson.D{{"bk_asst_obj_id", 1}, {"bk_asst_inst_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicIndexNamePrefix + "bkAsstID",
			Keys:                    bson.D{{"bk_asst_id", 1}},
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
		{
			Name:                    common.CCLogicUniqueIdxNamePrefix + "bkInstID_bkAsstInstID_bkObjAsstID",
			Keys:                    bson.D{{"bk_inst_id", 1}, {"bk_asst_inst_id", 1}, {"bk_obj_asst_id", 1}},
			Unique:                  true,
			Background:              true,
			PartialFilterExpression: make(map[string]interface{}),
		},
	}
	associationDefaultIndexes = []types.Index{
		{
			Name: common.CCLogicIndexNamePrefix + "bkObjId_bkInstID",
			Keys: bson.D{
				{"bk_obj_id", 1},
				{"bk_inst_id", 1},
			},
			Background: true,
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + "id",
			Keys:       bson.D{{"id", 1}},
			Unique:     true,
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkAsstObjId_bkAsstInstId",
			Keys: bson.D{
				{"bk_asst_obj_id", 1},
				{"bk_asst_inst_id", 1},
			},
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + "bkAsstID",
			Keys:       bson.D{{"bk_asst_id", 1}},
			Background: true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkInstID_bkAsstInstID_bkObjAsstID",
			Keys: bson.D{
				{common.BKInstIDField, 1},
				{common.BKAsstInstIDField, 1},
				{common.AssociationObjAsstIDField, 1},
			},
			Unique:     true,
			Background: true,
		},
	}
	tenantIndexes = []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "tenantID",
			Keys: bson.D{
				{common.TenantID, 1},
			},
			Unique:     true,
			Background: true,
		},
	}
	globalConfigIndexes = []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "tenantID",
			Keys: bson.D{
				{common.TenantID, 1},
			},
			Unique:     true,
			Background: true,
		},
	}
	defaultAreaHostIndexes = []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkCloudID_bkHostInnerIP",
			Keys: bson.D{
				{"bk_cloud_id", 1},
				{"bk_host_innerip", 1},
			},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				"bk_cloud_id":     map[string]string{common.BKDBType: "number"},
				"bk_host_innerip": map[string]string{common.BKDBType: "string"},
			},
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkCloudID_bkHostInnerIPV6",
			Keys: bson.D{
				{"bk_cloud_id", 1},
				{"bk_host_innerip_v6", 1},
			},
			Unique:     true,
			Background: true,
			PartialFilterExpression: map[string]interface{}{
				"bk_cloud_id":        map[string]string{common.BKDBType: "number"},
				"bk_host_innerip_v6": map[string]string{common.BKDBType: "string"},
			},
		},
		{
			Name: common.CCLogicIndexNamePrefix + "tenantID",
			Keys: bson.D{
				{common.TenantID, 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkHostID",
			Keys: bson.D{
				{"bk_host_id", 1},
			},
			Unique:     true,
			Background: true,
		},
	}
)

var templateIndexes = []types.Index{
	{
		Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
		Keys:                    bson.D{{common.BKFieldID, 1}},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + string(tenanttmp.TemplateTypeObjAttribute) +
			"_bkObjID_bkPropertyID",
		Keys:       bson.D{{"data.bk_obj_id", 1}, {"data.bk_property_id", 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			tenanttmp.BKTenantTemplateTypeField: tenanttmp.TemplateTypeObjAttribute,
		},
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + string(tenanttmp.TemplateTypePropertyGroup) +
			"_bkGroupName_bkObjID",
		Keys:       bson.D{{"data.bk_group_name", 1}, {"data.bk_obj_id", 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			tenanttmp.BKTenantTemplateTypeField: tenanttmp.TemplateTypePropertyGroup,
		},
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + string(tenanttmp.TemplateTypePropertyGroup) +
			"_bkObjID_bkGroupIndex",
		Keys:       bson.D{{"data.bk_obj_id", 1}, {"data.bk_group_index", 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			tenanttmp.BKTenantTemplateTypeField: tenanttmp.TemplateTypePropertyGroup,
		},
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + string(tenanttmp.TemplateTypeObjClassification) +
			"_bkClassificationID",
		Keys:       bson.D{{"data.bk_classification_id", 1}},
		Unique:     true,
		Background: false,
		PartialFilterExpression: map[string]interface{}{
			tenanttmp.BKTenantTemplateTypeField: tenanttmp.TemplateTypeObjClassification,
			"data.bk_classification_id":         map[string]string{common.BKDBType: "string"}},
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + string(tenanttmp.TemplateTypeObjClassification) +
			"_bkClassificationName",
		Keys:       bson.D{{"data.bk_classification_name", 1}},
		Unique:     true,
		Background: false,
		PartialFilterExpression: map[string]interface{}{
			tenanttmp.BKTenantTemplateTypeField: tenanttmp.TemplateTypeObjClassification,
			"data.bk_classification_name":       map[string]string{common.BKDBType: "string"}},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + string(tenanttmp.TemplateTypeObject) + "_bkObjID",
		Keys:       bson.D{{"data.bk_obj_id", 1}},
		Unique:     true,
		Background: false,
		PartialFilterExpression: map[string]interface{}{
			tenanttmp.BKTenantTemplateTypeField: tenanttmp.TemplateTypeObject,
			"data.bk_obj_id":                    map[string]string{common.BKDBType: "string"}},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + string(tenanttmp.TemplateTypeObject) + "_bkObjName",
		Keys:       bson.D{{"data.bk_obj_name", 1}},
		Unique:     true,
		Background: false,
		PartialFilterExpression: map[string]interface{}{
			tenanttmp.BKTenantTemplateTypeField: tenanttmp.TemplateTypeObject,
			"data.bk_obj_name":                  map[string]string{common.BKDBType: "string"}},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + string(tenanttmp.TemplateTypeBizSet) + "_bkBizSetName",
		Keys:       bson.D{{"data.bk_biz_set_name", 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			tenanttmp.BKTenantTemplateTypeField: tenanttmp.TemplateTypeBizSet,
		},
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + string(tenanttmp.TemplateTypeServiceCategory) + "_nameParentName",
		Keys:       bson.D{{"data.name", 1}, {"data.parent_name", 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			tenanttmp.BKTenantTemplateTypeField: tenanttmp.TemplateTypeServiceCategory,
		},
	},
}
