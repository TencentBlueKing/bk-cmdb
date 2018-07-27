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

package v3v0v8

import (
	"configcenter/src/common"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/validator"
)

// default group
var (
	groupBaseInfo = mCommon.BaseInfo
)

// Distribution init revision
var Distribution = "community" // could be community or enterprise

/*
	&metadata.Attribute{ObjectID: objID, PropertyID: "", PropertyName: "", IsRequired: , IsOnly: , PropertyGroup: , PropertyType: , Option: ""},
*/

// AppRow app structure
func AppRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDApp

	groupAppRole := mCommon.AppRole

	lifeCycleOption := []validator.EnumVal{{ID: "1", Name: "测试中", Type: "text"}, {ID: "2", Name: "已上线", Type: "text", IsDefault: true}, {ID: "3", Name: "停运", Type: "text"}}
	languageOption := []validator.EnumVal{{ID: "1", Name: "中文", Type: "text", IsDefault: true}, {ID: "2", Name: "English", Type: "text"}}
	dataRows := []*metadata.Attribute{
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_biz_name", PropertyName: "业务名", IsRequired: true, IsOnly: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "life_cycle", PropertyName: "生命周期", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: lifeCycleOption},

		//role
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKMaintainersField, PropertyName: "运维人员", IsRequired: true, IsOnly: false, IsEditable: true, PropertyGroup: groupAppRole, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKProductPMField, PropertyName: "产品人员", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupAppRole, PropertyType: common.FieldTypeUser, Option: ""},

		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKTesterField, PropertyName: "测试人员", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupAppRole, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_biz_developer", PropertyName: "开发人员", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupAppRole, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKOperatorField, PropertyName: "操作人员", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupAppRole, PropertyType: common.FieldTypeUser, Option: ""},

		&metadata.Attribute{ObjectID: objID, PropertyID: "time_zone", PropertyName: "时区", IsRequired: true, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTimeZone, Option: "", IsReadOnly: true},
		&metadata.Attribute{ObjectID: objID, PropertyID: "language", PropertyName: "语言", IsRequired: true, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: languageOption, IsReadOnly: true},
	}

	return dataRows

}

// SetRow set structure
func SetRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDSet
	serviceStatusOption := []validator.EnumVal{{ID: "1", Name: "开放", Type: "text", IsDefault: true}, {ID: "2", Name: "关闭", Type: "text"}}

	dataRows := []*metadata.Attribute{
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: validator.IntOption{}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_set_name", PropertyName: "集群名字", IsRequired: true, IsOnly: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_set_desc", PropertyName: "集群描述", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_set_env", PropertyName: "环境类型", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "测试", Type: "text"}, {ID: "2", Name: "体验", Type: "text"}, {ID: "3", Name: "正式", Type: "text", IsDefault: true}}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_service_status", PropertyName: "服务状态", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: serviceStatusOption},
		&metadata.Attribute{ObjectID: objID, PropertyID: "description", PropertyName: "备注", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_capacity", PropertyName: "设计容量", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: validator.IntOption{Min: "1", Max: "999999999"}},

		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKChildStr, PropertyName: "", IsRequired: false, IsOnly: false, IsSystem: true, PropertyType: "", Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKInstParentStr, PropertyName: "", IsSystem: true, IsRequired: true, IsOnly: true, PropertyType: common.FieldTypeInt, Option: validator.IntOption{}},
	}
	return dataRows
}

// ModuleRow module structure
func ModuleRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDModule
	moduleTypeOption := []validator.EnumVal{{ID: "1", Name: "普通", Type: "text", IsDefault: true}, {ID: "2", Name: "数据库", Type: "text"}}

	dataRows := []*metadata.Attribute{
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: validator.IntOption{}},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKSetIDField, PropertyName: "集群ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: validator.IntOption{}},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKModuleNameField, PropertyName: "模块名", IsRequired: true, IsOnly: true, IsEditable: true, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKChildStr, PropertyName: "", IsRequired: false, IsOnly: false, IsSystem: true, PropertyType: "", Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_module_type", PropertyName: "模块类型", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: moduleTypeOption},
		&metadata.Attribute{ObjectID: objID, PropertyID: "operator", PropertyName: "主要维护人", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_bak_operator", PropertyName: "备份维护人", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeUser, Option: ""},
	}
	return dataRows
}

// PlatRow plat structure
func PlatRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDPlat
	dataRows := []*metadata.Attribute{
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKCloudNameField, PropertyName: "云区域", IsRequired: true, IsOnly: true, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKOwnerIDField, PropertyName: "供应商", IsRequired: true, IsOnly: true, IsPre: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
	}
	return dataRows
}

