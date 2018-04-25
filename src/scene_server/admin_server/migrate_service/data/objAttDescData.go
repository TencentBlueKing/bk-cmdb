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

package data

import (
	"configcenter/src/common"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/validator"
	"configcenter/src/source_controller/api/metadata"
)

// default group
var (
	groupBaseInfo = mCommon.BaseInfo
)

// Distribution init revision
var Distribution = "community" // could be community or enterprise

/*
	&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "", PropertyName: "", IsRequired: , IsOnly: , PropertyGroup: , PropertyType: , Option: ""},
*/

// AppRow app structure
func AppRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDApp

	groupAppRole := mCommon.AppRole

	lifeCycleOption := []validator.EnumVal{{ID: "1", Name: "测试中", Type: "text"}, {ID: "2", Name: "已上线", Type: "text", IsDefault: true}, {ID: "3", Name: "停运", Type: "text"}}
	languageOption := []validator.EnumVal{{ID: "1", Name: "中文", Type: "text"}, {ID: "2", Name: "English", Type: "text"}}
	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_biz_name", PropertyName: "业务名", IsRequired: true, IsOnly: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "life_cycle", PropertyName: "生命周期", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: lifeCycleOption},

		//role
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKMaintainersField, PropertyName: "运维人员", IsRequired: true, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FiledTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKProductPMField, PropertyName: "产品人员", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FiledTypeUser, Option: ""},

		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKTesterField, PropertyName: "测试人员", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FiledTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_biz_developer", PropertyName: "开发人员", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FiledTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKOperatorField, PropertyName: "操作人员", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FiledTypeUser, Option: ""},
	}

	if Distribution == common.RevisionEnterprise {
		dataRows = append(dataRows,
			&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "time_zone", PropertyName: "时区", IsRequired: true, IsOnly: false, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTimeZone, Option: "", IsReadOnly: true},
			&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "language", PropertyName: "语言", IsRequired: true, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: languageOption, IsReadOnly: true},
		)
	}

	return dataRows

}

// SetRow set structure
func SetRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDSet
	serviceStatusOption := []validator.EnumVal{{ID: "1", Name: "开放", Type: "text", IsDefault: true}, {ID: "2", Name: "关闭", Type: "text"}}

	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: common.KvMap{}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_set_name", PropertyName: "集群名字", IsRequired: true, IsOnly: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_set_desc", PropertyName: "集群描述", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_set_env", PropertyName: "环境类型", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "测试", Type: "text"}, {ID: "2", Name: "体验", Type: "text"}, {ID: "3", Name: "正式", Type: "text", IsDefault: true}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_service_status", PropertyName: "服务状态", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: serviceStatusOption},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "description", PropertyName: "备注", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_capacity", PropertyName: "设计容量", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: common.KvMap{"min": "1", "max": "999999999"}},

		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKChildStr, PropertyName: "", IsRequired: false, IsOnly: false, PropertyType: "", Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKInstParentStr, PropertyName: "", IsSystem: true, IsRequired: true, IsOnly: true, PropertyType: common.FiledTypeInt, Option: ""},
	}
	return dataRows
}

// ModuleRow module structure
func ModuleRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDModule
	moduleTypeOption := []validator.EnumVal{{ID: "1", Name: "普通", Type: "text", IsDefault: true}, {ID: "2", Name: "数据库", Type: "text"}}

	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: common.KvMap{}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKSetIDField, PropertyName: "集群ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: common.KvMap{}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKModuleNameField, PropertyName: "模块名", IsRequired: true, IsOnly: true, Editable: true, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKChildStr, PropertyName: "", IsRequired: false, IsOnly: false, PropertyType: "", Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_module_type", PropertyName: "模块类型", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: moduleTypeOption},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "operator", PropertyName: "主要维护人", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_bak_operator", PropertyName: "备份维护人", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeUser, Option: ""},
	}
	return dataRows
}

// PlatRow plat structure
func PlatRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDPlat
	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKCloudNameField, PropertyName: "云区域", IsRequired: true, IsOnly: true, IsPre: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKOwnerIDField, PropertyName: "供应商", IsRequired: true, IsOnly: true, IsPre: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
	}
	return dataRows
}

