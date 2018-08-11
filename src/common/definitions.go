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

	// BKTrue the true definition
	BKTrue = 1

	// BKFalse the false definition
	BKFalse = 2

	// BKNoLimit no limit definition
	BKNoLimit = 999999999

	// BKDefaultLimit the default limit definition
	BKDefaultLimit = 20

	// BKParent the parent code
	BKParent = 1
	// BKChild the child code
	BKChild = 2

	// BKParentStr the parent name
	BKParentStr = "bk_parentid"

	// BKChildStr the child name
	BKChildStr = "bk_childid"

	// BKInstParentStr the inst parent name
	BKInstParentStr = "bk_parent_id"

	// BKDefaultOwnerID the default owner value
	BKDefaultOwnerID = "0"

	// BKDefaultSupplierID the default owner id
	BKDefaultSupplierID = 0

	// BKSuperOwnerID the super owner value
	BKSuperOwnerID = "supperowner"

	// BKDefaultDirSubArea the default dir subarea
	BKDefaultDirSubArea = 0

	// BKTimeTypeParseFlag the time flag
	BKTimeTypeParseFlag = "cc_time_type"

	// BKTopoBusinessLevelLimit the mainline topo level limit
	BKTopoBusinessLevelLimit = "level.businessTopoMax"

	// BKTopoBusinessLevelDefault the mainline topo level default level
	BKTopoBusinessLevelDefault = 7
)