// HostRow host structure
func HostRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDHost
	groupAgent := mCommon.HostAutoFields
	dataRows := []*metadata.Attribute{
		//基本信息分组
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKHostInnerIPField, PropertyName: "内网IP", IsRequired: true, IsOnly: true, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKHostOuterIPField, PropertyName: "外网IP", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.Attribute{ObjectID: objID, PropertyID: "operator", PropertyName: "主要维护人", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_bak_operator", PropertyName: "备份维护人", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeUser, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_asset_id", PropertyName: "固资编号", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_sn", PropertyName: "设备SN", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_comment", PropertyName: "备注", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_service_term", PropertyName: "质保年限", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: validator.IntOption{Min: "1", Max: "10"}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_sla", PropertyName: "SLA级别", IsRequired: false, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "L1", Type: "text"}, {ID: "2", Name: "L2", Type: "text"}, {ID: "3", Name: "L3", Type: "text"}}},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKCloudIDField, PropertyName: "云区域", IsRequired: false, IsOnly: true, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleAsst, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_state_name", PropertyName: "所在国家", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: stateEnum},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_province_name", PropertyName: "所在省份", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: provincesEnum},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_isp_name", PropertyName: "所属运营商", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: ispNameEnum},

		//自动发现分组
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_host_name", PropertyName: "主机名称", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKOSTypeField, PropertyName: "操作系统类型", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "Linux", Type: "text"}, {ID: "2", Name: "Windows", Type: "text"}}},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKOSNameField, PropertyName: "操作系统名称", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_os_version", PropertyName: "操作系统版本", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_os_bit", PropertyName: "操作系统位数", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_cpu", PropertyName: "CPU逻辑核心数", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeInt, Option: validator.IntOption{Min: "1", Max: "1000000"}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_cpu_mhz", PropertyName: "CPU频率", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeInt, Unit: "Hz", Option: validator.IntOption{Min: "1", Max: "100000000"}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_cpu_module", PropertyName: "CPU型号", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_mem", PropertyName: "内存容量", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeInt, Unit: "MB", Option: validator.IntOption{Min: "1", Max: "100000000"}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_disk", PropertyName: "磁盘容量", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeInt, Unit: "GB", Option: validator.IntOption{Min: "1", Max: "100000000"}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_mac", PropertyName: "内网MAC地址", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_outer_mac", PropertyName: "外网MAC", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FieldTypeSingleChar, Option: ""},
		//agent 没有分组
		&metadata.Attribute{ObjectID: objID, PropertyID: common.CreateTimeField, PropertyName: "录入时间", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeTime, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "import_from", PropertyName: "录入方式", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "excel", Type: "text"}, {ID: "2", Name: "agent", Type: "text"}, {ID: "3", Name: "api", Type: "text"}}},
	}

	return dataRows
}

// ProcRow proc structure
func ProcRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDProc
	groupPort := mCommon.ProcPort
	// groupGsekit := mCommon.Proc_gsekit_base_info
	// groupGsekitManage := mCommon.Proc_gsekit_manage_info
	dataRows := []*metadata.Attribute{
		//base info
		//&metadata.Attribute{ObjectID: objID, PropertyID: "bk_process_id", PropertyName: "进程ID", IsSystem: true, IsRequired: true, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: "{}"},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: true, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: validator.IntOption{}},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKProcessNameField, PropertyName: "进程名称", IsRequired: true, IsOnly: true, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "description", PropertyName: "进程描述", IsRequired: false, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},

		//监听端口分组
		&metadata.Attribute{ObjectID: objID, PropertyID: "bind_ip", PropertyName: "绑定IP", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupPort, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "127.0.0.1", Type: "text"}, {ID: "2", Name: "0.0.0.0", Type: "text"}, {ID: "3", Name: "第一内网IP", Type: "text"}, {ID: "4", Name: "第一外网IP", Type: "text"}}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "port", PropertyName: "端口", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupPort, PropertyType: common.FieldTypeSingleChar, Option: common.PatternMultiplePortRange, Placeholder: `单个端口：8080 </br>多个连续端口：8080-8089 </br>多个不连续端口：8080-8089,8199`},
		&metadata.Attribute{ObjectID: objID, PropertyID: "protocol", PropertyName: "协议", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupPort, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "TCP", Type: "text"}, {ID: "2", Name: "UDP", Type: "text"}}},

		//gsekit 基础信息
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_func_id", PropertyName: "功能ID", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_func_name", PropertyName: "功能名称", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "work_path", PropertyName: "工作路径", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "user", PropertyName: "启动用户", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "proc_num", PropertyName: "启动数量", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeInt, Option: validator.IntOption{Min: "1", Max: "1000000"}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "priority", PropertyName: "启动优先级", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeInt, Option: validator.IntOption{Min: "1", Max: "100"}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "timeout", PropertyName: "操作超时时长", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeInt, Option: validator.IntOption{Min: "1", Max: "1000000"}},

		//gsekit 进程信息
		&metadata.Attribute{ObjectID: objID, PropertyID: "start_cmd", PropertyName: "启动命令", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "stop_cmd", PropertyName: "停止命令", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "restart_cmd", PropertyName: "重启命令", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "face_stop_cmd", PropertyName: "强制停止命令", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "reload_cmd", PropertyName: "进程重载命令", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "pid_file", PropertyName: "PID文件路径", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "auto_start", PropertyName: "是否自动拉起", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeBool, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "auto_time_gap", PropertyName: "拉起间隔", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FieldTypeInt, Option: validator.IntOption{Min: "1", Max: "1000000"}},
	}
	return dataRows
}

