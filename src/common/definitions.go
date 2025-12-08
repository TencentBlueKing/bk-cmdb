/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package common

import (
	"math"
	"time"
)

const (
	// HTTPCreate create method
	HTTPCreate = "POST"

	// HTTPSelectPost select method
	HTTPSelectPost = "POST"

	// HTTPSelectGet select method
	HTTPSelectGet = "GET"

	// HTTPUpdate update method
	HTTPUpdate = "PUT"

	// HTTPDelete delete method
	HTTPDelete = "DELETE"

	// BKNoLimit no limit definition
	BKNoLimit = 999999999
	// BKMaxPageSize TODO
	// max limit of a page
	BKMaxPageSize = 1000

	// BKMaxLimitSize max limit of a page.
	BKMaxLimitSize = 500

	// BKMaxUpdateOrCreatePageSize maximum number of updates.
	BKMaxUpdateOrCreatePageSize = 100

	// BKMaxDeletePageSize max limit of a delete page
	BKMaxDeletePageSize = 500

	// BKMaxInstanceLimit TODO
	// max limit of instance count
	BKMaxInstanceLimit = 500

	// BKMaxRecordsAtOnce 一次最大操作记录数
	BKMaxRecordsAtOnce = 2000

	// BKDefaultLimit the default limit definition
	BKDefaultLimit = 20

	// BKAuditLogPageLimit the audit log page limit
	BKAuditLogPageLimit = 200

	// BKMaxExportLimit the limit to export
	BKMaxExportLimit = 500

	// BKInstMaxExportLimit the limit to instance export
	BKInstMaxExportLimit = 1000

	// BKMaxOnceExportLimit the limit once to export
	BKMaxOnceExportLimit = 30000

	// BKInstParentStr the inst parent name
	BKInstParentStr = "bk_parent_id"

	// BKDefaultTenantID the platform operation tenant value
	BKDefaultTenantID = "system"

	// BKSingleTenantID the tenant id while tenant model is off
	BKSingleTenantID = "default"

	// BKSuperTenantID the super tenant value
	BKSuperTenantID = "superadmin"

	// BKDefaultDirSubArea the default dir subarea
	BKDefaultDirSubArea = 0

	// BKTimeTypeParseFlag the time flag
	BKTimeTypeParseFlag = "cc_time_type"

	// BKTopoBusinessLevelLimit the mainline topo level limit
	BKTopoBusinessLevelLimit = "level.businessTopoMax"

	// BKTopoBusinessLevelDefault the mainline topo level default level
	BKTopoBusinessLevelDefault = 7

	// BKMaxSyncIdentifierLimit sync identifier max value
	BKMaxSyncIdentifierLimit = 200

	// BKMaxWriteOpLimit maximum limit of write operation.
	BKMaxWriteOpLimit = 200
	// BKWriteOpLimit default write operation limit
	BKWriteOpLimit = 200
)

const (

	// BKInnerObjIDBizSet the inner object
	BKInnerObjIDBizSet = "bk_biz_set_obj"

	// BKInnerObjIDApp the inner object
	BKInnerObjIDApp = "biz"

	// BKInnerObjIDSet the inner object
	BKInnerObjIDSet = "set"

	// BKInnerObjIDModule the inner object
	BKInnerObjIDModule = "module"

	// BKInnerObjIDHost the inner object
	BKInnerObjIDHost = "host"

	// BKInnerObjIDObject the inner object
	BKInnerObjIDObject = "object"

	// BKInnerObjIDProc the inner object
	BKInnerObjIDProc = "process"

	// BKInnerObjIDConfigTemp the inner object
	BKInnerObjIDConfigTemp = "config_template"

	// BKInnerObjIDTempVersion the inner object
	BKInnerObjIDTempVersion = "template_version"

	// BKInnerObjIDPlat the inner object
	BKInnerObjIDPlat = "plat"

	// BKInnerObjIDProject the inner object
	BKInnerObjIDProject = "bk_project"

	// BKInnerObjIDSwitch the inner object
	BKInnerObjIDSwitch = "bk_switch"
	// BKInnerObjIDRouter the inner object
	BKInnerObjIDRouter = "bk_router"
	// BKInnerObjIDBlance the inner object
	BKInnerObjIDBlance = "bk_load_balance"
	// BKInnerObjIDFirewall the inner object
	BKInnerObjIDFirewall = "bk_firewall"
	// BKInnerObjIDWeblogic the inner object
	BKInnerObjIDWeblogic = "bk_weblogic"
	// BKInnerObjIDTomcat the inner object
	BKInnerObjIDTomcat = "bk_tomcat"
	// BKInnerObjIDApache the inner object
	BKInnerObjIDApache = "bk_apache"
)

// BKInnerObjects is inner objects
var BKInnerObjects = []string{
	BKInnerObjIDApp,
	BKInnerObjIDModule,
	BKInnerObjIDProc,
	BKInnerObjIDHost,
	BKInnerObjIDProject,
	BKInnerObjIDBizSet,
	BKInnerObjIDPlat,
	BKInnerObjIDSet,
}

// Revision
const (
	RevisionEnterprise = "enterprise"
	RevisionCommunity  = "community"
	RevisionOpensource = "opensource"
)

const (
	// BKDBMULTIPLELike TODO
	// used only for host search
	BKDBMULTIPLELike = "$multilike"

	// BKDBIN the db operator
	BKDBIN = "$in"

	// BKDBOR the db operator
	BKDBOR = "$or"

	// BKDBNOR the not or db operator
	BKDBNOR = "$nor"

	// BKDBAND the db operator
	BKDBAND = "$and"

	// BKDBLIKE the db operator
	BKDBLIKE = "$regex"

	// BKDBOPTIONS the db operator,used with $regex
	// detail to see https://docs.mongodb.com/manual/reference/operator/query/regex/#op._S_options
	BKDBOPTIONS = "$options"

	// BKDBEQ the db operator
	BKDBEQ = "$eq"

	// BKDBNE the db operator
	BKDBNE = "$ne"

	// BKDBNIN the db oeprator
	BKDBNIN = "$nin"

	// BKDBLT the db operator
	BKDBLT = "$lt"

	// BKDBLTE the db operator
	BKDBLTE = "$lte"

	// BKDBGT the db operator
	BKDBGT = "$gt"

	// BKDBGTE the db opeartor
	BKDBGTE = "$gte"

	// BKDBExists the db opeartor
	BKDBExists = "$exists"

	// BKDBNot the db opeartor
	BKDBNot = "$not"

	// BKDBCount the db opeartor
	BKDBCount = "$count"

	// BKDBGroup the db opeartor
	BKDBGroup = "$group"

	// BKDBMatch the db opeartor
	BKDBMatch = "$match"

	// BKDBSum the db opeartor
	BKDBSum = "$sum"

	// BKDBPush the db opeartor
	BKDBPush = "$push"

	// BKDBUNSET the db opeartor
	BKDBUNSET = "$unset"

	// BKDBAddToSet The $addToSet operator adds a value to an array unless the value is already present, in which case $addToSet does nothing to that array.
	BKDBAddToSet = "$addToSet"

	// BKDBPull The $pull operator removes from an existing array all instances of a value or values that match a specified condition.
	BKDBPull = "$pull"

	// BKDBAll matches arrays that contain all elements specified in the query.
	BKDBAll = "$all"

	// BKDBProject passes along the documents with the requested fields to the next stage in the pipeline
	BKDBProject = "$project"

	// BKDBSize counts and returns the total number of items in an array
	BKDBSize = "$size"

	// BKDBType selects documents where the value of the field is an instance of the specified BSON type(s).
	// Querying by data type is useful when dealing with highly unstructured data where data types are not predictable.
	BKDBType = "$type"

	// BKDBSort the db operator
	BKDBSort = "$sort"

	// BKDBReplaceRoot the db operator
	BKDBReplaceRoot = "$replaceRoot"

	// BKDBLimit the db operator to limit return number of doc
	BKDBLimit = "$limit"

	// BKDBUnwind used to split values contained in an array field into separate doc
	BKDBUnwind = "$unwind"

	// BKDBLookup used to perform multi-table association operations
	BKDBLookup = "$lookup"

	// BKDBFrom collection to join
	BKDBFrom = "from"

	// BKDBLocalField field from the input documents
	BKDBLocalField = "localField"

	// BKDBForeignField field from the documents of the "from" collection
	BKDBForeignField = "foreignField"

	// BKDBAs output array field
	BKDBAs = "as"

	// BKDBSkip skip data index
	BKDBSkip = "$skip"

	// BKLimit data quantity limit
	BKLimit = "$limit"
)

