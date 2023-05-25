/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

export const MENU_ENTRY = 'menu_entry'
export const MENU_INDEX = 'menu_index'
export const MENU_BUSINESS = 'menu_business'
export const MENU_RESOURCE = 'menu_resource'
export const MENU_MODEL = 'menu_model'
export const MENU_ANALYSIS = 'menu_analysis'
export const MENU_ADMIN = 'menu_admin'
export const MENU_PLATFORM_MANAGEMENT = 'menu_platform_management'
export const MENU_BUSINESS_SET = 'menu_business_set'

/**
 * 业务集消费视图
 */
export const MENU_BUSINESS_SET_TOPOLOGY = 'menu_business_set_topology'

/**
 * 业务
 */
export const MENU_BUSINESS_HOST = 'menu_business_host'
export const MENU_BUSINESS_HOST_MANAGEMENT = 'menu_business_host_management'
export const MENU_BUSINESS_SERVICE = 'menu_business_service'
export const MENU_BUSINESS_SERVICE_TOPOLOGY = 'menu_business_service_topology'
export const MENU_BUSINESS_ADVANCED = 'menu_business_advanced'

export const MENU_BUSINESS_HOST_AND_SERVICE = 'menu_business_host_and_service'

// 服务模板
export const MENU_BUSINESS_SERVICE_TEMPLATE = 'menu_business_service_template'
export const MENU_BUSINESS_SERVICE_TEMPLATE_CREATE = 'menu_business_service_template_create'
export const MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS = 'menu_business_service_template_details'
export const MENU_BUSINESS_SERVICE_TEMPLATE_EDIT = 'menu_business_service_template_edit'

// 集群模板
export const MENU_BUSINESS_SET_TEMPLATE = 'menu_business_set_template'
export const MENU_BUSINESS_SET_TEMPLATE_CREATE = 'menu_business_set_template_create'
export const MENU_BUSINESS_SET_TEMPLATE_DETAILS = 'menu_business_set_template_details'
export const MENU_BUSINESS_SET_TEMPLATE_EDIT = 'menu_business_set_template_edit'
export const MENU_BUSINESS_SET_TEMPLATE_SYNC_HISTORY = 'menu_business_set_template_sync_history'

// 主机属性自动应用
export const MENU_BUSINESS_HOST_APPLY = 'menu_business_host_apply'
export const MENU_BUSINESS_HOST_APPLY_EDIT = 'menu_business_host_apply_edit'
export const MENU_BUSINESS_HOST_APPLY_CONFIRM = 'menu_business_host_apply_confirm'
export const MENU_BUSINESS_HOST_APPLY_RUN = 'menu_business_host_apply_run'
export const MENU_BUSINESS_HOST_APPLY_CONFLICT = 'menu_business_host_apply_conflict'
export const MENU_BUSINESS_HOST_APPLY_FAILED = 'menu_business_host_apply_failed'

export const MENU_BUSINESS_SERVICE_CATEGORY = 'menu_business_service_category'
export const MENU_BUSINESS_CUSTOM_QUERY = 'menu_business_custom_query'
export const MENU_BUSINESS_CUSTOM_FIELDS = 'menu_business_custom_fields'

/**
 * 资源
 */
export const MENU_RESOURCE_EVENTPUSH = 'menu_resource_eventpush'
export const MENU_RESOURCE_MANAGEMENT = 'menu_resource_management'
export const MENU_RESOURCE_BUSINESS = 'menu_resource_business'
export const MENU_RESOURCE_BUSINESS_HISTORY = 'menu_resource_business_history'
export const MENU_RESOURCE_BUSINESS_DETAILS = 'menu_resource_business_details'
export const MENU_RESOURCE_HOST = 'menu_resource_host'
export const MENU_RESOURCE_INSTANCE = 'menu_resource_instance'
export const MENU_RESOURCE_INSTANCE_DETAILS = 'menu_resource_instance_details'

/**
 * 项目
 */
export const MENU_RESOURCE_PROJECT = 'menu_resource_project'
export const MENU_RESOURCE_PROJECT_DETAILS = 'menu_resource_project_details'

/**
 * 模型
 */
export const MENU_MODEL_MANAGEMENT = 'menu_model_management'
export const MENU_MODEL_TOPOLOGY = 'menu_model_topology'
export const MENU_MODEL_TOPOLOGY_NEW = 'menu_model_topology_new'
export const MENU_MODEL_BUSINESS_TOPOLOGY = 'menu_model_business_topology'
export const MENU_MODEL_ASSOCIATION = 'menu_model_association'
export const MENU_MODEL_DETAILS = 'menu_model_details'

/**
 * 运营分析
 */
export const MENU_ANALYSIS_AUDIT = 'menu_analysis_audit'
export const MENU_ANALYSIS_OPERATION = 'menu_analysis_operation'
export const MENU_ANALYSIS_STATISTICS = 'menu_analysis_statistics'

/**
 * 平台管理
 */
export const MENU_PLATFORM_MANAGEMENT_GLOBAL_CONFIG = 'menu_platform_management_global_config'

// 判断收藏的目录id
export const MENU_RESOURCE_COLLECTION = 'menu_resource_collection'
export const MENU_RESOURCE_HOST_COLLECTION = 'menu_resource_host_collection'
export const MENU_RESOURCE_BUSINESS_COLLECTION = 'menu_resource_business_collection'
export const MENU_RESOURCE_BUSINESS_SET_COLLECTION = 'menu_resource_business_set_collection'
export const MENU_RESOURCE_PROJECT_COLLECTION = 'menu_resource_project_collection'

// 不同模式下不同资源的主机详情的id
export const MENU_RESOURCE_HOST_DETAILS = 'menu_resource_host_details'
export const MENU_RESOURCE_BUSINESS_HOST_DETAILS = 'menu_resource_business_host_details'
export const MENU_BUSINESS_HOST_DETAILS = 'menu_business_host_details'
export const MENU_BUSINESS_SET_HOST_DETAILS = 'menu_business_set_host_details'

// 转移主机
export const MENU_BUSINESS_TRANSFER_HOST = 'menu_business_transfer_host'

// 删除服务实例
export const MENU_BUSINESS_DELETE_SERVICE = 'menu_business_delete_service'

// 管控区域
export const MENU_RESOURCE_CLOUD_AREA = 'menu_resource_cloud_area'
// 云账户
export const MENU_RESOURCE_CLOUD_ACCOUNT = 'menu_resource_cloud_account'
// 云资源发现
export const MENU_RESOURCE_CLOUD_RESOURCE = 'menu_resource_cloud_resource'

// 业务集资源实例
export const MENU_RESOURCE_BUSINESS_SET = 'menu_resource_business_set'
export const MENU_RESOURCE_BUSINESS_SET_DETAILS = 'menu_resource_set_business_details'

// 业务(集)选择器中收藏key
export const BUSINESS_SELECTOR_COLLECTION = 'business_selector_collection'

// Pod
export const MENU_POD_DETAILS = 'menu_pod_details'
export const MENU_POD_CONTAINER_DETAILS = 'menu_pod_container_details'