// SwitchRow proc structure
func SwitchRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDSwitch
	dataRows := []*metadata.Attribute{
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_asset_id", PropertyName: "固资编号", IsRequired: true, IsOnly: true, IsPre: true, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_inst_name", PropertyName: "名称", IsRequired: true, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_sn", PropertyName: "SN", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_func", PropertyName: "用途", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_vendor", PropertyName: "厂商", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_model", PropertyName: "设备型号", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_admin_ip", PropertyName: "管理IP", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_operator", PropertyName: "维护人", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_os_detail", PropertyName: "操作系统详情", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_detail", PropertyName: "详细描述", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_biz_status", PropertyName: "运营状态", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "待运营", Type: "text"}, {ID: "2", Name: "运营中", Type: "text"}, {ID: "3", Name: "已下架", Type: "text"}}},
	}
	return dataRows
}

// RouterRow proc structure
func RouterRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDRouter
	dataRows := []*metadata.Attribute{
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_asset_id", PropertyName: "固资编号", IsRequired: true, IsOnly: true, IsPre: true, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_inst_name", PropertyName: "名称", IsRequired: true, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_sn", PropertyName: "SN", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_func", PropertyName: "用途", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_vendor", PropertyName: "厂商", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_model", PropertyName: "设备型号", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_admin_ip", PropertyName: "管理IP", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_operator", PropertyName: "维护人", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_os_detail", PropertyName: "操作系统详情", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_detail", PropertyName: "详细描述", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_biz_status", PropertyName: "运营状态", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "待运营", Type: "text"}, {ID: "2", Name: "运营中", Type: "text"}, {ID: "3", Name: "已下架", Type: "text"}}},
	}
	return dataRows
}

// LoadBalanceRow proc structure
func LoadBalanceRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDBlance
	dataRows := []*metadata.Attribute{
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_asset_id", PropertyName: "固资编号", IsRequired: true, IsOnly: true, IsPre: true, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_inst_name", PropertyName: "名称", IsRequired: true, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_sn", PropertyName: "SN", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_func", PropertyName: "用途", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_vendor", PropertyName: "厂商", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_model", PropertyName: "设备型号", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_admin_ip", PropertyName: "管理IP", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_operator", PropertyName: "维护人", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_os_detail", PropertyName: "操作系统详情", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_detail", PropertyName: "详细描述", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_biz_status", PropertyName: "运营状态", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "待运营", Type: "text"}, {ID: "2", Name: "运营中", Type: "text"}, {ID: "3", Name: "已下架", Type: "text"}}},
	}
	return dataRows
}

// FirewallRow proc structure
func FirewallRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDFirewall
	dataRows := []*metadata.Attribute{
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_asset_id", PropertyName: "固资编号", IsRequired: true, IsOnly: true, IsPre: true, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_inst_name", PropertyName: "名称", IsRequired: true, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_sn", PropertyName: "SN", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_func", PropertyName: "用途", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_vendor", PropertyName: "厂商", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_model", PropertyName: "设备型号", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_admin_ip", PropertyName: "管理IP", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_operator", PropertyName: "维护人", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_os_detail", PropertyName: "操作系统详情", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_detail", PropertyName: "详细描述", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_biz_status", PropertyName: "运营状态", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "待运营", Type: "text"}, {ID: "2", Name: "运营中", Type: "text"}, {ID: "3", Name: "已下架", Type: "text"}}},
	}
	return dataRows
}