const (
	// DefaultResModuleName the default idle module name
	DefaultResModuleName string = "空闲机"
	// DefaultFaultModuleName the default fault module name
	DefaultFaultModuleName string = "故障机"
	// DefaultRecycleModuleName the default fault module name
	DefaultRecycleModuleName string = "待回收"
)

const (
	// BKFieldDBID the db id definition
	BKFieldDBID = "_id"
	// BKFieldSeqID the sequence id definition
	BKFieldSeqID = "SequenceID"

	// BKFieldID the id definition
	BKFieldID = "id"
	// BKFieldName TODO
	BKFieldName = "name"

	// BKDefaultField the default field
	BKDefaultField = "default"

	// TenantID the tenant field
	TenantID = "tenant_id"

	// BKAppIDField the appid field
	BKAppIDField = "bk_biz_id"

	// BKIPArr the ip address
	BKIPArr = "ipArr"

	// BKAssetIDField  the asset id field
	BKAssetIDField = "bk_asset_id"

	// BKSNField  the sn  field
	BKSNField = "bk_sn"

	// BKHostInnerIPField the host innerip field
	BKHostInnerIPField = "bk_host_innerip"

	// BKHostCloudRegionField the host cloud region field
	BKHostCloudRegionField = "bk_cloud_region"

	// BKHostOuterIPField the host outerip field
	BKHostOuterIPField = "bk_host_outerip"

	// BKAddressingStatic the host addressing is static
	BKAddressingStatic = "static"

	// BKAddressingDynamic the host addressing is dynamic
	BKAddressingDynamic = "dynamic"

	// BKCloudInstIDField the cloud instance id field
	BKCloudInstIDField = "bk_cloud_inst_id"

	// BKCloudHostStatusField the cloud host status field
	BKCloudHostStatusField = "bk_cloud_host_status"

	// TimeTransferModel the time transferModel field
	TimeTransferModel = "2006-01-02 15:04:05"

	// TimeDayTransferModel the time transferModel field
	TimeDayTransferModel = "2006-01-02"

	// BKCloudTaskID the cloud sync task id
	BKCloudTaskID = "bk_task_id"

	// BKNewAddHost the cloud sync new add hosts
	BKNewAddHost = "new_add"

	// BKImportFrom the host import from field
	BKImportFrom = "import_from"

	// BKHostIDField the host id field
	BKHostIDField = "bk_host_id"

	// BKHostNameField the host name field
	BKHostNameField = "bk_host_name"

	// BKAppNameField the app name field
	BKAppNameField = "bk_biz_name"

	// BKSetIDField the setid field
	BKSetIDField = "bk_set_id"

	// BKSetNameField the set name field
	BKSetNameField = "bk_set_name"

	// BKModuleIDField the module id field
	BKModuleIDField = "bk_module_id"

	// BKModuleNameField the module name field
	BKModuleNameField = "bk_module_name"

	// HostApplyEnabledField TODO
	HostApplyEnabledField = "host_apply_enabled"

	// BKOSTypeField the os type field
	BKOSTypeField = "bk_os_type"

	// BKOSNameField the os name field
	BKOSNameField = "bk_os_name"

	// BKOSVersionField the host os version field
	BKOSVersionField = "bk_os_version"

	// BKOSBitField the host os bit field
	BKOSBitField = "bk_os_bit"

	// BKCpuField the host cpu field
	BKCpuField = "bk_cpu"

	// BKOsKernelVersionField the host os kernel version field
	BKOsKernelVersionField = "bk_os_kernel_version"

	// BKCpuModuleField the host cpu module field
	BKCpuModuleField = "bk_cpu_module"

	// BKMemField the host memory field
	BKMemField = "bk_mem"

	// BKDiskField the host disk field
	BKDiskField = "bk_disk"

	// BKMac the host inner mac field
	BKMac = "bk_mac"

	// BKOuterMac the host outer mac field
	BKOuterMac = "bk_outer_mac"

	// BKCpuArch the host cpu architecture field
	BKCpuArch = "bk_cpu_architecture"

	// BKHttpGet the http get
	BKHttpGet = "GET"

	// BKTencentCloudTimeOut the tencent cloud timeout
	BKTencentCloudTimeOut = 10

	// TencentCloudUrl the tencent cloud url
	TencentCloudUrl = "cvm.tencentcloudapi.com"

	// TencentCloudSignMethod the tencent cloud sign method
	TencentCloudSignMethod = "HmacSHA1"

	// BKCloudIDField the cloud id field
	BKCloudIDField = "bk_cloud_id"

	// BKCloudNameField the cloud name field
	BKCloudNameField = "bk_cloud_name"

	// BKObjIDField the obj id field
	BKObjIDField = "bk_obj_id"

	// BKObjNameField the obj name field
	BKObjNameField = "bk_obj_name"

	// ObjSortNumberField the sort number field
	ObjSortNumberField = "obj_sort_number"

	// BKObjIconField the obj icon field
	BKObjIconField = "bk_obj_icon"

	// BKInstIDField the inst id field
	BKInstIDField = "bk_inst_id"

	// BKInstNameField the inst name field
	BKInstNameField = "bk_inst_name"

	// ExportCustomFields the use custom display columns
	ExportCustomFields = "export_custom_fields"

	// BKProcIDField the proc id field
	BKProcIDField = "bk_process_id"

	// BKConfTempIdField is the config template id field
	BKConfTempIdField = "bk_conftemp_id"

	// BKProcNameField the proc name field
	BKProcNameField = "bk_process_name"

	// BKTemlateIDField the process template id field
	BKTemlateIDField = "template_id"

	// BKVersionIDField the version id field
	BKVersionIDField = "version_id"

	// BKTemplateNameField the template name field
	BKTemplateNameField = "template_name"

	// BKFileNameField the file name field
	BKFileNameField = "file_name"

	// BKPropertyIDField the propety id field
	BKPropertyIDField = "bk_property_id"

	// BKPropertyNameField the property name field
	BKPropertyNameField = "bk_property_name"

	// BKPropertyIndexField TODO
	BKPropertyIndexField = "bk_property_index"

	// BKPropertyTypeField the property type field
	BKPropertyTypeField = "bk_property_type"

	// BKPropertyGroupField TODO
	BKPropertyGroupField = "bk_property_group"

	// BKPropertyValueField the property value field
	BKPropertyValueField = "bk_property_value"

	// BKObjAttIDField the obj att id field
	BKObjAttIDField = "bk_object_att_id"

	// BKClassificationIDField the classification id field
	BKClassificationIDField = "bk_classification_id"

	// BKClassificationTypeField the classification type field
	BKClassificationTypeField = "bk_classification_type"

	// BKClassificationNameField the classification name field
	BKClassificationNameField = "bk_classification_name"

	// BKClassificationIconField the classification icon field
	BKClassificationIconField = "bk_classification_icon"

	// BKPropertyGroupIDField the property group id field
	BKPropertyGroupIDField = "bk_group_id"

	// BKPropertyGroupNameField the property group name field
	BKPropertyGroupNameField = "bk_group_name"

	// BKPropertyGroupIndexField the property group index field
	BKPropertyGroupIndexField = "bk_group_index"

	// BKAsstObjIDField the property obj id field
	BKAsstObjIDField = "bk_asst_obj_id"

	// BKAsstInstIDField the property inst id field
	BKAsstInstIDField = "bk_asst_inst_id"

	// BKOptionField the option field
	BKOptionField = "option"

	// BKAuditTypeField the audit type field
	BKAuditTypeField = "audit_type"

	// BKResourceTypeField the audit resource type field
	BKResourceTypeField = "resource_type"

	// BKOperateFromField the platform where operation from field
	BKOperateFromField = "operate_from"

	// BKOperationDetailField the audit operation detail field
	BKOperationDetailField = "operation_detail"

	// BKOperationTimeField the audit operation time field
	BKOperationTimeField = "operation_time"

	// BKResourceIDField the audit resource ID field
	BKResourceIDField = "resource_id"

	// BKResourceNameField the audit resource name field
	BKResourceNameField = "resource_name"

	// BKExtendResourceNameField the audit extend resource name field
	BKExtendResourceNameField = "extend_resource_name"

	// BKLabelField the audit resource name field
	BKLabelField = "label"

	// BKSetEnvField the set env field
	BKSetEnvField = "bk_set_env"

	// BKSetStatusField the set status field
	BKSetStatusField = "bk_service_status"

	// BKSetDescField the set desc field
	BKSetDescField = "bk_set_desc"

	// BKBizSetDescField the biz set desc field
	BKBizSetDescField = "bk_biz_set_desc"

	// BKSetCapacityField the set capacity field
	BKSetCapacityField = "bk_capacity"

	// BKPort the port
	BKPort = "port"

	// BKProcPortEnable whether enable port,  enable port use for monitor app. default value
	BKProcPortEnable = "bk_enable_port"

	// BKProcGatewayIP the process gateway ip
	BKProcGatewayIP = "bk_gateway_ip"

	// BKProcGatewayPort the process gateway port
	BKProcGatewayPort = "bk_gateway_port"

	// BKProcGatewayProtocol the process gateway protocol
	BKProcGatewayProtocol = "bk_gateway_protocol"

	// BKProcGatewayCity the process gateway city
	BKProcGatewayCity = "bk_gateway_city"

	// BKProcBindInfo the process bind info
	BKProcBindInfo = "bind_info"

	// BKUser the user
	BKUser = "user"

	// BKProtocol the protocol
	BKProtocol = "protocol"

	// BKIP the ip
	BKIP = "ip"

	// BKEnable the enable
	BKEnable = "enable"

	// BKProcessObjectName TODO
	// the process object name
	BKProcessObjectName = "process"

	// BKProcessIDField the process id field
	BKProcessIDField = "bk_process_id"

	// BKServiceInstanceIDField TODO
	BKServiceInstanceIDField = "service_instance_id"
	// BKServiceTemplateIDField TODO
	BKServiceTemplateIDField = "service_template_id"
	// BKProcessTemplateIDField TODO
	BKProcessTemplateIDField = "process_template_id"
	// BKServiceCategoryIDField TODO
	BKServiceCategoryIDField = "service_category_id"

	// BKSetTemplateIDField TODO
	BKSetTemplateIDField = "set_template_id"

	// HostApplyRuleIDField TODO
	HostApplyRuleIDField = "host_apply_rule_id"

	// BKParentIDField TODO
	BKParentIDField = "bk_parent_id"
	// BKRootIDField TODO
	BKRootIDField = "bk_root_id"

	// BKProcessNameField the process name field
	BKProcessNameField = "bk_process_name"

	// BKFuncIDField the func id field
	BKFuncIDField = "bk_func_id"

	// BKFuncName the function name
	BKFuncName = "bk_func_name"

	// BKStartParamRegex TODO
	BKStartParamRegex = "bk_start_param_regex"

	// BKBindIP the bind ip
	BKBindIP = "bind_ip"

	// BKWorkPath the work path
	BKWorkPath = "work_path"

	// BKIsPre the ispre field
	BKIsPre = "ispre"

	// BKObjectUniqueKeys object unique keys field
	BKObjectUniqueKeys = "keys"

	// BKIsIncrementField the isincrement field
	BKIsIncrementField = "is_increment"

	// BKIsCollapseField TODO
	BKIsCollapseField = "is_collapse"

	// BKProxyListField the proxy list field
	BKProxyListField = "bk_proxy_list"

	// BKIPListField the ip list field
	BKIPListField = "ip_list"

	// BKInvalidIPSField the invalid ips field
	BKInvalidIPSField = "invalid_ips"

	// BKGseProxyField the gse proxy
	BKGseProxyField = "bk_gse_proxy"

	// BKSubAreaField the sub area field
	BKSubAreaField = "bk_cloud_id"

	// BKProcField the proc field
	BKProcField = "bk_process"

	// BKMaintainersField the maintainers field
	BKMaintainersField = "bk_biz_maintainer"

	// BKProductPMField the product pm field
	BKProductPMField = "bk_biz_productor"

	// BKTesterField the tester field
	BKTesterField = "bk_biz_tester"

	// BKOperatorField the operator field
	BKOperatorField = "operator" // the operator of app of module, is means a job position

	// BKLifeCycleField the life cycle field
	BKLifeCycleField = "life_cycle"

	// BKDeveloperField the developer field
	BKDeveloperField = "bk_biz_developer"

	// BKLanguageField the language field
	BKLanguageField = "language"

	// BKBakOperatorField the bak operator field
	BKBakOperatorField = "bk_bak_operator"

	// BKTimeZoneField the time zone field
	BKTimeZoneField = "time_zone"

	// BKIsRequiredField the required field
	BKIsRequiredField = "isrequired"

	// BKModuleTypeField the module type field
	BKModuleTypeField = "bk_module_type"

	// BKOrgIPField the org ip field
	BKOrgIPField = "bk_org_ip"

	// BKDstIPField the dst ip field
	BKDstIPField = "bk_dst_ip"

	// BKDescriptionField the description field
	BKDescriptionField = "description"

	// BKIsOnlyField the isonly name field
	BKIsOnlyField = "isonly"
	// BKIsMultipleField the is multiple name field
	BKIsMultipleField = "ismultiple"
	// BKDefaultFiled the is default name field
	BKDefaultFiled = "default"
	// BKGseTaskIDField the gse taskid
	BKGseTaskIDField = "task_id"
	// BKTaskIDField the gse taskid
	BKTaskIDField = "task_id"
	// BKGseOpTaskIDField the gse taskid
	BKGseOpTaskIDField = "gse_task_id"
	// BKProcPidFile TODO
	BKProcPidFile = "pid_file"
	// BKProcStartCmd TODO
	BKProcStartCmd = "start_cmd"
	// BKProcStopCmd TODO
	BKProcStopCmd = "stop_cmd"
	// BKProcReloadCmd TODO
	BKProcReloadCmd = "reload_cmd"
	// BKProcRestartCmd TODO
	BKProcRestartCmd = "restart_cmd"
	// BKProcTimeOut TODO
	BKProcTimeOut = "timeout"
	// BKProcWorkPath TODO
	BKProcWorkPath = "work_path"
	// BKProcInstNum TODO
	BKProcInstNum = "proc_num"

	// BKActionField TODO
	BKActionField = "action"

	// BKGseOpProcTaskDetailField gse operate process return detail
	BKGseOpProcTaskDetailField = "detail"
	// BKGroupField TODO
	BKGroupField = "group"

	// BKAttributeIDField TODO
	BKAttributeIDField = "bk_attribute_id"

	// BKTokenField TODO
	BKTokenField = "token"
	// BKCursorField TODO
	BKCursorField = "cursor"
	// BKClusterTimeField TODO
	BKClusterTimeField = "cluster_time"
	// BKEventTypeField TODO
	BKEventTypeField = "type"
	// BKStartAtTimeField TODO
	BKStartAtTimeField = "start_at_time"
	// BKSubResourceField TODO
	BKSubResourceField = "bk_sub_resource"

	// BKBizSetIDField TODO
	BKBizSetIDField = "bk_biz_set_id"
	// BKBizSetNameField TODO
	BKBizSetNameField = "bk_biz_set_name"
	// BKBizSetScopeField TODO
	BKBizSetScopeField = "bk_scope"
	// BKBizSetMatchField TODO
	BKBizSetMatchField = "match_all"

	// BKHostInnerIPv6Field the host innerip field in the form of ipv6
	BKHostInnerIPv6Field = "bk_host_innerip_v6"

	// BKHostOuterIPv6Field the host outerip field in the form of ipv6
	BKHostOuterIPv6Field = "bk_host_outerip_v6"

	// BKAgentIDField the agent id field, used by agent to identify a host
	BKAgentIDField = "bk_agent_id"

	// BKCloudHostIdentifierField defines if the host is a cloud host that doesn't allow cross biz transfer
	BKCloudHostIdentifierField = "bk_cloud_host_identifier"

	// BKAddressingField the addressing field, defines the host addressing type
	BKAddressingField = "bk_addressing"
	// BKProjectIDField the project id field
	BKProjectIDField = "bk_project_id"
	// BKProjectNameField the project name field
	BKProjectNameField = "bk_project_name"
	// BKProjectCodeField the project code field
	BKProjectCodeField = "bk_project_code"
	// BKProjectDescField the project desc field
	BKProjectDescField = "bk_project_desc"
	// BKProjectTypeField the project type field
	BKProjectTypeField = "bk_project_type"
	// BKProjectSecLvlField the project sec lvl field
	BKProjectSecLvlField = "bk_project_sec_lvl"
	// BKProjectOwnerField the project owner field
	BKProjectOwnerField = "bk_project_owner"
	// BKProjectTeamField the project team field
	BKProjectTeamField = "bk_project_team"
	// BKProjectStatusField the project status field
	BKProjectStatusField = "bk_status"
	// BKProjectIconField the project icon field
	BKProjectIconField = "bk_project_icon"

	// BKSrcModelField source model field in model relationship table.
	BKSrcModelField = "src_model"
	// BKDestModelField destination model field in the model relationship table.
	BKDestModelField     = "dest_model"
	BKDestModelUUIDField = "dest_model_uuid"

	// ObjectIDField the object id field, it is an int type field and is used to associate with the model
	ObjectIDField = "object_id"

	// BKCloudRegionField is the cloud region field
	BKCloudRegionField = "bk_cloud_region"
	// BKCloudZoneField is the cloud zone field
	BKCloudZoneField = "bk_cloud_zone"
)