// HostRow host structure
func HostRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDHost
	groupAgent := mCommon.HostAutoFields
	dataRows := []*metadata.ObjectAttDes{
		//基本信息分组
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKHostInnerIPField, PropertyName: "内网IP", IsRequired: true, IsOnly: true, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKHostOuterIPField, PropertyName: "外网IP", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "operator", PropertyName: "主要维护人", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_bak_operator", PropertyName: "备份维护人", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_asset_id", PropertyName: "固资编号", IsRequired: false, IsOnly: false, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_sn", PropertyName: "设备SN", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_comment", PropertyName: "备注", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeLongChar, Option: ""},

		//自动发现分组
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKCloudIDField, PropertyName: "云区域", IsRequired: false, IsOnly: true, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleAsst, Option: common.KvMap{"value": "plat", "label": "云区域"}}, //common.FiledTypeInt, Option: common.KvMap{}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_host_name", PropertyName: "主机名称", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_service_term", PropertyName: "质保年限", IsRequired: false, IsOnly: false, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: common.KvMap{"min": "1", "max": "10"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_sla", PropertyName: "SLA级别", IsRequired: false, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "L1", Type: "text"}, {ID: "2", Name: "L2", Type: "text"}, {ID: "3", Name: "L3", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKOSTypeField, PropertyName: "操作系统类型", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "Linux", Type: "text"}, {ID: "2", Name: "Windows", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKOSNameField, PropertyName: "操作系统名称", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_os_version", PropertyName: "操作系统版本", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_os_bit", PropertyName: "操作系统位数", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_cpu", PropertyName: "CPU逻辑核心数", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeInt, Option: common.KvMap{"min": "1", "max": "1000000"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_cpu_mhz", PropertyName: "CPU频率", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeInt, Unit: "Hz", Option: common.KvMap{"min": "1", "max": "100000000"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_cpu_module", PropertyName: "CPU型号", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_mem", PropertyName: "内存容量", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeInt, Unit: "MB", Option: common.KvMap{"min": "1", "max": "100000000"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_disk", PropertyName: "磁盘容量", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeInt, Option: common.KvMap{"min": "1", "max": "100000000"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_mac", PropertyName: "内网MAC地址", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_outer_mac", PropertyName: "外网MAC", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		//agent 没有分组
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.CreateTimeField, PropertyName: "录入时间", IsRequired: false, IsOnly: false, Editable: false, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeTime, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "import_from", PropertyName: "录入方式", IsRequired: false, IsOnly: false, Editable: false, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "excel", Type: "text"}, {ID: "2", Name: "agent", Type: "text"}, {ID: "3", Name: "api", Type: "text"}}},
	}

	return dataRows
}

// ProcRow proc structure
func ProcRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDProc
	groupPort := mCommon.ProcPort
	// groupGsekit := mCommon.Proc_gsekit_base_info
	// groupGsekitManage := mCommon.Proc_gsekit_manage_info
	dataRows := []*metadata.ObjectAttDes{
		//base info
		//&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_process_id", PropertyName: "进程ID", IsSystem: true, IsRequired: true, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: "{}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: true, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: "{}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKProcessNameField, PropertyName: "进程名称", IsRequired: true, IsOnly: true, IsPre: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "description", PropertyName: "进程描述", IsRequired: false, IsOnly: false, IsPre: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},

		//监听端口分组
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bind_ip", PropertyName: "绑定IP", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupPort, PropertyType: common.FiledTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "127.0.0.1", Type: "text"}, {ID: "2", Name: "0.0.0.0", Type: "text"}, {ID: "3", Name: "第一内网IP", Type: "text"}, {ID: "4", Name: "第一外网IP", Type: "text"}}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "port", PropertyName: "端口", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupPort, PropertyType: common.FiledTypeSingleChar, Option: `^((\d+-\d+)|(\d+))(,((\d+)|(\d+-\d+)))*$`, Placeholder: `单个端口：8080 </br>多个连续端口：8080-8089 </br>多个不连续端口：8080-8089,8199`},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "protocol", PropertyName: "协议", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupPort, PropertyType: common.FiledTypeEnum, Option: []validator.EnumVal{{ID: "1", Name: "TCP", Type: "text"}, {ID: "2", Name: "UDP", Type: "text"}}},

		//gsekit 基础信息
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_func_id", PropertyName: "功能ID", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_func_name", PropertyName: "功能名称", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "work_path", PropertyName: "工作路径", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "user", PropertyName: "启动用户", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "proc_num", PropertyName: "启动数量", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeInt, Option: common.KvMap{"min": "1", "max": "1000000"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "priority", PropertyName: "启动优先级", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeInt, Option: common.KvMap{"min": "1", "max": "100"}},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "timeout", PropertyName: "操作超时时长", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeInt, Option: common.KvMap{"min": "1", "max": "1000000"}},

		//gsekit 进程信息
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "start_cmd", PropertyName: "启动命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "stop_cmd", PropertyName: "停止命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "restart_cmd", PropertyName: "重启命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "face_stop_cmd", PropertyName: "强制停止命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "reload_cmd", PropertyName: "进程重载命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "pid_file", PropertyName: "PID文件路径", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "auto_start", PropertyName: "是否自动拉起", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeBool, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "auto_time_gap", PropertyName: "拉起间隔", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeInt, Option: common.KvMap{"min": "1", "max": "1000000"}},
	}
	return dataRows
}