// WeblogicRow proc structure
func WeblogicRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDWeblogic
	dataRows := []*metadata.Attribute{
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_inst_key", PropertyName: "中间件标识", IsRequired: true, IsOnly: true, IsPre: true, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_inst_name", PropertyName: "名称", IsRequired: true, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_version", PropertyName: "版本", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_patch_version", PropertyName: "补丁版本", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_main_path", PropertyName: "主目录", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_log_path", PropertyName: "日志路径", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_vendor", PropertyName: "厂商", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_ip", PropertyName: "IP地址", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_port", PropertyName: "端口", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_detail", PropertyName: "详细描述", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_jdk_version", PropertyName: "JDK版本", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_jvm_free_mem", PropertyName: "JVM配置的最大空闲内存", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_jvm_capacity", PropertyName: "JVM堆的当前大小", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_jvm_used_mem", PropertyName: "JVM堆的当前可用的内存", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
	}
	return dataRows
}

// TomcatRow proc structure
func TomcatRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDTomcat
	dataRows := []*metadata.Attribute{
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_inst_key", PropertyName: "中间件标识", IsRequired: true, IsOnly: true, IsPre: true, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_inst_name", PropertyName: "名称", IsRequired: true, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_version", PropertyName: "版本", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_patch_version", PropertyName: "补丁版本", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_main_path", PropertyName: "主目录", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_log_path", PropertyName: "日志路径", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_vendor", PropertyName: "厂商", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_ip", PropertyName: "IP地址", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_port", PropertyName: "端口", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_detail", PropertyName: "详细描述", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_jdk_version", PropertyName: "JDK版本", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
	}
	return dataRows
}

// ApacheRow proc structure
func ApacheRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDApache
	dataRows := []*metadata.Attribute{
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_inst_key", PropertyName: "中间件标识", IsRequired: true, IsOnly: true, IsPre: true, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_inst_name", PropertyName: "名称", IsRequired: true, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_version", PropertyName: "版本", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_patch_version", PropertyName: "补丁版本", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_main_path", PropertyName: "主目录", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_log_path", PropertyName: "日志路径", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_vendor", PropertyName: "厂商", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_ip", PropertyName: "IP地址", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_port", PropertyName: "端口", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_detail", PropertyName: "详细描述", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_max_connect", PropertyName: "最大连接请求数", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: validator.IntOption{}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "bk_max_keepalive", PropertyName: "最大keepAlive请求数", IsRequired: false, IsOnly: false, IsPre: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: validator.IntOption{}},
	}
	return dataRows
}