const (
	// BKRequestField TODO
	BKRequestField = "bk_request_id"
	// BKTxnIDField TODO
	BKTxnIDField = "bk_txn_id"
)

const (
	// UserDefinedModules user define idle module key.
	UserDefinedModules = "user_modules"

	// SystemSetName user define idle module key.
	SystemSetName = "set_name"

	// SystemIdleModuleKey system idle module key.
	SystemIdleModuleKey = "idle"

	// SystemFaultModuleKey system define fault module name.
	SystemFaultModuleKey = "fault"

	// SystemRecycleModuleKey system define recycle module name.
	SystemRecycleModuleKey = "recycle"
)

// DefaultResSetName the inner module set
const DefaultResSetName string = "空闲机池"

// WhiteListAppName the white list app name
const WhiteListAppName = "蓝鲸"

// WhiteListSetName the white list set name
const WhiteListSetName = "公共组件"

// WhiteListModuleName the white list module name
const WhiteListModuleName = "gitserver"

// the inst record's logging information
const (
	// CreatorField the creator
	CreatorField = "creator"

	// CreateTimeField the create time field
	CreateTimeField = "create_time"

	// ConfirmTimeField the cloud resource confirm time filed
	ConfirmTimeField = "confirm_time"

	// StartTimeField the cloud sync start time field
	StartTimeFiled = "start_time"

	// ModifierField the modifier field
	ModifierField = "modifier"

	// LastTimeField the last time field
	LastTimeField = "last_time"

	// BKCreatedAt the model instance create time field
	BKCreatedAt = "bk_created_at"

	// BKCreatedBy the model instance creator field
	BKCreatedBy = "bk_created_by"

	// BKUpdatedAt the model instance update time field
	BKUpdatedAt = "bk_updated_at"

	// BKUpdatedBy the model instance modifier field
	BKUpdatedBy = "bk_updated_by"
)