const (
	// BKInnerObjIDApp the inner object
	BKInnerObjIDApp = "biz"

	// BKInnerObjIDSet the inner object
	BKInnerObjIDSet = "set"

	// BKInnerObjIDModule the inner object
	BKInnerObjIDModule = "module"

	// BKInnerObjIDHost the inner object
	BKInnerObjIDHost = "host"

	// BKINnerObjIDObject the inner object
	BKINnerObjIDObject = "object"

	// BKInnerObjIDProc the inner object
	BKInnerObjIDProc = "process"

	// BKInnerObjIDPlat the inner object
	BKInnerObjIDPlat = "plat"

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

// Revision
const (
	RevisionEnterprise = "enterprise"
	RevisionCommunity  = "community"
	RevisionOpensource = "opensource"
)

const (
	// BKDBIN the db operator
	BKDBIN = "$in"

	// BKDBOR the db operator
	BKDBOR = "$or"

	// BKDBLIKE the db operator
	BKDBLIKE = "$regex"

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
	// BKDBSortFieldSep the db sort field split char
	BKDBSortFieldSep = ","
)

const (
	// DefaultResModuleName the default idle module name
	DefaultResModuleName string = "空闲机"
	// DefaultFaultModuleName the default fault module name
	DefaultFaultModuleName string = "故障机"
)

const (
	// BKFIeldID the id definition
	BKFieldID = "id"

	// BKDefaultField the default field
	BKDefaultField = "default"

	// BKOwnerIDField the owner field
	BKOwnerIDField = "bk_supplier_account"

	// BKSupplierIDField the supplier id field
	BKSupplierIDField = "bk_supplier_id"

	// BKAppIDField the appid field
	BKAppIDField = "bk_biz_id"

	// BKIPAddr the ip address
	BKIPArr = "ipArr"

	// BKHostInnerIPField the host innerip field
	BKHostInnerIPField = "bk_host_innerip"

	// BKHostOuterIPField the host outerip field
	BKHostOuterIPField = "bk_host_outerip"

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

	// BKSubscriptionIDField the subscription id field
	BKSubscriptionIDField = "subscription_id"
	// BKSubscriptionNameField the subscription name field
	BKSubscriptionNameField = "subscription_name"

	// BKOSTypeField the os type field
	BKOSTypeField = "bk_os_type"

	// BKOSNameField the os name field
	BKOSNameField = "bk_os_name"

	// BKCloudIDField the cloud id field
	BKCloudIDField = "bk_cloud_id"

	// BKCloudNameField the cloud name field
	BKCloudNameField = "bk_cloud_name"

	// BKObjIDField the obj id field
	BKObjIDField = "bk_obj_id"

	// BKObjNameField the obj name field
	BKObjNameField = "bk_obj_name"

	// BKObjIconField the obj icon field
	BKObjIconField = "bk_obj_icon"

	// BKInstIDField the inst id field
	BKInstIDField = "bk_inst_id"

	// BKInstNameField the inst name field
	BKInstNameField = "bk_inst_name"

	// BKProcIDField the proc id field
	BKProcIDField = "bk_process_id"

	// BKConfTempId is the config template id field
	BKConfTempIdField = "bk_conftemp_id"

	// BKProcNameField the proc name field
	BKProcNameField = "bk_process_name"

	// BKPropertyIDField the propety id field
	BKPropertyIDField = "bk_property_id"

	// BKPropertyNameField the property name field
	BKPropertyNameField = "bk_property_name"

	// BKPropertyTypeField the property type field
	BKPropertyTypeField = "bk_property_type"

	BKPropertyGroupField = "bk_property_group"

	// BKPropertyValueField the property value field
	BKPropertyValueField = "bk_property_value"

	// BKObjAttIDField the obj att id field
	BKObjAttIDField = "bk_object_att_id"

	// BKClassificationIDField the classification id field
	BKClassificationIDField = "bk_classification_id"

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

	// BKPrivilegeField the privilege field
	BKPrivilegeField = "privilege"

	// BKUserGroupIDField the group id field
	BKUserGroupIDField = "group_id"

	// BKUserListField the user list field
	BKUserListField = "user_list"

	// BKContentField the content field
	BKContentField = "content"

	// BKExtKeyField the ext key field
	BKExtKeyField = "ext_key"

	// BKOpDescField the op desc field
	BKOpDescField = "op_desc"

	// BKOpTypeField the op type field
	BKOpTypeField = "op_type"

	// BKOpTargetField the op target field
	BKOpTargetField = "op_target"

	// BKOpTimeField the op time field
	BKOpTimeField = "op_time"

	// BKSetEnvField the set env field
	BKSetEnvField = "bk_set_env"

	// BKSetStatusField the set status field
	BKSetStatusField = "bk_service_status"

	// BKSetDescField the set desc field
	BKSetDescField = "bk_set_desc"

	// BKSetCapacityField the set capacity field
	BKSetCapacityField = "bk_capacity"

	// BKPort the port
	BKPort = "port"

	// BKUser the user
	BKUser = "user"

	// BKProtocol the protocol
	BKProtocol = "protocol"

	// BKProcessIDField the process id field
	BKProcessIDField = "bk_process_id"

	// BKProcessNameField the process name field
	BKProcessNameField = "bk_process_name"

	// BKFuncIDField the func id field
	BKFuncIDField = "bk_func_id"

	// BKFuncName the function name
	BKFuncName = "bk_func_name"

	// BKBindIP the bind ip
	BKBindIP = "bind_ip"

	// BKWorkPath the work path
	BKWorkPath = "work_path"

	// BKIsPre the ispre field
	BKIsPre = "ispre"

	// bkIsOnly the isonly field, is primary key
	BKIsOnly = "isonly"

	// BKIsIncrementField the isincrement field
	BKIsIncrementField = "is_increment"

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
	// BKGseTaskIdField the gse taskid
	BKGseTaskIdField = "task_id"
	// BKProcPidFile the process pid file
	BKProcPidFile    = "pid_file"
	BKProcStartCmd   = "start_cmd"
	BKProcStopCmd    = "stop_cmd"
	BKProcReloadCmd  = "reload_cmd"
	BKProcRestartCmd = "restart_cmd"
	BKProcTimeOut    = "timeout"
	BKProcWorkPath   = "work_path"
	BKProcInstNum    = "proc_num"

	// BKInstKeyField the inst key field for metric discover
	BKInstKeyField = "bk_inst_key"
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

	// ModifierField the modifier field
	ModifierField = "modifier"

	// LastTimeField the last time field
	LastTimeField = "last_time"
)

const (
	// ValidCreate valid create
	ValidCreate = "create"

	// ValidUpdate valid update
	ValidUpdate = "update"
)

// DefaultResSetFlag the default resource set flat
const DefaultResSetFlag int = 1

// DefaultAppFlag the default app flag
const DefaultAppFlag int = 1

// DefaultAppName the default app name
const DefaultAppName string = "资源池"

// BKAppName the default app name
const BKAppName string = "蓝鲸"

const BKMainLine = "mainline"

const (
	// DefaultResModuleFlag the default resource module flag
	DefaultResModuleFlag int = 1

	// DefaultFaultModuleFlag the default fault module flag
	DefaultFaultModuleFlag int = 2
)
const (
	// FieldTypeSingleChar the single char filed type
	FieldTypeSingleChar string = "singlechar"

	// FieldTypeLongChar the long char field type
	FieldTypeLongChar string = "longchar"

	// FieldTypeInt the int field type
	FieldTypeInt string = "int"

	// FieldTypeEnum the enum field type
	FieldTypeEnum string = "enum"

	// FieldTypeDate the date field type
	FieldTypeDate string = "date"

	// FieldTypeTime the time field type
	FieldTypeTime string = "time"

	// FieldTypeUser the user field type
	FieldTypeUser string = "objuser"

	// FieldTypeSingleAsst the single association
	FieldTypeSingleAsst string = "singleasst"

	// FieldTypeMultiAsst the multi association
	FieldTypeMultiAsst string = "multiasst"

	// FieldTypeTimeZone the timezone field type
	FieldTypeTimeZone string = "timezone"

	// FieldTypeBool the bool type
	FieldTypeBool string = "bool"

	// FieldTypeSingleLenChar the single char length limit
	FieldTypeSingleLenChar int = 256

	// FieldTypeLongLenChar the long char length limit
	FieldTypeLongLenChar int = 2000
)

const (
	// HostAddMethodExcel add a host method
	HostAddMethodExcel = "1"

	// HostAddMethodAgent add a  agent method
	HostAddMethodAgent = "2"

	// HostAddMethodAPI add api method
	HostAddMethodAPI = "3"

	// HostAddMethodExcelIndexOffset the height of the table header
	HostAddMethodExcelIndexOffset = 3

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

	// ExcelAsstPrimaryKeySplitChar split char
	ExcelAsstPrimaryKeySplitChar = "##"
	// ExcelAsstPrimaryKeyRowChar split char
	ExcelAsstPrimaryKeyRowChar = "\n"

	// ExcelDelAsstObjectRelation delete asst object relation
	ExcelDelAsstObjectRelation = "/"

	// ExcelDataValidationListLen excel dropdown list item count
	ExcelDataValidationListLen = 50

	// ExcelCommentSheetCotentLangPrefixKey excel comment sheet centent language prefixe key
	ExcelCommentSheetCotentLangPrefixKey = "import_comment"
)

const (
	// InputTypeExcel  data from excel
	InputTypeExcel = "excel"

	// InputTypeApiHostSync data from api for synchronize new host
	InputTypeApiNewHostSync = "api_sync_host"

	// BatchHostAddMaxRow batch sync add host max row
	BatchHostAddMaxRow = 128
)

const (
	// HTTPBKAPIErrorMessage apiserver error message
	HTTPBKAPIErrorMessage = "bk_error_msg"

	// HTTPBKAPIErrorCode apiserver error code
	HTTPBKAPIErrorCode = "bk_error_code"
)

// KvMap the map definition
type KvMap map[string]interface{}

const (
	// CCSystemOperatorUserName the system user
	CCSystemOperatorUserName  = "cc_system"
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
	EventCacheEventIDKey    = BKCacheKeyV3Prefix + "event:inst_id"
	EventCacheEventQueueKey = BKCacheKeyV3Prefix + "event:inst_queue"
	RedisSnapKeyPrefix      = BKCacheKeyV3Prefix + "snapshot:"
)

const (
	// LocalHostName the local host name definition
	LocalHostName = "localhost"

	// LocalHostIP the local host ip definition
	LocalHostIP = "127.0.0.1"
)

const (
	// BKHTTPHeaderUser current request http request header fields name for login user
	BKHTTPHeaderUser = "BK_User"
	// BKHTTPLanguage the language key word
	BKHTTPLanguage = "HTTP_BLUEKING_LANGUAGE"
	// BKHTTPOwnerID the owner id
	BKHTTPOwnerID = "HTTP_BLUEKING_SUPPLIER_ID"
	//BKHTTPOwnerID = "HTTP_BLUEKING_OWNERID"
	BKHTTPCookieLanugageKey = "blueking_language"
	BKSessionLanugageKey    = "language"

	BKHTTPCCRequestID = "rid"
)

const (
	// DefaultAppLifeCycleNormal  biz life cycle normal
	DefaultAppLifeCycleNormal = "2"
)

// Host OS type enumeration value
const (
	HostOSTypeEnumLinux   = "1"
	HostOSTypeEnumWindows = "2"
)

const (
	// InstAsstIDSplit instance associated fields, separators between mulitple ids of associated objects
	InstAsstIDSplit = ","
)

// integer const
const (
	MaxUint64 = ^uint64(0)
	MinUint64 = 0
	MaxInt64  = int64(MaxUint64 >> 1)
	MinInt64  = -MaxInt64 - 1
	MaxUint   = ^uint(0)
	MinUint   = 0
	MaxInt    = int(MaxUint >> 1)
	MinInt    = -MaxInt - 1
)

//flag
const HostCrossBizField = "hostcrossbiz"
const HostCrossBizValue = "e76fd4d1683d163e4e7e79cef45a74c1"

const (
	BKHTTPMIMEJSON = "application/json"
)

const (
	// APPConfigWaitTime application wait config from zookeeper time (unit sencend)
	APPConfigWaitTime = 15
)

const (
	// URLFilterWhiteList url filter white list not execute any filter
	// multiple url separeted by commas
	URLFilterWhiteListSuffix = "/healthz"

	URLFilterWhiteListSepareteChar = ","
)

type DataStatusFlag string

const (
	DataStatusDisabled DataStatusFlag = "disabled"
	DataStatusEnable   DataStatusFlag = "enable"
)

const (
	BKDataStatusField   = "bk_data_status"
	BKDataRecoverSuffix = "(recover)"
)