var stateEnum = []validator.EnumVal{
	{ID: "AR", Name: "阿根廷", Type: "text"},
	{ID: "AD", Name: "安道尔", Type: "text"},
	{ID: "AE", Name: "阿联酋", Type: "text"},
	{ID: "AF", Name: "阿富汗", Type: "text"},
	{ID: "AG", Name: "安提瓜和巴布达", Type: "text"},
	{ID: "AI", Name: "安圭拉", Type: "text"},
	{ID: "AL", Name: "阿尔巴尼亚", Type: "text"},
	{ID: "AM", Name: "亚美尼亚", Type: "text"},
	{ID: "AO", Name: "安哥拉", Type: "text"},
	{ID: "AQ", Name: "南极洲", Type: "text"},
	{ID: "AS", Name: "美属萨摩亚", Type: "text"},
	{ID: "AT", Name: "奥地利", Type: "text"},
	{ID: "AU", Name: "澳大利亚", Type: "text"},
	{ID: "AW", Name: "阿鲁巴", Type: "text"},
	{ID: "AX", Name: "奥兰群岛", Type: "text"},
	{ID: "AZ", Name: "阿塞拜疆", Type: "text"},
	{ID: "BA", Name: "波黑", Type: "text"},
	{ID: "BB", Name: "巴巴多斯", Type: "text"},
	{ID: "BD", Name: "孟加拉", Type: "text"},
	{ID: "BE", Name: "比利时", Type: "text"},
	{ID: "BF", Name: "布基纳法索", Type: "text"},
	{ID: "BG", Name: "保加利亚", Type: "text"},
	{ID: "BH", Name: "巴林", Type: "text"},
	{ID: "BI", Name: "布隆迪", Type: "text"},
	{ID: "BJ", Name: "贝宁", Type: "text"},
	{ID: "BL", Name: "圣巴泰勒米岛", Type: "text"},
	{ID: "BM", Name: "百慕大", Type: "text"},
	{ID: "BN", Name: "文莱", Type: "text"},
	{ID: "BO", Name: "玻利维亚", Type: "text"},
	{ID: "BQ", Name: "荷兰加勒比区", Type: "text"},
	{ID: "BR", Name: "巴西", Type: "text"},
	{ID: "BS", Name: "巴哈马", Type: "text"},
	{ID: "BT", Name: "不丹", Type: "text"},
	{ID: "BV", Name: "布韦岛", Type: "text"},
	{ID: "BW", Name: "博茨瓦纳", Type: "text"},
	{ID: "BY", Name: "白俄罗斯", Type: "text"},
	{ID: "BZ", Name: "伯利兹", Type: "text"},
	{ID: "CA", Name: "加拿大", Type: "text"},
	{ID: "CC", Name: "科科斯群岛", Type: "text"},
	{ID: "CD", Name: "刚果（金）", Type: "text"},
	{ID: "CF", Name: "中非", Type: "text"},
	{ID: "CG", Name: "刚果（布）", Type: "text"},
	{ID: "CH", Name: "瑞士", Type: "text"},
	{ID: "CI", Name: "科特迪瓦", Type: "text"},
	{ID: "CK", Name: "库克群岛", Type: "text"},
	{ID: "CL", Name: "智利", Type: "text"},
	{ID: "CM", Name: "喀麦隆", Type: "text"},
	{ID: "CN", Name: "中国", Type: "text"},
	{ID: "CO", Name: "哥伦比亚", Type: "text"},
	{ID: "CR", Name: "哥斯达黎加", Type: "text"},
	{ID: "CU", Name: "古巴", Type: "text"},
	{ID: "CV", Name: "佛得角", Type: "text"},
	{ID: "CW", Name: "库拉索", Type: "text"},
	{ID: "CX", Name: "圣诞岛", Type: "text"},
	{ID: "CY", Name: "塞浦路斯", Type: "text"},
	{ID: "CZ", Name: "捷克", Type: "text"},
	{ID: "DE", Name: "德国", Type: "text"},
	{ID: "DJ", Name: "吉布提", Type: "text"},
	{ID: "DK", Name: "丹麦", Type: "text"},
	{ID: "DM", Name: "多米尼克", Type: "text"},
	{ID: "DO", Name: "多米尼加", Type: "text"},
	{ID: "DZ", Name: "阿尔及利亚", Type: "text"},
	{ID: "EC", Name: "厄瓜多尔", Type: "text"},
	{ID: "EE", Name: "爱沙尼亚", Type: "text"},
	{ID: "EG", Name: "埃及", Type: "text"},
	{ID: "EH", Name: "西撒哈拉", Type: "text"},
	{ID: "ER", Name: "厄立特里亚", Type: "text"},
	{ID: "ES", Name: "西班牙", Type: "text"},
	{ID: "ET", Name: "埃塞俄比亚", Type: "text"},
	{ID: "FI", Name: "芬兰", Type: "text"},
	{ID: "FJ", Name: "斐济群岛", Type: "text"},
	{ID: "FK", Name: "马尔维纳斯群岛（福克兰）", Type: "text"},
	{ID: "FM", Name: "密克罗尼西亚联邦", Type: "text"},
	{ID: "FO", Name: "法罗群岛", Type: "text"},
	{ID: "FR", Name: "法国 法国", Type: "text"},
	{ID: "GA", Name: "加蓬", Type: "text"},
	{ID: "GB", Name: "英国", Type: "text"},
	{ID: "GD", Name: "格林纳达", Type: "text"},
	{ID: "GE", Name: "格鲁吉亚", Type: "text"},
	{ID: "GF", Name: "法属圭亚那", Type: "text"},
	{ID: "GG", Name: "根西岛", Type: "text"},
	{ID: "GH", Name: "加纳", Type: "text"},
	{ID: "GI", Name: "直布罗陀", Type: "text"},
	{ID: "GL", Name: "格陵兰", Type: "text"},
	{ID: "GM", Name: "冈比亚", Type: "text"},
	{ID: "GN", Name: "几内亚", Type: "text"},
	{ID: "GP", Name: "瓜德罗普", Type: "text"},
	{ID: "GQ", Name: "赤道几内亚", Type: "text"},
	{ID: "GR", Name: "希腊", Type: "text"},
	{ID: "GS", Name: "南乔治亚岛和南桑威奇群岛", Type: "text"},
	{ID: "GT", Name: "危地马拉", Type: "text"},
	{ID: "GU", Name: "关岛", Type: "text"},
	{ID: "GW", Name: "几内亚比绍", Type: "text"},
	{ID: "GY", Name: "圭亚那", Type: "text"},
	{ID: "HM", Name: "赫德岛和麦克唐纳群岛", Type: "text"},
	{ID: "HN", Name: "洪都拉斯", Type: "text"},
	{ID: "HR", Name: "克罗地亚", Type: "text"},
	{ID: "HT", Name: "海地", Type: "text"},
	{ID: "HU", Name: "匈牙利", Type: "text"},
	{ID: "ID", Name: "印尼", Type: "text"},
	{ID: "IE", Name: "爱尔兰", Type: "text"},
	{ID: "IL", Name: "以色列", Type: "text"},
	{ID: "IM", Name: "马恩岛", Type: "text"},
	{ID: "IN", Name: "印度", Type: "text"},
	{ID: "IO", Name: "英属印度洋领地", Type: "text"},
	{ID: "IQ", Name: "伊拉克", Type: "text"},
	{ID: "IR", Name: "伊朗", Type: "text"},
	{ID: "IS", Name: "冰岛", Type: "text"},
	{ID: "IT", Name: "意大利", Type: "text"},
	{ID: "JE", Name: "泽西岛", Type: "text"},
	{ID: "JM", Name: "牙买加", Type: "text"},
	{ID: "JO", Name: "约旦", Type: "text"},
	{ID: "JP", Name: "日本", Type: "text"},
	{ID: "KE", Name: "肯尼亚", Type: "text"},
	{ID: "KG", Name: "吉尔吉斯斯坦", Type: "text"},
	{ID: "KH", Name: "柬埔寨", Type: "text"},
	{ID: "KI", Name: "基里巴斯", Type: "text"},
	{ID: "KM", Name: "科摩罗", Type: "text"},
	{ID: "KN", Name: "圣基茨和尼维斯", Type: "text"},
	{ID: "KP", Name: "朝鲜", Type: "text"},
	{ID: "KR", Name: "韩国", Type: "text"},
	{ID: "KW", Name: "科威特", Type: "text"},
	{ID: "KY", Name: "开曼群岛", Type: "text"},
	{ID: "KZ", Name: "哈萨克斯坦", Type: "text"},
	{ID: "LA", Name: "老挝", Type: "text"},
	{ID: "LB", Name: "黎巴嫩", Type: "text"},
	{ID: "LC", Name: "圣卢西亚", Type: "text"},
	{ID: "LI", Name: "列支敦士登", Type: "text"},
	{ID: "LK", Name: "斯里兰卡", Type: "text"},
	{ID: "LR", Name: "利比里亚", Type: "text"},
	{ID: "LS", Name: "莱索托", Type: "text"},
	{ID: "LT", Name: "立陶宛", Type: "text"},
	{ID: "LU", Name: "卢森堡", Type: "text"},
	{ID: "LV", Name: "拉脱维亚", Type: "text"},
	{ID: "LY", Name: "利比亚", Type: "text"},
	{ID: "MA", Name: "摩洛哥", Type: "text"},
	{ID: "MC", Name: "摩纳哥", Type: "text"},
	{ID: "MD", Name: "摩尔多瓦", Type: "text"},
	{ID: "ME", Name: "黑山", Type: "text"},
	{ID: "MF", Name: "法属圣马丁", Type: "text"},
	{ID: "MG", Name: "马达加斯加", Type: "text"},
	{ID: "MH", Name: "马绍尔群岛", Type: "text"},
	{ID: "MK", Name: "马其顿", Type: "text"},
	{ID: "ML", Name: "马里", Type: "text"},
	{ID: "MM", Name: "缅甸", Type: "text"},
	{ID: "MN", Name: "蒙古国", Type: "text"},
	{ID: "MP", Name: "北马里亚纳群岛", Type: "text"},
	{ID: "MQ", Name: "马提尼克", Type: "text"},
	{ID: "MR", Name: "毛里塔尼亚", Type: "text"},
	{ID: "MS", Name: "蒙塞拉特岛", Type: "text"},
	{ID: "MT", Name: "马耳他", Type: "text"},
	{ID: "MU", Name: "毛里求斯", Type: "text"},
	{ID: "MV", Name: "马尔代夫", Type: "text"},
	{ID: "MW", Name: "马拉维", Type: "text"},
	{ID: "MX", Name: "墨西哥", Type: "text"},
	{ID: "MY", Name: "马来西亚", Type: "text"},
	{ID: "MZ", Name: "莫桑比克", Type: "text"},
	{ID: "NA", Name: "纳米比亚", Type: "text"},
	{ID: "NC", Name: "新喀里多尼亚", Type: "text"},
	{ID: "NE", Name: "尼日尔", Type: "text"},
	{ID: "NF", Name: "诺福克岛", Type: "text"},
	{ID: "NG", Name: "尼日利亚", Type: "text"},
	{ID: "NI", Name: "尼加拉瓜", Type: "text"},
	{ID: "NL", Name: "荷兰", Type: "text"},
	{ID: "NO", Name: "挪威", Type: "text"},
	{ID: "NP", Name: "尼泊尔", Type: "text"},
	{ID: "NR", Name: "瑙鲁", Type: "text"},
	{ID: "NU", Name: "纽埃", Type: "text"},
	{ID: "NZ", Name: "新西兰", Type: "text"},
	{ID: "OM", Name: "阿曼", Type: "text"},
	{ID: "PA", Name: "巴拿马", Type: "text"},
	{ID: "PE", Name: "秘鲁", Type: "text"},
	{ID: "PF", Name: "法属波利尼西亚", Type: "text"},
	{ID: "PG", Name: "巴布亚新几内亚", Type: "text"},
	{ID: "PH", Name: "菲律宾", Type: "text"},
	{ID: "PK", Name: "巴基斯坦", Type: "text"},
	{ID: "PL", Name: "波兰", Type: "text"},
	{ID: "PM", Name: "圣皮埃尔和密克隆", Type: "text"},
	{ID: "PN", Name: "皮特凯恩群岛", Type: "text"},
	{ID: "PR", Name: "波多黎各", Type: "text"},
	{ID: "PS", Name: "巴勒斯坦", Type: "text"},
	{ID: "PT", Name: "葡萄牙", Type: "text"},
	{ID: "PW", Name: "帕劳", Type: "text"},
	{ID: "PY", Name: "巴拉圭", Type: "text"},
	{ID: "QA", Name: "卡塔尔", Type: "text"},
	{ID: "RE", Name: "留尼汪", Type: "text"},
	{ID: "RO", Name: "罗马尼亚", Type: "text"},
	{ID: "RS", Name: "塞尔维亚", Type: "text"},
	{ID: "RU", Name: "俄罗斯", Type: "text"},
	{ID: "RW", Name: "卢旺达", Type: "text"},
	{ID: "SA", Name: "沙特阿拉伯", Type: "text"},
	{ID: "SB", Name: "所罗门群岛", Type: "text"},
	{ID: "SC", Name: "塞舌尔", Type: "text"},
	{ID: "SD", Name: "苏丹", Type: "text"},
	{ID: "SE", Name: "瑞典", Type: "text"},
	{ID: "SG", Name: "新加坡", Type: "text"},
	{ID: "SH", Name: "圣赫勒拿", Type: "text"},
	{ID: "SI", Name: "斯洛文尼亚", Type: "text"},
	{ID: "SJ", Name: "斯瓦尔巴群岛和扬马延岛", Type: "text"},
	{ID: "SK", Name: "斯洛伐克", Type: "text"},
	{ID: "SL", Name: "塞拉利昂", Type: "text"},
	{ID: "SM", Name: "圣马力诺", Type: "text"},
	{ID: "SN", Name: "塞内加尔", Type: "text"},
	{ID: "SO", Name: "索马里", Type: "text"},
	{ID: "SR", Name: "苏里南", Type: "text"},
	{ID: "SS", Name: "南苏丹", Type: "text"},
	{ID: "ST", Name: "圣多美和普林西比", Type: "text"},
	{ID: "SV", Name: "萨尔瓦多", Type: "text"},
	{ID: "SX", Name: "荷属圣马丁", Type: "text"},
	{ID: "SY", Name: "叙利亚", Type: "text"},
	{ID: "SZ", Name: "斯威士兰", Type: "text"},
	{ID: "TC", Name: "特克斯和凯科斯群岛", Type: "text"},
	{ID: "TD", Name: "乍得", Type: "text"},
	{ID: "TF", Name: "法属南部领地", Type: "text"},
	{ID: "TG", Name: "多哥", Type: "text"},
	{ID: "TH", Name: "泰国", Type: "text"},
	{ID: "TJ", Name: "塔吉克斯坦", Type: "text"},
	{ID: "TK", Name: "托克劳", Type: "text"},
	{ID: "TL", Name: "东帝汶", Type: "text"},
	{ID: "TM", Name: "土库曼斯坦", Type: "text"},
	{ID: "TN", Name: "突尼斯", Type: "text"},
	{ID: "TO", Name: "汤加", Type: "text"},
	{ID: "TR", Name: "土耳其", Type: "text"},
	{ID: "TT", Name: "特立尼达和多巴哥", Type: "text"},
	{ID: "TV", Name: "图瓦卢", Type: "text"},
	{ID: "TZ", Name: "坦桑尼亚", Type: "text"},
	{ID: "UA", Name: "乌克兰", Type: "text"},
	{ID: "UG", Name: "乌干达", Type: "text"},
	{ID: "UM", Name: "美国本土外小岛屿", Type: "text"},
	{ID: "UY", Name: "乌拉圭", Type: "text"},
	{ID: "UZ", Name: "乌兹别克斯坦", Type: "text"},
	{ID: "VA", Name: "梵蒂冈", Type: "text"},
	{ID: "VC", Name: "圣文森特和格林纳丁斯", Type: "text"},
	{ID: "VE", Name: "委内瑞拉", Type: "text"},
	{ID: "VG", Name: "英属维尔京群岛", Type: "text"},
	{ID: "VI", Name: "美属维尔京群岛", Type: "text"},
	{ID: "VN", Name: "越南", Type: "text"},
	{ID: "US", Name: "美国", Type: "text"},
	{ID: "VU", Name: "瓦努阿图", Type: "text"},
	{ID: "WF", Name: "瓦利斯和富图纳", Type: "text"},
	{ID: "WS", Name: "萨摩亚", Type: "text"},
	{ID: "YE", Name: "也门", Type: "text"},
	{ID: "YT", Name: "马约特", Type: "text"},
	{ID: "ZA", Name: "南非", Type: "text"},
	{ID: "ZM", Name: "赞比亚", Type: "text"},
	{ID: "ZW", Name: "津巴布韦", Type: "text"},
}