const (
	// ValidCreate valid create
	ValidCreate = "create"

	// ValidUpdate valid update
	ValidUpdate = "update"
)

// DefaultResSetFlag the default resource set flat
const DefaultResSetFlag int = 1

// DefaultFlagDefaultValue TODO
const DefaultFlagDefaultValue int = 0

// DefaultAppFlag the default app flag
const DefaultAppFlag int = 1

// DefaultAppName the default app name
const DefaultAppName string = "资源池"

// DefaultCloudName default area
const DefaultCloudName string = "Default Area"

// DefaultInstName TODO
const DefaultInstName string = "实例名"

// BKAppName the default app name
const BKAppName string = "蓝鲸"

// BKNetwork TODO
// bk_classification_id value
const BKNetwork = "bk_network"

const (
	// SNMPActionGet TODO
	SNMPActionGet = "get"

	// SNMPActionGetNext TODO
	SNMPActionGetNext = "getnext"
)

const (
	// DefaultResModuleFlag the default resource module flag
	DefaultResModuleFlag int = 1

	// DefaultFaultModuleFlag the default fault module flag
	DefaultFaultModuleFlag int = 2

	// NormalModuleFlag create module by user , default =0
	NormalModuleFlag int = 0

	// NormalSetDefaultFlag user create set default field value
	NormalSetDefaultFlag int64 = 0

	// DefaultRecycleModuleFlag default recycle module flag
	DefaultRecycleModuleFlag int = 3

	// DefaultResSelfDefinedModuleFlag the default resource self-defined module flag
	DefaultResSelfDefinedModuleFlag int = 4

	// DefaultUserResModuleFlag the default platform self-defined module flag.
	DefaultUserResModuleFlag int = 5
)