var provincesEnum = []validator.EnumVal{
	{ID: "110000", Name: "北京市", Type: "text"},
	{ID: "120000", Name: "天津市", Type: "text"},
	{ID: "130000", Name: "河北省", Type: "text"},
	{ID: "140000", Name: "山西省", Type: "text"},
	{ID: "150000", Name: "内蒙古自治区", Type: "text"},
	{ID: "210000", Name: "辽宁省", Type: "text"},
	{ID: "220000", Name: "吉林省", Type: "text"},
	{ID: "230000", Name: "黑龙江省", Type: "text"},
	{ID: "310000", Name: "上海市", Type: "text"},
	{ID: "320000", Name: "江苏省", Type: "text"},
	{ID: "330000", Name: "浙江省", Type: "text"},
	{ID: "340000", Name: "安徽省", Type: "text"},
	{ID: "350000", Name: "福建省", Type: "text"},
	{ID: "360000", Name: "江西省", Type: "text"},
	{ID: "370000", Name: "山东省", Type: "text"},
	{ID: "410000", Name: "河南省", Type: "text"},
	{ID: "420000", Name: "湖北省", Type: "text"},
	{ID: "430000", Name: "湖南省", Type: "text"},
	{ID: "440000", Name: "广东省", Type: "text"},
	{ID: "450000", Name: "广西壮族自治区", Type: "text"},
	{ID: "460000", Name: "海南省", Type: "text"},
	{ID: "500000", Name: "重庆市", Type: "text"},
	{ID: "510000", Name: "四川省", Type: "text"},
	{ID: "520000", Name: "贵州省", Type: "text"},
	{ID: "530000", Name: "云南省", Type: "text"},
	{ID: "540000", Name: "西藏自治区", Type: "text"},
	{ID: "610000", Name: "陕西省", Type: "text"},
	{ID: "620000", Name: "甘肃省", Type: "text"},
	{ID: "630000", Name: "青海省", Type: "text"},
	{ID: "640000", Name: "宁夏回族自治区", Type: "text"},
	{ID: "650000", Name: "新疆维吾尔自治区", Type: "text"},
	{ID: "710000", Name: "台湾省", Type: "text"},
	{ID: "810000", Name: "香港特别行政区", Type: "text"},
	{ID: "820000", Name: "澳门特别行政区", Type: "text"},
}

var ispNameEnum = []validator.EnumVal{
	{ID: "0", Name: "其他", Type: "text"},
	{ID: "1", Name: "电信", Type: "text"},
	{ID: "2", Name: "联通", Type: "text"}, {ID: "3", Name: "移动", Type: "text"}}