const (
	// DefaultModuleType TODO
	DefaultModuleType string = "1"
)

var FieldTypes = []string{FieldTypeSingleChar, FieldTypeLongChar, FieldTypeInt, FieldTypeFloat, FieldTypeEnum,
	FieldTypeEnumMulti, FieldTypeDate, FieldTypeTime, FieldTypeUser, FieldTypeOrganization, FieldTypeTimeZone,
	FieldTypeBool, FieldTypeList, FieldTypeTable, FieldTypeInnerTable, FieldTypeEnumQuote}

const (
	// FieldTypeSingleChar the single char filed type
	FieldTypeSingleChar string = "singlechar"

	// FieldTypeLongChar the long char field type
	FieldTypeLongChar string = "longchar"

	// FieldTypeInt the int field type
	FieldTypeInt string = "int"

	// FieldTypeFloat the float field type
	FieldTypeFloat string = "float"

	// FieldTypeEnum the enum field type
	FieldTypeEnum string = "enum"

	// FieldTypeEnumMulti the enum multi field type
	FieldTypeEnumMulti string = "enummulti"

	// FieldTypeEnumQuote the enum quote field type
	FieldTypeEnumQuote string = "enumquote"

	// FieldTypeDate the date field type
	FieldTypeDate string = "date"

	// FieldTypeTime the time field type
	FieldTypeTime string = "time"

	// FieldTypeUser the user field type
	FieldTypeUser string = "objuser"

	// FieldObject 此处只校验是否是对象。此校验是为了兼容biz_set中的bk_scope 的类型是 querybuilder，由于在 coreservice层解析出来的
	// 是map[string]interface,所以在此处只需要校验是否是对象，对于querybuilder的合法性应该放在场景层做校验。后续如果走的是object校验，
	// 都需要在场景层进行真正的校验
	FieldObject string = "object"

	// FieldTypeTimeZone the timezone field type
	FieldTypeTimeZone string = "timezone"

	// FieldTypeBool the bool type
	FieldTypeBool string = "bool"

	// FieldTypeList the list type
	FieldTypeList string = "list"

	// FieldTypeTable the table type, inner type.
	FieldTypeTable string = "table"

	// FieldTypeInnerTable the table type, the type
	// used for table fields in model reference scenarios
	FieldTypeInnerTable string = "innertable"

	// FieldTypeOrganization the organization field type
	FieldTypeOrganization string = "organization"

	// FieldTypeIDRule the id rule field type
	FieldTypeIDRule string = "id_rule"

	// FieldTypeSingleLenChar the single char length limit
	FieldTypeSingleLenChar int = 256

	// FieldTypeLongLenChar the long char length limit
	FieldTypeLongLenChar int = 2000

	// FieldTypeUserLenChar the user char length limit
	FieldTypeUserLenChar int = 2000

	// FieldTypeStrictCharRegexp the single char regex expression
	FieldTypeStrictCharRegexp string = `^[a-zA-Z]\w*$`

	// FieldTypeServiceCategoryRegexp the service category regex expression
	FieldTypeServiceCategoryRegexp string = `^([\w\p{Han}]|[:\-\(\)])+$`

	// FieldTypeMainlineRegexp the mainline instance name regex expression
	FieldTypeMainlineRegexp string = `^[^\\\|\/:\*,<>"\?#\s]+$`

	// FieldTypeSingleCharRegexp the single char regex expression
	FieldTypeSingleCharRegexp string = `\S`

	// FieldTypeLongCharRegexp the long char regex expression
	FieldTypeLongCharRegexp string = `\S`
)

const (
	// HostAddMethodExcel add a host method
	HostAddMethodExcel = "1"

	// HostAddMethodAgent add a  agent method
	HostAddMethodAgent = "2"

	// HostAddMethodAPI add api method
	HostAddMethodAPI = "3"

	// AddExcelDataIndexOffset the index of the add excel data
	AddExcelDataIndexOffset = 6

	// HostAddMethodExcelAssociationIndexOffset TODO
	HostAddMethodExcelAssociationIndexOffset = 2

	// HostAddMethodExcelDefaultIndex 生成表格数据起始索引，第一列为字段说明
	HostAddMethodExcelDefaultIndex = 1

	/*EXCEL color AARRGGBB :
	AA means Alpha
	RRGGBB means Red, in hex.
	GG means Red, in hex.
	BB means Red, in hex.
	*/

	// ExcelHeaderFirstRowColor cell bg color
	ExcelHeaderFirstRowColor = "FF92D050"
	// ExcelHeaderFirstRowFontColor  font color
	ExcelHeaderFirstRowFontColor = "00000000"
	// ExcelHeaderFirstRowRequireFontColor require font color
	ExcelHeaderFirstRowRequireFontColor = "FFFF0000"
	// ExcelHeaderOtherRowColor cell bg color
	ExcelHeaderOtherRowColor = "FFC6EFCE"
	// ExcelHeaderOtherRowFontColor font color
	ExcelHeaderOtherRowFontColor = "FF000000"
	// ExcelCellDefaultBorderColor black color
	ExcelCellDefaultBorderColor = "FFD4D4D4"
	// ExcelHeaderFirstColumnColor light gray
	ExcelHeaderFirstColumnColor = "fee9da"
	// ExcelFirstColumnCellColor dark gray
	ExcelFirstColumnCellColor = "fabf8f"
	// ExcelTableHeaderColor excel table header color
	ExcelTableHeaderColor = "d1e0b6"

	// ExcelAsstPrimaryKeySplitChar split char
	ExcelAsstPrimaryKeySplitChar = ","
	// ExcelAsstPrimaryKeyJoinChar split char
	ExcelAsstPrimaryKeyJoinChar = "="
	// ExcelAsstPrimaryKeyRowChar split char
	ExcelAsstPrimaryKeyRowChar = "\n"

	// ExcelDelAsstObjectRelation delete asst object relation
	ExcelDelAsstObjectRelation = "/"

	// ExcelDataValidationListLen excel dropdown list item count
	ExcelDataValidationListLen = 50

	// ExcelCommentSheetCotentLangPrefixKey excel comment sheet centent language prefixe key
	ExcelCommentSheetCotentLangPrefixKey = "import_comment"

	// ExcelFirstColumnFieldName export excel first column for tips
	ExcelFirstColumnFieldName = "field_name"
	// ExcelFirstColumnFieldType excel first column type field
	ExcelFirstColumnFieldType = "field_type"
	// ExcelFirstColumnFieldID excel first column id field
	ExcelFirstColumnFieldID = "field_id"
	// ExcelFirstColumnTableFieldName excel first column table name filed
	ExcelFirstColumnTableFieldName = "table_field_name"
	// ExcelFirstColumnFieldType excel first column table type field
	ExcelFirstColumnTableFieldType = "table_field_type"
	// ExcelFirstColumnFieldID excel first column table id field
	ExcelFirstColumnTableFieldID = "table_field_id"
	// ExcelFirstColumnInstData excel first column instance data field
	ExcelFirstColumnInstData = "inst_data"

	// ExcelFirstColumnAssociationAttribute TODO
	ExcelFirstColumnAssociationAttribute = "excel_association_attribute"
	// ExcelFirstColumnFieldDescription TODO
	ExcelFirstColumnFieldDescription = "excel_field_description"

	// ExcelCellIgnoreValue TODO
	// the value of ignored excel cell
	ExcelCellIgnoreValue = "--"
)

// BizSetConditionMaxDeep 业务集场景下querybuilder条件的最大深度不能超过2层
const BizSetConditionMaxDeep = 2
const (
	// InputTypeExcel  data from excel
	InputTypeExcel = "excel"

	// InputTypeApiNewHostSync data from api for synchronize new host
	InputTypeApiNewHostSync = "api_sync_host"

	// BatchHostAddMaxRow batch sync add host max row
	BatchHostAddMaxRow = 128

	// ExcelImportMaxRow excel import max row
	ExcelImportMaxRow = 1000
)

// deprecated old api response fields only for legacy api
const (
	// HTTPBKAPIErrorMessage apiserver error message
	HTTPBKAPIErrorMessage = "bk_error_msg"

	// HTTPBKAPIErrorCode apiserver error code
	HTTPBKAPIErrorCode = "bk_error_code"
)

// new api response fields in the standard of api gateway
const (
	// BkAPIErrorCode error code field name
	BkAPIErrorCode = "code"

	// BkAPIErrorMessage error message field name
	BkAPIErrorMessage = "message"
)

// KvMap the map definition
type KvMap map[string]interface{}

const (
	// CCSystemOperatorUserName the system user
	CCSystemOperatorUserName = "cc_system"
	// CCSystemCollectorUserName TODO
	CCSystemCollectorUserName = "cc_collector"
)

// APIRsp the result the http requst
type APIRsp struct {
	HTTPCode int         `json:"-"`
	Result   bool        `json:"result"`
	Code     int         `json:"code"`
	Message  interface{} `json:"message"`
	Data     interface{} `json:"data"`
}

const (
	// BKCacheKeyV3Prefix the prefix definition
	BKCacheKeyV3Prefix = "cc:v3:"
)

// event cache keys
const (
	EventCacheEventIDKey = BKCacheKeyV3Prefix + "event:inst_id"
	RedisSnapKeyPrefix   = BKCacheKeyV3Prefix + "snapshot:"

	// RedisHostSnapMsgDelayQueue the monitor reports data, and queries the name of the key placed in the delay queue in
	// the case of host failure.
	RedisHostSnapMsgDelayQueue = BKCacheKeyV3Prefix + "host_snap:delay_queue"
)
const (
	// RedisSentinelMode redis mode is sentinel
	RedisSentinelMode = "sentinel"
	// RedisSingleMode redis mode is single
	RedisSingleMode = "single"
)

// api cache keys
const (
	ApiCacheLimiterRulePrefix = BKCacheKeyV3Prefix + "api:limiter_rule:"
)

// ReadReferenceKey cmdb read preference key
const ReadReferenceKey = "Cc_Read_Preference"

// ReadPreferenceMode mongodb read preference mode
// https://www.mongodb.com/docs/manual/core/read-preference/
type ReadPreferenceMode string

// String 用于打印
func (r ReadPreferenceMode) String() string {
	return string(r)
}

// BKHTTPReadRefernceMode constants  这个位置对应的是mongodb 的read preference 的mode，如果driver 没有变化这里是不需要变更的，
// 新增mode 需要修改src/storage/dal/mongo/local/mongo.go 中的getCollectionOption 方法来支持
const (
	// NilMode not set
	NilMode ReadPreferenceMode = ""
	// PrimaryMode indicates that only a primary is
	// considered for reading. This is the default
	// mode.
	PrimaryMode ReadPreferenceMode = "1"
	// PrimaryPreferredMode indicates that if a primary
	// is available, use it; otherwise, eligible
	// secondaries will be considered.
	PrimaryPreferredMode ReadPreferenceMode = "2"
	// SecondaryMode indicates that only secondaries
	// should be considered.
	SecondaryMode ReadPreferenceMode = "3"
	// SecondaryPreferredMode indicates that only secondaries
	// should be considered when one is available. If none
	// are available, then a primary will be considered.
	SecondaryPreferredMode ReadPreferenceMode = "4"
	// NearestMode indicates that all primaries and secondaries
	// will be considered.
	NearestMode ReadPreferenceMode = "5"
)

// transaction related
const (
	TransactionIdHeader       = "cc_transaction_id_string"
	TransactionTimeoutHeader  = "cc_transaction_timeout"
	TransactionTenantIDHeader = "cc_transaction_tenant_id"

	// mongodb default transaction timeout is 1 minute.
	TransactionDefaultTimeout = 2 * time.Minute
)

const (
	// DefaultAppLifeCycleNormal  biz life cycle normal
	DefaultAppLifeCycleNormal = "2"
)

// Host OS type enumeration value
const (
	HostOSTypeEnumLinux   = "1"
	HostOSTypeEnumWindows = "2"
	HostOSTypeEnumAIX     = "3"
	HostOSTypeEnumUNIX    = "4"
	HostOSTypeEnumSolaris = "5"
	HostOSTypeEnumHpUX    = "6"
	HostOSTypeEnumFreeBSD = "7"
	HostOSTypeEnumMacOS   = "8"
)

// HostOSTypeName Host system enum and name association
var HostOSTypeName = map[string]string{
	HostOSTypeEnumLinux:   "linux",
	HostOSTypeEnumWindows: "windows",
	HostOSTypeEnumAIX:     "aix",
	HostOSTypeEnumUNIX:    "unix",
	HostOSTypeEnumSolaris: "solaris",
	HostOSTypeEnumHpUX:    "hp-ux",
	HostOSTypeEnumFreeBSD: "freebsd",
	HostOSTypeEnumMacOS:   "darwin",
}

const (
	// MaxProcessPrio TODO
	MaxProcessPrio = 10000
	// MinProcessPrio TODO
	MinProcessPrio = -100
)

// integer const
const (
	MaxUint64  = ^uint64(0)
	MinUint64  = 0
	MaxInt64   = int64(MaxUint64 >> 1)
	MinInt64   = -MaxInt64 - 1
	MaxUint    = ^uint(0)
	MinUint    = 0
	MaxInt     = int(MaxUint >> 1)
	MinInt     = -MaxInt - 1
	MaxFloat64 = math.MaxFloat64
	MinFloat64 = -math.MaxFloat64
)

// config admin
const (
	ConfigAdminID         = "configadmin"
	ConfigAdminValueField = "config"
)

const (
	// PlatformConfig platform config
	PlatformConfig = "platform_config"
)

const (
	// APPConfigWaitTime application wait config from zookeeper time (unit sencend)
	APPConfigWaitTime = 15
)

const (
	// URLFilterWhiteListSuffix url filter white list not execute any filter
	// multiple url separated by commas
	URLFilterWhiteListSuffix = "/healthz,/version,/monitor_healthz"

	// URLFilterWhiteListSepareteChar TODO
	URLFilterWhiteListSepareteChar = ","
)

// DataStatusFlag TODO
type DataStatusFlag string

const (
	// DataStatusDisabled TODO
	DataStatusDisabled DataStatusFlag = "disabled"
	// DataStatusEnable TODO
	DataStatusEnable DataStatusFlag = "enable"
)

const (
	// BKDataStatusField TODO
	BKDataStatusField = "bk_data_status"
	// BKDataRecoverSuffix TODO
	BKDataRecoverSuffix = "(recover)"
)

const (
	// Infinite TODO
	// period default value
	Infinite = "∞"
)

const (
	// BKBluekingLoginPluginVersion TODO
	// login type
	BKBluekingLoginPluginVersion = "blueking"
	// BKOpenSourceLoginPluginVersion TODO
	BKOpenSourceLoginPluginVersion = "opensource"
	// BKSkipLoginPluginVersion TODO
	BKSkipLoginPluginVersion = "skip-login"

	// BKNoopMonitorPlugin TODO
	// monitor plugin type
	BKNoopMonitorPlugin = "noop"
	// BKBluekingMonitorPlugin TODO
	BKBluekingMonitorPlugin = "blueking"

	// HTTPCookieBKToken TODO
	HTTPCookieBKToken = "bk_token"
	// HTTPCookieBKTicket is the bk ticket cookie name
	HTTPCookieBKTicket = "bk_ticket"
	// HTTPCookieLanguage is the blueking language cookie name
	HTTPCookieLanguage = "blueking_language"
	// HTTPCookieTenant is the tenant cookie name
	HTTPCookieTenant = "HTTP_BLUEKING_TENANT_ID"

	// WEBSessionUinKey TODO
	WEBSessionUinKey = "username"
	// WEBSessionChineseNameKey TODO
	WEBSessionChineseNameKey = "chName"
	// WEBSessionTenantUinKey is the tenant uin session key
	WEBSessionTenantUinKey = "tenant_uin"
	// WEBSessionTenantUinListeKey is the tenant uin list session key
	WEBSessionTenantUinListeKey = "tenant_uin_list"
	// WEBSessionAvatarUrlKey TODO
	WEBSessionAvatarUrlKey = "avatar_url"
	// WEBSessionMultiTenantKey multi tenant session key
	WEBSessionMultiTenantKey = "multitenant"
	// WEBSessionTimeZoneKey time zone session key
	WEBSessionTimeZoneKey = "time_zone"

	// LoginSystemMultiTenantTrue login system multi tenant flag
	LoginSystemMultiTenantTrue = "1"
	// LoginSystemMultiTenantFalse TODO
	LoginSystemMultiTenantFalse = "0"

	// LogoutHTTPSchemeCookieKey TODO
	LogoutHTTPSchemeCookieKey = "http_scheme"
	// LogoutHTTPSchemeHTTP TODO
	LogoutHTTPSchemeHTTP = "http"
	// LogoutHTTPSchemeHTTPS TODO
	LogoutHTTPSchemeHTTPS = "https"
)

// BKStatusField TODO
const BKStatusField = "status"

const (
	// BKProcInstanceOpUser TODO
	BKProcInstanceOpUser = "proc instance user"

	// BKIAMSyncUser TODO
	BKIAMSyncUser = "iam_sync_user"
)

// association fields
const (
	// the id of the association kind
	AssociationKindIDField    = "bk_asst_id"
	AssociationKindNameField  = "bk_asst_name"
	AssociationObjAsstIDField = "bk_obj_asst_id"
	AssociatedObjectIDField   = "bk_asst_obj_id"
)

// association
const (
	AssociationKindMainline = "bk_mainline"
	AssociationTypeBelong   = "belong"
	AssociationTypeGroup    = "group"
	AssociationTypeRun      = "run"
	AssociationTypeConnect  = "connect"
	AssociationTypeDefault  = "default"
)

const (
	// MetadataField data business key
	MetadataField = "metadata"
)

const (
	// BKBizDefault TODO
	BKBizDefault = "bizdefault"
)

const (
	// AttributePlaceHolderMaxLength TODO
	AttributePlaceHolderMaxLength = 2000
	// AttributeOptionMaxLength TODO
	AttributeOptionMaxLength = 2000
	// AttributeIDMaxLength TODO
	AttributeIDMaxLength = 128
	// AttributeNameMaxLength TODO
	AttributeNameMaxLength = 128
	// AttributeUnitMaxLength TODO
	AttributeUnitMaxLength = 20
	// AttributeOptionValueMaxLength TODO
	AttributeOptionValueMaxLength = 128
	// AttributeOptionArrayMaxLength TODO
	AttributeOptionArrayMaxLength = 200
	// ServiceCategoryMaxLength TODO
	ServiceCategoryMaxLength = 128
)

const (
	// NameFieldMaxLength TODO
	NameFieldMaxLength = 256
	// MainlineNameFieldMaxLength TODO
	MainlineNameFieldMaxLength = 256

	// ServiceTemplateIDNotSet 用于表示还未设置服务模板的情况，比如没有绑定服务模板
	ServiceTemplateIDNotSet = 0
	// SetTemplateIDNotSet TODO
	SetTemplateIDNotSet = 0

	// MetadataLabelBiz TODO
	MetadataLabelBiz = "metadata.label.bk_biz_id"

	// DefaultServiceCategoryName TODO
	DefaultServiceCategoryName = "Default"
)

const (
	// ContextRequestIDField TODO
	ContextRequestIDField = "request_id"
	// ContextRequestUserField TODO
	ContextRequestUserField = "request_user"
	// ContextRequestTenantField tenant request field
	ContextRequestTenantField = "request_tenant"
)

const (
	// TimerPattern TODO
	TimerPattern = "^[\\d]+\\:[\\d]+$"
	// BKTaskTypeField the api task type field
	BKTaskTypeField = "task_type"
	// SyncSetTaskFlag TODO
	SyncSetTaskFlag = "set_template_sync"
	// SyncModuleTaskFlag TODO
	SyncModuleTaskFlag = "service_template_sync"
	// SyncFieldTemplateTaskFlag field template synchronization task flag
	SyncFieldTemplateTaskFlag = "field_template_sync"
	// SyncModuleHostApplyTaskFlag module dimension host auto-apply async task flag.
	SyncModuleHostApplyTaskFlag = "module_host_apply_sync"
	// SyncServiceTemplateHostApplyTaskFlag  service template dimension host auto-apply async task flag.
	SyncServiceTemplateHostApplyTaskFlag = "service_template_host_apply_sync"
	// SyncInstIDRuleTaskFlag  instance id rule async task flag.
	SyncInstIDRuleTaskFlag = "inst_id_rule_sync"

	// BKHostState TODO
	BKHostState = "bk_state"
)

// LanguageType TODO
// multiple language support
type LanguageType string

const (
	// Chinese TODO
	Chinese LanguageType = "zh-cn"
	// English TODO
	English LanguageType = "en"
)

// cloud area const
const (
	BKCloudVendor     = "bk_cloud_vendor"
	BKCloudSyncTaskID = "bk_task_id"
	BKCreator         = "bk_creator"
	BKLastEditor      = "bk_last_editor"
	BKVpcID           = "bk_vpc_id"
	BKVpcName         = "bk_vpc_name"
	BKRegion          = "bk_region"
)

// configcenter
const (
	BKDefaultConfigCenter = "zookeeper"
)

const (
	// CCLogicUniqueIdxNamePrefix TODO
	CCLogicUniqueIdxNamePrefix = "bkcc_unique_"
	// CCLogicIndexNamePrefix TODO
	CCLogicIndexNamePrefix = "bkcc_idx_"
)

const (
	// DefaultResBusinessSetFlag the default resource business set flag
	DefaultResBusinessSetFlag = 1
)

const (
	// HostFavoriteType host query favorite condition type
	HostFavoriteType = "type"
)

// ModelQuoteType model quote type.
type ModelQuoteType string

const (
	Table ModelQuoteType = "table"
)

const (
	// BKTemplateID template id field
	BKTemplateID = "bk_template_id"
)

// ProcessPropertyName process property name
type ProcessPropertyName string

const (
	// ProcNumName process property procNum name
	ProcNumName ProcessPropertyName = "proc_num"

	// StopCmdName process property stopCmd name
	StopCmdName ProcessPropertyName = "stop_cmd"

	// RestartCmdName process property restartCmd name
	RestartCmdName ProcessPropertyName = "restart_cmd"

	// FaceStopCmdName process property faceStopCmd name
	FaceStopCmdName ProcessPropertyName = "face_stop_cmd"

	// BkFuncNameName process property bkFuncName name
	BkFuncNameName ProcessPropertyName = "bk_func_name"

	// WorkPathName process property workPath name
	WorkPathName ProcessPropertyName = "work_path"

	// BindIpName process property bindIp name
	BindIpName ProcessPropertyName = "bind_ip"

	// PriorityName process property priority name
	PriorityName ProcessPropertyName = "priority"

	// ReloadCmdName process property reloadCmd name
	ReloadCmdName ProcessPropertyName = "reload_cmd"

	// BkProcessName process property bkProcessNam
	BkProcessName ProcessPropertyName = "bk_process_name"

	// PortName process property port name
	PortName ProcessPropertyName = "port"

	// PidFileName process property pidFile name
	PidFileName ProcessPropertyName = "pid_file"

	// AutoStartName process property autoStart name
	AutoStartName ProcessPropertyName = "auto_start"

	// BkStartCheckSecsName process property autoStart name
	BkStartCheckSecsName ProcessPropertyName = "bk_start_check_secs"

	// StartCmdName process property startCmd name
	StartCmdName ProcessPropertyName = "start_cmd"

	// PropertyUserName process property user
	PropertyUserName ProcessPropertyName = "user"

	// TimeoutName process property timeout name
	TimeoutName ProcessPropertyName = "timeout"

	// ProtocolName process property protocol name
	ProtocolName ProcessPropertyName = "protocol"

	// DescriptionName process property description name
	DescriptionName ProcessPropertyName = "description"

	// BkStartParamRegexName process property bkStartParamRegex name
	BkStartParamRegexName ProcessPropertyName = "bk_start_param_regex"
)

const (
	// Field query field define
	Field = "field"

	// Operator query operator define
	Operator = "operator"

	// Value query value define
	Value = "value"

	// Condition query condition define
	Condition = "condition"
)

const (
	// TopoModuleName topo path name
	TopoModuleName = "topo_module_name"
)

// BuiltInCloudAreaIDs 内置管控区域, 后续如果有其他新增的id，需要添加到这个数组里
var BuiltInCloudAreaIDs = []int64{BKDefaultDirSubArea, UnassignedCloudAreaID}

const (
	// UnassignedCloudAreaID 未分配的管控区域id
	UnassignedCloudAreaID = -1

	// UnassignedCloudAreaName 未分配的管控区域名称
	UnassignedCloudAreaName = "未分配"
)

// Default default field value
type Default int64

const (
	// BuiltIn built-in value
	BuiltIn Default = 1
)

const (
	// GlobalIDRule global id rule flag
	GlobalIDRule = "global"

	// IDRulePrefix id rule self-increasing id prefix
	IDRulePrefix = "id_rule:incr_id:"

	// GlobalIncrIDVar global self-increasing id variable
	GlobalIncrIDVar = "global.incr_id"

	// LocalIncrIDVar model self-increasing id variable
	LocalIncrIDVar = "local.incr_id"

	// RandomIDVar random id variable
	RandomIDVar = "random_id"
)

const (
	// MongoMetaID is mongodb meta id field
	MongoMetaID = "_id"
)

// ShardingDBConfID is sharding db config identifier
const ShardingDBConfID = "sharding_db_config"

// SnapshotChannelName is host snapshot channel name
const SnapshotChannelName = "cmdb_snapshot"
